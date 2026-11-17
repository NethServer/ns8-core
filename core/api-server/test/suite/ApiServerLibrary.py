#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: AGPL-3.0-or-later
#

"""Robot Framework keywords to test the API server as a black box: a
fake agent that consumes tasks from Redis, a Redis MONITOR recorder,
a WebSocket client, JWT forging and TOTP codes."""

import base64
import hashlib
import hmac
import json
import queue
import re
import struct
import threading
import time

import redis
import websocket
from robot.api.deco import keyword, library


def _b64url(data):
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


@library(scope="GLOBAL")
class ApiServerLibrary:

    def __init__(self):
        self.redis_args = {}
        self.agents = {}
        self.monitor_thread = None
        self.monitor_lines = []
        self.monitor_lock = threading.Lock()
        self.ws = None

    #
    # Redis connection used by the library itself
    #

    @keyword
    def set_redis_credentials(self, username, password, host="127.0.0.1", port=6379):
        self.redis_args = dict(host=host, port=int(port), username=username,
                               password=password, decode_responses=True)

    def _redis(self, **kwargs):
        return redis.Redis(**self.redis_args, **kwargs)

    #
    # Fake agent: a Redis client named like the agent, which pops tasks
    # from its queue. The API server accepts tasks only for agents with
    # a recently active Redis connection.
    #

    @keyword
    def start_fake_agent(self, agent_id):
        stop = threading.Event()
        tasks = queue.Queue()
        rdb = self._redis(client_name=agent_id, socket_timeout=5)

        def run():
            while not stop.is_set():
                try:
                    item = rdb.brpop(f"{agent_id}/tasks", timeout=1)
                except redis.RedisError:
                    if stop.is_set():
                        break
                    time.sleep(0.2)
                    continue
                if item:
                    tasks.put(json.loads(item[1]))

        thread = threading.Thread(target=run, daemon=True)
        thread.start()
        self.agents[agent_id] = (stop, thread, tasks, rdb)
        time.sleep(0.5)

    @keyword
    def stop_fake_agents(self):
        for stop, thread, _tasks, rdb in self.agents.values():
            stop.set()
            thread.join(3)
            rdb.close()
        self.agents = {}

    @keyword
    def fake_agent_receives_task(self, agent_id, timeout=5):
        """Return the next task popped by the fake agent"""
        try:
            return self.agents[agent_id][2].get(timeout=float(timeout))
        except queue.Empty:
            raise AssertionError(f"no task received by {agent_id} within {timeout}s")

    @keyword
    def fake_agent_completes_task(self, agent_id, task_id, exit_code=0, output="", error=""):
        """Store the task result and publish its last progress message,
        as the agent does"""
        rdb = self._redis()
        prefix = f"task/{agent_id}/{task_id}"
        with rdb.pipeline() as pipe:
            pipe.set(f"{prefix}/output", output)
            pipe.set(f"{prefix}/error", error)
            pipe.set(f"{prefix}/exit_code", str(exit_code))
            pipe.publish(f"progress/{agent_id}/task/{task_id}", json.dumps({
                "id": task_id,
                "status": "completed" if int(exit_code) == 0 else "aborted",
                "progress": 100,
                "exit_code": int(exit_code),
            }))
            pipe.execute()
        rdb.close()

    #
    # Redis MONITOR recorder
    #

    @keyword
    def start_redis_monitor(self):
        self.stop_redis_monitor()
        with self.monitor_lock:
            self.monitor_lines = []
        rdb = self._redis()
        ready = threading.Event()

        def run():
            with rdb.monitor() as mon:
                ready.set()
                try:
                    for cmd in mon.listen():
                        with self.monitor_lock:
                            self.monitor_lines.append(cmd["command"])
                except (redis.RedisError, OSError, ValueError, AttributeError):
                    pass

        self.monitor_rdb = rdb
        self.monitor_thread = threading.Thread(target=run, daemon=True)
        self.monitor_thread.start()
        ready.wait(3)
        time.sleep(0.2)

    @keyword
    def stop_redis_monitor(self):
        if self.monitor_thread:
            try:
                self.monitor_rdb.connection_pool.disconnect()
            except Exception:
                pass
            self.monitor_thread.join(2)
            self.monitor_thread = None

    @keyword
    def redis_monitor_should_contain(self, pattern, timeout=3):
        """Wait for a command matching the regexp pattern; return it"""
        rx = re.compile(pattern)
        deadline = time.time() + float(timeout)
        while time.time() < deadline:
            with self.monitor_lock:
                for line in self.monitor_lines:
                    if rx.search(line):
                        return line
            time.sleep(0.1)
        with self.monitor_lock:
            seen = "\n".join(self.monitor_lines[-20:])
        raise AssertionError(f"no Redis command matches {pattern!r}; last commands:\n{seen}")

    @keyword
    def redis_monitor_should_not_contain(self, pattern, wait=1):
        time.sleep(float(wait))
        rx = re.compile(pattern)
        with self.monitor_lock:
            found = [line for line in self.monitor_lines if rx.search(line)]
        if found:
            raise AssertionError(f"unexpected Redis command(s) matching {pattern!r}: {found}")

    #
    # WebSocket client
    #

    @keyword
    def open_websocket(self, url="ws://127.0.0.1:9311/ws"):
        self.close_websocket()
        self.ws = websocket.create_connection(url, timeout=5)

    @keyword
    def websocket_send_json(self, message):
        if isinstance(message, str):
            message = json.loads(message)
        self.ws.send(json.dumps(message))

    @keyword
    def websocket_authorize(self, token):
        self.websocket_send_json({"action": "authorize", "payload": {"jwt": token}})

    @keyword
    def websocket_receive_json(self, timeout=3):
        """Return the next JSON message, or None if nothing arrives or
        the server closed the connection"""
        self.ws.settimeout(float(timeout))
        try:
            data = self.ws.recv()
        except (websocket.WebSocketTimeoutException, websocket.WebSocketConnectionClosedException):
            return None
        if not data:
            return None
        return json.loads(data)

    @keyword
    def websocket_should_be_closed(self, timeout=3):
        self.ws.settimeout(float(timeout))
        deadline = time.time() + float(timeout)
        while time.time() < deadline:
            try:
                data = self.ws.recv()
                if not data:
                    return
            except websocket.WebSocketConnectionClosedException:
                return
            except websocket.WebSocketTimeoutException:
                break
        raise AssertionError("the WebSocket connection is still open")

    @keyword
    def close_websocket(self):
        if self.ws:
            try:
                self.ws.close()
            except Exception:
                pass
            self.ws = None

    #
    # JWT and TOTP
    #

    @keyword
    def forge_token(self, secret, claims, alg="HS256"):
        """Sign claims (dict or JSON string) with the API server secret"""
        if isinstance(claims, str):
            claims = json.loads(claims)
        head = _b64url(json.dumps({"alg": alg, "typ": "JWT"}).encode())
        body = _b64url(json.dumps(claims).encode())
        digest = {"HS256": hashlib.sha256, "HS384": hashlib.sha384, "HS512": hashlib.sha512}[alg]
        sig = hmac.new(secret.encode(), f"{head}.{body}".encode(), digest).digest()
        return f"{head}.{body}.{_b64url(sig)}"

    @keyword
    def admin_claims(self, user, expires_in=3600, issued_ago=0, two_fa=False):
        """Claims as the API server issues them"""
        now = int(time.time())
        return {
            "id": user, "role": "", "actions": [], "2fa": bool(two_fa),
            "exp": now + int(expires_in), "orig_iat": now - int(issued_ago),
        }

    @keyword
    def decode_token_claims(self, token):
        payload = token.split(".")[1]
        return json.loads(base64.urlsafe_b64decode(payload + "=" * (-len(payload) % 4)))

    @keyword
    def totp_code(self, secret, offset=0):
        """Current TOTP code (RFC 6238, SHA1, 6 digits, 30s) of a base32 secret"""
        key = base64.b32decode(secret)
        counter = int(time.time()) // 30 + int(offset)
        mac = hmac.new(key, struct.pack(">Q", counter), hashlib.sha1).digest()
        pos = mac[-1] & 0x0F
        code = (struct.unpack(">I", mac[pos:pos + 4])[0] & 0x7FFFFFFF) % 1000000
        return f"{code:06d}"

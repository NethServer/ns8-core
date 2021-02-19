#!/opt/ctla/bin/python3

import redis
import os
import sys
import time

event_paths = ['/opt/ctla/share/events', '/var/opt/ctla/events']
channels = {}

def event_handler(message):
    handlers = {}
    for base in event_paths:
        try:
            path = base + '/' + message['channel']
            with os.scandir(path) as it:
                for entry in it:
                    handlers[entry.name] = f"{path}/{entry.name}"
        except Exception as ex:
            pass

    for h in sorted(handlers):
        try:
            print(f"[INFO] running handler {h}")
            exit_code = os.spawnlp(os.P_WAIT, handlers[h], h, message['data'])
            print("[{0}] handler exit with {1}".format("INFO" if exit_code == 0 else "ERROR", exit_code))
        except Exception as ex:
            print(f"[ERROR] {ex}")

for path in event_paths:
    with os.scandir(path) as it:
        for entry in it:
            if entry.is_dir():
                channels[entry.name] = event_handler

if not channels:
    print("[INFO] nothing to do: no channel handlers found", file=sys.stderr)
    exit(0)

r = redis.Redis(host='127.0.0.1', port=6379, db=0, decode_responses=True)
p = r.pubsub(ignore_subscribe_messages=True)
p.subscribe(**channels)
p.run_in_thread()

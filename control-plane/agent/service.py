#
# Agent
#
# !!!Note!!! pass event dirs as command line arguments
#

import redis
import os
import sys
import time
import socket

channel_prefix = sys.argv[1]
event_paths = sys.argv[2:]

def event_handler(message):
    global event_paths, channel_prefix

    handlers = {}
    for base in event_paths:
        try:
            event_name = (message['channel'])[len(channel_prefix) + 1:]
            path = base + '/' + event_name
            with os.scandir(path) as it:
                for entry in it:
                    handlers[entry.name] = f"{path}/{entry.name}"
        except Exception as ex:
            pass

    errors = 0
    for h in sorted(handlers):
        try:
            print(f"[INFO] running handler {h}")
            exit_code = os.spawnlp(os.P_WAIT, handlers[h], h, message['data'])
            if exit_code == 0:
                print(f"[INFO] handler {h} success")
            else:
                print(f"[ERROR] handler {h} exit code {exit_code}")
                errors += 1
        except Exception as ex:
            print(f"[ERROR] {ex}")

    return errors

def get_channels():
    global event_paths, channel_prefix

    channels = {}
    for path in filter(os.path.isdir, event_paths):
        with os.scandir(path) as it:
            for entry in it:
                if entry.is_dir():
                    channel_name = channel_prefix + ':' + entry.name
                    print(f"[INFO] configured event {entry.name}")
                    channels[channel_name] = event_handler

    return channels

def serve(channels):
    r = redis.Redis(host='127.0.0.1', port=6379, db=0, decode_responses=True)
    p = r.pubsub(ignore_subscribe_messages=True)
    p.subscribe(**channels)
    return p.run_in_thread(sleep_time=0.5)



if __name__ == '__main__':

    sa = os.getenv("NOTIFY_SOCKET")
    if sa:
        so = socket.socket(family=socket.AF_UNIX, type=socket.SOCK_DGRAM)
        def sd_send(msg):
            so.sendto(msg.encode(), sa)
    else:
        def sd_send(msg):
            print(f"[INFO] {msg}")

    startup_event = os.environ.get('AGENT_STARTUP_EVENT', False)
    if startup_event:
        sd_send(f"STATUS=Waiting for startup event {startup_event} to complete...\n")
        errors = event_handler({"channel": channel_prefix + ':' + startup_event, "data": ""})
        if errors > 0:
            sd_send(f"STATUS=Startup event {startup_event} completed with {errors} errors\n")
            exit(2)
        else:
            sd_send(f"STATUS=Startup event {startup_event} completed successfully\n")

    channels = get_channels()
    if not channels:
        print("[ERROR] nothing to do: no channel handlers found", file=sys.stderr)
        exit(3)

    t = serve(channels)
    sd_send("READY=1\n")

    t.join()

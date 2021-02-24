#
# Agent
#
# !!!Note!!! pass event dirs as command line arguments
#

import redis
import os
import sys
import time

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

    for h in sorted(handlers):
        try:
            print(f"[INFO] running handler {h}")
            exit_code = os.spawnlp(os.P_WAIT, handlers[h], h, message['data'])
            print("[{0}] handler exit with {1}".format("INFO" if exit_code == 0 else "ERROR", exit_code))
        except Exception as ex:
            print(f"[ERROR] {ex}")

def serve():
    global event_paths, channel_prefix

    channels = {}
    for path in filter(os.path.isdir, event_paths):
        with os.scandir(path) as it:
            for entry in it:
                if entry.is_dir():
                    channel_name = channel_prefix + ':' + entry.name
                    print(f"[INFO] configured event {entry.name}")
                    channels[channel_name] = event_handler

    if not channels:
        print("[ERROR] nothing to do: no channel handlers found", file=sys.stderr)
        exit(1)

    r = redis.Redis(host='127.0.0.1', port=6379, db=0, decode_responses=True)
    p = r.pubsub(ignore_subscribe_messages=True)
    p.subscribe(**channels)
    p.run_in_thread(sleep_time=0.5).join()


if __name__ == '__main__':
    serve()

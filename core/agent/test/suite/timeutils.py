import time

def get_timestamp_ms():
    return int(time.perf_counter_ns() // 1_000_000)

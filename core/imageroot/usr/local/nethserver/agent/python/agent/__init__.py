#
# Copyright (C) 2021 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# NethServer is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with NethServer.  If not, see COPYING.
#

import redis
import os
import subprocess
import sys
import uuid
import time
import json

__all__ = [
    "redis_connect",
    "run_helper",
    "set_env",
    "slurp_file",
]

def redis_connect(privileged=False, decode_responses=True, **kwargs):
    """Connect to the Redis DB with the right credentials
    """
    redis_host = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
    redis_port = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]
    if privileged:
        redis_username = os.getenv('REDIS_USER', os.getenv('AGENT_ID'))
        assert redis_username
        redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    else:
        redis_username = 'default'
        redis_password = 'nopass'

    return redis.Redis(
        host=redis_host,
        port=redis_port,
        db=0,
        username=redis_username,
        password=redis_password,
        decode_responses=decode_responses,
            #  we assume Redis keys and value strings are encoded UTF-8. Enabling this
            #  option implicitly converts to UTF-8 strings instead of binary strings 
            #  (e.g. {b'key': b'value'} != {'key':'value'})
        **kwargs
    )

def run_helper(*args):
    """Run the command and assert the exit code is 0

    If the exit code is non-zero raise an assertion error.
    The command output is redirected to stderr.
    """
    proc = subprocess.run(args, stdout=sys.stderr)
    assert proc.returncode == 0
    return proc

def set_env(name, value):
    fd = int(os.environ['AGENT_COMFD'])
    os.write(fd, f'set-env {name} {value}\n'.encode('utf-8'))

def slurp_file(file_name):
    with open(file_name) as f:
        return f.read().strip()

def run_subtask(redis_obj, agent_prefix, action, input_string="", input_obj=None):
    if input_obj is not None:
        input_string = json.dumps(input_obj)

    task_id = str(uuid.uuid4())
    task_obj = {"id": task_id, "action": action, "data": input_string}

    redis_obj.lpush(f'{agent_prefix}/tasks', json.dumps(task_obj))

    exit_code = None
    while True: # XXX infinite loop no timeout!
        exit_code = redis_obj.get(f'{agent_prefix}/task/{task_id}/exit_code')
        if exit_code is not None:
            break
        time.sleep(1)

    output = redis_obj.get(f'{agent_prefix}/task/{task_id}/output')
    error = redis_obj.get(f'{agent_prefix}/task/{task_id}/error')
    return int(exit_code), output, error

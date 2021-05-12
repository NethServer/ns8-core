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
import tempfile
import ipcalc
from envparse import env

def redis_connect(privileged=False, decode_responses=True, **kwargs):
    """Connect to the Redis DB with the right credentials
    """
    redis_host = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
    redis_port = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]
    if privileged:
        redis_username = os.getenv('REDIS_USER', os.getenv('AGENT_ID'))
        if not redis_username:
            print("redis_connect: REDIS_USER and AGENT_ID are not set in the environment!", file=sys.stderr)
            # Try to parse the node agent environment as fallback:
            env.read_envfile('/var/lib/nethserver/node/state/agent.env')
            redis_username = env('AGENT_ID', default='default')
            redis_password = env('REDIS_PASSWORD', default='nopass')
        else:
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

    The command output is redirected to stderr.
    """
    return subprocess.run(args, stdout=sys.stderr)

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

def save_wgconf(ipaddr, listen_port=55820, peers={}):

    private_key = slurp_file('/etc/nethserver/wg0.key')
    public_key = slurp_file('/etc/nethserver/wg0.pub')

    oldmask = os.umask(0o77)
    # Create a new file beside our target file path:
    wgconf = open('/etc/wireguard/wg0.conf.new', 'w')
    os.umask(oldmask)

    # Write the Interface head section:
    wgconf.write(f"[Interface]\n")
    wgconf.write(f"Address = {ipaddr}\n")
    wgconf.write(f"ListenPort = {listen_port}\n")
    wgconf.write(f"PrivateKey = {private_key}\n\n")

    # Append Peer sections:
    for pkey, peer in peers.items():
        if peer['public_key'] == public_key:
            continue # Skip record if it refers to the local node

        allowed_ips = { peer['ip_address'] }
        if 'routes' in peer:
            # The set avoids duplicate values:
            allowed_ips.update(peer['routes'])

        wgconf.write(f'[Peer]\n')
        wgconf.write(f"PublicKey = {peer['public_key']}\n")
        wgconf.write(f'AllowedIPs = {", ".join(allowed_ips)}\n')
        wgconf.write(f'PersistentKeepalive = 25\n')
        if 'endpoint' in peer and peer['endpoint'] != '':
            wgconf.write(f"Endpoint = {peer['endpoint']}\n")

    wgconf.close()
    # Overwrite the target file path:
    os.rename('/etc/wireguard/wg0.conf.new', '/etc/wireguard/wg0.conf')

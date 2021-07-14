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
import csv
import traceback
import shlex
import warnings

# Reference https://www.man7.org/linux/man-pages/man3/sd-daemon.3.html
SD_EMERG   = "<0>"  # system is unusable
SD_ALERT   = "<1>"  # action must be taken immediately
SD_CRIT    = "<2>"  # critical conditions
SD_ERR     = "<3>"  # error conditions
SD_WARNING = "<4>"  # warning conditions
SD_NOTICE  = "<5>"  # normal but significant condition
SD_INFO    = "<6>"  # informational
SD_DEBUG   = "<7>"  # debug-level messages

def redis_connect(privileged=False, decode_responses=True, **kwargs):
    """Connect to the Redis DB with the right credentials
    """
    redis_host = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
    redis_port = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]
    if privileged:
        redis_username = os.getenv('REDIS_USER', os.getenv('AGENT_ID'))
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

def read_envfile(file_path):
    """Read an environment file (e.g. agent.env) and return a dictionary of its contents
    """
    fo = open(file_path, 'r')
    lineno = 0
    env = {}
    for line in fo.readlines():
        lineno =+ 1
        try:
            variable, value = line.split("=")
        except ValueError:
            warnings.warn(f'read_envfile: Cannot parse line {lineno} in {file_path}', stacklevel=2)
            continue

        env[variable] = ''.join(list(shlex.shlex(value, posix=True)))

    return env

def run_helper(*args, log_command=True, **kwargs):
    """Run the command in args, writing the command line to stderr.

    The command output is redirected to stderr.
    """
    if log_command:
        print(shlex.join(args), file=sys.stderr)

    return subprocess.run(args, stdout=sys.stderr)

def __action(*args):
    # write to stderr if AGENT file descriptor is not available:
    # this is usefull when developing new actions
    fd = int(os.getenv('AGENT_COMFD',2))
    # The below fdopen does the following:
    # 1. it converts the file descriptor to an object for csv writer library
    # 2. it does not automatically close fd, otherwise next calls to __action will fail
    # 3. it makes sure everything is converted to utf-8 before writing
    fdobj = os.fdopen(fd, mode="w", encoding='utf-8', closefd=False)
    writer = csv.writer(fdobj, delimiter=' ', quotechar='"', quoting=csv.QUOTE_MINIMAL)
    writer.writerow(args)

def set_env(name, value):
    __action("set-env", name, value)

def unset_env(name):
    __action("unset-env", name)

def dump_env():
    __action("dump-env")

def set_status(value):
    __action("set-status", value)

def set_progress(value):
    __action("set-progress", value)

def set_weight(step_name, weight):
    __action("set-weight", step_name, str(weight))

def slurp_file(file_name):
    with open(file_name) as f:
        return f.read().strip()

def resolve_agent_id(agent_selector, node_id=None):
    """Resolves agent_selector to an agent_id
    If node_id is None (default), the NODE_ID value from the environment is considered.
    """
    rdb = redis_connect(privileged=False)
    agent_id = None

    if node_id is None:
        node_id = os.environ.get('NODE_ID')

    if agent_selector == 'node': # Expand to the agent id of the current node
        if node_id:
            agent_id = f'node/{node_id}'
    elif agent_selector.endswith('@node'): # Default module instance on the current node
        if node_id:
            default_instance = rdb.get(f'node/{node_id}/default_instance/{agent_selector[0:-len("@node")]}')
            if default_instance:
                agent_id = f'module/{default_instance}'
    elif agent_selector.endswith('@cluster'): # Default module instance for the whole cluster
        default_instance = rdb.get(f'cluster/default_instance/{agent_selector[0:-len("@cluster")]}')
        if default_instance:
            agent_id = f'module/{default_instance}'
    elif agent_selector.isnumeric(): # Convert to a node ID
        agent_id = f'node/{agent_selector}'
    elif agent_selector == 'cluster' \
        or agent_selector.startswith('module/') \
        or agent_selector.startswith('node/'): # Return value as-is
        agent_id = agent_selector
    else: # Consider the selector as a module ID. Add the "module/" prefix
        agent_id = f'module/{agent_selector}'

    if agent_id is None:
        warnings.warn(f"Could not resolve {agent_selector} to an agent_id", stacklevel=2)

    return agent_id

def run_subtask(redis_obj, agent_prefix, action, input_obj, nowait=False, progress_range=None):

    task_id = str(uuid.uuid4())
    task_obj = {"id": task_id, "action": action, "data": input_obj, "parent": os.getenv("AGENT_TASK_ID", "")}

    redis_obj.lpush(f'{agent_prefix}/tasks', json.dumps(task_obj))

    if nowait:
        return None, None, None

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

def assert_exp(exp, message='Assertion failed'):
    """Like the Python assert statement, this function asserts "exp" evaluates to True.
    If the assertion fails, the program is aborted and a stack trace is printed to stderr.
    """
    if exp:
        return
    else:
        try:
            raise Exception(message)
        except:
            if message: print(message, file=sys.stderr)
            traceback.print_stack(sys.exc_info()[2].tb_frame.f_back, file=sys.stderr)
            sys.exit(2)

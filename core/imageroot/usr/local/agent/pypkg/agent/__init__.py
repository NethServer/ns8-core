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
import ipaddress
import csv
import traceback
import shlex
import warnings
import agent.tasks
import socket
import errno
import urllib.request
import urllib.error

# Reference https://www.man7.org/linux/man-pages/man3/sd-daemon.3.html
SD_EMERG   = "<0>"  # system is unusable
SD_ALERT   = "<1>"  # action must be taken immediately
SD_CRIT    = "<2>"  # critical conditions
SD_ERR     = "<3>"  # error conditions
SD_WARNING = "<4>"  # warning conditions
SD_NOTICE  = "<5>"  # normal but significant condition
SD_INFO    = "<6>"  # informational
SD_DEBUG   = "<7>"  # debug-level messages

def redis_connect(privileged=False, use_replica=False, **kwargs):
    """Connect to the Redis DB. If no arguments are given
    the leader Redis instance with default read-only access rights
    credentials is selected.
    - Set `privileged=True` to modify Redis DB. Replica cannot be modified.
    - Set `use_replica=True` to discover service startup configuration from
      the local Redis replica.

    Any other keyword argument is passed to redis.Redis() constructor.
    """
    if use_replica:
        redis_host = os.getenv('REDIS_REPLICA_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
        redis_port = os.getenv('REDIS_REPLICA_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]
    else:
        redis_host = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
        redis_port = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]

    if privileged:
        redis_username = os.environ['REDIS_USER'] # Fatal if missing!
        redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    else:
        redis_username = 'default'
        redis_password = 'nopass'

    kwargs.setdefault('host', redis_host)
    kwargs.setdefault('port', redis_port)
    kwargs.setdefault('db', 0)
    kwargs.setdefault('username', redis_username)
    kwargs.setdefault('password', redis_password)
    #  we assume Redis keys and value strings are encoded UTF-8. Enabling this
    #  option implicitly converts to UTF-8 strings instead of binary strings
    #  (e.g. {b'key': b'value'} != {'key':'value'})
    kwargs.setdefault('decode_responses', True)

    return redis.Redis(**kwargs)

def read_envfile(file_path):
    """Read an environment file (e.g. agent.env) and return a dictionary of its contents
    """
    fo = open(file_path, 'r')
    lineno = 0
    env = {}
    for line in fo.readlines():
        lineno += 1
        try:
            record = line.rstrip("\n")
            if record == '' or record[0] == '#':
                continue # skip empty lines and comments
            variable, value = record.split("=", 1)
        except ValueError:
            warnings.warn(f'read_envfile: Cannot parse line {lineno} in {file_path}', stacklevel=2)
            continue

        env[variable] = value

    return env

def get_progress_callback(range_low, range_high):
    """Return a function that maps progress range 0-100 to range_low-range_high and
    calls internally the set_progress() function.
    """
    assert_exp(range_low < range_high)
    delta = range_high - range_low
    def cbk(progress):
        if progress >= 0 and progress <= 100:
            set_progress(range_low + int(delta * progress / 100))
    return cbk

def run_helper(*args, log_command=True, **kwargs):
    """Run the command in args, writing the command line to stderr with prio DEBUG.
    Set log_command=False, to suppress the command logging.

    The command output is redirected to stderr by default.
    Additional keyword arguments are passed to subprocess.Popen
    """
    if log_command:
        print(SD_DEBUG + shlex.join(args), file=sys.stderr)

    progress_callback = kwargs.pop('progress_callback', None)
    start_env = kwargs.pop('env', os.environ)
    pass_fds = kwargs.pop('pass_fds', [])
    outstream = kwargs.pop('stdout', sys.stderr)

    curr_progress = -1
    fdr, fdw = os.pipe()
    child_env = dict(start_env, AGENT_COMFD=str(fdw))
    proc = subprocess.Popen(args, stdout=outstream, pass_fds=(*pass_fds, fdw), env=child_env, **kwargs)
    os.close(fdw)
    with os.fdopen(fdr, mode="r", encoding='utf-8') as fdh:
        rows = csv.reader(fdh, delimiter=' ', quotechar='"', quoting=csv.QUOTE_MINIMAL)
        for row in rows:
            if row[0] == 'set-progress' and progress_callback and int(row[1]) > curr_progress:
                curr_progress = int(row[1])
                progress_callback(curr_progress)
    proc.wait()

    if progress_callback and curr_progress < 100:
        progress_callback(100)

    return subprocess.CompletedProcess(args, proc.returncode)

def run_restic(rdb, repository, repo_path, podman_args, restic_args, **kwargs):
    core_env = read_envfile('/etc/nethserver/core.env') # Import URLs of core images
    orepo = rdb.hgetall(f"cluster/backup_repository/{repository}")
    assert_exp(len(orepo) > 0) # Check the repository exists

    # Build the environment to run Restic against the given repository+repo_path
    restic_env = {}
    restic_env["RESTIC_PASSWORD"] = orepo['password']
    restic_env["RESTIC_CACHE_DIR"] = '/var/cache/restic'

    uschema, upath = orepo['url'].split(':', 1)
    if uschema == 's3':
        restic_env["RESTIC_REPOSITORY"] = orepo['url'] + "/" + repo_path
        restic_env["AWS_ACCESS_KEY_ID"] = orepo['aws_access_key_id']
        restic_env["AWS_SECRET_ACCESS_KEY"] = orepo['aws_secret_access_key']
        if orepo['provider'] == 'aws':
            restic_env['AWS_DEFAULT_REGION'] = orepo.get('aws_default_region', '')
    elif uschema == 'b2':
        restic_env["RESTIC_REPOSITORY"] = orepo['url'] + ":" + repo_path
        restic_env["B2_ACCOUNT_ID"] = orepo['b2_account_id']
        restic_env["B2_ACCOUNT_KEY"] = orepo['b2_account_key']
    else:
        raise Exception(f"Schema {uschema} not supported")

    # Build the Podman command line to run Restic
    podman_cmd = ['podman', 'run', '-i', '--rm', '--privileged', '--network=host', '--volume=restic-cache:/var/cache/restic']

    for envvar in restic_env:
        podman_cmd.extend(['-e', envvar]) # Import Restic environment variables

    podman_cmd.extend(podman_args) # Any argument is appended to podman invocation
    podman_cmd.append(core_env["RESTIC_IMAGE"])
    podman_cmd.extend(restic_args)

    penv = os.environ.copy()
    penv.update(restic_env)
    if os.getenv('DEBUG', False):
        print(*([f"{k}={v}" for k,v in restic_env.items()] + podman_cmd), file=sys.stderr)
    else:
        print("restic", *restic_args, file=sys.stderr)

    kwargs.setdefault('encoding', 'utf-8')
    kwargs.setdefault('env', penv)
    kwargs.setdefault('stdout', sys.stdout)
    kwargs.setdefault('stderr', sys.stderr)

    return subprocess.run(podman_cmd, **kwargs)

def get_existing_volume_args():
    """Return a list of --volume arguments for Podman run and similar. The argument values
    are built from the list of available volumes of the module.
    """
    volume_args = {}
    module_is_rootfull = os.geteuid() == 0
    module_id = os.environ['MODULE_ID']
    volume_args['_statedir'] = f"--volume={os.environ['AGENT_STATE_DIR']}:/srv/state"
    with subprocess.Popen(['podman', 'volume', 'ls', '--format=json'], stdout=subprocess.PIPE) as vproc:
        volumes = [vobj['Name'] for vobj in json.load(vproc.stdout)]

    for vname in volumes:
        if module_is_rootfull:
            if not vname.startswith(module_id + '-'):
                continue # skip volume if the prefix is not MODULE_ID
            vdir = vname.removeprefix(module_id + '-')
        else:
            vdir = vname # 1-1 mapping for rootless modules
        volume_args[vname] = f"--volume={vname}:/srv/volumes/{vdir}"

    return list(volume_args.values())

def get_state_volume_args():
    """Return a list of --volume arguments for Podman run and similar. The argument values
    are built from the module etc/state-include.conf file, if available, otherwise return
    the list of already existing volumes.
    """
    volume_args = {}
    install_dir = os.environ['AGENT_INSTALL_DIR']
    state_dir = os.environ['AGENT_STATE_DIR']
    module_is_rootfull = os.geteuid() == 0
    module_id = os.environ['MODULE_ID']
    volume_args['_statedir'] = f"--volume={state_dir}:/srv/state"
    try:
        # Mount Podman volumes by finding their names in the include file:
        with open(install_dir + '/etc/state-include.conf', 'r') as incfile:
            for iline in incfile:
                if iline.startswith('volumes/'):
                    sstart = len('volumes/')
                    send = iline.find('/', sstart)
                    svolume = iline[sstart:send]
                    if not svolume:
                        continue
                    elif module_is_rootfull:
                        dvolume = module_id + '-' + svolume
                    else:
                        dvolume = svolume

                    volume_args[dvolume] = f"--volume={dvolume}:/srv/volumes/{svolume}"
    except FileNotFoundError:
        pass

    return list(volume_args.values())


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
    print(SD_DEBUG + "dump_env() is deprecated. The environment is now automatically persisted to the agent state/ directory at the end of each action step. See also NethServer/core#324.", file=sys.stderr)

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
    elif agent_selector == 'self':
        agent_id = os.environ['AGENT_ID']
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

        allowed_ips = { peer['ip_address'] + '/32' }
        if 'destinations' in peer:
            # The set avoids duplicate values:
            allowed_ips.update({*peer['destinations']})

        wgconf.write(f'[Peer]\n')
        wgconf.write(f"PublicKey = {peer['public_key']}\n")
        wgconf.write(f'AllowedIPs = {", ".join(sorted(allowed_ips))}\n')
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

def save_acls(rdb):
    """
    Copy current ACLs to cluster/acls key
    This function can be executed only on the leader node
    """

    to_skip = ["default", "cluster", "api-server"]

    acl_list = rdb.acl_list()

    trx = rdb.pipeline()

    # Cleanup ACLs
    trx.delete("cluster/acls")

    # Skip ACLs which should be different on each node
    for acl in acl_list:
        user = acl.split(" ",2)[1]
        if user in to_skip:
            continue
        trx.sadd("cluster/acls", acl)

    trx.publish('cluster/event/acl-changed', 'cluster/acls')
    trx.execute()

def get_image_name_from_url(image_url):
    """Return the image name from its URL. e.g.:
    ghcr.io/nethserver/samba:v1.0.0 => samba
    ghcr.io/nethserver/dokuwiki@sha256:.... => dokuwiki
    """
    # Strip the leading url prefix
    _, image_nametag = image_url.rsplit('/', 1)
    # Strip the trailing tag or hash
    image_name, _ = image_nametag.replace('@', ':', 1).split(':', 1)
    return image_name

def add_public_service(name, ports):
    node_id = os.environ['NODE_ID']
    response = agent.tasks.run(
        agent_id=f'node/{node_id}',
        action='add-public-service',
        data={
            'service': name,
            'ports': ports
        }
    )
    return response['exit_code'] == 0

def add_custom_zone(name, interface, ports=[], rules=[]):
    node_id = os.environ['NODE_ID']
    response = agent.tasks.run(
        agent_id=f'node/{node_id}',
        action='add-custom-zone',
        data={
            'zone': name,
            'interface': interface,
            'ports': ports,
            'rules': rules
        }
    )
    return response['exit_code'] == 0


def remove_public_service(name):
    node_id = os.environ['NODE_ID']
    response = agent.tasks.run(
        agent_id=f'node/{node_id}',
        action='remove-public-service',
        data={'service': name}
    )
    return response['exit_code'] == 0

def remove_custom_zone(name):
    node_id = os.environ['NODE_ID']
    response = agent.tasks.run(
        agent_id=f'node/{node_id}',
        action='remove-custom-zone',
        data={'zone': name}
    )
    return response['exit_code'] == 0


def list_service_providers(rdb, service, transport='*', filters={}):
    """Look up the endpoint information about a given service. Filter
    results by transport protocol or any other exact attribute value"""

    results = []

    def match_filter(record, clauses):
        for xcl in clauses:
            if clauses[xcl] != record.get(xcl):
                return False
        else:
            return True

    for rkey in rdb.scan_iter(f'*/srv/{transport}/{service}'):
        rvalue = rdb.hgetall(rkey)

        if rkey.startswith('module/'):
            module_id = rkey[7:rkey.index('/', 7)]

            if not 'module_uuid' in rvalue:
                rvalue['module_uuid'] = rdb.hget(f'module/{module_id}/environment', 'MODULE_UUID')

            if not 'node' in rvalue:
                rvalue['node'] = rdb.hget(f'module/{module_id}/environment', 'NODE_ID')

            rvalue['transport'] = rkey.split('/')[-2]
            rvalue['module_id'] = module_id
            rvalue['ui_name'] = rdb.get(f'module/{module_id}/ui_name')

        if not match_filter(rvalue, filters):
            continue

        results.append(rvalue)

    return results

def get_smarthost_settings(rdb):
    conf = rdb.hgetall('cluster/smarthost')

    if not conf:
        data={
            "port": 587,
            "host": "",
            "username":"",
            "password": "",
            "enabled": False,
            "encrypt_smtp": "starttls",
            "tls_verify": True
        }
    else:
        data={
            "port": int(conf['port']),
            "host": conf['host'],
            "username": conf['username'],
            "password": conf['password'],
            "enabled": conf['enabled'] == "1",
            "encrypt_smtp": conf['encrypt_smtp'],
            "tls_verify": conf['tls_verify'] == "1"
        }
    return data

def http_route_in_use(domain=None, path=''):
    """we search if the domain (foo.com) or the webpath (/foo) is used or not by a service
    return true if the route is used, false is the route is not used"""
    try:
        req = urllib.request.Request('http://127.0.0.1'+path, method="HEAD")
        if domain:
            req.add_header('Host', domain)
        urllib.request.urlopen(req)
    except urllib.error.HTTPError as e:
        if e.code == 404:
            # Assume no path exists
            return False

    return True

def tcp_port_in_use(port):
    """ we search if the TCP port is used or not by a service
        return true if the post is used, false is the port is not used """
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.bind(('', port))
        sock.close()
        sock = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
        sock.bind(('', port))
        sock.close()
        # port is Free,  no errors
        return  False
    except OSError as ex:
        if ex.errno != errno.EADDRINUSE:
            raise
        # port is in use we have raise the error address in use
        return  True

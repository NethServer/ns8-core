#
# Copyright (C) 2023 Nethesis S.r.l.
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

import sys
import subprocess
import urllib.parse
import json
import uuid
import os
import time
import agent
import configparser
import tempfile
import base64
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

def get_default_backup_repository_name(provider, url, rid=""):
    """Suggest a default name for a backup repository"""
    name = ""
    try:
        if provider == "aws":
            name = "AWS " + url.split('/', 1)[1]
        elif provider == "backblaze":
            name = "Backblaze " + url.split(':', 1)[1]
    except:
        pass

    if not name:
        name = provider + " " + rid[:4]

    return name

def validate_schedule(schedule):
    """Check the given schedule string is good for Systemd timers"""
    return subprocess.run(['/usr/bin/systemd-analyze', 'calendar', schedule], stdout=subprocess.DEVNULL).returncode != 0

def run_rclone(rclone_args, temp_config=None, podman_args=[], **runargs):
    with tempfile.NamedTemporaryFile(suffix='.conf', prefix="tmp", dir='rclone') as tfile:
        if temp_config:
            tfile_name = os.path.basename(tfile.name)
            tfile.write(temp_config.encode('utf-8'))
            tfile.flush()
            os.fchmod(tfile.fileno(), 0o644)
            podman_args += [
                f"--env=RCLONE_CONFIG=/etc/rclone/{tfile_name}",
            ]
        if 'input' in runargs:
            podman_args.append("-i")
        run_command_args = \
            ["podman", "exec"] + podman_args + \
            ["rclone-gateway", "rclone"] + rclone_args
        return subprocess.run(run_command_args, **runargs)

def extract_region_code(hostname, position, default=""):
    """Obtain the region code for S3 backends by looking at the domain
    parts of the endpoint FQDN."""
    try:
        return hostname.split('.')[position]
    except (ValueError, IndexError):
        return default

def extract_rclone_basepath(url):
    """Extract the destination path component from a backup URL"""
    uschema, upath = url.split(':', 1)
    if uschema == 's3':
        if '/' in upath:
            _, upath = upath.split('/', 1)
    elif uschema == 'webdav':
        upath = urllib.parse.urlparse(upath).path
    return upath

def generate_rclone_conf(dest_uuid, url, provider, params):
    """Translate the input arguments in a rclone.conf-compatible string"""
    uschema, upath = url.split(':', 1)
    if uschema == 'b2':
        rclone_conf = (
            f"type = b2\n"
            f"account = {params['b2_account_id']}\n"
            f"key = {params['b2_account_key']}\n"
        )
    elif uschema == 's3':
        rclone_conf = (
            f"type = s3\n"
            f"env_auth = true\n"
            f"access_key_id = {params['aws_access_key_id']}\n"
            f"secret_access_key = {params['aws_secret_access_key']}\n"
        )
        s3_endpoint = ""
        if '/' in upath:
            s3_endpoint, _ = upath.split('/', 1)
            rclone_conf += f"endpoint = {s3_endpoint}\n"
        if provider == 'aws':
            rclone_conf += f"provider = AWS\n"
            rclone_conf += f"region = {params.get('aws_default_region', '')}\n"
        elif provider == 'generic-s3' and 'digitalocean' in s3_endpoint:
            rclone_conf += f"provider = DigitalOcean\n"
        elif provider == 'generic-s3' and 'ovh.net' in s3_endpoint:
            region = params.get("aws_default_region", extract_region_code(s3_endpoint, 1))
            rclone_conf += f"provider = Other\n"
            rclone_conf += f"region = {region}\n"
            rclone_conf += f"location_constraint = {region}\n"
        elif provider == 'generic-s3' and 'wasabi' in s3_endpoint:
            rclone_conf += f"provider = Wasabi\n"
            rclone_conf += f"region = {params.get('aws_default_region', extract_region_code(s3_endpoint, 1))}\n"
        elif provider == 'generic-s3' and 'synology' in s3_endpoint:
            rclone_conf += f"provider = Synology\n"
            rclone_conf += f"region = {params.get('aws_default_region', extract_region_code(s3_endpoint, 0))}\n"
            rclone_conf += f"no_check_bucket = true\n"
        elif provider == 'generic-s3' and 'your-objectstorage.com' in s3_endpoint:
            rclone_conf += f"provider = Other\n"
            rclone_conf += f"region = {params.get('aws_default_region', extract_region_code(s3_endpoint, 0))}\n"
        elif provider == 'generic-s3' and 'idrivee2' in s3_endpoint:
            rclone_conf += f"provider = IDrive\n"
            rclone_conf += f"no_check_bucket = true\n"
        elif provider == 'generic-s3' and 'cubbit.eu' in s3_endpoint:
            rclone_conf += f"provider = Other\n"
            rclone_conf += f"region = eu-west-1\n"
        elif provider == 'generic-s3' and 'scalablestorage.it' in s3_endpoint:
            rclone_conf += f"provider = Other\n"
            rclone_conf += f"no_check_bucket = true\n"
        else:
            rclone_conf += f"provider = Other\n"
    elif uschema == 'azure':
        rclone_conf = (
            f"type = azureblob\n"
            f"account = {params['azure_account_name']}\n"
            f"key = {params['azure_account_key']}\n"
        )
    elif uschema == 'smb':
        obscured_pass = rclone_obscure(params['smb_pass'])
        rclone_conf = (
            f"type = smb\n"
            f"host = {params['smb_host']}\n"
            f"user = {params['smb_user']}\n"
            f"pass = {obscured_pass}\n"
            f"domain = {params['smb_domain']}\n"
        )
    elif uschema == 'webdav':
        ourl = urllib.parse.urlparse(upath)
        rclone_conf = (
            f"type = webdav\n"
            f"url = {ourl.scheme}://{ourl.netloc}\n"
        )
    else:
        raise Exception(f"Schema {uschema} not supported")
    return f"[{dest_uuid}]\n" + rclone_conf

def stable_uuid_v5(data: dict) -> str:
    payload = json.dumps(
        data,
        sort_keys=True,              # stable order of keys
        separators=(",", ":"),       # remove whitespace for stability
        ensure_ascii=False,          # stable UTF-8 bytes for non-ascii
    ).encode("utf-8")
    return str(uuid.uuid5(uuid.NAMESPACE_URL, payload.decode("utf-8")))

def parse_rclone_params(rclone_conf):
    """Parse a rclone remote configuration string and return a flat dictionary.

    The returned dictionary contains the raw rclone keys from the first
    section of rclone_conf, plus aliases that map rclone-specific key names
    back to the original param names expected by generate_rclone_conf().
    This ensures that parse_rclone_params(generate_rclone_conf(...)) is
    a superset of the original params dict.
    """
    cp = configparser.ConfigParser()
    cp.read_string(rclone_conf)
    section = cp.sections()[0]
    result = dict(cp[section])

    rtype = result.get('type', '')
    if rtype == 'b2':
        result['b2_account_id'] = result.get('account', '')
        result['b2_account_key'] = result.pop('key', '')
    elif rtype == 's3':
        result['aws_access_key_id'] = result.get('access_key_id', '')
        result['aws_secret_access_key'] = result.get('secret_access_key', '')
        if 'region' in result:
            result['aws_default_region'] = result['region']
    elif rtype == 'azureblob':
        result['azure_account_name'] = result.get('account', '')
        result['azure_account_key'] = result.pop('key', '')
    elif rtype == 'smb':
        result['smb_host'] = result.get('host', '')
        result['smb_user'] = result.get('user', '')
        result['smb_pass'] = rclone_reveal(result.pop('pass'))
        result['smb_domain'] = result.get('domain', '')

    return result

def lookup_node_from_webdav_url(rdb, url):
    """Look up the node ID from its internal Rclone WebDAV URL, that
    contains its VPN IP address."""
    # Parse a url like "webdav:http://10.5.4.1:4694"
    try:
        node_ip = url.split(":")[2].strip('/')
    except IndexError:
        return None
    for node_id in set(rdb.hvals("cluster/module_node")):
        if node_ip == rdb.hget(f"node/{node_id}/vpn", "ip_address"):
            return node_id
    return None


# Same key hardcoded in rclone source
_CRYPT_KEY = bytes([
    0x9c, 0x93, 0x5b, 0x48, 0x73, 0x0a, 0x55, 0x4d,
    0x6b, 0xfd, 0x7c, 0x63, 0xc8, 0x86, 0xa9, 0x2b,
    0xd3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x8a, 0xfb,
    0xf4, 0xde, 0x16, 0x2b, 0x8b, 0x95, 0xf6, 0x38,
])

_AES_BLOCK_SIZE = 16  # aes.BlockSize in Go


def rclone_obscure(plaintext: str) -> str:
    """Obscure a value using AES-CTR, compatible with rclone's obscure format."""
    data = plaintext.encode()
    iv = os.urandom(_AES_BLOCK_SIZE)
    cipher = Cipher(algorithms.AES(_CRYPT_KEY), modes.CTR(iv))
    encryptor = cipher.encryptor()
    ciphertext = encryptor.update(data) + encryptor.finalize()
    return base64.urlsafe_b64encode(iv + ciphertext).rstrip(b"=").decode()


def rclone_reveal(obscured: str) -> str:
    """Reveal a value obscured by rclone (or obscure() above)."""
    # Restore padding stripped by rclone's RawURLEncoding
    padded = obscured + "=" * (-len(obscured) % 4)
    raw = base64.urlsafe_b64decode(padded)
    if len(raw) < _AES_BLOCK_SIZE:
        raise ValueError("input too short — is it actually obscured?")
    iv, ciphertext = raw[:_AES_BLOCK_SIZE], raw[_AES_BLOCK_SIZE:]
    cipher = Cipher(algorithms.AES(_CRYPT_KEY), modes.CTR(iv))
    decryptor = cipher.decryptor()
    return (decryptor.update(ciphertext) + decryptor.finalize()).decode()

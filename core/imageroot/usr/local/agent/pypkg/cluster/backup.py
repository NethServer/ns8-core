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

run_rclone_count = 0

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
    global run_rclone_count
    run_rclone_count += 1
    container_name = f"rclone{run_rclone_count}run" + str(os.getpid())
    tmp_file = None
    try:
        if temp_config:
            tmp_file = '.' + container_name + '.tmp'
            podman_args += [f"--volume=./{tmp_file}:/etc/rclone/rclone.conf"]
            with open(tmp_file, "w") as fp:
                fp.write(temp_config)
        run_command_args = [
            "podman", "run", "--network=host", "--rm", "--replace",
            "--log-driver=none", "--name=" + container_name,
            '--entrypoint=/usr/local/bin/rclone-wrapper'] + \
            podman_args + \
            [os.environ["RESTIC_IMAGE"]] + rclone_args
        return subprocess.run(run_command_args, **runargs)
    finally:
        if tmp_file:
            os.unlink(tmp_file)

def extract_region_code(hostname, position, default=""):
    """Obtain the region code for S3 backends by looking at the domain
    parts of the endpoint FQDN."""
    try:
        return hostname.split('.')[position]
    except (ValueError, IndexError):
        return default

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
            s3_endpoint, upath = upath.split('/', 1)
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
        else:
            rclone_conf += f"provider = Other\n"
    elif uschema == 'azure':
        rclone_conf = (
            f"type = azureblob\n"
            f"account = {params['azure_account_name']}\n"
            f"key = {params['azure_account_key']}\n"
        )
    elif uschema == 'smb':
        rclone_conf = (
            f"type = smb\n"
            f"host = {params['smb_host']}\n"
            f"user = {params['smb_user']}\n"
            f"pass = {params['smb_pass']}\n"
            f"domain = {params['smb_domain']}\n"
        )
    elif uschema == 'webdav':
        ourl = urllib.parse.urlparse(upath)
        upath = ourl.path
        rclone_conf = (
            f"type = webdav\n"
            f"url = {ourl.scheme}://{ourl.netloc}\n"
        )
    else:
        raise Exception(f"Schema {uschema} not supported")
    return f"[{dest_uuid}]\n" + rclone_conf, upath

def stable_uuid_v5(data: dict) -> str:
    payload = json.dumps(
        data,
        sort_keys=True,              # stable order of keys
        separators=(",", ":"),       # remove whitespace for stability
        ensure_ascii=False,          # stable UTF-8 bytes for non-ascii
    ).encode("utf-8")
    return str(uuid.uuid5(uuid.NAMESPACE_URL, payload.decode("utf-8")))

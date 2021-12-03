#
# Copyright (C) 2021 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redisclienttribute it and/or modify
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

import os
import sys
import agent
import os.path
import tempfile
import subprocess


class Restic:
    config = dict()
    restic_dir = ""
    restic_env = ""

    def prepare_env(self):
        # Create working directories
        os.makedirs(self.restic_dir, exist_ok=True)
        os.makedirs(self.cache_dir, exist_ok=True)

        # Prepare restic environment
        self.restic_env = f'{self.restic_dir}/restic.env'
        with open(self.restic_env, 'w') as efile:
            efile.write(f"RESTIC_PASSWORD={self.config['password']}\n")
            efile.write(f"RESTIC_CACHE_DIR=/cache\n")
            efile.write(f"RESTIC_REPOSITORY={self.config['url']}:{self.directory}\n")

            # Supported backend:
            # - S3 compatible services
            # - SFTP on standard ports
            # - Backblaze B2
            # - Azure

            parts = self.config['url'].split(":")
            if parts[0] == "s3":
                efile.write(f"AWS_ACCESS_KEY_ID={self.config['aws_access_key_id']}\n")
                efile.write(f"AWS_SECRET_ACCESS_KEY={self.config['aws_secret_access_key']}\n")
            elif parts[0] == "b2":
                efile.write(f"B2_ACCOUNT_ID={self.config['b2_account_id']}\n")
                efile.write(f"B2_ACCOUNT_KEY={self.config['b2_account_key']}\n")
            elif parts[0] == "azure":
                efile.write(f"AZURE_ACCOUNT_NAME={self.config['azure_account_name']}\n")
                efile.write(f"AZURE_ACCOUNT_KEY={self.config['azure_account_key']}\n")


class Restore(Restic):

    def __init__(self, module, repository):
        self.config['repository'] = repository
        rdb = agent.redis_connect()
        repo_config = rdb.hgetall(f'cluster/backup_repository/{self.config["repository"]}')
        for k in repo_config:
            self.config[k] = repo_config[k]
        rdb.close()
        self.config['module'] = module
        # The name is in the form <module_name>/<module_id>@<module_uuid>
        self.config['name_name'] = os.path.basename(module)
        self.config['module_id'], sep, self.config['module_uuid'] = self.config['name_name'].partition('@')

    def restore(self, source, destination):
        restic_env = self.prepare_env()
        cmd = ["restic", "--cache-dir", cache_dir, "-r", f"{self.config['url']}:{self.config['module']}"]
        with tempfile.TemporaryDirectory() as cache_dir:
            subprocess.run(cmd + ["restore", "latest", "--target", destination, "--include", source], env=restic_env)

    def dump_env(self):
        restic_env = self.prepare_env()
        rootless_env = f'/home/{self.config["module_id"]}/.config/state/environment'
        rootfool_env = f'/var/lib/nethserver/{self.config["module_id"]}/state/environment'
        with tempfile.TemporaryDirectory() as cache_dir:
            cmd = ["restic", "--cache-dir", cache_dir, "-r", f"{self.config['url']}:{self.config['module']}"]
            p_dump = None
            env = dict()
            # First try to access rootless env, if it fails, switch to rootfull one
            try:
                p_dump = subprocess.run(cmd + ["dump", "latest", rootless_env], env=restic_env, check=True, capture_output=True)
            except:
                p_dump = subprocess.run(cmd + ["dump", "latest", rootfool_env], env=restic_env, check=True, capture_output=True)
            for line in p_dump.stdout.splitlines():
                var, sep, val = line.decode().partition("=")
                # Skip vars which will be replaced on install
                if var not in ["NODE_ID", "TCP_PORT", "TCP_PORTS", "IMAGE_ID", "IMAGE_DIGEST", "IMAGE_REOPODIGEST", "MODULE_ID"]:
                    env[var] = val

        return env


class Backup(Restic):

    name = ""
    module_id = ""
    module_uuid = ""
    module_name = ""
    rootfull = False
    volumes = dict()
    paths = []

    def _get_name_from_url(self, url):
        url, sep, tag = url.rpartition(":")
        base_url, sep, name = url.rpartition("/")
        return name

    def add_volume(self, name):
        proc = subprocess.run(["podman", "volume", "inspect", "--format", "{{.Mountpoint}}", name], capture_output=True)
        self.volumes[name] = proc.stdout.decode().strip()

    def add_path(self, path):
        self.paths.append(path)

    def add_key(self, key, name = None):
        rdb = agent.redis_connect(privileged=False, decode_responses=False)
        dump = rdb.dump(key)
        rdb.close()

        if name == None:
            (prefix, sep, name) = key.rpartition('/')

        path = f'{self.dump_dir}/{name}.dump'
        with open(path, 'wb') as f:
            f.write(dump)
        self.paths.append(path)

    def __init__(self, name = None):
        # Try to read the backup name from command line if it's not passed as argument
        if name is None:
            if len(sys.argv) < 2:
                print("Missing backup name", file=sys.stderr)
                sys.exit(1)

            self.name = sys.argv[1]
        else:
            self.name = name

        self.rootfull = (os.geteuid() == 0)
        self.module_id = os.environ['MODULE_ID']

        self.module_name = self._get_name_from_url(os.environ['IMAGE_URL'])

        # UUID will be used to create unique directory inside the repository
        self.module_uuid = os.environ['UUID']
        self.directory = f"{self.module_name}/{self.module_id}@{self.module_uuid}"

        rdb = agent.redis_connect()

        # Check if backup exists
        if not rdb.exists(f"module/{self.module_id}/backup/{self.name}"):
            print(f"Configuration for backup '{self.name}' not found", file=sys.stderr)
            sys.exit(1)
       
        # Retrieve backup configuration
        self.config = rdb.hgetall(f"module/{self.module_id}/backup/{self.name}")
        repo_config = rdb.hgetall(f'cluster/backup_repository/{self.config["repository"]}')
        for k in repo_config:
            self.config[k] = repo_config[k]

        if self.rootfull:
            self.environment =  f"/var/lib/nethserver/{self.module_id}/state/environment"
            self.dump_dir = f"/var/lib/nethserver/{self.module_id}/dump/"
            self.restic_dir = f"/var/lib/nethserver/cache/restic"
            self.cache_dir = f"{self.restic_dir}/{self.directory}"
        else:
            self.environment = f'{os.path.expanduser("~")}/.config/state/environment'
            self.dump_dir = f'{os.path.expanduser("~")}/dump/'
            self.restic_dir = f"{os.path.expanduser('~')}/restic"
            self.cache_dir = f'{self.restic_dir}/{self.directory}'

        # Make sure dump dir exists, it's always included inside the backup
        os.makedirs(self.dump_dir, exist_ok=True)

        # Close unused Redis connection
        rdb.close()


    def run(self, prune = True):
        # Check if there is something to backup
        if not self.paths and not self.volumes:
            print("No file to backup", file=sys.stderr)
            sys.exit(1)

        self.prepare_env()

        # Dump environment key
        self.add_path(self.environment)

        # Prepare base command
        podman_cmd = ["podman", "run", "--rm", "--env-file", self.restic_env, "-v", f"{self.cache_dir}:/cache", "-v", f"{self.dump_dir}:/dump"]
        for volume in self.volumes:
            mount = self.volumes[volume]
            (prefix, sep, suffix) = mount.partition(volume)
            podman_cmd = podman_cmd + ["-v", f"{prefix}{volume}:/{volume}"]

        for path in self.paths:
            podman_cmd = podman_cmd + ["-v", f"{path}:/{os.path.basename(path)}"]

        cmd = podman_cmd + ["docker.io/restic/restic"]

        # Check if repository has been already initialized, if not, just do it
        check_init = subprocess.run(cmd + ["snapshots"], capture_output=True)
        if check_init.returncode > 0:
            init = subprocess.run(cmd + ["init"], check=True, capture_output=True)

        # Run the backup
        sub_cmd = ["backup"]
        for volume in self.volumes:
            sub_cmd.append(f"/{volume}")
        for path in self.paths:
            sub_cmd.append(f"/{os.path.basename(path)}")
        sub_cmd.append("/dump")

        subprocess.run(cmd + sub_cmd)

        if prune and self.config['retention']:
            subprocess.run(cmd + ["forget", "--prune", "--keep-within", self.config['retention']])

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
import subprocess


class Backup:

    name = ""
    module_id = ""
    module_uuid = ""
    module_name = ""
    rootfull = False
    config = None
    volumes = []
    paths = []
    excludes = []

    def _get_name_from_url(self, url):
        url, sep, tag = url.rpartition(":")
        base_url, sep, name = url.rpartition("/")
        return name

    def prepare_env(self):
        # Prepare restic environment
        restic_env = os.environ.copy()
        restic_env['RESTIC_PASSWORD'] = self.config['password']

        # Supported backend:
        # - S3 compatible services
        # - SFTP on standard ports
        # - Backblaze B2
        # - Azure
        parts = self.config['repository'].split(":")
        if parts[0] == "s3":
            restic_env['AWS_ACCESS_KEY_ID'] = self.config['aws_access_key_id']
            restic_env['AWS_SECRET_ACCESS_KEY'] = self.config['aws_secret_access_key']
        elif parts[0] == "b2":
            restic_env['B2_ACCOUNT_ID'] = self.config['b2_account_id']
            restic_env['B2_ACCOUNT_KEY'] = self.config['b2_account_key']
        elif parts[0] == "azure":
            restic_env['AZURE_ACCOUNT_NAME'] = self.config['azure_account_name']
            restic_env['AZURE_ACCOUNT_KEY'] = self.config['azure_account_key']

        return restic_env


    def add_volume(self, name):
        proc = subprocess.run(["podman", "volume", "inspect", "--format", "{{.Mountpoint}}", name], capture_output=True)
        self.volumes.append(proc.stdout.decode().strip())

    def add_path(self, path):
        self.paths.append(path)
    
    def add_exclude(self, exclude):
        self.exclude.append(exclude)

    #def add_key(self, key, name):
    #    self._save_dump(self.rdb.dump(key))
    #    return

    #def _save_dump(dump, name):
    #    os.makedirs(self.dump_dir, exist_ok=True)
    #    path = f'{self.dump_dir}/{name}.dump'
    #    with openp(path, 'wb') as f:
    #        f.write(dump)
    #    self.paths.append(paths)

    def __init__(self):
        """Read arguments from command line and initialize the object. Exit 1 if something is wrong
        """
        if len(sys.argv) < 2:
            print("Missing backup name", file=sys.stderr)
            sys.exit(1)

        self.name = sys.argv[1]

        self.rootfull = (os.getuid() == 0)
        self.module_id = os.environ['MODULE_ID']

        self.module_name = self._get_name_from_url(os.environ['IMAGE_URL'])

        # UUID will be used to create unique directory inside the repository
        self.module_uuid = os.environ['UUID']
        self.directory = f"{self.module_name}@{self.module_uuid}"

        rdb = agent.redis_connect()

        # Check if backup exists
        if not rdb.exists(f"module/{self.module_id}/backup/{self.name}"):
            print(f"Configuration for backup '{self.name}' not found", file=sys.stderr)
            sys.exit(1)
       
        # Retrieve backup configuration
        self.config = rdb.hgetall(f"module/{self.module_id}/backup/{self.name}")

        if self.rootfull:
            self.environment =  f"/var/lib/nethserver/{self.module_id}/state/environment"
            self.dump_dir = f"/var/lib/nethserver/{self.module_id}/dump/"
            self.cache_dir = f"/var/lib/nethserver/cache/restic/{self.directory}"
        else:
            self.environment = f'{os.path.expanduser("~")}/.config/state/environment'
            self.dump_dir = f'{os.path.expanduser("~")}/dump/'
            self.cache_dir = f'{os.path.expanduser("~")}/backup/cache/{self.directory}'

        # Close unused Redis connection
        rdb.close()


    def run(self, prune = True):
        # Check if there is something to backup
        if not self.paths and not self.volumes:
            print("No file to backup", file=sys.stderr)
            sys.exit(1)

        restic_env = self.prepare_env()

        # Dump environment key
        self.add_path(self.environment)

        # Prepare base command
        cmd = ["restic", "--cache-dir", self.cache_dir, "-r", f"{self.config['repository']}/{self.directory}"]

        # Check if repository has been already initialized, if not, just do it
        check_init = subprocess.run(cmd + ["snapshots"], env=restic_env, capture_output=True)
        if check_init.returncode > 0:
            init = subprocess.run(cmd + ["init"], env=restic_env, check=True, capture_output=True)

        # Run the backup
        for exclude in self.excludes:
            cmd = cmd + ["--exclude", exclude]

        subprocess.run(cmd + ["backup"] + self.volumes + self.paths, env=restic_env)

        if prune and self.config['retention']:
            subprocess.run(cmd + ["forget", "--prune", "--keep-within", self.config['retention']], env=restic_env)

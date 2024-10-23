---
layout: default
title: Backup & Restore
nav_order: 10
parent: Core
---

# Backup & Restore

The backup and restore is a procedure for disaster-recovery scenario: it can be used to save the data of an installed
module and restore it to a different node or cluster.

* TOC
{:toc}

## Design

The backup engine is [Restic](https://restic.net/) which runs inside a
container along with [Rclone](https://rclone.org/) used to inspect the
backup repository contents.

Backups are saved inside *backup destinations*, remote spaces where data
are stored. A backup destination can contain multiple backup instances,
each module instance has its own sub-directory to avoid conflicts. This
sub-directory is the root of the module instance Restic repository.

The system implements the common logic for backup inside the agent with `module-backup` command.
Each module can implement `module-dump-state` and `module-cleanup-state` to prepare/cleanup the data that has to be included in the backup.
The `state/environment` file is always included inside the backup, as it is required by the `cluster/restore-module` action.
The restore is implemented using a `restore-module` action inside the module agent, each module can extend it to implement specific restore steps.
The basic `10restore` step actually runs the Restic restore procedure in a temporary container.

All backups are scheduled by systemd timers. Given a backup with id `1`, it is possible to retrieve the time status with:
- rootless containers, eg. `dokuwiki1`, executed by `dokuwiki1` user: `systemctl --user status backup1.timer`
- rootfull containers, eg. `dnsmasq1`, executed by `root` user: `systemctl status backup1-dnsmasq1.timer`

## Include and exclude files

Whenever possible, containers should use volumes to avoid UID/GID
namespace mappings and SELinux issues during backup an restore.

Includes can be added to the `state-include.conf` file saved inside `AGENT_INSTALL_DIR/etc/`.
In the [source tree](modules/images/#source-tree), the file should be placed under `<module>/imageroot/etc/state-include.conf`.
On installed modules, the file will appear on different paths:
- rootless containers, eg. `dokuwiki1`, full path will be `/home/dokuwiki1/.config/etc/state-include.conf`
- rootfull containers, eg. `dnsmasq1`,  full path will be  `/var/lib/nethserver/dnsmasq1/etc/state-include.conf`

Lines are interpreted as path patterns. Only patterns referring to
volumes and the agent `state/` directory are considered.

Lines starting with `state/` refer to `AGENT_STATE_DIR` contents. Eg. to
include `mykey.dump` under the `AGENT_STATE_DIR` add

    state/mykey.dump

Lines starting with `volumes/` will be mapped to a volume name. Eg. to
include the whole `dokuwiki-data` volume add

    volumes/dokuwiki-data

Internally, volumes will be mapped as:

- `<volume_name>` (1-1) for rootless containers; eg. for module
  `dokuwiki1`, line prefix `volumes/dokuwiki-data` maps to volume name
  `dokuwiki-data`

- `<module_id>-<volume_name>` for rootfull containers; eg. for module
  `dnsmasq1`, line prefix `volumes/ data` maps to volume name `dnsmasq1-data`

Volumes listed in `state-include.conf` are automatically mounted (and
created if necessary) by the basic `10restore` step of the
`restore-module` action.

Excludes can be added to `state-exclude.conf` file saved inside the `AGENT_INSTALL_DIR`.

For a complete explanation of the patterns, like wildcard characters, see
the official Restic documentation to
[include](https://restic.readthedocs.io/en/stable/040_backup.html#including-files)
and
[exclude](https://restic.readthedocs.io/en/stable/040_backup.html#excluding-files)
files. Note that include and exclude patterns have a slight different
syntax.

## Save and restore Redis keys

### Dump key and state

To save a Redis key, you should:
- dump the key inside the `module-dump-state` command
- include the dump inside the backup

Given a module named `mymodule`, create the file `mymodule/imageroot/bin/module-dump-state` inside the module source tree:
```
#!/bin/bash
redis-dump module/mymodule1/mykey > mykey.dump
```

Make sure also `module-dump-state` is executable:
```
chmod a+x mymodule/imageroot/bin/module-dump-state
```

Then, add the key dump path to `mymodule/imageroot/etc/state-include.conf`:
```
state/mykey.dump
```

### Cleanup state

As best practice, the dump should be removed when the backup has completed.

Given a module named `mymodule`, create the file `mymodule/imageroot/bin/module-cleanup-state` inside the module source tree:
```
#!/bin/bash
rm -f mykey.dump
```

Make sure also `module-cleanup-state` is executable:
```
chmod a+x mymodule/imageroot/bin/module-cleanup-state
```

### Restore key

To restore a Redis key, you should add a step inside the `restore-module` action, after index 10.

Given a module named `mymodule`, create a file named `mymodule/imageroot/actions/restore-module/20loadkey` inside the module source tree:
```
#!/bin/bash
redis-restore mymodule1/mykey < mykey.dump
```

## Execute a backup

Before executing a backup, you must create a backup repository.
The flow to execute a backup will be something like:
- create a backup repository and retrieve its UUID
- retrieve the UUID of the repository
- configure a backup using the repository as target
- retrieve the backup ID
- execute the backup

1. Create the repository:
```
api-cli run add-backup-repository --data '{"name":"BackBlaze repo1","url":"b2:backupns8","parameters":{"b2_account_id":"xxxxxxxxxxxxxxxxxxxxxxxxx","b2_account_key":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"},"provider":"backblaze","password":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}'
```

2. The output will be something like, please note the `id` field:
```json
{"password": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "id": "48ce000a-79b7-5fe6-8558-177fd70c27b4"}
```

3. Create a new daily backup named `mybackup` with a retention of 3 snapshots (3 days) which includes `dokuwiki1` and `dnsmasq1` instances:
```
api-cli run add-backup --data '{"repository":"48ce000a-79b7-5fe6-8558-177fd70c27b4","schedule":"daily","retention":3,"instances":["dokuwiki1","dnsmasq1"],"enabled":true, "name":"mybackup"}'
```

4. The output will the id of the backup:
```
1
```

5. Run the backup with id `1`:
```
api-cli run run-backup --data '{"id":1}'
```

For debugging purposes, you can also launch systemd units:
- rootless container, eg. `dokuwiki1`: `runagent -m dokuwiki1 systemctl --user start backup1.service`
- rootfull container, eg. `dnsmasq1`: systemctl start backup1-dnsmasq1.service

To remove the backup use:
```
api-cli run remove-backup --data '{"id":1}'
```

To remove the backup repository:
```
api-cli run remove-backup-repository --data '{"id":"c7a9cfea-303c-5104-8ab7-39ac9f9842bd"}'
```

## Execute a restore

Before executing a restore, the backup repository should already be configured.

Restore the `dokuwiki1` instance at node `1` from repository `48ce000a-79b7-5fe6-8558-177fd70c27b4`:
```
api-cli run cluster/restore-module --data '{"node":1, "repository":"48ce000a-79b7-5fe6-8558-177fd70c27b4", "path":"dokuwiki/dokuwiki1@3792c7db-9450-4bd3-84a3-034cd0087839","snapshot":""}'
```

## Cluster configuration backup

The `cluster/download-cluster-backup` API returns a random URL path where an encrypted
archive is available for download.

    curl -O http://127.0.0.1:9311/backup/$(api-cli run cluster/download-cluster-backup | jq -r .path)

If the previous command is successful a file `dump.json.gz.gpg` is created
in the current directory.

The file with `.gpg` extension is encrypted with the password saved inside `/var/lib/nethserver/cluster/state/backup/passphrase`.
To decrypt it run a command like this:

    echo <passphrase> | \
        gpg --batch -d --passphrase-file /dev/stdin --pinentry-mode loopback -o dump.json.gz dump.json.gz.gpg

The restore procedure can be started from the UI of a new NS8
installation: upload the file and specify the password from the UI.

## The `restic-wrapper` command

The Restic binary is not installed in the host system. NS8 runs Restic
within a core container, preparing environment variables with values read
from the Redis DB and properly mounting the application Podman volumes.

The `restic-wrapper` command is designed to manually run Restic from the
command line. It can help to restore individual files and directories, or
run maintenance commands on remote Restic repositories.

The command can be invoked from any agent environment. Print its inline
help with this command:

    runagent restic-wrapper --help

Some options require a module backup ID, or backup destination UUID. Use
the `--show` option to list them. For example:

    runagent -m mail1 restic-wrapper --show

Example of output:

    Destinations:
    - dac5d576-ed63-5c4b-b028-c5e97022b27b OVH S3 destination (s3:s3.de.io.cloud.ovh.net/ns8-backups)
    - 14030a59-a4e6-57cc-b8ea-cd5f97fe44c8 BackBlaze repo1 (b2:ns8-backups)
    Scheduled backups:
    - 1 Backup to BackBlaze repo1, destination UUID 14030a59-a4e6-57cc-b8ea-cd5f97fe44c8

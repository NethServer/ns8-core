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

The backup engine is [restic](https://restic.net/) which runs inside a container along with [rclone](https://rclone.org/)
used to inspect the backup repository contents.

Backups are saved inside *backup repositories*, remote network-accessible spaces where data are stored.
No local storage backup is possible (eg. USB disks).
A backup repository can contain multiple backup instances, each module instance has its own sub-directory to avoid conflicts.

The system implements the common logic for backup inside the agent with `module-backup` command.
Each module can implement `module-dump-state` and `module-cleanup-state` to prepare/cleanup the data that has to be included in the backup.
The `state/environment` file is always included inside the backup, as it is required by the `cluster/restore-module` action.
The restore is implemented using a `restore-module` action inside the module agent, each module can extend it to implement specific restore steps.
The basic `10restore_module` actually invokes the `module-restore` command.

All backups are scheduled by systemd timers. Given a backup with id `1`, it is possible to retrieve the time status with:
- rootless containers, eg. `dokuwiki1`, executed by `dokuwiki1` user: `systemctl --user status backup1.timer`
- rootfull containers, eg. `samba1`, executed by `root` user: `systemctl status mybackup-samba1.timer`

## Include and exclude files

Whenever possible, containers should always use volumes to avoid SELinux issues during backup an restore.

Includes can be added to the `state-include.conf` file saved inside the `AGENT_INSTALL_DIR`.
In the [source tree](modules/images/#source-tree), the file should be placed under `<module>/imageroot/etc/state-include.conf`.
On installed modules, the file will appear on differenet paths:
- rootless containers, eg. `dokuwiki1`, full path will be `/home/dokuwiki1/.config/etc/state-include.conf`
- rootfull containers, eg. `samba1`,  fulle path will be  `/var/lib/nethserver/samba1/etc/stae-include.conf`

Lines starting with `volumes/` will be translated to volume name. Example:
```
volumes/dokuwiki-data
```

Internally, volumes will be mapped as:
- `<volume_name>` for rootless containers, eg. for `dokuwiki1` map to `dokuwiki-data`
- `<module_id>-<volume_name>` for rootfull containers, eg. for `samba1` map to `samba1-data` 

Excludes can be added to `state-exclude.conf` file saved inside the `AGENT_INSTALL_DIR`.

See [include](https://restic.readthedocs.io/en/stable/040_backup.html#including-files) and [exclude](https://restic.readthedocs.io/en/stable/040_backup.html#excluding-files).

## Save and restore Redis keys

### Dump key

To save a Redis key, you should:
- dump the key inside the `module-dump-state` command
- include the dump inside the backup

Given a module named `mymodule`, create the file `mymodule/imageroot/bin/module-dump-state` inside the module source tree:
```
#!/bin/bash
redis-dump mymodule1/mykey > mykey.dump
```

Then, add the key dump path to `mymodule/imageroot/etc/state-include.conf`:
```
state/mykey.dump
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

3. Create a new daily backup named `mybackup` with a retention of 3 snapshots (3 days) which includes `dokuwiki1` and `samba1` instances:
```
api-cli run add-backup --data '{"repository":"48ce000a-79b7-5fe6-8558-177fd70c27b4","schedule":"daily","retention":3,"instances":["dokuwiki1","samba1"],"enabled":true, "name":"mybackup"}'
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
- rootless container, eg. `dokuwiki1`: `ssh dokuwiki1@localhost "systemctl --user start backup1.service"`
- rootfull container, eg. `samba1`: systemctl start backup1-samba1.service

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

Not implemented yet, see the [card](https://trello.com/c/i5aIgxif/143-cluster-backup-restore) on project board.

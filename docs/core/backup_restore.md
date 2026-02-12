---
layout: default
title: Backup and Restore
nav_order: 10
parent: Core
---

# Backup and Restore

NS8 core provides a backup and restore infrastructure for applications,
and a procedure tailored for single node cluster disaster recovery.

* TOC
{:toc}

## Core architecture

- A backup **destination** is a remote filesystem accessible with Rclone,
  where backup data is saved. It's a cluster object persisted in Redis,
  under the restricted `private/nodes/*` key space: only nodes can read
  its contents. Secrets and configurations are stored in Rclone-format in
  the HASH key `private/nodes/backup_destination/rclone_conf`, indexed by
  destination ID. Other HASH keys are `cluster/backup_repository/<destination
  ID>`, with public backup destination attributes, and the HASH key
  `private/agents/backup_destination/restic_password`, indexed by
  destination ID: it contains the encryption password of all Restic
  repositories in the backup destination, shared among the applications
  that use it.

- A **backup** is the relation between a periodic schedule of the
  backup, a backup destination, and a set of applications to backup. It
  also has auxiliary attributes, like the backup snapshot retention. It's
  a cluster object persisted in Redis in HASH keys in the form
  `cluster/backup/<backup ID>`. The backup ID is generated from a
  sequence: `cluster/backup_sequence`.

- Every cluster node manages a pool of Systemd timers with the
  `backup-timers.service` unit, to start the backups at the scheduled
  times. The node `run-backup` helper command iterates over the configured
  applications and executes their `module-backup` procedure, by `runagent`
  impersonation. The `run-backup` command also updates the backup status
  of every application in Redis key `node/<node ID>/backup_status/<backup
  ID>`, and writes an overall backup status into
  `/run/node_exporter/backup<backup ID>.prom` files. It also uploads
  `.json` files with backup repositories metadata. When executed on the
  leader node, uploads a copy of the encrypted cluster backup.

- Every cluster node runs a `rclone-gateway` service, an authenticated
  HTTP proxy server. Access is granted to nodes and cluster applications
  providing their Redis credentials (environment variables `REDIS_USER`,
  `REDIS_PASSWORD`). The HTTP proxy server is a frontend to two internal
  `rclone serve` backends: `webdav` and `rest`. Nodes can access both,
  whilst applications can access a restricted set of operations of the
  Restic Rest server, and only on the data they own. Application backup
  jobs (Restic) access their repository through the REST endpoint. See
  also the [module backup documentation]({{site.baseurl}}/modules/backup_restore).


## Execute a backup

Before executing a backup, you must create a backup repository.
The flow to execute a backup will be something like:
- create a backup repository and retrieve its UUID
- retrieve the UUID of the repository
- configure a backup using the repository as target
- retrieve the backup ID
- execute the backup

1. Create the repository:

       api-cli run add-backup-repository --data '{"name":"BackBlaze repo1","url":"b2:backupns8","parameters":{"b2_account_id":"xxxxxxxxxxxxxxxxxxxxxxxxx","b2_account_key":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"},"provider":"backblaze","password":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}'

2. The output will be something like, please note the `id` field:

       {"password": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "id": "48ce000a-79b7-5fe6-8558-177fd70c27b4"}

3. Create a new daily backup named `mybackup` with a retention of 3
   snapshots (3 days) which includes `dokuwiki1` and `dnsmasq1` instances:

       api-cli run add-backup --data '{"repository":"48ce000a-79b7-5fe6-8558-177fd70c27b4","schedule":"daily","retention":3,"instances":["dokuwiki1","dnsmasq1"],"enabled":true, "name":"mybackup"}'


4. The output will the id of the backup, e.g. `1`

5. Run the backup with id `1`:

       runagent -m node run-backup --backup 1

To remove the backup use:

    api-cli run remove-backup --data '{"id":1}'


To remove the backup repository:

    api-cli run remove-backup-repository --data '{"id":"c7a9cfea-303c-5104-8ab7-39ac9f9842bd"}'

## Execute a restore

Before executing a restore, the backup repository should already be configured.

Restore latest snapshot of the `dokuwiki1` module on node `1` from
repository `48ce000a-79b7-5fe6-8558-177fd70c27b4`:

    api-cli run cluster/restore-module --data '{"node":1, "repository":"48ce000a-79b7-5fe6-8558-177fd70c27b4", "path":"dokuwiki/3792c7db-9450-4bd3-84a3-034cd0087839","snapshot":""}'

## Cluster configuration backup

The `run-backup` command on the leader node automatically uploads an
encrypted copy of the cluster backup to each backup destination,
alongside the application backups.

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

The Restic binary is not provided by the Linux distribution. Instead, NS8
runs Restic within a core container, preparing environment variables with
values read from the Redis DB and properly mounting the application Podman
volumes.

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

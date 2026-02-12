---
layout: default
title: Backup and Restore
nav_order: 60
parent: Modules
---

# Backup and Restore

* TOC
{:toc}

## General concepts

The module backup engine is [Restic](https://restic.net/) which runs
inside a container (see `RESTIC_IMAGE` in `/etc/nethserver/core.env`).

A module backup is saved under a reserved Restic repository inside one or
more *backup destinations*. A module has no direct access to destinations
and their secret configuration: it must access them through the local node
`rclone-gateway` HTTP proxy service, providing the module's Redis
credentials. The URL to access its Restic repository has this form:

    http://127.0.0.1:4694/<destination ID>/<app type>/<module UUID>

The URL path has three levels:

1. The destination ID, in UUID syntax.
2. The application type, the module ID without the trailing number, e.g. `traefik`.
3. The module UUID

## Entry points

The `module-backup` command is provided by NS8 core and is executed with
module's privileges to backup the module.

The module itself can implement two commands: `module-dump-state` and
`module-cleanup-state` to prepare/cleanup the data that has to be included
in the backup.

The module `state/environment` file is always included inside the backup,
as it is required by the cluster during restoration
(`cluster/restore-module` action).

The restore is implemented using the module `restore-module` action. The
core provides a basic implementation, and the module can extend it to
implement specific restore steps. The basic `10restore` step actually runs
the Restic restore procedure in a temporary container and extracts the
backup data to the `state/` directory and to Podman named volumes.

## Include and exclude files

Whenever possible, containers should use volumes to avoid UID/GID
namespace mappings and SELinux issues during backup an restore.

Includes can be added to the `state-include.conf` file saved inside `AGENT_INSTALL_DIR/etc/`.

In the [source tree](modules/images/#source-tree), the file should be
placed under `<module>/imageroot/etc/state-include.conf`. On installed
modules, the file will appear on different paths:

- rootless containers, eg. `dokuwiki1`, full path will be
  `/home/dokuwiki1/.config/etc/state-include.conf`

- rootfull containers, eg. `dnsmasq1`,  full path will be
  `/var/lib/nethserver/dnsmasq1/etc/state-include.conf`

Lines are interpreted as path patterns. Only patterns referring to volumes
and the agent `state/` directory are considered.

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

## Save Redis keys

To save a Redis key, you should:

- dump the key inside the `module-dump-state` command
- include the dump inside the backup

Given a module named `mymodule`, create the file
`mymodule/imageroot/bin/module-dump-state` inside the module source tree:

    #!/bin/bash
    redis-dump module/mymodule1/mykey > mykey.dump


Make sure also `module-dump-state` is executable:

    chmod a+x mymodule/imageroot/bin/module-dump-state


Then, add the key dump path to `mymodule/imageroot/etc/state-include.conf`:

    state/mykey.dump


## Cleanup state

As best practice, the dump should be removed when the backup has completed.

Given a module named `mymodule`, create the file `mymodule/imageroot/bin/module-cleanup-state` inside the module source tree:

    #!/bin/bash
    rm -f mykey.dump

Make sure also `module-cleanup-state` is executable:

    chmod a+x mymodule/imageroot/bin/module-cleanup-state


## Restore Redis keys

To restore a Redis key, you should add a step inside the `restore-module` action, after index 10.

Given a module named `mymodule`, create a file named `mymodule/imageroot/actions/restore-module/20loadkey` inside the module source tree:

    #!/bin/bash
    redis-restore mymodule1/mykey < mykey.dump

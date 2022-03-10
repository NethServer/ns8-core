---
layout: default
title: Cloning a module
nav_order: 10
parent: Core
---

# Cloning a module

A (source) module instance can be cloned (copied) by the
`cluster/clone-module` action. This action creates a new module instance
(destination) that is equivalent to the source one.

For example to create a clone/copy of instance `dokuwiki1` on node 1, run

    api-cli run clone-module --data '{"module":"dokuwiki1","replace":false,"node":1}'

The cluster **node** where the destination runs can be the same of the
source instance or not; generally if the services provided by the instance
do not require exclusive access to a particular system resource (i.e. bind
a given TCP port number) there should be no limitation on running multiple
instances of the same module on the same node. In the end, this kind of
limitations must be managed by the module itself.

In a nutshell, `clone-module`:

- Creates a new module instance for the destination using `cluster/add-module`
- Starts two parallel subtasks:
  - `transfer-state` on the source instance
  - `clone-module` on the destination instance
- When both subtasks are succesfully completed and if a full **replace**
  is required, the source instance is removed with `cluster/remove-module`
- The destination instance has the same MODULE_UUID if replace is required

Modules can implement `transfer-state` and `clone-module` to correctly
support module cloning. The core provides a partial implementation of
those actions that is explained in the following sections.

## Implementation of `transfer-state`

The core implementation of the `transfer-state` action runs a Rsync client
process that transfers the contents of the module volumes and the
AGENT_STATE_DIR (`state/`) directory to the destination instance. The
action input is like:

```json
{
    "replace": true,
    "credentials": ["dokuwiki1", "s3cr3t"],
    "address": "10.5.4.3",
    "port": 20021
}
```

Any step that runs after `05waitserver` can assume the Rsync server
process is ready to receive data.

At step `50sendstate` all volume contents are transferred to the
destination with a `rsync` client that runs in a privileged Podman
container.

This Podman container mounts all module volumes to correctly map numeric
UID/GIDs.

Note also that

* it is possible to run an additional rsync pass before `50sendstate`,
  expecially if there is a lot of data to transfer. Look at the action
  journal to grasp the rsync client invocation command line

* it is possible to add a Rsync filter to select what files/dirs are
  transferred: just drop a `state-filter-rules.rsync` in the
  AGENT_STATE_DIR. See `man 1 rsync` for the file syntax

If the transfer is successful the Rsync client sends a "terminate"
signal to the server.

The module could add service stop/start steps before and after
`50sendstate`. Consider the `replace` input boolean to decide if services
should be left stopped or not.

## Implementation of `clone-module`

The core implementation of the `clone-module` step runs a Rsync server
that receives volume and `state/` directory contents.

The action input could be:

```json
{
    "credentials": ["dokuwiki1", "s3cr3t"],
    "environment": {
        "IMAGE_URL": "ghcr.io/nethserver/dokuwiki:latest",
        "MODULE_UUID": "f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf"
    },
    "port": 20027,
    "volumes": [
        "dokuwiki-data"
    ],
    "replace": false
}
```

Note that the `environment` object (here partially redacted) is a copy of
the source instance module environment.

At step `05create_volumes` the `volumes` input array is used to re-create
the same set of volumes of the source module, if they are not already
created by `create-module`.

The `05replace` step replaces MODULE_UUID, if the `replace` flag is
`true`.

Finally `rsyncd` runs at step `10recvstate`. This step blocks until the
"terminate" signal is received from the client, or the task itself is
canceled.

Additional steps implemented by the module can then adjust the module and
finally start its services.
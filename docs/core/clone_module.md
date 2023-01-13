---
layout: default
title: Cloning a module
nav_order: 10
parent: Core
---

# Cloning a module

A (source) module instance can be cloned by the `cluster/clone-module`
action. This action creates a new module instance (destination) that is
equivalent to the source one.

For example to create a clone of instance `dokuwiki1` on node 1, run

    api-cli run clone-module --data '{"module":"dokuwiki1","replace":false,"node":1}'

The `replace` boolean specifies if the source instance should be removed when
`clone-module` completes with success.

- With `"replace":false` the cloning action makes a copy of the source
  instance. The destination instance differs only for its MODULE_UUID.
- With `"replace":true` the cloning action copies also the MODULE_UUID
  value in the destination instance. Finally it removes the source instance.

The cluster `node` where the destination instance runs can be the same of
the source instance. Generally there should be no limitation on running
multiple module instances on the same node. In some cases the services
provided by the instance might require exclusive access to a particular
system resource, as binding a fixed TCP port number. In this case, avoid
to bind the fixed port number in the action `create-module`, otherwise
cloning becomes impossible.

In a nutshell, `cluster/clone-module`:

- Creates a new module instance for the destination using
  `cluster/add-module`
- Starts two parallel subtasks:
  - `transfer-state` on the source instance
  - `clone-module` on the destination instance
- The destination instance copies the MODULE_UUID of the source, if
  `replace` is true
- When both subtasks are succesfully completed and if `replace` is true,
  the source instance is removed with `cluster/remove-module`

Modules can implement additional steps for the  `transfer-state` and
`clone-module` to correctly support module cloning. The core provides a
partial implementation of those actions that is explained in the following
sections.

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
process is ready to receive data.  The action implementation perform the
transfer in two similar steps to ensure data integrity. 

At steps `10sendpayload` and `50sendstate` all volumes are transferred to
the destination with a `rsync` client that runs in a privileged Podman
container. This container mounts all module volumes to correctly map
numeric UID/GIDs.

The step `10sendpayload` transfers data for the first time, thus it is
expected to be slower.

The step `50sendstate` stops the services, quickly transfers the latest
differences. If the transfer is successful the Rsync client sends a
"terminate" signal to the server. Then the step starts again the services
if the `replace` flag is false.

Additional notes:

* some files are never transfered by the core steps. Hardcoded
  high-priority exclude rules skip the following entries:
  - `state/agent.env`
  - `state/environment`
  - `state/apitoken.cache`

* some additional files with `.clone-module` extension are created and
  transfered. See `clone-module` for details.

* the automatic service stop-and-start logic works for rootless modules
  only. Rootfull modules must implement their steps to stop then start
  each service provided by the module.

* implementation of rootless stop-and-start logic is based on the Systemd
  `transfer-state.target` unit. To exclude a service from the
  stop-and-start cycle, add it to that target.

* it is possible to add a Rsync filter to select what files/dirs are
  transferred: just drop a `etc/state-filter-rules.rsync` under the
  AGENT_INSTALL_DIR. See `man 1 rsync` for the file syntax.

## Implementation of `clone-module`

The core implementation of the `clone-module` step runs a Rsync server
that receives volume and `state/` directory contents.

The action input could be:

```json
{
    "credentials": ["dokuwiki1", "s3cr3t"],
    "port": 20027,
    "volumes": [
        "dokuwiki-data"
    ],
    "replace": false
}
```

At step `05create_volumes` the `volumes` input array is used to re-create
the same set of volumes of the source module, if they are not already
created by `create-module`.

Then `rsyncd` server runs at step `10recvstate`. This server blocks the
step until the "terminate" signal is received from the client, or the task
itself is canceled.

If the state data transfer is successful this step merges the contents of
`environment.clone-module` (source module environment) with the
destination module `environment`. It also preserves the source variable
MODULE_UUID, if the `replace` flag is `true`.

Consider that even if the `state/` directory and all Podman volumes are
copied and merged properly by the core action steps, a module often needs
to invoke actions of other modules, too. For example, it might need to set
up an HTTP route, add a hostname to the TLS certificate or publish some
information in Redis. To implement this, further steps by the module can
self-invoke an action like `configure-module`, or symlink some step of it
directly under the `clone-module/` action directory.

Finally, the `90finalize` step enables and starts the **service and timer
Systemd units** as they are in the source module.

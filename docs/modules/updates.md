---
layout: default
title: Updates
nav_order: 82
parent: Modules
---

## Module updates

The core action `update-module` is the entry point of a module instance
update. It can run the update on multiple instances of the same module at
the same time.

The core then invokes the module `update-module` action. Each module
instance is responsible for updating itself. The core already provides an
extendable implementation of the `update-module` action for the module
agent. The base implementation

- pulls the module image and additional service images
- extracts and installs the module _imageroot_ under the `${AGENT_INSTALL_DIR}`
- runs `systemctl --user daemon-reload`, to refresh Systemd with new unit definitions
- updates the environment variables, keeping track of `PREV_*` values
- runs update scripts installed with the new image, under
  `${AGENT_INSTALL_DIR}/update-module.d/`
- removes stale container images

Modules can provide update scripts to handle upgrade paths like restarting
services and database schema migrations. Executable scripts can be placed
inside `imageroot/update-module.d` which will be then extracted to
`${AGENT_INSTALL_DIR}/update-module.d`. Scripts execution occurs in
alphabetical order. If a script aborts with an exit code, a warning is
printed. Example of a script that restarts a Systemd service,
`imageroot/update-module.d/20restart`:

    #!/bin/bash
    systemctl --user try-restart mymodule.service

Make sure the `20restart` is executable, otherwise it is ignored.

For a critical hook (for example, a required database migration), ship a
non-executable regular file with the same path plus `.fatal`, such as
`imageroot/update-module.d/20restart.fatal`. Its contents are ignored. If
the hook fails, the runner returns its exit code and skips later hooks
and action steps. Hooks without a marker keep the warning-only behavior.
The marker does not make a missing or non-executable hook run.

The core runner must support `.fatal` markers **before** the module update
starts. Once this support is released, set the module's
`org.nethserver.min-core` label to that release or later. Older runners
ignore the marker. A new module action override alone cannot change the
first update: the agent selects action steps before extracting the new
module files. Use normal non-zero exit codes for critical failures;
`agent.set_status()` only accepts `validation-failed`, which is intended
for input validation rather than migration or startup errors.

The runner does not undo a failed hook. Keep recovery data and leave
consumers stopped when schema or service initialization fails.

To update an instance from command line, use:

    update-module ghcr.io/nethserver/mymodule:1.0.0 mymodule1

Multiple instances of the same module can be updated together:

    update-module ghcr.io/nethserver/mymodule:1.0.0 mymodule1 mymodule2

If the given `module_url` is already present in the local Podman storage,
no remote download occurs and the local image is used instead. During
development this might be unwanted and to work around this behavior
pass `--force`:

    update-module ghcr.io/nethserver/mymodule:1.0.0 mymodule1 --force

The `update-module` command refuses to update an instance with an image of
a different module (e.g. updating `nextcloud1` with a `mattermost` image),
even with `--force`. This check is enforced by the CLI only: it can be
bypassed by invoking the `cluster/update-module` action directly, e.g. via
`api-cli` (see below).

The same action can also be invoked directly with `api-cli`, e.g. to script
around it:

    api-cli run update-module --data '{"module_url":"ghcr.io/nethserver/mymodule:1.0.0","instances":["mymodule1"],"force":true}'

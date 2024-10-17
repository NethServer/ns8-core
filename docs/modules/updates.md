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

To update an instance from command line, use:

     api-cli run update-module --data '{"module_url":"ghcr.io/nethserver/mymodule:latest","instances":["mymodule1"]}'

If the given `module_url` is already present in the local Podman storage,
no remote download occurs and the local image is used instead. During
development this might be unwanted and to work around this behavior
set the flag `"force": true` to the `update-module` action:

    api-cli run update-module --data '{"module_url":"ghcr.io/nethserver/mymodule:latest","instances":["mymodule1"],"force":true}'

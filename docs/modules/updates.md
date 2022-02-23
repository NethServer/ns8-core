---
layout: default
title: Updates
nav_order: 14
parent: Modules
---

## Module updates

Each module instance is responsible for updating itself.
The core provides an extendable implementation of the `update-module` action for the module agent.
The basic implementation 

- pulls the module image and additional service images
- updates the environment variables, keeping track of `PREV_*` values
- extracts and installs the module _imageroot_ under the `AGENT_INSTALL_DIR`
- performs `systemctl daemon-reload`
- removes no longer required files

After these steps, the module is running with the old version. To load the new version,
the service should be restarted using a custom script.

Modules can execute custom scripts upon update to handle upgrade paths like database schema migrations.
Executable scripts should be placed inside `imageroot/update-module.d` which will be then extracted to `${AGENT_INSTALL_DIR}/update-module.d`.

Scripts execution occurs in alphabetical order. Execution breaks when a script exits non-zero and the whole action is aborted.

Example of a script to load the updated version `imageroot/update-module.d/20restart`:

    #!/bin/bash
    systemctl --user restart mymodule

Make sure the `imageroot/update-module.d/20restart` is executable.

To update an instance from command line, use:

     api-cli run update-module --data '{"module_url":"ghcr.io/nethserver/mymodule:latest","instances":["mymodule1"]}'

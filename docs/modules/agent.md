---
layout: default
title: Agent
nav_order: 80
parent: Modules
---

# Module agent

* TOC
{:toc}

## Base actions

Agent task processing is defined by actions. An `action` is a directory
containing one ore more executable scripts called `steps`. See [the agent
README](https://github.com/NethServer/ns8-core/blob/main/core/agent/README.md#agent)
for the complete story.

The core provides a set of base actions defined in
`/usr/local/agent/actions` and inherited by module agents.

- `create-module`: executed upon installation, it automatically downloads the module image and all images
  listed inside `org.nethserver.images` label.
- `destroy-module`: executed on module removal; modules can add here
  scripts to clean up configuration, like Traefik routes, Firewalld,
  system-wide Systemd units...
- `get-status`: used by the UI, it returns the current module status, like the node where the application is running, used images and volumes, systemd units status.
  This action works correctly only for rootless modules.
- `list-service-providers`: Look up provider information for a given
  service name and other filtering criteria. See also [Service
  providers]({{site.baseurl}}/modules/service_providers)

Every agent has also a builtin `list-actions` action. This command lists the
actions provided by the agent of `module/traefik1`:

    api-cli run --agent module/traefik1 list-actions | jq

As a convention, all modules with a UI should also implement the following actions:

- `configure-module`: it validates the user input and applies the
  configuration, see [Dokuwiki
  example](https://github.com/NethServer/ns8-dokuwiki/tree/main/imageroot/actions/configure-module)
- `get-configuration`: it should return current configuration, the output
  should be equal or similar to configure-module input, see [Dokuwiki
  example](https://github.com/NethServer/ns8-dokuwiki/tree/main/imageroot/actions/get-configuration)


## Agent environment

Environment variables are copied to Redis at `module/<module_id>/environment`.

Each action has access to environment variables. They are:
1. initialized by `add-module`
2. changed by the agent at run time
3. inherited from the agent environment

Main environment variables are:
- `MODULE_ID`: the ID of the instance, i.e. `myapp1`
- `NODE_ID`: the node ID where the action is running, i.e. `1`
- `<name>_IMAGE`: see [Image download](#image-download)
- `TCP_PORT`, `TCP_PORTS`, `TCP_PORTS_RANGE`: see [Port allocation](#port-allocation)
- `UDP_PORT`, `UDP_PORTS`, `UDP_PORTS_RANGE`: see [Port allocation](#port-allocation)

Other available variables:
- `AGENT_COMFD`: the file descriptor to talk to the agent via [action commands](#action-commands)
- `AGENT_TASK_ID`: the unique ID of current task, i.e. `c0b8b976-9444-42d5-a40b-142a6a483a84`
- `AGENT_TASK_ACTION`: the name of executed action, i.e.: `create-module`
- `AGENT_TASK_USER`: the name of the user that created the action, if available. For example, `admin`
- `AGENT_INSTALL_DIR`
- `AGENT_STATE_DIR`
- `AGENT_ID`
- `IMAGE_ID`
- `IMAGE_URL`
- `IMAGE_DIGEST`
- `IMAGE_REPODIGEST`

Environment variables added by the `add-module` and defined at run time are
saved inside `/home/<module_id>/.config/state/environment`, e.g.
`/home/myapp1/.config/state/environment`.

To run a command within an agent environment, use the `runagent` wrapper. For instance
this command prints the agent environment of the rootfull module `modrootfull1`.

    runagent -m modrootfull1 env


## Agent commands

During the action execution, action steps can talk to the agent using a
simple protocol. For instance, it is possible to set/unset environment
variables that are persisted under the `AGENT_STATE_DIR`.

See all [available action commands](https://github.com/NethServer/ns8-core/blob/main/core/agent/README.md#action-commands).

The protocol is implemented by the Python `agent` package.

Python example:
```python
import agent

agent.set_weight('05longstep', 2)
```

## Validation

See [validation framework]({{site.baseurl}}/core/validation).

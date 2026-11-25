---
layout: default
title: TUN devices
nav_order: 13
parent: Core
---

# TUN devices

NS8 can manage TUN devices and pass them to rootless containers.

## Configuration

Each node is responsible for its own TUN devices.
The node is configured with a simple _tunadm_ API. A module must
be authorized to use it, by adding `node:tunadm` to the module image label
`org.nethserver.authorizations`. For instance, set

       org.nethserver.authorizations=node:tunadm

If the label already has a `node:` value, append the role to it with a
comma, for example `node:tunadm,portsadm`. See [Roles and
authorizations](../agents/#roles-and-authorizations).

Please note that for backward compatibility with core 3.16 the _tunadm_
authorization includes also the [_fwadm_](../firewall) one.

Then the module actions must use the `agent`
Python package to add/remove the tun device needed by the
module.

In `create-module`:

```python
import os
import agent
agent.assert_exp(agent.add_tun("tun1", "192.168.1.1/24"))
```

The tun is owned by the rootless module that calls `add-tun`, so the
module containers can open it. Recent kernels (for example Debian 13
since 6.12.117) refuse a rootless container to open a tun without owner.

The owner is applied when the tun is created, usually at boot. A tun
created by an older core has no owner: `/sys/class/net/<tun>/owner`
contains `-1`. On an affected kernel, the module must fix it in this
order:

1. stop the services that hold the tun,
2. call `remove_tun()` and then `add_tun()`,
3. start the services again.

In `destroy-module`:

```python
import os
import agent
agent.assert_exp(agent.remove_tun("tun1"))
```

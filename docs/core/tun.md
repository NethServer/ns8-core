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

Please note that the _tunadm_ authorization includes also the [_fwadm_](../firewall) one.

Then the module actions must use the `agent`
Python package to add/remove the tun device needed by the
module.

In `create-module`:

```python
import os
import agent
agent.assert_exp(agent.add_tun("tun1", "192.168.1.1/24"))
```

In `destroy-module`:

```python
import os
import agent
agent.assert_exp(agent.remove_tun("tun1"))
```

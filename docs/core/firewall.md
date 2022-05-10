---
layout: default
title: Firewall
nav_order: 13
parent: Core
---

# Firewall

* TOC
{:toc}

NS8 has a basic built-in firewall based on firewalld.

The firewall has 2 zones:
- trusted: it includes containers and wireguard interface
- public: everything else

Default policies:
- SSH port is always open
- all blocked packets are logged

Please note that rootfull containers with port mapping (DNAT) or private
networking are *not supported*.

## Configuration

To permit direct connections from the public zone to a service provided by
a module, the module itself must modify the node firewall configuration.

The node firewall is configured with a simple _fwadm_ API. A module must
be authorized to use it, by adding `node:fwadm` to the module image label
`org.nethserver.authorizations`. For instance, set

       org.nethserver.authorizations=node:fwadm

Then the `create-module` and `destroy-module` actions must use the `agent`
Python package to add/remove the node firewall configuration needed by the
module.

In `create-module`:

```python
import os
import agent
agent.assert_exp(agent.add_public_service(os.environ['MODULE_ID'], ["80/tcp", "443/tcp"]))
```

In `destroy-module`:

```python
import os
import agent
agent.assert_exp(agent.remove_public_service(os.environ['MODULE_ID']))
```

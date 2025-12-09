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

## Node Firewall Management

The following Python functions allow you to manage the node's firewall configuration. To use any of these functions, your module must be granted the `fwadm` role for the target node; otherwise, a permission error will be returned.

### Authorization

To permit direct connections from the public zone to a service provided by
a module, the module itself must modify the node firewall configuration.

The node firewall is configured with a simple _fwadm_ API. A module must
be authorized to use it, by adding `node:fwadm` to the module image label
`org.nethserver.authorizations`. For instance, set

       org.nethserver.authorizations=node:fwadm

### Managing Public Services

The `agent` Python package provides functions to add and remove public services by opening specific ports in the firewall.

#### Adding Public Services

In `create-module`:

```python
import os
import agent
# Raise an exception if add_public_service() returns False
agent.assert_exp(agent.add_public_service(os.environ['MODULE_ID'], ["9010/tcp", "9011/tcp"]), "Firewall service configuration has failed")
```

Function `agent.add_public_service()` can be later invoked with additional
ports, for example during an application update that implements a new
public service. The given port list is _added_ to the existing one. For example:

```python
import os
import agent

agent.add_public_service(os.environ['MODULE_ID'], ["9012/tcp"])
```

If you want to completely replace the port list, set `replace_ports=True`,
for example:

```python
import os
import agent

agent.add_public_service(os.environ['MODULE_ID'], ["9010/tcp","9012/tcp"], replace_ports=True)
```

#### Removing Public Services

In `destroy-module`:

```python
import os
import agent
# Ignore errors on service cleanup
agent.remove_public_service(os.environ['MODULE_ID'])
```

### Managing Rich Rules

For more advanced firewall configurations, you can use rich rules to define complex firewall rules including port forwarding, source address filtering, and more.

#### Adding Rich Rules

Use the `add_rich_rules()` function to apply firewall rich rules on the node:

```python
from agent import add_rich_rules

rules = [
    'rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 to-addr=192.168.1.100',
    'rule family=ipv4 source address=10.1.2.3 port port=22 protocol=tcp accept'
]
result = add_rich_rules(rules)
if not result:
    print("Failed to add rich rules")
```

Each element in the `rules` list should be a complete rich-rule string as accepted by firewall-cmd's `--add-rich-rule` option.

#### Removing Rich Rules

Use the `remove_rich_rules()` function to remove firewall rich rules from the node:

```python
from agent import remove_rich_rules

rules = [
    'rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 to-addr=192.168.1.100'
]
result = remove_rich_rules(rules)
if not result:
    print("Failed to remove rich rules")
```

The rule strings used for removal should match exactly the format used when adding them. Refer to the function docstrings and firewalld documentation for the complete rich rule syntax.

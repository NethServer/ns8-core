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

If the label already has a `node:` value, append the role to it with a
comma, for example `node:fwadm,portsadm`. See [Roles and
authorizations](../agents/#roles-and-authorizations).

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

### Managing Custom Zones

A module that brings its own network interface, for example a tun device
for a VPN, puts it in a custom zone. The zone accepts the given ports and
rich rules from that interface.

In `create-module` or `configure-module`:

```python
import agent

agent.assert_exp(agent.add_custom_zone("tun1", "tun1", ["53/udp", "443/tcp"]))
```

By default the node does not forward the traffic of a custom zone. A
module with a tun device created by `add_tun()` can set `masquerade` to
forward it, with the node address as source. Only what is listed passes:

```python
agent.assert_exp(agent.add_custom_zone("tun1", "tun1", ["53/udp", "443/tcp"], masquerade={
    "internet": True,                          # the Internet, private networks stay closed
    "destinations": ["10.5.5.0/24:443/tcp"],   # HTTPS of the other cluster nodes
}))
```

Each destination is an IPv4 address or network, with an optional port and
protocol, like `192.168.1.10:445/tcp` or `192.168.1.0/24`. `masquerade=True`
is the same as `{"internet": True}`. Do not list the whole cluster VPN
network without a port: the other nodes trust any traffic coming from the
node.

This creates a firewalld policy named after the zone with the `-masq`
suffix, so the zone name is limited to 13 characters. Call
`add_custom_zone()` again without `masquerade` to remove the policy.

In `destroy-module`:

```python
import agent

agent.remove_custom_zone("tun1")
```

### Managing Rich Rules

For more advanced firewall configurations, you can use rich rules to define complex firewall rules including port forwarding, source address filtering, and more.

#### Adding Rich Rules

Use the `add_rich_rules()` function to apply firewall rich rules on the node:

```python
import agent

rules = [
    'rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 to-addr=192.168.1.100',
    'rule family=ipv4 source address=10.1.2.3 port port=22 protocol=tcp accept'
]
result = agent.add_rich_rules(rules)
if not result:
    print("Failed to add rich rules")
```

Each element in the `rules` list should be a complete rich-rule string as accepted by firewall-cmd's `--add-rich-rule` option.

#### Removing Rich Rules

Use the `remove_rich_rules()` function to remove firewall rich rules from the node:

```python
import agent

rules = [
    'rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 to-addr=192.168.1.100'
]
result = agent.remove_rich_rules(rules)
if not result:
    print("Failed to remove rich rules")
```

The rule strings used for removal should match exactly the format stored by firewalld: use `list_rich_rules()` (below) to retrieve them. Refer to the function docstrings and firewalld documentation for the complete rich rule syntax.

#### Listing Rich Rules

Use the `list_rich_rules()` function to retrieve the rich rules currently
configured on the node:

```python
import agent

for rule in agent.list_rich_rules():
    print(rule)
```

The rules are read from the firewalld permanent configuration. Note that
firewalld normalizes the rule syntax when storing it, so the returned strings
may differ from the ones passed to `add_rich_rules()`; do not rely on an exact
string match to correlate them.

If the node cannot be reached the function returns an empty list and logs a
warning, so it never interrupts the calling action. Callers that reconcile a
desired rule set against the current one should pass `raise_on_error=True`
instead, since an empty list is otherwise indistinguishable from "no rules
configured":

```python
try:
    rules = agent.list_rich_rules(raise_on_error=True)
except Exception as e:
    print(f"Failed to list rich rules: {e}")
```

---
layout: default
title: Port allocation
nav_order: 50
parent: Modules
---

# Port allocation

Many web application modules need a TCP or UDP port to run a backend exposed by Traefik.
Such modules can set the `org.nethserver.tcp-ports-demand` and `org.nethserver.udp-ports-demand` which takes an integer number as value.
Example:
```
org.nethserver.tcp-ports-demand=1
org.nethserver.udp-ports-demand=1
```

The randomly-allocated TCP port number will be available inside the `TCP_PORT` environment variable and it will be
available to all step scripts and inside systemd units.
The available environment variables will be:
- `TCP_PORT`, `UDP_PORT`: it is always present and it contains always the first port, i.e. `20001`
- `TCP_PORTS_RANGE`, `UDP_PORTS_RANGE`: only if value is greater than 1, it contains the list of ports in range format,
  i.e `20001-20002`
- `TCP_PORTS`, `UDP_PORTS`: only if value is greater than 1 and less or equal than 8, it contains a comma separated list of
  ports like, i.e. `20001,20002,20003`

Currently, allocated ports are saved in an SQLite database file managed by the local node agent.

## Authorizations

The module requires an additional role to manage port allocation, which is
assigned by setting the `org.nethserver.authorizations` label on the
module image, as shown in the following example:

    org.nethserver.authorizations = node:portsadm

The following additional label values can be used to mix port allocations
with other existing node-related roles:

- `org.nethserver.authorizations = node:fwadm,portsadm`
- `org.nethserver.authorizations = node:tunadm,portsadm`

Note that the value must be exactly one of the above. Other combinations
like `node:portsadm,fwadm` are not valid.

The module will be granted execution permissions for the following actions
on the local node:
- `allocate-ports`
- `deallocate-ports`

These actions can be carried out using the agent library without making
direct node API calls, as explained in the next section.

## Agent library

The Python `agent` library provides a convenient interface for managing port allocation and deallocation, based on the node actions `allocate_ports` and `deallocate_ports`.

This interface dynamically allocate and deallocate ports based on the
module's needs without requiring direct interaction with the node's APIs.

It is optional to specify the `module_id` when calling the port allocation or deallocation functions. By default, if the `module_id` is not provided, the function will automatically use the value of the `MODULE_ID` environment variable. This simplifies the function calls in common scenarios, ensuring the correct module name is used without manual input. However, if needed, you can still explicitly pass the `module_id`.
Note that only the cluster agent can modify the port allocations of other modules.

### Allocate ports

Imagine an application module that initially requires only one TCP port. Later, a new feature is added, and it needs four TCP ports to handle more connections.

If ports are already allocated for this module, the previous allocation will be deallocated, and the new requested range of ports will be allocated. Here’s how this can be done:

```python
import agent

# Allocate 4 TCP ports for the module that is calling the function
allocated_ports = agent.allocate_ports(4, "tcp")
```
or
```python
import agent

# Allocate 4 UDP ports for "my_module" module
allocated_ports = agent.allocate_ports(4, "udp", "my_module")
```

If you want to preserve the previously allocated ports and add a new range of ports, you can use the `keep_existing` parameter:

```python
import agent

# Allocate 4 more TCP ports for the module that is calling the function
allocated_ports = agent.allocate_ports(4, "tcp", keep_existing=True)
```

### Deallocate ports

If the module no longer needs the allocated ports, such as when a feature is removed or disabled, the ports can be easily deallocated:

```python
import agent

# Deallocate TCP ports for the module that is calling the function
deallocated_port_ranges = agent.deallocate_ports("tcp")
```
or
```python
import agent

# Deallocate UDP ports for the "my_module" module
deallocated_port_ranges = agent.deallocate_ports("udp", "my_module")
```
By deallocating the ports, the module frees up the resources, allowing other modules to use those ports.

For more information about the low-level implementation of port allocation, see [Port allocation](../../core/port_allocation) in the core documentation.

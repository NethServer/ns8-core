---
layout: default
title: Port allocation
nav_order: 5
parent: Modules
---

# Port allocation

Many web application modules need a TCP or UDP port to run a backend exposed by Traefik.
Such modules can set the `org.nethserver.tcp-ports-demand` and `org.nethserver.tcp-ports-demand` which takes an integer number as value.
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

Currently allocated port is saved on a sqlite database file on each nodes.

Module can require additional roles to manage ports allocation.
This can be done setting label `org.nethserver.authorizations` on module image as the following example:
```
org.nethserver.authorizations = node:allocate-ports node:deallocate-ports
```

The `ports_manager` library provides functions for managing network ports used by different modules within an application. You can dynamically allocate and deallocate ports based on the module's requirements.

## Importing the Library

To use the `ports_manager` library, you need to import it into your Python script as follows:

```python
import node.ports_manager
```

## Available Functions

### `allocate_ports`

This function allows you to allocate a specific number of ports for a given module and protocol.

- **Parameters**:
  - `required_ports` (*int*): The number of ports required.
  - `module_name` (*str*): The name of the module requesting the ports.
  - `protocol` (*str*): The protocol for which the ports are required (e.g. "tcp" or "udp").

- **Usage Example**:

```python
allocated_ports = node.ports_manager.allocate_ports(5, "my_module", "tcp")
```

### `deallocate_ports`

This function allows you to deallocate all ports previously assigned to a specific module for a given protocol.

- **Parameters**:
  - `module_name` (*str*): The name of the module for which ports should be deallocated.
  - `protocol` (*str*): The protocol for which the ports were allocated (e.g., "tcp" or "udp").

- **Usage Example**:

```python
node.ports_manager.deallocate_ports("my_module", "udp")
```

## Additional Notes

- Ensure to handle exceptions that may be raised during the allocation or deallocation of ports.
- Ports allocated will remain reserved for the specified module until they are explicitly deallocated.
- When using the `allocate_ports` function, if the module already has allocated ports, they will first be deallocated and then reallocated.


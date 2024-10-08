---
layout: default
title: Port allocation
nav_order: 17
parent: Core
---

# Port allocation

The `ports_manager` library handles the low-level port allocation storage.
It registers and stores UDP and TCP port information in a SQLite database
for the local node within the cluster.

The SQLite database file, `ports.sqlite`, is stored in the node's `state/`
directory. To inspect its contents, use the following command:

    podman run -i --rm --volume /var/lib/nethserver/node/state/ports.sqlite:/srv/ports.sqlite:z alpine ash -c 'apk add sqlite ; sqlite3 /srv/ports.sqlite' <<<"SELECT *, 'tcp' FROM TCP_PORTS UNION SELECT *, 'udp' FROM UDP_PORTS;"

For more information, refer to [Port
allocation](../../modules/port_allocation) under the Modules section.


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

- **Return** a tuple (start_port, end_port) if allocation is successful, `None` otherwise.

- **Usage Example**:

```python
allocated_ports = node.ports_manager.allocate_ports(5, "my_module", "tcp")
print(f"my_module ports allocated: {allocated_ports}")
```

### `deallocate_ports`

This function allows you to deallocate all ports previously assigned to a specific module for a given protocol.

- **Parameters**:
  - `module_name` (*str*): The name of the module for which ports should be deallocated.
  - `protocol` (*str*): The protocol for which the ports were allocated (e.g., "tcp" or "udp").

- **Return** a tuple (start_port, end_port) if deallocation is successful, `None` otherwise.

- **Usage Example**:

```python
deallocated_ports = node.ports_manager.deallocate_ports("my_module", "udp")
print(f"my_module ports deallocated: {deallocated_ports}")
```

## Additional Notes

- Ensure to handle exceptions that may be raised during the allocation or deallocation of ports.
- Ports allocated will remain reserved for the specified module until they are explicitly deallocated.
- When using the `allocate_ports` function, if the module already has allocated ports, they will first be deallocated and then reallocated.
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

Currently last allocated port is saved inside Redis at `node/<node_id>/tcp_ports_sequence`, `node/<node_id>/udp_ports_sequence`.

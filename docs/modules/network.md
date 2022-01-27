---
layout: default
title: Network
nav_order: 7
parent: Modules
---

# Network

Since NS8 has no firewall, network access to applications should carefully restricted
using the correct IP binding.

As a general rule, any module which doesn't require a well-known port, should request a random port
using `org.nethserver.tcp-ports-demand` label and bind to that port only on 127.0.0.1,
i.e. `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:${TCP_PORT}:8080 ...`

Modules using well-known port, should bind to 127.0.0.1 and to all IPs where the service
must be exposed, like VPN IPs, i.e `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:19999:19999 $EXTRA_LISTEN ...`,
where `EXTRA_LISTEN` could be the bind to the VPN `-p 10.5.4.1:19999:19999`.


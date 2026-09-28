---
layout: default
title: Network
nav_order: 70
parent: Modules
---

# Network

The [system firewall]({{site.baseurl}}/core/firewall) has two active
zones:

* `public` -- every network interface is added to it; limited TCP/UDP
  ports are allowed from the public zone.
* `trusted` -- cluster VPN network; any connection is allowed from the
  trusted zone.

As a general rule, any module which doesn't require a well-known port,
should request a random port using `org.nethserver.tcp-ports-demand`
and `org.nethserver.udp-ports-demand` labels.

The following example creates a private network namespace and starts a TCP
proxy to connect port 8080 inside the container from `${TCP_PORT}` or `${UDP_PORT}`:

    /usr/bin/podman run ... --publish ${TCP_PORT}:8080 ...

Web applications are usually configured as backends for the local Traefik
HTTP proxy. They can bind only the loopback IP address:

    /usr/bin/podman run ... --publish 127.0.0.1:${TCP_PORT}:8080 ...

The next example does not use any TCP proxy and is more performant. It
requires to configure the listening service in the container to use
directly TCP port `${TCP_PORT}`. The container shares the network
namespace with host machine:

    /usr/bin/podman run ... --network=host ...

NS8 modules do not use Podman `bridge` networks. Rootless modules cannot
attach containers to a bridge of the host network. Rootful modules must
not create one, because it changes the host network configuration.

Modules using a well-known port, can bind any IP address for that port.
For instance:

    /usr/bin/podman run ... --publish 25:25

Such modules must be properly authorized to open the well-known port in
the system firewall. See [system
firewall]({{site.baseurl}}/core/firewall#configuration) for details.

## Reach the node from a container

A container with a `private` network, like the rootless default (Pasta),
cannot connect to services bound to the node loopback address,
`127.0.0.1`. To reach a node service from such a container, use the
`host.containers.internal` name. It resolves to an address of the node,
and the service must listen on it.

For instance, connect to the [LDAP proxy]({{site.baseurl}}/core/user_domains#ldap-service-discovery)
at `host.containers.internal:<port>`.

To reach an application behind the Traefik HTTP proxy of the same node,
map its host name to the node address:

    /usr/bin/podman run ... --add-host app.example.org:host-gateway ...

Podman resolves `host.containers.internal` to a special address, that is
translated to the node IP address as it was when the container started.
If the node IP address changes, restart the container to reach the node
again.

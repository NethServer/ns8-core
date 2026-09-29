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
`127.0.0.1`. By default Pasta also copies the node main IP address into
the container: a connection to that address, or to a name that resolves
to it, stays inside the container and does not reach the node. Other node
addresses, like the cluster VPN IP, and other hosts are reachable as
usual.

Pick the first case that fits the module. In a Podman pod, pass the
`--network` and `--add-host` options to `podman pod create`.

### Case 1: node services

To reach a node service, use the `host.containers.internal` name. It
resolves to a special address that Pasta translates to the node main IP
address. The service must listen on that address, not only on
`127.0.0.1`, and it sees connections coming from the node main IP
address.

For instance, connect to the [LDAP proxy]({{site.baseurl}}/core/user_domains#ldap-service-discovery)
at `host.containers.internal:<port>`. To keep an existing host name in
the application configuration, map it to the same address:

    /usr/bin/podman run ... --add-host accountprovider:host-gateway ...

The translation uses the node IP address as it was when the container
started. If the node IP address changes, restart the container to reach
the node again.

### Case 2: applications of the same node, with known host names

To reach an application behind the Traefik HTTP proxy of the same node,
like the module's own public host name, map each host name to the node
address:

    /usr/bin/podman run ... --add-host app.example.org:host-gateway ...

This also applies to the node FQDN: inside a container it resolves to a
loopback address copied from the node `/etc/hosts` file, unless it is
mapped as above.

### Case 3: host names not known in advance

If the module must reach any host name served by the same node, for
example a document server that calls back arbitrary web applications,
give the container a private address:

    /usr/bin/podman run ... --network=pasta:-a,10.0.2.100,-n,24,-g,10.0.2.2 ...

The node main IP address is then reachable, like any other address.
However, `host.containers.internal` and `host-gateway` no longer reach the
node. Connect to node services through the node cluster VPN IP address,
stored in the `node/{id}/vpn` key of [Redis]({{site.baseurl}}/core/database)
(`ip_address` field). For instance:

    /usr/bin/podman run ... --add-host accountprovider:${NODE_VPN_IP} ...

Here `NODE_VPN_IP` is an example environment variable that the module
must set itself. The VPN IP address does not change when the node main
IP address changes.

In every case, published ports keep the client source IP address, and the
node loopback interface is not exposed to the container. Do not use
`slirp4netns:allow_host_loopback=true` or the Pasta `--map-gw` option: they
expose every service bound to the node loopback address.

Some software expects a network interface named `eth0`, while Pasta
copies the node interface name. Rename it with the `-I` option, for
instance `--network=pasta:-I,eth0`.

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

### Node services

To reach a node service, use the `cluster-localnode` name. The core adds
it to the node `/etc/hosts` file, resolving to the node cluster VPN IP
address, and Podman copies that file into the container. The name works
with any container network, and needs no `podman run` options. The
service must listen on the VPN IP address, or on all addresses, and it
sees connections coming from the VPN IP address.

For instance, connect to the [LDAP proxy]({{site.baseurl}}/core/user_domains#ldap-service-discovery)
at `cluster-localnode:<port>`.

The VPN IP address does not change when the node main IP address changes.
It is assigned when the cluster is created or the node joins it. Until
then `cluster-localnode` resolves to `127.0.0.1`, that is the container
itself. Application modules are installed later, but core modules may
run before: they must restart their containers after the node gets its
VPN IP address, because Podman copies `/etc/hosts` when a container is
created.

Do not use the `cluster-leader` name: on the leader node it resolves to
`127.0.0.1`. Prefer `cluster-localnode` also to the Podman
`host.containers.internal` name: it does not work with the private
address of the second case below, and it follows the node main IP
address only as it was when the container started.

### Applications of the same node

Applications behind the Traefik HTTP proxy of the same node, like the
module's own public host name, resolve to the node main IP address. Pick
one of the following cases.

**Case 1: known host names.** Keep the default network and map each host
name to the node address:

    /usr/bin/podman run ... --add-host app.example.org:host-gateway ...

This also applies to the node FQDN: inside a container it resolves to a
loopback address copied from the node `/etc/hosts` file, unless it is
mapped as above. The mapping uses the node main IP address as it was when
the container started: if it changes, restart the container.

**Case 2: host names not known in advance.** If the module must reach any
host name served by the same node, for example a document server that
calls back arbitrary web applications, give the container a private
address:

    /usr/bin/podman run ... --network=pasta:-a,10.0.2.100,-n,24,-g,10.0.2.2 ...

The node main IP address is then reachable, like any other address. The
container cannot reach hosts of the `10.0.2.0/24` network, and
`host-gateway` no longer reaches the node.

### Other notes

- Without `--network`, `podman run` and `podman pod create` use the
  rootless default network, Pasta: add the option only when a case above
  requires it. Containers of a pod share the pod network: pass
  `--network` and `--add-host` to `podman pod create`, never to the
  `podman run` command of a pod container. Podman rejects `--add-host`
  there, and a container started with any `--network` value gets its own
  network, separated from the pod.
- Published ports keep the client source IP address, and the node
  loopback interface is not exposed to the container. Do not use
  `slirp4netns:allow_host_loopback=true` or the Pasta `--map-gw` option:
  they expose every service bound to the node loopback address.
- Some software expects a network interface named `eth0`, while Pasta
  copies the node interface name. Rename it with the `-I` option, for
  instance `--network=pasta:-I,eth0`.

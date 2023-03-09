---
layout: default
title: VPN
nav_order: 8
parent: Core
---

# VPN

Each node is connected to the leader using
[WireGuard](https://www.wireguard.com/) in a star network topology, where
the leader node is the hub of the VPN links so it is connected to any
other node.

The VPN has its own network, the default is `10.5.4.0/24` and it listens on port 55820.
Make sure this address is not already used inside your existing network and the UDP port 55820 is open on your firewall.
You can change both configurations during cluster initialization.

## Promote a worker to leader

Changing the cluster leader implies all VPN connections must be
reconfigured to point to the new VPN hub. The future leader node must
provide a VPN endpoint that is reachable by any other node. For this
reason, **be prepared in advance** to deal with a leader failure by
configuring a reachable endpoint address on one or more worker nodes!

For instance if node with ID `3` has a public and static IP address
`1.2.3.4` log on the leader node and run the following command:

    redis-cli HSET node/3/vpn endpoint 1.2.3.4:55820

A DNS host name can be set too, if there is one available:

    redis-cli HSET node/3/vpn endpoint node3.example.com:55820

When the promotion time comes and depending on the old leader node state,
an automatic or manual procedure executes the new leader promotion.

- If the old leader node is alive, it can notify the cluster nodes to do
  the change automatically.

- If the old leader node is dead, a manual change is needed on every node.

### Alive leader

Assuming `3` is the node ID of the new leader,

1. log on to the old leader node,

1. run the `promote-node` action:

       api-cli run promote-node --data '{"node_id":3}'

In a few seconds the change is propagated to all nodes and the old leader
becomes a worker.

### Dead leader

Assuming `3` is the node ID of the new leader, for each worker node

1. log on to the worker node
1. run the following command for each of them:

       switch-leader 3

The order of the nodes is not important but it is wise to start with the
new leader node.

## Wireguerd configuration

The VPN is configured in every node to use the interface `wg0`. To inspect
the VPN status run

    wg show wg0

Example:
```
[root@cs1 ~]# wg show wg0
interface: wg0
  public key: 4fYZdrhcdjQo3cyvLIJC1TaFR8G6KRT7hZgZNPRoGzo=
  private key: (hidden)
  listening port: 55820

peer: 3Ft+OeV/fdRXnWKcyD8e+emsMxzQGZ1grWqR+/FCoEc=
  endpoint: 178.128.169.93:55820
  allowed ips: 10.5.4.2/32
  latest handshake: 1 minute, 28 seconds ago
  transfer: 96.23 MiB received, 50.01 MiB sent
  persistent keepalive: every 25 seconds

peer: nvu810d2Q3LGUwJNmJawXj7u16+JUukgreQo4NCcACc=
  endpoint: 104.248.130.240:55820
  allowed ips: 10.5.4.3/32
  latest handshake: 1 minute, 34 seconds ago
  transfer: 12.14 MiB received, 13.57 MiB sent
  persistent keepalive: every 25 seconds
```

Configuration file is saved inside `/etc/wireguard/wg0.conf`, public and private keys are saved inside `/etc/nethserver/wg0.{pub,key}`.

The configuration file is initialized by `create-cluster` and `join-node`
actions. Its contents reflect the `wg0` VPN runtime settings, which are
dumped by the `apply-vpn-routes` command (handler of the "vpn-changed" and
"leader-changed" events) and by `systemctl stop wg-quick@wg0`, using the
`wg-quick save` command.

Local modification to the configuration file can be merged into the
runtime state with the following command:

    systemctl reload wg-quick@wg0

Runtime changes can be persisted to the configuration file with the
following command, until the next `apply-vpn-routes` run:

    wg-quick save wg0

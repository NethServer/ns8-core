---
layout: default
title: VPN
nav_order: 8
parent: Core
---

# VPN

Each node is connected to the leader using [WireGuard](https://www.wireguard.com/) in a star network topology.

The VPN has its own network, the default is `10.5.4.0/24` and it listens on port 55820.
Make sure this address is not already used inside your existing network and the UDP port 55820 is open on your firewall.
You can change both configurations during cluster intialization.

To inspect the VPN status you can use both `systemctl status wg-quick@wg0.service` and `wg show` commands.
When the VPN is correctly configured, the leader now always takes the first IP for the above CIDR.

Example:
```
# ip address show wg0
4: wg0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 qdisc noqueue state UNKNOWN group default qlen 1000
    link/none 
    inet 10.5.4.1/32 scope global wg0
       valid_lft forever preferred_lft forever
```

Configuration file is saved inside `/etc/wireguard/wg0.conf`, public and private keys are saved inside `/etc/nethserver/wg0.{pub,key}`.

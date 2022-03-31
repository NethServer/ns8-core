---
layout: default
title: Firewall
nav_order: 13
parent: Core
---

# Built-in firewall

* TOC
{:toc}

NS8 has a basic built-in firewall based on firewalld.

The firewall has 2 zones:
- trusted: it includes containers and wireguard interface
- public: everything else

Default policies:
- SSH port is always open
- all blocked packets are logged

FIXME: describe how to use it


Podman firewalld integration is only partial:
- podman 4, in CentOS Stream, automatically configure itself for firewalld using netavark
- pdoman 3.x, on Debian and Ubuntu, needs to explicitly configure the CNI plugin

Please note that rootfull containers with port mapping (DNAT) are *not supported*.

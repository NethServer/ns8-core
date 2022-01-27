---
layout: default
title: Filesystem
nav_order: 1
parent: Core
---

# Filesystem locations

The core executable scripts and binaries are installed under
`/usr/local/bin` and `/usr/local/sbin`.

* The `agent` and `api-server` commands are implemented in
  [Go](https://golang.org/).
* Python scripts run within a special environment known as "agent environment".
  It is installed inside `/usr/local/agent/pyenv`.

Most NS8 core files reside in `/etc/nethserver`,
`/var/lib/nethserver/` and `/usr/local/agent/`. 

See `/var/lib/nethserver/node/state/coreimage.lst` for a complete list.

Rootless modules are totally contained inside UNIX user home directory, like `/home/trafeik1`.
Rootfull modules are homed under `/var/lib/nethserver/samba1`.

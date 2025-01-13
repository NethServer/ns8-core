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

## Custom base path for home directories

NS8 core uses the common `/home` path for users' home directories, but you
can configure a different path. Follow these steps to configure a node's
agent for this purpose:

1. Create the alternative base path:

       mkdir -m 0755 /home1

2. For systems with SELinux, make it equivalent to `/home`:

       semanage fcontext -a -e /home /home1

3. Fix SELinux labels for the new path:

       restorecon -R -v /home1

4. Edit the node's environment file:

       runagent -m node vi environment

   Add or modify the `HOME_BASEDIR` line in the environment file:

       HOME_BASEDIR=/home1

From now on, new module instances will use `/home1` as their base
directory. Existing modules are not affected and will retain their current
home directory.

Refer to the `semanage-fcontext` man page for additional SELinux-related
information.

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
can configure a different path if it is a device mount point. Follow these
steps to configure a node's agent for this purpose:

1. Create the alternative base path:

       mkdir -m 0755 /home1

2. Mount the device on the new path:

       mount /dev/some /home1

   To make the mount persistent, either edit `/etc/fstab` or create a
   systemd `.mount` unit. Ensure the device is correctly mounted after a
   system reboot.

3. Configure the node agent to use `/home1` as base directory for new
   modules:

       runagent -m node configure-home-basedir --set /home1

From now on, new module instances will use `/home1` as their base
directory. Existing modules are not affected and will retain their current
home directory.

The `configure-home-basedir` command modifies SELinux configuration. Run
the following command to inspect the current customizations:

    semanage fcontext -l -C

Refer to the `semanage-fcontext` man page for additional SELinux-related
information.

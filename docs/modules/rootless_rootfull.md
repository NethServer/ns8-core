---
layout: default
title: Rootless vs Rootfull
nav_order: 2
parent: Modules
---

## Rootless vs Rootfull

The main difference between *rootless* and *rootfull* modules, as
suggested by the adjective, is the Unix user running the module processes
and its system privileges.

**Rootfull** modules run as `root` (EUID 0).

**Rootless** modules run as a normal Unix user. The Unix user account is
created by the node agent when the module instance is added. It has
session lingering enabled: it automatically starts a persistent Systemd
user manager instance at system boot.

The two types of modules have a similar filesystem structure. Rootless
modules are installed to `/home/<module_id>/.config`, whilst rootfull are
installed to `/var/lib/nethserver/<module_id>`.

Rootless modules also have Systemd user units installed under
`~/.config/systemd/user`. Recall that system-wide user units are installed
by the core under `/etc/systemd/user`.

To inspect and modify a rootless module start a SSH session. SSH is
preferred to `su - <user>` because the latter does not properly initialize
the Systemd session environment. For instance, to check if Traefik is running:

    ssh traefik1@localhost
    systemctl --user status


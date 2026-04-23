---
layout: default
title: Rootless vs Rootfull
nav_order: 20
parent: Modules
---

# Rootless vs Rootfull

In the Podman world, non-privileged users can run containers on their own:
those are the so-called *rootless* containers.  In NethServer 8, we borrow
the same word from Podman and use it in the context of modules, together
with its opposite, *rootfull*.

To inspect a rootless module start Bash with the `runagent` command to
properly initialize the Systemd runtime environment. For instance, to
check if Traefik is running:

    runagent -m traefik1
    # A Bash prompt from traefik1's session appears
    systemctl --user status traefik
    # Systemctl prints its output

If you type the above commands in the same line, `runagent` executes
`systemctl` without forking an interactive Bash shell:

    runagent -m traefik1 systemctl --user status traefik
    # Systemctl prints its output from traefik1's session

Let's see the differences of rootless modules vs rootfull modules.

## Unix user

The main difference between *rootless* and *rootfull* modules, as
suggested by the adjective, is the Unix user running the module processes
and its system privileges.

**Rootfull**: module containers run as `root` (EUID 0).

**Rootless**: module containers run as a normal Unix user. The Unix user account is
created by the node agent when the module instance is added. It has
session lingering enabled: it automatically starts a persistent Systemd
user manager instance at system boot.

To check if a module is rootless or not, in Python write:

```python
import os
if os.geteuid() == 0:
    print('ROOTFULL')
else:
    print('ROOTLESS')
```

Same thing, in Bash:

```bash
if [[ $EUID == 0 ]]; then
    echo ROOTFULL
else
    echo ROOTLESS
fi
```

As alternative print the effective user ID with

    id -u

To protect against system configuration errors, starting with Core 3.20,
Unix users created for rootless modules are assigned the `/sbin/nologin`
shell by default. To upgrade existing modules to the new shell policy, run
a command like this:

```bash
runagent -l | xargs -l -t -r -- usermod -s /sbin/nologin
```

Note that the above command may raise harmless errors for `cluster`,
`node` and other rootfull module names, which have no corresponding Unix
user.

An interactive shell is not required by NS8 rootless modules. However, if
you need to enable `/bin/bash` for any reason (for example, to log in as
the module user and transfer files with `ssh+rsync`), execute:

```bash
usermod -s /bin/bash myapp1
```

If the new default policy is too strict, customize the shell assigned to
new modules by `add-module` action. Run this command on each cluster node:

```bash
runagent -m node python3 -c 'import agent ; agent.set_env("LOGIN_SHELL", "/bin/bash")'
```

The above command sets the custom shell executable path in the node's
`LOGIN_SHELL` agent environment.


## Filesystem paths

The two types of modules have a similar filesystem structure. Rootless
modules are installed to `/home/<module_id>/.config`, whilst rootfull are
installed to `/var/lib/nethserver/<module_id>`.

## Systemd units

Rootless modules also have Systemd user units installed under
`~/.config/systemd/user`. Recall that some system-wide user units are
installed by the core under `/etc/systemd/user`. When running the
`systemctl` command, add the `--user` flag. Eg.

    systemctl --user status traefik

Rootfull modules share the system-wide Systemd directories. Their units
are installed under `/etc/systemd/system`. As they share the same
directory, unit files must be named as `<module_id>` or they must use the
`<module_id>-` prefix to avoid naming clashes with other instances of the
same module. Eg:

    cat /etc/systemd/system/samba1.service
    systemctl start samba1

## Volumes

Rootfull modules share the same Podman volumes namspace. As consequence,
rootfull modules must use the `<module_id>-` prefix for their volume names
to avoid volume naming clashes. Eg.

    samba1-data
    samba1-config

Rootless modules can use any volume name because the Podman volume
namespace is private for the module.

This command prints out the filesystem path where Podman stores volumes
data:

{% raw %}
    podman system info --format='{{.Store.VolumePath}}'
{% endraw %}

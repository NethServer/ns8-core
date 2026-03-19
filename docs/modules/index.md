---
layout: default
title: Modules
nav_order: 5
has_children: true
---

# Modules: implementation details

A module is an NS8 application which is usually composed by a backend and
a frontend.

Every time a module instance is added to the cluster, the new instance
is named as the module itself followed by a progressive number
starting from 1. Given a module named `myapp`, instances will be named
`myapp1`, `myapp2`, etc. This uniquely identifies the instance inside the cluster.

Modules can be managed using these commands:

- `add-module <module> <node_id>`: install a module on the given node with ID `node_id`; search for `module` inside enabled repositories and install
  latest available version. If `module` is a image registry URL, just install the module straight from it; this method can be used to install
  customized images.
  - `--nologin`: set `/usr/sbin/nologin` as the login shell for the new rootless module Unix user (see [Login shell hardening](#login-shell-hardening)). Applicable to **rootless modules only**; passing it for a rootfull module is an error.
- `remove-module [--no-preserve] <module>`: remove an installed module; if `--no-preserve` is given, erase also all module data

## Login shell hardening

By default, rootless module users are created with `/bin/bash` as their login shell.
Passing `--nologin` to `add-module` sets `/usr/sbin/nologin` instead.
This option is **only applicable to rootless modules**; attempting to use it with a rootfull module will produce an error.

```
add-module traefik 1 --nologin
```

This prevents direct interactive login to the module Unix user via SSH or console,
which is a useful security hardening measure aligned with the principle of least privilege.

### How it works

- The `useradd` call in the node `add-module` action uses `/usr/sbin/nologin` instead of
  `/bin/bash` when `nologin` is set.
- Systemd user units (lingering) and the `runagent` tool are **not affected**: they do not
  rely on the login shell to start or manage module processes.
- Rootless containers managed by Podman under the module user continue to work normally.

### Support access

Setting `nologin` does **not** prevent support access to the module environment.
The standard procedure remains:

1. Connect as `root` (SSH or console).
2. Use `runagent -m <module_id> bash` to enter the module agent environment.

Direct SSH access to the module user account (e.g. `ssh mail1@node`) will be
blocked and print `This account is currently not available.` This is expected
behavior.

### Risks and caveats

- **Not tested extensively**: this option has not been subject to in-depth
  integration testing on all module types. Unexpected behavior is possible;
  revert by running `usermod -s /bin/bash <module_id>` if problems arise.
- **Existing modules are not affected**: `--nologin` applies only at install
  time. To harden an already-installed module, run
  `usermod -s /usr/sbin/nologin <module_id>` manually.
- **rsync-over-SSH transfers**: some internal operations (e.g. replication,
  clone) that use rsync over SSH tunnels between nodes go through the `root`
  user, not the module user, so they are not affected.
- **Reversible**: the shell can be restored at any time with
  `usermod -s /bin/bash <module_id>`.


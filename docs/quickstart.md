---
layout: default
title: Quickstart
nav_order: 2
---

# Installation & First steps

This is just a prototype, **do not** use in production.
{: .label .label-yellow}

* TOC
{:toc}

## System requirements

NethServer 8 can be deployed on one ore more nodes in a cluster scenario.

Supported distributions:
- CentOS Stream 9
- Debian 11 (bullseye)

Hardware requirements for a single node installation:
- 2 vCPU/cores
- 2GB RAM
- 20GB disk

More nodes can be added later. When adding a new node, you should use
similar hardware and the same distribution installed on the leader node.

Always install NethServer 8 on a clean server machine, do not install it on a desktop.

## Core installation

Pick your preferred distribution between supported ones and make sure it's up to date. 
Also ensure that the system firewall is not blocking any connection.

Start the installation procedure as `root`:
```
curl https://raw.githubusercontent.com/NethServer/ns8-scratchpad/main/core/install.sh | bash
```

At the end of the install script the UI is available at **https://\<server_ip_or_fqdn\>/cluster-admin/**:

- default user: `admin`
- default password: `Nethesis,1234`

Some configuration tasks can be completed also by invoking additional
commands: follow on-screen instructions printed by `install.sh`.

Run either new cluster initialization (`create-cluster`) or joining an existing cluster (`join-cluster`).

### Install a development branch

Developers may prefer to run `install.sh` with one or more images from a
development branch. The first argument selects the branch name, following
arguments the module names to be pulled from that branch.

    bash install.sh mybranch module1 module2 ...

## Applications

Core applications installed by default:
- [Traefik](https://github.com/NethServer/ns8-scratchpad/blob/main/traefik/README.md) -- see how to request a Let's Encrypt [certificate for the FQDN](https://github.com/NethServer/ns8-scratchpad/blob/main/traefik/README.md#set-certificate)
- [Loki](https://github.com/NethServer/ns8-scratchpad/blob/main/loki/REDME.md) (only on the leader node)
- [Promtail](https://github.com/NethServer/ns8-scratchpad/blob/main/promtail/README.md)
- [LDAP proxy](https://github.com/NethServer/ns8-scratchpad/blob/main/ldapproxy/README.md)

Other available core applications:
- [Samba](https://github.com/NethServer/ns8-scratchpad/blob/main/samba/README.md)


### Application installation

Applications can be installed directly from the user interface.

If you prefer the command line, use the `add-module` command:
```
add-module <module> <node_id>
```
The  above command will try to install the latest stable version of the module.

Example to install Dokuwiki on node 1:
```
add-module dokuwiki 1
```

You can also install an image which is not still available inside the repository by using
its URL.

Example to install Dokuwiki directly from the image registry:
```
add-module ghrc.io/nethserver/dokuwiki:mydev 1
```

Many applications need a configuration step after install, for more info, 
please refer to the README of each application.

Partial list of available applications:

- [Dokuwiki](https://github.com/NethServer/ns8-scratchpad/blob/main/dokuwiki/README.md)
- [Nextcloud](https://github.com/NethServer/ns8-scratchpad/blob/main/nextcloud/README.md)
- [Mail](https://github.com/NethServer/ns8-scratchpad/blob/main/netdata/README.md)

## Uninstall

The `uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*` and `/var/lib/nethserver/*`!

    bash /var/lib/nethserver/node/uninstall.sh


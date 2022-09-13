---
layout: default
title: Quickstart
nav_order: 2
---

# Installation & First steps

This is just a prototype, **do not** yet use in production.
{: .label .label-yellow}

* TOC
{:toc}

## System requirements

NethServer 8 can be deployed on one ore more nodes in a cluster scenario.

Supported distributions and architectures:
- Rocky Linux 9.0 (Blue Onyx) x86-64
- CentOS Stream 9 x86-64
- Debian 11 (bullseye) x86-64

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

First, ensure `curl` is installed:
- for Debian: `apt-get install -y curl`
- for CentOS: `dnf install -y curl`

Start the installation procedure as `root`:
```
curl https://raw.githubusercontent.com/NethServer/ns8-core/main/core/install.sh | bash
```

At the end of the install script the UI is available at **https://\<server_ip_or_fqdn\>/cluster-admin/**:

- default user: `admin`
- default password: `Nethesis,1234`

Some configuration tasks can be completed also by invoking additional
commands: follow on-screen instructions printed by `install.sh`.

Run either new cluster initialization (`create-cluster`) or joining an existing cluster (`join-cluster`).

### Install custom images

Developers may prefer to run `install.sh` with one or more images from a
development branch or alternative registries.

    bash install.sh ghcr.io/nethserver/core:latest ghcr.io/nethserver/traefik:mybranch

### Install customization

The install script also accepts the following environment variables:
- `TESTING`: override testing flag inside the `default` repository. It can be `0` (disabled) or `1` (enabled), default is `0`
- `REPMOD`: override `default` software repository URL, it could be something like `https://mycustomrrepo.server.test/repomd`
- `ADMIN_PASSWORD`: override default admin password

## Applications

Core applications installed by default:
- [Traefik](https://github.com/NethServer/ns8-core/blob/main/traefik/README.md) -- see how to request a Let's Encrypt [certificate for the FQDN](https://github.com/NethServer/ns8-core/blob/main/traefik/README.md#set-certificate)
- [Loki](https://github.com/NethServer/ns8-core/blob/main/loki/REDME.md) (only on the leader node)
- [Promtail](https://github.com/NethServer/ns8-core/blob/main/promtail/README.md)
- [LDAP proxy](https://github.com/NethServer/ns8-core/blob/main/ldapproxy/README.md)

Other available core applications:
- [Samba](https://github.com/NethServer/ns8-core/blob/main/samba/README.md)


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
add-module ghcr.io/nethserver/dokuwiki:mydev 1
```

If the given image is already present in the local Podman storage, no
remote download occurs and the local image is used instead. During
development this might be unwanted and to work around this behavior
execute the following command in every cluster node, before `add-module`:

    podman rmi ghcr.io/nethserver/dokuwiki:mydev

Many applications need a configuration step after install, for more info, 
please refer to the README of each application.

Partial list of available applications:

- [Dokuwiki](https://github.com/NethServer/ns8-dokuwiki/blob/main/README.md)
- [Nextcloud](https://github.com/NethServer/ns8-core/blob/main/nextcloud/README.md)

## Uninstall

The `uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*` and `/var/lib/nethserver/*`!

    bash /var/lib/nethserver/node/uninstall.sh


---
layout: default
title: Quickstart
nav_order: 2
---

# Installation & first steps

* TOC
{:toc}

See the [Administrator manual]({{site.admin_manual}}) for system requirements and installation instructions.

## Core: install custom images

Developers may prefer to run `install.sh` with one or more images from a
development branch or alternative registries.

    curl https://raw.githubusercontent.com/NethServer/ns8-core/main/core/install.sh > install.sh 
    bash install.sh ghcr.io/nethserver/core:latest ghcr.io/nethserver/traefik:mybranch

See also [module override](#module-override).

The install script also accepts the following environment variables:
- `TESTING`: This flag allows you to override the testing setting within
  the `default` repository. It can be set to `0` (disabled) or `1`
  (enabled), with a default value of `0`. Note that this flag will
  automatically reset to `0` after the next core update.
- `REPMOD`: override `default` software repository URL, it could be something like `https://mycustomrrepo.server.test/repomd`
- `ADMIN_PASSWORD`: override default admin password

## Applications: install from CLI

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
please refer to [Administrator manual]({{site.admin_manual}}).

### Module override

Sometimes you may need to install a specific tag for a module even from the UI.
As an example, this can be usefull to override the image used by a
user domain provider.

Command to override the Samba module:
```
redis-cli hset cluster/override/modules samba ghcr.io/nethserver/samba:rootless
```

You can now provision the user domain directly from the web interface.

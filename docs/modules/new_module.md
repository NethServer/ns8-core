---
layout: default
title: New module tutorial
nav_order: 12
parent: Modules
---

# How to create a rootless module

Each module is distributed using a container image.
This document will guide you, a developer, on how to create a new rootless module.
The module will be a simple NethServer package for [Dokuwiki](https://www.dokuwiki.org/dokuwiki).

Before moving further, please take a look to [how NethServer is designed](design.md).
See also [Technical details](details.md).

## Step 0: install NS8

To develop a new module first install NS8 by following the Installation page.
Then, make sure `buildah` is installed to the system.

For Debian:
```
apt-get install buildah
```

For CentOS:
```
dnf install buildah
```

To push push images to the registry, you must configure the authentication.
Create a [GitHub Personal Access Token (PAT)](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
for the **ghcr.io** registry (for read-only access `read:packages private` scope should be enough) then run the following command, specifying
your GitHub user name and providing the generated PAT as password:
```
buildah login ghcr.io
```

## Step 1: module structure

A module is usually composed by 2 main parts:
- `imageroot` directory, contains all scripts to configure the module
- `build-images.sh`, a script to manually build one or more images of the module and push them inside the image registry.
- `ui` directory, contains all user interface files
- `README.md` describes module implementation and usage

### build-images.sh

The script creates an empty container image and publish it to the registry on user request.

The main parts to look for are:
- `repobase`: the URL of the remote image registry
- `reponame`: the name of the module
- labels:
  - `org.nethserver.tcp-ports-demand=1`: the container will use a single TCP port accessible using Traefik as proxy
  - `org.nethserver.rootfull=0`: this is a rootless module
  - `org.nethserver.authorizations=traefik@any:routeadm`: the module will automatically configure Traefik
  - `org.nethserver.images=docker.io/bitnami/dokuwiki:20200729.0.0-debian-10-r299`: use a specific upstream image
    to make sure all installations will use the same software version

See [Dokuwiki build images script](https://github.com/NethServer/ns8-dokuwiki/blob/main/build-images.sh).

### imageroot

The `imageroot` directory will be extracted to the system during the module install.
It contains 2 main paths:

- `imageroot/actions` contains all actions for the module agent, the directory will be copied to `/home/dokuwiki1/.config/actions/`
- `imageroot/systemd/user/` contains all systemd `.service` files, the directory will be copied to `/dokuwiki/.config/systemd/user/`

Usually, a module contains also a `configure-module` action to gather user input and configure the module accordingly.
The `configure-module` action should:

- [validate](../dokuwiki/imageroot/actions/configure-module/validate-input.json) user input
- [set environment variables](../dokuwiki/imageroot/actions/configure-module/20configure) and eventually expand the configuration
- [enable and start the systemd unit](../dokuwiki/imageroot/actions/configure-module/60systemd)
- [setup Traefik proxy](../dokuwiki/imageroot/actions/configure-module/30traefik), if needed

Feel free to use `dokuwiki` as scaffold module: just copy it and change it accordingly to your needs.

### User Interface

See [Modules UI](/modules/ui).

## Step 2: push the image to registry

The `build-images.sh` script will output the command to push the images to
the registry. Example:
```
buildah push ghcr.io/nethserver/dokuwiki docker://ghcr.io/nethserver/dokuwiki:latest
```

If the package should be added to a repository, make sure to specify a valid [semantic version](https://semver.org/) as tag.
Example:

```
buildah push ghcr.io/nethserver/dokuwiki docker://ghcr.io/nethserver/dokuwiki:0.0.1
```

## Step 3: start the module

First, we have to tell the cluster to instantiate a new module. The node agent will
create a new Unix user that runs the rootless module and schedule the `create-module`
action to start soon.
Just install the module straight from the image registry, without using NethServer package repository:
```
add-module ghcr.io/nethserver/dokuwiki:latest 1
```

In case of error, see `journalctl`.

When the `create-module` has been completed, start the [configuration](../dokuwiki/README.md#configure).

## Step 5: publish the package

If you want to publish the package, first make sure the package as a tag which is a valid semantic version.
Then create a new Pull Request (PR) to [ns8-repomd](https://github.com/NethServer/ns8-repomd/) repository.
The PR should contain:
- a new [directory](https://github.com/NethServer/ns8-repomd/tree/main/dokuwiki) with the name of the module
- the directory should contain:
  - a file named [`metadata.json`](thttps://github.com/NethServer/ns8-repomd/blob/main/dokuwiki/metadata.json) which contains all information regarding the module itself
  - a [logo](https://github.com/NethServer/ns8-repomd/blob/main/dokuwiki/logo.png)
  - a [screenshots](https://github.com/NethServer/ns8-repomd/tree/main/dokuwiki/screenshots) directory with one ore more screenshots (optional)

When the PR has been merged, [repository metadata](https://github.com/NethServer/ns8-repomd/tree/repomd) will be automatically updated accordingly.

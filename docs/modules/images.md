---
layout: default
title: Images
nav_order: 3
parent: Modules
---

# Images

* TOC
{:toc}

## General

NS8 modules are distributed as packages that are implemented as Podman
container images. A package is pulled from a remote software repository,
implemented by a registry service, like [GitHub
Packages](https://ghcr.io). Each package can be seen much like a *tar
archive*.

When installing a rootless module, the corresponding image is extracted
inside the home of the module Unix user, i.e. `/home/<module_id>`.
Rootfull modules are extracted under `/var/lib/nethserver/<module_id>`.

In both cases the module UI is extracted to
`/var/lib/nethserver/cluster/ui/apps/<module_id>`.

## Source tree

The sources repository of a module can be structured as follow:

- `imageroot/`: it contains module scripts and configuration. Everything inside this directory is copied under the module installation directory. Common subdirs include:
  * `systemd/user`: where Systemd units are stored.
  * `actions/`: each subdirectory implements an *action*.
  * `bin/`: it contains additional binaries for the module. It is added to PATH in the agent environment.
  * `etc/`: to store additional configuration for Backup and Restore
  * `pypkg/`: path for module Python packages, added to `PYTHONPATH`
- `ui/`: it contains all UI source code of the module
- `build-images.sh`: a script to manually build one or more images of the module and push them inside the image registry.
- `README.md`: a [Markdown](https://guides.github.com/features/mastering-markdown/) file describing the module purpose and implementation

## Image labels

Module images can use a list of well-known labels to configure the system:

- `org.nethserver.tcp-ports-demand`: see [Port allocation](../port_allocation)
- `org.nethserver.images`: see [Service images](#service-images)
- `org.nethserver.rootfull`: can be `0` or `1`, if set to `0` the module will run podman in rootless mode,
  if set to `1` the module will run podman in rootfull mode. See [Rootless vs Rootfull](../rootless_rootfull)
- `org.nethserver.authorizations`: authorize the module to run actions of other modules. For instance `traefik@any:routeadm` allows setting Traefik routes.
  See [Roles and authorizations]({{site.baseurl}}/core/agents/#roles-and-authorizations).
- `org.nethserver.flags`: a space-separated list of well-known flags. Currently available flags are:
  - `core_module`: if present, it marks a module as core
  - `account_provider`: if present, it marks a module as an account provider
  - `no_data_backup`: if present, the modules will need no data backup
  - `rootless`: if present, the module is rootless (calculated from `org.nethserver.rootfull` label)
  - `rootfull`: if present, the module is rootfull (calculated from `org.nethserver.rootfull` label)

Labels are set by `build-images.sh`, when the images are built.

## Service images

Most modules run software from additional Podman images. The
`org.nethserver.images` takes a space-separated list of image URLs that
will be automatically downloaded by the `create-module` base action.

Information about the downloaded images are stored in the agent environment.

Environment variables names are set as follow:
- one variable for each image
- variable name is the uppercase value of the image name
- symbols are mapped to `_` (underscore)
- `_IMAGE` suffix is appended

Example:
- `docker.io/library/mysql:10.3-alpine` becomes `MYSQL_IMAGE=docker.io/library/mysql:10.3-alpine`



---
layout: default
title: Images
nav_order: 30
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
- `org.nethserver.max-per-node`: maximum number of module instances installed on the same node
- `org.nethserver.min-from`: the image can be used to install a new
  application instance, or to update an existing instance provided it has
  a version greater than or equal to the label value. E.g. if the image
  label `org.nethserver.min-from` has value `2.0.0`, an existing instance
  with version `1.3.0` cannot be updated with it.
- `org.nethserver.min-core`: the image can be used to install a new
  application instance, or update an existing one, if the core version is
  greater than or equal to the label value. E.g. if the image label
  `org.nethserver.min-core` has value `2.7.0` it cannot be installed if
  the leader node running core has version `2.6.2`.

Labels are set by `build-images.sh`, when the images are built.

## Service images

Most modules run software from additional Podman images. The
`org.nethserver.images` takes a space-separated list of image URLs that
will be automatically downloaded by the `create-module` base action.

Information about the downloaded images are stored in the agent
environment, so they can be referenced in unit `.service` files and action
scripts.

Environment variables names are set as follow:
- one variable for each image
- variable name is the uppercase value of the image name
- symbols are mapped to `_` (underscore)
- if the image name begins with a digit, a `I` is prepended
- `_IMAGE` suffix is appended

Examples:
- `docker.io/library/mysql:10.3-alpine` becomes `MYSQL_IMAGE=docker.io/library/mysql:10.3-alpine`
- `quay.io/prometheus/node-exporter:v1.5.0` becomes `NODE_EXPORTER_IMAGE=quay.io/prometheus/node-exporter:v1.5.0`
- `docker.io/2fauth/2fauth:5.2.0` becomes `I2FAUTH_IMAGE=docker.io/2fauth/2fauth:5.2.0`

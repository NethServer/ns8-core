---
layout: default
title: New module tutorial
nav_order: 120
parent: Modules
---

# How to create a rootless module

This document will guide you, a developer, on how to create a new rootless module.
The module will be a simple NethServer package containing an HTTP echo server.

Before moving further, please take a look to [how NethServer is designed]({{site.baseurl}}/design).
See also [core]({{site.baseurl}}/core) and [modules]({{site.baseurl}}/modules).

* TOC
{:toc}

## Step 0: requirements

This tutorial will use some GitHub features to simplify the module setup
process. While you do not strictly need a GitHub account to develop a
module, you will need one if you want to follow this guide step-by-step.

Before proceeding, please [log in to GitHub](https://github.com/login) or
[create a new account](https://github.com/signup).

On your machine you will also need [GIT](https://git-scm.com/) and a text editor.

## Step 1: clone the kickstart repository

To develop a new module access the [ns8-kickstart](https://github.com/NethServer/ns8-kickstart) GitHub
repository and click on [`Use this template`](https://github.com/NethServer/ns8-kickstart/generate) to
generate a new repository.

1. Name your repo with `ns8-` prefix (e.g. `ns8-mymodule`). 
   Do not end your module name with a number, like ~~`ns8-baaad2`~~!

1. Clone the repository and follow the instructions inside the README

## Step 2: edit the code

It's now time to modify the code and implement your module.

Start easy and just edit the `README.md` file, by replacing this section with your module
description.

If you just want to see your package built and published go to [Step 3](#step-3-commit-and-push).
Otherwise, read about the module structure [below](#module-structure).

### Module structure

Each module is distributed using a container image.

A module is usually composed by the following main parts:
- `imageroot` directory, contains all scripts to configure the module
- `ui` directory, contains all user interface files
- `build-images.sh`, a script to manually build one or more images of the module and eventually push them inside the image registry
- `README.md` describes module implementation and usage

The module is something like an RPM or DEB file: it contains the configuration of the application that will run inside NS8.
During the installation process the system will download the real application (eg. Wordpress, Samba) and configure it
using the [actions]({{site.baseurl}}/modules/agent/) from the module.

#### imageroot

The `imageroot` directory will be extracted to the system during the module install.
It contains 2 main paths:

- `actions` contains all actions for the module agent, the directory will be copied to `/home/mymodule1/.config/actions/`
- `systemd/user/` contains all systemd `.service` files, the directory will be copied to `/home/mymodule1/.config/systemd/user/`

Usually, a module contains also a `configure-module` action to gather user input and configure the module accordingly.
The `configure-module` action should:

- validate user input
- set environment variables and eventually expand the configuration
- enable and start the systemd unit(s)
- setup Traefik proxy routes, if needed

Make sure that all parameters used by `configure-module` are saved inside the environment.
Such parameters could than be retrieved by a `get-configuration` action, used to display the configuration inside the UI.

#### build-images.sh

The script creates an empty container image and publish it to the registry on user request.

The main parts to look for are:
- `repobase`: the URL of the remote image registry
- `reponame`: the name of the module
- [labels]({{site.baseurl}}/modules/images/#image-labels) for basic module setup

#### User Interface

No matter on which node the module will run, the web user interface will be automatically imported inside the leader node
and it will be available from the admin URL `https://your.server.fqdn/cluster-admin`.

See how to [setup and develop the UI]({{site.baseurl}}/ui/modules#modules-user-interface).

Also make sure to edit [`imageroot/ui/public/metadata.json`]({{site.baseurl}}/modules/metadata) and replace the logo at `imageroot/ui/src/assets/module_default_logo.png`.

## Step 3: commit and push

Commit your local changes, than push it to GitHub repository.

After the push, GitHub will build the module as a container image and publish it to GitHub registry.
The module will be available under your namespace like `https://github.com/myuser/ns8-mymodule/pkgs/container/mymodule`.
Published image will contain the name of the branch, or tag, as version like `mymodule:main`, `mymodule:1.0.0`.
The `latest` version will always point to the latest build from the main branch, like `mymodule:latest`.

Another workflow will also build the API documentation from JSON schema and publish the result inside the `apidoc-mybranch`
branch where `mybranch` is the current pushed git branch or tag. Default is `apidoc-main`.

You're ready to [install the module](#step-4-try-the-module), unless you want to know more about [how to manually build and publish](#manual-builds)
your module image.

### Manual builds

You can build your modules locally using any distribution which supports `podman` and `buildah`.
To build the image, just execute `bash build-images.sh` script.

Before pushing images to the registry, you must configure the authentication.
Create a [GitHub Personal Access Token (PAT)](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
for the **ghcr.io** registry (for read-only access `read:packages private` scope should be enough) then run the following command, specifying
your GitHub user name and providing the generated PAT as password:
```
buildah login ghcr.io
```

The `build-images.sh` script will output the command to push the images to
the registry. Example:
```
buildah push ghcr.io/myuser/mymodule docker://ghcr.io/myuser/mymodule:latest
```

If the package should be available inside the [NethServer software repository](#step-5-publish-to-ns8-software-repository), make sure to specify a valid [semantic version](https://semver.org/) as tag.

## Step 4: try the module

If you didn't it already, it's now time to [install NS8]({{site.baseurl}}/quickstart).

You are going an SSH connection to the cluster leader to follow below steps.

First, we have to tell the cluster leader to instantiate a new module. The node agent will
create a new Unix user that runs the rootless module and schedule the `create-module`
action to start soon.

Just install the module straight from the image registry, without using NethServer package repository:
```
add-module ghcr.io/myuser/mymodule:latest 1
```

In case of error, see `journalctl`.

You can now access the module configuration at `https://your.server.fqdn/cluster-admin` or navigate directly
to the application at `https://your.server.fqdn/mymodule`.

Repeat steps from 2 to 4 until you're module is ready for prime time!

## Step 5: publish to an NS8 software repository

This step is optional. First, ensure the image is properly tagged in your
container image registry. The tag must adhere to our [Semver syntax
subset](../../development_process#module-version-numbering-rules).

- Refer to [Software repositories](../../core/software_repositories/) to
  understand how an application is listed in the NS8 Software Center and
  how to set up a personal application repository.

- To publish the application in NethForge, refer to the [Certification
  process](../certification).

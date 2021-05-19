# How to create a rootless module

Each module is distributed using a container image.
This document will guide you, a developer, on how to create a new rootless module.
The module will be a simple web application named mymodule.

## Step 0: install NS8

To develop a new module first install NS8 by following the Installation page.
Then, make sure `buildah` is installed to the system.

For Debian:
```
apt-get install buildah
```

For Fedora:
```
dnf install buildah
```

## Step 1: add to list catalog

Each module can be installed using the `install-module` script.
The install script searches for the given module inside `/var/lib/nethserver/cluster/state/images-catalog.txt` and save them
inside Redis.

When creating a new module, save its metadata directly to Redis:
```
echo 'HSET image/<key> rootfull <0|1> url <URL:tag> name "<name>"' >> /var/lib/nethserver/cluster/state/images-catalog.txt
```
Where:
- `key` is the unique module name
- `rootfull` can be: `0` if the module should run in rootless moode, `1` if the module should run in rootfull mode. If unsure, set to `0`
- `url`: the URL of the image module inside a container registry with a tag. Usually the tag is set to `latest`
- `name`: a name or a description for the module

Example:
```
redis-cli HSET image/mymodule rootfull 0 url ghcr.io/nethserver/mymodule:latest name "MyModule"
```

## Step 2: module structure

A module is usually composed by 2 main parts:
- `imageroot` directory, contains all scripts to configure the module
- `build-image.sh`, a script to generate the container image and push it to the image registry

### build-image.sh

The script creates an empty container image and publish it to the registry on user request.

The main parts to look for are:
- `repobase`: the URL of the remote image registry
- `reponame`: the name of the module
- `org.nethserver/tcp_ports_demand` label: describe how many ports the module needs to bind. The label takes an integer number as value, like `org.nethserver/tcp_ports_demand=3`

Example of `build-image.sh` file:
```
#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

reponame="mymodule"
container=$(buildah from scratch)

buildah add "${container}" imageroot /
buildah config --entrypoint=/ --label="org.nethserver/tcp_ports_demand=1" "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")
printf "%s\n" "${messages[@]}"
```

### imageroot

The `imageroot` directory will be extracted to the system during the module install.
It contains 2 main paths:

- `imageroot/actions` contains all actions for the module agent, the directory will be copied to `/home/mymodule1/.config/actions/`
- `imageroot/systemd/user/` contains all systemd `.service` files, the directory will be copied to `/home/mymodule1/.config/systemd/user/`

Every module **must** define at least a `create-module` action.
The `create-module` should contain a step to pull the image.

Something like `imageroot/actions/create-module/10pull`:
```
#!/bin/bash

set -e

# Redirect any output to the journal (stderr)
exec 1>&2

IMAGE=docker.io/bitnami/dokuwiki:latest

# Talk with agent using file descriptor:
# - save image name to environment

echo "set-env IMAGE ${IMAGE}" >&${AGENT_COMFD}

podman pull ${IMAGE}
```

Usually, a module contains also a `configure-module` action to gather user input and configure the module accordingly.
The `configure-module` action should:

- validate user input
- set environment variables and eventually expand the configuration
- enable and start the systemd unit


Feel free to use `dokuwiki` as scaffold module: just copy it and change it accordingly to your needs.

#### Port allocation

A module can require one or more TCP ports allocation with the
"org.nethserver/tcp_ports_demand" image label. For instance
```
 org.nethserver/tcp_ports_demand=3
```

It means the module receives a range of three TCP port numbers in the
usual module environment, where the following variables are set

- `TCP_PORT=1001` (always set to the range first port)
- `TCP_PORTS_RANGE=1001-1003` (the full range, if more than one port is allocated)
- `TCP_PORTS=1001,1002,1003` (comma-separated list of port numbers, if no more than 8 ports are allocated)

This command adds the label to the image:
```
buildah config --label="org.nethserver/tcp_ports_demand=3" "${container}"
```

## Step 3: start the module

First, we have to tell to the cluster agent to istantiate a new module. The cluster agent will
allocate a new unix user and call the `create-module`:
```
add-module mymodule 1
```

The output will look like:
```
{"rootfull": false, "mid": "mymodule1"}
```

Then, wait for the `create-module` to complete. You can inspect `journalctl` to monitor the task progess.
The task will have a name like `module/mymodule1/e6707401-44b7-4283-8fe9-f0aa189d3a23`.
Example:
```
journalctl | grep module/mymodule1
```

You can also inspect the task:
- read the exit code: `redis-cli get module/mymodule11/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/exit_code`
- read the standard output: `redis-cli get module/mymodule1/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/output`
- read the standard error: `redis-cli get module/mymodule1/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/error`


When the `create-module` has been completed, fire the `configure-module` action:
```
redis-cli LPUSH module/mymodule1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "mymodule 8181\n"}'
```

Where `data` contains the name of the module `mymodule` and the `port` where podman should bind.


## Step 4: add traefik routes

The web application is listening only to 127.0.0.1. To make it reachable from other network hosts,
a new traefik route is required.

Add host based route:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": "mymodule1 http://127.0.0.1:8182 mymodule.myhost.org\n"}'
```

The data field contains 3 parameters:
- the service name
- the bind URL, set the port accordingly using the `TCP_PORT` variable from `/home/mymodule11/.config/state/environment`
- a fully qualified host name representing the virtualhost for the module


To remove a host based route:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "mymodule1\n"}'
```

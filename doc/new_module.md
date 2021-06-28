# How to create a rootless module

Each module is distributed using a container image.
This document will guide you, a developer, on how to create a new rootless module.
The module will be a simple web application named `mymodule`.

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
The install script reads modules metadata from `/var/lib/nethserver/cluster/state/images-catalog.txt` and save them
inside Redis.

When creating a new module, save its metadata directly to Redis:
```
redis-cli HSET image/mymodule url ghcr.io/nethserver/mymodule:latest name "MyModule"
```

Remember to put the same command in the `images-catalog.txt` so other people can find the new module in future installations.

The record is created by a Redis command in this generic form:
```
HSET image/<key> url <URL:tag> name "<name>"
```
Where:
- `key`: the module unique identifier, i.e. `mymodule`
- `url`: the URL of the image module inside a container registry with a tag. Usually the tag is set to `latest`, e.g. `ghcr.io/nethserver/mymodule:latest`
- `name`: the name of the module, e.g. `MyModule`


## Step 2: module structure

A module is usually composed by 2 main parts:
- `imageroot` directory, contains all scripts to configure the module
- `build-image.sh`, a script to generate the container image and push it to the image registry

### build-image.sh

The script creates an empty container image and publish it to the registry on user request.

The main parts to look for are:
- `repobase`: the URL of the remote image registry
- `reponame`: the name of the module
- `org.nethserver.tcp-ports-demand` label: describe how many ports the module needs to bind. The label takes an integer number as value, like `org.nethserver.tcp-ports-demand=3`

Example of `build-image.sh` file with output for GitHub Actions:
```
#!/bin/bash

set -e
images=()
repobase="ghcr.io/nethserver"

reponame="mymodule"
container=$(buildah from scratch)

buildah add "${container}" imageroot /
buildah config --entrypoint=/ --label="org.nethserver.tcp-ports-demand=1" "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
images+=("${repobase}/${reponame}")

if [[ -n "${CI}" ]]; then
    # Set output value for Github Actions
    printf "::set-output name=images::%s\n" "${images[*]}"
else
    printf "Publish the images with:\n\n"
    for image in "${images[@]}"; do printf "  buildah push %s docker://%s:latest\n" "${image}" "${image}" ; done
    printf "\n"
fi
```

### imageroot

The `imageroot` directory will be extracted to the system during the module install.
It contains 2 main paths:

- `imageroot/actions` contains all actions for the module agent, the directory will be copied to `/home/mymodule1/.config/actions/`
- `imageroot/systemd/user/` contains all systemd `.service` files, the directory will be copied to `/home/mymodule1/.config/systemd/user/`

Every module **must** define at least a `create-module` action.
The `create-module` should contain a step to pull the web service image.

Something like `imageroot/actions/create-module/10pull`:
```
#!/bin/bash

set -e

# Redirect any output to the journal (stderr)
exec 1>&2

SVC_IMAGE=docker.io/bitnami/dokuwiki:latest

# Talk with agent using file descriptor:
# - save image name to environment

echo "set-env SVC_IMAGE ${SVC_IMAGE}" >&${AGENT_COMFD}

echo "I got the following ports: $TCP_PORTS"

podman pull ${SVC_IMAGE}
```

Usually, a module contains also a `configure-module` action to gather user input and configure the module accordingly.
The `configure-module` action should:

- validate user input
- set environment variables and eventually expand the configuration
- enable and start the systemd unit


Feel free to use `dokuwiki` as scaffold module: just copy it and change it accordingly to your needs.

#### Port allocation

A module can require one or more TCP ports allocation with the
"org.nethserver.tcp-ports-demand" image label. For instance
```
 org.nethserver.tcp-ports-demand=3
```

It means the module receives a range of three TCP port numbers in the
usual module environment, where the following variables are set

- `TCP_PORT=1001` (always set to the range first port)
- `TCP_PORTS_RANGE=1001-1003` (the full range, if more than one port is allocated)
- `TCP_PORTS=1001,1002,1003` (comma-separated list of port numbers, if no more than 8 ports are allocated)

This command adds the label to the image:
```
buildah config --label="org.nethserver.tcp-ports-demand=3" "${container}"
```

## Step 3: start the module

First, we have to tell the cluster to instantiate a new module. The node agent will
create a new Unix user that runs the rootless module and schedule the `create-module`
action to start soon:
```
add-module mymodule 1
```

> The above command means: "add a new instance of module *mymodule* to node *1*" 

The output will look like:
```
{"rootfull": false, "mid": "mymodule1"}
```

> Note: `mid` is an abbreviation for *module-identifier*.

Then, wait for the `create-module` to complete. You can inspect `journalctl` to monitor the task progess.
The task will store its result under a random key prefix, e.g. `module/mymodule1/task/e6707401-44b7-4283-8fe9-f0aa189d3a23`.
Example:
```
journalctl | grep module/mymodule1
```

You can also inspect the task results:
- read the exit code: `redis-cli get module/mymodule11/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/exit_code`
- read the standard output: `redis-cli get module/mymodule1/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/output`
- read the standard error: `redis-cli get module/mymodule1/task/e6707401-44b7-4283-8fe9-f0aa189d3a23/error`


When the `create-module` has been completed, fire the `configure-module` action:
```
redis-cli LPUSH module/mymodule1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "some input\n"}'
```

This action should use the environment variable `$TCP_PORT` and the input data to configure the web application.

The environment is stored also in a file, e.g. `/home/mymodule1/.config/state/environment`. It is safe to load it
in Systemd unit files with the `EnvironmentFile=` directive.

Every action executable step runs with `~/.config/state` as working directory.

## Step 4: add traefik routes

The web application is listening only to 127.0.0.1.  To make it reachable from other network hosts,
a new traefik route is required.

To get the TCP port number among the environment variables used by the `create-module` run:

    redis-cli HGETALL module/mymodule1/environment


Assuming the port is 20001, add a hostname-based HTTP route:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "mymodule1", "url": "http://127.0.0.1:20001", "host": "mymodule.myhost.org", "lets_encrypt": true, "http2https": false}}'
```

The data field contains 3 parameters:
- the service name
- the bind URL, set the port accordingly using the `TCP_PORT` variable from `/home/mymodule1/.config/state/environment`
- a fully qualified host name representing the virtualhost for the module


To remove a host based route:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "mymodule1\n"}'
```

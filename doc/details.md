# Technical details

Before reading this chapter, make sure to take a look at [Design & Architecture](design.md).

## General

NS8 relies on default system configuration as much as possible.
All logs are available using `journalctl` and services can be inspected using `systemctl` command.

While agents are implemented using [Go](https://golang.org/), most of the actions are implemented with [Python](https://www.python.org/).
Python scripts run on a special environment installed inside `/usr/local/nethserver/agent/python`.
The environment is automatically enabled upon root login by executing `source /etc/profile`.

Most NethServer files are saved inside `/var/lib/nethserver/` and `/usr/local/nethserver/`.

### Database: Redis

NethServer configuration is saved inside Redis which runs as rootfull container.
To see Redis service status, use `systemctl status redis`.

Once the core has been initialized, you can access Redis with one of the following commands:

    redis-cli <<EOF
    PING
    EOF

The `redis-cli` command is faster when invoked by the root user, because it runs a command in
the existing `redis` rootfull container. For non-root users a temporary container is created and destroyed.

The `redis-cli` command attempts to connect with the higher available Redis privileges, by reading `agent.env` files.

An experimental Python-based helper script is available too and it better suits
Bash action scripts as it does not start any container.
This client is synchronous as it waits for the server response.

    redis-exec <<EOF
    PING
    EOF

An alternative, synchronous invocation relies on the `nc` command, provided by the `nmap-ncat` RPM (good for development environments):

    nc 127.0.0.1 6379 <<EOF
    PING
    EOF

For completeness here is a pure Bash invocation that does not wait the server response:

    cat >/dev/tcp/127.0.0.1/6379 <<EOF
    PING
    EOF


#### DB schema and queues

TODO:
- redis queues
- schema: https://docs.google.com/document/d/13bo5guJN588HGzP5EPe9e5heHBmPt88U9AoyV0OCs0g/edit#heading=h.ahj72tll1jvk

### Proxy and certificates: Traefik

Traefik is the proxy server used to access all HTTP/HTTPS services.
It also handles Let's Encrypt certificate requests and renewal.

See [Traefik README](../traefik/README.md) for more info.

## Packages

NethServer packages are container images which will be unpacked after download.
Each package can be seen much like a tar archive.

The installation process just installs all required dependencies,
and extract the files from the core image to system path, like `/var/lib/nethserver`.

When installing a rootless module, the package will be extracted inside the home of the module Unix user.

## API server and UI

The API server is a daemon implemented using [Go](https://golang.org).

The UI is a web application developed with [VueJs](https://vuejs.org/) and based on [Carbon Design System](https://www.carbondesignsystem.com/).

TODO:

- UI and API server paths
- UI core parts
- how to create user and roles
- UI default credentials
- auditing

## Package repository

The systems supports a virtually infinite number of repositories. Each repository contains the list of available modules
and can be enabled or disabled.
If the repository has the `testing` flag enabled, it will list also non-stable module versions.

Repository metadata are downloaded and cached inside redis for 3600 seconds.
To expire an existing repository cache, execute: `redis-cli del cluster/repository_cache/<repository_name>`

A NethServer repository must contain a file named `repodata.json` which describes the content of the repository.
The [ns8-repomd](https://github.com/NethServer/ns8-repomd/) it's the default repository and contains also
the `createrepo.sh` command which creates the metadata.

### createrepo.sh

The `createrepo.sh` command analyzes a directory passed as parameter, if no path is given, it walks the current working directory.
For each directory found, it uses the directory name to query the remote image registry with
[skopeo](https://github.com/containers/skopeo) and retrieve image tags.
If tags are valid [semantic versions](https://semver.org/) they are added to the list of available module version.

Each module that needs to be published inside the repository should have a directory named after the module itself.
The directory should contain:

- a file named `metadata.json`: it contains all module metadata like the name and the URL inside the image registry
- a file named `logo.png`: a PNG file, 256x256 pixels
- a directory named `screenshots`: it can contain one ore more image, each image must be in PNG format,
  with a resolution of 1024x768 pixels

See [dokuwiki directory](https://github.com/NethServer/ns8-repomd/tree/main/dokuwiki) for an example of a complete module
metadata directory.

Execute `createrepo.sh` to generate the `repodata.json`. The format is described [here](https://github.com/NethServer/ns8-scratchpad/blob/main/core/imageroot/var/lib/nethserver/cluster/repodata-schema.json).

## Module architecture

A NethServer module carries all scripts to configure an upstream container image
and start a podman process.

Modules are installed by the core `node-agent` service, when it receives a `add-module` event.

Every time a module instance is spawn, the instance will be named as the module itself
followed by a progressive number starting from 1.
Given a module named `myapp`, instances will be named `myapp1`, `myapp2`, etc.
The Unix user account has session lingering enabled: it automatically starts a persistent Systemd user manager instance.

A module provides a bunch of event handlers and Systemd unit definitions.
An event is handled by one or more *action* scripts, stored under `$HOME/.config/actions`. 
The local Sysadmin can extend and/or override them by putting their action scripts under `$HOME/module-events`.
Module-local systemd units are installed under `$HOME/.config/systemd/user`. System-wide units are installed under `/etc/systemd/user`.

The `module-agent.service` Systemd unit executes the event handlers. The agent daemon runs in the Python virtual
environment installed in `/usr/local/share/agent/`: action scripts inherit the same environment. Additional binaries
can be installed under `/usr/local/share/agent/bin/`.

To inspect and modify a rootless module start a SSH session. SSH is preferred to `su - <user>` because the latter
does not properly initialize the Systemd session environment.
Example: check if Traefik is running:
```
ssh traefik1@localhost
systemctl --user status
```

Modules can be managed using these commands:

- `add-module <module> <node_id>`: install a module on the given node with ID `node_id`; search for `module` inside enabled repositories and install
  latest available version. If `module` is a image registry URL, just install the module straight from it; this method can be used to install
  customized images.
- `remove-module [--no-preserve] <module>`: remove an installed module; if `--no-preserve` is given, erase also all module data

### Internals

NethServer defines some conventions that should ease the creation of new modules.

#### Filesystem structure

Each module has should be self-contained inside a directory named as the module itself.
Such directory should contain:

- `imageroot` directory: it contains all module scripts and configuration, it's composed by two sub-directories:
  - `systemd`: everything inside this directory is copied under `/home/<instance>/.config/systemd/` directory.
  It should contains a sub-directory named `user` where systemd units are stored.
  - `actions`: a list of directories, each directory implements an `action`
- `ui` directory: it contains all UI source code of the module
- `build-image.sh`: a script to manually build the image of the module and push it inside the image registry.
- `README.md`: a [Markdown](https://guides.github.com/features/mastering-markdown/) file describing the module purpose and implementation

#### Image labels

Module images can use a list of well-known labels to configure the system:

- `org.nethserver.tcp-ports-demand`: see [Port allocation](#port-allocation)
- `org.nethserver.images`: see [Image download](#image-download)
- `org.nethserver.rootfull`: can be `0` or `1`, if set to `0` the module will run podman in rootless mode,
  if set to `1` the module will run podman in rootfull mode
- `org.nethserver.authorizations`: see agent API below 

##### Image download

Most applications will use one or more upstream images containing the software to be executed.
The `org.nethserver.images` takes a space-separated list of image URLs that will be automatically
downloaded by the node agent.

Environment variables names are set as follow:
- one variable for each image
- variable name is the uppercase value of the image name
- symbols are mapped to `_` (underscore)
- `_IMAGE` suffix is appended

Example:
- `docker.io/library/mysql:10.3-alpine` becomes `MYSQL_IMAGE=docker.io/library/mysql:10.3-alpine`

##### Port allocation

Many web application modules need to listen on a TCP port which will be than exposed using Traefik.
Such modules can set the `org.nethserver.tcp-ports-demand` which takes an integer number as value.
Example:
```
org.nethserver.tcp-ports-demand=1
```

The randomly-allocated TCP port will be available inside the `TCP_PORT` environment variable and it will be 
available to all step scripts and inside systemd units.
The available environment variables will be:
- `TCP_PORT`: it is always present and it contains always the first port, i.e. `1000`
- `TCP_PORTS_RANGE`: only if value is greater than 1, it contains the list of ports in range format,
  i.e `1000-1002`
- `TCP_PORTS`: only if value is greater than 1 and less or equal than 8, it contains a comma separated list of
  ports like, i.e. `1000,10001,10002`

Currently last allocated port is saved inside Redis at `'node/<node_id>/tcp_ports_sequence`.

#### Systemd units

A module should always have at least one systemd unit to keep the service running.
The systemd unit should:
- make sure the environment is always set: `EnvironmentFile=%S/state/environment`
- make sure the working directory is always set: `WorkingDirectory=%S/state`
- start and stop podman processes
- ensure restart at least on failure
- bind to the network with possible lowest privilege
- enable the service on boot

See [Dokuwiki unit](../dokuwiki/imageroot/systemd/user/dokuwiki.service) as an example.

##### Network binding

Because NethServer has no firewall, network access to applications should carefully restricted
using the correct IP binding.

As a general rule, any module which doesn't require a well-known port, should request a random port
using `org.nethserver.tcp-ports-demand` label and bind to that port only on 127.0.0.1,
i.e. `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:${TCP_PORT}:8080 ...`

Modules using well-known port, should bind to 127.0.0.1 and to all IPs where the service
must be exposed, like VPN IPs, i.e `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:19999:19999 $EXTRA_LISTEN ...`,
where `EXTRA_LISTEN` could be the bind to the VPN `-p 10.5.4.1:19999:19999`.

#### Actions

An `action` is a directory containing one ore more executable scripts called `steps`.
Each action is invoked by a cluster/node/module agent when an event is pushed inside the Redis queue.
Steps are executed in lexicographic order, if one step fails, the action execution will be stopped.
Actions must always take a valid JSON as input from standard input and always generate a valid JSON as output
inside the standard output.
There are no conventions on the format of output on standard error.

See also [Actions execution model](../core/agent/README.md#actions-execution-model).

Example of action invocation at cluster level:
```
UUID=$(uuidgen); redis-cli LPUSH cluster/tasks '{"id": "'$UUID'", "action": "list-modules", "data": null}'
```

To retrieve the output:
```
redis-cli get cluster/task/$UUID/output
```

The following actions are automatically available for all modules:

- `create-module`: executed upon installation, it automatically downloads the module image and all images
  listed inside `org.nethserver.images` label.
- `remove-module`: executed on module removal, modules can add here scripts to remove configuration executed on
  other modules, like Traefik routes

Each module can add extra step to the above actions by creating a directory with the of the action itself inside the module implementation,
e.g. `imageroot/actions/myapp/create-module/20mystep`.

If a module needs user input to have a valid default configuration, it should declare a `configure-module` action.

##### Environment variables

Environment variables are saved inside Redis at 'module/<module_id>/environment`.

Each action has access to all environment variables:
1. added by `add-module`
2. defined from agent at runtime
3. inherited by the agent

Main environment variables are:
- `MODULE_ID`: the ID of the instance, i.e. `myapp1`
- `NODE_ID`: the node ID where the action is running, i.e. `1`
- `<name>_IMAGE`: see [Image download](#image-download)
- `TCP_PORT`, `TCP_PORTS`, `TCP_PORTS_RANGE`: see [Port allocation](#port-allocation)

Other available variables:
- `AGENT_COMFD`: the file descriptor to talk to the agent via [action commands](#action-commands)
- `AGENT_TASK_ID`: the unique ID of current task, i.e. `c0b8b976-9444-42d5-a40b-142a6a483a84`
- `AGENT_TASK_ACTION`: the name of executed action, i.e.: `create-module`
- `AGENT_INSTALL_DIR`
- `AGENT_STATE_DIR`
- `AGENT_ID`
- `IMAGE_ID`
- `IMAGE_URL`
- `IMAGE_DIGEST`
- `IMAGE_REPODIGEST`

Environment variables added by the `add-module` are also saved inside `/home/<module_id>/.config/state/environment` file,
i.e `/home/myapp1/.config/state/environment`.

##### Action commands

Actions can talk to their own agent using a simple protocol.
An action can set/unset environment variable, dump the environment to a well known file, etc.

See all [available action commands](../core/agent/README.md#action-commands).

Python example:
```python
import agent

agent.set_env('MYVAR', 'value')
agent.dump_env()
```

##### Database access

Each module has granted full read-only access to the Redis database.
To enable write access, the module should authenticate itself.
Everything is already available inside the Python library.

Access Redis in read-only:
```python
import agent

rdb = agent.redis_connect(privileged=False)
repo = rdb.hgetall('module/myapp1/environment')
```

Access Redis in read-write mode:
```python
import agent

rdb = agent.redis_connect(privileged=True)
repo = rdb.hset('module/myapp1/environment', {'myvar': 'myvalue'})
```

##### Exit codes

Step scripts must exit with status `0` on success.
Any other exit code will cause the action failure and next steps will not be executed.

See also [Action status and outcomes](/core/agent/README.md#action-status-and-outcomes).

##### Validation

Actions should always have at least one validation step.
Formal validation of input and output can be done using [JSON validation](#JSON-validation).
Semantic validation should be implemented in a step executed early in the action, like `10validation`.
If the validation fails, the step must return a non-zero exit code.

###### JSON validation

Input and output can be validated using [JSON Schema](https://json-schema.org/).
JSON validation is implement inside the core and available to all actions.
To validate the input just drop a `validate-input.json` file inside the action directory,
to validate the output just create a `validate-output.json` file.

As general rule, schema files are loaded from special paths:
- `${AGENT_INSTALL_DIR}/validator-definitions.json`: this is always picked up by the agent from all actions. It can contain the common data format definitions
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-input.json`: schema applied to validate the input data
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-output.json`: schema applied to validate the output data

#### Agent API

TODO

#### Roles and authorizations

TODO

#### Automatic builds

Every time a commit is pushed to the repository, a new build can
be automatically started.
To enable automatic builds, create a yaml file like `.github/workflows/<module>.yaml`.

Use [dokuwki Github workflow](../.github/workflows/dokuwiki.yaml) as a template, make sure to replace all occurrences of `dokuwiki`
with the name of the new module.

##### API doc generation

TODO

---
Next: [How to create a new module](new_module.md)

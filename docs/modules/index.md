---
layout: default
title: Modules
nav_order: 5
has_children: true
---

# Modules: implementation details

Every time a module instance is added to the cluster, the new instance
is named as the module itself followed by a progressive number
starting from 1. Given a module named `myapp`, instances will be named
`myapp1`, `myapp2`, etc.

Modules can be managed using these commands:

- `add-module <module> <node_id>`: install a module on the given node with ID `node_id`; search for `module` inside enabled repositories and install
  latest available version. If `module` is a image registry URL, just install the module straight from it; this method can be used to install
  customized images.
- `remove-module [--no-preserve] <module>`: remove an installed module; if `--no-preserve` is given, erase also all module data

## Images

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

## Rootless vs Rootfull

The main difference between *rootless* and *rootfull* modules, as
suggested by the adjective, is the Unix user running the module processes
and its system privileges.

**Rootfull** modules run as `root` (EUID 0).

**Rootless** modules run as a normal Unix user. The Unix user account is
created by the node agent when the module instance is added. It has
session lingering enabled: it automatically starts a persistent Systemd
user manager instance at system boot.

The two types of modules have a similar filesystem structure. Rootless
modules are installed to `/home/<module_id>/.config`, whilst rootfull are
installed to `/var/lib/nethserver/<module_id>`.

Rootless modules also have Systemd user units installed under
`~/.config/systemd/user`. Recall that system-wide user units are installed
by the core under `/etc/systemd/user`.

To inspect and modify a rootless module start a SSH session. SSH is
preferred to `su - <user>` because the latter does not properly initialize
the Systemd session environment. For instance, to check if Traefik is running:

    ssh traefik1@localhost
    systemctl --user status

## Sources structure

The sources repository of a module can be structured as follow:

- `imageroot/`: it contains module scripts and configuration. Everything inside this directory is copied under the module installion directory. Common subdirs include:
  * `systemd/user`: where Systemd units are stored.
  * `actions/`: each subdirectory implements an *action*.
  * `bin/`: it contains additional binaries for the module. It is added to PATH in the agent environment.
- `ui/`: it contains all UI source code of the module
- `build-image.sh`: a script to manually build the image of the module and push it inside the image registry.
- `README.md`: a [Markdown](https://guides.github.com/features/mastering-markdown/) file describing the module purpose and implementation

## Image labels

Module images can use a list of well-known labels to configure the system:

- `org.nethserver.tcp-ports-demand`: see [Port allocation](#port-allocation)
- `org.nethserver.images`: see [Service images](#service-images)
- `org.nethserver.rootfull`: can be `0` or `1`, if set to `0` the module will run podman in rootless mode,
  if set to `1` the module will run podman in rootfull mode. See [Rootless vs Rootfull](#rootless-vs-rootfull)
- `org.nethserver.authorizations`: see [Roles and authorizations](#roles-and-authorizations)

Labels are set by `build-image.sh`, when the image is built.

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

## Port allocation

Many web application modules need a TCP port to run a backend exposed by Traefik.
Such modules can set the `org.nethserver.tcp-ports-demand` which takes an integer number as value.
Example:
```
org.nethserver.tcp-ports-demand=1
```

The randomly-allocated TCP port number will be available inside the `TCP_PORT` environment variable and it will be
available to all step scripts and inside systemd units.
The available environment variables will be:
- `TCP_PORT`: it is always present and it contains always the first port, i.e. `1000`
- `TCP_PORTS_RANGE`: only if value is greater than 1, it contains the list of ports in range format,
  i.e `1000-1002`
- `TCP_PORTS`: only if value is greater than 1 and less or equal than 8, it contains a comma separated list of
  ports like, i.e. `1000,10001,10002`

Currently last allocated port is saved inside Redis at `'node/<node_id>/tcp_ports_sequence`.

## Systemd units

A module should have at least one Systemd unit to keep the service running.
The Systemd unit should:
- enable the service on boot
- start and stop podman processes
- ensure restart at least on failure

The Systemd unit can:
- change the working directory to ease finding additional files: `WorkingDirectory=%S/state`
- execute ancillary commands in `Exec*` steps with the `runagent` helper: `ExecStartPre=/usr/local/bin/runagent expand-templates`
- import and use additional environment variables: `EnvironmentFile=%S/state/environment`

When `EnvironmentFile=%S/state/environment` is used in Systemd service
units, remember that the file syntax is not designed to be compatible with
Systemd. Values are stored as raw values, i.e. special characters are not
escaped or quoted:

1. Limit the use of this technique to variables with controlled values, e.g. `IMAGE_URL`.
2. Avoid storing and reading arbitrary strings, like randomly generated
   tokens, and values coming from the user's input.

See [Dokuwiki
unit](https://github.com/NethServer/ns8-scratchpad/blob/main/dokuwiki/README.md#dokuwiki)
as an example.

## Network binding

Since NS8 has no firewall, network access to applications should carefully restricted
using the correct IP binding.

As a general rule, any module which doesn't require a well-known port, should request a random port
using `org.nethserver.tcp-ports-demand` label and bind to that port only on 127.0.0.1,
i.e. `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:${TCP_PORT}:8080 ...`

Modules using well-known port, should bind to 127.0.0.1 and to all IPs where the service
must be exposed, like VPN IPs, i.e `ExecStart=/usr/bin/podman run ... -p 127.0.0.1:19999:19999 $EXTRA_LISTEN ...`,
where `EXTRA_LISTEN` could be the bind to the VPN `-p 10.5.4.1:19999:19999`.

## Base actions

Agent task processing is defined by actions. An `action` is a directory
containing one ore more executable scripts called `steps`. See [the agent
README](https://github.com/NethServer/ns8-scratchpad/blob/main/core/agent/README.md#agent)
for the complete story.

The core provides a set of base actions defined in
`/usr/local/agent/actions` and inherited by module agents.

- `create-module`: executed upon installation, it automatically downloads the module image and all images
  listed inside `org.nethserver.images` label.
- `remove-module`: executed on module removal, modules can add here scripts to remove configuration executed on
  other modules, like Traefik routes
- `get-status`: used by the UI, it returns the current module status, like the node where the application is running, used images and volumes, systemd units status.
  This action works correctly only for rootless modules.

Every agent has also a builtin `list-actions` action. This command lists the
actions provided by the agent of `module/traefik1`:

    api-cli run --agent module/traefik1 list-actions | jq

As a convention, all modules with a UI should also implement the following actions:

- `configure-module`: it validates the user input and applies the
  configuration, see [Dokuwiki
  example](https://github.com/NethServer/ns8-scratchpad/blob/main/dokuwiki/imageroot/actions/configure-module)
- `get-configuration`: it should return current configuration, the output
  should be equal or similar to configure-module input, see [Dokuwiki
  example](https://github.com/NethServer/ns8-scratchpad/blob/main/dokuwiki/imageroot/actions/get-configuration)


## Agent environment

Environment variables are copied to Redis at `module/<module_id>/environment`.

Each action has access to environment variables. They are:
1. initialized by `add-module`
2. changed by the agent at runtime
3. inherited from the agent environment

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

Environment variables added by the `add-module` and defined at runtime are
saved inside `/home/<module_id>/.config/state/environment`, e.g.
`/home/myapp1/.config/state/environment`.

To run a command within an agent environment, use the `runagent` wrapper. For instance
this command prints the agent environment of the rootfull module `promtail1`.

    runagent -m promtail1 env

## Events

Events are messages published on specific Redis channels. Modules must
subscribe to the channel they want to be notified from. Usually channel
address is under a module name, something like
`module/<module_id>/event/<event_name>`.

Events should respect the following rules:
- use past tense inside the name, like `account-provider-credentials-updated`
- accept a parameter in JSON format
- the parameter should contain minimal required info about the event
- optionally, the parameter can contain extra data which can ease the event usage

Well known events:
- `account-provider-changed`: the event should be fired by account providers to inform about configuration changes
- `account-provider-credentials-updated`: the event can be fired by account providers or LDAP proxies,
  it notifies the change of LDAP credentials (eg. `ldapservice` user) to applications

### Signaling events

To signal an event, just [PUBLISH](https://redis.io/commands/PUBLISH) a message to the channel.
All events take a JSON object as a parameter.

Example:
```
redis-cli PUBLISH module/traefik1/event/test '{"field": "value"}'
```

### Handling events

Events are handled by the `eventsgw` service.
Multiple service instances run for both rootfull and rootless modules.

The service starts only if the `eventsgw.conf` configuration file exists.
Rootfull modules must create the the file inside `/var/lib/nethserver/<module_id>/state/eventsgw.conf`,
i.e `/var/lib/nethserver/samba1/state/eventsgw.conf`.
Rootless module must create the file inside the state dir `/home/<module_id>/.config/state/eventsgw.conf`,
i.e. `/home/openldap1/.config/state/eventsgw.conf`.

The `eventsgw.conf` is a INI file with an `actions` section.
Each section has a key-value format:
- the key on the left is the Redis channel to subscribe
- the value on the right is the command or action to execute

Example:
```ini
[actions]
# Execute an action on the module when the event is received. The event data is passed to the action `data` argument.
module/traefik*/event/certificate-renew = renew-certificate
```

Modules should always use an action upon receiving an event.

Still the `eventgsw` service can also execute a shell comand instead of a full action.
Use this mode just for testing.
Example:
```ini
[commands]
# Execute a shell command when the event is received. Event data is piped to STDIN.
module/traefik1/event/certificate-update = redis-exec INCR mykey
```

## Database access

Everyone can access the Redis database with read-only privileges. The
Redis user `default` (password is not important) can be used for this
purpose.

To get write access a module must provide the Redis credentials stored in
its `agent.env` file. The complete path is `~/.config/state/agent.env`.
Write access is restricted to Redis keys and channels with prefix
`module/{module_id}/*`. The same credentials allow read access of keys
with the same prefix.

The above rules are already implemented by the Python `agent` module.

Access Redis in read-only:
```python
import agent

rdb = agent.redis_connect(privileged=False)
somehash = rdb.hgetall('cluster/somehash')
```

Access Redis in read-write mode:
```python
import agent

rdb = agent.redis_connect(privileged=True)
rdb.hset('module/myapp1/myhash', mapping={'myvar': 'myvalue'})
```

## Agent commands

During the action execution, action steps can talk to the agent using a simple protocol.
It is possible to set/unset environment variables, dump the environment to a well known file, etc.

See all [available action commands](https://github.com/NethServer/ns8-scratchpad/blob/main/core/agent/README.md#action-commands).

The protocol is implemented by the Python `agent` package.

Python example:
```python
import agent

agent.set_env('MYVAR', 'value')
agent.dump_env()
```

## Automatic builds

Every time a commit is pushed to the repository, a new build can
be automatically started.
To enable automatic builds, create a yaml file like `.github/workflows/<module>.yaml`.

Use [dokuwki Github workflow](../.github/workflows/dokuwiki.yaml) as a template, make sure to replace all occurrences of `dokuwiki`
with the name of the new module.


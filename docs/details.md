# Core details

Before reading this chapter, make sure to take a look at [Design & Architecture](design.md).

NS8 relies on default distro configurations as much as possible.
All logs are available using `journalctl` and services can be inspected using `systemctl` command.

## Filesystem locations

The core executable scripts and binaries are installed under
`/usr/local/bin` and `/usr/local/sbin`.

* The `agent` and `api-server` commands are implemented in
  [Go](https://golang.org/).
* Python scripts run within a special environment installed inside
  `/usr/local/agent/pyenv`.

Most NS8 core files reside in `/etc/nethserver`,
`/var/lib/nethserver/` and `/usr/local/agent/`. 

See `/var/lib/nethserver/node/state/coreimage.lst` for a complete list.

## Database: Redis

The Redis database runs as a Podman rootfull container. Get the Redis
service status:

    systemctl status redis

If any change occurs in the Redis database, it is persisted on the disk
within 5 seconds. The following command prints the database path:

    podman volume inspect redis-data | jq -r .[].Mountpoint

Access Redis with one of the following commands:

    redis-cli PING

* The `redis-cli` command is faster when invoked by the root user, because
  it runs a command in the existing `redis` rootfull container. For
  non-root users a temporary container is created and destroyed.

* The `redis-cli` command attempts to connect with the higher available
  Redis privileges, by reading `agent.env` files.

An Python-based helper script, `redis-exec` is available and it better
suits Bash action scripts as it does not start any container. This client
is synchronous as it waits for the server response:

    runagent redis-exec PING

Note that the `runagent` wrapper is needed because `redis-exec` is
designed to run inside an "agent environment".

An alternative, synchronous invocation relies on the `nc` command,
provided by the `nmap-ncat` RPM (good for development environments):

    nc 127.0.0.1 6379 <<EOF
    PING
    EOF

For completeness here is a pure Bash invocation that does not wait the
server response:

    cat >/dev/tcp/127.0.0.1/6379 <<EOF
    PING
    EOF

Redis is a key/value database. Some key naming rules are enforced because
access rights are based on key name patterns. Both admin UI users and
agents have their own Redis credentials for authentication and
authorization. 

Some key rule examples:

* `module/{module_id}/environment` The key type is HASH. It stores
  a copy of environment variables of the module *module_id*.

* `node/{node_id}/vpn` The key type is HASH. It stores the cluster VPN
  configuration applied by the node agent *node_id*

* `cluster/node_sequence` The key value is an integer. As the name
  suggests it is used to generate a node ID, when a new node is added to
  the cluster.

The complete Redis schema is frequently updated. It is temporarly published at the following URL:

https://docs.google.com/document/d/13bo5guJN588HGzP5EPe9e5heHBmPt88U9AoyV0OCs0g

Redis provides blocking pop operations on lists. The BRPOP command is used
by agents to wait on a queue for incoming tasks. Task queues keys look
like `module/{module_id}/tasks`.

Redis provides PUB/SUB operations on *channels*. Channels are used to
implement the observer pattern, for instance in the following situations:

1. An agent wants to notify the UI about the progress of a running task.
   The channel name in this case looks like
   `progress/module/{module_id}/task/{uuid}`.

2. An agent wants to communicate to other agents that a particular event
   occurred (e.g. a TLS certificate was automatically renewed). The
   channel name in this case could be like
   `module/{module_id}/event/{event_name}`

## Proxy and certificates: Traefik

Traefik is the proxy server used to access all HTTP/HTTPS services.
It also handles Let's Encrypt certificate requests and renewal.

See [Traefik README](https://github.com/NethServer/ns8-scratchpad/blob/main/traefik/README.md) for more info.

Let's Encrypt certificates are automatically exported to Redis upon request and renewal.
Certificates are saved under an hash named `/module/traefik<X>/certificate/<domain>` key,
i.e `/module/traefik1/certificate/server.nethserver.org`.
The certificate is saved inside the `cert` field, while the key is saved inside the `key` field.

When a certicate is exported, Traefik module fires the `certificate-update` event with
a JSON messages. The JSON message contains the following fields:

- `key`: the X509 certificate private key, PEM format encoded with base64
- `certificate`: the X509 certificate, PEM format encoded with base64
- `node`: the node id where the traefik instance is running
- `module`: the module id of the traefik instance
- `domain`: the FQDN for the certificate

Example:
```json
{
    "key": "AAa...",
    "certificate": "BBb..",
    "node": "1",
    "module": "traefik1",
    "domain": "server.nethserver.org"
}
```

The event is published under Redis channel `module/traefik<X>/event/certificate-update`.

See also [how to fire and handle events](#events).

## Users and groups: Ldapproxy

Users and groups are stored in an LDAP database, served by one **account
provider module**. Multiple modules can work together to serve the same
LDAP database as replicas of it. An LDAP database represents an account
**domain**.

A NS8 cluster can host multiple account domains from different
implementations. It is possible to configure external LDAP services as
hosted ones. Supported LDAP schemas are

1. Active Directory
2. RFC2307

A module can discover the list of available account domains with the
`cluster.ldapproxy` Python module. The following command dumps a list of
parameters required to connect with an LDAP database on cluster node 1.

    runagent python3 -mcluster.ldapproxy

Returned TCP endpoints are local (`host` is `127.0.0.1`) and do not
require TLS. The port number depends on the LDAP domain.

Those ports are held by the [Ldapproxy
module](https://github.com/NethServer/ns8-scratchpad/blob/main/ldapproxy/README.md).
It is a L4 proxy that relays the TCP connection to an LDAP backend server,
enabling TLS and handling backend failures as needed.

If the LDAP client module runs in a Podman container with a
`private` network, add the following arguments to the `podman run`
command:

    --network=slirp4netns:allow_host_loopback=true

Then replace `127.0.0.1` with the special `10.0.2.2` IP address, that is
translated by Podman back to the loopback device, 127.0.0.1 on the root
network namespace.

Python code example:

```python
from cluster.ldapproxy import Ldapproxy
lp = Ldapproxy()
domains = lp.get_domains_list()
print(domains)
domain = lp.get_domain("mydomain")
print(domain)
```

The module can add the following settings to `eventsgw.conf` to listen to
account provider changes:

```ini
[commands]
*/event/account-provider-changed = mycommand && systemctl --user reload mymodule.service
```

See also [how to fire and handle events](#events).


## API server

The API server is a daemon implemented using [Go](https://golang.org).

This component is used to send command from UI to Redis, using HTTP Rest API and Redis Pub/Sub protocol.

### API Paths
To view the API paths and documentation follow these steps:
- Access the site https://petstore.swagger.io
- In the top bar paste https://raw.githubusercontent.com/NethServer/ns8-scratchpad/swagdoc/swagger.json
- Click on "Explore"
- Browse accordions to view API structure and data

### Audit
Every request made to the server, using its APIs or WebSocket, is logged inside an audit db. The audit db is store in a file using a SQLite database schema. Each record is composed by:
- `ID`: the unique id of the record, autoincrement
- `User`: the name of the user that made the specific action
- `Action`: the name of the action made by the user
- `Data`: the payload of the action (if present)
- `Timestamp`: the time when the specific action is executed

## User Interface

NS8 user interface is a web application developed with [VueJs](https://vuejs.org/) and based on [Carbon Design System](https://www.carbondesignsystem.com/).

NS8 UI can be accessed at `https://leader_node/cluster-admin/`. Default username is `admin` and default password is `Nethesis,1234`. The web application source is structured in multiple VueJs projects:
- Core UI
- A distinct UI project for every module (e.g. Dokuwiki)
- A UI library that includes a set of reusable UI components and functions (VueJs mixins) used by core and modules UI

### Core UI

Core UI includes the following components:
- Login page
- Shell
  - Side menu
  - Top header
  - Global search
  - Notification drawer
  - App launcher
- Cluster status page
- Software center
- Cluster logs (auditing)
- Cluster settings

Core UI also includes a [Storybook](https://storybook.js.org/) to explore and test the reusable components included in the UI library.
To launch Storybook webapp:

```bash
cd core/ui
yarn storybook
```

### Modules UI

Every NS8 module follows the same UI guidelines in order to provide a uniform user experience. Almost every module has at least these standard pages:
- Status
- Settings
- Logs
- About

Status page is the landing page of the module, it should provide a dashboard displaying the current status of the module, including instance name, installation node and information about module services.

Settings page should contain a form to set and review module configuration.

Logs page should display log entries relevant to the module, in order to support troubleshooting activity.

About page should provide module meta-information such as documentation URL, source code and author information.

NS8 modules make use of UI library components and functions, e.g. by including the shared `NsButton` component in a module form, or requesting the execution of a cluster task.

### UI library

Since core and modules UI share multiple components and features, these facilities have been aggregated into NPM package `@nethserver/ns8-ui-lib`.
The reusable UI components included in the library are easily recognizable since their name starts with `Ns` prefix, e.g. `NsButton`, `NsInlineNotification`.
The library also includes a set of VueJs mixins used by core and modules to access utility functions and perform common tasks.

### Shortcuts

The core and the applications can add quickly accessible actions called shortcuts.
Shortcuts will be available from the search bar inside the UI.
Examples of shortcuts could be something like "Add a user" or "Show applicationX logs".

Shortcuts are described inside a `shortcuts.json` file which must be inside the `ui` path:
- core: `/var/lib/nethserver/cluster/ui/shortcuts.json`
- applications: `/var/lib/nethserver/cluster/ui/apps/<app>/shortcuts.json`

A `shorcuts.json` must always contain the following fields:
- `name`: the action name
- `description`: a text which describes the action
- `tags`: a list of tags to ease the research
- `path`: the UI path to open when the shortcut has been selected. The path **must** be relative.

The `list-shortcuts` API will read all `shortcuts.json` files and combine them on a single array and it will generate absolute UI paths.


## Software repositories

Module packages are pulled from software repositories. Multiple
repositories can be configured in a NS8 cluster. Each repository contains
the list of available modules
and can be enabled or disabled.
If the repository has the `testing` flag enabled, it will list also non-stable module versions.

Repository metadata are downloaded and cached inside Redis for 3600
seconds.
To expire an existing repository cache, execute: `redis-cli del cluster/repository_cache/<repository_name>`

A NS8 repository must contain a file named `repodata.json` which describes
the content of the repository. The
[ns8-repomd](https://github.com/NethServer/ns8-repomd/) is the default
repository and contains also
the `createrepo.sh` command which creates the metadata.

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

## Agents

The core `cluster` and `node` agent, and every module instance run their
own `agent` process. The process is started by the `agent.service` Systemd
user unit and the `agent@.service` Systemd system template unit provided
by the core.

See also the [Agent documentation](https://github.com/NethServer/ns8-scratchpad/blob/main/core/agent/README.md).

Similarly `eventsgw.service` and `eventsgw@.service` units are provided by
the core to run event handler services. See 

### Tasks processing

Agents wait for tasks on a Redis list. Only the `cluster` agent and the
`api-server` are allowed to push a new task for an arbitrary module agent. The
task contain two main attributes: the action name (`action`), and the
action input (`data`).

The following Redis command run the `list-actions` action on the cluster agent:

    LPUSH cluster/tasks '{"id":"1a4a3965-d8d2-4c22-99d5-e17e6a5db36b","action":"list-actions","data":{}}'

Read the action output:

    GET cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/output

The action output is a string in JSON format.

```json
[
  "list-actions",
  "destroy-module",
  "add-module",
  "add-user",
  "remove-module",
  "remove-repository",
  "update-routes",
  "alter-user",
  "grant-actions",
  "remove-user",
  "add-repository",
  "get-module-info",
  "join-cluster",
  "list-installed-modules",
  "list-updates",
  "revoke-actions",
  "create-module",
  "get-status",
  "add-node",
  "alter-repository",
  "create-cluster",
  "list-modules",
  "list-repositories"
]
```

Read the collected stderr data from the action steps:

    GET cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/error

Read the originally submitted task payload:

    GET cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/context

Read the action exit code:

    GET cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/exit_code

The above keys are transient. After a few hours they are evicted. This
command return the remaining key TTL (Time To Live), in seconds.

    TTL cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/error

While the action is running, some messages are sent through the progress
channel. It is possible to get them by subscribing the channel.

    SUBSCRIBE progress/cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b

### Tasks submission

The Python package `agent.tasks` implements the task processing, handling
communication errors and re-connection/retrying logic at every stage, like:

- Task submission
- Task progress tracking
- Task response parsing

Example:

```python
import agent.tasks
print(agent.tasks.run(agent_id='cluster', action='list-actions', endpoint='redis://cluster-leader'))
```

Note the `endpoint=` argument: it uses the `redis://` protocol, that
requires `cluster` privileges to work.

As said above, only the `cluster` agent and `api-server` are allowed to
LPUSH new tasks in Redis. Module agents must submit tasks for other agents
as HTTP requests, through the `api-server`. 

Basically this can be done in two steps:

1. obtain a JWT token with the module Redis credentials
2. send the task request

The Python module `agent.tasks` implements the two steps internally.

It has a few assumptions, that are always satisfied by an action step environment:

* It obtain the module credentials from the environment from `REDIS_USER`
  and `REDIS_PASSWORD`
* It stores the JWT token in the current working directory (the file is
  `./apitoken.cache`)

This is an example call, note that the `agent_id` is obtained from the
environment.

```python
import agent.tasks, os
print(agent.tasks.run(agent_id=os.environ['AGENT_ID'], action='list-actions'))
```

The above call fails, because agents are not generally authorized to run
`list-actions`. However, just for this example, it is possible to override
the Redis credentials received from the environment.

```python
import agent.tasks, os
os.environ['REDIS_USER'] = 'admin'
os.environ['REDIS_PASSWORD'] = 'Nethesis,1234'
print(agent.tasks.run(agent_id=os.environ['AGENT_ID'], action='list-actions'))
```

The next section explains how to authorize agents to run actions through
roles assignment.

### Roles and authorizations

In normal conditions module actions must be assigned to `roles`, and roles
must be granted to agents and users.

In Redis, a SET contains the actions assigned to a role. If default
assignments are not enough the set can be modified by the module during
the `create-module` action.

The following Bash command in the `create-module` action allows agents
with the `actionsreader` role to run `list-actions` on `AGENT_ID` (e.g.
`module/mymod1`).

```sh
    redis-exec SADD "${AGENT_ID}/roles/actionsreader" "list-actions"
```

Now it is possible to grant the `actionsreader` on `module/mymod1` to
another agent (e.g. `module/authmod2`). This is the corresponding Redis command:

    HSET roles/module/authmod2 module/mymod1 actionsreader

Roles are always granted by the `cluster` agent.

Users with the `owner` role can use the following actions to manage users
and their roles:

- `add-user`
- `remove-user`
- `grant-actions`
- `revoke-actions`

Those actions are implemented with the `cluster.grants` Python module.

A module can require additional roles to other modules when it is
instantiated. Set the label `org.nethserver.authorizations` of the module
image to a space separated list of values. Each value describes what is
the required role and what are the modules. The value syntax is `{agent
selector}:{role}`. For instance

    org.nethserver.authorizations = mymod@node:actionsreader funmod@cluster:rolex coolmod@any:roley

If the installed module is `authmod2`, the above label is interpreted in
the following way:

1. grant `authmod2` role `actionsreader` on the _default instance_ of
   module `mymod` running on the same node (`@node` suffix). The first
   instance of `mymod` installed on the node of `authmod2` is usually the
   default one.

2. grant `authmod2` role `rolex` on the _default instance_ of `funmod` at
   cluster level (`@cluster` suffix). That means the first instance of `funmod`
   installed in the cluster on any node.

3. grant `authmod2` role `roley` on any existing instance of module `coolmod`.

Other possible *agent selector* values:

- `cluster` - selects the cluster agent
- `node` - selects the node agent where the module is running

Depending on the running node ID the value resolves to an agent ID. The above rules
are implemented in Python by the `agent.resolve_agent_id()` function

```python
import agent, os
print(agent.resolve_agent_id('traefik@node', node_id=os.getenv('NODE_ID', 1)))
```

The `node_id` argument is optional.

---
Next [Modules](modules.md)

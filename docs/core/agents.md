---
layout: default
title: Agents
nav_order: 3
parent: Core
---

# Agents

Agents are daemons that listens to Redis queues and execute tasks inside the system.

* TOC
{:toc}

## Agent types

The core `cluster` and `node` agent, and every module instance run their
own `agent` process. The process is started by the `agent.service` Systemd
user unit and the `agent@.service` Systemd system template unit provided
by the core.

Systemd units for agents are:

- `agent@cluster.service` running as root. Its
  actions are defined in `/var/lib/nethserver/cluster/actions`

- `agent@node.service` running as root. Its actions
  are defined in `/var/lib/nethserver/node/actions`

- `agent.service` running as non-privileged Unix user
  for each rootless module instance. See the "Additional modules" section
  below for more details

See also the [Agent documentation](https://github.com/NethServer/ns8-core/blob/main/core/agent/README.md).

## Tasks processing

Agents wait for tasks on a Redis list. Only the `cluster` agent and the
`api-server` are allowed to push a new task for an arbitrary module agent. The
task contain two main attributes: the action name (`action`), and the
action input (`data`).

The following Redis command run the `list-actions` action on the cluster agent:

    LPUSH cluster/tasks '{"id":"1a4a3965-d8d2-4c22-99d5-e17e6a5db36b","action":"list-actions","data":{}}'

Read the action output:

    GET task/cluster/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/output

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

    GET task/cluster/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/error

Read the originally submitted task payload:

    GET task/cluster/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/context

Read the action exit code:

    GET task/cluster/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/exit_code

The above keys are transient. After a few hours they are evicted. This
command return the remaining key TTL (Time To Live), in seconds.

    TTL task/cluster/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b/error

While the action is running, some messages are sent through the progress
channel. It is possible to get them by subscribing the channel.

    SUBSCRIBE progress/cluster/task/1a4a3965-d8d2-4c22-99d5-e17e6a5db36b

## Tasks submission

The Python package `agent.tasks` implements the task processing, handling
communication errors and re-connection/retrying logic at every stage, like:

- Task submission
- Task progress tracking
- Task response parsing
- Task cancellation when the TERM signal is received

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
roles assignment, and typical error messages if they are not properly set
up.

## Roles and authorizations

In normal conditions module actions must be assigned to `roles`, and roles
must be granted to agents and users. Symptoms of misconfigured roles and
authorizations are `api-server` log lines like

    [AUTH] error retrieving user authorizations: redis: nil

The task execution typically fails with a 403 error. Sometimes the UI
displays an "Authorization error" or "Permission error" notification.

In Redis, a SET defines the actions assigned to a role. If default
assignments are not enough the set can be modified by the module during
the `create-module` action.  A set element represent an exact action name
or an action name pattern, expressed in [glob-pattern
syntax](https://pkg.go.dev/path/filepath#Match) (e.g. `get-*`).

The following Bash command in the `create-module` action allows agents
with the `actionsreader` role to run `list-actions` on `AGENT_ID` (e.g.
`module/mymod1`).

```sh
redis-exec SADD "${AGENT_ID}/roles/actionsreader" "list-actions"
```

Same as above, but to allow also any other action with name prefix `list-`:

```sh
redis-exec SADD "${AGENT_ID}/roles/actionsreader" "list-*"
```

Now it is possible to grant the role `actionsreader` on `module/mymod1` to
another agent (e.g. `module/authmod2`). This is the corresponding Redis command:

    HSET roles/module/authmod2 module/mymod1 actionsreader

Granted permissions are stored in Redis HASH keys with `roles/`
prefix. For example the `roles/admin` key stores the granted permissions
for the `admin` user, you can inspect it with this Redis
command:

    HGETALL roles/admin

It works the same for a module. A module key is for example
`roles/module/traefik1`:

    HGETALL roles/module/traefik1

If the module instance needs to run its own actions, extend the builtin
role `selfadm`. A basic role definition is set up by the builtin step
`10selfadm_role`. It adds `configure-module` to the `selfadm` role during
the module instance creation. If this is undesirable, override the builtin
step.

For instance add the following command in a Bash step of
`configure-module`:

```sh
redis-exec SADD "${AGENT_ID}/roles/selfadm" "get-configuration"
```

Roles are always granted by the `cluster` agent.  Users with the `owner`
role on `cluster` (cluster administrators) can use the following cluster
actions to manage users and their roles:

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

- `self` - selects the module itself. Note that the module instance is
  implicitly granted the role `selfadm` on itself.
- `cluster` - selects the cluster agent
- `node` - selects the node agent where the module is running

Depending on the running node ID the value resolves to an agent ID. The above rules
are implemented in Python by the `agent.resolve_agent_id()` function

```python
import agent, os
print(agent.resolve_agent_id('traefik@node', node_id=os.getenv('NODE_ID', 1)))
```

The `node_id` argument is optional.

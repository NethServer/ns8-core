---
layout: default
title: Database
nav_order: 2
parent: Core
---

# Database

The Redis database runs as a Podman rootfull container, bound to TCP port 6379.

* TOC
{:toc}

Redis is not used only as a database:
  * it exchanges messages among agents and the API server
  * it helps modules to discover information about the core and other modules
  * it provides a signaling bus based on its PUB/SUB feature
  * it stores (a copy of) the system and modules configuration
  * it is an authentication backend for the management UI and the agents

Get the Redis service status:

    systemctl status redis

If any change occurs in the Redis database, it is persisted on the disk
within 5 seconds. The following command prints the database path:

    podman volume inspect redis-data | jq -r .[].Mountpoint

Redis provides blocking *pop* operations on lists. The BRPOP command is
used by agents to wait on a queue for incoming tasks. Task queues keys
look like `module/{module_id}/tasks`.

Redis provides PUB/SUB operations on *channels*. Channels are used to
implement the observer pattern, for instance in the following situations:

1. An agent wants to notify the UI about the progress of a running task.
   The channel name in this case looks like
   `progress/module/{module_id}/task/{uuid}`.

2. An agent wants to communicate to other agents that a particular event
   occurred (e.g. a TLS certificate was automatically renewed). The
   channel name in this case could be like
   `module/{module_id}/event/{event_name}`

## Access

Access Redis with one of the following commands:

    redis-cli PING

* The `redis-cli` command is faster when invoked by the root user, because
  it runs a command in the existing `redis` rootfull container. For
  non-root users a temporary container is created and destroyed.

* The `redis-cli` command attempts to connect with the higher available
  Redis privileges, by reading `agent.env` files.

A Python-based helper script, `redis-exec` is available and it better
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

## Roles

In a NS8 cluster the leader node runs Redis with *master* role. All
**write operations** must be executed on the leader instance. Worker nodes run
Redis with *replica* role, so they allow only **read operations**.

Changes are propagated from the leader instance to the worker instances.
The following command on the leader node describes the replication status:

    redis-cli ROLE

Example of the command output:

    1) "master"
    2) (integer) 57693
    3) 1) 1) "10.5.4.2"
          2) "6379"
          3) "57693"

The same command in a worker node has a different output. Note that the
deprecated term _slave_ may still used by the Redis protocol to describe
the replica role:

    1) "slave"
    2) "10.5.4.1"
    3) (integer) 6379
    4) "connected"
    5) (integer) 57721

Refer to the Redis official documentation for the exact meaning of the
command output.

The above command identifies the Redis leader instance by IP address and
port number. To know the leader node ID number run this command instead

    redis-cli HGET cluster/environment NODE_ID

To get back the IP address of an arbitrary node, given its node ID, e.g. 3:

    redis-cli HGET node/3/vpn ip_address

## Schema

Redis is a key/value database. Some key naming rules are enforced because
access rights are based on key name patterns. Both admin UI users and
agents have their own Redis credentials for authentication and
authorization.

Some key rule examples:

* `cluster/environment`, `node/{node_id}/environment` The key type is
  HASH. It stores a copy of environment variables of respectively the
  cluster agent running on the leader node, and the agent environment of
  *node_id*.

* `module/{module_id}/environment` The key type is HASH. It stores
  a copy of environment variables of the module *module_id*.

* `node/{node_id}/vpn` The key type is HASH. It stores the cluster VPN
  configuration applied by the node agent *node_id*

* `cluster/node_sequence` The key value is an integer. As the name
  suggests it is used to generate a node ID, when a new node is added to
  the cluster.

Next sections describe the value stored for each **database key**. Key names
follow this notation:

- `/` logical key namespace separator
- `{var}` (curly braces with name) placeholder for any string value
- ` ` (single space) hash field separator

A similar pattern is used also for channel names of
[events]({{site.baseurl}}/core/events).

### cluster/

Key prefix for the cluster agent key namespace. Keys with this prefix are
publicly accessible. As many keys has this prefix, see also the next
subsections for more information.

|key|type|description|
|---|----|-----------|
|cluster/network                        |CIDR       |WireGuard VPN network, default: `10.5.4.0/24` |
|cluster/default_instance/{image}       |STRING     |a module ID for the given image ID |
|cluster/node_sequence                  |INTEGER    |generate node IDs, default: `0` |
|cluster/module_sequence/{image}        |INTEGER    |module sequence to generate instances of image ID, default: `0` |
|cluster/module_node                    |HASH       |The module-node association, used for roles assignment|
|cluster/module_domains                 |HASH       |Store relation of modules with user domains. Hash key is MODULE_ID, hash value is a list of domains separated by spaces, e.g. `mydom1.tld mydom2.tld mydom3.tld`|
|cluster/authorizations/{agent_id}      |SET        |Authorization labels persistence, to enforce labels on future modules|
|cluster/roles/{role}                   |SET        |glob patterns matching the actions that {role} can run. {role} is one of "owner", "reader"...|
|cluster/environment                    |HASH       |Cluster environment variables|
|cluster/environment NODE_ID            |INTEGER    |The node ID of the leader node |
|cluster/ui_name                        |STRING     |UI label for the cluster|
|cluster/uuid                           |STRING     |Generated UUID that identifies the cluster|
|cluster/override/modules               |HASH       |Maps (image name) => (image URL), overriding the standard image name/URL resolution function |
|cluster/subscription                   |HASH       |[Subscription]({{site.baseurl}}/core/subscription) attributes in key/value pairs|
|cluster/apply_updates                  |HASH       |[Scheduled updates]({{site.baseurl}}/core/updates) attributes in key/value pairs|
|cluster/tcp_ports_demand               |HASH       |Max TCP ports a module can request. Hash key is MODULE_ID, hash value is the max TCP ports.|
|cluster/udp_ports_demand               |HASH       |Max UDP ports a module can request. Hash key is MODULE_ID, hash value is the max UDP ports.|

#### cluster/smarthost

This is not a key prefix, it is a key of type HASH by itself and contains
setting storage for [SMTP smarthost]({{site.baseurl}}/core/smarthost).

|key|type|description|
|---|----|-----------|
|cluster/smarthost                      |HASH   |Settings and credentials of a SMTP smarthost provider|
|cluster/smarthost port                 |INTEGER|Public TCP port of the smtp server|
|cluster/smarthost enabled              |0\|1   |(0 disabled, 1 enabled) enable the smarthost provider|
|cluster/smarthost tls                  |0\|1   |(0 disabled, 1 enabled) enable STARTTLS|
|cluster/smarthost tls_verify           |0\|1   |(0 disabled, 1 enabled) verify if the certificate is valid and if the hostname is associated to the certificate|
|cluster/smarthost username             |STRING |username for smtp credentials|
|cluster/smarthost password             |STRING |password for smtp credentials|

#### cluster/repository/

Setting storage for [Software repositories]({{site.baseurl}}/core/software_repositories).

|key|type|description|
|---|----|-----------|
|cluster/repository/{repo}              |HASH   |A software {repo}|
|cluster/repository/{repo} url          |STRING |HTTPS URL of repository|
|cluster/repository/{repo} status       |0\|1   |(0 disabled, 1 enabled)|
|cluster/repository/{repo} testing      |0\|1   |(0 disabled, 1 enabled)|
|cluster/repository_cache/{repo}        |HASH   |Cache of remote repo metadata, with a TTL|
|cluster/repository_cache/{repo} data   |STRING |It contains a JSON string from repodata.json|
|cluster/repository_cache/{repo} updates|STRING |Most recent modification date and time of remote repodata.json|

#### cluster/backup_repository/

Each backup repository is identified with a UUID-based key of type HASH,
under the key prefix `cluster/backup_repository/`.

|key|type|description|
|---|----|-----------|
|cluster/backup_repository/{repo}                       |HASH   |A backup repository UUID|
|cluster/backup_repository/{repo} url                   |STRING |Restic URL of repository, eg. `b2:<bucket>`, `s3:s3.amazonaws.com/<bucket>`|
|cluster/backup_repository/{repo} password              |STRING |Restic encrypt password|
|cluster/backup_repository/{repo} name                  |STRING |Custom repository name|
|cluster/backup_repository/{repo} provider              |STRING |Provider name, can be: `backblaze`, `aws`|
|cluster/backup_repository/{repo} b2_account_id         |STRING |Backblaze B2 account ID, present only if provider is `backblaze`|
|cluster/backup_repository/{repo} b2_account_key        |STRING |Backblaze B2 account KEY, present only if provider is `backblaze`|
|cluster/backup_repository/{repo} aws_access_key_id     |STRING |AWS S3 access key, present only if provider is `aws`|
|cluster/backup_repository/{repo} aws_secret_access_key |STRING |AWS S3 secret key, present only if provider is `aws`|
|cluster/backup_repository/{repo} aws_access_key_id     |STRING |AWS S3 region, present only if provider is `aws`|

#### cluster/backup/

The attribute of a backup schedule are stored in an HASH key under the
`cluster/backup/` key prefix.

|key|type|description|
|---|----|-----------|
|cluster/backup/{id}                       |HASH    |A backup repository numeric ID|
|cluster/backup/{id} name                  |STRING  |Backup name|
|cluster/backup/{id} enabled               |0\|1    |(0 disabled, 1 enabled)|
|cluster/backup/{id} repository            |STRING  |UUID of backup_repository|
|cluster/backup/{id} retention             |INTEGER |Number of snapshots to keep|
|cluster/backup/{id} schedule              |STRING  |Schedule time using `onCalendar` systemd syntax|
|cluster/backup/{id} schedule_hint         |STRING  |Schedule in JSON format for the UI|

#### cluster/user_domain/

Settings storage for external user domains, LDAP databases not hosted by
the cluster.

|key|type|description|
|---|----|-----------|
|cluster/user_domain/ldap/{domain}/conf                 |HASH   |An external LDAP domain|
|cluster/user_domain/ldap/{domain}/conf schema          |STRING |It can be `ad` or `rfc2307`|
|cluster/user_domain/ldap/{domain}/conf bind_dn         |STRING |LDAP bind DN|
|cluster/user_domain/ldap/{domain}/conf bind_password   |STRING |LDAP bind password|
|cluster/user_domain/ldap/{domain}/conf base_dn         |STRING |LDAP base DN|
|cluster/user_domain/ldap/{domain}/conf tls             |STRING |Can be `on` or `off`|
|cluster/user_domain/ldap/{domain}/conf tls_verify      |STRING |Can be `on` or `off`|
|cluster/user_domain/ldap/{domain}/providers            |LIST   |List of domain provider hosts|
|cluster/user_domain/ldap/{domain}/ui_names             |HASH   |UI labels for domain providers. The key is the provider name, eg `ldap.ns.test:389` => `mylabel`|

#### cluster/counters_cache/

Store the counters of users and groups for each user domain.

|key|type|description|
|---|----|-----------|
|cluster/counters_cache/users/{domain}                  |INT    |Cached number of users in the domain, updated by list APIs|
|cluster/counters_cache/groups/{domain}                 |INT    |Cached number of groups in the domain, updated by list APIs|

### node/

Keys with this key prefix are publicly readable. Each node agent has write
access to its keys under `node/{id}/`, where `{id}` is the integer identifier
of the node (e.g. `node/1/...`).

|key|type|description|
|---|----|-----------|
|node/{id}/vpn                      |HASH||
|node/{id}/vpn ip_address           |STRING     |IP address in cluster/network|
|node/{id}/vpn public_key           |STRING     |Public WireGuard VPN key|
|node/{id}/vpn destinations         |STRING     |List of IP addresses routed through the VPN|
|node/{id}/vpn endpoint             |STRING     |Public IP or host name of the VPN endpoint with :port suffix|
|node/{id}/flags                    |SET        |Flags to mark a node. Supported flags: `nomodules`, the node can't run any module|
|node/{id}/roles/{role}             |SET        |glob patterns matching the actions that {role} can run. {role} is one of "owner", "reader"...|
|node/{id}/tcp_ports_sequence       |INTEGER    |Lower boundary for the next TCP_PORTS_RANGE allocation, default: `20000`|
|node/{id}/udp_ports_sequence       |INTEGER    |Lower boundary for the next UDP_PORTS_RANGE allocation, default: `20000`|
|node/{id}/default_instance/{image} |STRING     |A module ID for the given image ID|
|node/{id}/environment              |HASH       |Node environment variables|
|node/{id}/ui_name                  |STRING     |UI label for the node|


### module/

Keys with this key prefix are publicly readable. Each node agent has write
access to its keys under `module/{id}/`, where `{id}` is the module
 instance identifier (e.g. `module/openldap1/...`).

|key|type|description|
|---|----|-----------|
|module/{id}/environment            |HASH|
|module/{id}/roles/{role}           |SET        |glob patterns matching the actions that {role} can run. {role} is one of "owner", "reader"...|
|module/{id}/backups                |SET        |List of backup numeric IDs|
|module/{id}/flags                  |SET        |Images flags copied from `org.nethserver.flags` image label|
|module/{id}/srv/{transport}/{service} | HASH   |Service discovery information for other modules. See [Service providers]({{site.baseurl}}/modules/service_providers) |
|module/{id}/backup_status/{backup_id} | HASH   |Backup status information with Restic output |

### task/

Keys with this prefix hold sensitive information about a task execution.
The `{agent_id}` identifies the agent that can read/write the keys. For
instance `cluster`, `node/2`, `module/traefik1`...

|key|type|description|
|---|----|-----------|
|task/{agent_id}/context    |STRING  |JSON representation of a task|
|task/{agent_id}/error      |STRING  |stderr data generated by the action execution|
|task/{agent_id}/output     |STRING  |task output, JSON-encoded|
|task/{agent_id}/exit_code  |INTEGER |action exit code|

### user/

|key|type|description|
|---|----|-----------|
|user/{username}                |HASH||
|user/{username} display_name   |STRING|Name and surname|
|user/{username} email          |STRING|Mail address|


### roles/

Each user (or agent) can assume at most one role when it interacts with
another agent. This relation is stored in Redis in a key of type HASH. A
key of the hash is the agent identifier and the corresponding value is the
name of the role.

|key|type|description|
|---|----|-----------|
|roles/{username}             |HASH   ||
|roles/{username} {agent_id}  |STRING |Name of the role - {agent_id} is for instance `module/mail1`, `cluster`...|

### module/traefikX/

Each node runs a Traefik module instance. TLS certificates are stored under its module key prefix.

|key|type|description|
|---|----|-----------|
|module/traefik{X}/certificate/{name}       |HASH|  |
|module/traefik{X}/certificate/{name} key   |STRING |Base64 key of a x509 certificate|
|module/traefik{X}/certificate/{name} cert  |STRING |Base64 x509 certificate|

### AGENT_ID/tasks

Each [agent]({{site.baseurl}}/core/agents) runs a wait loop to run tasks
assigned to it. It waits on a task queue implemented with a Redis LIST key
with the form `{AGENT_ID}/tasks`.

|key|type|description|
|---|----|-----------|
|cluster/tasks                      |LIST      |Task items are JSON encoded objects|
|node/{id}/tasks                    |LIST      |Task items are JSON encoded objects|
|module/{id}/tasks                  |LIST      |Task items are JSON encoded objects|

An element of the task queue is a JSON-encoded string representing an
object with the following basic attributes:

|attribute|type|description|
|---|----|-----------|
|id        |string |UUID generated by the API server |
|action    |string |e.g. `list-users`, `install`, `create-domain`...|
|data      |any    |any JSON type (object, list...)|

For instance

```json
{
  "id": "66b73f7a-8998-4262-a784-36639fc4b2c1",
  "action": "list-actions",
  "data": ["cancel-task", "list-actions"]
}
```

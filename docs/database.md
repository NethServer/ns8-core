# Database: Redis

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

---
Next: [Proxy and certificates](proxy_certificates.md)

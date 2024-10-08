---
layout: default
title: Service providers
nav_order: 100
parent: Modules
---

# Service providers

A module can become the provider of some kind of service for other modules
running in the same cluster. To ease the service discovery, information
about how to reach the service and any other application parameter must be
written in the Redis database as explained in the next sections.

## `srv` keys

A module can own one or more `srv` Redis keys. Each one identifies a
service endpoint on a given transport protocol. For instance

    module/mail1/srv/tcp/imap
    module/mail1/dovecot/srv/http/openmetrics
    module/mail1/rspamd/srv/http/openmetrics

The generic key pattern is

    module/{MODULE_ID}[/{OPTIONAL_QUALIFIER}]/srv/{TRANSPORT}/{SERVICE_NAME}

Each key must be of type HASH. The hash content depends on the service
requirements, typically the `host` and IP `port` entries are set.

## Raising events

The service provider **should raise** one or more distinct events whenever
the information in its `srv` keys is changed. The event name and payload
must be documented by the service provider. A common convention for the
event name is:

    service-{SERVICE_NAME}-changed

The event payload is a JSON-encoded string representing an object. The
object can contain the needed attributes to identify the source of the event.
They are necessary to event listeners to identify the event
source and retrieve further information as needed, an example payload can be:

```json
{
    "key": "module/mail1/srv/tcp/imap",
    "module_uuid": "8d257122-0a7f-49c7-a620-08961a68cfa0",
    "module_id": "mail1"
}
```

# Service discovery

Service consumers can look up a given service name with the
`agent.list_service_providers()` function, or the corresponding
`list-service-providers` base module action.

For instance, to find IMAP servers on node 1:

    api-cli run module/mail1/list-service-providers --data '{"service":"imap", "transport":"tcp", "filter":{"node":"1"}}' | jq

Yields

```json
[
  {
    "port": "143",
    "host": "10.5.4.1",
    "node": "1",
    "user_domain": "dp.nethserver.net",
    "module_uuid": "8d257122-0a7f-49c7-a620-08961a68cfa0",
    "module_id": "mail1",
    "transport": "tcp",
    "ui_name": null
  }
]
```

The input JSON object attributes are:

- `service` (mandatory): look up this service name
- `transport` (optional): look up records with the given transport, or any
  transport if the attribute is missing
- `filter`: an object where each entry represents an additional filter
  condition. Only services matching all the given conditions are returned.
  The actual attributes depend on the record structure: service providers
  must document which attributes they define. The `list-service-providers`
  base implementation in any case, provides at least the following
  attributes for convenience:
    - `module_id` copy of `HGET module/{id}/environment MODULE_ID`
    - `module_uuid` copy of `HGET module/{id}/environment MODULE_UUID`
    - `ui_name`, copy of `GET module/{id}/ui_name`)
    - `transport`, copy of `transport` of the Redis key name

The following Python code snippet invokes `agent.list_service_providers()`
to discover an IMAP server. It connects to the local Redis replica, so it
works even if the leader is temporarly unreachable: this makes the service
startup more robust.

```python
import agent
rdb = agent.redis_connect(use_replica=True) # connect the local Redis replica
print(agent.list_service_providers(rdb, 'imap', 'tcp'))
```

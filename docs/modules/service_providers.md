---
layout: default
title: Service providers
nav_order: 10
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

The service provider **must raise** one or more distinct events whenever
the information in its `srv` keys is changed. The event name and payload
must be documented by the service provider. A common convention for the
event name is:

    service-{SERVICE_NAME}-changed

The event payload is a JSON-encoded string representing an object. The
object contains at least the `key`, `module_uuid` and `module_id`
attributes. They are necessary to event listeners to identify the event
source and retrieve further information as needed:

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

Python code snippet to invoke the `agent.list_service_providers()` function:

```python
import agent
rdb = agent.redis_connect() # full read-only access on every key
print(agent.list_service_providers(rdb, 'imap', 'tcp'))
```

Note that the function always adds the following attributes to the response records:

- `module_id`
- `module_uuid`
- `ui_name`
- `transport`

---
layout: default
title: Events
nav_order: 8
parent: Core
---

# Events

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

## Signaling events

To signal an event, just [PUBLISH](https://redis.io/commands/PUBLISH) a message to the channel.
All events take a JSON object as a parameter.

Example:
```
redis-cli PUBLISH module/traefik1/event/test '{"field": "value"}'
```

## Handling events

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

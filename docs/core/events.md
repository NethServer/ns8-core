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
- use past tense inside the name, like `account-provider-changed`,
  `user-domain-changed`, `module-added`...
- accept a parameter in JSON format
- the parameter should contain minimal required info about the event
- optionally, the parameter can contain extra data which can ease the event usage

Well known events:
- `user-domain-changed`: the configuration of a user domain has changed
- `account-provider-changed`: an account provider instance was removed or added to a user-domain

Events fired by the `cluster` agent (i.e. channel is `cluster/event/<event name>`):
- `module-added`: the event is fired at the end of the add-module process to inform other modules that a new module has been installed on the cluster
- `module-removed`: the event is fired at the end of the remove-module process to inform other modules that a module has been removed on the cluster

## Signaling events

To signal an event, just [PUBLISH](https://redis.io/commands/PUBLISH) a message to the channel.
All events take a JSON object as a parameter.

Example:
```
redis-cli PUBLISH module/traefik1/event/test '{"field": "value"}'
```

## Handling events

Events are handled like actions by the `agent` service.

Modules can define their event handlers under a dedicated directory:
`${AGENT_INSTALL_DIR}/events/`. The handler directory contains one or more
step scripts that must be executable files. For instance:

- **rootfull module**: `/var/lib/nethserver/mymodule/events/some-event-happened/10handler`.
- **rootless module**: `/home/mymodule/.config/events/some-event-happened/10handler`

Like actions, the event payload can be read by the step scripts from the
standard input file descriptor.

The event handlers receive two additional environment variables:

- `AGENT_EVENT_SOURCE` the agent ID that originated the event
- `AGENT_EVENT_NAME` the name of the event (it is the name same of the
  directory where of the event handler)

Lines sent to standard error and standard output are relayed to the system
log.

Similarly to actions, if an handler step exit code is non-zero, the
handler execution is aborted. Subsequent steps are not executed and an
error message is sent to the log.

---
layout: default
title: Subscription
parent: Core
---

# Subscription

The core provides builtin features to subscribe both free and paid support services.

The Redis HASH key `cluster/subscription` holds the attributes of the subscription.

- `provider` The subscription provider identifier. Possible values: `nsent`, `nscom`
- `system_id` System subscription identifier
- `auth_token` Authentication token for subscription APIs
- `vpn_cert_cn` X509 Common Name for server certificate validation
- `vpn_peer_host` (default "sos.nethesis.it")
- `vpn_peer_port` (default "1194")
- `support_user` User name added to cluster-admin when a support session
  is started on the leader node. Password is set to a random UUID value.
  The account is removed when the support session is stopped

When a cluster has the Redis HASH key `cluster/subscription` some Systemd
services and timers are started and enabled at boot. The unit
`check-subscription.service` starts at boot and keep all the other unit
status aligned with the Redis key value.  This is a list of the services
managed by `check-subscription`:

- `send-inventory` (leader only) Run at night, send the cluster inventory to the provider
- `send-heartbeat` (leader only) Run every 10 minutes to signal the cluster liveness
- `send-backup` (leader only) Run at night, send the cluster encrypted backup to the provider
- `node-monitor` Check the node loadavg, root/boot space, free swap and raises alarms

The subscription status and running services are checked when:
- the node boots
- the node joins the cluster
- the cluster is restored from a cluster backup file
- the cluster subscription is enabled or disabled
- the cluster leader node changes

## APIs

Cluster:

- `get-subscription`
- `set-subscription`

Node:

- `get-support-session`
- `start-support-session`
- `stop-support-session`

## Events

- `subscription-changed`
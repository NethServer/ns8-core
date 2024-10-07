---
layout: default
title: Subscription
nav_order: 16
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
- `apply-updates` (leader only) Run at night, apply core, modules and OS
  updates according to configuration in Redis key `cluster/apply_updates`.
  See also [core updates]({{site.baseurl}}/core/updates)

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

## Remote support

Support sessions are started/stopped by each cluster node separately. Relevant node actions are

- `start-support-session`
- `stop-support-session`

For example, it is possible to start/stop a support session for any node
of the cluster. Run on the leader node this command:

    api-cli run node/7/start-support-session

The action prints the session ID, required to ask support:

    {"session_id": "811630de-67f4-5b6d-8b2f-7cc01f592b1b"}

As alternative, log in on the node and start manually the support service
with this command:

    systemctl start support
    systemctl status support

The second command prints a journal excerpt that should contain the session ID.

When a node support session is started

- support SSH key is added to the node /root/.ssh/authorized_keys file.
  Only connections from the node itself are allowed with that key
- support credentials are enabled for cluster-admin access. The user name
  is defined in Redis. Type `HGET cluster/subscription support_user` to
  see its value, by default it is `nethsupport`. The password is set to
  the support session ID value
- an OpenVPN tunnel is established with the support server, allowing
  connections to sshd and cluster-admin. The support server address and
  port can be overridden in the Redis `cluster/subscription` key

When the support session is terminated

- the OpenVPN tunnel is closed
- cluster-admin access is removed
- SSH key is removed

## OpenVPN support client

The Systemd unit `support.service` controls an OpenVPN client running in a
container, within a separate network namespace where the TUN device is
created.  This network namespace is configured with `nft` to forward some
TCP ports from the VPN tunnel to services of the node. For example

- Forward VPN ports 981 to 127.0.0.1:22
- Forward VPN ports 22, 9090, 80 and 443 to the respective port on 127.0.0.1

See `nft.conf` for the port forwarding rules.

## Inventory format

The inventory format is compatible with NS7 inventory and merges data
collected on the leader node with data collected by calls to the
`get-facts` action of each module instance.

If a module implements the `get-facts` action, its output is merged in the
inventory data, no matter the node hosting it.

This is an excerpt of `get-facts` data merged with user and group counters
for the LDAP domain `ad.example.org`:

```json
{
  "cluster_module_domain_table": [
    {
      "active_users": 65,
      "domain_id": "ad.example.org",
      "instance_of": "mail",
      "module_id": "mail8",
      "node_id": 5,
      "total_groups": 31,
      "total_users": 86
    },
    {
      "active_users": 65,
      "domain_id": "ad.example.org",
      "instance_of": "ejabberd",
      "module_id": "ejabberd1",
      "node_id": 5,
      "total_groups": 31,
      "total_users": 86
    },
    {
      "active_users": 41,
      "domain_id": "ad.example.org",
      "instance_of": "nextcloud",
      "module_id": "nextcloud1",
      "node_id": 1,
      "total_groups": 31,
      "total_users": 92
    },
    {
      "active_users": 65,
      "domain_id": "ad.example.org",
      "instance_of": "webtop",
      "module_id": "webtop7",
      "node_id": 5,
      "total_groups": 31,
      "total_users": 86
    },
    {
      "active_users": 0,
      "domain_id": null,
      "instance_of": "mattermost",
      "module_id": "mattermost1",
      "node_id": 1,
      "total_groups": null,
      "total_users": 116
    }
  ],
  "cluster_user_domains_table": [
    {
      "active_users": 65,
      "domain_id": "ad.example.org",
      "ejabberd_count": 1,
      "mail_count": 1,
      "nextcloud_count": 1,
      "total_groups": 31,
      "total_users": 86,
      "webtop_count": 1
    }
  ]
}
```

Notes:

- The Mattermost instance `mattermost1` is not bound to the LDAP domain,
  users are stored into the module database. Its `get-facts` action
  returns `active_users` and `total_users` counters with data stored
  inside its database.
- The Nextcloud instance `nextcloud1` implements `get-facts`. Its output
  overrides the `active_users` and `total_users` counters inherited from
  the `ad.example.org` user domain.
- Instances `mail8`, `ejabberd1` and `webtop7` inherit counters from
  domain `ad.example.org`

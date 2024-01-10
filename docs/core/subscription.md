---
layout: default
title: Subscription
parent: Core
---

# Subscription

The core provides builtin features to subscribe both free and paid support services

The Redis HASH key `cluster/subscription` holds the attributes of the subscription.

- `system_id` System subscription identifier
- `auth_token` Authentication token for subscription APIs
- `vpn_cert_cn` X509 Common Name for server certificate validation
- `support_user` User name added to cluster-admin when a support session
  is started on the leader node. Password is set to a random UUID value.
  The account is removed when the support session is stopped


## Events

- `subscription-changed`
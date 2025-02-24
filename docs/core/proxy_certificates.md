---
layout: default
title: Proxy and certificates
nav_order: 5
parent: Core
---

# Proxy and certificates

Traefik is the HTTP proxy server used to access HTTP and HTTPS services,
and it also manages TLS certificates of the local cluster node. See
[Traefik
README](https://github.com/NethServer/ns8-traefik/blob/main/README.md) for
more information.

When a TLS certificate is uploaded from the cluster-admin UI, or it was
requested with ACME protocols and obtained from Let's Encrypt, the Traefik
module notifies the `certificate-changed` event, with a JSON message that
contains the following fields:

- `node_id` (integer): the node id where the Traefik instance is running,
- `module_id` (string): the ID of Traefik itself,
- `names` (array): list of DNS names (also wildcard patterns) that are
  affected by the changed certificate.

Example:
```json
{
    "node_id": 2,
    "module_id": "traefik2",
    "names": ["*.example.org", "www.example.com"]
}
```

The event is published in the Redis channel
`module/traefik<X>/event/certificate-changed`.

Applications that listen for this event can obtain the certificate and its
private key by calling the `get-certificate` action on the module that
generated the event.

For more information, refer to [TLS certificates](../../modules/certificates/)
under the Modules section.

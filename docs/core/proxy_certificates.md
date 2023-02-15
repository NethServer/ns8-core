---
layout: default
title: Proxy and certificates
nav_order: 5
parent: Core
---

# Proxy and certificates: Traefik

Traefik is the proxy server used to access all HTTP/HTTPS services.
It also handles Let's Encrypt certificate requests and renewal.

See [Traefik README](https://github.com/NethServer/ns8-traefik/blob/main/README.md) for more info.

Let's Encrypt certificates are automatically exported to Redis upon request and renewal.
Certificates are saved under an hash named `/module/traefik<X>/certificate/<domain>` key,
i.e `/module/traefik1/certificate/server.nethserver.org`.
The certificate is saved inside the `cert` field, while the key is saved inside the `key` field.

When a certificate is exported, Traefik module fires the `certificate-updated` event with
a JSON messages. The JSON message contains the following fields:

- `key`: the X509 certificate private key, PEM format encoded with base64
- `certificate`: the X509 certificate, PEM format encoded with base64
- `node`: the node id where the traefik instance is running
- `module`: the module id of the traefik instance
- `domain`: the FQDN for the certificate

Example:
```json
{
    "key": "AAa...",
    "certificate": "BBb..",
    "node": "1",
    "module": "traefik1",
    "domain": "server.nethserver.org"
}
```

The event is published under Redis channel `module/traefik<X>/event/certificate-updated`.

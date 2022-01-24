# Proxy and certificates: Traefik

Traefik is the proxy server used to access all HTTP/HTTPS services.
It also handles Let's Encrypt certificate requests and renewal.

See [Traefik README](https://github.com/NethServer/ns8-scratchpad/blob/main/traefik/README.md) for more info.

Let's Encrypt certificates are automatically exported to Redis upon request and renewal.
Certificates are saved under an hash named `/module/traefik<X>/certificate/<domain>` key,
i.e `/module/traefik1/certificate/server.nethserver.org`.
The certificate is saved inside the `cert` field, while the key is saved inside the `key` field.

When a certicate is exported, Traefik module fires the `certificate-update` event with
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

The event is published under Redis channel `module/traefik<X>/event/certificate-update`.

See also [how to fire and handle events](#events).

## Users and groups: Ldapproxy

Users and groups are stored in an LDAP database, served by one **account
provider module**. Multiple modules can work together to serve the same
LDAP database as replicas of it. An LDAP database represents an account
**domain**.

A NS8 cluster can host multiple account domains from different
implementations. It is possible to configure external LDAP services as
hosted ones. Supported LDAP schemas are

1. Active Directory
2. RFC2307

A module can discover the list of available account domains with the
`cluster.ldapproxy` Python module. The following command dumps a list of
parameters required to connect with an LDAP database on cluster node 1.

    runagent python3 -mcluster.ldapproxy

Returned TCP endpoints are local (`host` is `127.0.0.1`) and do not
require TLS. The port number depends on the LDAP domain.

Those ports are held by the [Ldapproxy
module](https://github.com/NethServer/ns8-scratchpad/blob/main/ldapproxy/README.md).
It is a L4 proxy that relays the TCP connection to an LDAP backend server,
enabling TLS and handling backend failures as needed.

If the LDAP client module runs in a Podman container with a
`private` network, add the following arguments to the `podman run`
command:

    --network=slirp4netns:allow_host_loopback=true

Then replace `127.0.0.1` with the special `10.0.2.2` IP address, that is
translated by Podman back to the loopback device, 127.0.0.1 on the root
network namespace.

Python code example:

```python
from cluster.ldapproxy import Ldapproxy
lp = Ldapproxy()
domains = lp.get_domains_list()
print(domains)
domain = lp.get_domain("mydomain")
print(domain)
```

The module can add the following settings to `eventsgw.conf` to listen to
account provider changes:

```ini
[commands]
*/event/account-provider-changed = mycommand && systemctl --user reload mymodule.service
```

See also [how to fire and handle events](#events).

--
- [Proxy and certificates](proxy_certificates.md)

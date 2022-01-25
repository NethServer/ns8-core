---
layout: default
title: User domains
nav_order: 6
parent: Core
---

# Users and groups: Ldapproxy

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



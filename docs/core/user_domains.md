---
layout: default
title: User domains
nav_order: 6
parent: Core
---

# User domains

Users and groups are stored in an LDAP database, served by one **account
provider module**. Multiple modules can work together to serve the same
LDAP database as replicas of it. An LDAP database represents an account
**domain**.

A NS8 cluster can host multiple account domains from different
implementations. It is possible to configure and connect external LDAP
services, too. Supported LDAP schema are

1. `ad` - Active Directory
2. `rfc2307`

## LDAP service discovery

A module can discover the list of available account domains with the
`agent.ldapproxy` Python module. The following command dumps a list of
parameters required to connect with an LDAP database on cluster node 1.

    runagent python3 -magent.ldapproxy

Returned TCP endpoints are local (`host` is `127.0.0.1`) and do not
require TLS. The port number depends on the LDAP domain.

Those ports are held by the [Ldapproxy
module](https://github.com/NethServer/ns8-core/blob/main/ldapproxy/README.md).
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
from agent.ldapproxy import Ldapproxy
lp = Ldapproxy()
domains = lp.get_domains_list()
print(domains)
domain = lp.get_domain("mydomain")
print(domain)
```

The module can handle the user domain configuration changes by defining an
event handler. Create an executable script with path
`${AGENT_INSTALL_DIR}/events/user-domain-changed/10handler`. For instance:

```shell
read -r domain < <(jq -r .domain)
if [[ "${domain}" == "mydomain" ]]; then
    systemctl --user reload mymodule.service
fi
```

## List users and groups

Once LDAP connection parameters are retrieved with the `agent.ldapproxy`
Python package, it is easy to get user and group listings with the
`agent.ldapclient` package.

This is an excerpt from the `cluster/list-domain-groups` API implementation:

```python
from agent.ldapproxy import Ldapproxy
from agent.ldapclient import Ldapclient

domparams = Ldapproxy().get_domain('mydom.test')
groups = Ldapclient.factory(**domparams).list_groups()
```

For complete examples see the API implementation of

- `cluster/list-domain-groups`
- `cluster/list-domain-users`
- `cluster/get-domain-user`
- `cluster/get-domain-group`

## Hidden users and groups

Some users and/or groups can be hidden to UI and other applications.

Applications might need to build LDAP search filters to configure user and
groups. The `Ldapproxy` library provides some methods that return filter
strings that honor the user and group lists used by the core. For example:

```python
from agent.ldapproxy import Ldapproxy
lp = Ldapproxy()
users_filter = lp.get_ldap_users_search_filter_clause("mydomain")
print(users_filter)
groups_filter = lp.get_ldap_groups_search_filter_clause("mydomain")
print(groups_filter)
```

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

## Bind modules and account domains

If a module wants to use an account domain it must be granted API
permissions. Add the `accountconsumer` role to the
`org.nethserver.authorizations` label of the module image. For instance
set

    org.nethserver.authorizations=cluster:accountconsumer

The module can now execute a bind procedure, so the core is aware of
existing relations between modules and account domains. When such
relations are formally established the core can

- limit/grant access to LDAP resources
- show the relations in the web user interfaces

For example, a module that uses one domain at a time can unbind the old
domain and bind the new one with a script like this:

```python
import agent

ldap_user_domain = "dept1.example.org"

# Bind the new domain, overriding previous values (unbind)
agent.bind_user_domains([ldap_user_domain])
```

At any time, retrieve the list of domains currently bound:

```python
import agent
rdb = agent.redis_connect(use_replica=True)
domlist = agent.get_bound_domain_list(rdb)
```

When the module or the domain is removed from the cluster, the relation
cleanup occurs automatically.

If the module wants to be notified of any change to the relation between
modules and user domains it can subscribe the `module-domain-changed`
event. For instance, this is the payload of such event:

```json
{
    "modules": ["mymodule1"],
    "domains": ["mydomain.test"]
}
```

The event payload contains a list of module and domains that were affected
by the relation change. Modules and domains can be either added or
removed: they are listed to ease the implementation of event handlers in
both account provider and account client modules.

For instance, the following Python excerpt checks if the module domain was
changed:

```python
event = json.load(sys.stdin)
if not os.environ["LDAP_USER_DOMAIN"] in event["domains"]:
    sys.exit(0) # nothing to do if our domain is among affected domains

# Handle the event by some means, for example
# - rewrite some config file
# - reload some service running in a container
```

## Password expiration warning

The node leader can send email notifications to users when their password
is about to expire.
The notification is available only for internal user domains.

This features has 2 requirements:

- password aging must be enabled on the user domain provider
- the cluster must be configured to send mail notifications using an internal or external [SMTP server](smarthost.md)

A timer named `password-warning.service` runs daily on all nodes to check for expiring passwords, but it actually
sends the email only on the leader node.

Configuration is saved inside Redis so it is available to all nodes:
- `cluster/password_warning/<domain>`
- `cluster/password_warning/templates`

See [database schema](database.md) for details.

## User destination mail

The mail of the user can be obtained in 2 different ways:

1. from an OpenLDAP or Samba field: the address is saved inside the `mail` field of the user object; this field has the highest priority
2. using an internal mail server instance: if there is an internal mail server, the address is automatically set to `<user>@<user_domain>`, the local Postfix
   instance will deliver the mail to the user mailbox

Please note that if the cluster is configured to send mail notifications using an external SMTP server,
the mail field must be set in the user object because the `user_domain` is not known to the external server.

To check for expiring passwords and immediately send notifications, run the following command on the leader node:

```bash
systemctl start password-warning.service
```

## Default templates

Default templates are stored inside `/etc/nethserver/password_warning` directory.
Available templates:

- `default_en.tmpl`: English mail template
- `default_it.tmpl`: Italian mail template
- `default_en.sbj`: English mail subject
- `default_it.sbj`: Italian mail subject

### Configure password warning

The configuration can be done using the cluster web interface.
If you want to configure it from command line, see the `cluster/set-password-warning` API.

The API takes the following parameters:

- `domain`: the domain name
- `notification`: enable or disable the notification
- `days`: the number of days before the password expiration
- `mail_from`: the email address of the sender
- `mail_template_name`: the name of the template to use, it could be `default_en`, `default_it` or `custom`. If set to custom, you must provide also `mail_template_content` and `mail_subject`.
- `mail_template_content`: the content of the template, must be base64 encoded
- `mail_subject`: the subject of the email in plain text. 

#### Using a Custom Template

Field `mail_template_content` and `mail_subject` use [Python string template](https://docs.python.org/3/library/string.html#template-strings).
The string representation can contain the following placeholders: `$user`, `$name`, `$domain`, `$days`, `$portal_url`.

1. Create a file named `custom.txt` with your custom content:
    ```
    Dear $user ($name) of domain $domain.
    Your password is going to expire in $days days.
    Change it here: $portal_url
    ```
2. Use the custom template:
    ```bash
    api-cli run cluster/set-password-warning --data '{"domain": "leader.cluster0.gs.nethserver.net", "notification": true, "days": 1000, "mail_from": "no-reply@nethserver.org", "mail_template_name": "custom", "mail_template_content": "'$(cat custom.txt | base64 -w0)'", "mail_subject": "my expiration"}'
    ```

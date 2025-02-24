---
layout: default
title: TLS certificates
nav_order: 80
parent: Modules
---

TLS certificates
================

X509 certificates for TLS are managed by the Traefik instance running in
the local node. For more details see
[Traefik]({{site.baseurl}}/core/proxy_certificates) in the core
documentation.

Web applications
----------------

Web applications should not care about the TLS certificate, as the Traefik
HTTP proxy listens on TCP port 443 for HTTPS requests. HTTP backends must
be configured only for clear text HTTP connections.

If the module wants to set up an host-name-based HTTP route, it is
possible to instruct Traefik to obtain a TLS certificate for it in a
unique action.

```python
agent.tasks.run(
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g. module/traefik1
    action='set-route',
    data={
        'instance': os.environ['MODULE_ID'],
        'url': f'http://127.0.0.1:{os.environ["TCP_PORT"]}',
        'host': host_fqdn,
        'lets_encrypt': True, # request cetificate for host_fqdn
        'http2https': True, # redirect http:// URL scheme to https://
    },
)
```

Note that `instance` must be a unique route identifier. If the module uses
multiple HTTP routes, add a prefix to `MODULE_ID` to distinguish them.


Other applications
------------------

Applications that need `.pem`, `.key`, `.crt` and similar files for their
configuration must listen the `certificate-changed` event. This event
occurs when a TLS certificate is uploaded from the UI or obtained/renewed
from Let's Encrypt.

A certificate request for Let's Encrypt can be issued with the
`set-cerfificate` action.

On `certificate-changed` event, invoke the `get-certificate` action to
obtain the PEM-encoded key and certificate.

For example, to issue a Let's Encrypt certificate request for
`SERVICE_FQDN` run:

```python
agent.tasks.run_nowait( # _nowait runs in background
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g. module/traefik1
    action='set-certificate',
    data={
        "fqdn": os.environ["SERVICE_FQDN"],
    },
    parent='', # Run as a root task
    extra={
        'title': 'set-certificate',
        'isNotificationHidden': False, # Show the task progress in the UI
    },
)
```

Then to handle the event, create an executable script, under
`$AGENT_INSTALL_DIR/events/certificate-changed`

```python
event = json.load(sys.stdin)

if str(event['node']) != os.environ['NODE_ID']:
    sys.exit(0) # ignore events from other cluster nodes

try
    domain_suffix = os.environ["SERVICE_FQDN"].split('.', 1)[1]
except IndexError:
    domain_suffix = ''

domain_wildcard = '*.' + domain_suffix 
if not os.environ["SERVICE_FQDN"] in event['names'] and not domain_wildcard in event['names']:
    sys.exit(0) # nothing to do for certificates of other services

agent.run_helper('systemctl', '--user', 'reload', 'myservice', check=True)
```

The handler script reloads `myservice` when the certificate of the service
is uploaded, obtained, renewed, or changed by other means.

When the service Systemd unit reloads (or restarts) the managed service,
run the `get-certificate` command... TODO


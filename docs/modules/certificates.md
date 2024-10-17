---
layout: default
title: Certificates
nav_order: 80
parent: Modules
---

Certificates
============

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
configuration can obtain the certificate ansynchronously:

1. issue the certificate request

2. handle the `certificate-updated` event

For instance, to issue an asynchronous Let's Encrypt certificate request
for `SERVICE_FQDN` run:

```python
agent.tasks.run_nowait( # _nowait runs in background
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g. module/traefik1
    action='set-certificate',
    data={
        "fqdn": os.environ["SERVICE_FQDN"],
        "sync": True, # For the UI task progress
    },
    parent='', # Run as a root task
    extra={
        'title': 'set-certificate',
        'isNotificationHidden': False, # Show the task progress in the UI
    },
)
```

Then to handle the event, create an executable script under
`$AGENT_INSTALL_DIR/events/certificate-updated`

```python
event = json.load(sys.stdin)

if str(event['node']) != os.environ['NODE_ID']:
    sys.exit(0) # ignore events from other cluster nodes

if event["domain"] != os.environ["SERVICE_FQDN"]:
    sys.exit(0) # nothing to do for certificates of other services

agent.run_helper('systemctl', '--user', 'reload', 'myservice', check=True)
```

The handler script reloads `myservice` when the certificate of the service
is obtained, renewed, or changed by other means.

While the Let's Encrypt certificate request is pending the service could
be temporarily started with a self-signed certificate, for instance by
using the `openssl` command. Generating such certificate is out of the
scope of this tutorial.

In any case, before starting the service, the module can retrieve the
certificate from Redis. Certificates are saved in a Redis HASH key like
`module/traefik1/certificates/www.example.org`. See [Redis DB
reference]({{site.baseurl}}/core/database/#moduletraefikx) for details.

For example

```bash
mtraefik=$(redis-exec GET "node/${NODE_ID}/default_instance/traefik")
redis-exec HGET "module/${mtraefik}/certificate/${SERVICE_FQDN}" key | base64 -d > server.key
redis-exec HGET "module/${mtraefik}/certificate/${SERVICE_FQDN}" cert | base64 -d > server.crt
```


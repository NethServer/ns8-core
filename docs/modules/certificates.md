---

layout: default
title: TLS certificates
nav_order: 80
parent: Modules

---

TLS certificates
================

X.509 certificates for TLS are managed by the Traefik instance running on
the local node. For more details, see
[Proxy and certificates](../../core/proxy_certificates) in the core
documentation.

Web applications
----------------

Web applications should not handle TLS certificates directly, as the
Traefik HTTP proxy listens on TCP port 443 for HTTPS requests. HTTP
backends must be configured only for plaintext HTTP connections.

If the module needs to set up a hostname-based HTTP route, it is
possible to instruct Traefik to obtain a TLS certificate for it in a
single action.

```python
agent.tasks.run(
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g., module/traefik1
    action='set-route',
    data={
        'instance': os.environ['MODULE_ID'],
        'url': f'http://127.0.0.1:{os.environ["TCP_PORT"]}',
        'host': host_fqdn,
        'lets_encrypt': True, # request certificate for host_fqdn
        'http2https': True, # redirect http:// URL scheme to https://
    },
)
```

Note that `instance` must be a unique route identifier. If the module
uses multiple HTTP routes, add a prefix to `MODULE_ID` to distinguish them.

TLS for other applications
--------------------------

Applications that require `.pem`, `.key`, `.crt`, or similar files for
their configuration must be granted Traefik's `fulladm` role to run
privileged actions. Add the following authorization request to your app's
image labels:

     org.nethserver.authorizations = traefik@node:fulladm

A certificate request for Let's Encrypt can be issued using the
`set-certificate` action.

The application can then listen for the `certificate-changed` event. This
event is triggered when a TLS certificate is uploaded via the UI or
obtained/renewed from Let's Encrypt.

Upon receiving a `certificate-changed` event, invoke the `get-certificate`
action to obtain the PEM-encoded key and certificate.

For example, to issue a Let's Encrypt certificate request for
`SERVICE_FQDN` in the background, run:

```python
agent.tasks.run_nowait( # _nowait runs in the background
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g., module/traefik1
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

If you prefer a synchronous call, run:

```python
response = agent.tasks.run(
    agent_id=agent.resolve_agent_id('traefik@node'), # e.g., module/traefik1
    action='set-certificate',
    data={
        "fqdn": os.environ["SERVICE_FQDN"],
        "sync_timeout": 60, # default timeout is 30
    },
)
# Check the action exit code:
if response['exit_code'] != 0:
    # decide what to do if the certificate was not obtained
```

To handle the event, create an executable script under
`$AGENT_INSTALL_DIR/events/certificate-changed`. The event handler:

1. Checks if the event is relevant to the running service.
2. Restarts the service.

The following Python code implements these steps using the `agent`
library and the `get-certificate` command:

```python
import agent
import json
import sys
event = json.load(sys.stdin)
myname = os.environ["SERVICE_FQDN"]

if not agent.certificate_event_matches(event, myname):
    sys.exit(0) # ignore event if not relevant

agent.run_helper('systemctl', '--user', 'try-restart', 'myservice')
```

It is recommended to run `get-certificate` before the service starts or is
fully restarted from the handler, to ensure the correct certificate is
used by the service. There can be different ways to achieve this behavior,
but it is recommended to define a separate Systemd unit, to avoid repeated
`get-certificate` invocations if the service enters a crash-loop.

This is an example of `get-certificate.service` definition that also
generates a PEM file with both private key and the complete certificate
chain:

```ini
[Unit]
Description=Get TLS certificate from Traefik

[Service]
Type=oneshot
EnvironmentFile=%S/state/environment
WorkingDirectory=%S/state
SyslogIdentifier=%N
ExecStart=-mkdir -vp tls-certs
ExecStart=-runagent get-certificate --cert-file=tls-certs/server.pem --key-file=tls-certs/server.key $POSTFIX_HOSTNAME
ExecStart=-bash -c 'cd tls-certs ; umask 177 ; cat server.key server.pem > fullchain.pem ; touch -r server.pem fullchain.pem'
```

Add the following lines to `myservice.service` Systemd unit, to run
`get-certificate.service` before the service is started or restarted:

```ini
[Unit]
# ...existing lines
Wants=get-certificate.service
After=get-certificate.service
```

To avoid a full service restart when the certificate changes, the
`myservice.service` unit can provide a reload operation. In this case,
modify the `certificate-changed` event handler as follows.

1. Checks if the event is relevant to the running service.
2. Obtains the certificate and its private key for the service FQDN.
3. If the certificate was changed, reloads the running service.

```python
import agent
import json
import sys
event = json.load(sys.stdin)
myname = os.environ["SERVICE_FQDN"]

if not agent.certificate_event_matches(event, myname):
    sys.exit(0) # ignore event if not relevant

# Obtain the changed certificate and its private key
get_proc = agent.run_helper('get-certificate', '--cert-file=service.pem', '--key-file=service.key', myname)

# Reload the running service if the certificate differs
if get_proc.returncode == 0:
    agent.run_helper('systemctl', '--user', 'reload', 'myservice')
```

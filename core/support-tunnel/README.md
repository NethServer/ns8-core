# Support Tunnel Client

WebSocket-based remote support client, started on demand by the cluster
administrator.

## Binary

The `tunnel-client-linux-amd64` binary is a build artifact from the
[MY project](https://github.com/NethServer/my) (`services/support/cmd/tunnel-client`).

To build it:

```bash
cd /path/to/my/services/support
make build-tunnel-client
cp build/tunnel-client /path/to/ns8-core/core/support-tunnel/tunnel-client-linux-amd64
cd /path/to/ns8-core/core/support-tunnel
sha256sum tunnel-client-linux-amd64 > tunnel-client-linux-amd64.sha256
```

`build-image.sh` fails if the binary or its `.sha256` file is missing or
the checksum does not match.

## Configuration

`start-tunnel-client` writes `/var/lib/nethserver/support-tunnel/support.env`
(mode 0600, directory 0700); the client reads everything from that
environment, no secret is passed on the command line.

Credentials default to the my subscription (`cluster/subscription`
`system_id` and `auth_token`, provider `nsent`). The optional
`cluster/support_tunnel` hash overrides them and sets other options:

- `url`: `wss://` tunnel URL (default `wss://support.my.nethesis.it/api/tunnel`, QA `wss://support.qa.my.nethesis.it/api/tunnel`)
- `key`, `secret`: system credentials, both required to override
- `exclude_patterns`: comma-separated globs of services to hide

## Ephemeral accounts cleanup

The accounts created for a session are recorded in
`/var/lib/nethserver/support-tunnel/users-state.json`. They are removed by
the client on stop, by `tunnel-client --cleanup-only` in `ExecStopPost`
if the client was killed, and at boot by `support-tunnel-cleanup.service`
(enabled by the first `start-tunnel-client`).

## Plugin directories

- `users.d/` - Ephemeral user provisioning scripts (setup/teardown lifecycle)
- `diagnostics.d/` - System health check scripts

## Coexistence with don

The tunnel runs alongside the OpenVPN-based "don" support (`support.service`),
each with its own service, state and support accounts:

- don: `/var/lib/nethserver/support/`, `support.service`, `start-support-session`
- tunnel: `/var/lib/nethserver/support-tunnel/`, `support-tunnel.service`, `start-tunnel-client`

During the beta the tunnel follows don: `start-support-session` and
`stop-support-session` also start and stop it, and `support-tunnel.service`
is `PartOf=support.service`, so it stops whenever don does (expiry included).
A tunnel failure never makes the don session fail. `start-tunnel-client` and
`stop-tunnel-client` still drive the tunnel alone.

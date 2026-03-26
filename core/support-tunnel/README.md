# Support Tunnel Client

WebSocket-based remote support system that replaces the OpenVPN-based "don" support.

## Binary

The `tunnel-client-linux-amd64` binary is a build artifact from the
[MY project](https://github.com/NethServer/my) (`services/support/cmd/tunnel-client`).

To build it:

```bash
cd /path/to/my/services/support
make build-tunnel-client
cp build/tunnel-client /path/to/ns8-core/core/support-tunnel/tunnel-client-linux-amd64
```

## Plugin directories

- `users.d/` - Ephemeral user provisioning scripts (setup/teardown lifecycle)
- `diagnostics.d/` - System health check scripts

## Coexistence with don

This system coexists with the existing OpenVPN-based "don" support (`support.service`).
The two systems use separate paths and service names:
- don: `/var/lib/nethserver/support/`, `support.service`
- tunnel: `/var/lib/nethserver/support-tunnel/`, `support-tunnel.service`

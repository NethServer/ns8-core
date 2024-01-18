# support

This container runs an OpenVPN client that connects with a remote support
system.

The server must provide a certificate matching `VPN_CERT_COMMON_NAME` and
signed by `VPN_CERT_FILE`.

## Environment

- `VPN_PEER_HOST`
- `VPN_PEER_PORT`
- `VPN_USERNAME`
- `VPN_PASSWORD`
- `VPN_CERT_FILE`
- `VPN_CERT_COMMON_NAME`
- `VPN_INACTIVE`

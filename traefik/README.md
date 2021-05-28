# Traefik

This module implements a proxy for web applications using [Traefik](https://doc.traefik.io/traefik/).

The module exposes 2 actions:
- `create-host`
- `delete-host`

## create-host

This action creates a virtual host and a proxy pass to expose web applications instances.
The action takes 5 arguments:
- `instance`: the instance name, which is unique inside the cluster
- `url`:the backend target URL
- `host`: a fully qualified domain named as virtual host
- `lets_encrypt`: can be `true` or `false`, if set to `true` request a valid Let's Encrypt certificate
- `http2https` can be `true` or `false`, if set to `true` HTTP will be redirect to HTTPS

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": {"instance": "module1", "url": "http://127.0.0.0:2000", "host": "module.example.org", "lets_encrypt": true, "http2https": true}}'
```

## delete-host

This action delets an existing virtual host. It can be used when removing a module instance.
The action takes 1 parameter:
- `instance`: the instance name

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": {"instance": "module1"}}
```

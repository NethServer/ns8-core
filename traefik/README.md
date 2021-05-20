# Traefik

This module implements a proxy for web applications using [Traefik](https://doc.traefik.io/traefik/).

The module exposes 2 actions:
- `create-host`
- `delete-host`

## create-host

This action creates a virtual host and a proxy pass to expose web applications instances.
The action takes 5 arguments:
- `N`: the instance name, which is unique inside the cluster
- `URL`:the backend target URL
- `HOST`: a fully qualified domain named as virtual host
- `LETS_ENCRYPT`: can be 0 or 1, if set to 1 request a valid Let's Encrypt certificate
- `FORCE_HTTPS` can be 0 or 1, if set to 1 HTTP will be redirect to HTTPS

All parameters should be separated by a space and terminated with a whiteline character.

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": "dokuwiki1 http://127.0.0.1:20001 dokuwiki.test.local 1 1\n"}'
```

## delete-host

This action delets an existing virtual host. It can be used when removing a module instance.
The action takes 1 parameter:
- `N`: the instance name

The parameter should be terminated with a whiteline character.

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "dokuwiki1\n"}
```

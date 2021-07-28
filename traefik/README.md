# Traefik

This module implements a proxy for web applications using [Traefik](https://doc.traefik.io/traefik/).

The module exposes 2 actions:
- `set-route`
- `delete-route`
- `set-certificate`
- `delete-certificate`

## set-route

This action creates HTTP routes based on a combination of host and path, is possible to define three type
of rules:

* only `host`: These rules will capture all the requests directed to a specific host
* `host` and `path`: These rules will capture all the requests directed to a specific combination of host and path prefix
* only `path`: These rules will capture all the requests directed to a specific path prefix, regardless of the host.

This is the priority of the rules type evaluation (top-down):

1. only `host`
1. `host` and `path`
1. only `path`

### Parameters

- `instance`: the instance name, which is unique inside the cluster, mandatory
- `url`: the backend target URL, mandatory
- `host`: a fully qualified domain name as virtual host
- `path`: a path prefix, the matching evaluation will be performed whit and without the trailing slash, eg `/foo` will match `/foo` and `/foo/*`, also `/foo/` will match `/foo` and `/foo/*`
- `lets_encrypt`: can be `true` or `false`, if set to `true` request a valid Let's Encrypt certificate, mandatory
- `http2https` can be `true` or `false`, if set to `true` HTTP will be redirect to HTTPS, mandatory

### Examples

Only `host`
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "module1", "url": "http://127.0.0.0:2000", "host": "module.example.org", "lets_encrypt": true, "http2https": true}}'
```

`host` and `path`
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "module1", "url": "http://127.0.0.0:2000", "host": "module.example.org", "path": "/foo", "lets_encrypt": true, "http2https": true}}'
```
Only `path`

```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "module1", "url": "http://127.0.0.0:2000", "path":"/foo", "lets_encrypt": true, "http2https": true}}'
```

## delete-route

This action delets an existing route. It can be used when removing a module instance.
The action takes 1 parameter:
- `instance`: the instance name

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-route", "data": {"instance": "module1"}}
```

## set-certificate

This action explicitly requests a Let's Encrypt certificate. It can be used when there is no hostname (or hostname + path) route configured on traefik module or if the service is not make accessible via traefik.

The action takes 1 parameter:
- `fqdn`: the fqdn of the requested certificate

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-certificate", "data": {"fqdn": "'$(hostname -f)'"}}'
```
## delete-certificate

This action deletes an existing route used for explicit request a certificate.

NB. The certificate will **not** actually be removed from traefik and if the conditions will remain in place it will be renewed.

The action takes 1 parameter:
- `fqdn`: the fqdn of the requested certificate

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-certificate", "data": {"fqdn": "'$(hostname -f)"}}'
```

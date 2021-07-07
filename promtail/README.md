# Promtail

Start and configure a Promtail instance.
The module use the [Promtail official docker image](https://github.com/grafana/loki/releases).

## Install

Instantiate the module, example:
```
add-module ghcr.io/nethserver/promtail:latest 1
```

The output of the command will return the instance name.
Output example:
```
{"module_id": "promtail1", "image_name": "promtail", "image_url": "ghcr.io/nethserver/promtail:latest"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

As default, the Loki server used will be the default instance of the cluster.

## Configure

Let's assume that the Promtail instance is named `promtail1`.

Then launch `configure-module`, by setting the following parameters:
- `loki_instance`, *string*: Module id of the Loki server.
- `loki_url`, *string*: Address of the remote Loki server.
- `auth`, *object*: Authentication credentials, only valid if `loki_url` is used:
  - `token`, *string*: Bearer token to send to the server.
  - `basic`, *object*: Basic auth username and password for authenticate to the server.
    - `username`, *string*: Username.
    - `Password`, *string*: Password.

All parameters must be set inside the `data` field as a JSON object.

Example:
```
redis-cli LPUSH module/promtail1/tasks '{"id":"'$(uuidgen)'","action":"configure-module","data": {"loki_instance":"loki1"}}'
```

```
redis-cli LPUSH module/promtail1/tasks '{"id":"'$(uuidgen)'","action":"configure-module","data": {"loki_url":"https://example.com/loki/api/v1/push"}}'
```

```
redis-cli LPUSH module/promtail1/tasks '{"id":"'$(uuidgen)'","action":"configure-module","data": {"loki_url": ""https://example.com/loki/api/v1/push", "auth": {"token": "XXXXXXXXXXXXXXXXXX"}}}'
```

```
redis-cli LPUSH module/promtail1/tasks '{"id":"'$(uuidgen)'","action":"configure-module","data": {"loki_url": "http://10.135.0.3:3100/loki/api/v1/push", "auth": {"basic": {"user": "user", "password": "password"}}}}'
```

## Uninstall

To uninstall the instance:
```
remove-module --no-preserve promtail1
```

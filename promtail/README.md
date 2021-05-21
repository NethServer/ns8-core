# Promtail

Start and configure a Promtail instance.
The module use the [Promtail official docker image](https://github.com/grafana/loki/releases).

## Install

Instantiate the module:
```
add-module promtail 1
```

The output of the command will return the instance name.
Output example:
```
{"rootfull": true, "mid": "promtail1"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

## Configure

Let's assume that the Promtail instance is named `promtail1`.

Then launch `configure-module`, by setting the following parameters:
- `loki_url`, *string*: Address of the remote Loki server.

All parameters must be set inside the `data` field as a JSON object.

Example:
```
redis-cli LPUSH module/promtail1/tasks '{"id":"'$(uuidgen)'","action":"configure-module","data":"{\"loki_url\":\"https://example.com/loki/api/v1/push\"}"}'
```

## Uninstall

To uninstall the instance:
```
remove-module --no-preserve promtail1
```

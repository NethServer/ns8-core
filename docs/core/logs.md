---
layout: default
title: Logs
nav_order: 9
parent: Core
---

# Logs

Almost everything is logged inside the system journal.
Recent logs are available using `journalctl` and services can be inspected using `systemctl` command.

As root use `journalctl` to see messages from agents, rootfull and rootless modules.
As rootless UNIX user (eg. `traefik1`), use `journalctl --user` to see messages only from systemd user session.

By default a Loki instance is installed inside the leader node, it
collects the logs from all cluster nodes. A rootfull Promtail container
runs as a node service on all nodes, including the leader one. It sends
all logs to the Loki server.

Logs of nodes, modules and the whole cluster can be read from the Logs
page, or with the `api-server-logs` command. The following command prints
the last lines of log produced by module `traefik1` and waits for further
lines, like a `tail -f`:

    api-server-logs logs -e module -n traefik1

See `api-server-logs -h` for more command invocation styles.

It is also possible to query the logs of all nodes and modules with the
`logcli` command.

Get the list of existing node identifiers:

    # logcli labels node_id -q
    1
    2

Get the list of modules:

    # logcli labels module_id -q
    ldapproxy1
    loki1
    traefik1

Example log search for module `ldapproxy1`:

<!-- {% raw %} -->

    # logcli query -q --no-labels '{module_id="ldapproxy1"} | json | line_format "{{.MESSAGE}}"'

<!-- {% endraw %} -->

The string wrapped by `''` is the log query. It is expressed in
[LogQL](https://grafana.com/docs/loki/latest/logql/log_queries/).

## Forwarding

On an active Loki instance, it is possible to enable two forwarding services. 
When forwarders are active and the leader node changes, they are automatically migrated to the new leader node.

By default, these forwarders are disabled.

### Syslog server

This service can be enabled either from the cluster admin interface or through the API using the following examples:
```
# api-cli run module/loki1/set-syslog-forwarder --data '{"active": true, "address": "127.0.0.1", "port": "514", "protocol": "udp", "format": "rfc3164", "start_time": "2024-01-01T00:00:00.00Z", "filter": "security"}'
```

**Configuration parameters**
- **active**: Boolean to enable or disable the forwarder (e.g., `true` or `false`).
- **address**: The IP address of the Syslog server.
- **port**: The port on which the Syslog server listens.
- **protocol**: Protocol to use (e.g., `udp` or `tcp`).
- **format**: Log format (e.g., `rfc3164` or `rfc5424`).
- **start_time**: The time from which logs should be forwarded. (`ISO 8601`)
- **filter**: The filter to apply to the logs. The only available filters are `"security"` and `""` (all logs).

It can be also stopped through API:
```
# api-cli run module/loki1/set-syslog-forwarder --data '{"active": false}'
```

### Cloud Log Manager

Cloud Log Manager forwarder is available only with a subscription, and if the subscription expires or is removed, the forwarder is automatically stopped.

Similarly, this service can be enabled from the cluster admin interface or through the API with examples like the following:
```
# api-cli run module/loki1/set-clm-forwarder --data '{"active": true, "address": "https://nar.nethesis.it", "tenant": "example", "start_time": "2024-07-12T00:00:00.000000Z"}'
```

**Configuration parameters**
- **active**: Boolean to enable or disable the forwarder (e.g., `true` or `false`).
- **address**: The URL of the Cloud Log Manager.
- **tenant**: The tenant identifier.
- **start_time**: The time from which logs should be forwarded. (`ISO 8601`)
- **filter**: The filter to apply to the logs. The only available filters are `"security"` and `""` (all logs). By default, the `security` filter is applied and the parameter is optional.

It can be also stopped through API:
```
# api-cli run module/loki1/set-clm-forwarder --data '{"active": false}'
```

When Loki is installed, it generates a unique identifier that is used on the [Cloud Log Manager](https://naradmin.nethesis.it/).
You can find this identifier on Cloud Log Manager settings panel or using the following command:
```
# redis-cli hget module/$(redis-cli --raw get cluster/default_instance/loki)/environment CLOUD_LOG_MANAGER_HOSTNAME
```

### Troubleshooting

- Check the status of forwarders using loki agent:
```
runagent -m loki1
systemctl --user status *forwarder.service
```

See also:
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/)
- [ns8-loki](https://github.com/NethServer/ns8-loki) module repository

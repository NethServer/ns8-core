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

    # logcli  query -q --no-labels '{module_id="ldapproxy1"} | json | line_format "{{.MESSAGE}}"'

The string wrapped by `''` is the log query. It is expressed in
[LogQL](https://grafana.com/docs/loki/latest/logql/log_queries/).

See also:
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/)
- [ns8-loki](https://github.com/NethServer/ns8-loki) module repository

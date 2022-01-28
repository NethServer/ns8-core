---
layout: default
title: Logs
nav_order: 9
parent: Core
---

# Logs

Almost everything is logged inside journalctl.
Recent logs are available using `journalctl` and services can be inspected using `systemctl` command.
As root use `journalctl` to see messages from agents, rootfull and rootless modules.
As rootless UNIX user (eg. `traefik1`), use `journalctl --user` to see messages only from systemd user session.

By default, [Grafana Loki](https://grafana.com/oss/loki/) is installed inside the leader node, it collects the logs
from all cluster nodes.
A rootfull [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) container runs on all nodes,
including the leader one. It sends all logs to the Loki server.

From the leader node, it is possible to query the logs of all nodes.
First, get the list of existing nodes:
```
# logcli labels nodename -q
bullseye
leader
```

Then, search for the logs of a node (`leader` in this example):
{% raw %}
```
# logcli  query -q --no-labels '{nodename="leader"} | json | line_format "{{.MESSAGE}}"'
2021-05-28T15:49:27Z Created slice cgroup user-libpod_pod_d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6.slice.
2021-05-28T15:49:27Z 2021-05-28 15:49:27.621079036 +0000 UTC m=+0.106351212 container create 6e51520a9fa4f63ac0f3dbbf89ef2ef075041dd90c5df952a35206a14691c654 (image=k8s.gcr.io/pause:3.2, name=d32e9a0f2237-infra)
2021-05-28T15:49:27Z 2021-05-28 15:49:27.62160088 +0000 UTC m=+0.106873062 pod create d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6 (image=, name=loki)
2021-05-28T15:49:27Z d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6
2021-05-28T15:49:27Z loki.service: Found left-over process 19008 (podman pause) in control group while starting unit. Ignoring.
2021-05-28T15:49:27Z This usually indicates unclean termination of a previous run, or service implementation deficiencies.
```
{% endraw %}

See also [Loki](https://github.com/NethServer/ns8-core/tree/main/loki) and [Promtail](https://github.com/NethServer/ns8-core/tree/main/promtail) READMEs.

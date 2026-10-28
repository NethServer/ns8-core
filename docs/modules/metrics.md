---
layout: default
title: Metrics
nav_order: 105
parent: Modules
---

# Metrics

A module can publish Prometheus metrics with a command that prints them.
No HTTP server and no open port are needed: the node serves them through
`node_exporter`, which the `metrics` module already scrapes.

Use it for values the module already knows, like the number of app active
users. A module that runs its own Prometheus exporter publishes it as
a target instead, with the `metrics-target-changed` event. See the
[ns8-metrics README](https://github.com/NethServer/ns8-metrics#provisioning-prometheus).

## The export-metrics command

Add an executable `imageroot/bin/export-metrics` to the module. Every 5
minutes the node runs it as:

```text
runagent -m <module_id> export-metrics
```

Example in Python, with the `prometheus_client` library included in
the core:

```python
#!/usr/bin/env python3

import os
import sys
import agent
import prometheus_client as prom

rdb = agent.redis_connect()
# A value the module keeps up to date, here in its own Redis key
users = rdb.hget(f"module/{os.environ['MODULE_ID']}/stats", "active_users") or 0

# Drop the <name>_created series of counters: it changes on every run
prom.disable_created_metrics()
registry = prom.CollectorRegistry()
active_users = prom.Gauge(
    "myapp_active_users",
    "Users active in the last 24 hours",
    registry=registry,
)
active_users.set(users)
sys.stdout.buffer.write(prom.generate_latest(registry))
```

Commands in other languages can print the text format directly, or use
the official client library of the language.

The node adds a `module` label to every sample:

```text
myapp_active_users{module="mail1"} 42
```

## Rules

| Topic | Rule |
|---|---|
| Output | Valid [Prometheus text format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format) on stdout, without timestamps. The node does not check it: review it with the module code. |
| Speed | Must finish within 30 seconds. Read values that are already there, for example in Redis. Compute expensive values in the module's own jobs. |
| Names | Prefix metric names with the application name, like `nethvoice_`. Never use `node_`, `ns8_`, `go_` or `process_`. |
| Labels | The node adds `module="<module_id>"`. Output that sets `module` or `module_id` is rejected. |
| Errors | Log to stderr. The node copies each line to its journal, prefixed with the module ID. |

## Failures

If the command fails, times out, or sets a reserved label:

- the node keeps the last output of the module;
- the error is logged:

  ```text
  journalctl -u refresh-module-metrics.service
  ```

The 30 seconds start when the command starts. When they expire, the node
sends `SIGKILL` to the process group of the command, which includes its
child processes. There is no `SIGTERM` first, so the command cannot
catch the signal or run cleanup code. Its stdout and stderr are
discarded. Do not leave state that needs cleanup: for example, write
files to a temporary name and rename them when they are complete.

If the output is not valid, `node_exporter` skips it. The other modules
are not affected. It logs the error and sets the node-wide
`node_textfile_scrape_error` gauge to `1`:

```text
journalctl -u node_exporter.service
```

When the command is removed, its metrics are removed at the next run.

## Alerts

A module ships its own alert rules in the Redis hash
`module/<module_id>/metrics_alert_rules`, one field per rule set. After
each change it emits the `metrics-alert-rules-changed` event, and
Prometheus reloads its rules without a restart. See
"metrics-alert-rules-changed event" in the
[ns8-metrics README](https://github.com/NethServer/ns8-metrics#metrics-alert-rules-changed-event).

- Prometheus copies the `module` label to `module_id`, the label module
  alert rules are scoped by. See "Per-module node series" in the
  [ns8-metrics README](https://github.com/NethServer/ns8-metrics#per-module-node-series).
- Module rules get the `module_id` matcher of their own module
  automatically.
- Give each rule a `severity` label, `warning` or `critical`, and the
  `summary_en`, `summary_it`, `description_en` and `description_it`
  annotations.
- Rules of a removed module are deleted automatically.
- The core defines no alert on module metrics.

Example: more active users than the license allows, in `license.yml`:

```yaml
alert: MyappTooManyUsers
expr: myapp_active_users > 100
for: 1h
labels:
  severity: warning
annotations:
  summary_en: Too many active users
  summary_it: Troppi utenti attivi
  description_en: Active users exceed the license limit
  description_it: Gli utenti attivi superano il limite della licenza
```

Publish it for `mail1`:

```text
redis-cli -x HSET module/mail1/metrics_alert_rules license < license.yml
redis-cli PUBLISH module/mail1/event/metrics-alert-rules-changed '{}'
```

## Collector

- Command: `refresh-module-metrics`, run by the node agent.
- Schedule: every 5 minutes, by `refresh-module-metrics.timer`.
  `node_exporter.service` wants the timer.
- Modules found as in the [volume
  metrics]({{site.baseurl}}/core/volume_metrics):
  - rootless: a user whose home holds `.config/state/agent.env`;
  - rootfull: a `/var/lib/nethserver/<name><number>` directory.
- Runs `export-metrics` only when `<install dir>/bin/export-metrics`
  exists and is executable. Up to 4 modules run at the same time, in
  random order, each after a random delay of up to 5 seconds.
- Output: `/var/lib/nethserver/node/state/node_exporter/module-<module_id>.prom`,
  one file per module, so invalid output does not hide the others.
  Written atomically, only when it changes.

Run it at once, without waiting for the timer:

```text
runagent -m node refresh-module-metrics
```

If another run is in progress, the command fails at once.

---
layout: default
title: Container metrics
nav_order: 18
parent: Core
---

# Container metrics

NS8 exposes per-container CPU, memory, PIDs, block I/O and network metrics
on each node's existing `node_exporter` endpoint. The metrics are produced
by `refresh-container-metrics.service`, a collector that runs on every
node, writes a Prometheus text file into `node_exporter`'s textfile
collector directory every 60 seconds, and is bound to the lifecycle of
`node_exporter.service` itself. No new port is opened and no new scrape
job is needed: the series appear on the node's `/metrics` endpoint like
any other `node_exporter` metric and are collected into the `metrics`
module's Prometheus on the leader node with the rest of the cluster's
metrics.

The collector reads cgroup v2 accounting files directly, instead of
querying the systemd unit that started the container. NS8 starts
containers with `podman run --detach --cgroups=no-conmon`, so the
systemd unit's own cgroup only accounts for the `conmon` supervisor
process, not the container payload (on a test node: 770 KB reported by
the unit vs. 165 MB in the container's real `libpod-<CID>.scope`
cgroup). The collector reads that scope cgroup instead.

## Metrics

Counters:

| Metric | Labels | Meaning |
|---|---|---|
| `ns8_container_cpu_seconds_total` | `module`, `container`, `mode="user\|system"` | Container CPU time spent, in seconds |
| `ns8_container_blkio_bytes_total` | `module`, `container`, `device`, `op="read\|write"` | Container block I/O transferred, in bytes |
| `ns8_container_blkio_ops_total` | `module`, `container`, `device`, `op="read\|write"` | Container block I/O operations |
| `ns8_container_network_receive_bytes_total` | `module`, `container`, `device` | Container bytes received |
| `ns8_container_network_transmit_bytes_total` | `module`, `container`, `device` | Container bytes transmitted |
| `ns8_container_network_receive_packets_total` | `module`, `container`, `device` | Container packets received |
| `ns8_container_network_transmit_packets_total` | `module`, `container`, `device` | Container packets transmitted |
| `ns8_container_oom_kills_total` | `module`, `container` | Processes killed by the container OOM killer |

Gauges:

| Metric | Labels | Meaning |
|---|---|---|
| `ns8_container_memory_usage_bytes` | `module`, `container` | Container memory usage |
| `ns8_container_memory_peak_bytes` | `module`, `container` | Container peak memory usage |
| `ns8_container_memory_limit_bytes` | `module`, `container` | Container memory limit |
| `ns8_container_memory_swap_bytes` | `module`, `container` | Container swap usage |
| `ns8_container_memory_anon_bytes` | `module`, `container` | Container anonymous memory |
| `ns8_container_memory_file_bytes` | `module`, `container` | Container page cache memory |
| `ns8_container_pids` | `module`, `container` | Processes running in the container |
| `ns8_container_pids_limit` | `module`, `container` | Container process limit |
| `ns8_container_start_time_seconds` | `module`, `container` | Container creation time |
| `ns8_container_info` | `module`, `container`, `id`, `image`, `unit`, `rootless` | Container metadata, always `1` |
| `ns8_container_collector_duration_seconds` | none | Duration of the last container metrics collection |
| `ns8_container_collector_skipped_containers` | none | Containers that could not be collected in the last collection |
| `ns8_container_collector_last_success_timestamp_seconds` | none | Unix timestamp of the last successful collection |

## Labels

`container` is the container name, and `module` is the id of the module
that owns it. Core containers that are not part of an installed module,
such as `redis`, `promtail`, `node_exporter` and `rclone-gateway`, are
reported with `module="core"`.

A container whose owner cannot be determined is reported with
`module="unknown"` instead. That happens while a container is exiting,
when one was started outside a service unit (`podman run` by hand), or
when a rootless container belongs to a local account that is not a
module. Module users are recognised by the agent state file
`~/.config/state/agent.env` (the same test `runagent` uses), not by
owning a container storage directory — any account that has ever run
rootless podman has one of those too. `module="unknown"` is kept
separate from `module="core"` so a failed attribution is never mistaken
for a core container.

Modules that run their containers in a pod also report the pod's infra
container, named `<pod>-infra`. It is a real container with its own
cgroup, and it is attributed to the module that owns the pod.

`ns8_container_start_time_seconds` is podman's own creation timestamp
for the container. NS8 units start their containers with `--replace`, so
a container is recreated on every restart and its creation time is also
its start time. The cgroup directory mtime is used only as a fallback,
when the index carries no timestamp — it's weaker, since a cgroup's
mtime also advances whenever a child cgroup appears or goes away.

The container id (truncated to 12 characters) appears only on
`ns8_container_info`, not on every series: a container gets a new id on
each restart, so putting it everywhere would churn the label set and
break `rate()` across restarts. Use the stable `container` name to
select or aggregate a container across restarts.

## Block I/O

`ns8_container_blkio_bytes_total` and `ns8_container_blkio_ops_total`
depend on the `io` cgroup controller being delegated to the container's
cgroup. Rootfull containers, which run under `machine.slice`, already
have the `io` controller delegated, so block I/O is reported
immediately.

Rootless modules run under `user@<uid>.service`, which by default
delegates only `memory` and `pids`. The core update that ships this
feature installs a drop-in,
`/etc/systemd/system/user@.service.d/50-ns8-delegate.conf`, that adds
`io` gives disk stats access to rootless modules.

The drop-in applies to every local user manager on the node, not only
NethServer module users — a `user@.service.d` drop-in can't be scoped
to particular user ids. Root's own manager is the exception, masked by
an empty `/etc/systemd/system/user@0.service.d/50-ns8-delegate.conf`.

## Network

`ns8_container_network_*` is reported once per network namespace, not
once per container.

Every container in a pod shares the pod's network namespace, so they
all read the same interface counters. To avoid multiplying the pod's
traffic by its container count in a `sum by (module)`, it's reported
once, with `container` set to `<pod>-infra` (podman's name for the
pod's infra container). Other pod members carry no
`ns8_container_network_*` series. Group by `module`, or select
`<pod>-infra`, to get a pod's traffic.

Beyond that, the family is reported only for a container that has its
own network namespace. Most NS8 containers run with `--network=host`
and share the node's namespace instead, so their traffic shows up in
the node's own `node_network_*` metrics, not here.

Loopback (`lo`) is excluded: it's the container talking to itself, not
real traffic in or out, and on a test node it was 98% of all bytes
received — large enough to hide everything else. When a container's
network family is absent from `ns8_container_network_*`, that means
"not applicable", not "zero traffic".

More generally, the collector omits a metric rather than reporting it as
zero whenever the value could not be measured, for example when a
cgroup file is missing. Prometheus treats an absent series and a zero
series differently, and reporting zero would assert a measurement that
was never actually taken.

## Example query

Top 5 modules by CPU usage over the last 5 minutes:

```
topk(5, sum by (module) (rate(ns8_container_cpu_seconds_total[5m])))
```

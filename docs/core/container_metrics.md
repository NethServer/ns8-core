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
collector directory every 30 seconds, and is bound to the lifecycle of
`node_exporter.service` itself. No new port is opened and no new scrape
job is needed: the series appear on the node's `/metrics` endpoint like
any other `node_exporter` metric and are collected into the `metrics`
module's Prometheus on the leader node with the rest of the cluster's
metrics.

The collector reads cgroup v2 accounting files directly, instead of
querying the systemd unit that started the container. NS8 starts
containers with `podman run --detach --cgroups=no-conmon`, so the
systemd unit's own cgroup only accounts for the `conmon` supervisor
process, not the container payload. On the reference development node
the unit reported 770 KB of memory while the container's own
`libpod-<CID>.scope` cgroup reported 165 MB. The collector reads the
scope cgroup instead, which carries the real numbers.

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
| `ns8_container_start_time_seconds` | `module`, `container` | Container start time |
| `ns8_container_info` | `module`, `container`, `id`, `image`, `unit`, `rootless` | Container metadata, always `1` |
| `ns8_container_collector_duration_seconds` | none | Duration of the last container metrics collection |
| `ns8_container_collector_last_success_timestamp_seconds` | none | Unix timestamp of the last successful collection |

## Labels

`container` is the container name, and `module` is the id of the module
that owns it. Core containers that are not part of an installed module,
such as `redis`, `promtail`, `node_exporter` and `rclone-gateway`, are
reported with `module="node"`.

A container whose owning systemd unit cannot be determined is reported
with `module="unknown"` instead. That happens while a container is
exiting, or when one was started outside a service unit, for example by
running `podman run` by hand. It is kept distinct from `module="node"`
so that a failed attribution is never mistaken for a core container.

Modules that run their containers in a pod also report the pod's infra
container, named `<pod>-infra`. It is a real container with its own
cgroup, and it is attributed to the module that owns the pod.

The container id appears only on `ns8_container_info`, and is truncated
to 12 characters. It is deliberately not a label on every series: a
container is recreated on every module restart and gets a new id each
time, so putting the id on every series would churn the label set on
every restart and break `rate()` across that restart. The stable
`container` name is what should be used to select or aggregate a
container across restarts.

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
`cpu cpuset io` to the delegated controller set. That drop-in only takes
effect the next time a rootless module's user manager restarts, which
happens on the module's own restart or on a node reboot. Until then,
`ns8_container_blkio_*` is absent for that module's containers, while
CPU, memory and PID metrics are complete. This is expected, not a fault:
an update deliberately does not bounce every rootless container on a
node just to enable a monitoring feature.

The drop-in applies to every local user manager on the node, not only to
NethServer module users: a `user@.service.d` drop-in cannot be scoped to
particular user ids. To roll it back, remove
`/etc/systemd/system/user@.service.d/50-ns8-delegate.conf`, run
`systemctl daemon-reload`, then restart the affected user managers or
reboot the node. Block I/O metrics disappear again for rootless modules
when it is removed.

## Network

`ns8_container_network_*` is reported only for a container that has its
own network namespace. Most NS8 containers run with `--network=host`
and share the node's network namespace, so their traffic is already
covered by the node's own `node_network_*` metrics instead of a
per-container series. When a container's network family is absent from
`ns8_container_network_*`, that means "not applicable to this
container", not "zero traffic".

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

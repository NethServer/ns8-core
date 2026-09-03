---
layout: default
title: Volume metrics
nav_order: 19
parent: Core
---

# Volume metrics

NS8 exposes per-volume and per-module disk usage on each node's existing
`node_exporter` endpoint. The metrics are produced by
`refresh-volume-metrics.service`, a collector that runs on every node,
writes a Prometheus text file into `node_exporter`'s textfile collector
directory every 15 minutes, and is bound to the lifecycle of
`node_exporter.service` itself. Like the [container
metrics](container_metrics.md), the series appear on the node's
`/metrics` endpoint and are collected into the `metrics` module's
Prometheus on the leader node, with no new port and no new scrape job.

The collector walks the directory trees, because there's nothing
cheaper to read: podman volumes are plain directories, and the root
filesystem is mounted `noquota`, so a size can only be obtained by
visiting every inode. That walk is expensive on a node holding mail or
file shares, which is why it runs as its own unit on a slow interval
instead of another pass inside `refresh-container-metrics.service` — a
slow walk never delays the container samples, and a failure in one
collector doesn't take the other down. The unit also runs with
`Nice=10` and `IOSchedulingClass=idle`, so a monitoring pass never
makes a busy server feel slow.

Sizes are the space actually occupied on disk, `st_blocks * 512`, not
the apparent file size, and a hardlinked file is charged once per
module. Both match `du -s -B1`; `ns8_volume_files` matches
`du -s --inodes`.

The image-store walk stays on one filesystem, like `du -x` (compare
against `du -sx -B1`). It skips a running rootfull container's union
view, mounted inside the store at `overlay/<layer>/merged` — counting
it would re-count the whole container rootfs on top of the layers it's
built from, and hardlink dedupe doesn't catch it since lower-layer
files have a link count of one. The volume and state walks do cross
filesystems, so a volume backed by its own mount is still measured.

## Metrics

All gauges:

| Metric | Labels | Meaning |
|---|---|---|
| `ns8_volume_size_bytes` | `module`, `volume` | Volume disk usage |
| `ns8_volume_files` | `module`, `volume` | Filesystem objects in the volume, directories included |
| `ns8_module_volumes_bytes` | `module` | Named volume disk usage of the module |
| `ns8_module_images_bytes` | `module` | Container image store disk usage |
| `ns8_module_state_bytes` | `module` | Module state and configuration disk usage |
| `ns8_module_total_bytes` | `module` | Module disk usage, volumes and images included |
| `ns8_volume_collector_duration_seconds` | none | Duration of the last volume metrics collection |
| `ns8_volume_collector_skipped_modules` | none | Modules that could not be collected in the last collection |
| `ns8_volume_collector_last_success_timestamp_seconds` | none | Unix timestamp of the last successful collection |

## Buckets

Each module is walked once, and the result is split into three
non-overlapping buckets, reported per module by
`ns8_module_volumes_bytes`, `ns8_module_images_bytes` and
`ns8_module_state_bytes`. `ns8_module_total_bytes` is their sum, emitted
because adding three families up in PromQL is awkward.

`ns8_module_volumes_bytes` is the same figure as
`sum by (module) (ns8_volume_size_bytes)`, kept as its own series so
that a per-module view needs no aggregation.

Hardlink accounting is per module: each module's buckets share one set
of visited inodes, so a file with several names inside one module is
charged once, and two modules that share a hardlinked file are each
charged for it. Summing `ns8_module_total_bytes` across modules can
therefore exceed the space actually occupied, which is the same
behaviour as running `du -s` on each module separately.

| Bucket | Rootless module | Rootfull module |
|---|---|---|
| volumes | `~/.local/share/containers/storage/volumes/<volume>` | `/var/lib/containers/storage/volumes/<volume>` |
| images | `~/.local/share/containers/storage`, volumes excluded | `/var/lib/containers/storage`, volumes excluded |
| state | the home directory, container storage excluded | `/var/lib/nethserver/<module>` |

A rootless module is any local user whose home holds the agent state
file `.config/state/agent.env` (the same test `runagent` uses), not
just any owner of a `.local/share/containers/storage` directory —
an administrator who has ever run rootless podman owns one of those
too. This test also keeps a module counted while all its containers
are stopped, unlike a discovery based on running cgroups.

A rootfull module is a `/var/lib/nethserver/<name><number>` directory.
The core's own directories -- `cluster`, `node`, `api-server` -- carry
no trailing digits, so the shape alone tells them apart. Those three are
not modules, but their contents are real disk usage: they are summed
into `ns8_module_state_bytes{module="core"}`.

## Rootfull attribution

Rootfull modules share a single storage root, which constrains what can
be attributed to whom.

The image store is one overlay tree with shared layers, so no rootfull
module can be charged for its own part of it. It is reported once, as
`ns8_module_images_bytes{module="core"}`, rather than split by a guess.
Rootfull modules therefore have no `ns8_module_images_bytes` series of
their own, and their `ns8_module_total_bytes` covers their volumes and
state only.

Rootfull volumes are attributed by name, longest match first:
`crowdsec1-data` belongs to `crowdsec1`. Podman's volume-to-container
link lives only in its own database, unreadable without libpod, so the
name is the only durable clue available. A volume matching no module
(`redis-data`, `rclone-cache`) belongs to a core container and is
reported with `module="core"`, same as the container metrics.

The core's own directories are not modules, so `core` has no
`ns8_module_state_bytes` series: its total is the shared image store
plus the volumes of the core containers.

Volume names are unique per module, not per node: `openldap1` and
`samba1` both own a volume named `data`. Always select on both labels.

## Missing values

As with the container metrics, the collector omits a metric rather than
reporting it as zero whenever a value could not be measured. Prometheus
treats an absent series and a zero series differently, and reporting
zero would assert a measurement that was never taken. A zero in
`ns8_volume_size_bytes` is a real measurement: an empty volume whose
directory fits inside its inode occupies no blocks at all.

If one of a module's trees can't be walked, that bucket is left out and
the module is counted once in `ns8_volume_collector_skipped_modules`.
The module's other buckets are still reported — this gauge is what
tells a partial pass from a complete one, since an absent
`ns8_module_state_bytes` alone could just mean the module has no state
tree. A pass that fails outright leaves the previous file in place, so
metrics go stale rather than disappearing; check
`ns8_volume_collector_last_success_timestamp_seconds` to tell the two
cases apart.

## Alerting

`ns8_volume_collector_skipped_modules` above zero means the last pass
could not measure something, and the module it belongs to is named in
the journal of `refresh-volume-metrics.service`.

`ns8_volume_collector_duration_seconds` approaching the 900-second
interval is the other signal worth watching. The walk has no deadline —
abandoning it would leave the largest module (the one whose size
matters most) either unmeasured or reported as a shrinking partial sum.
Instead, the collector bounds how much of the wall clock it may
occupy: it rests at least as long as the pass took, so a slow pass
stretches the sampling period rather than crowding the node, and the
walk never takes more than half the time. Figures stay complete, just
sampled less often than every 15 minutes; this gauge shows when that's
happening. The walk also yields to real work throughout, via
`Nice=10` and `IOSchedulingClass=idle`.

```
ns8_volume_collector_duration_seconds > 900
```

## Example queries

Top 5 modules by disk usage:

```
topk(5, ns8_module_total_bytes)
```

Volumes that grew more than 1 GB in the last day:

```
(ns8_volume_size_bytes - ns8_volume_size_bytes offset 1d) > 1e9
```

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

"""Collect per-container resource usage from cgroup v2.

Podman starts NS8 containers with ``--cgroups=no-conmon``, so the systemd unit
cgroup accounts only for the conmon supervisor. The container itself lives in a
``libpod-<CID>.scope`` cgroup, which is what this module reads.

Every entry point takes explicit filesystem roots so the module can be tested
against a synthetic tree.
"""

import glob
import json
import os
import os.path
import pwd
import re
import sys

DEFAULT_CGROUP_ROOT = "/sys/fs/cgroup"
DEFAULT_PROC_ROOT = "/proc"
DEFAULT_SYS_DEV_BLOCK = "/sys/dev/block"
ROOTFULL_STORAGE_ROOT = "/var/lib/containers/storage"

# Container cgroups turn up at several depths and under several shapes, so
# discovery is recursive rather than a set of fixed paths:
#
#   machine.slice/libpod-<cid>.scope                       rootfull
#   user@<uid>.service/user.slice/libpod-<cid>.scope       rootless
#   .../user-libpod_pod_<pod>.slice/libpod-<cid>.scope     rootless, in a pod
#   system.slice/<unit>.service/libpod-payload-<cid>       --cgroups=split
#
# A fixed-depth glob silently missed the last two, which made every
# pod-based module invisible.
_SCOPE_GLOB = os.path.join("**", "libpod-*")

_SCOPE_RE = re.compile(r"/libpod-(?:payload-)?([0-9a-f]{64})(?:\.scope)?$")
_UID_RE = re.compile(r"/user-(\d+)\.slice/")


def _read_text(path):
    """Read a small file, returning None when it does not exist."""
    try:
        with open(path) as fp:
            return fp.read()
    except OSError:
        return None


def discover_scopes(cgroup_root=DEFAULT_CGROUP_ROOT):
    """Return the live container scopes found under cgroup_root."""
    scopes = []
    for path in glob.glob(os.path.join(cgroup_root, _SCOPE_GLOB), recursive=True):
        match = _SCOPE_RE.search(path)
        if match is None or not os.path.isdir(path):
            continue
        uid_match = _UID_RE.search(path)
        scopes.append(
            {
                "cid": match.group(1),
                "path": path,
                "rootless": uid_match is not None,
                "uid": int(uid_match.group(1)) if uid_match else None,
            }
        )
    scopes.sort(key=lambda scope: scope["path"])
    return scopes


def read_containers_json(storage_root):
    """Map container id to name and image, from podman's storage index."""
    path = os.path.join(storage_root, "overlay-containers", "containers.json")
    try:
        with open(path) as fp:
            entries = json.load(fp)
    except (OSError, ValueError):
        return {}
    if not isinstance(entries, list):
        return {}
    index = {}
    for entry in entries:
        if not isinstance(entry, dict):
            continue
        cid = entry.get("id")
        if not cid:
            continue
        names = entry.get("names")
        if not isinstance(names, list):
            names = []
        try:
            metadata = json.loads(entry.get("metadata") or "{}")
        except (ValueError, TypeError):
            metadata = {}
        if not isinstance(metadata, dict):
            metadata = {}
        # Both values become label values in the exposition, so anything that
        # is not a string is treated as absent: a non-string reaching
        # escape_label() would raise inside render(), long after collect()
        # succeeded, and would then fail identically on every later cycle.
        name = names[0] if names else ""
        image = metadata.get("image-name", "")
        index[cid] = {
            "name": name if isinstance(name, str) else "",
            "image": image if isinstance(image, str) else "",
        }
    return index


def storage_roots(uids, rootfull_storage_root=ROOTFULL_STORAGE_ROOT, passwd_lookup=pwd.getpwuid):
    """Rootfull storage root first, then one per distinct module user."""
    roots = [rootfull_storage_root]
    for uid in sorted(set(uids)):
        try:
            home = passwd_lookup(uid).pw_dir
        except KeyError:
            continue
        roots.append(os.path.join(home, ".local/share/containers/storage"))
    return roots


NETHSERVER_ROOT = "/var/lib/nethserver"
RESERVED_MODULE_DIRS = frozenset(["cluster", "node", "api-server"])

_UNIT_PATTERNS = (
    "system.slice/*.service",
    "user.slice/user-*.slice/user@*.service/app.slice/*.service",
)

_CONMON_CID_RE = re.compile(r"\x00-c\x00([0-9a-f]{64})\x00")


def map_units(cgroup_root=DEFAULT_CGROUP_ROOT, proc_root=DEFAULT_PROC_ROOT):
    """Map container id to the systemd unit whose conmon supervises it."""
    units = {}
    for pattern in _UNIT_PATTERNS:
        for unit_path in glob.glob(os.path.join(cgroup_root, pattern)):
            unit = os.path.basename(unit_path)
            procs = _read_text(os.path.join(unit_path, "cgroup.procs"))
            if not procs:
                continue
            for pid in procs.split():
                try:
                    with open(os.path.join(proc_root, pid, "cmdline"), "rb") as fp:
                        cmdline = fp.read().decode("utf-8", "replace")
                except OSError:
                    continue
                if "conmon" not in cmdline.split("\x00")[0]:
                    continue
                match = _CONMON_CID_RE.search(cmdline)
                if match is not None:
                    units[match.group(1)] = unit
    return units


def list_module_ids(nethserver_root=NETHSERVER_ROOT):
    """Rootfull module ids, longest first so prefix matching is unambiguous."""
    try:
        names = os.listdir(nethserver_root)
    except OSError:
        return []
    ids = [
        name
        for name in names
        if name not in RESERVED_MODULE_DIRS
        and os.path.isdir(os.path.join(nethserver_root, name))
    ]
    ids.sort(key=lambda name: (-len(name), name))
    return ids


def resolve_module(scope, unit, module_ids, passwd_lookup=pwd.getpwuid):
    """Return the NS8 module id owning this container, or "node" for core.

    An empty unit means attribution failed outright -- the conmon was not
    found, because the container is exiting or was started outside a service
    unit -- and yields "unknown". Reporting it as "node" instead would make it
    indistinguishable from a genuine core container such as redis, and would
    break rate() continuity every time a container flipped between the two.
    """
    if scope["rootless"]:
        try:
            return passwd_lookup(scope["uid"]).pw_name
        except KeyError:
            return "unknown"
    if not unit:
        return "unknown"
    name = unit[: -len(".service")] if unit.endswith(".service") else unit
    for module_id in module_ids:
        if name == module_id or name.startswith(module_id + "-"):
            return module_id
    return "node"


def _read_int(path):
    """Read a single-value cgroup file. "max" and missing files give None."""
    text = _read_text(path)
    if text is None:
        return None
    text = text.strip()
    if text == "max":
        return None
    try:
        return int(text)
    except ValueError:
        return None


def _read_keyed(path):
    """Read a "key value" cgroup file into a dict of ints."""
    text = _read_text(path)
    if text is None:
        return {}
    values = {}
    for line in text.splitlines():
        fields = line.split()
        if len(fields) != 2:
            continue
        try:
            values[fields[0]] = int(fields[1])
        except ValueError:
            continue
    return values


def read_stats(scope_path):
    """Read the CPU, memory, PID and OOM counters of one container scope."""
    cpu = _read_keyed(os.path.join(scope_path, "cpu.stat"))
    memory = _read_keyed(os.path.join(scope_path, "memory.stat"))
    events = _read_keyed(os.path.join(scope_path, "memory.events"))
    try:
        start_time = int(os.stat(scope_path).st_mtime)
    except OSError:
        start_time = None
    return {
        "cpu_user_usec": cpu.get("user_usec"),
        "cpu_system_usec": cpu.get("system_usec"),
        "memory_current": _read_int(os.path.join(scope_path, "memory.current")),
        "memory_peak": _read_int(os.path.join(scope_path, "memory.peak")),
        "memory_max": _read_int(os.path.join(scope_path, "memory.max")),
        "memory_swap": _read_int(os.path.join(scope_path, "memory.swap.current")),
        "memory_anon": memory.get("anon"),
        "memory_file": memory.get("file"),
        "pids_current": _read_int(os.path.join(scope_path, "pids.current")),
        "pids_max": _read_int(os.path.join(scope_path, "pids.max")),
        "oom_kills": events.get("oom_kill"),
        "start_time": start_time,
    }


_IO_FIELDS = ("rbytes", "wbytes", "rios", "wios")


def resolve_block_device(devno, sys_dev_block=DEFAULT_SYS_DEV_BLOCK):
    """Turn a "major:minor" pair into a kernel device name, e.g. "vda"."""
    try:
        return os.path.basename(os.readlink(os.path.join(sys_dev_block, devno)))
    except OSError:
        return devno


def read_io(scope_path, sys_dev_block=DEFAULT_SYS_DEV_BLOCK):
    """Per-device block I/O counters, or None when io.stat is absent."""
    text = _read_text(os.path.join(scope_path, "io.stat"))
    if text is None:
        return None
    devices = []
    for line in text.splitlines():
        fields = line.split()
        if not fields:
            continue
        counters = {}
        for field in fields[1:]:
            key, _, value = field.partition("=")
            if key not in _IO_FIELDS:
                continue
            try:
                counters[key] = int(value)
            except ValueError:
                continue
        if len(counters) != len(_IO_FIELDS):
            continue
        counters["device"] = resolve_block_device(fields[0], sys_dev_block)
        devices.append(counters)
    return devices


def container_pids(scope_path):
    """PIDs of the container payload, innermost cgroup first."""
    pids = []
    for name in ("container/cgroup.procs", "cgroup.procs"):
        text = _read_text(os.path.join(scope_path, name))
        if text:
            pids.extend(text.split())
    return pids


def _read_link(path):
    try:
        return os.readlink(path)
    except OSError:
        return None


def parse_net_dev(text):
    """Parse /proc/<pid>/net/dev into per-interface counters."""
    interfaces = []
    for line in text.splitlines()[2:]:
        name, separator, rest = line.partition(":")
        if not separator:
            continue
        fields = rest.split()
        if len(fields) < 16:
            continue
        try:
            interfaces.append(
                {
                    "device": name.strip(),
                    "receive_bytes": int(fields[0]),
                    "receive_packets": int(fields[1]),
                    "transmit_bytes": int(fields[8]),
                    "transmit_packets": int(fields[9]),
                }
            )
        except ValueError:
            continue
    return interfaces


def read_network(scope_path, proc_root=DEFAULT_PROC_ROOT):
    """Interface counters, or None for host-network containers."""
    host_ns = _read_link(os.path.join(proc_root, "1", "ns", "net"))
    if host_ns is None:
        # Without a trustworthy comparison basis, a private-looking
        # namespace can't be told apart from the host's: treat the sample
        # as absent rather than risk reporting the host's own counters.
        return None
    for pid in container_pids(scope_path):
        container_ns = _read_link(os.path.join(proc_root, pid, "ns", "net"))
        if container_ns is None:
            continue
        if container_ns == host_ns:
            return None
        text = _read_text(os.path.join(proc_root, pid, "net", "dev"))
        if text is None:
            continue
        return parse_net_dev(text)
    return None


def collect(
    cgroup_root=DEFAULT_CGROUP_ROOT,
    proc_root=DEFAULT_PROC_ROOT,
    sys_dev_block=DEFAULT_SYS_DEV_BLOCK,
    rootfull_storage_root=ROOTFULL_STORAGE_ROOT,
    nethserver_root=NETHSERVER_ROOT,
    passwd_lookup=pwd.getpwuid,
):
    """Build one record per live container on this node."""
    scopes = discover_scopes(cgroup_root)
    units = map_units(cgroup_root, proc_root)
    module_ids = list_module_ids(nethserver_root)

    uids = [scope["uid"] for scope in scopes if scope["uid"] is not None]
    names = {}
    for root in storage_roots(uids, rootfull_storage_root, passwd_lookup):
        names.update(read_containers_json(root))

    records = []
    for scope in scopes:
        # One unreadable container must never cost us the whole node's
        # metrics: log it, skip it, keep collecting the rest.
        try:
            unit = units.get(scope["cid"], "")
            meta = names.get(scope["cid"], {})
            records.append(
                {
                    "cid": scope["cid"],
                    "name": meta.get("name") or scope["cid"][:12],
                    "image": meta.get("image", ""),
                    "unit": unit,
                    "rootless": scope["rootless"],
                    "module": resolve_module(scope, unit, module_ids, passwd_lookup),
                    "stats": read_stats(scope["path"]),
                    "io": read_io(scope["path"], sys_dev_block),
                    "network": read_network(scope["path"], proc_root),
                }
            )
        except Exception as ex:
            print(
                "collect(): skipping container %s: %r" % (scope["cid"][:12], ex),
                file=sys.stderr,
            )
    records.sort(key=lambda record: (record["module"], record["name"]))
    return records


def escape_label(value):
    """Escape a label value for the Prometheus text exposition format."""
    return value.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")


class _Exposition(object):
    """Accumulate samples, keeping every family's samples contiguous."""

    def __init__(self):
        self._families = {}
        self._order = []

    def add(self, name, kind, help_text, labels, value):
        if name not in self._families:
            self._families[name] = {"kind": kind, "help": help_text, "samples": []}
            self._order.append(name)
        if labels:
            rendered = ",".join(
                '%s="%s"' % (key, escape_label(val)) for key, val in labels
            )
            sample = "%s{%s} %s" % (name, rendered, value)
        else:
            sample = "%s %s" % (name, value)
        self._families[name]["samples"].append(sample)

    def render(self):
        lines = []
        for name in self._order:
            family = self._families[name]
            lines.append("# HELP %s %s" % (name, family["help"]))
            lines.append("# TYPE %s %s" % (name, family["kind"]))
            lines.extend(family["samples"])
        return "\n".join(lines) + "\n"


_GAUGES = (
    ("memory_current", "ns8_container_memory_usage_bytes", "Container memory usage"),
    ("memory_peak", "ns8_container_memory_peak_bytes", "Container peak memory usage"),
    ("memory_max", "ns8_container_memory_limit_bytes", "Container memory limit"),
    ("memory_swap", "ns8_container_memory_swap_bytes", "Container swap usage"),
    ("memory_anon", "ns8_container_memory_anon_bytes", "Container anonymous memory"),
    ("memory_file", "ns8_container_memory_file_bytes", "Container page cache memory"),
    ("pids_current", "ns8_container_pids", "Processes running in the container"),
    ("pids_max", "ns8_container_pids_limit", "Container process limit"),
    ("start_time", "ns8_container_start_time_seconds", "Container start time"),
)


def render(records, duration_seconds, timestamp):
    """Render collected records as a Prometheus text exposition document."""
    exposition = _Exposition()

    for record in records:
        base = (("module", record["module"]), ("container", record["name"]))
        stats = record["stats"]
        for mode, key in (("user", "cpu_user_usec"), ("system", "cpu_system_usec")):
            if stats[key] is None:
                continue
            exposition.add(
                "ns8_container_cpu_seconds_total",
                "counter",
                "Container CPU time spent, in seconds",
                base + (("mode", mode),),
                "%.6f" % (stats[key] / 1000000.0),
            )

    for record in records:
        base = (("module", record["module"]), ("container", record["name"]))
        for device in record["io"] or []:
            labels = base + (("device", device["device"]),)
            for op, key in (("read", "rbytes"), ("write", "wbytes")):
                exposition.add(
                    "ns8_container_blkio_bytes_total",
                    "counter",
                    "Container block I/O transferred, in bytes",
                    labels + (("op", op),),
                    "%d" % device[key],
                )
        for device in record["io"] or []:
            labels = base + (("device", device["device"]),)
            for op, key in (("read", "rios"), ("write", "wios")):
                exposition.add(
                    "ns8_container_blkio_ops_total",
                    "counter",
                    "Container block I/O operations",
                    labels + (("op", op),),
                    "%d" % device[key],
                )

    for direction, help_text in (
        ("receive", "Container bytes received"),
        ("transmit", "Container bytes transmitted"),
    ):
        for record in records:
            base = (("module", record["module"]), ("container", record["name"]))
            for interface in record["network"] or []:
                exposition.add(
                    "ns8_container_network_%s_bytes_total" % direction,
                    "counter",
                    help_text,
                    base + (("device", interface["device"]),),
                    "%d" % interface["%s_bytes" % direction],
                )

    for direction, help_text in (
        ("receive", "Container packets received"),
        ("transmit", "Container packets transmitted"),
    ):
        for record in records:
            base = (("module", record["module"]), ("container", record["name"]))
            for interface in record["network"] or []:
                exposition.add(
                    "ns8_container_network_%s_packets_total" % direction,
                    "counter",
                    help_text,
                    base + (("device", interface["device"]),),
                    "%d" % interface["%s_packets" % direction],
                )

    for record in records:
        if record["stats"]["oom_kills"] is None:
            continue
        exposition.add(
            "ns8_container_oom_kills_total",
            "counter",
            "Processes killed by the container OOM killer",
            (("module", record["module"]), ("container", record["name"])),
            "%d" % record["stats"]["oom_kills"],
        )

    for key, name, help_text in _GAUGES:
        for record in records:
            value = record["stats"][key]
            if value is None:
                continue
            exposition.add(
                name,
                "gauge",
                help_text,
                (("module", record["module"]), ("container", record["name"])),
                "%d" % value,
            )

    for record in records:
        exposition.add(
            "ns8_container_info",
            "gauge",
            "Container metadata, always 1",
            (
                ("module", record["module"]),
                ("container", record["name"]),
                ("id", record["cid"][:12]),
                ("image", record["image"]),
                ("unit", record["unit"]),
                ("rootless", "true" if record["rootless"] else "false"),
            ),
            "1",
        )

    exposition.add(
        "ns8_container_collector_duration_seconds",
        "gauge",
        "Duration of the last container metrics collection",
        (),
        "%.6f" % duration_seconds,
    )
    exposition.add(
        "ns8_container_collector_last_success_timestamp_seconds",
        "gauge",
        "Unix timestamp of the last successful collection",
        (),
        "%d" % timestamp,
    )
    return exposition.render()


def write_atomic(path, text):
    """Write text to path through a temporary file, then rename it in place."""
    temporary = "%s.%d" % (path, os.getpid())
    with open(temporary, "w") as fp:
        fp.write(text)
    os.rename(temporary, path)

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
import os
import os.path
import re

DEFAULT_CGROUP_ROOT = "/sys/fs/cgroup"
DEFAULT_PROC_ROOT = "/proc"
DEFAULT_SYS_DEV_BLOCK = "/sys/dev/block"
ROOTFULL_STORAGE_ROOT = "/var/lib/containers/storage"

_SCOPE_PATTERNS = (
    "machine.slice/libpod-*.scope",
    "system.slice/libpod-*.scope",
    "user.slice/user-*.slice/user@*.service/user.slice/libpod-*.scope",
)

_SCOPE_RE = re.compile(r"/libpod-([0-9a-f]{64})\.scope$")
_UID_RE = re.compile(r"/user-(\d+)\.slice/")


def discover_scopes(cgroup_root=DEFAULT_CGROUP_ROOT):
    """Return the live container scopes found under cgroup_root."""
    scopes = []
    for pattern in _SCOPE_PATTERNS:
        for path in glob.glob(os.path.join(cgroup_root, pattern)):
            match = _SCOPE_RE.search(path)
            if match is None:
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

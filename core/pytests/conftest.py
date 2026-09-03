#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import json
import os
import sys

import pytest

sys.path.insert(
    0,
    os.path.join(
        os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
        "imageroot/usr/local/agent/pypkg",
    ),
)

CID_A = "a" * 64
CID_B = "b" * 64
CID_C = "c" * 64


class FakeRoot:
    """Builds a synthetic /sys/fs/cgroup, /proc and container storage tree."""

    def __init__(self, base):
        self.base = str(base)
        self.cgroup = os.path.join(self.base, "cgroup")
        self.proc = os.path.join(self.base, "proc")
        self.storage = os.path.join(self.base, "storage")
        self.dev_block = os.path.join(self.base, "dev-block")
        os.makedirs(self.cgroup)
        os.makedirs(self.proc)
        os.makedirs(self.dev_block)
        self.add_process("1", "/sbin/init", netns="net:[4026531840]")

    def _write(self, path, content):
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as fp:
            fp.write(content)

    def add_rootfull_scope(self, cid, files=None):
        path = os.path.join(self.cgroup, "machine.slice", "libpod-%s.scope" % cid)
        return self._add_scope(path, files)

    def add_rootless_scope(self, cid, uid, files=None):
        path = os.path.join(
            self.cgroup,
            "user.slice",
            "user-%d.slice" % uid,
            "user@%d.service" % uid,
            "user.slice",
            "libpod-%s.scope" % cid,
        )
        return self._add_scope(path, files)

    def _add_scope(self, path, files):
        os.makedirs(os.path.join(path, "container"), exist_ok=True)
        for name, content in (files or {}).items():
            self._write(os.path.join(path, name), content)
        return path

    def add_unit(self, unit, pids, uid=None):
        """Create a service cgroup holding the given pids."""
        if uid is None:
            path = os.path.join(self.cgroup, "system.slice", unit)
        else:
            path = os.path.join(
                self.cgroup,
                "user.slice",
                "user-%d.slice" % uid,
                "user@%d.service" % uid,
                "app.slice",
                unit,
            )
        self._write(os.path.join(path, "cgroup.procs"), "".join(p + "\n" for p in pids))
        return path

    def add_process(self, pid, *argv, **kwargs):
        """Create /proc/<pid> with a NUL separated cmdline and a netns link."""
        path = os.path.join(self.proc, pid)
        os.makedirs(os.path.join(path, "ns"), exist_ok=True)
        with open(os.path.join(path, "cmdline"), "wb") as fp:
            fp.write(("\0".join(argv) + "\0").encode())
        netns = kwargs.get("netns")
        if netns is not None:
            # A regular file standing in for the /proc/<pid>/ns/net symlink is
            # not enough: the code calls os.readlink(), so create a real link.
            link = os.path.join(path, "ns", "net")
            if not os.path.lexists(link):
                os.symlink(netns, link)
        net_dev = kwargs.get("net_dev")
        if net_dev is not None:
            self._write(os.path.join(path, "net", "dev"), net_dev)
        return path

    def add_conmon(self, pid, cid):
        return self.add_process(
            pid, "/usr/bin/conmon", "--api-version", "1", "-c", cid, "-u", cid
        )

    def write_containers_json(self, entries, uid=None, home=None):
        if uid is None:
            root = os.path.join(self.storage, "rootfull")
        else:
            root = os.path.join(home, ".local/share/containers/storage")
        self._write(
            os.path.join(root, "overlay-containers", "containers.json"),
            json.dumps(entries),
        )
        return root


@pytest.fixture
def fake_root(tmp_path):
    return FakeRoot(tmp_path)

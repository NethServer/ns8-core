#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

from conftest import CID_A, CID_B

import node.containers as containers


def test_discover_scopes_finds_rootfull_and_rootless(fake_root):
    rootfull = fake_root.add_rootfull_scope(CID_A)
    rootless = fake_root.add_rootless_scope(CID_B, 1004)

    found = containers.discover_scopes(fake_root.cgroup)

    by_cid = {scope["cid"]: scope for scope in found}
    assert set(by_cid) == {CID_A, CID_B}
    assert by_cid[CID_A] == {
        "cid": CID_A,
        "path": rootfull,
        "rootless": False,
        "uid": None,
    }
    assert by_cid[CID_B] == {
        "cid": CID_B,
        "path": rootless,
        "rootless": True,
        "uid": 1004,
    }


def test_discover_scopes_ignores_pause_and_non_scope_dirs(fake_root):
    fake_root.add_rootfull_scope(CID_A)
    import os

    os.makedirs(
        os.path.join(fake_root.cgroup, "user.slice/user-1004.slice/user@1004.service/user.slice/podman-pause-55177911.scope")
    )
    os.makedirs(os.path.join(fake_root.cgroup, "machine.slice/libpod-nothex.scope"))

    found = containers.discover_scopes(fake_root.cgroup)

    assert [scope["cid"] for scope in found] == [CID_A]

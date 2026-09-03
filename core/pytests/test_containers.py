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


def test_read_containers_json_maps_name_and_image(fake_root):
    root = fake_root.write_containers_json(
        [
            {
                "id": CID_A,
                "names": ["crowdsec1"],
                "image": "b62b3268e6a5",
                "metadata": '{"image-name":"docker.io/crowdsecurity/crowdsec:v1.7.8-debian"}',
            }
        ]
    )

    assert containers.read_containers_json(root) == {
        CID_A: {
            "name": "crowdsec1",
            "image": "docker.io/crowdsecurity/crowdsec:v1.7.8-debian",
        }
    }


def test_read_containers_json_tolerates_missing_and_broken_data(fake_root):
    assert containers.read_containers_json("/nonexistent") == {}

    root = fake_root.write_containers_json(
        [
            {"id": CID_A, "names": [], "metadata": "not json"},
            {"names": ["orphan"]},
        ]
    )
    assert containers.read_containers_json(root) == {CID_A: {"name": "", "image": ""}}


def test_storage_roots_lists_rootfull_first_then_module_homes(fake_root):
    class Passwd(object):
        def __init__(self, pw_dir, pw_name):
            self.pw_dir = pw_dir
            self.pw_name = pw_name

    homes = {1004: "/home/metrics1", 1001: "/home/traefik1"}

    def lookup(uid):
        if uid not in homes:
            raise KeyError(uid)
        return Passwd(homes[uid], "user%d" % uid)

    roots = containers.storage_roots(
        [1004, 1001, 1004, 9999],
        rootfull_storage_root="/var/lib/containers/storage",
        passwd_lookup=lookup,
    )

    assert roots == [
        "/var/lib/containers/storage",
        "/home/traefik1/.local/share/containers/storage",
        "/home/metrics1/.local/share/containers/storage",
    ]

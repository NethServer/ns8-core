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


def test_read_containers_json_tolerates_non_object_metadata(fake_root):
    root = fake_root.write_containers_json(
        [{"id": CID_A, "names": ["crowdsec1"], "metadata": "[1,2,3]"}]
    )

    assert containers.read_containers_json(root) == {
        CID_A: {"name": "crowdsec1", "image": ""}
    }


def test_read_containers_json_skips_non_dict_entries(fake_root):
    root = fake_root.write_containers_json(
        ["not-a-dict", {"id": CID_A, "names": ["crowdsec1"]}]
    )

    assert containers.read_containers_json(root) == {
        CID_A: {"name": "crowdsec1", "image": ""}
    }


def test_read_containers_json_tolerates_non_list_names(fake_root):
    root = fake_root.write_containers_json(
        [{"id": CID_A, "names": "crowdsec1"}]
    )

    assert containers.read_containers_json(root) == {CID_A: {"name": "", "image": ""}}


def test_read_containers_json_tolerates_non_list_top_level(fake_root):
    root = fake_root.write_containers_json({"id": CID_A})

    assert containers.read_containers_json(root) == {}


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


def test_map_units_pairs_container_with_its_conmon_unit(fake_root):
    fake_root.add_conmon("4242", CID_A)
    fake_root.add_unit("crowdsec1.service", ["4242"])
    fake_root.add_conmon("5151", CID_B)
    fake_root.add_unit("prometheus.service", ["5151"], uid=1004)

    assert containers.map_units(fake_root.cgroup, fake_root.proc) == {
        CID_A: "crowdsec1.service",
        CID_B: "prometheus.service",
    }


def test_map_units_ignores_non_conmon_processes(fake_root):
    fake_root.add_process("7", "/usr/local/bin/agent", "--agentid=module/crowdsec1")
    fake_root.add_unit("agent@crowdsec1.service", ["7"])

    assert containers.map_units(fake_root.cgroup, fake_root.proc) == {}


def test_list_module_ids_skips_reserved_directories(tmp_path):
    for name in ("cluster", "node", "api-server", "crowdsec1", "mail2"):
        (tmp_path / name).mkdir()
    (tmp_path / "notadir").write_text("")

    assert containers.list_module_ids(str(tmp_path)) == ["crowdsec1", "mail2"]


def test_resolve_module_rootless_uses_owning_user():
    class Passwd(object):
        pw_name = "metrics1"
        pw_dir = "/home/metrics1"

    scope = {"cid": CID_B, "path": "/x", "rootless": True, "uid": 1004}

    assert (
        containers.resolve_module(scope, "prometheus.service", [], lambda uid: Passwd())
        == "metrics1"
    )


def test_resolve_module_rootfull_matches_unit_prefix():
    scope = {"cid": CID_A, "path": "/x", "rootless": False, "uid": None}
    module_ids = ["crowdsec1"]

    assert containers.resolve_module(scope, "crowdsec1.service", module_ids) == "crowdsec1"
    assert (
        containers.resolve_module(scope, "crowdsec1-firewall-bouncer.service", module_ids)
        == "crowdsec1"
    )


def test_resolve_module_core_containers_belong_to_node():
    scope = {"cid": CID_A, "path": "/x", "rootless": False, "uid": None}

    assert containers.resolve_module(scope, "redis.service", ["crowdsec1"]) == "node"
    assert containers.resolve_module(scope, "", ["crowdsec1"]) == "node"

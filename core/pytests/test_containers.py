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


def test_read_stats_parses_every_cgroup_file(fake_root):
    scope = fake_root.add_rootfull_scope(
        CID_A,
        files={
            "cpu.stat": "usage_usec 387935381\nuser_usec 329902588\nsystem_usec 58032793\n",
            "memory.current": "165306368\n",
            "memory.peak": "268435456\n",
            "memory.max": "536870912\n",
            "memory.swap.current": "4096\n",
            "memory.stat": "anon 246710272\nfile 32485376\nkernel 3911680\n",
            "pids.current": "12\n",
            "pids.max": "2048\n",
            "memory.events": "low 0\nhigh 0\nmax 0\noom 3\noom_kill 2\n",
        },
    )

    stats = containers.read_stats(scope)

    assert stats["cpu_user_usec"] == 329902588
    assert stats["cpu_system_usec"] == 58032793
    assert stats["memory_current"] == 165306368
    assert stats["memory_peak"] == 268435456
    assert stats["memory_max"] == 536870912
    assert stats["memory_swap"] == 4096
    assert stats["memory_anon"] == 246710272
    assert stats["memory_file"] == 32485376
    assert stats["pids_current"] == 12
    assert stats["pids_max"] == 2048
    assert stats["oom_kills"] == 2
    assert stats["start_time"] > 0


def test_read_stats_maps_unlimited_to_none(fake_root):
    scope = fake_root.add_rootfull_scope(
        CID_A, files={"memory.max": "max\n", "pids.max": "max\n"}
    )

    stats = containers.read_stats(scope)

    assert stats["memory_max"] is None
    assert stats["pids_max"] is None


def test_read_stats_returns_none_for_missing_files(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)

    stats = containers.read_stats(scope)

    assert stats["cpu_user_usec"] is None
    assert stats["memory_current"] is None
    assert stats["oom_kills"] is None


def test_read_io_parses_per_device_counters(fake_root):
    import os

    os.symlink(
        "../../devices/pci0000:00/0000:00:06.0/virtio4/block/vda",
        os.path.join(fake_root.dev_block, "252:0"),
    )
    scope = fake_root.add_rootfull_scope(
        CID_A,
        files={
            "io.stat": "252:0 rbytes=21557248 wbytes=1190886400 rios=315 wios=91632 dbytes=0 dios=0\n"
        },
    )

    assert containers.read_io(scope, fake_root.dev_block) == [
        {
            "device": "vda",
            "rbytes": 21557248,
            "wbytes": 1190886400,
            "rios": 315,
            "wios": 91632,
        }
    ]


def test_read_io_returns_none_when_controller_not_delegated(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)

    assert containers.read_io(scope, fake_root.dev_block) is None


def test_read_io_falls_back_to_device_number(fake_root):
    scope = fake_root.add_rootfull_scope(
        CID_A, files={"io.stat": "8:16 rbytes=1 wbytes=2 rios=3 wios=4\n"}
    )

    assert containers.read_io(scope, fake_root.dev_block)[0]["device"] == "8:16"


NET_DEV = """Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo:     100       2    0    0    0     0          0         0      100       2    0    0    0     0       0          0
  eth0: 3479743296 1654108    0    0    0     0          0         0 1023429731 1524055    0    0    0     0       0          0
"""


def test_read_network_skips_host_namespace_containers(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)
    fake_root.add_process("900", "/usr/bin/crowdsec", netns="net:[4026531840]", net_dev=NET_DEV)
    with open(scope + "/container/cgroup.procs", "w") as fp:
        fp.write("900\n")

    assert containers.read_network(scope, fake_root.proc) is None


def test_read_network_reports_private_namespace_interfaces(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)
    fake_root.add_process("901", "/usr/bin/app", netns="net:[4026532999]", net_dev=NET_DEV)
    with open(scope + "/container/cgroup.procs", "w") as fp:
        fp.write("901\n")

    assert containers.read_network(scope, fake_root.proc) == [
        {
            "device": "lo",
            "receive_bytes": 100,
            "receive_packets": 2,
            "transmit_bytes": 100,
            "transmit_packets": 2,
        },
        {
            "device": "eth0",
            "receive_bytes": 3479743296,
            "receive_packets": 1654108,
            "transmit_bytes": 1023429731,
            "transmit_packets": 1524055,
        },
    ]


def test_container_pids_prefers_the_container_child_cgroup(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)
    with open(scope + "/cgroup.procs", "w") as fp:
        fp.write("70936\n")
    with open(scope + "/container/cgroup.procs", "w") as fp:
        fp.write("70938\n70939\n")

    assert containers.container_pids(scope) == ["70938", "70939", "70936"]


def test_read_network_returns_none_when_host_netns_is_unreadable(fake_root):
    import os

    scope = fake_root.add_rootfull_scope(CID_A)
    os.remove(os.path.join(fake_root.proc, "1", "ns", "net"))
    fake_root.add_process("902", "/usr/bin/app", netns="net:[4026532999]", net_dev=NET_DEV)
    with open(scope + "/container/cgroup.procs", "w") as fp:
        fp.write("902\n")

    assert containers.read_network(scope, fake_root.proc) is None


def test_read_network_skips_pid_that_no_longer_exists(fake_root):
    scope = fake_root.add_rootfull_scope(CID_A)
    fake_root.add_process("904", "/usr/bin/app", netns="net:[4026533000]", net_dev=NET_DEV)
    with open(scope + "/container/cgroup.procs", "w") as fp:
        fp.write("999999\n904\n")

    assert containers.read_network(scope, fake_root.proc) == [
        {
            "device": "lo",
            "receive_bytes": 100,
            "receive_packets": 2,
            "transmit_bytes": 100,
            "transmit_packets": 2,
        },
        {
            "device": "eth0",
            "receive_bytes": 3479743296,
            "receive_packets": 1654108,
            "transmit_bytes": 1023429731,
            "transmit_packets": 1524055,
        },
    ]


def test_parse_net_dev_skips_malformed_lines():
    text = NET_DEV + "  eth1: 1 2 3\n  eth2: bad 1654108 0 0 0 0 0 0 1023429731 1524055 0 0 0 0 0 0\n"

    assert containers.parse_net_dev(text) == [
        {
            "device": "lo",
            "receive_bytes": 100,
            "receive_packets": 2,
            "transmit_bytes": 100,
            "transmit_packets": 2,
        },
        {
            "device": "eth0",
            "receive_bytes": 3479743296,
            "receive_packets": 1654108,
            "transmit_bytes": 1023429731,
            "transmit_packets": 1524055,
        },
    ]


def test_collect_builds_one_record_per_container(fake_root, tmp_path):
    nethserver = tmp_path / "nethserver"
    (nethserver / "crowdsec1").mkdir(parents=True)
    (nethserver / "node").mkdir()

    home = str(tmp_path / "home" / "metrics1")

    rootfull = fake_root.add_rootfull_scope(
        CID_A, files={"memory.current": "100\n", "cpu.stat": "user_usec 5\nsystem_usec 6\n"}
    )
    rootless = fake_root.add_rootless_scope(
        CID_B, 1004, files={"memory.current": "200\n"}
    )
    assert rootfull and rootless

    fake_root.add_conmon("4242", CID_A)
    fake_root.add_unit("crowdsec1.service", ["4242"])
    fake_root.add_conmon("5151", CID_B)
    fake_root.add_unit("prometheus.service", ["5151"], uid=1004)

    rootfull_storage = fake_root.write_containers_json(
        [{"id": CID_A, "names": ["crowdsec1"], "metadata": '{"image-name":"crowdsec:1"}'}]
    )
    fake_root.write_containers_json(
        [{"id": CID_B, "names": ["prometheus"], "metadata": '{"image-name":"prom:2"}'}],
        uid=1004,
        home=home,
    )

    class Passwd(object):
        pw_name = "metrics1"
        pw_dir = home

    records = containers.collect(
        cgroup_root=fake_root.cgroup,
        proc_root=fake_root.proc,
        sys_dev_block=fake_root.dev_block,
        rootfull_storage_root=rootfull_storage,
        nethserver_root=str(nethserver),
        passwd_lookup=lambda uid: Passwd(),
    )

    assert [(r["module"], r["name"]) for r in records] == [
        ("crowdsec1", "crowdsec1"),
        ("metrics1", "prometheus"),
    ]
    assert records[0]["image"] == "crowdsec:1"
    assert records[0]["unit"] == "crowdsec1.service"
    assert records[0]["rootless"] is False
    assert records[0]["stats"]["memory_current"] == 100
    assert records[0]["io"] is None
    assert records[0]["network"] is None
    assert records[1]["rootless"] is True
    assert records[1]["stats"]["memory_current"] == 200


def test_collect_falls_back_to_short_id_when_name_is_unknown(fake_root, tmp_path):
    fake_root.add_rootfull_scope(CID_A)

    records = containers.collect(
        cgroup_root=fake_root.cgroup,
        proc_root=fake_root.proc,
        sys_dev_block=fake_root.dev_block,
        rootfull_storage_root=str(tmp_path / "missing"),
        nethserver_root=str(tmp_path / "missing"),
    )

    assert len(records) == 1
    assert records[0]["name"] == CID_A[:12]
    assert records[0]["module"] == "node"

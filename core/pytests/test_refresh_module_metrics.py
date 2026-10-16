#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

"""Unit tests for refresh-module-metrics.

Run with: python3 -m unittest discover -s core/pytests
"""

import importlib.machinery
import importlib.util
import io
import os
import sys
import tempfile
import time
import types
import unittest
import unittest.mock

CORE = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..")
SCRIPT = os.path.join(
    CORE, "imageroot/var/lib/nethserver/node/bin/refresh-module-metrics"
)
PYPKG = os.path.join(CORE, "imageroot/usr/local/agent/pypkg")


def read(path):
    with open(path) as fp:
        return fp.read()


def import_agent():
    # Stub redis and aiohttp, imported by agent but unused here
    for name in ("redis", "aiohttp"):
        if name not in sys.modules:
            try:
                __import__(name)
            except ImportError:
                sys.modules[name] = types.ModuleType(name)
    if not hasattr(sys.modules["redis"], "asyncio"):
        sys.modules["redis"].asyncio = types.ModuleType("redis.asyncio")
    if PYPKG not in sys.path:
        sys.path.insert(0, PYPKG)
    import agent

    return agent


def load_script():
    loader = importlib.machinery.SourceFileLoader("refresh_module_metrics", SCRIPT)
    spec = importlib.util.spec_from_loader(loader.name, loader)
    module = importlib.util.module_from_spec(spec)
    loader.exec_module(module)
    return module


import_agent()
rmm = load_script()

GOOD = """\
# HELP app_jobs_total Jobs run
# TYPE app_jobs_total counter
app_jobs_total{job="backup"} 3
app_jobs_total{job="sync"} 5
"""


class Fixture(unittest.TestCase):
    """Temp folders and a captured stderr."""

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        stderr = unittest.mock.patch.object(sys, "stderr", io.StringIO())
        self.stderr = stderr.start()
        self.addCleanup(stderr.stop)
        self.output_dir = os.path.join(self.tmp.name, "node_exporter")
        os.makedirs(self.output_dir)

    def log(self):
        return self.stderr.getvalue()

    def install(self, module_id, mode=0o755):
        """Create a module install dir with an export-metrics command."""
        install_dir = os.path.join(self.tmp.name, module_id)
        os.makedirs(os.path.join(install_dir, "bin"))
        path = os.path.join(install_dir, "bin", rmm.COMMAND)
        with open(path, "w") as fp:
            fp.write("#!/bin/sh\n")
        os.chmod(path, mode)
        return (module_id, install_dir)

    def write(self, name, text):
        with open(os.path.join(self.output_dir, name), "w") as fp:
            fp.write(text)


class TestRelabel(unittest.TestCase):
    def test_module_label_added(self):
        self.assertEqual(
            rmm.relabel(GOOD, "mail1").splitlines(),
            [
                "# HELP app_jobs_total Jobs run",
                "# TYPE app_jobs_total counter",
                'app_jobs_total{module="mail1",job="backup"} 3',
                'app_jobs_total{module="mail1",job="sync"} 5',
            ],
        )

    def test_samples_without_labels(self):
        self.assertEqual(
            rmm.relabel("app_a 1\napp_b{} 2\napp_c\t3\n", "mail1").splitlines(),
            ['app_a{module="mail1"} 1', 'app_b{module="mail1"} 2', 'app_c{module="mail1"}\t3'],
        )

    def test_label_value_with_braces(self):
        self.assertEqual(
            rmm.relabel('app_x{a="}{"} 1\n', "mail1"), 'app_x{module="mail1",a="}{"} 1\n'
        )

    def test_identity_labels_rejected(self):
        for labels in ('module="samba1"', 'a="1",module_id="samba1"', 'a="1", module ="x"'):
            with self.assertRaises(rmm.MetricsError) as ctx:
                rmm.relabel("app_x 1\napp_x{%s} 1\n" % labels, "mail1")
            self.assertIn("line 2: module and module_id labels are reserved", str(ctx.exception))
        # Not a label name
        rmm.relabel('app_x{mymodule="1",module_name="2"} 1\n', "mail1")


class TestListModules(Fixture):
    def test_rootfull_and_rootless(self):
        root = os.path.join(self.tmp.name, "nethserver")
        for name in ("crowdsec1", "node", "cluster", "api-server"):
            os.makedirs(os.path.join(root, name))
        homes = os.path.join(self.tmp.name, "home")
        users = []
        for name, marker in (("mail1", True), ("admin", True), ("samba1", False)):
            home = os.path.join(homes, name)
            os.makedirs(os.path.join(home, ".config/state"))
            if marker:
                open(os.path.join(home, rmm.MODULE_MARKER), "w").close()
            users.append(types.SimpleNamespace(pw_name=name, pw_dir=home))
        self.assertEqual(
            rmm.list_modules(root, lambda: users),
            [
                ("crowdsec1", os.path.join(root, "crowdsec1")),
                ("mail1", os.path.join(homes, "mail1", ".config")),
            ],
        )


class TestRunExport(Fixture):
    def run_script(self, script, module_id="mail1"):
        # sh -c SCRIPT sh MODULE_ID export-metrics
        return rmm.run_export(module_id, command=("sh", "-c", script, "sh"))

    def test_output_relabeled(self):
        text = self.run_script('test "$1 $2" = "mail1 export-metrics" && echo "app_x 1"')
        self.assertEqual(text, 'app_x{module="mail1"} 1\n')

    def test_stderr_logged_with_module(self):
        self.run_script("echo app_x 1; echo some detail >&2")
        self.assertEqual(self.log(), "module mail1: some detail\n")

    def test_failure(self):
        with self.assertRaises(rmm.MetricsError) as ctx:
            self.run_script("echo app_x 1; exit 3")
        self.assertIn("failed with exit code 3", str(ctx.exception))

    def test_timeout_kills_every_process(self):
        marker = os.path.join(self.tmp.name, "late")
        started = time.monotonic()
        with unittest.mock.patch.object(rmm, "TIMEOUT", 1):
            with self.assertRaises(rmm.MetricsError) as ctx:
                self.run_script("(sleep 2; touch %s) & sleep 10" % marker)
        self.assertLess(time.monotonic() - started, 5)
        self.assertIn("timed out after 1 seconds", str(ctx.exception))
        time.sleep(2.5)
        self.assertFalse(os.path.exists(marker))

    def test_timeout_with_escaped_process(self):
        pidfile = os.path.join(self.tmp.name, "pid")
        started = time.monotonic()
        try:
            with unittest.mock.patch.object(rmm, "TIMEOUT", 1), \
                    unittest.mock.patch.object(rmm, "KILL_WAIT", 1):
                with self.assertRaises(rmm.MetricsError) as ctx:
                    # setsid leaves the session but keeps stdout open
                    self.run_script("setsid sleep 30 & echo $! > %s; sleep 30" % pidfile)
            self.assertLess(time.monotonic() - started, 5)
            self.assertIn("timed out", str(ctx.exception))
        finally:
            os.kill(int(read(pidfile)), 9)


class TestCollect(Fixture):
    def test_modules_with_and_without_command(self):
        mail = self.install("mail1")
        samba = (("samba1", os.path.join(self.tmp.name, "samba1")))
        called = []

        def export(module_id):
            called.append(module_id)
            return 'app_x{module="%s"} 1\n' % module_id

        outputs, keep = rmm.collect([mail, samba], export)
        self.assertEqual(called, ["mail1"])
        self.assertEqual(outputs, {"module-mail1.prom": 'app_x{module="mail1"} 1\n'})
        self.assertEqual(keep, set())

    def test_not_executable_ignored(self):
        mail = self.install("mail1", mode=0o644)
        outputs, keep = rmm.collect([mail], lambda module_id: self.fail("called"))
        self.assertEqual((outputs, keep), ({}, set()))
        self.assertIn("module mail1: ", self.log())
        self.assertIn("export-metrics is not executable: ignored", self.log())

    def test_failure_keeps_last_output(self):
        def export(module_id):
            if module_id == "mail1":
                raise rmm.MetricsError("export-metrics failed with exit code 1")
            return "app_x 1\n"

        outputs, keep = rmm.collect([self.install("mail1"), self.install("samba1")], export)
        self.assertEqual(outputs, {"module-samba1.prom": "app_x 1\n"})
        self.assertEqual(keep, {"module-mail1.prom"})
        self.assertIn(
            "module mail1: export-metrics failed with exit code 1: last output kept",
            self.log(),
        )


class TestPublish(Fixture):
    def test_write_keep_and_remove(self):
        # Files of other collectors are left alone
        self.write("containers.prom", "x 1\n")
        self.write("volumes.prom.tmp", "x 1\n")
        self.write("module-mail1.prom", "old mail\n")
        self.write("module-samba1.prom", "old samba\n")
        self.write("module-gone1.prom", "old gone\n")
        self.write("module-gone1.prom.tmp", "x")
        rmm.publish(
            self.output_dir, {"module-samba1.prom": "new samba\n"}, {"module-mail1.prom"}
        )
        self.assertEqual(
            sorted(os.listdir(self.output_dir)),
            ["containers.prom", "module-mail1.prom", "module-samba1.prom", "volumes.prom.tmp"],
        )
        self.assertEqual(read(os.path.join(self.output_dir, "module-mail1.prom")), "old mail\n")
        self.assertEqual(read(os.path.join(self.output_dir, "module-samba1.prom")), "new samba\n")

    def test_unchanged_file_not_rewritten(self):
        path = os.path.join(self.output_dir, "module-mail1.prom")
        rmm.publish(self.output_dir, {"module-mail1.prom": "a 1\n"})
        os.utime(path, (1, 1))
        rmm.publish(self.output_dir, {"module-mail1.prom": "a 1\n"})
        self.assertEqual(os.stat(path).st_mtime, 1)
        rmm.publish(self.output_dir, {"module-mail1.prom": "a 2\n"})
        self.assertNotEqual(os.stat(path).st_mtime, 1)
        self.assertEqual(read(path), "a 2\n")

    def test_utf8_output_without_utf8_locale(self):
        real_open = open
        real_fdopen = os.fdopen

        # Simulate a node whose locale is not UTF-8
        def ascii_default_open(file, mode="r", *args, **kwargs):
            if "b" not in mode and kwargs.get("encoding") is None:
                kwargs["encoding"] = "ascii"
            return real_open(file, mode, *args, **kwargs)

        def ascii_default_fdopen(fd, mode="r", *args, **kwargs):
            if "b" not in mode and kwargs.get("encoding") is None:
                kwargs["encoding"] = "ascii"
            return real_fdopen(fd, mode, *args, **kwargs)

        with unittest.mock.patch("builtins.open", ascii_default_open), \
                unittest.mock.patch.object(os, "fdopen", ascii_default_fdopen):
            rmm.publish(self.output_dir, {"module-mail1.prom": 'a{c="Forlì ☃"} 1\n'})
        with open(os.path.join(self.output_dir, "module-mail1.prom"), encoding="utf-8") as fp:
            self.assertIn("Forlì ☃", fp.read())


if __name__ == "__main__":
    unittest.main()

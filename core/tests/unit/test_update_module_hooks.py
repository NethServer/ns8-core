# SPDX-License-Identifier: GPL-3.0-or-later
"""Exercise the real update hook runner without an agent or containers."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[3]
RUNNER = ROOT / 'core/imageroot/usr/local/agent/actions/update-module/50run_scriptdir'


class UpdateHooksTest(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.root = Path(temporary.name)
        self.hooks = self.root / 'update-module.d'
        self.hooks.mkdir()

    def hook(self, name, code=0, fatal=False, executable=True):
        path = self.hooks / name
        path.write_text(f'#!/bin/bash\nset -e\necho {name} >> ran\nexit {code}\n')
        path.chmod(0o755 if executable else 0o644)
        if fatal:
            (self.hooks / (name + '.fatal')).touch()

    def run_hooks(self):
        result = subprocess.run(['bash', str(RUNNER)], cwd=self.root,
                                env=dict(os.environ, AGENT_INSTALL_DIR=str(self.root),
                                         HOME=str(self.root)), capture_output=True, text=True, timeout=5)
        ran = (self.root / 'ran').read_text().splitlines() if (self.root / 'ran').exists() else []
        return result, ran

    def test_unmarked_failures_still_warn_and_continue(self):
        self.hook('10first', code=17)
        self.hook('20second')
        result, ran = self.run_hooks()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(ran, ['10first', '20second'])
        self.assertIn('10first has failed (exit 17)', result.stderr)
        self.assertEqual(result.stdout, '')

    def test_marked_failure_returns_exact_status_and_skips_later_hooks(self):
        for status in [1, 17, 124]:
            with self.subTest(status=status):
                (self.root / 'ran').unlink(missing_ok=True)
                self.hook('10optional', code=23)
                self.hook('20critical', code=status, fatal=True)
                self.hook('30later')
                result, ran = self.run_hooks()
                self.assertEqual(result.returncode, status, result.stderr)
                self.assertEqual(ran, ['10optional', '20critical'])
                self.assertIn('Critical update hook failed', result.stderr)

    def test_successful_marked_hook_continues(self):
        self.hook('10critical', fatal=True)
        self.hook('20later')
        result, ran = self.run_hooks()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(ran, ['10critical', '20later'])

    def test_marker_does_not_make_nonexecutable_hook_run(self):
        self.hook('10disabled', code=17, fatal=True, executable=False)
        self.hook('20later')
        result, ran = self.run_hooks()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(ran, ['20later'])

    def test_missing_or_empty_hook_directory_is_compatible(self):
        result, ran = self.run_hooks()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(ran, [])
        self.hooks.rmdir()
        result, ran = self.run_hooks()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(ran, [])


if __name__ == '__main__':
    unittest.main()

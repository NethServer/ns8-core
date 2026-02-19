#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import fcntl
import tempfile
from contextlib import contextmanager

@contextmanager
def safe_open(path, mode="w", encoding=None):
    if "w" not in mode and "x" not in mode:
        raise ValueError("Only write modes supported")

    directory = os.path.dirname(path) or "."

    fd = None
    tmp_path = None
    try:
        # Unique temp file in same directory
        fd, tmp_path = tempfile.mkstemp(
            prefix="." + os.path.basename(path) + "-",
            dir=directory
        )

        if "b" in mode:
            f = os.fdopen(fd, mode)
        else:
            f = os.fdopen(fd, mode, encoding=encoding)

        try:
            yield f

            # Crash-safe commit
            f.flush()
            os.fsync(f.fileno())
            f.close()

            os.replace(tmp_path, path)

            # Ensure rename durability
            dirfd = os.open(directory, os.O_DIRECTORY)
            try:
                os.fsync(dirfd)
            finally:
                os.close(dirfd)

        except Exception:
            f.close()
            raise

    finally:
        if tmp_path and os.path.exists(tmp_path):
            try:
                os.unlink(tmp_path)
            except FileNotFoundError:
                pass

@contextmanager
def exclusive_file_lock(path):
    directory = os.path.dirname(path) or "."
    lock_path = directory + '/.' + os.path.basename(path) + '.lock'
    fd = os.open(lock_path, os.O_CREAT | os.O_RDWR, 0o644)
    try:
        fcntl.flock(fd, fcntl.LOCK_EX)
        yield
    finally:
        fcntl.flock(fd, fcntl.LOCK_UN)
        os.close(fd)

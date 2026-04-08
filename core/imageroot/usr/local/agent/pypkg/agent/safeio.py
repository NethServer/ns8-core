#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import stat
import fcntl
import tempfile
from contextlib import contextmanager

@contextmanager
def safe_open(path, mode="w", encoding=None):
    if "w" not in mode and "x" not in mode:
        raise ValueError("Only write modes supported")

    directory = os.path.dirname(path) or "."
    filename = os.path.basename(path)

    fd = None
    f = None
    tmp_path = None
    try:
        try:
            # Read original file permissions
            ref_mode = stat.S_IMODE(os.stat(path).st_mode)
            # Unique temp file in same directory
            fd, tmp_path = tempfile.mkstemp(
                prefix="." + filename + ".tmp-",
                dir=directory
            )
            # Copy original file permissions to temp file
            os.fchmod(fd, ref_mode)
        except FileNotFoundError:
            attempt = 0
            base_flags = os.O_RDWR if '+' in mode else os.O_WRONLY
            while True:
                # Generate a unique temporary file name
                tmp_path = os.path.join(directory, ".{0}.tmp-{1}-{2}".format(filename, os.getpid(), attempt))
                try:
                    # Try to create a new file on tmp_path. Kernel
                    # applies current process umask on file creation
                    fd = os.open(tmp_path, base_flags | os.O_CREAT | os.O_EXCL, 0o666)
                    break # exit retry-loop
                except FileExistsError:
                    # tmp_path is not unique, try a different path
                    attempt += 1

        if "b" in mode:
            f = os.fdopen(fd, mode)
        else:
            f = os.fdopen(fd, mode, encoding=encoding)
        fd = None  # fd is now owned by f

        try:
            yield f

            # Crash-safe commit
            f.flush()
            os.fsync(f.fileno())
            f.close()
            f = None  # Mark as closed

            os.replace(tmp_path, path)

            # Ensure rename durability
            dirfd = os.open(directory, os.O_DIRECTORY)
            try:
                os.fsync(dirfd)
            finally:
                os.close(dirfd)

        except Exception:
            if f is not None:
                f.close()
            raise

    finally:
        if fd is not None:
            os.close(fd)
        if tmp_path and os.path.exists(tmp_path):
            try:
                os.unlink(tmp_path)
            except FileNotFoundError:
                pass

@contextmanager
def exclusive_file_lock(path):
    directory = os.path.dirname(path) or "."
    lock_path = directory + '/.' + os.path.basename(path).lstrip(".") + '.lock'
    fd = os.open(lock_path, os.O_CREAT | os.O_RDWR, 0o644)
    try:
        fcntl.flock(fd, fcntl.LOCK_EX)
        yield
    finally:
        fcntl.flock(fd, fcntl.LOCK_UN)
        os.close(fd)

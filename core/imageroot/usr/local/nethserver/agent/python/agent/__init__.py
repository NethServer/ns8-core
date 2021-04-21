#
# Copyright (C) 2021 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# NethServer is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with NethServer.  If not, see COPYING.
#

import redis
import os
import subprocess
import sys

def redis_connect():
    """Connect to the Redis DB with the right credentials
    """
    redis_host = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[0]
    redis_port = os.getenv('REDIS_ADDRESS', '127.0.0.1:6379').split(':', 1)[1]
    redis_username = os.getenv('REDIS_USERNAME', 'default')
    redis_password = os.getenv('REDIS_PASSWORD', 'letmein')

    return redis.Redis(
        host=redis_host,
        port=redis_port,
        db=0,
        username=redis_username,
        password=redis_password,
        decode_responses=True,
            #  we assume Redis keys and value strings are encoded UTF-8. Enabling this
            #  option implicitly converts to UTF-8 strings instead of binary strings 
            #  (e.g. {b'key': b'value'} != {'key':'value'})
    )


def run_helper(*args):
    """Run the command and assert the exit code is 0

    If the exit code is non-zero raise an assertion error.
    The command output is redirected to stderr.
    """
    proc = subprocess.run(args, stdout=sys.stderr)
    assert proc.returncode == 0
    return proc

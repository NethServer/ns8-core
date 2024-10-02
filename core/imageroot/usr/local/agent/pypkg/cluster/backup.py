#
# Copyright (C) 2023 Nethesis S.r.l.
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

import sys
import subprocess

def get_default_backup_repository_name(provider, url, rid=""):
    """Suggest a default name for a backup repository"""
    name = ""
    try:
        if provider == "aws":
            name = "AWS " + url.split('/', 1)[1]
        elif provider == "backblaze":
            name = "Backblaze " + url.split(':', 1)[1]
    except:
        pass

    if not name:
        name = provider + " " + rid[:4]

    return name

def validate_schedule(schedule):
    """Check the given schedule string is good for Systemd timers"""
    return subprocess.run(['/usr/bin/systemd-analyze', 'calendar', schedule], stdout=subprocess.DEVNULL).returncode != 0

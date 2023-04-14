#!/bin/sh

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

set -e

echo "${RSYNCD_USER}:${RSYNCD_PASSWORD}" >/etc/rsyncd.secrets
chmod 600 /etc/rsyncd.secrets

cat >/etc/rsyncd.conf <<EOF
#
# rsyncd configuration for module state transferring
#
address = ${RSYNCD_ADDRESS}
port = ${RSYNCD_PORT}
syslog tag = ${RSYNCD_SYSLOG_TAG}
reverse lookup = no
forward lookup = no
uid = 0
gid = 0
auth users = ${RSYNCD_USER}
secrets file = /etc/rsyncd.secrets
# If chroot is enabled, -a/--perms flag does not work
use chroot = no

[data]
path = /srv
read only = no
numeric ids = yes
munge symlinks = no
hosts allow = 127.0.0.1 ::1 ${RSYNCD_NETWORK:-0.0.0.0/0}

[terminate]
path = /var/empty
read only = yes
list = no
post-xfer exec = kill -TERM 1
EOF


# Start the container CMD
exec "${@}"

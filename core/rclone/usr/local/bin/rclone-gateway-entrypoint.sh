#!/bin/ash

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

# shellcheck shell=dash disable=SC3045

set -e

rclone ${DEBUG:+-vvvv} rcd --rc-no-auth --rc-addr "${RCLONE_UNIX_SOCKET}" &
echo "Initialize backends..."
rclone-serve-restart
echo "Initialize frontend..."
haproxy ${DEBUG:+-V} \
    -W \
    -db \
    -p /var/lib/rclone/haproxy.pid \
    -f /etc/haproxy/haproxy.cfg &

wait -n

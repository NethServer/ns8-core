#!/bin/ash

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

# shellcheck shell=dash disable=SC3045

set -e

rclone ${DEBUG:+-vvvv} \
    rcd \
    --rc-no-auth \
    --rc-addr "${RCLONE_UNIX_SOCKET}" \
    &
sleep 0.2
until rclone rc core/version >/dev/null 2>&1; do
    echo "Still waiting for rclone startup..."
    sleep 1
done
echo "Initialize backends..."
rclone-serve-restart
echo "Initialize frontend..."
haproxy ${DEBUG:+-V} \
    -W \
    -db \
    -p /var/lib/rclone/haproxy.pid \
    -f /etc/haproxy/haproxy.cfg \
    2>/dev/null &

wait -n

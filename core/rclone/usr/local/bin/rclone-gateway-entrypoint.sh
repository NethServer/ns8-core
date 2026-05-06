#!/bin/sh

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e

rclone ${DEBUG:+-vvvv} rcd --rc-no-auth --rc-addr "${RCLONE_UNIX_SOCKET}" &
rclone-serve-restart
exec haproxy ${DEBUG:+-V} -W -db -f /etc/haproxy/haproxy.cfg

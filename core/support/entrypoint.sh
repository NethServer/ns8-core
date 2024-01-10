#!/bin/sh

#
# Copyright (C) 2024 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e

set -a
# Mandatory vars
: "${VPN_PEER_HOST:?}"
: "${VPN_USERNAME:?}"
: "${VPN_PASSWORD:?}"
: "${VPN_CERT_FILE:?}"
: "${VPN_CERT_COMMON_NAME:?}"

# Optional UDP port number
: "${VPN_PEER_PORT:=1194}"
: "${VPN_INACTIVE:=0}"
set +a

if [ $# -eq 0 ]; then
    cd /var/lib/openvpn
    ( umask 077 ; printf "%s\n%s\n" "${VPN_USERNAME}" "${VPN_PASSWORD}" > credentials )
    envsubst < /usr/local/templates/openvpn.conf > openvpn.conf
    exec openvpn --config openvpn.conf
else
    exec "${@}"
fi

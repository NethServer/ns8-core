#!/bin/bash

#
# Copyright (C) 2022 Nethesis S.r.l.
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

provision_type=${1:?The mandatory entrypoint argument is missing}

if [[ $provision_type == "run-dc" ]]; then
    extra_args=()
    if [[ -f "/etc/dns_forwarder" ]]; then
        extra_args+=("--option=dns forwarder=$(</etc/dns_forwarder)")
    fi
    exec samba -i --debug-stderr "${extra_args[@]}"
elif [[ $provision_type == "join-domain" ]]; then
    echo "Joining domain ${REALM}..."
    rm -f /etc/samba/smb.conf
    # Read adminuser and pass from standard input
    read -p "adminuser: " -r ADMINUSER
    read -s -p "adminpass: " -r ADMINPASS
    kinit "${ADMINUSER}" <<<"${ADMINPASS}"
    samba-tool domain join "${REALM,,}" DC \
        -k yes \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    kdestroy
elif [[ $provision_type == "new-domain" ]]; then
    echo "Starting domain provisioning procedure..."
    rm -f /etc/samba/smb.conf
    # Read adminuser and pass from standard input
    read -p "adminuser: " -r ADMINUSER
    read -s -p "adminpass: " -r ADMINPASS
    if [[ "${ADMINUSER}" == "administrator" ]]; then
        usebuiltinadmin=1
    fi
    samba-tool domain provision \
        --server-role=dc "${usebuiltinadmin:+--adminpass=${ADMINPASS}}" \
        "--domain=${NBDOMAIN}" "--realm=${REALM}" \
        "--host-ip=${IPADDRESS}" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    if [[ -z "${usebuiltinadmin}" ]]; then
        samba-tool user create "${ADMINUSER}" <<<"${ADMINPASS}"$'\n'"${ADMINPASS}"
        samba-tool user disable administrator
    fi
    samba-tool user setexpiry --noexpiry "${ADMINUSER}"
    samba-tool user create "${SVCUSER}" <<<"${SVCPASS}"$'\n'"${SVCPASS}"
    samba-tool user setexpiry --noexpiry "${SVCUSER}"
else
    echo "[WARNING] Provision type $provision_type is not recognized: now going to sleep +INF"
    exec sleep +INF
fi

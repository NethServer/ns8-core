#!/bin/bash

set -e

if grep -q -F "${REALM}" /etc/samba/smb.conf; then
    : # pass
elif [[ $1 == "join-domain" ]]; then
    echo "[NOTICE] Joining domain ${REALM}..."
    rm -f /etc/samba/smb.conf
    kinit ${ADMINUSER:-administrator} <<<"${ADMINPASS}"
    samba-tool domain join ${REALM,,} DC \
        -k yes \
        "--option=dns forwarder = 127.0.0.53" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    kdestroy
elif [[ $1 == "new-domain" ]]; then
    echo "[NOTICE] Starting domain provisioning procedure..."
    rm -f /etc/samba/smb.conf
    samba-tool domain provision \
        --server-role=dc "--domain=${NBDOMAIN}" "--realm=${REALM}" \
        "--adminpass=${ADMINPASS:-Nethesis,1234}" \
        "--host-ip=${IPADDRESS}" \
        "--option=dns forwarder = 127.0.0.53" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
else
    echo "[NOTICE] arg $1 not recognized: now going to sleep +INF"
    exec sleep +INF
fi

exec samba -i --debug-stderr

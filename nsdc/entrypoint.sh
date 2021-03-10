#!/bin/bash

if ! grep -q -F "${REALM}" /etc/samba/smb.conf; then
    echo "[NOTICE] Starting domain provisioning procedure..."
    rm -f /etc/samba/smb.conf
    samba-tool domain provision \
        --server-role=dc "--domain=${NBDOMAIN}" "--realm=${REALM}" \
        "--adminpass=${ADMINPASS:-Nethesis,1234}" \
        "--host-ip=${IPADDRESS}" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    cp -av /var/lib/samba/private/krb5.conf /etc/krb5.conf
fi

exec samba -i --debug-stderr

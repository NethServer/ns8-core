#!/bin/bash

#
# Override the provisioning procedure:
# 1. add --use-ntvfs to run in a rootless container
# 2. bind the main IP address only (and 127.0.0.1) to avoid resolved conflicts
#

if ! grep -q -F ${NSDC_REALM} /etc/samba/smb.conf; then
    echo "[NOTICE] Starting domain provisioning procedure..."
    rm -f /etc/samba/smb.conf
    samba-tool domain provision \
        --server-role=dc --domain=${NSDC_NBDOMAIN} --realm=${NSDC_REALM} \
        --adminpass=${NSDC_ADMINPASS:-Nethesis,1234} \
        --host-ip=${NSDC_IPADDRESS} \
        --option="bind interfaces only = yes" \
        --option="interfaces = 127.0.0.1 ${NSDC_IPADDRESS}"
    cp -av /var/lib/samba/private/krb5.conf /etc/krb5.conf
fi

exec samba -i --debug-stderr

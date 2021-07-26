#!/bin/bash

set -e

provision_type=${1:?The mandatory entrypoint argument is missing}

if [[ $provision_type == "run-dc" ]]; then
    exec samba -i --debug-stderr
elif [[ $provision_type == "join-domain" ]]; then
    echo "Joining domain ${REALM}..."
    rm -f /etc/samba/smb.conf
    kinit "${ADMINUSER:-administrator}" <<<"${ADMINPASS}"
    samba-tool domain join ${REALM,,} DC \
        -k yes \
        "--option=dns forwarder = 127.0.0.53" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    kdestroy
elif [[ $provision_type == "new-domain" ]]; then
    echo "Starting domain provisioning procedure..."
    rm -f /etc/samba/smb.conf
    samba-tool domain provision \
        --server-role=dc "--domain=${NBDOMAIN}" "--realm=${REALM}" \
        "--adminpass=${ADMINPASS:-Nethesis,1234}" \
        "--host-ip=${IPADDRESS}" \
        "--option=dns forwarder = 127.0.0.53" \
        "--option=bind interfaces only = yes" \
        "--option=interfaces = 127.0.0.1 ${IPADDRESS}"
    samba-tool user setexpiry --noexpiry "${ADMINUSER:-administrator}"
    samba-tool user create "${SVCUSER}" <<<"${SVCPASS}"$'\n'"${SVCPASS}"
    samba-tool user setexpiry --noexpiry "${SVCUSER}"
else
    echo "[WARNING] Provision type $provision_type is not recognized: now going to sleep +INF"
    exec sleep +INF
fi

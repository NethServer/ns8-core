#!/bin/bash

set -e

provision_type=${1:?The mandatory entrypoint argument is missing}

if [[ $provision_type == "run-dc" ]]; then
    exec samba -i --debug-stderr
elif [[ $provision_type == "join-domain" ]]; then
    echo "Joining domain ${REALM}..."
    rm -f /etc/samba/smb.conf
    # Read adminuser and pass from standard input
    read -p "adminuser: " -r ADMINUSER
    read -s -p "adminpass: " -r ADMINPASS
    kinit "${ADMINUSER}" <<<"${ADMINPASS}"
    samba-tool domain join "${REALM,,}" DC \
        -k yes \
        "--option=dns forwarder = 127.0.0.53" \
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
        "--option=dns forwarder = 127.0.0.53" \
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

#!/bin/bash

#
# Copyright (C) 2023 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e

echo "Checking port 80 and 443 are not already in use"
for port in 80 443
do
    if ss -H -l "( sport = :${port} )" | grep -q .; then
        echo "Installation failed: port ${port} is already in use."
        exit 1
    fi
done

echo "Restart journald:"
systemctl restart systemd-journald.service

core_url='ghcr.io/nethserver/core:ns8-stable'
source /etc/os-release

echo "Install dependencies:"
if [[ "${PLATFORM_ID}" == "platform:el9" ]]; then
    dnf update -y # Fix SELinux issues with basic packages
    dnf install -y wireguard-tools podman jq openssl firewalld pciutils python3.11
    systemctl enable --now firewalld
elif [[ "${ID}" == "debian" && "${VERSION_ID}" == "12" ]]; then
    apt-get update
    apt-get -y install gnupg2
    apt-get update
    apt-get -y install python3-venv podman wireguard uuid-runtime jq openssl psmisc firewalld pciutils wget
else
    echo "System not supported"
    exit 1
fi

# Search for core image and prepare modules for next install script
modules=()
for arg in "${@}"; do
    module=$(basename ${arg} | sed s/:.*//)
    if [ "$module" == "core" ]; then
        core_url=${arg}
    else
        modules+=(${arg})
    fi
done

echo "Extracting core sources from ${core_url}:"
cid=$(podman create "${core_url}")
podman export ${cid} | tar --totals -C / --no-overwrite-dir --no-same-owner --exclude=.gitignore --exclude-caches-under -x -v -f - | LC_ALL=C sort | tee coreimage.lst
mv -v coreimage.lst /var/lib/nethserver/node/state/coreimage.lst
podman rm -f ${cid}

/var/lib/nethserver/node/install-core.sh "${modules[@]}"

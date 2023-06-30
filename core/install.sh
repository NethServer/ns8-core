#!/bin/bash

#
# Copyright (C) 2021 Nethesis S.r.l.
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

echo "Checking machine hostname"
fqdn=$(hostname -f)
if [[ -z "${fqdn}" || "${fqdn}" == *localhost* ]]; then
    echo "Current hostname '$fqdn' is not valid. The hostname must not contain 'localhost'."
    echo "Please set a valid FQDN like 'myserver.nethserver.org'."
    exit 2
fi

if ! grep -q -E "\b${fqdn}\b" /etc/hosts ; then
    echo "The machine hostname should be inside /etc/hosts, otherwise the system will be unresponsive."
    echo "You can fix it with the following command:"
    echo
    echo "    echo '127.0.0.1 $fqdn' >> /etc/hosts"
    echo
    echo "Restart the installation when the problem has been solved."
    exit 3
fi

echo "Checking port 80 and 443 are not already in use"
for port in 80 443
do
    if timeout 5 bash -c "< /dev/tcp/localhost/$port" &> /dev/null; then 
        echo "The TCP port $port is already used on this host, we cannot install, we exit"
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
    dnf install -y wireguard-tools podman jq openssl firewalld pciutils
    systemctl enable --now firewalld
elif [[ "${ID}" == "debian" && "${VERSION_ID}" == "11" ]]; then
    apt-get update
    apt-get -y install gnupg2
    # Add repo for podman => 3.4.2
    echo 'deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_11/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
    wget -O - https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_11/Release.key | apt-key add -
    apt-get update
    apt-get -y install python3-venv podman wireguard uuid-runtime jq openssl psmisc firewalld pciutils wget
elif [[ "${ID}" == "ubuntu" && "${VERSION_ID}" == "20.04" && "${CI}" == "true" && "${GITHUB_ACTIONS}" == "true" ]]; then
    apt-get update
    apt-get -y install wireguard firewalld pciutils
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
mkdir -pv /var/lib/nethserver/node/state
cid=$(podman create "${core_url}")
podman export ${cid} | tar --totals -C / --no-overwrite-dir --no-same-owner -x -v -f - | LC_ALL=C sort | tee /var/lib/nethserver/node/state/coreimage.lst
podman rm -f ${cid}

/var/lib/nethserver/node/install-core.sh "${modules[@]}"

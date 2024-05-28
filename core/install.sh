#!/bin/bash

#
# Copyright (C) 2023 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e

print_ns_yum_config ()
{
    cat <<'EOF'
#
# nethserver.repo -- YUM repositories for Rocky Linux, added by NS8 install.sh
#

[ns-baseos]
name=NS8 Rocky Linux $releasever - BaseOS
mirrorlist=http://mirrorlist.nethserver.org/rockylinux?arch=$basearch&repo=BaseOS-$releasever$rltype
gpgcheck=1
enabled=1
countme=1
metadata_expire=6h
gpgkey=http://mirrorlist.nethserver.org/rpm-gpg-key-ns8

[ns-appstream]
name=NS8 Rocky Linux $releasever - AppStream
mirrorlist=http://mirrorlist.nethserver.org/rockylinux?arch=$basearch&repo=AppStream-$releasever$rltype
gpgcheck=1
enabled=1
countme=1
metadata_expire=6h
gpgkey=http://mirrorlist.nethserver.org/rpm-gpg-key-ns8
EOF
}

if [[ $EUID != 0 ]]; then
    echo "This script must be executed with root privileges."
    exit 1
fi

# ensure /usr/local/bin is in the PATH environment variable
if [[ ! "$PATH" =~ (^|:)/usr/local/bin(:|$) ]]; then
    export PATH=$PATH:/usr/local/bin
fi

# ensure /usr/local/sbin is in the PATH environment variable
if [[ ! "$PATH" =~ (^|:)/usr/local/sbin(:|$) ]]; then
    export PATH=$PATH:/usr/local/sbin
fi

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
    if [[ "${ID}" == rocky ]]; then
        print_ns_yum_config > /etc/yum.repos.d/nethserver.repo
        dnf config-manager --save --set-disabled appstream baseos extras
        rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-Rocky-9
    fi
    dnf update -y # Fix SELinux issues with basic packages
    dnf install -y wireguard-tools podman curl jq openssl firewalld pciutils python3.11
    systemctl enable --now firewalld
elif [[ "${ID}" == "debian" && "${VERSION_ID}" == "12" ]]; then
    apt-get update
    apt-get -y install gnupg2
    apt-get update
    apt-get -y install python3-venv podman wireguard uuid-runtime curl jq openssl psmisc firewalld pciutils wget
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
mkdir -vp /var/lib/nethserver/node/state /var/lib/nethserver/cluster/state
chmod -c 0700 /var/lib/nethserver/node/state /var/lib/nethserver/cluster/state
mv -v coreimage.lst /var/lib/nethserver/node/state/coreimage.lst
podman rm -f ${cid}

/var/lib/nethserver/node/install-core.sh "${modules[@]}"

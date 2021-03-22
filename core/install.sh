#!/bin/bash

set -e

distro=$(awk -F = '/^ID=/ { print $2 }' /etc/os-release)

echo "Install dependencies:"
if [[ ${distro} == "fedora" ]]; then
    dnf install -y wireguard-tools podman jq
elif [[ ${distro} == "debian" ]]; then
    # Install podman
    apt-get -y install gnupg2 python3-venv
    grep -q 'http://deb.debian.org/debian buster-backports' /etc/apt/sources.list || echo 'deb http://deb.debian.org/debian buster-backports main' >> /etc/apt/sources.list
    echo 'deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_10/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
    wget -O - https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_10/Release.key | apt-key add -
    apt-get update
    apt-get -y -t buster-backports install libseccomp2 podman

    # Install wireguard
    apt install linux-headers-$(uname -r) -y
    apt install wireguard -y

    # Enable access to journalctl --user
    grep  -e "^#Storage=persistent" /etc/systemd/journald.conf || echo "Storage=persistent" >> /etc/systemd/journald.conf
    systemctl restart systemd-journald
fi

if [ -f /usr/local/etc/registry.json ]; then
    echo "Registry auth found."
    export REGISTRY_AUTH_FILE=/usr/local/etc/registry.json
    chmod -c 644 ${REGISTRY_AUTH_FILE}
fi

echo "Set kernel parameters:"
sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 | tee /etc/sysctl.d/80-nethserver.conf
if [[ ${distro} == "debian" ]]; then
    sysctl -w kernel.unprivileged_userns_clone=1 | tee -a /etc/sysctl.d/80-nethserver.conf
fi

INSTALL_DIR="/usr/local/share"
AGENT_DIR="${INSTALL_DIR}/agent"

echo "Extracting core sources:"
podman pull ghcr.io/nethserver/core:${IMAGE_TAG:-latest}
cid=$(podman create ghcr.io/nethserver/core:${IMAGE_TAG:-latest})
podman export ${cid} | tar -C ${INSTALL_DIR} -x -v -f -
podman rm -f ${cid}

cp -f ${AGENT_DIR}/node-agent.service      /etc/systemd/system/node-agent.service
cp -f ${AGENT_DIR}/redis.service           /etc/systemd/system/redis.service
cp -f ${AGENT_DIR}/module-agent@.service   /etc/systemd/system/module-agent@.service
cp -f ${AGENT_DIR}/module-init@.service    /etc/systemd/system/module-init@.service
cp -f ${AGENT_DIR}/module-agent.service    /etc/systemd/user/module-agent.service
cp -f ${AGENT_DIR}/module-init.service     /etc/systemd/user/module-init.service
cp -f ${AGENT_DIR}/nethserver              /usr/local/bin/nethserver

chmod a+x /usr/local/bin/nethserver

if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to module skeleton dir:"
install -d -m 700 /usr/local/share/module.skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /usr/local/share/module.skel/.ssh/authorized_keys

echo "Setup agent:"
python3 -mvenv ${AGENT_DIR}
${AGENT_DIR}/bin/pip3 install redis

echo "NODE_PREFIX=$(hostname -s)" > /usr/local/etc/node-agent.env

echo "Setup registry:"
if [[ ! -f /usr/local/etc/registry.json ]] ; then
    echo '{"auths":{}}' > /usr/local/etc/registry.json
fi

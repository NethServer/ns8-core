#!/bin/bash

set -e

distro=$(awk -F = '/^ID=/ { print $2 }' /etc/os-release)

echo "Install dependencies:"
if [[ ${distro} == "fedora" ]]; then
    dnf install -y wireguard-tools podman jq
elif [[ ${distro} == "debian" ]]; then
    apt-get -y install gnupg2 python3-venv
    echo 'deb http://deb.debian.org/debian buster-backports main' >> /etc/apt/sources.list
    echo 'deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_10/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
    wget -O - https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_10/Release.key | apt-key add -
    apt-get update
    apt-get -y -t buster-backports install libseccomp2 podman

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

installdir="/usr/local/share"
agentdir="${installdir}/agent"

echo "Extracting core sources:"
podman pull ghcr.io/nethserver/core:${IMAGE_TAG:-latest}
cid=$(podman create ghcr.io/nethserver/core:${IMAGE_TAG:-latest})
podman export ${cid} | tar -C ${installdir} -x -v -f -
podman rm -f ${cid}

cp -f ${agentdir}/node-agent.service      /etc/systemd/system/node-agent.service
cp -f ${agentdir}/redis.service           /etc/systemd/system/redis.service
cp -f ${agentdir}/module-agent@.service   /etc/systemd/system/module-agent@.service
cp -f ${agentdir}/module-init@.service    /etc/systemd/system/module-init@.service
cp -f ${agentdir}/module-agent.service    /etc/systemd/user/module-agent.service
cp -f ${agentdir}/module-init.service     /etc/systemd/user/module-init.service

echo "Setup agent:"
python3 -mvenv ${agentdir}
${agentdir}/bin/pip3 install redis

echo "NODE_PREFIX=$(hostname -s)" > /usr/local/etc/node-agent.env
systemctl enable --now node-agent.service redis.service

if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to module skeleton dir:"
install -d -m 700 /usr/local/share/module.skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /usr/local/share/module.skel/.ssh/authorized_keys

if [[ ! -f /usr/local/etc/registry.json ]] ; then
    echo '{"auths":{}}' > /usr/local/etc/registry.json
    chmod -c 644 /usr/local/etc/registry.json
    echo
    echo "[INFO] Container registry configuration is missing."
    echo "If you want to push images, execute login to image registry:"
    echo "  podman login --authfile /usr/local/etc/registry.json ghcr.io"
    echo "  export REGISTRY_AUTH_FILE=/usr/local/etc/registry.json"
    echo "  chmod -c 644 /usr/local/etc/registry.json"
fi

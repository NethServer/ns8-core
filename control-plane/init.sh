#!/bin/bash

set -e

if [[ ! -f /usr/local/etc/registry.json ]] ; then
    echo "[ERROR] missing the registry access configuration. Copy it to /usr/local/etc/registry.json"
    exit 1
fi

echo "Registry auth found."
export REGISTRY_AUTH_FILE=/usr/local/etc/registry.json
chmod -c 644 ${REGISTRY_AUTH_FILE}

echo "Set kernel parameters:"
sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 | tee /etc/sysctl.d/80-nethserver.conf

echo "Install dependencies:"
dnf install -y wireguard-tools podman jq

installdir="/usr/local/share"
agentdir="${installdir}/agent"
cplanedir="${installdir}/cplane"

echo "Extracting control plane sources:"
podman pull ghcr.io/nethserver/control-plane:latest
cid=$(podman create ghcr.io/nethserver/control-plane:latest)
podman export ${cid} | tar -C ${installdir} -x -v -f -
podman rm -f ${cid}

cp -f ${agentdir}/node-agent.service      /etc/systemd/system/node-agent.service
cp -f ${agentdir}/module-agent.service    /etc/systemd/user/module-agent.service

echo "Setup agent:"
python3 -mvenv ${agentdir}
${agentdir}/bin/pip3 install redis

echo "Starting control plane:"
useradd -m -k ${cplanedir}/skel cplane
loginctl enable-linger cplane

systemctl enable --now node-agent.service

if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to data plane home skeleton dir:"
install -d -m 700 /usr/local/share/dplane/skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /usr/local/share/dplane/skel/.ssh/authorized_keys

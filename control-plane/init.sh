#!/bin/bash

set -e

sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 | tee /etc/sysctl.d/80-nethserver.conf

dnf install -y wireguard-tools podman

if ! podman pull registry.digitalocean.com/nethserver/control-plane:latest ; then
    echo "[ERROR] run: export REGISTRY_AUTH_FILE=path-to/docker-config.json"
    exit 1
fi

installdir="/usr/local/share"
agentdir="${installdir}/agent"
cplanedir="${installdir}/cplane"

cid=$(podman create registry.digitalocean.com/nethserver/control-plane:latest)
podman cp ${cid}:/ ${installdir}
podman rm -f ${cid}

cp -f ${agentdir}/node-agent.service      /etc/systemd/system/node-agent.service
cp -f ${agentdir}/module-agent.service    /etc/systemd/user/module-agent.service

python3 -mvenv ${agentdir}
${agentdir}/bin/pip3 install redis

useradd -m -k ${cplanedir}/skel cplane
loginctl enable-linger cplane

systemctl enable --now node-agent.service

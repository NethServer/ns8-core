#!/bin/bash

set -e

sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 | tee /etc/sysctl.d/80-nethserver.conf

dnf install -y wireguard-tools podman

if ! podman pull registry.digitalocean.com/nethserver/control-plane:latest ; then
    echo "[ERROR] run: export REGISTRY_AUTH_FILE=path-to/docker-config.json"
    exit 1
fi

agentdir="/usr/local/share/agent"

python3 -mvenv ${agentdir}
${agentdir}/bin/pip3 install redis

cid=$(podman create --entrypoint=/ registry.digitalocean.com/nethserver/control-plane:latest)
podman cp ${cid}:agent/service.py ${agentdir}/service.py
podman cp ${cid}:agent/node-events ${agentdir}/node-events
podman cp ${cid}:agent/node-agent.service /etc/systemd/system/node-agent.service
podman cp ${cid}:agent/service-agent.service /etc/systemd/user/service-agent.service
podman cp ${cid}:cplane/skel /etc/opt/cplane/skel
podman rm -f ${cid}

useradd -m -k /etc/opt/cplane/skel cplane
loginctl enable-linger cplane

systemctl enable --now node-agent.service

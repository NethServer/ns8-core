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

installdir="/usr/local/share"
agentdir="${installdir}/agent"

echo "Extracting core sources:"
cid=$(podman create ghcr.io/nethserver/core:${IMAGE_TAG:-latest})
podman export ${cid} | tar -C ${installdir} -x -v -f -
podman rm -f ${cid}

cp -f ${installdir}/etc/containers/containers.conf /etc/containers/containers.conf

cp -f ${agentdir}/node-agent.service      /etc/systemd/system/node-agent.service
cp -f ${agentdir}/redis.service           /etc/systemd/system/redis.service
cp -f ${agentdir}/module-agent@.service   /etc/systemd/system/module-agent@.service
cp -f ${agentdir}/module-init@.service    /etc/systemd/system/module-init@.service
cp -f ${agentdir}/module-agent.service    /etc/systemd/user/module-agent.service
cp -f ${agentdir}/module-init.service     /etc/systemd/user/module-init.service

# Setp redis first, the agent will connect to it
echo "Setup redis:"
systemctl enable --now redis.service
echo -e "Waiting for redis: "
until [ "$(podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PING 2>/dev/null)" == "PONG" ]; do
    sleep 1
done
echo "OK"

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

echo "Setup traefik:"
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik ''
HSET module/traefik0/module.env LE_EMAIL root@$(hostname -f) EVENTS_IMAGE ghcr.io/nethserver/traefik:latest
PUBLISH $(hostname -s):module.init traefik0
EOF

echo "Setup restic server:"
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET module/restic0/module.env EVENTS_IMAGE ghcr.io/nethserver/restic-server:latest
PUBLISH $(hostname -s):module.init restic0
EOF

echo "Setup WireGuard VPN:"
private_key=$(wg genkey)
public_key=$(echo $private_key | wg pubkey)
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET node/$(hostname -s)/vpn PUBLIC_KEY $public_key 
HSET node/$(hostname -s)/vpn IP_ADDRESS 10.5.4.1
HSET node/$(hostname -s)/vpn LISTEN_PORT 55820
PUBLISH $(hostname -s):vpn-master.init $private_key
EOF

echo
echo "WireGuard server public key: $public_key"
echo


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

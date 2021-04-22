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

echo "Set kernel parameters:"
sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 | tee /etc/sysctl.d/80-nethserver.conf
if [[ ${distro} == "debian" ]]; then
    sysctl -w kernel.unprivileged_userns_clone=1 | tee -a /etc/sysctl.d/80-nethserver.conf
fi

echo "Extracting core sources:"
mkdir -pv /var/lib/nethserver/node/state
cid=$(podman create ghcr.io/nethserver/core:latest)
podman export ${cid} | tar -C / -x -v -f - | tee /var/lib/nethserver/node/state/image.lst
podman rm -f ${cid}

if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    echo "Generating a new RSA key pair for SSH:"
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to module skeleton dir:"
install -d -m 700 /etc/nethserver/skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /etc/nethserver/skel/.ssh/authorized_keys

echo "Setup agent:"
agent_dir=/usr/local/nethserver/agent
python3 -mvenv ${agent_dir}
${agent_dir}/bin/pip3 install -U pip
${agent_dir}/bin/pip3 install -r /etc/nethserver/pythonreq.txt

echo "Setup registry:"
if [[ ! -f /etc/nethserver/registry.json ]] ; then
    echo '{"auths":{}}' > /etc/nethserver/registry.json
fi

echo "Generating WireGuard VPN key pair:"
(umask 0077; wg genkey | tee /etc/nethserver/wg0.key | wg pubkey) | tee /etc/nethserver/wg0.pub

echo "Generating API server SECRET:"
(umask 0077; echo "SECRET=$(tr -dc '[:alnum:]' < /dev/urandom | dd bs=4 count=8 status=none)" >/etc/nethserver/api-server.env)

echo "Pull core images":
podman pull docker.io/redis:6-alpine

echo "Start Redis DB and core agents:"
systemctl enable --now redis

# Drop and initialize the whole Redis DB. Just add the keys for the node bootstrap
podman exec -i redis redis-cli >/dev/null <<EOF
FLUSHALL
MULTI
HSET user/default cluster admin
SADD cluster/roles/admin create-cluster add-node
SET cluster/leader 1
SET cluster/node_sequence 1
HSET image/traefik name "Traefik edge proxy" rootfull 0 url ghcr.io/nethserver/traefik:latest
LPUSH node/1/tasks '{"id":"add-traefik1","action":"add-module","data":"{\"image\":\"traefik\"}"}'
EXEC
EOF

echo "Start API server and core agents:"
echo "AGENT_ID=node/1" > /etc/nethserver/node.env
systemctl enable --now api-server.service agent@cluster.service agent@node.service


cat - <<EOF

NethServer 8 scratchpad
--------------------------------------------------------------------------

Congratulations!

This node is now ready to run as a single-node cluster instance of NS8


A. To join this node to an alredy existing cluster run:

      join-cluster <cluster_url>

   For instance:

      join-cluster https://admin:Nethesis,1234@cluster.example.com

B. To initialize this node as a cluster leader run:

      create-cluster <vpn_endpoint_address>:<vpn_endpoint_port> [vpn_cidr] [admin_password]

   For instance:

      create-custer $(hostname -f):55820 10.5.4.0/24 Nethesis,1234

EOF

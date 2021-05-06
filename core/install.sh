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
sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 -w net.ipv4.ip_forward=1 | tee /etc/sysctl.d/80-nethserver.conf
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

if ! grep -q ' cluster-leader$' /etc/hosts; then
    echo "Add /etc/hosts entries:"
    echo "127.0.0.1 cluster-leader" >> /etc/hosts
fi

echo "Generate WireGuard VPN key pair:"
(umask 0077; wg genkey | tee /etc/nethserver/wg0.key | wg pubkey) | tee /etc/nethserver/wg0.pub

echo "Start Redis DB:"
systemctl enable --now redis

echo "Generating cluster password:"
cluster_password=$(podman exec redis redis-cli ACL GENPASS)
cluster_pwhash=$(echo -n "${cluster_password}" | sha256sum | awk '{print $1}')
(umask 0077; exec >/var/lib/nethserver/cluster/state/agent.env
    printf "AGENT_ID=cluster\n"
    printf "REDIS_PASSWORD=%s\n" "${cluster_password}"
    printf "REDIS_ADDRESS=127.0.0.1:6379\n" # Override the cluster-leader /etc/hosts record
)

echo "Generating api-server password:"
apiserver_password=$(podman exec redis redis-cli ACL GENPASS)
apiserver_pwhash=$(echo -n "${apiserver_password}" | sha256sum | awk '{print $1}')
(umask 0077; exec >/etc/nethserver/api-server.env
    printf "REDIS_PASSWORD=%s\n" "${apiserver_password}"
    printf "REDIS_USER=api-server\n"
    printf "REDIS_ADDRESS=127.0.0.1:6379\n" # Override the cluster-leader /etc/hosts record
)

echo "Generating node password:"
node_password=$(podman exec redis redis-cli ACL GENPASS)
node_pwhash=$(echo -n "${node_password}" | sha256sum | awk '{print $1}')
(umask 0077; exec >/var/lib/nethserver/node/state/agent.env
    printf "AGENT_ID=node/1\n"
    printf "REDIS_PASSWORD=%s\n" "${node_password}"
)

(
    # Add the keys for the cluster bootstrap
    cat <<EOF
SADD cluster/roles/admin add-module
SET cluster/leader 1
SET cluster/node_sequence 1
LPUSH cluster/tasks '{"id":"addtraef-ik1x-xxxx-xxxx-xxxxxxxxxxxx","action":"add-module","data":"{\"image\":\"traefik\",\"node\":\"1\"}"}'
EOF

    # Load module images metadata. XXX this is a temporary implementation
    grep '^HSET image/' /var/lib/nethserver/cluster/state/images-catalog.txt

    # Setup initial ACLs
    cat <<EOF
ACL SETUSER cluster ON #${cluster_pwhash} ~* &* +@all
AUTH cluster "${cluster_password}"
ACL SETUSER default ON nopass ~* &* nocommands +@read +@connection +subscribe +psubscribe
ACL SETUSER api-server ON #${apiserver_pwhash} ~* &* nocommands +@read +@pubsub +lpush +@transaction +@connection
ACL SETUSER node/1 ON #${node_pwhash} resetkeys ~node/1/* resetchannels &progress/task/* nocommands +@read +@write +@transaction +@connection +publish +psync +replconf +ping
ACL SAVE
SAVE
EOF

) | podman exec -i redis redis-cli

echo "Start API server and core agents:"
systemctl enable --now api-server.service agent@cluster.service agent@node.service


cat - <<EOF

NethServer 8 Scratchpad
--------------------------------------------------------------------------

Congratulations!

This node is now ready to run as a standalone instance of NS8 Scratchpad

Open a new login shell or type the following command to fix the environment:

    source /etc/profile.d/nethserver.sh


A. To join this node to an alredy existing cluster run:

      join-cluster <cluster_url> <jwt_auth>

   For instance:

      join-cluster https://cluster.example.com eyJhbGc...NiIsInR5c

B. To initialize this node as a cluster leader run:

      create-cluster <vpn_endpoint_address>:<vpn_endpoint_port> [vpn_cidr] [admin_password]

   For instance:

      create-cluster $(hostname -f):55820 10.5.4.0/24 Nethesis,1234

EOF

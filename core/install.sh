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
repobase=${REPOBASE:-ghcr.io/nethserver}

source /etc/os-release

echo "Install dependencies:"
if [[ ${ID} == "centos" && "${PLATFORM_ID}" == "platform:el9" ]]; then
    dnf install -y wireguard-tools podman jq openssl
    systemctl disable --now firewalld || :
elif [[ "${ID}" == "debian" && "${VERSION_ID}" == "11" ]]; then
    apt-get update
    apt-get -y install gnupg2 python3-venv podman wireguard uuid-runtime jq openssl
elif [[ "${ID}" == "ubuntu" && "${VERSION_ID}" == "20.04" && "${CI}" == "true" && "${GITHUB_ACTIONS}" == "true" ]]; then
    apt-get update
    apt-get -y install wireguard
else
    echo "System not supported"
    exit 1
fi

if [[ $(journalctl --output cat -n 1) == "" ]]; then
    echo "Journal is empty, attempting a systemd-journald restart:"
    systemctl restart systemd-journald.service
fi

echo "Set kernel parameters:"
sysctl -w net.ipv4.ip_unprivileged_port_start=23 -w user.max_user_namespaces=28633 -w net.ipv4.ip_forward=1 | tee /etc/sysctl.d/80-nethserver.conf

# Pull modules from the given development branch
imagetag="$1"
if [[ -n "${imagetag}" ]]; then
    shift
    for module in "${@}"; do
        podman pull "${repobase}/${module}:${imagetag}"
        echo "Tagging development branch ${module}:${imagetag} => ${module}:ns8-stable"
        podman tag "${repobase}/${module}:${imagetag}" "${repobase}/${module}:ns8-stable"
    done
    shift $# # Discard all arguments
fi

echo "Extracting core sources:"
mkdir -pv /var/lib/nethserver/node/state
cid=$(podman create "${repobase}/core:ns8-stable")
podman export ${cid} | tar --totals -C / --no-overwrite-dir --no-same-owner -x -v -f - | sort | tee /var/lib/nethserver/node/state/coreimage.lst
podman rm -f ${cid}

source /etc/nethserver/core.env
echo "Pulling rclone image ${RCLONE_IMAGE}:"
podman pull "${RCLONE_IMAGE}"


if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    echo "Generating a new RSA key pair for SSH:"
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to module skeleton dir:"
install -d -m 700 /etc/nethserver/skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /etc/nethserver/skel/.ssh/authorized_keys

if ! id "api-server" &>/dev/null; then
    echo "Create the api-server user:"
    useradd -r -m -d /var/lib/nethserver/api-server api-server
fi

echo "Setup Python virtual environment for agents:"
core_dir=/usr/local/agent/pyenv
python3 -mvenv ${core_dir} --upgrade-deps
${core_dir}/bin/pip3 install -r /etc/nethserver/pythonreq.txt
echo "/usr/local/agent/pypkg" >$(${core_dir}/bin/python3 -c "import sys; print(sys.path[-1] + '/pypkg.pth')")

echo "Setup registry:"
if [[ ! -f /etc/nethserver/registry.json ]] ; then
    echo '{"auths":{}}' > /etc/nethserver/registry.json
fi

if ! grep -q ' cluster-leader$' /etc/hosts; then
    echo "Add /etc/hosts entries:"
    echo "127.0.0.1 cluster-leader" >> /etc/hosts
    echo "127.0.0.1 cluster-localnode" >> /etc/hosts
fi

echo "Generate WireGuard VPN key pair:"
(umask 0077; wg genkey | tee /etc/nethserver/wg0.key | wg pubkey) | tee /etc/nethserver/wg0.pub

echo "Start Redis DB:"
systemctl enable --now redis

echo "Generating cluster password:"
cluster_password=$(podman exec redis redis-cli ACL GENPASS)
cluster_pwhash=$(echo -n "${cluster_password}" | sha256sum | awk '{print $1}')
(umask 0077; exec >/var/lib/nethserver/cluster/state/agent.env
    printf "AGENT_BASEACTIONS_DIR=\n"
    printf "AGENT_ID=cluster\n" # Override value from agent@.service
    printf "REDIS_USER=cluster\n"
    printf "REDIS_PASSWORD=%s\n" "${cluster_password}"
    printf "REDIS_ADDRESS=127.0.0.1:6379\n" # Override the cluster-leader /etc/hosts record
)
printf "NODE_ID=1\n" > /var/lib/nethserver/cluster/state/environment

echo "Generating api-server password:"
apiserver_password=$(podman exec redis redis-cli ACL GENPASS)
apiserver_pwhash=$(echo -n "${apiserver_password}" | sha256sum | awk '{print $1}')
apiserver_jwt_secret=$(uuidgen)
(umask 0077; exec >/etc/nethserver/api-server.env
    printf "REDIS_PASSWORD=%s\n" "${apiserver_password}"
    printf "REDIS_USER=api-server\n"
    printf "REDIS_ADDRESS=127.0.0.1:6379\n" # Override the cluster-leader /etc/hosts record
    printf "SECRET=%s\n" "${apiserver_jwt_secret}"
)

echo "Generating node password:"
node_password=$(podman exec redis redis-cli ACL GENPASS)
node_pwhash=$(echo -n "${node_password}" | sha256sum | awk '{print $1}')
(umask 0077; exec >/var/lib/nethserver/node/state/agent.env
    printf "AGENT_BASEACTIONS_DIR=\n"
    printf "AGENT_ID=node/1\n" # Override value from agent@.service
    printf "REDIS_USER=node/1\n"
    printf "REDIS_PASSWORD=%s\n" "${node_password}"
)
printf "NODE_ID=1\n" > /var/lib/nethserver/node/state/environment

(
    # Add the keys for the cluster bootstrap
    cat <<EOF
SET cluster/node_sequence 1
SET node/1/tcp_ports_sequence 20000
EOF

    # Configure default module repository
    # FIXME: remove testing flag before official release
    cat <<EOF
HSET cluster/repository/default url https://raw.githubusercontent.com/NethServer/ns8-repomd/${REPOBRANCH:-repomd}/ status 1 testing 1
EOF

    # Setup initial ACLs
    cat <<EOF
ACL SETUSER cluster ON #${cluster_pwhash} ~* &* +@all
AUTH cluster "${cluster_password}"
ACL SETUSER default ON nopass ~* &* nocommands +@read +@connection +subscribe +psubscribe +psync +replconf +ping
ACL SETUSER api-server ON #${apiserver_pwhash} ~* &* nocommands +@read +@pubsub +lpush +@transaction +@connection +role
ACL SETUSER node/1 ON #${node_pwhash} resetkeys ~node/1/* resetchannels &progress/node/1/* &node/1/event/* nocommands +@read +@write +@transaction +@connection +publish
ACL SAVE
SAVE
EOF

) | redis-cli

echo "Start API server and core agents:"
systemctl enable --now api-server.service agent@cluster.service agent@node.service

echo "Grant initial permissions:"
runagent python3 <<'EOF'
import agent
import cluster.grants
rdb = agent.redis_connect(privileged=True)
cluster.grants.grant(rdb, action_clause="*",      to_clause="owner",  on_clause='cluster')
cluster.grants.grant(rdb, action_clause="list-*", to_clause="reader", on_clause='cluster')
cluster.grants.grant(rdb, action_clause="get-*",  to_clause="reader", on_clause='cluster')
cluster.grants.grant(rdb, action_clause="show-*", to_clause="reader", on_clause='cluster')
cluster.grants.grant(rdb, action_clause="read-*", to_clause="reader", on_clause='cluster')
cluster.grants.grant(rdb, action_clause="*",      to_clause="owner",  on_clause='node/1')
cluster.grants.grant(rdb, action_clause="list-*", to_clause="reader", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="get-*",  to_clause="reader", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="show-*", to_clause="reader", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="read-*", to_clause="reader", on_clause='node/1')
EOF

echo "Install Traefik:"
add-module ${repobase}/traefik:ns8-stable 1

echo "Setting default admin password:"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-Nethesis,1234}"
add-user --role owner --password "${ADMIN_PASSWORD}" admin

echo "Enable the events gateway for the node agent:"
systemctl enable --now eventsgw@node

cat - <<EOF

NethServer 8 Scratchpad
----------------------------------------------------------------------------

Finish the cluster configuration by running one of the following procedures.

A. To join this node to an already existing cluster run:

      join-cluster <cluster_url> <jwt_auth>

   For instance:

      join-cluster https://cluster.example.com eyJhbGc...NiIsInR5c

B. To initialize this node as a cluster leader run:

      create-cluster <vpn_endpoint_address>:<vpn_endpoint_port> [vpn_cidr] [admin_password]

   For instance:

      create-cluster $(hostname -f):55820 10.5.4.0/24 ${ADMIN_PASSWORD}

Finally, access the administration UI at:

   https://$(hostname -f)/cluster-admin/

   User: admin
   Password: ${ADMIN_PASSWORD}
EOF

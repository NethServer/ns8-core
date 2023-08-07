#!/bin/bash

#
# Copyright (C) 2023 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e

#
# Finalize NS8 installation
#

source /etc/nethserver/core.env

if [[ ! -f ~/.ssh/id_rsa.pub ]] ; then
    echo "Generating a new RSA key pair for SSH:"
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa
fi

echo "Adding id_rsa.pub to module skeleton dir:"
install -d -m 700 /etc/nethserver/skel/.ssh
install -m 600 -T ~/.ssh/id_rsa.pub /etc/nethserver/skel/.ssh/authorized_keys

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

(
    # Add the keys for the cluster bootstrap
    cat <<EOF
SET cluster/node_sequence 1
SET node/1/tcp_ports_sequence 20000
EOF

    # Configure default module repository
    cat <<EOF
HSET cluster/repository/default url ${REPOMD:-https://distfeed.nethserver.org/ns8/updates/}/ status 1 testing ${TESTING:-0}
EOF

    # Setup initial ACLs
    cat <<EOF
ACL SETUSER cluster ON #${cluster_pwhash} ~* &* +@all
AUTH cluster "${cluster_password}"
ACL SETUSER default ON nopass resetkeys ~cluster/* ~node/* ~module/* resetchannels &* nocommands +@read +@connection +subscribe +psubscribe -@admin
ACL SETUSER api-server ON #${apiserver_pwhash} ~* &* nocommands +@read +@pubsub +lpush +@transaction +@connection +role +hset
ACL SETUSER node/1 ON #${node_pwhash} resetkeys ~node/1/* ~task/node/1/* resetchannels &progress/node/1/* &node/1/event/* nocommands +@read +@write +@transaction +@connection +publish -@admin +psync +replconf
ACL SAVE
SAVE
CONFIG SET masteruser node/1
CONFIG SET masterauth ${node_password}
CONFIG REWRITE
EOF

    printf 'SET cluster/ui_name "%s"\n' "${CLUSTER_NAME:-NethServer 8}"

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
cluster.grants.grant(rdb, action_clause="add-public-service",  to_clause="fwadm", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="remove-public-service",  to_clause="fwadm", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="add-custom-zone",  to_clause="fwadm", on_clause='node/1')
cluster.grants.grant(rdb, action_clause="remove-custom-zone",  to_clause="fwadm", on_clause='node/1')
EOF

for arg in "${@}"; do
    module=$(basename ${arg} | sed s/:.*//)
    echo "Adding module override:  ${module} -> ${arg}"
    redis-cli hset cluster/override/modules ${module} ${arg}
done

echo "Install Traefik:"
add-module traefik 1

echo "Setting default admin password:"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-Nethesis,1234}"
add-user --role owner --password "${ADMIN_PASSWORD}" admin

cat - <<EOF

NethServer 8 Core
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

   https://<hostname_or_IP>/cluster-admin/

For instance, if $(hostname -f) is resolvable 

   https://$(hostname -f)/cluster-admin/

Enter the following credentials:

   User: admin
   Password: ${ADMIN_PASSWORD}
EOF

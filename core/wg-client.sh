#!/bin/bash

# Usage example: wg-client.sh xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx 1.2.3.4:55820 10.5.4.2
# Then, add to the server:
# 
# [Peer]
# PublicKey = <client_public_ip>
# AllowedIPs = 10.5.4.2/32
#

set -e 

if [[ -z "$1" ]] || [[ -z "$2" ]] || [[ -z "$3" ]]; then
    echo "Usage: $0 <server_public_key> <server_uri> <vpn_ip>"
    exit 1
fi

server_key=$1
server_uri=$2
vpn_ip=$3

VPN_OUT_DIR=/etc/wireguard
mkdir -p $VPN_OUT_DIR
chmod 700 $VPN_OUT_DIR

if [ ! -f $VPN_OUT_DIR/privatekey ]; then
    umask 0077
    wg genkey | tee $VPN_OUT_DIR/privatekey | wg pubkey > $VPN_OUT_DIR/publickey
    umask 0022
fi

net=$(echo $vpn_ip | cut -d'.' -f1,2,3)".0"
cat > $VPN_OUT_DIR/wg0.conf << EOF
[Interface]
Address = $vpn_ip/32
PrivateKey = $(cat $VPN_OUT_DIR/privatekey)

[Peer]
PublicKey = $server_key
Endpoint = $server_uri
AllowedIPs   = $net/24
PersistentKeepalive = 30
EOF

systemctl enable --now wg-quick@wg0

echo
echo "CLIENT PUBLIC KEY: $(cat $VPN_OUT_DIR/publickey)"

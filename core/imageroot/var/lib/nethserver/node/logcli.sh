#!/bin/bash

#
#Set as default server the vpn ip address of the leader node
#

if [[ -z "${LOKI_ADDR}" ]]; then
	leader_node=$(echo  HGET cluster/environment NODE_ID | redis-cli | tr -d '"')
	export LOKI_ADDR="http://$(echo HGET node/${leader_node}/vpn ip_address | redis-cli | tr -d '"' ):3100"
fi

/usr/local/sbin/logcli "${@}"

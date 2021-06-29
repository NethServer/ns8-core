#!/bin/bash

#
#Set as default server of the cluster
#

if [[ -z "${LOKI_ADDR}" ]]; then
	loki_id=$(REDIS_USER=default redis-exec GET cluster/default_instance/loki)
	export LOKI_ADDR="http://$(REDIS_USER=default redis-exec HGET module/${loki_id}/environment LOKI_ADDR):3100"
fi

/usr/local/sbin/logcli "${@}"

#!/usr/local/bin/runagent bash

#
#Set as default server of the cluster
#

if [[ -z "${LOKI_ADDR}" ]]; then
	loki_id=$(REDIS_USER=default redis-exec GET cluster/default_instance/loki)
	loki_http_port=$(REDIS_USER=default redis-exec HGET module/${loki_id}/environment LOKI_HTTP_PORT)
	export LOKI_ADDR="http://$(REDIS_USER=default redis-exec HGET module/${loki_id}/environment LOKI_ADDR):${loki_http_port}"
	export LOKI_USERNAME=$(REDIS_USER=default redis-exec HGET module/${loki_id}/environment LOKI_API_AUTH_USERNAME)
	export LOKI_PASSWORD=$(REDIS_USER=default redis-exec HGET module/${loki_id}/environment LOKI_API_AUTH_PASSWORD)
fi

/usr/local/sbin/logcli "${@}"

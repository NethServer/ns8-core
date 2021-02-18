#!/bin/bash

#
# Initialize redis with
#
# HSET nsdc.env NSDC_HOSTNAME nsdc1.ad.dp.nethserver.net NSDC_NBDOMAIN AD NSDC_REALM AD.DP.NETHSERVER.NET NSDC_ADMINPASS Nethesis,1234 NSDC_IPADDRESS 10.135.0.2
#

set -e

image="docker.io/nowsci/samba-domain"

podman volume create nsdc-data
podman volume create nsdc-config

podman image pull "${image}"

source <(ctlp-hgetall nsdc.env)

#
# Override the provisioning procedure:
# 1. add --use-ntvfs to run in a rootless container
# 2. bind the main IP address only (and 127.0.0.1) to avoid resolved conflicts
#

cat - > nsdc-init.sh <<EOF
#!/bin/bash
if ! grep -q -F ${NSDC_REALM} /etc/samba/smb.conf; then
    rm -f /etc/samba/smb.conf
    samba-tool domain provision \
        --server-role=dc --domain=${NSDC_NBDOMAIN} --realm=${NSDC_REALM} \
        --adminpass=${NSDC_ADMINPASS:-Nethesis,1234} \
        --use-ntvfs \
        --host-ip=${NSDC_IPADDRESS} \
        --option="bind interfaces only = yes" \
        --option="interfaces = 127.0.0.1 ${NSDC_IPADDRESS}"
    cp -av /var/lib/samba/private/krb5.conf /etc/krb5.conf
fi

samba -i --debug-stderr
EOF
chmod +x nsdc-init.sh

podman create \
    --name=samba \
    --dns="${NSDC_IPADDRESS}" \
    --add-host="${NSDC_HOSTNAME}:${NSDC_IPADDRESS}" \
    --network=host \
    --hostname=${NSDC_HOSTNAME} \
    --env-file=<(ctlp-hgetall nsdc.env) \
    --privileged=true \
    --volume=nsdc-data:/var/lib/samba \
    --volume=nsdc-config:/etc/samba \
    --volume=${PWD}/nsdc-init.sh:/init.sh \
    "${image}" /init.sh

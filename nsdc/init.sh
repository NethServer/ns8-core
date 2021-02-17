#!/bin/bash

set -e

image="docker.io/nowsci/samba-domain"

podman volume create nsdc-data
podman volume create nsdc-config

podman image pull "${image}"

# Why not rootless:
#
# 1. podman pod create --network=host doesn't work when rootless 
# 2. samba provisioning crashes: still need to run rootful (needs to investigate xattr support)

podman create \
    --network=host \
    --hostname=$(ctlp-get nsdc.hostname) \
    --name=samba \
    --env-file=<(ctlp-hgetall nsdc.env) \
    --privileged=true \
    --volume=nsdc-data:/var/lib/samba \
    --volume=nsdc-config:/etc/samba/external \
    "${image}"

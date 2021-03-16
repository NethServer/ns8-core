#!/bin/bash

set -e

image="mail-dovecot"
container=$(buildah from docker.io/library/alpine:3.13)

buildah run "${container}" /bin/sh <<EOF
addgroup -S vmail
adduser -G vmail -h /var/lib/vmail -S vmail
apk add --no-cache dovecot dovecot-ldap dovecot-pigeonhole-plugin dovecot-submissiond
find /var/cache/apk /var/lib/apk -type f -delete
EOF
buildah copy "${container}" config /srv/mail/.config
buildah add "${container}" dovecot/ /etc/dovecot/
buildah config \
    --cmd='["/bin/sh", "/etc/dovecot/start.sh"]' \
    --label "org.nethserver.initroot=/srv/mail" \
    "${container}"
buildah commit "${container}" "${image}"

echo
echo "Publish the image with:"
echo
echo " buildah push ${image} docker://ghcr.io/nethserver/${image}:latest"

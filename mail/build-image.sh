#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

#
# Root image
#
reponame="mail"
container=$(buildah from scratch)
buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

#
# Postfix addtional image
#
reponame="mail-postfix"
container=$(buildah from docker.io/library/alpine:3.13)
buildah run "${container}" /bin/sh <<EOF
apk add --no-cache postfix gettext
find /var/cache/apk /var/lib/apk -type f -delete
EOF
buildah add "${container}" postfix/ /etc/postfix/
buildah config --cmd='["/bin/sh", "/etc/postfix/start.sh"]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

#
# Dovecot additional image
# 
reponame="mail-dovecot"
container=$(buildah from docker.io/library/alpine:3.13)

buildah run "${container}" /bin/sh <<EOF
addgroup -S vmail
adduser -G vmail -h /var/lib/vmail -S vmail
apk add --no-cache dovecot dovecot-ldap dovecot-pigeonhole-plugin dovecot-submissiond dovecot-pop3d dovecot-lmtpd
find /var/cache/apk /var/lib/apk -type f -delete
EOF
buildah add "${container}" dovecot/ /etc/dovecot/
buildah config --cmd='["/bin/sh", "/etc/dovecot/start.sh"]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

printf "%s\n" "${messages[@]}"

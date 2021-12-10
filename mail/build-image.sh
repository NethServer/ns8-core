#!/bin/bash

set -e
images=()
repobase="${REPOBASE:-ghcr.io/nethserver}"

#
# Root image
#
reponame="mail"
container=$(buildah from scratch)
buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui /ui
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
images+=("${repobase}/${reponame}")

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
images+=("${repobase}/${reponame}")

#
# Dovecot additional image
# 
reponame="mail-dovecot"
container=$(buildah from docker.io/library/alpine:3.13)
buildah run "${container}" /bin/sh <<EOF
addgroup -S vmail
adduser -G vmail -h /var/lib/vmail -S vmail
apk add --no-cache dovecot dovecot-ldap dovecot-pigeonhole-plugin dovecot-submissiond dovecot-pop3d dovecot-lmtpd gettext
find /var/cache/apk /var/lib/apk -type f -delete
EOF
buildah add "${container}" dovecot/ /etc/dovecot/
buildah config --cmd='["/bin/sh", "/etc/dovecot/start.sh"]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
images+=("${repobase}/${reponame}")

#
#
#

if [[ -n "${CI}" ]]; then
    # Set output value for Github Actions
    printf "::set-output name=images::%s\n" "${images[*]}"
else
    printf "Publish the images with:\n\n"
    for image in "${images[@]}"; do printf "  buildah push %s docker://%s:latest\n" "${image}" "${image}" ; done
    printf "\n"
fi

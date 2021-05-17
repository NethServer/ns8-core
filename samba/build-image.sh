#!/bin/bash

set -e
messages=("Publish the image with:" "")

repobase="ghcr.io/nethserver"
reponame="ubuntu-samba"

if ! buildah inspect --type image "${repobase}/${reponame}" &>/dev/null; then
    container=$(buildah from docker.io/library/ubuntu:rolling)
    buildah run "${container}" -- bash <<EOF
set -e
apt-get update
apt-get -y install samba winbind krb5-user iputils-ping
apt-get clean
find /var/lib/apt/lists/ -type f -delete
EOF
    buildah commit "${container}" "${repobase}/${reponame}"
fi

container=$(buildah from "${repobase}/${reponame}")

reponame="samba"
buildah add "${container}" imageroot/ /srv/imageroot/
buildah add "${container}" scripts/entrypoint.sh /entrypoint.sh
buildah config --label 'org.nethserver/imageroot=/srv/imageroot' "${container}"
buildah config --cmd='' "${container}"
buildah config --entrypoint='[ "/bin/bash", "/entrypoint.sh" ]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"

messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

printf "%s\n" "${messages[@]}"

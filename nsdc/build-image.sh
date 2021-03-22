#!/bin/bash

set -e
repobase="ghcr.io/nethserver"
reponame="samba"

if ! buildah inspect --type image "${repobase}/${reponame}" &>/dev/null; then
    container=$(buildah from docker.io/library/ubuntu:rolling)
    buildah run "${container}" -- bash <<EOF
set -e
apt-get update
apt-get -y install samba winbind
apt-get clean
find /var/lib/apt/lists/ -type f -delete
EOF
    buildah commit "${container}" "${repobase}/${reponame}"
fi

container=$(buildah from "${repobase}/${reponame}")
reponame="nsdc"
buildah copy "${container}" module-events /srv/module-events
buildah copy "${container}" entrypoint.sh /entrypoint.sh
buildah config --label 'org.nethserver.initroot=/srv' "${container}"
buildah config --cmd='' "${container}"
buildah config --entrypoint='[ "/bin/bash", "/entrypoint.sh" ]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"

echo
echo "Publish the images with:"
echo
echo " buildah push "${repobase}/${reponame}" docker://${repobase}/${reponame}"

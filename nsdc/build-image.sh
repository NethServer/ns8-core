#!/bin/bash

set -e

image1="samba"
if ! buildah inspect --type image localhost/${image1} &>/dev/null; then
    container=$(buildah from docker.io/library/ubuntu:rolling)
    buildah run "${container}" -- bash <<EOF
set -e
apt-get update
apt-get -y install samba winbind
apt-get clean
find /var/lib/apt/lists/ -type f -delete
EOF
    buildah commit "${container}" "${image1}"
fi

image2="nsdc"
container=$(buildah from localhost/${image1})
buildah copy "${container}" module-events /srv/module-events
buildah copy "${container}" entrypoint.sh /entrypoint.sh
buildah config --label 'org.nethserver.initroot=/srv' "${container}"
buildah config --cmd='' "${container}"
buildah config --entrypoint='[ "/bin/bash", "/entrypoint.sh" ]' "${container}"
buildah commit "${container}" "${image2}"

echo
echo "Publish the images with:"
echo
echo " buildah push ${image2} docker://ghcr.io/nethserver/${image2}:latest"

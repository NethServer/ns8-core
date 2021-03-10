#!/bin/bash

set -e

image1="cplane-nsdc"
container=$(buildah from scratch)
buildah copy "${container}" module-events /module-events
buildah config --entrypoint='["/"]' "${container}"
buildah commit "${container}" "${image1}"

image2="nsdc"
container=$(buildah from docker.io/library/ubuntu:rolling)
buildah run "${container}" -- bash <<EOF
set -e
apt-get update
apt-get -y install samba winbind
apt-get clean
find /var/lib/apt/lists/ -type f -delete
EOF
buildah config --cmd='' --entrypoint='["/entrypoint.sh"]' "${container}"
buildah copy "${container}" entrypoint.sh /entrypoint.sh
buildah commit "${container}" "${image2}"

echo
echo "Publish the images with:"
echo
echo " buildah push ${image1} docker://ghcr.io/nethserver/${image1}:latest"
echo " buildah push ${image2} docker://ghcr.io/nethserver/${image2}:latest"

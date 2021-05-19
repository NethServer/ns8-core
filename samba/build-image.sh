#!/bin/bash

set -e
images=()

repobase="ghcr.io/nethserver"
reponame="ubuntu-samba"

container="ubuntu-working-container"
# Prepare a local Ubuntu-based samba image
if ! buildah inspect --type container "${container}" &>/dev/null; then
    container=$(buildah from --name "${container}" docker.io/library/ubuntu:rolling)
    buildah run "${container}" -- bash <<EOF
set -e
apt-get update
apt-get -y install samba winbind krb5-user iputils-ping
apt-get clean
find /var/lib/apt/lists/ -type f -delete
EOF
    buildah commit "${container}" "${repobase}/${reponame}"
fi

#
# Add our entrypoint script to build samba-dc
#
container=$(buildah from "${repobase}/${reponame}")
reponame="samba-dc"
buildah add "${container}" scripts/entrypoint.sh /entrypoint.sh
buildah config --cmd='' --entrypoint='[ "/bin/bash", "/entrypoint.sh" ]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
images+=("${repobase}/${reponame}")

#
# Imageroot samba
#
container=$(buildah from scratch)
reponame="samba"
buildah add "${container}" imageroot /
buildah config --label 'org.nethserver/rootfull=1' --entrypoint=/ "${container}"
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

#!/bin/bash

set -e
images=()
repobase="ghcr.io/nethserver"

reponame="ldapproxy"
container=$(buildah from scratch)

buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui /ui
buildah config --label 'org.nethserver/register_default_instance=1' --entrypoint=/ "${container}"
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

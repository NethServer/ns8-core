#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

reponame="ldapproxy"
container=$(buildah from scratch)

buildah add "${container}" imageroot /
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

#
#
#

printf "%s\n" "${messages[@]}"

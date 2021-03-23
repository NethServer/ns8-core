#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

#
# Root image
#
reponame="ldapproxy"
container=$(buildah from scratch)
buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

#
#
#

printf "%s\n" "${messages[@]}"

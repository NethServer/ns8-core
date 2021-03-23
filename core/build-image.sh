#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

container=$(buildah from scratch)
reponame="core"

buildah copy ${container} etc /etc
buildah copy ${container} agent /agent
buildah copy ${container} module.skel /module.skel
buildah config --entrypoint=/ ${container}
buildah commit "${container}" "${repobase}/${reponame}"
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

#
#
#

printf "%s\n" "${messages[@]}"

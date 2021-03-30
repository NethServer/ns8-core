#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

reponame="netdata"
container=$(buildah from scratch)

buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest

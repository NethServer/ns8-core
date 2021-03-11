#!/bin/bash

set -e

image="traefik"
container=$(buildah from scratch)

buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${image}"

echo
echo "Publish the image with:"
echo
echo " buildah push ${image} docker://ghcr.io/nethserver/${image}:latest"

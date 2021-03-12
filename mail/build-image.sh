#!/bin/bash

set -e

image="mail"
container=$(buildah from docker.io/library/python:3.9-alpine)

buildah copy "${container}" config /.config
buildah config \
    "${container}"
buildah commit "${container}" "${image}"

echo
echo "Publish the image with:"
echo
echo " buildah push ${image} docker://ghcr.io/nethserver/${image}:latest"

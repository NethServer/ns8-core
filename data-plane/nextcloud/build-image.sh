#!/bin/bash

set -e

image="nextcloud"
container=$(buildah from scratch)

buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${image}"

echo
echo "Access DigitalOcean control panel and download docker-config.json file from Container Registry."
echo "Then publish the image with:"
echo
echo " export REGISTRY_AUTH_FILE=path-to/docker-config.json"
echo " buildah push ${image} docker://ghcr.io/nethserver/${image}:latest"

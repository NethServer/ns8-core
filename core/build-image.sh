#!/bin/bash

set -e

# export IMAGE_TAG=$(git rev-parse --abbrev-ref HEAD)
image_tag=${IMAGE_TAG:-latest}

container=$(buildah from scratch)

buildah copy ${container} agent /agent
buildah copy ${container} redis.skel /redis.skel
buildah copy ${container} module.skel /module.skel
buildah config --entrypoint=/ ${container}
buildah commit ${container} core

echo
echo "Publish the image with:"
echo
echo " buildah push core docker://ghcr.io/nethserver/core:${image_tag}"

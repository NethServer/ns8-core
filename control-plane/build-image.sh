#!/bin/bash

set -e

container=$(buildah from scratch)

buildah copy ${container} agent /agent
buildah copy ${container} cplane /cplane
buildah copy ${container} dplane /dplane
buildah config --entrypoint=/ ${container}
buildah commit ${container} control-plane

echo
echo "Publish the image with:"
echo
echo " buildah push control-plane docker://ghcr.io/nethserver/control-plane:latest"

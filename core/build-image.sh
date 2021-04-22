#!/bin/bash

set -e
messages=("Publish the image with:" "")
repobase="ghcr.io/nethserver"

# Reuse existing gobuilder-core container, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q gobuilder-core; then
    buildah from --name gobuilder-core -v ${PWD}:/usr/src/core:Z golang:1.16-alpine
fi

# Build statically linked Go binaries based on Musl
buildah run gobuilder-core sh -c "cd /usr/src/core/agent      && CGO_ENABLED=0 go build -v ."
buildah run gobuilder-core sh -c "cd /usr/src/core/api-server && CGO_ENABLED=0 go build -v ."

# Build the core image
container=$(buildah from scratch)
reponame="core"
buildah add ${container} imageroot /
buildah add ${container} agent/agent /usr/local/bin/agent
buildah add ${container} api-server/api-server /usr/local/bin/api-server
buildah config --entrypoint=/ ${container}
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm ${container}
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

printf "%s\n" "${messages[@]}"

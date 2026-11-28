#!/bin/bash

set -e

if ! command -v buildah; then
    # The package manager depends on the distro
    # Pick up the first one that is available
    if command -v apt; then
        apt install buildah
        trap "apt remove -y buildah" EXIT
    elif command -v dnf; then
        dnf install -q -y buildah
        trap "dnf remove -q -y buildah" EXIT
    else
        echo "Error: cannot install buildah! Package manager not found." 1>&2
        exit 1
    fi
fi

image_tag="${1:-latest}"
authorizations="${2:-}"

container=$(buildah from scratch)

buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui /ui
buildah config --entrypoint='["/"]' "${container}"
if [[ -n "${authorizations}" ]]; then
    buildah config --label="org.nethserver.authorizations=${authorizations}" "${container}"
fi
buildah commit "${container}" "localhost/dummy:${image_tag}"

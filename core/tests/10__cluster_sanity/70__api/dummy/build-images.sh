#!/bin/bash -x

container=$(buildah from scratch)

buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui /ui
buildah config --entrypoint='["/"]' "${container}"
buildah commit "${container}" "localhost/dummy:latest"

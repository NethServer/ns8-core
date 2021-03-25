#!/bin/bash

set -e

image="nextcloud"
container=$(buildah from scratch)

buildah copy "${container}" config /.config
buildah config --entrypoint=/ "${container}"
buildah commit "${container}" "${image}"

buildah --authfile /usr/local/etc/registry.json push nextcloud docker://ghcr.io/nethserver/nextcloud:latest

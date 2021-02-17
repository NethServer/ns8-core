#!/bin/bash

set -e

image="docker.io/library/redis:latest"

podman volume create redis-data

podman image pull "${image}"
podman pod create --name=ctlp --publish 6379:6379
podman create --pod=ctlp --name=redis --volume redis-data:/data "${image}"
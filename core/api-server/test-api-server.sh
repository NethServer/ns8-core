#!/bin/bash

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: AGPL-3.0-or-later
#

#
# Black-box tests of the API server. Run this script in this directory with
#
#     bash test-api-server.sh [ROBOT ARG]...
#
# The API server binary is built in a disposable container, then it runs
# in the test container, connected to a Redis server of the same version
# shipped by the core image.
#

set -e

# GOLANG_IMAGE and GOCACHE_VOLUME (a volume name or a host directory) can
# be overridden, e.g. to reuse images and caches already available
golang_image=${GOLANG_IMAGE:-docker.io/library/golang:1.26.10-bookworm}
gocache_volume=${GOCACHE_VOLUME:-gocache-api-server}
redis_image=docker.io/library/redis:8.8.4-alpine

# Build the API server binary. The Go module and build caches persist in
# the cache volume, to speed up subsequent builds.
podman run -i -q --rm \
    --volume="${PWD}:/usr/src:z" \
    --volume="${gocache_volume}:/root/.cache:z" \
    --env=GOMODCACHE=/root/.cache/mod \
    --env=GOCACHE=/root/.cache/build \
    --env=GOTOOLCHAIN="${GOTOOLCHAIN:-local}" \
    --env=CGO_ENABLED=1 \
    --workdir=/usr/src \
    "${golang_image}" \
    go build -tags sqlite_omit_load_extension -o api-server api-server.go

# Start a fresh Redis server in a separate network namespace
podman run -d -q --rm --replace \
    --name test-api-server-redis \
    --network=private \
    "${redis_image}" \
    --port 6379 --save '' --protected-mode no

#
# Hint: run the following command on a separate terminal to receive a detailed
# Redis command trace:
#
#     podman exec -ti test-api-server-redis redis-cli --user cluster --pass cluster-test-password MONITOR
#

set +e

# Start the test suite by joining the Redis network namespace. Script
# arguments are forwarded to the "robot" command. The Python virtual env
# is cached in the "rfenv-api-server" volume.
podman run -i -q --rm --replace --name test-api-server-suite \
    --network=container:test-api-server-redis \
    --requires=test-api-server-redis \
    --workdir="/usr/src" \
    --volume="rfenv-api-server:/usr/local/rfenv:z" \
    --volume="tstate-api-server:/srv:z" \
    --volume="${PWD}:/opt/src:z,ro" \
    docker.io/library/python:3.11 \
bash -l -s - "${@}" <<'EOF'
set -e
cp -a /opt/src/api-server /opt/src/test .
if [[ ! -d /usr/local/rfenv/bin ]]; then
    python3 -mvenv /usr/local/rfenv --upgrade
    /usr/local/rfenv/bin/python3 -m pip install --upgrade pip
fi
PATH="/usr/local/rfenv/bin:$PATH"
VIRTUAL_ENV=/usr/local/rfenv
export PATH VIRTUAL_ENV
pip3 install -q -r test/requirements.txt
mkdir -p /srv/log /srv/output
robot \
    --name "API server tests" \
    -d /srv/output \
    --console=dotted "${@}" \
    test/suite/
EOF
rc=$?

podman stop -i test-api-server-redis >/dev/null
exit ${rc}

#!/bin/bash

#
# Copyright (C) 2021 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# NethServer is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with NethServer.  If not, see COPYING.
#

#
# Run this script in this directory with
#
# bash test-agent.sh [ROBOT ARG]...
#

set -e

bcontainer="gobuilder-agent"

# Reuse existing ${bcontainer}, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep "${bcontainer}" >/dev/null; then
    buildah from --name "${bcontainer}" -v ${PWD}:/usr/src:Z docker.io/golang:1.18-alpine
fi

# Rebuild the "agent" Golang binary
buildah run "${bcontainer}" sh -c "cd /usr/src && CGO_ENABLED=0 go build -v ."

# If not already running, start a detached Redis server in a separate
# network namespace.
if ! podman container exists test-agent-redis; then
    podman run -d -q --rm --replace \
        --name test-agent-redis \
        --network=slirp4netns \
        docker.io/redis:6-alpine \
        --port 6379 --save 5 1 --protected-mode no --loglevel debug
fi

#
# Hint: run the following command on a separate terminal to receive a detailed
# Redis command trace:
#
#     podman exec -ti test-agent-redis nc 127.0.0.1 6379 <<<$'MONITOR\r\n'
#

set +e

# Start the test suite by joining the Redis network namespace. Script
# arguments are forwarded to the "robot" command! Python virtual env is
# cached in the "rfenv" volume, to speed-up subsequent runs.
exec podman run -i -q --rm --replace --name test-agent-suite \
    --network=container:test-agent-redis \
    --requires=test-agent-redis \
    --env=AGENT_ID=module/t1000 \
    --workdir="/usr/src" \
    --volume="rfenv:/usr/local/rfenv:z" \
    --volume="tstate:/srv:z" \
    --volume="${PWD}:/opt/src:z,ro" \
    docker.io/python:3.11 \
bash -l -s - "${@}" <<'EOF'
set -e
cp -a /opt/src/. .
if [[ ! -d /usr/local/rfenv/bin ]]; then
    python3 -mvenv /usr/local/rfenv --upgrade
    /usr/local/rfenv/bin/python3 -m pip install --upgrade pip
fi
PATH="/usr/local/rfenv/bin:$PATH"
VIRTUAL_ENV=/usr/local/rfenv
export PATH VIRTUAL_ENV
pip3 install -q -r test/requirements.txt
robot \
    --name "Agent tests" \
    -d /srv/output \
    --console=dotted "${@}" \
    --skiponfailure unstable \
    test/suite/
EOF

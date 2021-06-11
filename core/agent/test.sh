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

set -e

export AGENT_INSTALL_DIR=.

container="gobuilder-agent"
redis_image="docker.io/redis:6-alpine"

# Reuse existing ${container}, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q "${container}"; then
    buildah from --name "${container}" -v ${PWD}:/usr/src:Z golang:1.16-alpine
fi

buildah run "${container}" sh -c "cd /usr/src && CGO_ENABLED=0 go build -v ."

podman run --volume=redis-data:/data --rm --network=host --name redis --detach $redis_image --port 6380 --save 5 1 --bind 127.0.0.1 ::1 --protected-mode no
trap 'echo "Stopping processes..." ; kill $agentpid ; podman stop redis ; rm -f environment' EXIT

podman exec -i redis redis-cli -p 6380 <<EOF
HSET t/environment MYVARIABLE MYVALUE
LPUSH t/tasks '{"id":"l","action":"list-actions","data":""}'
LPUSH t/tasks '{"id":"t","action":"test","data":${1:-\"OK\"}}'
EOF

podman exec -i redis redis-cli -p 6380 monitor &
sleep 1

REDIS_ADDRESS=127.0.0.1:6380 ./agent t . &
agentpid=$!

sleep 3

printf "Environment file contents:\n%s\n" "$(cat environment)"

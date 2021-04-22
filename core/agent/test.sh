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

container="gobuilder-agent"

# Reuse existing ${container}, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q "${container}"; then
    buildah from --name "${container}" -v ${PWD}:/usr/src:Z golang:1.16-alpine
fi

buildah run "${container}" sh -c "cd /usr/src && CGO_ENABLED=0 go build -v ."

podman run --volume=redis-data:/data --replace --network=host --name redis --detach docker.io/redis:6-alpine --save 5 1 --bind 127.0.0.1 ::1 --protected-mode no
trap "podman stop redis && podman rm redis" EXIT

podman exec -i redis redis-cli <<EOF
HSET t/environment MYVARIABLE MYVALUE
LPUSH t/tasks '{"id":"t","action":"test","data":""}'
EOF

./agent t . &
agentpid=$!
sleep 1
kill $agentpid

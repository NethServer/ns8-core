#!/bin/bash

set -e
messages=("Publish the images with:" "")
repobase="ghcr.io/nethserver"

# Reuse existing gobuilder-core container, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q gobuilder-core; then
    echo "Pulling Golang runtime..."
    buildah from --name gobuilder-core -v ${PWD}:/usr/src/core:Z docker.io/library/golang:1.16-alpine
fi

echo "Build statically linked Go binaries based on Musl..."
echo "1/2 agent..."
buildah run gobuilder-core sh -c "cd /usr/src/core/agent      && CGO_ENABLED=0 go build -v ."
echo "2/2 api-server..."
buildah run gobuilder-core sh -c "cd /usr/src/core/api-server && CGO_ENABLED=0 go build -v ."

echo "Building the Core image..."
container=$(buildah from scratch)
reponame="core"
buildah add ${container} imageroot /
buildah add ${container} agent/agent /usr/local/bin/agent
buildah add ${container} api-server/api-server /usr/local/bin/api-server
buildah config --entrypoint=/ ${container}
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm ${container}
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

echo "Building the Redis image..."
container=$(buildah from docker.io/redis:6-alpine)
reponame="redis"
# Reset upstream volume configuration: it is necessary to modify /data contents with our .conf file.
buildah config --volume=/data- "${container}"
buildah run "${container}" sh <<EOF
mkdir etc

cat >etc/redis.acl <<EOR
user default on nopass sanitize-payload ~* &* +@all
EOR

cat >etc/redis.conf <<EOR
#
# redis.conf -- do not edit manually!
#
bind 127.0.0.1 ::1
port 6379
protected-mode no
save 5 1
aclfile "/data/etc/redis.acl"
dir "/data"
masteruser default
masterauth nopass
EOR

EOF
buildah config --volume=/data '--cmd=[ "redis-server", "/data/etc/redis.conf" ]' ${container}
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm ${container}
messages+=(" buildah push ${repobase}/${reponame} docker://${repobase}/${reponame}:latest")

printf "%s\n" "${messages[@]}"

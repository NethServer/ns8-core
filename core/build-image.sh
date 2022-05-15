#!/bin/bash

set -e
cleanup_list=()
trap 'rm -rf "${cleanup_list[@]}"' EXIT
images=()
repobase="${REPOBASE:-ghcr.io/nethserver}"

# Reuse existing gobuilder-core container, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q gobuilder-core; then
    echo "Pulling Golang runtime..."
    golang_cache_path="${PWD}/.golang-cache"
    mkdir -vp "${golang_cache_path}/{mcache,bcache}"
    buildah from --name gobuilder-tmp docker.io/library/golang:1.18-alpine
    buildah run gobuilder-tmp apk add g++ gcc
    buildah config --env GOCACHE=/var/lib/misc/bcache --env GOMODCACHE=/var/lib/misc/mcache gobuilder-tmp
    buildah commit --rm gobuilder-tmp gobuilder-image
    buildah from --name gobuilder-core \
        -v "${golang_cache_path}:/var/lib/misc:z" \
        -v "${PWD}:/usr/src/core:z" \
        localhost/gobuilder-image
fi

# Reuse existing nodebuilder-core container, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q nodebuilder-core; then
    echo "Pulling NodeJS runtime..."
    buildah from --name nodebuilder-core -v "${PWD}:/usr/src/core:z" docker.io/library/node:lts-slim
fi

echo "Build statically linked Go binaries based on Musl..."
echo "1/2 agent..."
buildah run gobuilder-core sh -c "cd /usr/src/core/agent && CGO_ENABLED=0 go build -v ."

echo "2/2 api-server and api-server-logs..."
# Statically link libraries and disable Sqlite extensions that expect a dynamic loader (not portable across distros)
# Ref https://www.arp242.net/static-go.html
buildah run gobuilder-core sh -c "cd /usr/src/core/api-server && go build -v -ldflags=\"-extldflags=-static\" -tags sqlite_omit_load_extension api-server.go"
buildah run gobuilder-core sh -c "cd /usr/src/core/api-server && go build -v -ldflags=\"-extldflags=-static\" -tags sqlite_omit_load_extension api-server-logs.go"

echo "Build static UI files with node..."
buildah run nodebuilder-core sh -c "cd /usr/src/core/ui && yarn install && yarn build"

echo "Install tidy..."
buildah run nodebuilder-core sh -c "apt-get update && apt-get install -y tidy"

echo "Provide core style to external modules..."
buildah run nodebuilder-core sh -c "cd /usr/src/core/ui/dist/css/ && tidy ../index.html | grep -oP 'link href=\"css/app~.+\.css\" rel=\"stylesheet\"' | grep -oP 'app~.+\.css' | paste -sd ' ' - > css_files"
buildah run nodebuilder-core sh -c 'cd /usr/src/core/ui/dist/css/ && cat $(cat /usr/src/core/ui/dist/css/css_files) > /usr/src/core/ui/dist/css/core.css'

echo "Download Logcli..."
logcli_tmp_dir=$(mktemp -d)
cleanup_list+=("${logcli_tmp_dir}")
(
    cd "${logcli_tmp_dir}"
    curl -L -O https://github.com/grafana/loki/releases/download/v2.2.1/logcli-linux-amd64.zip
    python -mzipfile -e logcli-linux-amd64.zip .
    chmod -c 755 logcli-linux-amd64
)

echo "Building the Core image..."
container=$(buildah from scratch)
reponame="core"
buildah add "${container}" imageroot /
buildah add "${container}" "${logcli_tmp_dir}/logcli-linux-amd64" /usr/local/bin/logcli.bin
buildah add "${container}" agent/agent /usr/local/bin/agent
buildah add "${container}" api-server/api-server /usr/local/bin/api-server
buildah add "${container}" api-server/api-server-logs /usr/local/bin/api-server-logs
buildah add "${container}" ui/dist /var/lib/nethserver/cluster/ui
buildah add "${container}" install.sh /var/lib/nethserver/node/install.sh
core_env_file=$(mktemp)
cleanup_list+=("${core_env_file}")
printf "CORE_IMAGE=${repobase}/core:%s\n" "${IMAGETAG:-latest}" >> "${core_env_file}"
printf "REDIS_IMAGE=${repobase}/redis:%s\n" "${IMAGETAG:-latest}" >> "${core_env_file}"
printf "RCLONE_IMAGE=docker.io/rclone/rclone:1.57.0\n" >> "${core_env_file}"
printf "RSYNC_IMAGE=${repobase}/rsync:%s\n" "${IMAGETAG:-latest}" >> "${core_env_file}"
printf "RESTIC_IMAGE=${repobase}/restic:%s\n" "${IMAGETAG:-latest}" >> "${core_env_file}"
chmod -c 644 "${core_env_file}"
source "${core_env_file}"
buildah add "${container}" ${core_env_file} /etc/nethserver/core.env
buildah config \
    --label="org.nethserver.images=${REDIS_IMAGE} ${RCLONE_IMAGE}" \
    --label="org.nethserver.flags=core_module" \
    --entrypoint=/ "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm "${container}"
images+=("${repobase}/${reponame}")

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
port 6379
supervised no
protected-mode no
save 5 1
aclfile "/data/etc/redis.acl"
dir "/data"
masteruser default
masterauth nopass
notify-keyspace-events AKE
EOR

EOF
buildah config --volume=/data '--cmd=[ "redis-server", "/data/etc/redis.conf" ]' "${container}"
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm "${container}"
images+=("${repobase}/${reponame}")

echo "Building the restic image..."
container=$(buildah from alpine)
reponame="restic"
buildah run ${container} -- apk add --no-cache restic
buildah config --cmd [] ${container}
buildah config --entrypoint '["/usr/bin/restic"]' ${container}
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm "${container}"
images+=("${repobase}/${reponame}")

echo "Building the rsync image..."
container=$(buildah from alpine)
reponame="rsync"
buildah run ${container} -- apk add --no-cache rsync
buildah add "${container}" rsync/entrypoint.sh /entrypoint.sh
buildah config \
    --entrypoint='["/entrypoint.sh"]' \
    --cmd='["rsync", "--daemon", "--no-detach"]' \
    ${container}
buildah commit "${container}" "${repobase}/${reponame}"
buildah rm "${container}"
images+=("${repobase}/${reponame}")

if [[ -n "${CI}" ]]; then
    # Set output value for Github Actions
    printf "::set-output name=images::%s\n" "${images[*]}"
else
    printf "Publish the images with:\n\n"
    for image in "${images[@]}"; do printf "  buildah push %s docker://%s:latest\n" "${image}" "${image}" ; done
    printf "\n"
fi

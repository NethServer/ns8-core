#!/bin/bash

# Terminate on error
set -e

# Prepare variables for later user
images=()
# The image willbe pushed to GitHub image registry under nethserver organization
repobase="ghcr.io/nethserver"
# Configure the image name
reponame="mariadb"

# Create a new empty container image
container=$(buildah from scratch)

# Reuse existing nodebuilder-mariadb container, to speed up builds
if ! buildah containers --format "{{.ContainerName}}" | grep -q nodebuilder-mariadb; then
    echo "Pulling NodeJS runtime..."
    buildah from --name nodebuilder-mariadb -v "${PWD}:/usr/src/mariadb:Z" docker.io/library/node:lts
fi

echo "Build static UI files with node..."
buildah run nodebuilder-mariadb sh -c "cd /usr/src/mariadb/ui       && npm install && npm run build"

# Add imageroot directory to the container image
buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui/dist /ui
# Setup the entrypoint, ask to reserve one TCP port with the label and set a rootless container
buildah config --entrypoint=/ \
    --label="org.nethserver.tcp-ports-demand=1" \
    --label="org.nethserver.rootfull=0" \
    --label="org.nethserver.authorizations=traefik@any:routeadm" \
    --label="org.nethserver.images=docker.io/mariadb:10.6.5 docker.io/phpmyadmin:5.1.1" \
    "${container}"
# Commit everything
buildah commit "${container}" "${repobase}/${reponame}"

images+=("${repobase}/${reponame}")

# Setup CI when pushing to Github
if [[ -n "${CI}" ]]; then
    # Set output value for Github Actions
    printf "::set-output name=images::%s\n" "${images[*]}"
else
    # Just print info for manual push
    printf "Publish the images with:\n\n"
    for image in "${images[@]}"; do printf "  buildah push %s docker://%s:0.0.1\n" "${image}" "${image}" ; done
    printf "\n"
fi

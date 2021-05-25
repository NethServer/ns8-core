#!/bin/bash

# Terminate on error
set -e

# Prepare variables for later user
images=()
# The image willbe pushed to GitHub image registry under nethserver organization
repobase="ghcr.io/nethserver"
# Configure the image name
reponame="dokuwiki"

# Create a new empty container image
container=$(buildah from scratch)

# Add imageroot directory to the container image
buildah add "${container}" imageroot /imageroot
buildah add "${container}" ui /ui
# Setup the entrypoint, ask to reserve one TCP port with the label and set a rootless container
buildah config --entrypoint=/ --label="org.nethserver/tcp_ports_demand=1" --label="org.nethserver/rootfull=0" "${container}"
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
    for image in "${images[@]}"; do printf "  buildah push %s docker://%s:latest\n" "${image}" "${image}" ; done
    printf "\n"
fi

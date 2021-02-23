#!/bin/bash

container=$(buildah from scratch)

buildah copy ${container} agent /agent
buildah copy ${container} cplane /cplane

echo "Manually download the DigitalOcean docker-config.json file, then publish the image with:"
echo
echo " buildah commit --authfile docker-config.json ${container} docker://registry.digitalocean.com/nethserver/control-plane:latest"

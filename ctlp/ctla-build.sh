#!/bin/bash

container=$(buildah from docker.io/library/python:3.9.1-alpine3.13)

buildah run ${container} /bin/ash <<EOF
python3 -mvenv /opt/ctla
/opt/ctla/bin/pip3 install redis
mkdir -p /opt/ctla/share/events
mkdir -p /var/opt/ctla/events
mkdir -p /etc/opt/ctla/systemd/system
EOF

buildah copy ${container} ctla.py /opt/ctla/bin/ctla

buildah run ${container} /bin/ash <<EOF
chmod +x /opt/ctla/bin/ctla
EOF

buildah config --shell /bin/ash --entrypoint '[ "/opt/ctla/bin/ctla" ]' ${container}
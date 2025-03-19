#!/bin/bash

LEADER_NODE=$1
SSH_KEYFILE=${SSH_KEYFILE:-$HOME/.ssh/id_rsa}
COREMODULES=$(echo ${COREMODULES} | tr ' ' ',')

ssh_key="$(cat $SSH_KEYFILE)"

wget -nv -P /tmp/ https://raw.githubusercontent.com/microsoft/playwright/master/utils/docker/seccomp_profile.json

podman run -i \
    -v ..:/home/pwuser/ns8-scratchpad:z \
    --security-opt seccomp=/tmp/seccomp_profile.json \
    --ipc=host \
    --volume=site-packages:/home/pwuser/.local/lib/python3.8/site-packages:Z \
    --name rf-core-runner ghcr.io/marketsquare/robotframework-browser/rfbrowser-stable:19.4.0 \
    bash -l -s <<EOF
    set -e
    echo "$ssh_key" > /home/pwuser/ns8-key
    set -x
    pip install -r /home/pwuser/ns8-scratchpad/core/tests/pythonreq.txt
    mkdir ~/outputs
    cd /home/pwuser/ns8-scratchpad/core/
    robot -v NODE_ADDR:${LEADER_NODE} \
        -v SSH_KEYFILE:/home/pwuser/ns8-key \
	-v COREMODULES:${COREMODULES} \
	--name core \
	--skiponfailure unstable \
    --exclude frontend \
	--console dotted \
	-d ~/outputs /home/pwuser/ns8-scratchpad/core/tests/
EOF

tests_res=$?

podman cp rf-core-runner:/home/pwuser/outputs tests/
podman stop rf-core-runner
podman rm rf-core-runner
rm -f /tmp/seccomp_profile.json

exit ${tests_res}

#!/bin/bash

LEADER_NODE=$1
SSH_KEYFILE=${SSH_KEYFILE:-$HOME/.ssh/id_rsa}

ssh_key="$(cat $SSH_KEYFILE)"

podman run -i \
    -v ..:/home/pwuser/ns8-scratchpad:z \
    --name rf-core-runner ghcr.io/marketsquare/robotframework-browser/rfbrowser-stable:v10.0.3 \
    bash -l -s <<EOF
    set -e -x
    echo "$ssh_key" > /home/pwuser/ns8-key
    pip install -r /home/pwuser/ns8-scratchpad/core/tests/pythonreq.txt
    mkdir ~/outputs
    cd /home/pwuser/ns8-scratchpad/core/
    robot -v NODE_ADDR:${LEADER_NODE} \
        -v SSH_KEYFILE:/home/pwuser/ns8-key \
	-v COREBRANCH:${COREBRANCH} \
	-v COREMODULES:${COREMODULES} \
	-d ~/outputs /home/pwuser/ns8-scratchpad/core/tests/
EOF

podman cp rf-core-runner:/home/pwuser/outputs tests/
podman stop rf-core-runner
podman rm rf-core-runner

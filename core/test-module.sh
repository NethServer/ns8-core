#!/bin/bash

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

set -e -a

# ////
_script_start=$(date +%s)

# Resolve the directory containing this script so volume mounts are correct
# regardless of the working directory from which the script is invoked.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

SSH_KEYFILE=${SSH_KEYFILE:-$HOME/.ssh/id_rsa}
LEADER_NODE="${1:?missing LEADER_NODE argument}"
shift

ssh_key="$(< $SSH_KEYFILE)"
venvroot=/usr/local/venv

echo "Test! RUN_UI_TESTS=${RUN_UI_TESTS} ////"

# Select the container image and Python requirements file based on whether UI tests are enabled.
# UI tests require the Playwright image (Debian-based, includes browser binaries).
# Non-UI tests use a lightweight Alpine Python image.
if [ "${RUN_UI_TESTS}" = "true" ]; then
    container_image="mcr.microsoft.com/playwright:v1.51.0-noble"
    container_shell="bash"
    pythonreq="/srv/source/tests/pythonreq-ui.txt"
    cache_volume="rftest-cache-ui"
else
    container_image="docker.io/python:3.11-alpine"
    container_shell="ash"
    pythonreq="/srv/source/tests/pythonreq.txt"
    cache_volume="rftest-cache"
fi

podman run -i \
    --network=host \
    --volume="${SCRIPT_DIR}":/srv/source:z \
    --volume=${cache_volume}:${venvroot}:z \
    --replace --name=rftest \
    --env=ssh_key \
    --env=venvroot \
    --env=LEADER_NODE \
    --env=COREMODULES \
    --env=RUN_UI_TESTS \
    --env=pythonreq \
    --env=_script_start \
    "${container_image}" \
    ${container_shell} -l -s -- "${@}" <<'EOF'
set -e

echo "$ssh_key" > /tmp/idssh

# Install the Python venv and Robot Framework dependencies if not already cached.
# Cache is invalidated when the MD5 checksum of the requirements file changes, ensuring
# that dependency updates are always picked up even on reused (self-hosted) runners.
pythonreq_checksum_file="${venvroot}/.pythonreq.md5"
pythonreq_current_checksum=$(md5sum "${pythonreq}" | cut -d' ' -f1)
pythonreq_cached_checksum=$(cat "${pythonreq_checksum_file}" 2>/dev/null || true)

if [ ! -x "${venvroot}/bin/robot" ] || [ "${pythonreq_current_checksum}" != "${pythonreq_cached_checksum}" ] ; then
    if command -v apt-get > /dev/null 2>&1; then
        # mcr.microsoft.com/playwright:*-noble has npm pre-installed but no Python
        apt-get update -q
        apt-get install -y -q python3 python3-venv
        python3 -mvenv "${venvroot}"
    else
        # Alpine image already has Python; --upgrade refreshes pip/setuptools in-place
        python3 -mvenv "${venvroot}" --upgrade
    fi
    ${venvroot}/bin/pip3 install -q -r "${pythonreq}"
    # Save the checksum so future runs can detect requirement changes
    echo "${pythonreq_current_checksum}" > "${pythonreq_checksum_file}"
    # Invalidate the rfbrowser sentinel so it is re-initialized with the new packages
    rm -f "${venvroot}/.rfbrowser_initialized"
fi

# Initialize the Playwright browser binaries for rfbrowser (UI tests only).
# A sentinel file prevents re-running this expensive step on cache hits.
if [ "${RUN_UI_TESTS}" = "true" ] && [ ! -f "${venvroot}/.rfbrowser_initialized" ] ; then
    ${venvroot}/bin/rfbrowser init
    touch "${venvroot}/.rfbrowser_initialized"
fi

cd /srv/source
mkdir -vp tests/outputs/

# Exclude UI tests if RUN_UI_TESTS is not set to "true"
if [ "${RUN_UI_TESTS}" = "true" ]; then
    ui_tag_filter=""
else
    ui_tag_filter="--exclude ui"
fi

echo "DEBUG[core]: $(( $(date +%s) - _script_start ))s elapsed from script start to robot launch ////"

exec ${venvroot}/bin/robot \
    -v "NODE_ADDR:${LEADER_NODE}" \
    ${COREMODULES:+-v "COREMODULES:${COREMODULES}"} \
    -v SSH_KEYFILE:/tmp/idssh \
    -v "RUN_UI_TESTS:${RUN_UI_TESTS}" \
    --name core \
    --skiponfailure unstable \
    ${ui_tag_filter} \
    -d tests/outputs "${@}" tests/
EOF

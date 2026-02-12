#!/bin/bash

# NOTE: bash is required for correct wait -n semantics

#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#


set -eu

HAPROXY_CFG=/usr/local/haproxy/haproxy.cfg

RCLONE_WEBDAV_SOCK=/var/lib/rclone/backend-webdav.sock
RCLONE_REST_SOCK=/var/lib/rclone/backend-rest.sock

log() {
    echo "[entrypoint] $*"
}

# ---------------------------
# HAProxy functions
# ---------------------------
start_haproxy() {
    log "Starting haproxy frontend"
    haproxy -W -db -f "$HAPROXY_CFG" &
    HAPROXY_MASTER_PID=$!
    log "haproxy master pid=$HAPROXY_MASTER_PID"
}

# ---------------------------
# Rclone functions
# ---------------------------
start_rclone() {
    if ! grep -F '[combined]' "${RCLONE_CONFIG}" ; then
        log "Cannot start Rclone backends: combined remote is empty."
        return
    fi
    log "Starting rclone webdav backend"
    rm -f "$RCLONE_WEBDAV_SOCK"
    rclone serve webdav combined: \
        --addr "unix://$RCLONE_WEBDAV_SOCK" \
        --vfs-cache-mode writes &
    RCLONE_WEBDAV_PID=$!
    log "webdav pid=$RCLONE_WEBDAV_PID"

    log "Starting rclone restic backend"
    rm -f "$RCLONE_REST_SOCK"
    rclone serve restic combined: \
        --addr "unix://$RCLONE_REST_SOCK" &
    RCLONE_REST_PID=$!
    log "restic pid=$RCLONE_REST_PID"
}

# ---------------------------
# Signal handling
# ---------------------------
PROCESS_RELOAD=0
PROCESS_SHUTDOWN=0
reload_handler()   { PROCESS_RELOAD=1 ; }
shutdown_handler() { PROCESS_SHUTDOWN=1 ; }
trap reload_handler   HUP USR1 USR2
trap shutdown_handler TERM INT

# ---------------------------
# Start services
# ---------------------------
start_rclone
start_haproxy

log "Container ready"

# ----------------------------------------
# Keep container alive and process signals
# ----------------------------------------
wait -n || exit_code=$?
while [ ${PROCESS_RELOAD} == 1 ] ; do
    log "Reloading haproxy"
    kill -USR2 "$HAPROXY_MASTER_PID" || {
        log "haproxy master not responding, starting fresh"
        start_haproxy
    }
    log "Reloading rclone backends"
    kill "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    kill "$RCLONE_REST_PID" 2>/dev/null || true
    wait "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    wait "$RCLONE_REST_PID" 2>/dev/null || true
    start_rclone
    PROCESS_RELOAD=0
    wait -n || exit_code=$?
done

if [ ${PROCESS_SHUTDOWN} == 1 ]; then
    exit_code=0
    log "Shutting down"
    kill "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    kill "$RCLONE_REST_PID" 2>/dev/null || true
    kill "$HAPROXY_MASTER_PID" 2>/dev/null || true
    wait || true
fi

exit ${exit_code}

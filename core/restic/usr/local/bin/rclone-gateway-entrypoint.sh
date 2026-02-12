#!/bin/sh
set -eu

HAPROXY_CFG=/etc/haproxy/haproxy.cfg
HAPROXY_PIDFILE=/run/haproxy.pid

RCLONE_WEBDAV_SOCK=/run/rclone-webdav.sock
RCLONE_REST_SOCK=/run/rclone-rest.sock

RCLONE_CFG=/etc/rclone/rclone.conf

log() {
    echo "[entrypoint] $*"
}

# ---------------------------
# HAProxy functions
# ---------------------------
start_haproxy() {
    log "Starting haproxy"
    haproxy -W -db -f "$HAPROXY_CFG" -p "$HAPROXY_PIDFILE" &
    HAPROXY_MASTER_PID=$!
    log "haproxy master pid=$HAPROXY_MASTER_PID"
}

# ---------------------------
# Rclone functions
# ---------------------------
start_rclone() {
    log "Starting rclone webdav"
    rm -f "$RCLONE_WEBDAV_SOCK"
    rclone serve webdav combined: \
        --config "$RCLONE_CFG" \
        --addr "unix://$RCLONE_WEBDAV_SOCK" &
    RCLONE_WEBDAV_PID=$!
    log "webdav pid=$RCLONE_WEBDAV_PID"

    log "Starting rclone restic"
    rm -f "$RCLONE_REST_SOCK"
    rclone serve restic combined: \
        --config "$RCLONE_CFG" \
        --addr "unix://$RCLONE_REST_SOCK" &
    RCLONE_REST_PID=$!
    log "restic pid=$RCLONE_REST_PID"
}

# ---------------------------
# Reload
# ---------------------------
reload() {
    log "Reloading haproxy"
    if [ -f "$HAPROXY_PIDFILE" ]; then
        OLD_PID=$(cat "$HAPROXY_PIDFILE")
        haproxy -W -db -f "$HAPROXY_CFG" -p "$HAPROXY_PIDFILE" -sf "$OLD_PID" &
        HAPROXY_MASTER_PID=$!
        log "haproxy reloaded new master pid=$HAPROXY_MASTER_PID"
    else
        log "pidfile missing, starting fresh haproxy"
        start_haproxy
    fi

    log "Reloading rclone"
    kill "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    kill "$RCLONE_REST_PID" 2>/dev/null || true
    wait "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    wait "$RCLONE_REST_PID" 2>/dev/null || true
    start_rclone
}


# ---------------------------
# Shutdown
# ---------------------------
shutdown() {
    log "Shutting down"
    kill "$RCLONE_WEBDAV_PID" 2>/dev/null || true
    kill "$RCLONE_REST_PID" 2>/dev/null || true
    if [ -f "$HAPROXY_PIDFILE" ]; then
        kill "$(cat "$HAPROXY_PIDFILE")" 2>/dev/null || true
    fi
    wait || true
    exit 0
}

# ---------------------------
# Signal handling
# ---------------------------
trap reload HUP USR1 USR2
trap shutdown TERM INT

# ---------------------------
# Start services
# ---------------------------
start_rclone
start_haproxy

log "Container ready"

# ------------------------------
# Keep container alive for years
# ------------------------------
sleep 1000000000

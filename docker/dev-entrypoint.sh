#!/usr/bin/env bash

set -euo pipefail

TARGET_USER="${DEV_CONTAINER_USER:-agent}"
DOCKER_SOCKET="${DOCKER_SOCKET:-/var/run/docker.sock}"

if [ "$(id -u)" -eq 0 ] && [ -S "$DOCKER_SOCKET" ]; then
    # Match the mounted socket's group ID so the non-root dev user can use it.
    socket_gid="$(stat -c '%g' "$DOCKER_SOCKET")"

    if [ "$socket_gid" != "0" ]; then
        socket_group="$(getent group "$socket_gid" | cut -d: -f1 || true)"

        if [ -z "$socket_group" ]; then
            socket_group="host-docker"
            groupadd --gid "$socket_gid" "$socket_group" 2>/dev/null || true
            socket_group="$(getent group "$socket_gid" | cut -d: -f1 || true)"
        fi

        if [ -n "$socket_group" ] && ! id -nG "$TARGET_USER" | tr ' ' '\n' | grep -qx "$socket_group"; then
            usermod -aG "$socket_group" "$TARGET_USER"
        fi
    fi

    exec runuser -u "$TARGET_USER" -- "$@"
fi

exec "$@"

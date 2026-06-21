#!/bin/sh
set -e

if [ "$(stat -c '%u' /app/data)" != "65534" ]; then
    chown -R 65534:65534 /app/data
fi

exec su-exec 65534:65534 /app/server
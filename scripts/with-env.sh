#!/bin/sh
set -e
if [ -f .env ]; then
    set -a
    . ./.env
    set +a
elif [ -z "${DATABASE_URL:-}" ]; then
    printf '%s\n' 'Missing database configuration: copy .env.example to .env and configure your database first.' >&2
    exit 1
fi
exec "$@"

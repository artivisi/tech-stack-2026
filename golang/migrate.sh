#!/usr/bin/env bash
set -eu
cd "$(dirname "$0")"

if [ ! -f ../.env ]; then
    echo "ERROR: ../.env not found" >&2
    exit 1
fi

set -a
. ../.env
set +a

exec go run ./cmd/migrate "$@"

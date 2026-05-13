#!/usr/bin/env bash
# Source ../.env then invoke mvnw. Mirrors the expressjs shell-sourcing pattern.

set -eu
cd "$(dirname "$0")"

if [ ! -f ../.env ]; then
    echo "ERROR: ../.env not found" >&2
    exit 1
fi

set -a
. ../.env
set +a

exec ./mvnw spring-boot:run "$@"

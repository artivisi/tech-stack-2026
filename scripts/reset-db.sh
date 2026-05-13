#!/usr/bin/env bash
# Wipe and recreate the postgres container + volume.
# Use between verification runs (especially across stacks) for a clean slate.

set -eu
cd "$(dirname "$0")/.."

echo "tearing down..."
docker compose down -v

echo "starting fresh db..."
docker compose up -d db

echo "waiting for healthy..."
until docker inspect --format '{{.State.Health.Status}}' tech-stack-2026-db 2>/dev/null | grep -q healthy; do
  sleep 1
done

echo "db ready"

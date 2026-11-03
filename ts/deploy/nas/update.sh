#!/bin/sh
# Rebuild and restart ojeommu from the newest main. Same shape as the other NAS update.sh scripts,
# so a webhook hook can call it later; until then run it by hand.
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/../../.." && pwd)
cd "$root"
[ "$(git branch --show-current)" = main ] || { echo 'main checkout에서만 갱신합니다.' >&2; exit 1; }
git pull --ff-only origin main
compose() { docker compose -f "$root/ts/deploy/nas/docker-compose.yml" "$@"; }
# set -e keeps the running container when the build fails.
compose build
compose up -d

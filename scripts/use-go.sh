#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_ROOT="${GUTENBERG_GOROOT:-${BLACK_FORGE_GOROOT:-/tmp/black-forge-go/go}}"
DEFAULT_GOCACHE="/tmp/gutenberg-go-cache"
DEFAULT_GOMODCACHE="/tmp/gutenberg-go-modcache"

if [ ! -d "$DEFAULT_GOCACHE" ] && [ -d /tmp/black-forge-go-cache ]; then
  DEFAULT_GOCACHE="/tmp/black-forge-go-cache"
fi

if [ ! -f "$DEFAULT_GOMODCACHE/cache/download/modernc.org/sqlite/@v/v1.39.1.zip" ] && [ -f /tmp/black-forge-go-modcache/cache/download/modernc.org/sqlite/@v/v1.39.1.zip ]; then
  DEFAULT_GOMODCACHE="/tmp/black-forge-go-modcache"
fi

if [ ! -x "$GO_ROOT/bin/go" ]; then
  echo "Go is not installed at $GO_ROOT/bin/go" >&2
  echo "Download and extract the official Go archive, then rerun this script." >&2
  exit 1
fi

export GOROOT="$GO_ROOT"
export GOCACHE="${GUTENBERG_GOCACHE:-${BLACK_FORGE_GOCACHE:-$DEFAULT_GOCACHE}}"
export GOMODCACHE="${GUTENBERG_GOMODCACHE:-${BLACK_FORGE_GOMODCACHE:-$DEFAULT_GOMODCACHE}}"
export GOFLAGS="${GOFLAGS:--buildvcs=false}"
export PATH="$GO_ROOT/bin:$ROOT/.tools/bin:$PATH"
mkdir -p "$GOCACHE" "$GOMODCACHE"
go "$@"

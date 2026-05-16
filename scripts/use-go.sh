#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEFAULT_GOCACHE="/tmp/gutenberg-go-cache"
DEFAULT_GOMODCACHE="/tmp/gutenberg-go-modcache"
CONFIGURED_GO_ROOT="${GUTENBERG_GOROOT:-${BLACK_FORGE_GOROOT:-}}"
FALLBACK_GO_ROOT="/tmp/black-forge-go/go"

if [ ! -d "$DEFAULT_GOCACHE" ] && [ -d /tmp/black-forge-go-cache ]; then
  DEFAULT_GOCACHE="/tmp/black-forge-go-cache"
fi

if [ ! -f "$DEFAULT_GOMODCACHE/cache/download/modernc.org/sqlite/@v/v1.39.1.zip" ] && [ -f /tmp/black-forge-go-modcache/cache/download/modernc.org/sqlite/@v/v1.39.1.zip ]; then
  DEFAULT_GOMODCACHE="/tmp/black-forge-go-modcache"
fi

if [ -n "$CONFIGURED_GO_ROOT" ]; then
  if [ ! -x "$CONFIGURED_GO_ROOT/bin/go" ]; then
    echo "Go is not installed at $CONFIGURED_GO_ROOT/bin/go" >&2
    echo "Set GUTENBERG_GOROOT to a Go installation root, or install Go on PATH." >&2
    exit 1
  fi
  GO_ROOT="$CONFIGURED_GO_ROOT"
  GO_BIN="$GO_ROOT/bin/go"
elif command -v go >/dev/null 2>&1; then
  GO_BIN="$(command -v go)"
  GO_ROOT="$(cd "$(dirname "$GO_BIN")/.." && pwd)"
elif [ -x "$FALLBACK_GO_ROOT/bin/go" ]; then
  GO_ROOT="$FALLBACK_GO_ROOT"
  GO_BIN="$GO_ROOT/bin/go"
else
  echo "Go is not installed." >&2
  echo "Install Go on PATH, or set GUTENBERG_GOROOT to a Go installation root." >&2
  exit 1
fi

export GOROOT="$GO_ROOT"
export GOCACHE="${GUTENBERG_GOCACHE:-${BLACK_FORGE_GOCACHE:-$DEFAULT_GOCACHE}}"
export GOMODCACHE="${GUTENBERG_GOMODCACHE:-${BLACK_FORGE_GOMODCACHE:-$DEFAULT_GOMODCACHE}}"
export GOFLAGS="${GOFLAGS:--buildvcs=false}"
export PATH="$GO_ROOT/bin:$ROOT/.tools/bin:$PATH"
mkdir -p "$GOCACHE" "$GOMODCACHE"
"$GO_BIN" "$@"

#!/usr/bin/env bash
# Compares two connector manifests (kei.connector-manifest/v1) and fails on
# any breaking change of proposal 0002, section 6. Prints one line per
# breaking change:
#
#   breaking: <op id>: <reason>
#
# Usage: scripts/manifest-diff.sh <old manifest.json> <new manifest.json>
# Exit status: 0 no breaking change, 1 breaking changes, 2 usage or input
# error. A rename passes only with the old id kept, marked deprecated, and
# replaced_by naming the new id.
set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "usage: scripts/manifest-diff.sh <old manifest.json> <new manifest.json>" >&2
  exit 2
fi

root="$(cd "$(dirname "$0")/.." && pwd)"
abs() {
  if [ ! -f "$1" ]; then
    echo "manifest-diff: $1: no such file" >&2
    exit 2
  fi
  echo "$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
}
old="$(abs "$1")"
new="$(abs "$2")"

# Build, then exec: go run would turn every non-zero exit into 1.
bin="$(mktemp -d)"
trap 'rm -rf "$bin"' EXIT
go -C "$root/tools/protoc-gen-kei" build -o "$bin/manifest-diff" ./cmd/manifest-diff
"$bin/manifest-diff" "$old" "$new"

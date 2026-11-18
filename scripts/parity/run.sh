#!/usr/bin/env bash
# Runs one language's parity runner for gen/go/conformance/parity_test.go.
#
# Usage: scripts/parity/run.sh <go|ts|python|rust> <cases.json> <base url> <token>
#
# A runner reads the cases (one generated operation call each: op, connector,
# method, request_type, response_type, request), calls every operation through
# that language's generated code with an HTTP bearer connection whose base URL
# is "<base url>/<op>", and prints one JSON line per case:
#
#   {"op": "<op id>", "response": <protobuf JSON of the result>}
#   {"op": "<op id>", "error": "<message>"}
#
# When a runner is not available it prints "SKIP: <reason>" and exits 3.
#
# Adding a language: add a branch here and a runner next to this script, then
# add the language to requiredParityLanguages in
# gen/go/conformance/parity_test.go.
set -euo pipefail

root="$(cd "$(dirname "$0")/../.." && pwd)"

if [ "$#" -ne 4 ]; then
  echo "usage: scripts/parity/run.sh <go|ts|python|rust> <cases.json> <base url> <token>" >&2
  exit 2
fi
lang=$1
if [ ! -f "$2" ]; then
  echo "parity: $2: no such file" >&2
  exit 2
fi
# Runners may change directory (go -C), so pass the cases as an absolute path.
cases="$(cd "$(dirname "$2")" && pwd)/$(basename "$2")"
set -- "$cases" "$3" "$4"

skip() {
  echo "SKIP: $*"
  exit 3
}

case "$lang" in
  go)
    exec go -C "$root/gen/go" run ./conformance/parityrunner "$@"
    ;;
  ts)
    command -v node >/dev/null || skip "node is not installed"
    [ -d "$root/gen/ts/node_modules" ] || skip "gen/ts dependencies are not installed (npm --prefix gen/ts ci)"
    [ -f "$root/gen/ts/dist/kei/connection/index.js" ] || skip "gen/ts is not built (npm --prefix gen/ts run build)"
    exec node "$root/scripts/parity/runner.mjs" "$@"
    ;;
  python)
    # Task 7 adds scripts/parity/runner.py and the Python generated operations.
    [ -f "$root/scripts/parity/runner.py" ] || skip "no Python runner yet (scripts/parity/runner.py)"
    exec env PYTHONPATH="$root/gen/python/src" python3 -I "$root/scripts/parity/runner.py" "$@"
    ;;
  rust)
    # Task 9 adds scripts/parity/rust (a cargo binary) and the Rust operations.
    [ -f "$root/scripts/parity/rust/Cargo.toml" ] || skip "no Rust runner yet (scripts/parity/rust)"
    exec cargo run --quiet --manifest-path "$root/scripts/parity/rust/Cargo.toml" -- "$@"
    ;;
  *)
    echo "unknown language: $lang" >&2
    exit 2
    ;;
esac

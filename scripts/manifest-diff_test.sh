#!/usr/bin/env bash
# Tests scripts/manifest-diff.sh against the fixtures in
# scripts/testdata/manifest-diff: removing an operation fails with a
# "breaking:" line; renaming it (new id, old id deprecated with replaced_by)
# passes; an unchanged manifest passes; bad arguments fail.
set -uo pipefail

cd "$(dirname "$0")/.."
data=scripts/testdata/manifest-diff
failures=0

expect() { # expect <want exit> <want output substring> <args...>
  local want_code=$1 want_out=$2
  shift 2
  local out code
  out=$(bash scripts/manifest-diff.sh "$@" 2>&1)
  code=$?
  if [ "$code" -ne "$want_code" ] || ! grep -qF -- "$want_out" <<<"$out"; then
    echo "FAIL: manifest-diff.sh $*: exit $code, want $want_code with '$want_out'; output:" >&2
    echo "$out" >&2
    failures=$((failures + 1))
  else
    echo "ok: manifest-diff.sh $* -> exit $code"
  fi
}

expect 1 "breaking: acme.things.get: removed" "$data/old.json" "$data/removed.json"
expect 0 "OK: no breaking changes" "$data/old.json" "$data/renamed.json"
expect 0 "OK: no breaking changes" "$data/old.json" "$data/old.json"
expect 2 "usage" "$data/old.json"
expect 2 "no such file" "$data/old.json" "$data/missing.json"
expect 2 "manifest-diff:" "$data/old.json" "$data/manifest-diff_test.sh-not-json.txt"

if [ "$failures" -gt 0 ]; then
  echo "$failures manifest-diff case(s) failed" >&2
  exit 1
fi
echo "OK: manifest-diff.sh behaves"

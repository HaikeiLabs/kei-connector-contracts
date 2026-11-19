#!/usr/bin/env bash
# Copyright 2026 Haikei Labs
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Tests scripts/no-ticket-ids.sh in a throwaway git repository: a tracked
# file containing a ticket id fails with the match listed; a clean repo
# passes; a ticket id under .github/workflows is excluded.
#
# The id is built from pieces below so this test file itself never matches
# the HAI-[0-9]+ pattern the check greps for.
set -uo pipefail

cd "$(dirname "$0")/.."
check=scripts/no-ticket-ids.sh
id="HAI-""123"
failures=0

expect() { # expect <want exit> <want output substring> <repo>
  local want_code=$1 want_out=$2 repo=$3
  local out code
  out=$(bash "$check" "$repo" 2>&1)
  code=$?
  if [ "$code" -ne "$want_code" ] || ! grep -qF -- "$want_out" <<<"$out"; then
    echo "FAIL: no-ticket-ids.sh in $repo: exit $code, want $want_code with '$want_out'; output:" >&2
    echo "$out" >&2
    failures=$((failures + 1))
  else
    echo "ok: no-ticket-ids.sh in $repo -> exit $code"
  fi
}

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

new_repo() { # new_repo <name>
  local dir="$tmp/$1"
  mkdir -p "$dir"
  git -C "$dir" init -q
  git -C "$dir" config user.email test@example.com
  git -C "$dir" config user.name test
}

# 1. A tracked file with a ticket id fails and lists the match.
new_repo dirty
printf 'see %s for details\n' "$id" > "$tmp/dirty/README.md"
git -C "$tmp/dirty" add -A && git -C "$tmp/dirty" commit -qm init
expect 1 "$id" "$tmp/dirty"

# 2. A clean repo passes.
new_repo clean
printf 'no ids here\n' > "$tmp/clean/README.md"
git -C "$tmp/clean" add -A && git -C "$tmp/clean" commit -qm init
expect 0 "OK: no internal ticket ids" "$tmp/clean"

# 3. A ticket id under .github/workflows is excluded.
new_repo workflow
mkdir -p "$tmp/workflow/.github/workflows"
printf 'ref: %s\n' "$id" > "$tmp/workflow/.github/workflows/ci.yaml"
git -C "$tmp/workflow" add -A && git -C "$tmp/workflow" commit -qm init
expect 0 "OK: no internal ticket ids" "$tmp/workflow"

if [ "$failures" -gt 0 ]; then
  echo "FAIL: $failures test(s) failed" >&2
  exit 1
fi
echo "ok: all no-ticket-ids.sh tests passed"

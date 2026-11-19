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
# Fails if any tracked file contains an internal ticket id (HAI-<number>).
# Internal ticket ids must not leak into the public repository.
# .github/workflows is excluded: workflow changes require owner oversight and
# the last reference there is removed by a separate owner-run step.
#
# Usage:
#   scripts/no-ticket-ids.sh [repo-dir]   # default: current directory
set -uo pipefail

cd "${1:-.}"

matches="$(git grep -nE 'HAI-[0-9]+' -- . ':!.github/workflows' || true)"

if [ -n "$matches" ]; then
  echo "ERROR: internal ticket ids found:" >&2
  echo "$matches" >&2
  echo >&2
  echo "Remove the ticket ids (keep the surrounding text readable)." >&2
  exit 1
fi

echo "OK: no internal ticket ids"

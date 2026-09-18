#!/usr/bin/env bash
# Adds or verifies the Apache 2.0 header on every Go file.
#
# Apache 2.0 does not strictly require per-file headers — the appendix says
# "we recommend" — and the binding obligations are the LICENSE file and
# preserving notices on redistribution. Headers are still worth having on a
# public repository: a single file copied into someone else's project carries
# its licence with it, which a repository-level LICENSE alone does not achieve.
#
# Usage:
#   scripts/license-header.sh          # add the header where missing
#   scripts/license-header.sh --check  # exit 1 if any file is missing it
set -euo pipefail

cd "$(dirname "$0")/.."

HEADER='// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.'

check_only=false
if [ "${1:-}" = "--check" ]; then
  check_only=true
fi

missing=()

# -print0/read -d handles paths with spaces. Generated files are excluded:
# a header cannot survive regeneration.
while IFS= read -r -d '' file; do
  if head -1 "$file" | grep -q '^// Code generated .* DO NOT EDIT\.$'; then
    continue
  fi

  if head -3 "$file" | grep -q 'Licensed under the Apache License'; then
    continue
  fi

  missing+=("$file")

  if [ "$check_only" = false ]; then
    tmp="$(mktemp)"
    printf '%s\n\n' "$HEADER" > "$tmp"
    cat "$file" >> "$tmp"
    mv "$tmp" "$file"
    echo "added: $file"
  fi
done < <(find . -name '*.go' -not -path './vendor/*' -print0)

if [ "$check_only" = true ]; then
  if [ ${#missing[@]} -gt 0 ]; then
    echo "ERROR: ${#missing[@]} Go file(s) missing the Apache 2.0 header:" >&2
    printf '  %s\n' "${missing[@]}" >&2
    echo >&2
    echo "Run scripts/license-header.sh to add it." >&2
    exit 1
  fi
  echo "OK: every Go file carries the Apache 2.0 header"
else
  if [ ${#missing[@]} -eq 0 ]; then
    echo "OK: every Go file already carries the header"
  else
    echo "added the header to ${#missing[@]} file(s)"
  fi
fi

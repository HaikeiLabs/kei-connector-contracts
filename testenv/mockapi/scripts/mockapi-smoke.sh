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

# Starts the contract mock, drives a few requests through it and reports
# the results. Run from this directory or from the repository root.
set -euo pipefail
cd "$(dirname "$0")/.."

PORT="${MOCKAPI_PORT:-18099}"
ADDR="127.0.0.1:${PORT}"

go build ./cmd/mockapi
./mockapi -addr "${ADDR}" -manifest ../../manifest/connector-manifest.v1.json -fixtures ../../fixtures/operations &
MOCK_PID=$!
trap 'kill "${MOCK_PID}" 2>/dev/null || true' EXIT

for _ in $(seq 1 50); do
  if curl -s -o /dev/null "http://${ADDR}/repos/acme/widgets/issues/1347"; then break; fi
  sleep 0.1
done

echo "== issues.get (fixture, lossless id above 2^53)"
curl -s -H "Authorization: Bearer smoke" "http://${ADDR}/repos/acme/widgets/issues/1347" | head -c 200; echo
echo "== X-Mock-Fail: 429 (Retry-After expected)"
curl -s -i -H "Authorization: Bearer smoke" -H "X-Mock-Fail: 429" "http://${ADDR}/repos/acme/widgets/issues/1347" | head -3
echo "== wrong credential placement (400)"
curl -s -o /dev/null -w '%{http_code}\n' "http://${ADDR}/repos/acme/widgets/issues/1?token=smoke"
echo "== unescaped path (400)"
curl -s -o /dev/null -w '%{http_code}\n' -H "Authorization: Bearer smoke" "http://${ADDR}/repos/acme(ops)/widgets/issues/1"
echo "== graphql root field routing"
curl -s -H "Authorization: Bearer smoke" -H "Content-Type: application/json" \
  -d '{"query":"query($id: ID!) { issue(id: $id) { identifier } }","variables":{"id":"00000000-0000-4000-8000-000000000101"}}' \
  "http://${ADDR}/graphql" | head -c 200; echo
echo "== load (2s at 20 rps)"
go run ./cmd/mockload -addr "http://${ADDR}" -rps 20 -duration 2s

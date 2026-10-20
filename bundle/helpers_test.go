// Copyright 2026 Haikei Labs
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
// limitations under the License.

package bundle

import (
	"time"
)

// The v2 bundle document the v2 tests build. It is the catalog's shape
// without the canonical JSON Schema documents: those stay with the service
// that compiles them, so this package's own tests use the document directly.
const (
	v2Installation = "inst-v2"
	v2Org          = "org-v2"
	v2Workspace    = "ws-v2"
)

func v2BundleDoc(now time.Time) map[string]any {
	return map[string]any{
		"schema": "kei.policy-bundle/v2", "bundle_id": "bundle-v2", "bundle_version": 2, "policy_revision": 7,
		"issuer": Issuer, "audience": map[string]any{"installation_id": v2Installation, "org_id": v2Org, "workspace_id": v2Workspace},
		"issued_at": now.UTC().Format(time.RFC3339Nano), "not_after": now.Add(12 * time.Hour).UTC().Format(time.RFC3339Nano),
		"refresh":            map[string]any{"poll_interval_seconds": 60, "stale_after_seconds": 300},
		"subject_resolution": map[string]any{"mode": "bundle", "max_cache_seconds": 30},
		"policy_set": map[string]any{
			"schema": "kei.policy-set/v2", "org_id": v2Org, "workspace_id": v2Workspace, "match_semantics": "kei.match/v2",
			"evaluation": map[string]any{"precedence": []string{"priority_desc", "deny_before_permit", "policy_id_asc"}, "checks": "per_capability_first_match_and", "default": "deny"},
			"tools":      []any{}, "checks": []any{}, "tool_policies": []any{},
			"harness_match_semantics": "kei.harness-match/v1", "harness_policies": []any{},
		},
		"agent_connector_bindings": []any{},
		"identities":               []any{map[string]any{"source": "discord", "external_id": "ext-alice", "user_id": "user-alice", "status": "linked", "groups": []string{"members"}}},
	}
}

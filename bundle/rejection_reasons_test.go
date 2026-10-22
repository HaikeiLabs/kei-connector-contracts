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
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRejectionReasons(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)

	// validV1 is a minimal valid v1 bundle skeleton used as a base for
	// rejection-path cases that must reach a specific check.
	validV1 := func() map[string]any {
		return map[string]any{
			"schema": "kei.policy-bundle/v1", "bundle_id": "b-1", "bundle_version": 1, "policy_revision": 1,
			"issuer":                   Issuer,
			"audience":                 map[string]any{"installation_id": "inst", "org_id": "org", "workspace_id": "ws"},
			"issued_at":                now.Add(-1 * time.Hour).UTC().Format(time.RFC3339Nano),
			"not_after":                now.Add(12 * time.Hour).UTC().Format(time.RFC3339Nano),
			"refresh":                  map[string]any{"poll_interval_seconds": 60, "stale_after_seconds": 300},
			"subject_resolution":       map[string]any{"mode": "online"},
			"policy_set":               map[string]any{"match_semantics": "kei.match/v1", "evaluation": map[string]any{"order": "array_order_first_match", "default": "deny"}, "groups": []any{}, "agents": []any{}, "tools": []any{}, "policies": []any{}},
			"agent_connector_bindings": []any{},
		}
	}

	// validV2 is a minimal valid v2 bundle skeleton.
	validV2 := func() map[string]any {
		return map[string]any{
			"schema": "kei.policy-bundle/v2", "bundle_id": "b-v2", "bundle_version": 1, "policy_revision": 1,
			"issuer":                   Issuer,
			"audience":                 map[string]any{"installation_id": "inst", "org_id": "org", "workspace_id": "ws"},
			"issued_at":                now.Add(-1 * time.Hour).UTC().Format(time.RFC3339Nano),
			"not_after":                now.Add(12 * time.Hour).UTC().Format(time.RFC3339Nano),
			"refresh":                  map[string]any{"poll_interval_seconds": 60, "stale_after_seconds": 300},
			"subject_resolution":       map[string]any{"mode": "bundle", "max_cache_seconds": 30},
			"policy_set":               map[string]any{"schema": "kei.policy-set/v2", "org_id": "org", "workspace_id": "ws", "match_semantics": "kei.match/v2", "evaluation": map[string]any{"precedence": []string{"priority_desc"}, "checks": "per_capability_first_match_and", "default": "deny"}, "tools": []any{}, "checks": []any{}, "tool_policies": []any{}, "harness_match_semantics": "kei.harness-match/v1", "harness_policies": []any{}},
			"agent_connector_bindings": []any{},
			"identities":               []any{map[string]any{"source": "discord", "external_id": "ext-a", "user_id": "user-a", "status": "linked"}},
		}
	}

	cases := []struct {
		name   string
		bundle func() map[string]any
		// validateFn is either Validate or ValidateV2
		validateFn func([]byte, string, string, string, int64, time.Time) (*Verified, error)
		wantReason string // substring expected in the error message
	}{
		{
			name: "v1 schema not v1 or v2",
			bundle: func() map[string]any {
				d := validV1()
				d["schema"] = "kei.policy-bundle/v3"
				return d
			},
			validateFn: Validate,
			wantReason: "schema",
		},
		{
			name: "v1 subject_resolution.mode hybrid",
			bundle: func() map[string]any {
				d := validV1()
				// Use a mode that is neither online nor bundle
				d["subject_resolution"] = map[string]any{"mode": "hybrid"}
				return d
			},
			validateFn: Validate,
			wantReason: "subject_resolution.mode",
		},
		{
			name: "v2 subject_resolution.mode not bundle",
			bundle: func() map[string]any {
				d := validV2()
				d["subject_resolution"] = map[string]any{"mode": "online"}
				return d
			},
			validateFn: Validate,
			wantReason: "subject_resolution.mode",
		},
		{
			name: "v2 identities nil",
			bundle: func() map[string]any {
				d := validV2()
				d["identities"] = nil
				return d
			},
			validateFn: Validate,
			wantReason: "identities",
		},
		{
			name: "v2 policy_set.schema mismatch",
			bundle: func() map[string]any {
				d := validV2()
				ps := d["policy_set"].(map[string]any)
				ps["schema"] = "kei.policy-set/v1"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.schema",
		},
		{
			name: "v2 policy_set.match_semantics mismatch",
			bundle: func() map[string]any {
				d := validV2()
				ps := d["policy_set"].(map[string]any)
				ps["match_semantics"] = "kei.match/v1"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.match_semantics",
		},
		{
			name: "v2 policy_set.evaluation.checks mismatch",
			bundle: func() map[string]any {
				d := validV2()
				ps := d["policy_set"].(map[string]any)
				eval := ps["evaluation"].(map[string]any)
				eval["checks"] = "first_match"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.evaluation.checks",
		},
		{
			name: "v2 policy_set.evaluation.default mismatch",
			bundle: func() map[string]any {
				d := validV2()
				ps := d["policy_set"].(map[string]any)
				eval := ps["evaluation"].(map[string]any)
				eval["default"] = "permit"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.evaluation.default",
		},
		{
			name: "v1 policy_set.match_semantics mismatch",
			bundle: func() map[string]any {
				d := validV1()
				ps := d["policy_set"].(map[string]any)
				ps["match_semantics"] = "kei.match/v2"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.match_semantics",
		},
		{
			name: "v1 policy_set.evaluation.order mismatch",
			bundle: func() map[string]any {
				d := validV1()
				ps := d["policy_set"].(map[string]any)
				eval := ps["evaluation"].(map[string]any)
				eval["order"] = "priority_desc"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.evaluation.order",
		},
		{
			name: "v1 policy_set.evaluation.default mismatch",
			bundle: func() map[string]any {
				d := validV1()
				ps := d["policy_set"].(map[string]any)
				eval := ps["evaluation"].(map[string]any)
				eval["default"] = "permit"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set.evaluation.default",
		},
		{
			name: "v2 policy_set strictJSON failure",
			bundle: func() map[string]any {
				d := validV2()
				// Introduce an unknown field to fail strict JSON decoding
				ps := d["policy_set"].(map[string]any)
				ps["unknown_field"] = "value"
				return d
			},
			validateFn: Validate,
			wantReason: "policy_set",
		},
		{
			name: "v2 plan validator rejects tool",
			bundle: func() map[string]any {
				d := validV2()
				ps := d["policy_set"].(map[string]any)
				ps["tools"] = []any{map[string]any{
					"name": "rejected_tool", "source": "github", "operation_class": "read",
					"route":                 map[string]any{"connector_binding": map[string]any{"connector_id": "conn-a"}},
					"required_capabilities": []any{"pull_request.read"},
					"resource_types":        []any{map[string]any{"type": "pull_request", "parent_type": "repository"}},
					"plan": map[string]any{
						"args_schema":    map[string]any{"type": "object", "additionalProperties": false},
						"context_schema": map[string]any{"type": "object", "additionalProperties": false},
						"operations":     []any{},
					},
				}}
				return d
			},
			validateFn: Validate,
			wantReason: "tool rejected_tool",
		},
	}

	// Install a schema validator for v1 and v2. The v2 path also needs a plan
	// validator; install a rejecting one so the plan-rejection case works.
	SetSchemaValidator(func([]byte) error { return nil })
	SetV2SchemaValidator(func([]byte) error { return nil })
	SetPlanValidator(func(tool V2Tool) error {
		if len(tool.Plan) == 0 {
			return nil
		}
		return errors.New("plan rejected by test validator")
	})
	t.Cleanup(func() {
		SetSchemaValidator(nil)
		SetV2SchemaValidator(nil)
		SetPlanValidator(nil)
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(tc.bundle())
			if err != nil {
				t.Fatal(err)
			}
			_, err = tc.validateFn(raw, "inst", "org", "ws", 0, now)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !errors.Is(err, ErrBundleUnsupported) {
				t.Errorf("errors.Is(err, ErrBundleUnsupported) = false; err = %v", err)
			}
			if !strings.Contains(err.Error(), tc.wantReason) {
				t.Errorf("error %q does not contain reason %q", err.Error(), tc.wantReason)
			}
		})
	}
}

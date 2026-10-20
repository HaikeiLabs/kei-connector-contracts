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
	"testing"
	"time"
)

// planBundleDoc is a v2 bundle whose policy set registers one connector tool
// with an operation plan.
func planBundleDoc(t *testing.T, now time.Time) []byte {
	t.Helper()
	doc := v2BundleDoc(now)
	set := doc["policy_set"].(map[string]any)
	set["tools"] = []any{map[string]any{
		"name": "github.review_pr", "source": "github", "operation_class": "read",
		"route":                 map[string]any{"connector_binding": map[string]any{"connector_id": "conn-a"}},
		"required_capabilities": []any{"pull_request.read"},
		"resource_types":        []any{map[string]any{"type": "pull_request", "parent_type": "repository"}},
		"plan": map[string]any{
			"args_schema":    map[string]any{"type": "object", "properties": map[string]any{"number": map[string]any{"type": "integer"}}, "required": []any{"number"}, "additionalProperties": false},
			"context_schema": map[string]any{"type": "object", "properties": map[string]any{"repository": map[string]any{"type": "string"}}, "required": []any{"repository"}, "additionalProperties": false},
			"operations": []any{map[string]any{
				"id": "read-pull-request", "capability": "pull_request.read",
				"resource":                   map[string]any{"type": "pull_request", "id": map[string]any{"from": "args", "pointer": "/number", "type": "integer"}, "parent": map[string]any{"type": "repository", "id": map[string]any{"from": "context", "field": "repository", "type": "string"}}},
				"provider_resource_template": "repos/{parent.id}/pulls/{resource.id}", "provider_input": map[string]any{},
			}},
		},
	}}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// useV2SchemaValidator installs a v2 schema validator for the lifetime of a
// test. The canonical JSON Schema document stays with the service that
// compiles it; this package's own tests only need a validator that accepts.
func useV2SchemaValidator(t *testing.T) {
	t.Helper()
	previous := schemaValidatorV2
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(previous) })
}

// usePlanValidator installs a plan validator for the lifetime of a test.
func usePlanValidator(t *testing.T, v PlanValidator) {
	t.Helper()
	previous := planValidator
	SetPlanValidator(v)
	t.Cleanup(func() { SetPlanValidator(previous) })
}

// Without a plan validator the plan is not checked: this package carries no
// JSON Schema implementation, so the plan profile is the caller's business.
func TestPlanValidationIsSkippedWithoutAValidator(t *testing.T) {
	useV2SchemaValidator(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	if _, err := ValidateV2(planBundleDoc(t, now), v2Installation, v2Org, v2Workspace, 1, now); err != nil {
		t.Fatalf("ValidateV2 without a plan validator: %v", err)
	}
}

// A plan validator that rejects the plan makes the whole bundle unsupported,
// exactly as an invalid plan did before the validator was pluggable.
func TestPlanValidatorRejectsTheBundle(t *testing.T) {
	useV2SchemaValidator(t)
	usePlanValidator(t, func(V2Tool) error { return errors.New("plan is invalid") })
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	_, err := ValidateV2(planBundleDoc(t, now), v2Installation, v2Org, v2Workspace, 1, now)
	if !errors.Is(err, ErrBundleUnsupported) {
		t.Fatalf("ValidateV2 with a rejecting plan validator = %v, want ErrBundleUnsupported", err)
	}
}

// A plan validator that accepts the plan leaves the bundle valid.
func TestPlanValidatorAcceptsTheBundle(t *testing.T) {
	useV2SchemaValidator(t)
	var validated int
	usePlanValidator(t, func(V2Tool) error {
		validated++
		return nil
	})
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	if _, err := ValidateV2(planBundleDoc(t, now), v2Installation, v2Org, v2Workspace, 1, now); err != nil {
		t.Fatalf("ValidateV2 with an accepting plan validator: %v", err)
	}
	if validated == 0 {
		t.Fatal("the plan validator was not called")
	}
}

// A tool whose plan the validator rejects is not a registered connector tool.
func TestRejectedPlanIsNotARegisteredConnectorTool(t *testing.T) {
	useV2SchemaValidator(t)
	usePlanValidator(t, func(V2Tool) error { return nil })
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	v, err := ValidateV2(planBundleDoc(t, now), v2Installation, v2Org, v2Workspace, 1, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := v.V2ConnectorTool("github.review_pr", now); !ok {
		t.Fatal("a tool with an accepted plan did not resolve as a registered connector tool")
	}
	SetPlanValidator(func(V2Tool) error { return errors.New("plan is invalid") })
	if _, ok := v.V2ConnectorTool("github.review_pr", now); ok {
		t.Fatal("a tool with a rejected plan resolved as a registered connector tool")
	}
}

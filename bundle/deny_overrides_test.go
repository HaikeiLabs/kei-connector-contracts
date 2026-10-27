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
	"testing"
	"time"
)

// v1Verified builds a v1 bundle whose single registered tool "t" is governed by
// the given policies, in array order, under the given evaluation.order.
func v1Verified(t *testing.T, order string, policies []map[string]any) *Verified {
	t.Helper()
	for _, p := range policies {
		p["scope"] = map[string]any{}
	}
	set := map[string]any{
		"match_semantics": "kei.match/v1", "evaluation": map[string]any{"order": order, "default": "deny"},
		"groups": []any{}, "agents": []any{}, "policies": policies,
		"tools": []any{map[string]any{"name": "t", "service": "svc", "action": "read"}},
	}
	raw, err := json.Marshal(set)
	if err != nil {
		t.Fatal(err)
	}
	return &Verified{Payload: Payload{PolicySet: raw}}
}

func TestV1DenyOverrides(t *testing.T) {
	permit := map[string]any{"id": "p-permit", "src_pattern": "*", "dst_pattern": "tool:t", "action": "permit", "priority": 100}
	deny := map[string]any{"id": "p-deny", "src_pattern": "*", "dst_pattern": "tool:t", "action": "deny", "priority": 1}
	call := Call{Tool: "t"}

	t.Run("deny_overrides: a later deny beats an earlier permit", func(t *testing.T) {
		d := v1Verified(t, "deny_overrides", []map[string]any{permit, deny}).Evaluate(call, nil)
		if d.Outcome != "deny" || d.PolicyID != "p-deny" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("deny_overrides: a permit with no deny allows", func(t *testing.T) {
		d := v1Verified(t, "deny_overrides", []map[string]any{permit}).Evaluate(call, nil)
		if d.Outcome != "allow" || d.PolicyID != "p-permit" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("deny_overrides: no match denies", func(t *testing.T) {
		d := v1Verified(t, "deny_overrides", nil).Evaluate(call, nil)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("array_order_first_match is unchanged", func(t *testing.T) {
		d := v1Verified(t, "array_order_first_match", []map[string]any{permit, deny}).Evaluate(call, nil)
		if d.Outcome != "allow" || d.PolicyID != "p-permit" {
			t.Fatalf("decision = %+v", d)
		}
	})
}

func TestV1ValidateAcceptsDenyOverrides(t *testing.T) {
	for _, order := range []string{"array_order_first_match", "deny_overrides"} {
		raw, _ := json.Marshal(map[string]any{
			"match_semantics": "kei.match/v1", "evaluation": map[string]any{"order": order, "default": "deny"},
			"groups": []any{}, "agents": []any{}, "tools": []any{}, "policies": []any{},
		})
		if err := validatePolicyDialect(raw); err != nil {
			t.Errorf("order %q rejected: %v", order, err)
		}
	}
}

func TestV2DenyOverridesPrecedence(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}
	policies := []mcpTestPolicy{
		{id: "p-permit", src: "*", dst: "mcp:slack", action: "permit", priority: 100},
		{id: "p-deny", src: "*", dst: "mcp:slack.list_channels", action: "deny", priority: 1},
	}

	t.Run("priority_desc keeps the higher-priority permit", func(t *testing.T) {
		v := v2MCPBundle(t, now, policies)
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p-permit" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("deny_overrides: a lower-priority deny wins", func(t *testing.T) {
		v := setV2Precedence(t, v2MCPBundle(t, now, policies), "deny_overrides")
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "deny" || d.PolicyID != "p-deny" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("deny_overrides: harness evaluation puts the deny first", func(t *testing.T) {
		v := setV2Precedence(t, v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-permit", src: "*", dst: "shell:git status", action: "permit", priority: 100},
			{id: "p-deny", src: "*", dst: "shell:git", action: "deny", priority: 1},
		}), "deny_overrides")
		outcome, id, _ := v.EvaluateHarness(HarnessCall{Kind: "claude_code", Argv: []string{"git", "status"}})
		if outcome != "deny" || id != "p-deny" {
			t.Fatalf("outcome=%s id=%s", outcome, id)
		}
	})
}

// setV2Precedence rewrites the bundle's evaluation.precedence to the given id.
func setV2Precedence(t *testing.T, v *Verified, id string) *Verified {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(v.PayloadBytes, &payload); err != nil {
		t.Fatal(err)
	}
	set := payload["policy_set"].(map[string]any)
	set["evaluation"].(map[string]any)["precedence"] = []string{id}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var p Payload
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	return &Verified{Payload: p, PayloadBytes: raw, Digest: v.Digest}
}

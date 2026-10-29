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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	pol "github.com/HaikeiLabs/kei-connector-contracts/policy"
)

func TestMatchTargetExact(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		pattern string
		want    bool
	}{
		{"agent exact", "agent:abc", "agent:abc", true},
		{"capability exact", "capability:issue.read", "capability:issue.read", true},
		{"connector exact", "connector:jira", "connector:jira", true},
		{"tool exact", "tool:search_wiki", "tool:search_wiki", true},
		{"user exact", "user:alice", "user:alice", true},
		{"group match", "group:engineering", "group:engineering", true},

		{"agent suffix must not match", "agent:abc", "agent:xabc", false},
		{"agent suffix must not match 2", "agent:abc", "agent:bc", false},
		{"capability suffix must not match", "capability:issue.read", "capability:read", false},
		{"capability v1 target", "capability:issue.read", "capability:read", false},
		{"connector suffix must not match", "connector:myjira", "connector:jira", false},
		{"connector v1 target", "connector:myjira", "connector:jira", false},
		{"connector prefix must not match", "connector:jira", "connector:jira_v2", false},
		{"tool suffix must not match", "tool:search_wiki_admin", "tool:wiki", false},
		{"user suffix must not match", "user:u-123", "user:evil-u-123", false},

		{"star-suffix is literal", "otherleads/42", "*leads/42", false},
		{"star-suffix scheme", "capability:issue.read", "*read", false},

		{"scheme wildcard", "agent:anything", "agent:*", true},
		{"scheme wildcard capability", "capability:issue.read", "capability:*", true},
		{"scheme wildcard tool", "tool:search_wiki", "tool:*", true},
		{"scheme wildcard no cross-scheme", "agent:a-1", "service:*", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pol.MatchTarget(tc.pattern, tc.input); got != tc.want {
				t.Errorf("MatchTarget(%q, %q) = %v, want %v", tc.pattern, tc.input, got, tc.want)
			}
		})
	}
}

func TestMatchPolicyPatternExact(t *testing.T) {
	cases := []struct {
		name    string
		input1  string
		input2  string
		pattern string
		want    bool
	}{
		{"user exact on input1", "alice", "alice@corp.com", "user:alice", true},
		{"email exact on input2", "u-1", "alice@corp.com", "email:alice@corp.com", true},
		{"email exact on input1", "alice@corp.com", "", "email:alice@corp.com", true},

		{"user suffix on input1", "evil-alice", "", "user:alice", false},
		{"email suffix on input2", "u-1", "malice@corp.com", "email:alice@corp.com", false},

		{"domain on input2", "u-1", "alice@corp.com", "*@corp.com", true},
		{"domain on input1", "alice@corp.com", "", "*@corp.com", true},
		{"domain case-insensitive", "u-1", "alice@CORP.com", "*@corp.com", true},
		{"domain evil suffix rejected", "u-1", "alice@evilcorp.com", "*@corp.com", false},
		{"domain subdomain rejected", "u-1", "alice@mail.corp.com", "*@corp.com", false},
		{"domain no local part rejected", "", "@corp.com", "*@corp.com", false},
		{"domain case-insensitive reverse", "ALICE@corp.com", "", "*@CORP.COM", true},
		{"domain extra labels rejected", "alice@corp.com.evil.io", "", "*@corp.com", false},

		{"email domain form", "u-1", "alice@corp.com", "email:*@corp.com", true},
		{"email domain case-insensitive", "u-1", "Alice@Corp.COM", "email:*@corp.com", true},
		{"email domain lookalike rejected", "u-1", "alice@evilcorp.com", "email:*@corp.com", false},

		{"scheme wildcard user", "user:anything", "", "user:*", true},
		{"scheme wildcard needs value", "user:", "", "user:*", false},

		{"star-suffix is literal in input1", "evil-alice", "", "*alice", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pol.MatchPrincipal(tc.pattern, tc.input1, tc.input2, nil); got != tc.want {
				t.Errorf("MatchPrincipal(%q, %q, %q, nil) = %v, want %v",
					tc.pattern, tc.input1, tc.input2, got, tc.want)
			}
		})
	}
}

func v2ToolBundle(t *testing.T, now time.Time, checks []any) *Verified {
	t.Helper()
	doc := v2BundleDoc(now)
	doc["agent_connector_bindings"] = []any{map[string]any{"agent_id": "agent-a", "connector_id": "conn-a", "capabilities": []string{"pull_request.read", "check.read"}}}
	set := doc["policy_set"].(map[string]any)
	set["tools"] = []any{map[string]any{
		"name": "github.review_pr", "source": "github", "operation_class": "read",
		"route":                 map[string]any{"connector_binding": map[string]any{"connector_id": "conn-a"}},
		"required_capabilities": []string{"pull_request.read", "check.read"},
		"resource_types":        []any{map[string]any{"type": "pull_request", "parent_type": "repository"}},
		"plan": map[string]any{
			"args_schema":    map[string]any{"type": "object", "properties": map[string]any{"number": map[string]any{"type": "integer"}}, "required": []any{"number"}, "additionalProperties": false},
			"context_schema": map[string]any{"type": "object", "properties": map[string]any{"repository": map[string]any{"type": "string"}}, "required": []any{"repository"}, "additionalProperties": false},
			"operations": []any{
				map[string]any{"id": "read-pull-request", "capability": "pull_request.read", "resource": map[string]any{"type": "pull_request", "id": map[string]any{"from": "args", "pointer": "/number", "type": "integer"}, "parent": map[string]any{"type": "repository", "id": map[string]any{"from": "context", "field": "repository", "type": "string"}}}, "provider_resource_template": "repos/{parent.id}/pulls/{resource.id}", "provider_input": map[string]any{}},
				map[string]any{"id": "read-checks", "capability": "check.read", "resource": map[string]any{"type": "pull_request", "id": map[string]any{"from": "args", "pointer": "/number", "type": "integer"}, "parent": map[string]any{"type": "repository", "id": map[string]any{"from": "context", "field": "repository", "type": "string"}}}, "provider_resource_template": "repos/{parent.id}/pulls/{resource.id}", "provider_input": map[string]any{}},
			},
		},
	}}
	set["checks"] = checks
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return &Verified{Payload: func() Payload {
		var payload Payload
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatal(err)
		}
		return payload
	}(), PayloadBytes: raw, Digest: "sha256:test-v2"}
}

func v2Permit(policyID, capability string) any {
	return map[string]any{
		"policy_id": policyID, "name": policyID, "effect": "permit", "priority": 10,
		"scope":    map[string]any{"org_id": "org-v2", "workspace_id": "ws-v2", "agent_id": "agent-a"},
		"selector": map[string]any{"source_principal": "user:user-alice", "source": "github", "connector_id": "conn-a", "tool": "github.review_pr", "capability": capability, "resource_type": "pull_request", "resource_id": "42", "resource_parent": "acme/project"},
	}
}

func v2TestPlanOps() []Operation {
	sum := sha256.Sum256([]byte(`{}`))
	digest := "sha256:" + hex.EncodeToString(sum[:])
	return []Operation{
		{ID: "read-pull-request", ConnectorID: "conn-a", Capability: "pull_request.read", Action: "read", Resource: "repos/acme/project/pulls/42", InputDigest: digest, ProviderInput: json.RawMessage(`{}`), ResourceType: "pull_request", ResourceID: "42", ResourceParent: "acme/project", ResourceParentType: "repository"},
		{ID: "read-checks", ConnectorID: "conn-a", Capability: "check.read", Action: "read", Resource: "repos/acme/project/pulls/42", InputDigest: digest, ProviderInput: json.RawMessage(`{}`), ResourceType: "pull_request", ResourceID: "42", ResourceParent: "acme/project", ResourceParentType: "repository"},
	}
}

func v2HarnessBundle(t *testing.T, now time.Time, status, effect string) *Verified {
	t.Helper()
	v := v2ToolBundle(t, now, nil)
	var doc map[string]any
	if err := json.Unmarshal(v.PayloadBytes, &doc); err != nil {
		t.Fatal(err)
	}
	set := doc["policy_set"].(map[string]any)
	tools := set["tools"].([]any)
	tool := tools[0].(map[string]any)
	tool["source"] = "discord"
	tool["route"] = map[string]any{"harness_executor": map[string]any{"executor": "discord", "registration": "web_search"}}
	delete(tool, "required_capabilities")
	delete(tool, "resource_types")
	delete(tool, "plan")
	set["tool_policies"] = []any{map[string]any{
		"policy_id": "harness-policy", "name": "harness-policy", "effect": effect, "priority": 10,
		"scope":    map[string]any{"org_id": "org-v2", "workspace_id": "ws-v2"},
		"selector": map[string]any{"source_principal": "user:user-alice", "source": "discord", "tool": "github.review_pr"},
	}}
	identities := doc["identities"].([]any)
	identities[0].(map[string]any)["status"] = status
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var payload Payload
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	return &Verified{Payload: payload, PayloadBytes: raw, Digest: "sha256:harness-test"}
}

func TestEvaluateV2HarnessToolRequiresFreshLinkedSessionAndPolicy(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}
	for status, reason := range map[string]string{"unlinked": "subject_unlinked", "pending": "subject_pending", "revoked": "subject_revoked", "cross_workspace": "subject_cross_workspace"} {
		t.Run(status, func(t *testing.T) {
			v := v2HarnessBundle(t, now, status, "permit")
			if d := v.EvaluateV2HarnessTool("github.review_pr", "", subject, now); d.ReasonCode != reason || d.Outcome != "deny" {
				t.Fatalf("decision = %+v", d)
			}
		})
	}
	t.Run("permit", func(t *testing.T) {
		v := v2HarnessBundle(t, now, "linked", "permit")
		if d := v.EvaluateV2HarnessTool("github.review_pr", "", subject, now); d.Outcome != "allow" || d.PolicyID != "harness-policy" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("deny", func(t *testing.T) {
		v := v2HarnessBundle(t, now, "linked", "deny")
		if d := v.EvaluateV2HarnessTool("github.review_pr", "", subject, now); d.Outcome != "deny" || d.ReasonCode != "policy_deny" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("stale snapshot is now fresh", func(t *testing.T) {
		v := v2HarnessBundle(t, now, "linked", "permit")
		if d := v.EvaluateV2HarnessTool("github.review_pr", "", subject, now.Add(30*time.Second)); d.Outcome != "allow" || d.PolicyID != "harness-policy" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("connector route is not harness registered", func(t *testing.T) {
		v := v2ToolBundle(t, now, nil)
		if d := v.EvaluateV2HarnessTool("github.review_pr", "", subject, now); d.ReasonCode != "tool_not_registered" {
			t.Fatalf("decision = %+v", d)
		}
	})
}

func TestEvaluateV2PlanRequiresANDForAllDeclaredCapabilities(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}
	ops := v2TestPlanOps()
	t.Run("both permits allow", func(t *testing.T) {
		v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read"), v2Permit("p2", "check.read")})
		if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, ops, now); d.Outcome != "allow" {
			t.Fatalf("decision = %+v", d)
		}
	})
	t.Run("missing capability check denies", func(t *testing.T) {
		v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read")})
		if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, ops, now); d.Outcome != "deny" {
			t.Fatalf("decision = %+v", d)
		}
	})
	for name, mutate := range map[string]func([]Operation){
		"extra capability":       func(plan []Operation) { plan[1].Capability = "secret.write" },
		"parent type mismatch":   func(plan []Operation) { plan[0].ResourceParentType = "organization" },
		"missing parent context": func(plan []Operation) { plan[0].ResourceParent = "" },
	} {
		t.Run(name, func(t *testing.T) {
			v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read"), v2Permit("p2", "check.read")})
			plan := append([]Operation(nil), ops...)
			mutate(plan)
			if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, plan, now); d.Outcome != "deny" || d.ReasonCode != "operation_plan_invalid" {
				t.Fatalf("decision = %+v", d)
			}
		})
	}
	t.Run("resource id mismatch is denied by structured selector", func(t *testing.T) {
		v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read"), v2Permit("p2", "check.read")})
		plan := append([]Operation(nil), ops...)
		plan[0].ResourceID = "43"
		if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, plan, now); d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("resource id mismatch decision=%+v", d)
		}
	})
	t.Run("parent mismatch is denied by structured selector", func(t *testing.T) {
		v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read"), v2Permit("p2", "check.read")})
		plan := append([]Operation(nil), ops...)
		plan[0].ResourceParent = "acme/other"
		if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, plan, now); d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("parent mismatch decision=%+v", d)
		}
	})
}

func TestV2DirectCapabilityPermitDoesNotImplyToolPermit(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	v := v2ToolBundle(t, now, []any{
		map[string]any{"policy_id": "direct-read", "name": "direct", "effect": "permit", "priority": 100, "scope": map[string]any{"org_id": "org-v2", "workspace_id": "ws-v2", "agent_id": "agent-a"}, "selector": map[string]any{"source_principal": "user:user-alice", "capability": "pull_request.read"}},
		map[string]any{"policy_id": "direct-check", "name": "direct", "effect": "permit", "priority": 100, "scope": map[string]any{"org_id": "org-v2", "workspace_id": "ws-v2", "agent_id": "agent-a"}, "selector": map[string]any{"source_principal": "user:user-alice", "capability": "check.read"}},
	})
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}
	ops := v2TestPlanOps()
	if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, ops, now); d.Outcome != "deny" {
		t.Fatalf("direct capability permits authorized tool call: %+v", d)
	}
}

func TestEvaluateV2PlanAcceptsSnapshotAndRefreshSubstitution(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	v := v2ToolBundle(t, now, []any{v2Permit("p1", "pull_request.read"), v2Permit("p2", "check.read")})
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}
	ops := v2TestPlanOps()
	// HAI-452: snapshot is fresh for the bundle lifetime, so plan evaluation
	// proceeds and all permits allow.
	if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, ops, now.Add(30*time.Second)); d.Outcome != "allow" {
		t.Fatalf("fresh snapshot decision = %+v", d)
	}
	v.Payload.Identities[0].ExternalID = "changed-identity"
	if d := v.EvaluateV2Plan("github.review_pr", "agent-a", subject, ops, now); d.ReasonCode != "subject_unlinked" {
		t.Fatalf("substituted identity decision = %+v", d)
	}
}

// mcpTestPolicy describes one harness_policy entry for MCP test bundles.
type mcpTestPolicy struct {
	id       string
	src      string
	dst      string
	action   string
	priority int
	agentID  string
}

func v2MCPBundle(t *testing.T, now time.Time, policies []mcpTestPolicy) *Verified {
	t.Helper()
	doc := v2BundleDoc(now)
	set := doc["policy_set"].(map[string]any)
	var hps []any
	for _, p := range policies {
		var agent any
		if p.agentID != "" {
			agent = p.agentID
		}
		hps = append(hps, map[string]any{
			"policy_id": p.id, "name": "mcp-" + p.id, "src_pattern": p.src, "dst_pattern": p.dst, "action": p.action, "priority": p.priority,
			"scope": map[string]any{"org_id": v2Org, "workspace_id": v2Workspace, "agent_id": agent},
		})
	}
	set["harness_policies"] = hps
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var payload Payload
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	return &Verified{Payload: payload, PayloadBytes: raw, Digest: "sha256:mcp-test"}
}

func TestEvaluateV2MCP(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}

	t.Run("permit mcp:<slug>.<tool>", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "mcp:slack.list_channels", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("permit mcp:<slug>", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "mcp:slack", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("deny wins at equal priority", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-permit", src: "*", dst: "mcp:slack", action: "permit", priority: 10},
			{id: "p-deny", src: "*", dst: "mcp:slack.list_channels", action: "deny", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "deny" || d.PolicyID != "p-deny" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("group source via membership", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "group:members", dst: "mcp:slack.list_channels", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("no match denies", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "mcp:other_server.read", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("dotted tool name", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "mcp:api.github.com.list_repos", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("api.github.com", "list_repos", "agent-a", "claude_code", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("harness kind source", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-codex", src: "harness:codex", dst: "mcp:slack", action: "permit", priority: 10},
		})
		d := v.EvaluateV2MCP("slack", "list_channels", "agent-a", "codex", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p-codex" {
			t.Fatalf("codex session = %+v, want allow via p-codex", d)
		}
		d = v.EvaluateV2MCP("slack", "list_channels", "agent-a", "claude_code", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("claude_code session = %+v, want deny no_matching_policy", d)
		}
	})
}

// TestEvaluateV2Shell exercises the shell-command path of the v2 harness
// evaluation. It shares the v2MCPBundle helper because both paths consume the
// same harness_policies wire shape.
func TestEvaluateV2Shell(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}

	t.Run("permit shell:<argv prefix>", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "shell:git status", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status", "--short"}, "agent-a", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("deny wins at equal priority", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-permit", src: "*", dst: "shell:git status", action: "permit", priority: 10},
			{id: "p-deny", src: "*", dst: "shell:git", action: "deny", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "deny" || d.PolicyID != "p-deny" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("group source via membership", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "group:members", dst: "shell:git status", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("no match denies", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "shell:other command", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("shell:* permit is not a shell grant", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "shell:*", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("registered tool grant is not a shell grant", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "tool:pi.bash", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("harness kind source", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-pi", src: "harness:pi", dst: "shell:git status", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "allow" || d.PolicyID != "p-pi" {
			t.Fatalf("decision = %+v", d)
		}
		v = v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p-codex", src: "harness:codex", dst: "shell:git status", action: "permit", priority: 10},
		})
		d = v.EvaluateV2Shell([]string{"git", "status"}, "agent-a", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("decision = %+v", d)
		}
	})

	t.Run("empty agent id denies", func(t *testing.T) {
		v := v2MCPBundle(t, now, []mcpTestPolicy{
			{id: "p1", src: "*", dst: "shell:git status", action: "permit", priority: 10},
		})
		d := v.EvaluateV2Shell([]string{"git", "status"}, "", subject, now)
		if d.Outcome != "deny" || d.ReasonCode != "agent_not_assigned" {
			t.Fatalf("decision = %+v", d)
		}
	})
}

func TestV2SessionPathRequiresPreflight(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	v := v2ToolBundle(t, now, nil)
	d := v.EvaluateSessionAt(Call{Tool: "github.review_pr"}, &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: SubjectLinked}, nil, now)
	if d.Outcome != "deny" || d.ReasonCode != "preflight_required" {
		t.Fatalf("v2 regular authorize path = %+v", d)
	}
}

// scope.agent_id is deprecated: producers omit it (a policy applies to every
// agent in its workspace; agents get access through group membership).
// Consumers keep honoring it when a producer still sends it, so older bundles
// evaluate as before.
func TestV2ScopeAgentIDIsOptionalAndStillHonored(t *testing.T) {
	type scope struct {
		OrgID       *string `json:"org_id"`
		WorkspaceID *string `json:"workspace_id"`
		AgentID     *string `json:"agent_id"`
	}
	decode := func(raw string) scope {
		t.Helper()
		var s scope
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	for name, raw := range map[string]string{
		"key omitted": `{"org_id":"org-v2","workspace_id":"ws-v2"}`,
		"null":        `{"org_id":"org-v2","workspace_id":"ws-v2","agent_id":null}`,
	} {
		if !v2ScopeMatches(decode(raw), "any-agent", "org-v2", "ws-v2") {
			t.Errorf("%s: a scope without agent_id must match every agent", name)
		}
	}
	scoped := decode(`{"org_id":"org-v2","workspace_id":"ws-v2","agent_id":"agent-a"}`)
	if !v2ScopeMatches(scoped, "agent-a", "org-v2", "ws-v2") || v2ScopeMatches(scoped, "agent-b", "org-v2", "ws-v2") {
		t.Error("a scope that still carries agent_id must keep matching only that agent")
	}
}

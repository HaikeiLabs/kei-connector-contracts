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
	"strings"
	"testing"
	"time"
)

// whenDoc builds a v2 bundle whose harness_policies carry the given rows, with
// deny-overrides declared. Each row is {id, action, when} where when may be nil.
func whenDoc(now time.Time, rows ...map[string]any) map[string]any {
	doc := v2BundleDoc(now)
	doc["harness_match_semantics"] = "kei.harness-match/v1"
	doc["harnesses"] = []any{map[string]any{"agent_id": "agent-1", "kind": "viapdp"}}
	doc["harness_kinds"] = []any{pdpKind()}
	set := doc["policy_set"].(map[string]any)
	set["evaluation"] = map[string]any{"precedence": []string{"deny_overrides"}, "checks": "per_capability_first_match_and", "default": "deny"}
	hps := []any{}
	for _, r := range rows {
		hp := map[string]any{
			"policy_id": r["id"], "name": r["id"], "src_pattern": "*", "dst_pattern": r["dst"], "action": r["action"], "priority": 0,
			"scope": map[string]any{"org_id": nil, "workspace_id": nil, "agent_id": nil},
		}
		if r["when"] != nil {
			hp["when"] = r["when"]
		}
		hps = append(hps, hp)
	}
	set["harness_policies"] = hps
	return doc
}

func whenValidate(t *testing.T, doc map[string]any, now time.Time) (*Verified, error) {
	t.Helper()
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	raw, _ := json.Marshal(doc)
	return ValidateV2(raw, v2Installation, v2Org, v2Workspace, 1, now)
}

var whenSubject = &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: "linked"}

// ADR-038 D-6: a harness policy may carry when {harness: kind}. The closed
// parser must accept exactly that shape and reject anything else.
func TestValidateV2AcceptsAWhenHarnessCondition(t *testing.T) {
	now := time.Now().UTC()
	doc := whenDoc(now, map[string]any{"id": "p1", "action": "permit", "dst": "shell:git status", "when": map[string]any{"harness": "viapdp"}})
	if _, err := whenValidate(t, doc, now); err != nil {
		t.Fatalf("a harness policy with when {harness} must validate, got %v", err)
	}
}

func TestValidateV2RejectsAMalformedWhen(t *testing.T) {
	now := time.Now().UTC()
	cases := map[string]any{
		"an unknown condition key":   map[string]any{"harness": "viapdp", "region": "eu"},
		"an empty harness":           map[string]any{"harness": ""},
		"an empty condition object":  map[string]any{},
		"a non-string harness":       map[string]any{"harness": 7},
		"a condition that is a list": []any{"viapdp"},
	}
	for name, when := range cases {
		t.Run(name, func(t *testing.T) {
			doc := whenDoc(now, map[string]any{"id": "p1", "action": "permit", "dst": "shell:git status", "when": when})
			if _, err := whenValidate(t, doc, now); err == nil {
				t.Fatal("must be rejected")
			}
		})
	}
}

func TestEvaluateV2ShellHonorsWhenHarness(t *testing.T) {
	now := time.Now().UTC()
	argv := []string{"git", "status"}
	onlyViapdp := map[string]any{"harness": "viapdp"}

	t.Run("a permit applies to the named harness only", func(t *testing.T) {
		v, err := whenValidate(t, whenDoc(now, map[string]any{"id": "p1", "action": "permit", "dst": "shell:git status", "when": onlyViapdp}), now)
		if err != nil {
			t.Fatal(err)
		}
		if d := v.EvaluateV2ShellForKind("viapdp", argv, "agent-1", whenSubject, now); d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("named harness = %+v, want allow by p1", d)
		}
		if d := v.EvaluateV2ShellForKind("pi", argv, "agent-1", whenSubject, now); d.Outcome != "deny" || d.ReasonCode != "no_matching_policy" {
			t.Fatalf("other harness = %+v, want no_matching_policy", d)
		}
	})

	t.Run("a deny with a matching when still beats a permit", func(t *testing.T) {
		v, err := whenValidate(t, whenDoc(now,
			map[string]any{"id": "p1", "action": "permit", "dst": "shell:git"},
			map[string]any{"id": "p2", "action": "deny", "dst": "shell:git status", "when": onlyViapdp},
		), now)
		if err != nil {
			t.Fatal(err)
		}
		if d := v.EvaluateV2ShellForKind("viapdp", argv, "agent-1", whenSubject, now); d.Outcome != "deny" || d.PolicyID != "p2" {
			t.Fatalf("named harness = %+v, want deny by p2", d)
		}
	})

	t.Run("a deny whose when does not match cannot deny", func(t *testing.T) {
		v, err := whenValidate(t, whenDoc(now,
			map[string]any{"id": "p1", "action": "permit", "dst": "shell:git"},
			map[string]any{"id": "p2", "action": "deny", "dst": "shell:git status", "when": onlyViapdp},
		), now)
		if err != nil {
			t.Fatal(err)
		}
		if d := v.EvaluateV2ShellForKind("pi", argv, "agent-1", whenSubject, now); d.Outcome != "allow" || d.PolicyID != "p1" {
			t.Fatalf("other harness = %+v, want allow by p1", d)
		}
	})
}

func TestEvaluateV2MCPHonorsWhenHarness(t *testing.T) {
	now := time.Now().UTC()
	doc := whenDoc(now,
		map[string]any{"id": "p1", "action": "permit", "dst": "mcp:docs", "when": map[string]any{"harness": "viapdp"}},
	)
	v, err := whenValidate(t, doc, now)
	if err != nil {
		t.Fatal(err)
	}
	if d := v.EvaluateV2MCP("docs", "search", "agent-1", "viapdp", whenSubject, now); d.Outcome != "allow" {
		t.Fatalf("named harness = %+v, want allow", d)
	}
	if d := v.EvaluateV2MCP("docs", "search", "agent-1", "pi", whenSubject, now); d.Outcome == "allow" {
		t.Fatalf("other harness = %+v, want not allowed", d)
	}
}

func TestEvaluateHarnessHonorsWhenHarness(t *testing.T) {
	now := time.Now().UTC()
	build := func(when any) *Verified {
		row := map[string]any{"id": "p1", "action": "permit", "dst": "shell:git status"}
		if when != nil {
			row["when"] = when
		}
		doc := whenDoc(now, row)
		delete(doc, "harness_kinds") // built-in kinds, as in the deny-overrides harness test
		raw, _ := json.Marshal(doc)
		var payload Payload
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatal(err)
		}
		return &Verified{Payload: payload, PayloadBytes: raw, Digest: "sha256:when-test"}
	}
	call := func(kind string) HarnessCall { return HarnessCall{Kind: kind, Argv: []string{"git", "status"}} }

	if outcome, id, _ := build(nil).EvaluateHarness(call("claude_code")); outcome != "permit" || id != "p1" {
		t.Fatalf("control, no when: outcome=%s id=%s, want permit by p1", outcome, id)
	}
	scoped := build(map[string]any{"harness": "claude_code"})
	if outcome, id, _ := scoped.EvaluateHarness(call("claude_code")); outcome != "permit" || id != "p1" {
		t.Fatalf("named harness: outcome=%s id=%s, want permit by p1", outcome, id)
	}
	if outcome, id, _ := scoped.EvaluateHarness(call("codex")); id == "p1" || outcome == "permit" {
		t.Fatalf("other harness: outcome=%s id=%s, want the claude_code-only policy not to apply", outcome, id)
	}
}

func TestV2WhenWireShapeIsJustTheHarness(t *testing.T) {
	now := time.Now().UTC()
	doc := whenDoc(now, map[string]any{"id": "p1", "action": "permit", "dst": "shell:git", "when": map[string]any{"harness": "viapdp"}})
	raw, _ := json.Marshal(doc)
	if !strings.Contains(string(raw), `"when":{"harness":"viapdp"}`) {
		t.Fatalf("wire shape = %s", raw)
	}
}

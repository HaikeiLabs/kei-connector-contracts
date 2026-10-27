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

func harnessKindsDoc(now time.Time, kinds any, harnessKind string) []byte {
	doc := v2BundleDoc(now)
	doc["harness_match_semantics"] = "kei.harness-match/v1"
	doc["harnesses"] = []any{map[string]any{"agent_id": "agent-1", "kind": harnessKind}}
	if kinds != nil {
		doc["harness_kinds"] = kinds
	}
	raw, _ := json.Marshal(doc)
	return raw
}

func exampleKind() map[string]any {
	return map[string]any{
		"kind": "example", "short_name": "ex", "tool_prefix": "example.", "shell_tool": "example.shell", "gateable": true,
		"native_tools": []any{map[string]any{"tool_name": "example.read", "native_name": "Read"}},
	}
}

func pdpKind() map[string]any {
	return map[string]any{"kind": "viapdp", "tool_prefix": "viapdp.", "shell_tool": "viapdp.shell", "gateable": true, "decides_via": "pdp"}
}

func TestBundleHarnessKindsExtendTheBuiltInKinds(t *testing.T) {
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	now := time.Now().UTC()

	v, err := ValidateV2(harnessKindsDoc(now, []any{exampleKind()}, "example"), v2Installation, v2Org, v2Workspace, 1, now)
	if err != nil {
		t.Fatalf("kind named by the bundle: %v", err)
	}
	if k := v.Payload.HarnessKinds; len(k) != 1 || k[0].ShellTool != "example.shell" || k[0].NativeTools[0].NativeName != "Read" || !k[0].Gateable {
		t.Fatalf("HarnessKinds = %+v", k)
	}
	if !v.Payload.KnownHarnessKind("example") || v.Payload.KnownHarnessKind("other") {
		t.Fatal("KnownHarnessKind must follow the bundle's own list")
	}

	// Without the field the built-in kinds still validate and nothing else does.
	if _, err := ValidateV2(harnessKindsDoc(now, nil, "codex"), v2Installation, v2Org, v2Workspace, 1, now); err != nil {
		t.Fatalf("built-in kind with no registry: %v", err)
	}
	if _, err := ValidateV2(harnessKindsDoc(now, nil, "example"), v2Installation, v2Org, v2Workspace, 1, now); err == nil {
		t.Fatal("an unlisted kind must be rejected when the bundle names none")
	}
	// With a registry, a kind it does not list is rejected.
	if _, err := ValidateV2(harnessKindsDoc(now, []any{exampleKind()}, "codex"), v2Installation, v2Org, v2Workspace, 1, now); err == nil {
		t.Fatal("a kind missing from the bundle's registry must be rejected")
	}
}

func TestBundleHarnessKindsAreValidated(t *testing.T) {
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	now := time.Now().UTC()
	with := func(edit func(map[string]any)) []any {
		k := exampleKind()
		edit(k)
		return []any{k}
	}
	for name, kinds := range map[string][]any{
		"bad kind": with(func(k map[string]any) { k["kind"] = "Bad Kind" }),
		"tool outside prefix": with(func(k map[string]any) {
			k["native_tools"] = []any{map[string]any{"tool_name": "other.read", "native_name": "Read"}}
		}),
		"shell tool off prefix": with(func(k map[string]any) { k["shell_tool"] = "other.shell" }),
		"empty native name": with(func(k map[string]any) {
			k["native_tools"] = []any{map[string]any{"tool_name": "example.read", "native_name": ""}}
		}),
		"duplicate kind": {exampleKind(), exampleKind()},
	} {
		_, err := ValidateV2(harnessKindsDoc(now, kinds, "example"), v2Installation, v2Org, v2Workspace, 1, now)
		if err == nil || !strings.Contains(err.Error(), "harness_kinds") {
			t.Errorf("%s: err = %v, want a harness_kinds rejection", name, err)
		}
	}
}

func TestEvaluateHarnessUsesTheBundleKinds(t *testing.T) {
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	now := time.Now().UTC()
	v, err := ValidateV2(harnessKindsDoc(now, []any{exampleKind()}, "example"), v2Installation, v2Org, v2Workspace, 1, now)
	if err != nil {
		t.Fatal(err)
	}
	outcome, _, reason := v.EvaluateHarness(HarnessCall{Kind: "example", Tool: "example.read"})
	if reason == HarnessReasonUnsupportedKind {
		t.Fatalf("outcome %s reason %s: a kind the bundle lists must not be unsupported", outcome, reason)
	}
	if _, _, reason := v.EvaluateHarness(HarnessCall{Kind: "unlisted", Tool: "x"}); reason != HarnessReasonUnsupportedKind {
		t.Fatalf("an unlisted kind: reason %q, want %q", reason, HarnessReasonUnsupportedKind)
	}
}

func TestHarnessKindDecidesVia(t *testing.T) {
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	now := time.Now().UTC()
	native := exampleKind()
	native["decides_via"] = "native"
	v, err := ValidateV2(harnessKindsDoc(now, []any{native, pdpKind()}, "viapdp"), v2Installation, v2Org, v2Workspace, 1, now)
	if err != nil {
		t.Fatal(err)
	}
	if v.Payload.HarnessKindDecidesViaPDP("example") || !v.Payload.HarnessKindDecidesViaPDP("viapdp") {
		t.Fatal("decides_via must follow the registry entry")
	}
	if _, reason := evalReason(v, "viapdp"); reason != HarnessReasonEvaluateViaPDP {
		t.Fatalf("a pdp kind: reason %q, want %q", reason, HarnessReasonEvaluateViaPDP)
	}
	if outcome, _ := evalReason(v, "example"); outcome != HarnessUnmatched {
		t.Fatalf("a native kind with no policy: outcome %q, want unmatched", outcome)
	}
	bad := exampleKind()
	bad["decides_via"] = "telepathy"
	if _, err := ValidateV2(harnessKindsDoc(now, []any{bad}, "example"), v2Installation, v2Org, v2Workspace, 1, now); err == nil || !strings.Contains(err.Error(), "decides_via") {
		t.Fatalf("an unknown decides_via must be rejected, got %v", err)
	}
	// With no registry entry value, the built-in kinds keep their meaning.
	var none Payload
	if !none.HarnessKindDecidesViaPDP("pi") || !none.HarnessKindDecidesViaPDP("custom") || none.HarnessKindDecidesViaPDP("codex") {
		t.Fatal("built-in decides_via: pi and custom via the PDP, codex native")
	}
}

func evalReason(v *Verified, kind string) (outcome, reason string) {
	outcome, _, reason = v.EvaluateHarness(HarnessCall{Kind: kind, Tool: kind + ".read"})
	return
}

func TestShellToolNamesItsKind(t *testing.T) {
	var legacy Payload
	if k, ok := legacy.HarnessKindForShellTool("pi.bash"); !ok || k.Kind != "pi" {
		t.Fatalf("with no registry pi.bash is the legacy Pi shell tool, got %+v ok=%v", k, ok)
	}
	if k, ok := legacy.HarnessKindForShellTool("viapdp.shell"); ok || k.Kind != "" {
		t.Fatalf("no match must return the zero kind, got %+v ok=%v", k, ok)
	}
	if _, ok := legacy.HarnessKindForShellTool("claude_code.bash"); ok {
		t.Fatal("a native kind's shell tool is not decided by the PDP")
	}
	p := Payload{HarnessKinds: []HarnessKind{
		{Kind: "viapdp", ShellTool: "viapdp.shell", DecidesVia: "pdp"},
		{Kind: "example", ShellTool: "example.shell", DecidesVia: "native"},
	}}
	if k, ok := p.HarnessKindForShellTool("viapdp.shell"); !ok || k.Kind != "viapdp" {
		t.Fatalf("registry shell tool = %+v ok=%v", k, ok)
	}
	if _, ok := p.HarnessKindForShellTool("example.shell"); ok {
		t.Fatal("a native kind's shell tool is not decided by the PDP")
	}
	if _, ok := p.HarnessKindForShellTool("pi.bash"); ok {
		t.Fatal("with a registry only its own shell tools count")
	}
}

func TestEvaluateV2ShellForKindAttributesTheCallToItsKind(t *testing.T) {
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	now := time.Now().UTC()
	doc := v2BundleDoc(now)
	doc["harness_match_semantics"] = "kei.harness-match/v1"
	doc["harnesses"] = []any{map[string]any{"agent_id": "agent-1", "kind": "viapdp"}}
	doc["harness_kinds"] = []any{pdpKind()}
	set := doc["policy_set"].(map[string]any)
	set["harness_policies"] = []any{map[string]any{
		"policy_id": "p1", "name": "p1", "src_pattern": "harness:viapdp", "dst_pattern": "shell:git status", "action": "permit", "priority": 10,
		"scope": map[string]any{"org_id": nil, "workspace_id": nil, "agent_id": nil},
	}}
	raw, _ := json.Marshal(doc)
	v, err := ValidateV2(raw, v2Installation, v2Org, v2Workspace, 1, now)
	if err != nil {
		t.Fatal(err)
	}
	subject := &SessionSubject{UserID: "user-alice", Source: "discord", ExternalID: "ext-alice", Status: "linked"}
	argv := []string{"git", "status"}
	if d := v.EvaluateV2ShellForKind("viapdp", argv, "agent-1", subject, now); d.Outcome != "allow" || d.PolicyID != "p1" {
		t.Fatalf("shell for its own kind = %+v, want allow by p1", d)
	}
	if d := v.EvaluateV2ShellForKind("pi", argv, "agent-1", subject, now); d.Outcome == "allow" {
		t.Fatalf("a harness:viapdp policy must not match a pi call, got %+v", d)
	}
	if d := v.EvaluateV2PiShell(argv, "agent-1", subject, now); d.Outcome == "allow" {
		t.Fatalf("EvaluateV2PiShell is the pi kind, got %+v", d)
	}
}

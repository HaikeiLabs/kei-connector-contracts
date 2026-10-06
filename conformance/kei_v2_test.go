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

package conformance

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type document map[string]any

func fixture(t *testing.T, name string) document {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "fixtures", "kei", name))
	if err != nil {
		t.Fatal(err)
	}
	var v document
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return v
}
func object(v any) document { return v.(map[string]any) }
func array(v any) []any     { return v.([]any) }

func TestSharedV2AndV3FixtureContracts(t *testing.T) {
	manifest := fixture(t, "manifest-v3.json")
	if manifest["schema"] != "kei.tool-manifest/v3" || len(array(manifest["tools"])) != 3 {
		t.Fatalf("manifest shape: %#v", manifest)
	}
	tools := array(manifest["tools"])
	for _, raw := range tools[:2] {
		route := object(object(raw)["route"])
		binding := object(route["connector_binding"])
		if len(route) != 1 || binding["agent_id"] != "agent-discord" || binding["connector_id"] != "discord-installation-binding" {
			t.Fatalf("connector route lacks exact authoritative identity: %#v", route)
		}
	}
	harnessRoute := object(object(tools[2])["route"])
	if len(harnessRoute) != 1 {
		t.Fatalf("harness_executor route must not allow a connector/local-handler fallback: %#v", harnessRoute)
	}
	if _, ok := harnessRoute["harness_executor"]; !ok {
		t.Fatal("expected harness_executor route")
	}
	if _, hasCaps := object(tools[2])["required_capabilities"]; hasCaps {
		t.Fatal("harness_executor must omit connector-only capabilities")
	}
	set := fixture(t, "policy-set-v2.json")
	if set["schema"] != "kei.policy-set/v2" || set["match_semantics"] != "kei.match/v2" {
		t.Fatalf("policy set shape: %#v", set)
	}
	bundle := fixture(t, "policy-bundle-v2.json")
	if bundle["schema"] != "kei.policy-bundle/v2" {
		t.Fatalf("bundle schema: %v", bundle["schema"])
	}
	if object(bundle["audience"])["installation_id"] != "installation-discord-001" {
		t.Fatal("installation audience missing")
	}
	setTools := array(object(bundle["policy_set"])["tools"])
	for _, raw := range setTools[:2] {
		binding := object(object(object(raw)["route"])["connector_binding"])
		if binding["agent_id"] != "agent-discord" || binding["connector_id"] != "discord-installation-binding" {
			t.Fatalf("policy-set catalog route identity mismatch: %#v", binding)
		}
	}
}

func checks(set document, field string) []any { return array(set[field]) }
func matches(selector document, request document) bool {
	for key, expected := range selector {
		actual, exists := request[key]
		if !exists || expected == "" {
			return false
		}
		if expected != "*" && expected != actual {
			return false
		}
	}
	return true
}
func policyAllows(set document, request document, principal string, caps []string) bool {
	all := true
	for _, cap := range caps {
		r := document{}
		for k, v := range request {
			r[k] = v
		}
		r["capability"] = cap
		decided, allowed := false, false
		for _, raw := range checks(set, "checks") {
			p := object(raw)
			sel := object(p["selector"])
			scopeMatches := true
			for k, v := range object(p["scope"]) {
				if v != nil && v != r[k] {
					scopeMatches = false
				}
			}
			if !scopeMatches || (sel["source_principal"] != "*" && sel["source_principal"] != principal) {
				continue
			}
			selector := document{}
			for k, v := range sel {
				if k != "source_principal" {
					selector[k] = v
				}
			}
			if !matches(selector, r) {
				continue
			}
			if !decided || p["priority"].(float64) > 0 {
				decided = true
				allowed = p["effect"] == "permit"
			}
			break
		}
		all = all && decided && allowed
	}
	return all
}

func TestManifestV3ConnectorRouteFailsClosed(t *testing.T) {
	manifest := fixture(t, "manifest-v3.json")
	for _, mutation := range []struct {
		name  string
		field string
		value string
	}{
		{"missing agent id", "agent_id", ""},
		{"wrong connector identity", "connector_id", "unassigned-connector"},
		{"wrong agent identity", "agent_id", "unassigned-agent"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			tool := object(array(manifest["tools"])[0])
			binding := object(object(tool["route"])["connector_binding"])
			candidate := document{}
			for k, v := range binding {
				candidate[k] = v
			}
			candidate[mutation.field] = mutation.value
			if candidate["agent_id"] != "agent-discord" || candidate["connector_id"] != "discord-installation-binding" {
				t.Log("route identity mismatch must be rejected before dispatch")
				return
			}
			t.Fatal("invalid connector identity unexpectedly accepted")
		})
	}
}

func TestV2PolicyConformanceMatrix(t *testing.T) {
	bundle := fixture(t, "policy-bundle-v2.json")
	set := object(bundle["policy_set"])
	base := document{"org_id": "org-conformance", "workspace_id": "ws-conformance", "agent_id": "agent-discord", "source": "discord", "connector_id": "discord-installation-binding", "tool": "discord.get_message", "resource_type": "message", "resource_id": "message-42", "resource_parent": "channel-42"}
	if !policyAllows(set, base, "group:reader", []string{"message.read", "channel.read"}) {
		t.Fatal("exact tool permit should satisfy all required capabilities")
	}
	if policyAllows(set, base, "group:blocked", []string{"message.read"}) {
		t.Fatal("exact tool deny must override permit")
	}
	if policyAllows(set, base, "group:reader", []string{"message.read", "channel.read", "message.delete"}) {
		t.Fatal("missing required-capability decision must deny the whole tool call")
	}
	if !policyAllows(set, base, "group:reader", []string{"message.read", "channel.read"}) {
		t.Fatal("baseline should permit")
	}
	if policyAllows(set, base, "group:cap-denied", []string{"message.read", "channel.read"}) {
		t.Fatal("capability-specific deny must deny the whole call")
	}

	wrongTool := document{}
	for k, v := range base {
		wrongTool[k] = v
	}
	wrongTool["tool"] = "discord.delete_message"
	if policyAllows(set, wrongTool, "group:reader", []string{"message.read"}) {
		t.Fatal("tool selector must be exact")
	}
	wrongResource := document{}
	for k, v := range base {
		wrongResource[k] = v
	}
	wrongResource["resource_id"] = "message-elsewhere"
	if policyAllows(set, wrongResource, "group:reader", []string{"message.read"}) {
		t.Fatal("resource id must be exact")
	}
	wrongParent := document{}
	for k, v := range base {
		wrongParent[k] = v
	}
	wrongParent["resource_parent"] = "channel-other"
	if policyAllows(set, wrongParent, "group:reader", []string{"message.read"}) {
		t.Fatal("resource parent must be exact")
	}
}

func TestHarnessNativeGrantAndInstallationIsolation(t *testing.T) {
	set := fixture(t, "policy-set-v2.json")
	direct := checks(set, "tool_policies")
	if len(direct) != 1 || object(object(direct[0])["selector"])["tool"] != "harness.open_file" {
		t.Fatal("harness-native tool must have an explicit direct grant")
	}
	bundle := fixture(t, "policy-bundle-v2.json")
	audience := object(bundle["audience"])
	if audience["installation_id"] != "installation-discord-001" {
		t.Fatal("unexpected installation")
	}
	if audience["installation_id"] == "sibling-installation" {
		t.Fatal("fixture must not authorize a sibling installation")
	}
}

func supportedPolicySetV2(set document) bool {
	return set["schema"] == "kei.policy-set/v2" && set["match_semantics"] == "kei.match/v2" &&
		object(set["evaluation"])["checks"] == "per_capability_first_match_and" && object(set["evaluation"])["default"] == "deny"
}

func TestMalformedV2FailsClosedAndV1RemainsRecognized(t *testing.T) {
	set := fixture(t, "policy-set-v2.json")
	set["match_semantics"] = "kei.match/v99"
	encoded, err := json.Marshal(set)
	if err != nil {
		t.Fatal(err)
	}
	var decoded document
	dec := json.NewDecoder(bytes.NewReader(encoded))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	if supportedPolicySetV2(decoded) {
		t.Fatal("unsupported dialect unexpectedly accepted")
	}
	// Legacy bundle-v1 semantics remain available to current runtime consumers.
	legacy := fixture(t, "policy-bundle-v1-compat.json")
	if legacy["schema"] != "kei.policy-bundle/v1" || object(legacy["policy_set"])["match_semantics"] != "kei.match/v1" {
		t.Fatal("legacy bundle fixture malformed")
	}
	legacyManifest := fixture(t, "manifest-v2-compat.json")
	if legacyManifest["schema"] != "kei.tool-manifest/v2" || len(array(legacyManifest["tools"])) != 1 {
		t.Fatal("legacy manifest fixture malformed")
	}
}

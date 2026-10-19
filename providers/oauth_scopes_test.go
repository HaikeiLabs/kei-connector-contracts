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

package providers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/policy"
)

// fixtureMapping is a representative capability→scope mapping used by tests
// that do not need to supply their own. It mirrors the Google Drive subset of
// CapabilityOAuthScopes.
var fixtureMapping = map[string][]string{
	"drive.search":        {"https://www.googleapis.com/auth/drive.readonly"},
	"drive.metadata.read": {"https://www.googleapis.com/auth/drive.readonly"},
	"docs.read":           {"https://www.googleapis.com/auth/drive.readonly"},
	"drive.doc.create":    {"https://www.googleapis.com/auth/drive.file"},
	"drive.doc.update":    {"https://www.googleapis.com/auth/drive.file"},
}

func TestRequiredOAuthScopes(t *testing.T) {
	for _, tc := range []struct {
		name         string
		mapping      map[string][]string
		capabilities []contract.Capability
		want         []string
		wantErr      string // empty means success
	}{
		{
			name:         "empty capabilities return empty scopes",
			mapping:      fixtureMapping,
			capabilities: nil,
			want:         []string{},
		},
		{
			name:         "empty capabilities with nil mapping",
			mapping:      nil,
			capabilities: nil,
			want:         []string{},
		},
		{
			name:         "single capability resolves to one scope",
			mapping:      fixtureMapping,
			capabilities: []contract.Capability{{Name: "drive.search"}},
			want:         []string{"https://www.googleapis.com/auth/drive.readonly"},
		},
		{
			name:    "union of scopes from multiple capabilities",
			mapping: fixtureMapping,
			capabilities: []contract.Capability{
				{Name: "drive.search"},
				{Name: "drive.doc.create"},
			},
			want: []string{
				"https://www.googleapis.com/auth/drive.file",
				"https://www.googleapis.com/auth/drive.readonly",
			},
		},
		{
			name:    "duplicate scopes are deduplicated",
			mapping: fixtureMapping,
			capabilities: []contract.Capability{
				{Name: "drive.search"},
				{Name: "drive.metadata.read"},
				{Name: "docs.read"},
			},
			want: []string{"https://www.googleapis.com/auth/drive.readonly"},
		},
		{
			name:    "unknown capability fails closed",
			mapping: fixtureMapping,
			capabilities: []contract.Capability{
				{Name: "admin.raw"},
			},
			wantErr: `OAuth scopes are not declared for capability "admin.raw"`,
		},
		{
			name:    "unknown capability among known caps fails closed",
			mapping: fixtureMapping,
			capabilities: []contract.Capability{
				{Name: "drive.search"},
				{Name: "admin.raw"},
			},
			wantErr: `OAuth scopes are not declared for capability "admin.raw"`,
		},
		{
			name:    "non-empty capabilities with nil mapping fails closed",
			mapping: nil,
			capabilities: []contract.Capability{
				{Name: "drive.search"},
			},
			wantErr: "OAuth scope mappings are not declared",
		},
		{
			name:    "non-empty capabilities with empty mapping fails closed",
			mapping: map[string][]string{},
			capabilities: []contract.Capability{
				{Name: "drive.search"},
			},
			wantErr: "OAuth scope mappings are not declared",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RequiredOAuthScopes(tc.mapping, tc.capabilities)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("RequiredOAuthScopes = %v, wanted error containing %q", got, tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("RequiredOAuthScopes error = %q, want substring %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("RequiredOAuthScopes = %v, %v", got, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("RequiredOAuthScopes = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMissingOAuthScopes(t *testing.T) {
	for _, tc := range []struct {
		name     string
		required []string
		granted  []string
		want     []string
	}{
		{
			name:     "all granted",
			required: []string{"scope-a", "scope-b"},
			granted:  []string{"scope-a", "scope-b"},
			want:     []string{},
		},
		{
			name:     "some missing",
			required: []string{"scope-a", "scope-b", "scope-c"},
			granted:  []string{"scope-a"},
			want:     []string{"scope-b", "scope-c"},
		},
		{
			name:     "all missing",
			required: []string{"scope-a"},
			granted:  nil,
			want:     []string{"scope-a"},
		},
		{
			name:     "nothing required",
			required: nil,
			granted:  []string{"scope-a"},
			want:     []string{},
		},
		{
			name:     "both empty",
			required: nil,
			granted:  nil,
			want:     []string{},
		},
		{
			name:     "extra granted scopes ignored",
			required: []string{"scope-a"},
			granted:  []string{"scope-a", "scope-b", "scope-c"},
			want:     []string{},
		},
		{
			name:     "result sorted",
			required: []string{"z-scope", "a-scope", "m-scope"},
			granted:  []string{},
			want:     []string{"a-scope", "m-scope", "z-scope"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := MissingOAuthScopes(tc.required, tc.granted)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("MissingOAuthScopes = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValidateGrantedOAuthScopes(t *testing.T) {
	caps := []contract.Capability{{Name: "drive.search"}}
	mapping := fixtureMapping

	// Full grant passes.
	if err := ValidateGrantedOAuthScopes(mapping, caps, []string{"https://www.googleapis.com/auth/drive.readonly"}); err != nil {
		t.Fatalf("full grant rejected: %v", err)
	}
	// Empty grant fails with missing-scope error.
	if err := ValidateGrantedOAuthScopes(mapping, caps, nil); err == nil {
		t.Fatal("empty grant accepted")
	} else if !strings.Contains(err.Error(), "missing required scopes") {
		t.Fatalf("wrong error for missing grant: %v", err)
	}
	// Unknown capability fails with its own distinct error, not "missing".
	if err := ValidateGrantedOAuthScopes(mapping, []contract.Capability{{Name: "admin.raw"}}, []string{"anything"}); err == nil {
		t.Fatal("unknown capability accepted")
	} else if !strings.Contains(err.Error(), `not declared for capability "admin.raw"`) {
		t.Fatalf("wrong error for unknown capability: %v", err)
	}
	// Nil mapping with capabilities fails closed.
	if err := ValidateGrantedOAuthScopes(nil, caps, []string{"anything"}); err == nil {
		t.Fatal("nil mapping accepted")
	} else if !strings.Contains(err.Error(), "mappings are not declared") {
		t.Fatalf("wrong error for nil mapping: %v", err)
	}
}

func TestOAuthReconsentScopes(t *testing.T) {
	mapping := map[string][]string{
		"cap.a": {"scope-1"},
		"cap.b": {"scope-2"},
		"cap.c": {"scope-3"},
	}

	for _, tc := range []struct {
		name         string
		capabilities []contract.Capability
		granted      []string
		want         []string
		wantErr      string
	}{
		{
			name:         "no new scopes needed",
			capabilities: []contract.Capability{{Name: "cap.a"}, {Name: "cap.b"}},
			granted:      []string{"scope-1", "scope-2"},
			want:         []string{},
		},
		{
			name:         "new capability needs additional scope",
			capabilities: []contract.Capability{{Name: "cap.a"}, {Name: "cap.c"}},
			granted:      []string{"scope-1"},
			want:         []string{"scope-3"},
		},
		{
			name:         "multiple new scopes needed",
			capabilities: []contract.Capability{{Name: "cap.a"}, {Name: "cap.b"}, {Name: "cap.c"}},
			granted:      []string{"scope-1"},
			want:         []string{"scope-2", "scope-3"},
		},
		{
			name:         "nothing granted, everything needed",
			capabilities: []contract.Capability{{Name: "cap.a"}},
			granted:      nil,
			want:         []string{"scope-1"},
		},
		{
			name:         "unknown capability fails closed",
			capabilities: []contract.Capability{{Name: "cap.a"}, {Name: "unknown.cap"}},
			granted:      []string{"scope-1"},
			wantErr:      `not declared for capability "unknown.cap"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := OAuthReconsentScopes(mapping, tc.capabilities, tc.granted)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("OAuthReconsentScopes = %v, wanted error containing %q", got, tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("OAuthReconsentScopes error = %q, want substring %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("OAuthReconsentScopes = %v, %v", got, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("OAuthReconsentScopes = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOAuthReconsentScopesNilMapping(t *testing.T) {
	if _, err := OAuthReconsentScopes(nil, []contract.Capability{{Name: "cap.a"}}, []string{"scope-1"}); err == nil {
		t.Fatal("nil mapping accepted")
	} else if !strings.Contains(err.Error(), "mappings are not declared") {
		t.Fatalf("wrong error for nil mapping: %v", err)
	}
}

func TestGlobalCapabilityOAuthScopesSmoke(t *testing.T) {
	// Verify that the global map still resolves every capability it declares
	// (a smoke test for the deprecated global).
	for provider, mappings := range CapabilityOAuthScopes {
		for capability := range mappings {
			got, err := RequiredOAuthScopes(mappings, []contract.Capability{{Name: capability}})
			if err != nil {
				t.Errorf("%s/%s: RequiredOAuthScopes = %v, %v", provider, capability, got, err)
			}
			if len(got) == 0 {
				t.Errorf("%s/%s: resolved to zero scopes", provider, capability)
			}
		}
		// Unknown capability in the global map fails closed.
		if _, err := RequiredOAuthScopes(mappings, []contract.Capability{{Name: "nonexistent"}}); err == nil {
			t.Errorf("%s: unknown capability was resolved", provider)
		}
	}
}

func TestRequiredOAuthScopesLinearNarrowScopes(t *testing.T) {
	linearMappings := CapabilityOAuthScopes[contract.ProviderLinear]
	for _, tc := range []struct {
		capabilities []string
		want         []string
	}{
		{[]string{"team.read", "issue.read"}, []string{LinearReadScope}},
		{[]string{"issue.read", "issue.create"}, []string{LinearIssuesCreateScope, LinearReadScope}},
		{[]string{"comment.create"}, []string{LinearCommentsCreateScope}},
		{[]string{"issue.update"}, []string{LinearWriteScope}},
	} {
		caps := make([]contract.Capability, 0, len(tc.capabilities))
		for _, name := range tc.capabilities {
			caps = append(caps, contract.Capability{Name: name})
		}
		got, err := RequiredOAuthScopes(linearMappings, caps)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%v: scopes = %v, %v; want %v", tc.capabilities, got, err, tc.want)
		}
	}
	// Every defined Linear capability resolves.
	for _, capability := range contract.CapabilitiesFor(contract.ProviderLinear) {
		if _, err := RequiredOAuthScopes(linearMappings, []contract.Capability{capability}); err != nil {
			t.Errorf("Linear capability %q has no scope mapping: %v", capability.Name, err)
		}
	}
}

// DeclaresOAuthScopes is the interim preflight rule (owner decision 2026-10-10,
// option C): only a mapped provider whose connector scopes name its OAuth
// scopes is checked.
func TestDeclaresOAuthScopes(t *testing.T) {
	for _, tc := range []struct {
		name     string
		provider contract.Provider
		scopes   []string
		want     bool
	}{
		{"linear provider scopes", contract.ProviderLinear, []string{LinearReadScope}, true},
		{"gmail provider scope", contract.ProviderGmail, []string{GmailReadonlyScope}, true},
		{"capability names from kei-cli", contract.ProviderLinear, []string{"team.read", "issue.create"}, false},
		{"runtime session placeholder", contract.ProviderLinear, []string{"session"}, false},
		{"github has no map", contract.ProviderGitHub, []string{"repo"}, false},
	} {
		if got := DeclaresOAuthScopes(contract.Metadata{Provider: tc.provider, Scopes: tc.scopes}); got != tc.want {
			t.Errorf("%s: DeclaresOAuthScopes = %t, want %t", tc.name, got, tc.want)
		}
	}
}

// TestGuardSkipsPreflightForDescriptiveScopes is preserved from the existing
// test suite to document that the Guard preflight semantics are unchanged.
func TestGuardSkipsPreflightForDescriptiveScopes(t *testing.T) {
	for _, tc := range []struct {
		provider contract.Provider
		scopes   []string
		inv      contract.Invocation
	}{
		{contract.ProviderGitHub, []string{"read"}, invocation("repository.read", contract.ActionRead, "github/acme/repo")},
		{contract.ProviderLinear, []string{"session"}, invocation("team.read", contract.ActionRead, "linear/team/KEI")},
		{contract.ProviderLinear, []string{"team.read", "issue.create"}, invocation("team.read", contract.ActionRead, "linear/team/KEI")},
	} {
		meta := metaFor(t, tc.provider, []string{"*"}, nil)
		meta.CredentialSource = contract.CredentialSourceOAuth
		meta.AccountModel = contract.AccountModelPerUser
		meta.Scopes = tc.scopes
		if err := Guard(meta, tc.inv); err != nil {
			t.Errorf("%s scopes %v: Guard = %v, want no preflight", tc.provider, tc.scopes, err)
		}
	}
	meta := metaFor(t, contract.ProviderLinear, []string{"*"}, nil)
	meta.CredentialSource = contract.CredentialSourceOAuth
	meta.AccountModel = contract.AccountModelPerUser
	meta.Scopes = []string{LinearReadScope}
	if err := Guard(meta, invocation("team.read", contract.ActionRead, "linear/team/KEI")); err == nil {
		t.Error("provider scopes that miss issue.create/comment.create/issue.update passed the preflight")
	}
}

// TestDistinctErrors documents that unknown-capability and missing-scope errors
// are different, confirming the fail-closed contract.
func TestDistinctErrors(t *testing.T) {
	mapping := map[string][]string{"good.cap": {"scope-a"}}

	_, err := RequiredOAuthScopes(mapping, []contract.Capability{{Name: "bad.cap"}})
	if err == nil {
		t.Fatal("unknown capability did not error")
	}
	unknownMsg := err.Error()

	err = ValidateGrantedOAuthScopes(mapping, []contract.Capability{{Name: "good.cap"}}, nil)
	if err == nil {
		t.Fatal("missing grant did not error")
	}
	missingMsg := err.Error()

	if unknownMsg == missingMsg {
		t.Fatalf("unknown capability error %q matches missing grant error %q; they must be distinct", unknownMsg, missingMsg)
	}
}

// TestGuardOAuthPreflightOutcomes keeps re-consent, unknown mapping, and
// policy-deny outcomes distinct. Policy evaluation remains outside Guard.
func TestGuardOAuthPreflightOutcomes(t *testing.T) {
	mapping := map[string][]string{"good.cap": {"scope-a", "scope-b"}}
	metaWithOAuth := func(t *testing.T, scopes []string) contract.Metadata {
		t.Helper()
		m := contract.Metadata{
			ID:               "c-1",
			TenantID:         "t-1",
			WorkspaceID:      "w-1",
			Name:             "test",
			Provider:         contract.ProviderLinear,
			Status:           contract.StatusActive,
			CredentialRef:    "vault/tenant/t-1/test",
			CredentialSource: contract.CredentialSourceOAuth,
			AccountModel:     contract.AccountModelPerUser,
			Scopes:           scopes,
			Resources:        []string{"*"},
			Capabilities:     contract.CapabilitiesFor(contract.ProviderLinear),
			CreatedBy:        "u-1",
		}
		if err := m.Validate(); err != nil {
			t.Fatalf("fixture metadata is invalid: %v", err)
		}
		return m
	}
	inv := contract.Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: "c-1", Capability: "team.read", Action: contract.ActionRead, Resource: "test-resource", TraceID: "trace-1"}
	type result struct {
		err      error
		decision *policy.Outcome
	}
	cases := []struct {
		name  string
		run   func(*testing.T) result
		check func(*testing.T, result)
	}{
		{
			name: "missing scope is a typed re-consent error",
			run: func(t *testing.T) result {
				meta := metaWithOAuth(t, []string{LinearReadScope})
				return result{err: Guard(meta, inv)}
			},
			check: func(t *testing.T, got result) {
				var preflight *OAuthPreflightError
				if !errors.As(got.err, &preflight) {
					t.Fatalf("Guard error = %v (%T), want *OAuthPreflightError", got.err, got.err)
				}
				if !errors.Is(got.err, &OAuthPreflightError{}) {
					t.Fatalf("Guard error %v is not errors.Is-compatible with *OAuthPreflightError", got.err)
				}
				if len(preflight.MissingScopes) == 0 || !reflect.DeepEqual(preflight.MissingScopes, MissingOAuthScopes([]string{LinearReadScope, LinearIssuesCreateScope, LinearCommentsCreateScope, LinearWriteScope}, []string{LinearReadScope})) {
					t.Fatalf("MissingScopes = %v, want the bounded missing Linear scope list", preflight.MissingScopes)
				}
			},
		},
		{
			name: "unknown capability mapping is a distinct plain error",
			run: func(t *testing.T) result {
				return result{err: ValidateGrantedOAuthScopes(mapping, []contract.Capability{{Name: "unknown.cap"}}, []string{"anything"})}
			},
			check: func(t *testing.T, got result) {
				var preflight *OAuthPreflightError
				if got.err == nil || errors.As(got.err, &preflight) || !strings.Contains(got.err.Error(), `not declared for capability "unknown.cap"`) {
					t.Fatalf("unknown mapping error = %v (%T), want a distinct plain mapping error", got.err, got.err)
				}
			},
		},
		{
			name: "policy deny is a policy outcome outside Guard",
			run: func(t *testing.T) result {
				meta := metaWithOAuth(t, []string{LinearReadScope, LinearIssuesCreateScope, LinearCommentsCreateScope, LinearWriteScope})
				if err := Guard(meta, inv); err != nil {
					return result{err: err}
				}
				decision := policy.Decide([]policy.Rule{{ID: "deny-1", Principal: "u-1", Target: "*", Effect: policy.Effect("deny"), Enabled: true}}, policy.Request{User: "u-1", Target: inv.Resource})
				return result{decision: &decision}
			},
			check: func(t *testing.T, got result) {
				if got.err != nil || got.decision == nil || got.decision.Decision != "deny" {
					t.Fatalf("policy result = %+v, Guard error = %v; want policy.Outcome deny with Guard passing", got.decision, got.err)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { tc.check(t, tc.run(t)) })
	}
}

// ExampleRequiredOAuthScopes documents the expected usage with a supplied
// mapping from template data.
func ExampleRequiredOAuthScopes() {
	// Mapping supplied from template data required_oauth_scopes.
	mapping := map[string][]string{
		"drive.search": {"https://www.googleapis.com/auth/drive.readonly"},
	}
	scopes, err := RequiredOAuthScopes(mapping, []contract.Capability{{Name: "drive.search"}})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(scopes)
	// Output: [https://www.googleapis.com/auth/drive.readonly]
}

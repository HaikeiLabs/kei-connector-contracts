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
	"reflect"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func TestOAuthScopesLeastPrivilegeAndPreflight(t *testing.T) {
	caps := []contract.Capability{{Name: "drive.search", Action: contract.ActionRead}}
	required, err := RequiredOAuthScopes(contract.ProviderGoogle, caps)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"https://www.googleapis.com/auth/drive.readonly"}
	if !reflect.DeepEqual(required, want) {
		t.Fatalf("required = %v, want %v", required, want)
	}
	if err := ValidateGrantedOAuthScopes(contract.ProviderGoogle, caps, required); err != nil {
		t.Fatal(err)
	}
	if err := ValidateGrantedOAuthScopes(contract.ProviderGoogle, caps, nil); err == nil {
		t.Fatal("preflight accepted missing grant")
	}
}

func TestOAuthScopesIncrementalReconsentAndDeclaredCoverage(t *testing.T) {
	granted := []string{GmailReadonlyScope}
	for _, provider := range contract.DefinedProviders() {
		mappings, declared := CapabilityOAuthScopes[provider]
		if !declared {
			continue
		}
		for capability := range mappings {
			if _, err := RequiredOAuthScopes(provider, []contract.Capability{{Name: capability}}); err != nil {
				t.Errorf("%s/%s: %v", provider, capability, err)
			}
		}
	}
	extra, err := OAuthReconsentScopes(contract.ProviderGmail, []contract.Capability{{Name: "message.search"}, {Name: "message.get"}}, granted)
	if err != nil {
		t.Fatal(err)
	}
	if len(extra) != 0 {
		t.Fatalf("re-consent asked for already granted scopes: %v", extra)
	}
	if _, err := RequiredOAuthScopes(contract.ProviderGoogle, []contract.Capability{{Name: "admin.raw"}}); err == nil {
		t.Fatal("unknown capability inferred a scope")
	}
	if _, err := RequiredOAuthScopes(contract.ProviderGitHub, []contract.Capability{{Name: "repository.read"}}); err == nil {
		t.Fatal("undeclared provider mapping accepted")
	}
}

func TestLinearOAuthScopesAreNarrowPerWrite(t *testing.T) {
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
		got, err := RequiredOAuthScopes(contract.ProviderLinear, caps)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%v: scopes = %v, %v; want %v", tc.capabilities, got, err, tc.want)
		}
	}
	for _, capability := range contract.CapabilitiesFor(contract.ProviderLinear) {
		if _, err := RequiredOAuthScopes(contract.ProviderLinear, []contract.Capability{capability}); err != nil {
			t.Errorf("Linear capability %q has no scope mapping: %v", capability.Name, err)
		}
	}
}

// The interim preflight rule (owner decision 2026-10-10, option C): only a
// mapped provider whose connector scopes name its OAuth scopes is checked.
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

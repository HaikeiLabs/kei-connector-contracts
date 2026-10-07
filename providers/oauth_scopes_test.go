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

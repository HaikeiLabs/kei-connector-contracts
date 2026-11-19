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

package contract

import "testing"

// Postgres (Supabase included) is a read-only shared-secret
// provider. Its one capability, query.run, is database.query.run (class
// READ); the DSN for a read-only role sits behind an opaque_ref.
func TestPostgresIsReadOnlyWithOpaqueCredential(t *testing.T) {
	if ProviderPostgres != "postgres" {
		t.Fatalf("ProviderPostgres = %q, want postgres (the catalog template id)", ProviderPostgres)
	}
	if !ValidProvider(ProviderPostgres) {
		t.Fatal("postgres is not a valid provider")
	}
	caps := CapabilitiesFor(ProviderPostgres)
	if len(caps) != 1 || caps[0] != (Capability{Name: "query.run", Action: ActionRead}) {
		t.Fatalf("postgres capabilities = %+v, want [query.run read]", caps)
	}
	m := frozenMetadata(ProviderPostgres, caps)
	if err := m.Validate(); err != nil {
		t.Fatalf("postgres metadata rejected: %v", err)
	}
	if err := m.ValidateExecutable(); err != nil {
		t.Fatalf("postgres metadata not executable: %v", err)
	}
	oauth := m
	oauth.CredentialSource, oauth.Subject = CredentialSourceOAuth, "subject-1"
	if err := oauth.Validate(); err == nil {
		t.Error("postgres connector with an oauth credential source was accepted")
	}
	if len(AccountModelsFor(ProviderPostgres)) != 0 {
		t.Error("postgres has account models")
	}
	for _, write := range []Capability{
		{Name: "query.run", Action: ActionCreate},
		{Name: "query.run", Action: ActionUpdate},
		{Name: "query.run", Action: ActionDelete},
		{Name: "query.write", Action: ActionCreate},
		{Name: "database.query.run", Action: ActionRead},
	} {
		if err := frozenMetadata(ProviderPostgres, []Capability{write}).ValidateExecutable(); err == nil {
			t.Errorf("postgres capability %q/%q was accepted", write.Name, write.Action)
		}
	}
}

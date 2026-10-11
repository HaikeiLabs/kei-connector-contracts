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

// A skills repository is read with the invoking user's OAuth token or with a
// shared-secret token behind an opaque_ref. A shared OAuth account and the
// legacy single-subject oauth binding are not offered.
func TestSkillsRepoCredentialBinding(t *testing.T) {
	perUser := frozenMetadata(ProviderSkillsRepo, CapabilitiesFor(ProviderSkillsRepo))
	perUser.CredentialSource, perUser.AccountModel = CredentialSourceOAuth, AccountModelPerUser
	if err := perUser.Validate(); err != nil {
		t.Fatalf("per_user oauth rejected: %v", err)
	}
	opaque := frozenMetadata(ProviderSkillsRepo, CapabilitiesFor(ProviderSkillsRepo))
	if err := opaque.Validate(); err != nil {
		t.Fatalf("opaque_ref rejected: %v", err)
	}
	shared := frozenMetadata(ProviderSkillsRepo, CapabilitiesFor(ProviderSkillsRepo))
	shared.CredentialSource, shared.AccountModel, shared.Subject = CredentialSourceOAuth, AccountModelShared, SharedSubject(shared.ID)
	if err := shared.Validate(); err == nil {
		t.Fatal("shared account model was accepted")
	}
}

// Both capabilities are reads; a write declared on a skills_repo connector is
// not executable.
func TestSkillsRepoIsReadOnly(t *testing.T) {
	for _, c := range CapabilitiesFor(ProviderSkillsRepo) {
		if c.Action != ActionRead {
			t.Errorf("capability %q has action %q", c.Name, c.Action)
		}
	}
	for _, write := range []Capability{{Name: "skills.fetch", Action: ActionCreate}, {Name: "skills.publish", Action: ActionCreate}} {
		if err := frozenMetadata(ProviderSkillsRepo, []Capability{write}).ValidateExecutable(); err == nil {
			t.Errorf("capability %q/%q was executable", write.Name, write.Action)
		}
	}
}

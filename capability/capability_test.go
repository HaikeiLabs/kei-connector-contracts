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

package capability

import (
	"reflect"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func TestLookupCapability(t *testing.T) {
	cap, ok := LookupCapability(contract.ProviderLinear, "issue.read")
	if !ok {
		t.Fatal("expected issue.read to be defined for linear")
	}
	if cap.Name != "issue.read" || cap.Action != contract.ActionRead {
		t.Fatalf("unexpected capability: %+v", cap)
	}

	if _, ok := LookupCapability(contract.ProviderLinear, "admin.raw_sql"); ok {
		t.Fatal("undefined capability must not be found")
	}
	if _, ok := LookupCapability("unknown", "issue.read"); ok {
		t.Fatal("unknown provider must not yield capabilities")
	}
}

// TestCRMIsNotAProvider (HAI-384): there is no dedicated crm provider; a CRM
// API is an http_api connector the admin names.
func TestCRMIsNotAProvider(t *testing.T) {
	if contract.ValidProvider("crm") {
		t.Fatal(`"crm" must not be a valid provider`)
	}
	if _, err := CapabilitiesForProvider("crm"); err == nil {
		t.Fatal(`"crm" must have no capability set`)
	}
	for _, p := range Providers() {
		if p == "crm" {
			t.Fatal(`Providers() must not list "crm"`)
		}
	}
}

func TestCapabilityForErrors(t *testing.T) {
	if _, err := CapabilityFor(contract.ProviderGitHub, "repository.read"); err != nil {
		t.Fatalf("defined capability errored: %v", err)
	}
	if _, err := CapabilityFor(contract.ProviderGitHub, "repository.purge"); err == nil {
		t.Fatal("undefined capability must error")
	}
}

func TestCapabilitiesForProviderCopiesAndRejectsUnknown(t *testing.T) {
	got, err := CapabilitiesForProvider(contract.ProviderGoogle)
	if err != nil {
		t.Fatalf("CapabilitiesForProvider(google_drive): %v", err)
	}
	if len(got) == 0 {
		t.Fatal("google_drive must declare capabilities")
	}
	if _, err := CapabilitiesForProvider("atlassian"); err == nil {
		t.Fatal("unknown provider must error, not return an empty set")
	}
}

func TestProvidersSortedAndComplete(t *testing.T) {
	got := Providers()
	want := []contract.Provider{contract.ProviderDiscord, contract.ProviderFreshBooks, contract.ProviderGitHub, contract.ProviderGmail, contract.ProviderGoogle, contract.ProviderGrafana, contract.ProviderHTTPAPI, contract.ProviderLinear, contract.ProviderMercury, contract.ProviderNotion, contract.ProviderS3, contract.ProviderSkillsRepo, contract.ProviderTito}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Providers() = %v, want %v", got, want)
	}
}

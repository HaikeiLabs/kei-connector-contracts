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

package governance

import (
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// TestDecideHasNoCapabilityList (HAI-258): the control plane allows any
// well-formed capability the connector declares, and the allow carries the
// connector's declared capability, description included. Policy decides.
func TestDecideHasNoCapabilityList(t *testing.T) {
	m := decisionFixture()
	m.Capabilities = []contract.Capability{{Name: "deal.read", Action: contract.ActionRead, Description: "Read a deal"}}
	env := decisionEnvelope()
	env.Invocation.Capability = "deal.read"
	got := Decide(m, env)
	if got.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow", got.Decision, got.Reason)
	}
	if got.Capability != m.Capabilities[0] {
		t.Fatalf("capability = %+v, want the declared %+v", got.Capability, m.Capabilities[0])
	}

	m = decisionFixture()
	m.Capabilities[0].Description = "Read a lead"
	if got := Decide(m, decisionEnvelope()); got.Capability.Description != "Read a lead" {
		t.Fatalf("defined capability = %+v, want the declared description", got.Capability)
	}
}

// TestDecideAuthorizesHTTPAPIWithoutRequest (HAI-256): an http_api envelope
// carries no request details; the method comes from the capability name.
func TestDecideAuthorizesHTTPAPIWithoutRequest(t *testing.T) {
	m := fixture()
	m.Provider = contract.ProviderHTTPAPI
	m.Capabilities = []contract.Capability{{Name: "http.get", Action: contract.ActionRead}}
	m.HTTPAPI = &contract.HTTPAPI{RegisteredBaseURL: "https://api.example.test", RegisteredHostname: "api.example.test", AllowedMethods: []contract.HTTPMethod{contract.HTTPMethodGet}, AllowedPathPatterns: []string{"/contacts/*"}, RequestMaxBytes: 1, ResponseMaxBytes: 1, TimeoutMS: 1, CredentialInjection: contract.CredentialInjection{Location: contract.CredentialInjectionHeader, Name: "X-Key"}, Redaction: contract.RedactionRules{Headers: []string{"X-Key"}}}
	env := decisionEnvelope()
	env.Invocation.Capability, env.Invocation.Resource = "http.get", "/contacts/42"
	if got := Decide(m, env); got.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow", got.Decision, got.Reason)
	}
	env.Invocation.Resource = "/admin"
	if got := Decide(m, env); got.Decision != DecisionDeny || got.Reason != "HTTP path is not allowed" {
		t.Fatalf("decision = %+v, want deny on the path", got)
	}
}

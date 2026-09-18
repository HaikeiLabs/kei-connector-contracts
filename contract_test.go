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

package connectors

import "testing"

func fixture() Metadata {
	return Metadata{ID: "c-1", TenantID: "t-1", WorkspaceID: "w-1", Name: "CRM", Provider: ProviderCRM, Status: StatusActive, CredentialRef: "vault/tenant/t-1/crm", Scopes: []string{"leads:read"}, Resources: []string{"leads/*"}, Policy: PolicyAttributes{AllowedActions: []Action{ActionRead, ActionUpdate}, AllowedResources: []string{"leads"}}, Capabilities: []Capability{{Name: "lead.read", Action: ActionRead}}, CreatedBy: "u-1"}
}

func TestValidateCallEnforcesTenantWorkspaceAndCapability(t *testing.T) {
	m := fixture()
	in := Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: "c-1", Capability: "lead.read", Action: ActionRead, Resource: "leads/42", TraceID: "trace-1"}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("valid call rejected: %v", err)
	}
	in.TenantID = "other"
	if err := ValidateCall(m, in); err == nil || err.Error() != "connector not found" {
		t.Fatalf("cross-tenant call error = %v", err)
	}
}

// TestValidateCallMutationsAreStructurallyAllowed verifies the decoupling: the
// connector contract no longer requires an approval for mutating capabilities.
// Approval gating moved to the ABAC policy layer (governance pivot).
func TestValidateCallMutationsAreStructurallyAllowed(t *testing.T) {
	m := fixture()
	m.Capabilities = []Capability{{Name: "lead.update", Action: ActionUpdate}}
	in := Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: "c-1", Capability: "lead.update", Action: ActionUpdate, Resource: "leads/42", TraceID: "trace-1"}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("mutation rejected by the contract: %v", err)
	}
}

func TestCredentialRefRejectsInlineAndURLs(t *testing.T) {
	for _, ref := range []string{"", "https://example.com/token", "vault\nsecret"} {
		if err := ValidateCredentialRef(ref); err == nil {
			t.Errorf("accepted unsafe ref %q", ref)
		}
	}
}

func TestProviderDefinitionsRejectUnknownCapabilities(t *testing.T) {
	m := fixture()
	m.Capabilities = []Capability{{Name: "admin.raw_sql", Action: ActionRead}}
	if err := m.Validate(); err == nil {
		t.Fatal("unknown provider capability was accepted")
	}
}

func httpFixture() Metadata {
	h := &HTTPAPI{RegisteredBaseURL: "https://api.example.test/v1", RegisteredHostname: "api.example.test", AllowedMethods: []HTTPMethod{HTTPMethodGet, HTTPMethodHead}, AllowedPathPatterns: []string{"/contacts/*", "/health"}, QueryAllowlist: []string{"limit", "cursor"}, RequestMaxBytes: 1 << 20, ResponseMaxBytes: 2 << 20, TimeoutMS: 5000, AllowRedirects: true, MaxRedirects: 2, CredentialInjection: CredentialInjection{Location: CredentialInjectionHeader, Name: "X-API-Key"}, Redaction: RedactionRules{Headers: []string{"X-API-Key", "Authorization"}, QueryParameters: []string{"cursor"}}}
	return Metadata{ID: "c-http", TenantID: "t-1", WorkspaceID: "w-1", Name: "HTTP", Provider: ProviderHTTPAPI, Status: StatusActive, CredentialRef: "vault/tenant/t-1/http", Scopes: []string{"contacts:read"}, Resources: []string{"contacts/*"}, Policy: PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"contacts"}}, Capabilities: CapabilitiesFor(ProviderHTTPAPI), CreatedBy: "u-1", HTTPAPI: h}
}

func TestHTTPAPIContractAcceptsValidInvocation(t *testing.T) {
	m := httpFixture()
	in := Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: "c-http", Capability: "http.get", Action: ActionRead, Resource: "contacts/42", TraceID: "trace-1", HTTP: &HTTPInvocation{Method: HTTPMethodGet, Path: "/contacts/42", Query: map[string][]string{"limit": {"10"}}}}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("valid HTTP call rejected: %v", err)
	}
}

func TestHTTPAPIContractRejectsUnsafeMetadata(t *testing.T) {
	base := httpFixture()
	tests := []struct {
		name   string
		mutate func(*Metadata)
	}{
		{"localhost", func(m *Metadata) {
			m.HTTPAPI.RegisteredBaseURL = "http://localhost:8080"
			m.HTTPAPI.RegisteredHostname = "localhost"
		}},
		{"private IP", func(m *Metadata) {
			m.HTTPAPI.RegisteredBaseURL = "https://169.254.169.254"
			m.HTTPAPI.RegisteredHostname = "169.254.169.254"
		}},
		{"metadata hostname", func(m *Metadata) {
			m.HTTPAPI.RegisteredBaseURL = "https://metadata.internal"
			m.HTTPAPI.RegisteredHostname = "metadata.internal"
		}},
		{"cross host redirects", func(m *Metadata) { m.HTTPAPI.AllowCrossHostRedirects = true }},
		{"query credential injection", func(m *Metadata) { m.HTTPAPI.CredentialInjection.Location = "query" }},
		{"path traversal", func(m *Metadata) { m.HTTPAPI.AllowedPathPatterns = []string{"/contacts/../admin"} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := base
			tt.mutate(&m)
			if err := m.Validate(); err == nil {
				t.Fatal("unsafe HTTP metadata was accepted")
			}
		})
	}
}

func TestHTTPAPIContractRejectsDisallowedInvocation(t *testing.T) {
	m := httpFixture()
	for _, in := range []HTTPInvocation{{Method: HTTPMethodHead, Path: "/admin"}, {Method: "POST", Path: "/contacts/42"}, {Method: HTTPMethodGet, Path: "/contacts/../admin"}, {Method: HTTPMethodGet, Path: "/contacts/42", Query: map[string][]string{"secret": {"x"}}}} {
		if err := m.HTTPAPI.ValidateInvocation(in); err == nil {
			t.Errorf("unsafe invocation %+v was accepted", in)
		}
	}
}

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

import (
	"strings"
	"testing"
)

// The control plane authorizes on structure and on what the connector
// declares; ABAC policy decides access (HAI-258). Whether the tenant runtime
// has code for a capability is a separate, executable check (ValidateExecutable)
// that the runtime runs through setup.ValidateMetadata.

func TestValidateAcceptsAnyWellFormedDeclaredCapability(t *testing.T) {
	m := fixture()
	m.Capabilities = []Capability{{Name: "deal.read", Action: ActionRead}, {Name: "investor.stage.update", Action: ActionUpdate}}
	if err := m.Validate(); err != nil {
		t.Fatalf("connector-declared capability rejected: %v", err)
	}
}

func TestValidateStillChecksCapabilityShape(t *testing.T) {
	for name, c := range map[string]Capability{
		"empty name":      {Name: "", Action: ActionRead},
		"space in name":   {Name: "investor list", Action: ActionRead},
		"leading dot":     {Name: ".investor", Action: ActionRead},
		"control char":    {Name: "investor\n.list", Action: ActionRead},
		"unknown action":  {Name: "investor.list", Action: Action("egress")},
		"missing action":  {Name: "investor.list"},
		"secret material": {Name: "investor.list", Action: ActionRead, Description: "Bearer abc"},
	} {
		t.Run(name, func(t *testing.T) {
			m := fixture()
			m.Capabilities = []Capability{c}
			if err := m.Validate(); err == nil {
				t.Fatalf("malformed capability %+v was accepted", c)
			}
		})
	}
}

func TestValidateExecutableRequiresDefinedCapabilities(t *testing.T) {
	m := fixture()
	if err := m.ValidateExecutable(); err != nil {
		t.Fatalf("defined capability rejected: %v", err)
	}
	m.Capabilities = []Capability{{Name: "lead.read", Action: ActionRead}, {Name: "admin.raw_sql", Action: ActionRead}}
	err := m.ValidateExecutable()
	if err == nil || err.Error() != `capability "admin.raw_sql" is not defined for provider "crm"` {
		t.Fatalf("undefined capability error = %v", err)
	}
	// The definition is name and action: a defined name under another action
	// is not executable.
	m.Capabilities = []Capability{{Name: "lead.read", Action: ActionDelete}}
	if err := m.ValidateExecutable(); err == nil {
		t.Fatal("defined name with the wrong action was accepted")
	}
}

// TestValidateCallKeepsRuntimeSemantics pins that ValidateCall, which the
// tenant runtime calls before executing, still refuses a capability it has no
// code for and still needs the full HTTP request for http_api.
func TestValidateCallKeepsRuntimeSemantics(t *testing.T) {
	m := fixture()
	m.Capabilities = []Capability{{Name: "deal.read", Action: ActionRead}}
	in := validCall("c-1", "deal.read", "leads/1")
	if err := ValidateCall(m, in); err == nil || !strings.Contains(err.Error(), "is not defined for provider") {
		t.Fatalf("ValidateCall undefined capability error = %v", err)
	}
	h := httpFixture()
	if err := ValidateCall(h, validCall("c-http", "http.get", "/contacts/42")); err == nil || err.Error() != "http invocation is required" {
		t.Fatalf("ValidateCall without request error = %v", err)
	}
}

func validCall(connectorID, capability, resource string) Invocation {
	return Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: connectorID, Capability: capability, Action: ActionRead, Resource: resource, TraceID: "trace-1"}
}

func TestAuthorizeCallAllowsDeclaredCapability(t *testing.T) {
	m := fixture()
	m.Capabilities = []Capability{{Name: "deal.read", Action: ActionRead}}
	if err := AuthorizeCall(m, validCall("c-1", "deal.read", "deals/1")); err != nil {
		t.Fatalf("declared capability rejected: %v", err)
	}
	if err := AuthorizeCall(m, validCall("c-1", "deal.write", "deals/1")); err == nil || err.Error() != "capability is not allowed" {
		t.Fatalf("undeclared capability error = %v", err)
	}
	in := validCall("c-1", "deal.read", "deals/1")
	in.Action = ActionUpdate
	if err := AuthorizeCall(m, in); err == nil || err.Error() != "capability action mismatch" {
		t.Fatalf("action mismatch error = %v", err)
	}
	s3 := fixture()
	s3.Provider = ProviderS3
	s3.Capabilities = []Capability{{Name: "object.read", Action: ActionRead}}
	if err := AuthorizeCall(s3, validCall("c-1", "object.read", "s3://billing/1.json")); err != nil {
		t.Fatalf("s3 connector rejected: %v", err)
	}
	s3.Provider = Provider("mainframe")
	if err := AuthorizeCall(s3, validCall("c-1", "object.read", "x")); err == nil || !strings.Contains(err.Error(), "unsupported provider") {
		t.Fatalf("unknown provider error = %v", err)
	}
}

// TestAuthorizeCallDerivesHTTPRequestFromCapability (HAI-256): the tenant
// runtime keeps the request local, so the control plane authorizes the method
// the capability names (http.get -> GET) and the resource, which is the
// request path, against the connector's configured methods and patterns.
func TestAuthorizeCallDerivesHTTPRequestFromCapability(t *testing.T) {
	m := httpFixture() // GET, HEAD; /contacts/*, /health
	m.Capabilities = append(m.Capabilities, Capability{Name: "http.delete", Action: ActionDelete}, Capability{Name: "contacts.list", Action: ActionRead})
	for _, tc := range []struct{ capability, resource string }{
		{"http.get", "/contacts/42"},
		{"http.head", "/contacts/42"},
		{"http.get", "/health"},
	} {
		if err := AuthorizeCall(m, validCall("c-http", tc.capability, tc.resource)); err != nil {
			t.Errorf("%s %s rejected: %v", tc.capability, tc.resource, err)
		}
	}

	del := validCall("c-http", "http.delete", "/contacts/42")
	del.Action = ActionDelete
	for name, tc := range map[string]struct {
		in   Invocation
		want string
	}{
		"method not configured":      {del, "HTTP method is not allowed"},
		"capability with no method":  {validCall("c-http", "contacts.list", "/contacts/42"), "capability has no HTTP method"},
		"path outside patterns":      {validCall("c-http", "http.get", "/admin/users"), "HTTP path is not allowed"},
		"prefix without wildcard":    {validCall("c-http", "http.get", "/healthz"), "HTTP path is not allowed"},
		"traversal":                  {validCall("c-http", "http.get", "/contacts/../admin"), ""},
		"encoded traversal":          {validCall("c-http", "http.get", "/contacts/%2e%2e/admin"), ""},
		"query in resource":          {validCall("c-http", "http.get", "/contacts/42?secret=1"), ""},
		"absolute URL resource":      {validCall("c-http", "http.get", "https://evil.example/contacts/1"), ""},
		"resource without a slash":   {validCall("c-http", "http.get", "contacts/42"), ""},
		"undeclared http capability": {validCall("c-http", "http.options", "/contacts/42"), ""},
	} {
		t.Run(name, func(t *testing.T) {
			err := AuthorizeCall(m, tc.in)
			if err == nil {
				t.Fatal("allowed")
			}
			if tc.want != "" && err.Error() != tc.want {
				t.Fatalf("error = %q, want %q", err, tc.want)
			}
		})
	}

	// A full request, when one is present, is checked as sent.
	full := validCall("c-http", "http.get", "/contacts/42")
	full.HTTP = &HTTPInvocation{Method: HTTPMethodGet, Path: "/contacts/42", Query: map[string][]string{"api_key": {"x"}}}
	if err := AuthorizeCall(m, full); err == nil || err.Error() != `query parameter "api_key" is not allowed` {
		t.Fatalf("full request error = %v", err)
	}
}

// TestValidateCredentialRefNamesTheLegacyNamespace: the colon form of the
// legacy kei-oauth namespace is reported as the namespace, not as a scheme.
func TestValidateCredentialRefNamesTheLegacyNamespace(t *testing.T) {
	for _, ref := range []string{"kei-oauth:github", "kei-oauth-abc123", "KEI-OAUTH:x"} {
		if err := ValidateCredentialRef(ref); err == nil || err.Error() != "credential_ref must not use the legacy kei-oauth namespace" {
			t.Errorf("%q error = %v", ref, err)
		}
	}
	if err := ValidateCredentialRef("https://vault.example/x"); err == nil || err.Error() != "credential_ref must be an opaque reference" {
		t.Errorf("url error = %v", err)
	}
}

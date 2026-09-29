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

func oauthConnector(provider Provider, model AccountModel, subject string) Metadata {
	m := fixture()
	m.Provider = provider
	m.Capabilities = CapabilitiesFor(provider)[:1]
	m.CredentialSource = CredentialSourceOAuth
	m.CredentialRef = "cref-0001"
	m.AccountModel = model
	m.Subject = subject
	return m
}

func TestAccountModelsPerProvider(t *testing.T) {
	want := map[Provider][]AccountModel{
		ProviderGmail:  {AccountModelPerUser, AccountModelShared, AccountModelDomainDelegation},
		ProviderGoogle: {AccountModelPerUser, AccountModelShared, AccountModelDomainDelegation},
		ProviderLinear: {AccountModelPerUser, AccountModelShared},
		ProviderGitHub: {AccountModelPerUser, AccountModelShared},
		ProviderTito:   nil,
		ProviderCRM:    nil,
	}
	for provider, models := range want {
		got := AccountModelsFor(provider)
		if len(got) != len(models) {
			t.Fatalf("%s account models = %v, want %v", provider, got, models)
		}
		for i := range models {
			if got[i] != models[i] {
				t.Fatalf("%s account models = %v, want %v", provider, got, models)
			}
		}
	}
}

// A per-user connector acts for whoever invokes it, so it declares no subject,
// and any invoking subject passes the structural check (policy gates access).
func TestPerUserConnectorHasNoSubjectAndAcceptsAnyInvoker(t *testing.T) {
	m := oauthConnector(ProviderGmail, AccountModelPerUser, "")
	if err := m.Validate(); err != nil {
		t.Fatalf("per_user connector rejected: %v", err)
	}
	withSubject := oauthConnector(ProviderGmail, AccountModelPerUser, "user-1")
	if err := withSubject.Validate(); err == nil {
		t.Fatal("per_user connector with a declared subject was accepted")
	}
	in := Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "kei-user-7", AgentID: "a-1", ConnectorID: "c-1", Capability: "message.search", Action: ActionRead, Resource: "gmail/messages", TraceID: "trace-1"}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("per_user invocation rejected: %v", err)
	}
	if got := CredentialSubject(m, in); got != "kei-user-7" {
		t.Fatalf("CredentialSubject = %q, want the invoking user", got)
	}
}

// A shared connector's token row belongs to the connector itself, so its
// subject must be connector:<id>, and invokers are gated by policy only.
func TestSharedConnectorIsBoundToItsServiceSubject(t *testing.T) {
	m := oauthConnector(ProviderLinear, AccountModelShared, SharedSubject("c-1"))
	if err := m.Validate(); err != nil {
		t.Fatalf("shared connector rejected: %v", err)
	}
	for _, subject := range []string{"", "admin-1", SharedSubject("c-2")} {
		bad := oauthConnector(ProviderLinear, AccountModelShared, subject)
		if err := bad.Validate(); err == nil {
			t.Errorf("shared connector with subject %q was accepted", subject)
		}
	}
	in := Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "kei-user-7", AgentID: "a-1", ConnectorID: "c-1", Capability: "team.read", Action: ActionRead, Resource: "linear/team/ENG", TraceID: "trace-1"}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("shared invocation by another user rejected: %v", err)
	}
	if got := CredentialSubject(m, in); got != "connector:c-1" {
		t.Fatalf("CredentialSubject = %q, want connector:c-1", got)
	}
}

// Domain-wide delegation uses a service-account key through the shared-secret
// path and names the impersonated mailbox in config.
func TestDomainDelegationConnectorUsesOpaqueRefAndImpersonatedEmail(t *testing.T) {
	m := fixture()
	m.Provider = ProviderGmail
	m.Capabilities = CapabilitiesFor(ProviderGmail)
	m.CredentialSource = CredentialSourceOpaqueRef
	m.AccountModel = AccountModelDomainDelegation
	m.Config = map[string]any{"impersonate_email": "events@example.com"}
	if err := m.Validate(); err != nil {
		t.Fatalf("delegation connector rejected: %v", err)
	}
	if got := CredentialSubject(m, Invocation{Subject: "kei-user-7"}); got != "connector:c-1" {
		t.Fatalf("CredentialSubject = %q, want connector:c-1", got)
	}

	oauth := m
	oauth.CredentialSource = CredentialSourceOAuth
	if err := oauth.Validate(); err == nil {
		t.Error("delegation connector with an oauth credential source was accepted")
	}
	missing := m
	missing.Config = map[string]any{}
	if err := missing.Validate(); err == nil {
		t.Error("delegation connector without impersonate_email was accepted")
	}
	github := m
	github.Provider = ProviderGitHub
	github.Capabilities = CapabilitiesFor(ProviderGitHub)[:1]
	if err := github.Validate(); err == nil {
		t.Error("delegation connector for github was accepted")
	}
}

func TestAccountModelRules(t *testing.T) {
	cases := map[string]Metadata{
		"unknown model": oauthConnector(ProviderGmail, AccountModel("team"), ""),
		"model on a non-oauth tito": func() Metadata {
			m := fixture()
			m.Provider = ProviderTito
			m.Capabilities = CapabilitiesFor(ProviderTito)
			m.AccountModel = AccountModelShared
			return m
		}(),
		"per_user over opaque_ref": func() Metadata {
			m := oauthConnector(ProviderGitHub, AccountModelPerUser, "")
			m.CredentialSource = CredentialSourceOpaqueRef
			return m
		}(),
		"shared over opaque_ref": func() Metadata {
			m := oauthConnector(ProviderGitHub, AccountModelShared, SharedSubject("c-1"))
			m.CredentialSource = CredentialSourceOpaqueRef
			return m
		}(),
	}
	for name, m := range cases {
		if err := m.Validate(); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

// Legacy oauth connectors (no account model) keep the v0.1.0 catalog rule: the
// credential subject is the declared subject.
func TestCredentialSubjectForLegacyAndOpaqueConnectors(t *testing.T) {
	legacy := oauthConnector(ProviderGitHub, "", "gh:12345")
	if got := CredentialSubject(legacy, Invocation{Subject: "gh:12345"}); got != "gh:12345" {
		t.Fatalf("legacy CredentialSubject = %q", got)
	}
	opaque := fixture()
	if got := CredentialSubject(opaque, Invocation{Subject: "kei-user-7"}); got != "connector:c-1" {
		t.Fatalf("opaque_ref CredentialSubject = %q, want connector:c-1", got)
	}
}

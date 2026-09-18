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

import (
	"strings"
	"testing"
	"time"
)

func frozenMetadata(provider Provider, capabilities []Capability) Metadata {
	return Metadata{
		ID: "connector-1", TenantID: "tenant-1", WorkspaceID: "workspace-1",
		Name: "frozen", Provider: provider, Status: StatusActive,
		CredentialRef: "vault/tenant-1/frozen", Scopes: []string{"read"},
		Resources: []string{"resource-1"}, Capabilities: capabilities,
		CreatedBy: "subject-1",
	}
}

func frozenInvocation(capability string, action Action) Invocation {
	return Invocation{
		TenantID: "tenant-1", WorkspaceID: "workspace-1", Subject: "subject-1",
		AgentID: "agent-1", ConnectorID: "connector-1", Capability: capability,
		Action: action, Resource: "resource-1", TraceID: "trace-1", IdempotencyKey: "idem-1",
	}
}

func TestConnectorContractFreezeProviderCatalog(t *testing.T) {
	want := map[Provider]map[string]Action{
		ProviderGitHub: {
			"repository.read": ActionRead, "issue.read": ActionRead,
			"pull_request.read": ActionRead, "check.read": ActionRead,
			"workflow.read": ActionRead, "issue.create": ActionCreate,
			"issue.update": ActionUpdate, "pull_request.create": ActionCreate,
			"pull_request.update": ActionUpdate, "issue.comment": ActionComment,
		},
		ProviderGoogle: {"drive.search": ActionRead, "drive.metadata.read": ActionRead, "docs.read": ActionRead},
		ProviderNotion: {"search": ActionRead, "page.read": ActionRead, "database.query": ActionRead},
		ProviderLinear: {"team.read": ActionRead, "project.read": ActionRead, "cycle.read": ActionRead, "issue.read": ActionRead, "issue.create": ActionCreate, "issue.update": ActionUpdate},
		// The finance providers are frozen as read-only. If a future change
		// adds a mutation capability here, this test fails and the reviewer
		// has to justify letting an agent write to a ledger or a bank.
		ProviderFreshBooks: {"invoice.read": ActionRead, "expense.read": ActionRead, "payment.read": ActionRead, "client.read": ActionRead},
		ProviderMercury:    {"account.read": ActionRead, "transaction.read": ActionRead, "balance.read": ActionRead},
	}
	for provider, expected := range want {
		got := CapabilitiesFor(provider)
		if len(got) != len(expected) {
			t.Fatalf("%s capability count = %d, want %d", provider, len(got), len(expected))
		}
		for _, capability := range got {
			if action, ok := expected[capability.Name]; !ok || action != capability.Action {
				t.Errorf("%s capability %q has action %q; want %q", provider, capability.Name, capability.Action, action)
			}
		}
	}
}

func TestConnectorContractFreezeBindingAndCredentialRules(t *testing.T) {
	m := frozenMetadata(ProviderGoogle, CapabilitiesFor(ProviderGoogle))
	if err := m.Validate(); err != nil {
		t.Fatalf("valid metadata rejected: %v", err)
	}
	for _, ref := range []string{"", "https://secrets.example/ref", "Bearer token", "ghp_inline", "secret=value"} {
		if err := ValidateCredentialRef(ref); err == nil {
			t.Errorf("credential ref %q was accepted", ref)
		}
	}
	for _, ref := range []string{"vault/tenant-1/google-drive", "kv-v2/data/connectors/drive"} {
		if err := ValidateCredentialRef(ref); err != nil {
			t.Errorf("opaque credential ref %q rejected: %v", ref, err)
		}
	}

	valid := frozenInvocation("drive.search", ActionRead)
	if err := ValidateCall(m, valid); err != nil {
		t.Fatalf("valid invocation rejected: %v", err)
	}
	for name, mutate := range map[string]func(*Invocation){
		"tenant":    func(in *Invocation) { in.TenantID = "other-tenant" },
		"workspace": func(in *Invocation) { in.WorkspaceID = "other-workspace" },
		"connector": func(in *Invocation) { in.ConnectorID = "other-connector" },
	} {
		in := valid
		mutate(&in)
		if err := ValidateCall(m, in); err == nil || err.Error() != "connector not found" {
			t.Errorf("%s mismatch error = %v, want connector not found", name, err)
		}
	}
}

func TestConnectorContractFreezeEnvelopeAuditIdempotencyAndLegacyPolicy(t *testing.T) {
	m := frozenMetadata(ProviderGitHub, CapabilitiesFor(ProviderGitHub))
	m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"legacy"}}
	in := frozenInvocation("repository.read", ActionRead)
	first := NewClient().Invoke(m, in)
	second := NewClient().Invoke(m, in)
	if first.Decision.Decision != DecisionAllow || second.Decision.Decision != DecisionAllow {
		t.Fatalf("valid invocation decisions = %q, %q", first.Decision.Decision, second.Decision.Decision)
	}
	if first.Envelope.Version != EnvelopeVersion1 || first.Envelope.MintedBy != MintedByControlPlane {
		t.Fatalf("envelope = %+v", first.Envelope)
	}
	if first.AuditRecord.SpanID != "idem-1" || second.AuditRecord.SpanID != "idem-1" {
		t.Fatalf("audit span IDs = %q, %q", first.AuditRecord.SpanID, second.AuditRecord.SpanID)
	}
	if first.AuditRecord.ToolName != "connector/connector-1/repository.read" {
		t.Fatalf("audit tool name = %q", first.AuditRecord.ToolName)
	}
	if first.Envelope.IssuedAt.IsZero() || first.Envelope.IssuedAt.Equal(time.Time{}) {
		t.Fatal("envelope timestamp was not issued")
	}

	denied := in
	denied.Subject = ""
	if got := Decide(m, Envelope{Version: EnvelopeVersion1, Invocation: denied, MintedBy: MintedByControlPlane, IssuedAt: time.Now()}); got.Decision != DecisionDeny || !strings.Contains(got.Reason, "subject") {
		t.Fatalf("missing subject decision = %+v", got)
	}
}

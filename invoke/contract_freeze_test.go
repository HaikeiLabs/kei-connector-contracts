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

package invoke

import (
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/envelope"
	"github.com/HaikeiLabs/kei-connector-contracts/governance"
)

// These helpers mirror contract/contract_freeze_test.go; this freeze test
// spans envelope, governance, and audit, so it lives in the top package.
func frozenMetadata(provider contract.Provider, capabilities []contract.Capability) contract.Metadata {
	return contract.Metadata{
		ID: "connector-1", TenantID: "tenant-1", WorkspaceID: "workspace-1",
		Name: "frozen", Provider: provider, Status: contract.StatusActive,
		CredentialRef: "vault/tenant-1/frozen", Scopes: []string{"read"},
		Resources: []string{"resource-1"}, Capabilities: capabilities,
		CreatedBy: "subject-1",
	}
}

func frozenInvocation(capability string, action contract.Action) contract.Invocation {
	return contract.Invocation{
		TenantID: "tenant-1", WorkspaceID: "workspace-1", Subject: "subject-1",
		AgentID: "agent-1", ConnectorID: "connector-1", Capability: capability,
		Action: action, Resource: "resource-1", TraceID: "trace-1", IdempotencyKey: "idem-1",
	}
}

func TestConnectorContractFreezeEnvelopeAuditIdempotencyAndLegacyPolicy(t *testing.T) {
	m := frozenMetadata(contract.ProviderGitHub, contract.CapabilitiesFor(contract.ProviderGitHub))
	m.Policy = contract.PolicyAttributes{AllowedActions: []contract.Action{contract.ActionRead}, AllowedResources: []string{"legacy"}}
	in := frozenInvocation("repository.read", contract.ActionRead)
	first := NewClient().Invoke(m, in)
	second := NewClient().Invoke(m, in)
	if first.Decision.Decision != governance.DecisionAllow || second.Decision.Decision != governance.DecisionAllow {
		t.Fatalf("valid invocation decisions = %q, %q", first.Decision.Decision, second.Decision.Decision)
	}
	if first.Envelope.Version != envelope.EnvelopeVersion1 || first.Envelope.MintedBy != envelope.MintedByControlPlane {
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
	if got := governance.Decide(m, envelope.Envelope{Version: envelope.EnvelopeVersion1, Invocation: denied, MintedBy: envelope.MintedByControlPlane, IssuedAt: time.Now()}); got.Decision != governance.DecisionDeny || !strings.Contains(got.Reason, "subject") {
		t.Fatalf("missing subject decision = %+v", got)
	}
}

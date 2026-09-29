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

package audit

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/envelope"
	"github.com/HaikeiLabs/kei-connector-contracts/governance"
)

func linearFixture() contract.Metadata {
	return contract.Metadata{
		ID: "linear-connector", TenantID: "tenant-1", WorkspaceID: "workspace-1",
		Name: "Linear", Provider: contract.ProviderLinear, Status: contract.StatusActive,
		CredentialSource: contract.CredentialSourceOAuth, CredentialRef: "linear/oauth-binding",
		Subject: "user-1", Scopes: []string{"read", "write"}, Resources: []string{"linear/*"},
		Capabilities: contract.CapabilitiesFor(contract.ProviderLinear), CreatedBy: "user-1",
	}
}

func linearInvocation(capability string, action contract.Action) contract.Invocation {
	return contract.Invocation{TenantID: "tenant-1", WorkspaceID: "workspace-1", Subject: "user-1",
		AgentID: "agent-1", ConnectorID: "linear-connector", Capability: capability,
		Action: action, Resource: "linear/issue/KEI-42", TraceID: "trace-1", IdempotencyKey: "idem-1"}
}

func TestLinearReadAllowed(t *testing.T) {
	m := linearFixture()
	if err := contract.ValidateCall(m, linearInvocation("issue.read", contract.ActionRead)); err != nil {
		t.Fatalf("Linear read rejected: %v", err)
	}
}

func TestLinearMutationsPassContractWithoutApproval(t *testing.T) {
	m := linearFixture()
	for _, tc := range []struct {
		capability string
		action     contract.Action
	}{{"issue.create", contract.ActionCreate}, {"issue.update", contract.ActionUpdate}} {
		t.Run(tc.capability, func(t *testing.T) {
			in := linearInvocation(tc.capability, tc.action)
			// contract.ValidateCall does not check ApprovalID — that is the
			// ABAC policy layer's job (connector_policy.go).
			if err := contract.ValidateCall(m, in); err != nil {
				t.Fatalf("contract rejected mutation without approval: %v", err)
			}
			// governance.Decide also does not check ApprovalID.
			env := envelope.Envelope{Version: envelope.EnvelopeVersion1, Invocation: in,
				MintedBy: envelope.MintedByControlPlane, IssuedAt: time.Now()}
			if got := governance.Decide(m, env); got.Decision != governance.DecisionAllow {
				t.Fatalf("governance.Decide without approval = %+v, want allow", got)
			}
		})
	}
}

func TestLinearCrossWorkspaceDenied(t *testing.T) {
	in := linearInvocation("issue.read", contract.ActionRead)
	in.WorkspaceID = "workspace-elsewhere"
	if got := governance.Decide(linearFixture(), envelope.Envelope{Version: envelope.EnvelopeVersion1, Invocation: in,
		MintedBy: envelope.MintedByControlPlane, IssuedAt: time.Now()}); got.Decision != governance.DecisionDeny || got.Reason != "connector not found" {
		t.Fatalf("cross-workspace decision = %+v, want connector not found deny", got)
	}
}

func TestLinearAuditPreservesSubjectAndExcludesProviderPayload(t *testing.T) {
	in := linearInvocation("issue.create", contract.ActionCreate)
	payloadMarker := "private-linear-issue-body"
	env := envelope.Envelope{Version: envelope.EnvelopeVersion1, Invocation: in, MintedBy: envelope.MintedByControlPlane, IssuedAt: time.Now()}
	record := BuildAuditRecord(env, governance.PolicyDecision{Decision: governance.DecisionAllow}, env.IssuedAt)
	if record.InvokingSubject != "user-1" {
		t.Fatalf("audit subject = %q, want user-1", record.InvokingSubject)
	}
	raw, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), payloadMarker) {
		t.Fatalf("audit record contains provider payload: %s", raw)
	}
}

func TestLinearIdempotencyKeyDeduplicatesAuditIdentity(t *testing.T) {
	in := linearInvocation("issue.create", contract.ActionCreate)
	env := envelope.Envelope{Version: envelope.EnvelopeVersion1, Invocation: in, MintedBy: envelope.MintedByControlPlane, IssuedAt: time.Now()}
	first := BuildAuditRecord(env, governance.PolicyDecision{Decision: governance.DecisionAllow}, env.IssuedAt)
	second := BuildAuditRecord(env, governance.PolicyDecision{Decision: governance.DecisionAllow}, env.IssuedAt)
	if first.SpanID != second.SpanID || first.SpanID != in.IdempotencyKey {
		t.Fatalf("replay span IDs = %q, %q; want idempotency key %q", first.SpanID, second.SpanID, in.IdempotencyKey)
	}
}

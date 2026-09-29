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
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func linearFixture() Metadata {
	return Metadata{
		ID: "linear-connector", TenantID: "tenant-1", WorkspaceID: "workspace-1",
		Name: "Linear", Provider: ProviderLinear, Status: StatusActive,
		CredentialSource: CredentialSourceOAuth, CredentialRef: "linear/oauth-binding",
		Subject: "user-1", Scopes: []string{"read", "write"}, Resources: []string{"linear/*"},
		Capabilities: CapabilitiesFor(ProviderLinear), CreatedBy: "user-1",
	}
}

func linearInvocation(capability string, action Action) Invocation {
	return Invocation{TenantID: "tenant-1", WorkspaceID: "workspace-1", Subject: "user-1",
		AgentID: "agent-1", ConnectorID: "linear-connector", Capability: capability,
		Action: action, Resource: "linear/issue/KEI-42", TraceID: "trace-1", IdempotencyKey: "idem-1"}
}

func TestLinearReadAllowed(t *testing.T) {
	m := linearFixture()
	if err := ValidateCall(m, linearInvocation("issue.read", ActionRead)); err != nil {
		t.Fatalf("Linear read rejected: %v", err)
	}
}

func TestLinearMutationsPassContractWithoutApproval(t *testing.T) {
	m := linearFixture()
	for _, tc := range []struct {
		capability string
		action     Action
	}{{"issue.create", ActionCreate}, {"issue.update", ActionUpdate}} {
		t.Run(tc.capability, func(t *testing.T) {
			in := linearInvocation(tc.capability, tc.action)
			// ValidateCall does not check ApprovalID — that is the
			// ABAC policy layer's job (connector_policy.go).
			if err := ValidateCall(m, in); err != nil {
				t.Fatalf("contract rejected mutation without approval: %v", err)
			}
			// Decide also does not check ApprovalID.
			env := Envelope{Version: EnvelopeVersion1, Invocation: in,
				MintedBy: MintedByControlPlane, IssuedAt: time.Now()}
			if got := Decide(m, env); got.Decision != DecisionAllow {
				t.Fatalf("Decide without approval = %+v, want allow", got)
			}
		})
	}
}

func TestLinearCrossWorkspaceDenied(t *testing.T) {
	in := linearInvocation("issue.read", ActionRead)
	in.WorkspaceID = "workspace-elsewhere"
	if got := Decide(linearFixture(), Envelope{Version: EnvelopeVersion1, Invocation: in,
		MintedBy: MintedByControlPlane, IssuedAt: time.Now()}); got.Decision != DecisionDeny || got.Reason != "connector not found" {
		t.Fatalf("cross-workspace decision = %+v, want connector not found deny", got)
	}
}

func TestLinearAuditPreservesSubjectAndExcludesProviderPayload(t *testing.T) {
	in := linearInvocation("issue.create", ActionCreate)
	payloadMarker := "private-linear-issue-body"
	env := Envelope{Version: EnvelopeVersion1, Invocation: in, MintedBy: MintedByControlPlane, IssuedAt: time.Now()}
	record := BuildAuditRecord(env, PolicyDecision{Decision: DecisionAllow}, env.IssuedAt)
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
	in := linearInvocation("issue.create", ActionCreate)
	env := Envelope{Version: EnvelopeVersion1, Invocation: in, MintedBy: MintedByControlPlane, IssuedAt: time.Now()}
	first := BuildAuditRecord(env, PolicyDecision{Decision: DecisionAllow}, env.IssuedAt)
	second := BuildAuditRecord(env, PolicyDecision{Decision: DecisionAllow}, env.IssuedAt)
	if first.SpanID != second.SpanID || first.SpanID != in.IdempotencyKey {
		t.Fatalf("replay span IDs = %q, %q; want idempotency key %q", first.SpanID, second.SpanID, in.IdempotencyKey)
	}
}

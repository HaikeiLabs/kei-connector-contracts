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
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/envelope"
)

func decisionFixture() contract.Metadata {
	m := fixture()
	m.Capabilities = []contract.Capability{
		{Name: "lead.read", Action: contract.ActionRead},
		{Name: "lead.update", Action: contract.ActionUpdate},
	}
	return m
}

func decisionEnvelope() envelope.Envelope {
	return envelope.Envelope{
		Version: envelope.EnvelopeVersion1,
		Invocation: contract.Invocation{
			TenantID:       "t-1",
			WorkspaceID:    "w-1",
			Subject:        "u-1",
			AgentID:        "a-1",
			ConnectorID:    "c-1",
			Capability:     "lead.read",
			Action:         contract.ActionRead,
			Resource:       "leads/42",
			TraceID:        "trace-1",
			IdempotencyKey: "idem-1",
		},
		MintedBy: envelope.MintedByControlPlane,
		IssuedAt: time.Now().UTC(),
	}
}

func TestDecideAllowsGovernedRead(t *testing.T) {
	decision := Decide(decisionFixture(), decisionEnvelope())
	if decision.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow", decision.Decision, decision.Reason)
	}
	if decision.Capability.Name != "lead.read" {
		t.Fatalf("capability = %q, want lead.read", decision.Capability.Name)
	}
}

func TestDecideDeniesInsteadOfErroring(t *testing.T) {
	m := decisionFixture()
	env := decisionEnvelope()

	cases := []struct {
		name   string
		mutate func(*contract.Metadata, *envelope.Envelope)
		want   string
	}{
		{"inactive connector", func(m *contract.Metadata, _ *envelope.Envelope) { m.Status = contract.StatusSuspended }, "connector is not active"},
		{"cross-tenant", func(_ *contract.Metadata, e *envelope.Envelope) { e.Invocation.TenantID = "other" }, "connector not found"},
		{"cross-workspace", func(_ *contract.Metadata, e *envelope.Envelope) { e.Invocation.WorkspaceID = "other" }, "connector not found"},
		{"wrong connector", func(_ *contract.Metadata, e *envelope.Envelope) { e.Invocation.ConnectorID = "other" }, "connector not found"},
		{"undeclared undefined capability", func(_ *contract.Metadata, e *envelope.Envelope) { e.Invocation.Capability = "admin.raw_sql" }, "capability is not allowed"},
		{"undeclared capability", func(m *contract.Metadata, e *envelope.Envelope) {
			m.Capabilities = []contract.Capability{{Name: "lead.read", Action: contract.ActionRead}}
			e.Invocation.Capability = "lead.update"
		}, "capability is not allowed"},
		{"capability action mismatch", func(_ *contract.Metadata, e *envelope.Envelope) {
			e.Invocation.Capability = "lead.update"
		}, "capability action mismatch"},
		{"malformed envelope", func(_ *contract.Metadata, e *envelope.Envelope) { e.Invocation.TraceID = "" }, "invalid invocation envelope"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mm, ee := m, env
			tc.mutate(&mm, &ee)
			decision := Decide(mm, ee)
			if decision.Decision != DecisionDeny {
				t.Fatalf("decision = %s, want deny", decision.Decision)
			}
			if !strings.Contains(decision.Reason, tc.want) {
				t.Fatalf("reason = %q, want it to contain %q", decision.Reason, tc.want)
			}
			if decision.Capability.Name != "" {
				t.Fatalf("deny must not carry a capability, got %+v", decision.Capability)
			}
		})
	}
}

// TestDecideAllowsDeclaredMutation documents that a mutation is decided by
// structure and capability declaration alone; access is ABAC policy.
func TestDecideAllowsDeclaredMutation(t *testing.T) {
	env := decisionEnvelope()
	env.Invocation.Capability = "lead.update"
	env.Invocation.Action = contract.ActionUpdate

	decision := Decide(decisionFixture(), env)
	if decision.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow for declared mutation", decision.Decision, decision.Reason)
	}
	if decision.Capability.Name != "lead.update" {
		t.Fatalf("capability = %+v, want lead.update", decision.Capability)
	}
}

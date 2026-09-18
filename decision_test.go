package connectors

import (
	"strings"
	"testing"
	"time"
)

func decisionFixture() Metadata {
	m := fixture()
	m.Capabilities = []Capability{
		{Name: "lead.read", Action: ActionRead},
		{Name: "lead.update", Action: ActionUpdate},
	}
	return m
}

func decisionEnvelope() Envelope {
	return Envelope{
		Version: EnvelopeVersion1,
		Invocation: Invocation{
			TenantID:       "t-1",
			WorkspaceID:    "w-1",
			Subject:        "u-1",
			AgentID:        "a-1",
			ConnectorID:    "c-1",
			Capability:     "lead.read",
			Action:         ActionRead,
			Resource:       "leads/42",
			TraceID:        "trace-1",
			IdempotencyKey: "idem-1",
		},
		MintedBy: MintedByControlPlane,
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
		mutate func(*Metadata, *Envelope)
		want   string
	}{
		{"inactive connector", func(m *Metadata, _ *Envelope) { m.Status = StatusSuspended }, "connector is not active"},
		{"cross-tenant", func(_ *Metadata, e *Envelope) { e.Invocation.TenantID = "other" }, "connector not found"},
		{"cross-workspace", func(_ *Metadata, e *Envelope) { e.Invocation.WorkspaceID = "other" }, "connector not found"},
		{"wrong connector", func(_ *Metadata, e *Envelope) { e.Invocation.ConnectorID = "other" }, "connector not found"},
		{"undefined capability", func(_ *Metadata, e *Envelope) { e.Invocation.Capability = "admin.raw_sql" }, "not defined for provider"},
		{"undeclared capability", func(m *Metadata, e *Envelope) {
			m.Capabilities = []Capability{{Name: "lead.read", Action: ActionRead}}
			e.Invocation.Capability = "lead.update"
		}, "capability is not allowed"},
		{"capability action mismatch", func(_ *Metadata, e *Envelope) {
			e.Invocation.Capability = "lead.update"
		}, "capability action mismatch"},
		{"malformed envelope", func(_ *Metadata, e *Envelope) { e.Invocation.TraceID = "" }, "invalid invocation envelope"},
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

// TestDecideAllowsMutationWithoutApproval documents the decoupling: a mutation
// is decided by structure and capability declaration alone; approval gating
// lives in the ABAC policy layer.
func TestDecideAllowsMutationWithoutApproval(t *testing.T) {
	env := decisionEnvelope()
	env.Invocation.Capability = "lead.update"
	env.Invocation.Action = ActionUpdate

	decision := Decide(decisionFixture(), env)
	if decision.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow for unapproved mutation", decision.Decision, decision.Reason)
	}
	if decision.Capability.Name != "lead.update" {
		t.Fatalf("capability = %+v, want lead.update", decision.Capability)
	}
}

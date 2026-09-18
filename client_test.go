package connectors

import (
	"testing"
	"time"
)

func TestClientInvokeMintsGovernedEnvelope(t *testing.T) {
	m := decisionFixture()
	in := Invocation{
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
	}
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	c := &Client{now: func() time.Time { return now }}
	outcome := c.Invoke(m, in)

	if outcome.Envelope.Version != EnvelopeVersion1 {
		t.Fatalf("envelope version = %q", outcome.Envelope.Version)
	}
	if outcome.Envelope.MintedBy != MintedByControlPlane {
		t.Fatalf("minted_by = %q, want the control plane", outcome.Envelope.MintedBy)
	}
	if !outcome.Envelope.IssuedAt.Equal(now) {
		t.Fatalf("issued_at = %v, want %v", outcome.Envelope.IssuedAt, now)
	}
	if outcome.Decision.Decision != DecisionAllow {
		t.Fatalf("decision = %s (%s), want allow", outcome.Decision.Decision, outcome.Decision.Reason)
	}
	if outcome.AuditRecord.SpanID != "idem-1" {
		t.Fatalf("audit span_id = %q, want the idempotency key", outcome.AuditRecord.SpanID)
	}
}

func TestClientInvokeDeniesWithoutErroring(t *testing.T) {
	c := &Client{now: func() time.Time { return time.Now() }}

	// A structural violation (cross-workspace) must be a structured deny, never
	// an error.
	outcome := c.Invoke(decisionFixture(), Invocation{
		TenantID:    "t-1",
		WorkspaceID: "w-other",
		Subject:     "u-1",
		AgentID:     "a-1",
		ConnectorID: "c-1",
		Capability:  "lead.read",
		Action:      ActionRead,
		Resource:    "leads/42",
		TraceID:     "trace-2",
	})
	if outcome.Decision.Decision != DecisionDeny {
		t.Fatalf("decision = %s, want deny", outcome.Decision.Decision)
	}
	if outcome.Decision.Reason != "connector not found" {
		t.Fatalf("reason = %q, want connector not found", outcome.Decision.Reason)
	}
	if outcome.AuditRecord.Decision != string(DecisionDeny) {
		t.Fatalf("audit decision = %q, want deny", outcome.AuditRecord.Decision)
	}
}

package connectors

import (
	"strings"
	"testing"
	"time"
)

func validEnvelope() Envelope {
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

func TestEnvelopeValidateAcceptsValidEnvelope(t *testing.T) {
	if err := validEnvelope().Validate(); err != nil {
		t.Fatalf("valid envelope rejected: %v", err)
	}
}

func TestEnvelopeValidateFailsClosed(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Envelope)
		want   string
	}{
		{"unknown version", func(e *Envelope) { e.Version = "2.0" }, "unsupported envelope version"},
		{"missing minted_by", func(e *Envelope) { e.MintedBy = "" }, "minted by the control plane"},
		{"foreign minted_by", func(e *Envelope) { e.MintedBy = "harness" }, "minted by the control plane"},
		{"missing issued_at", func(e *Envelope) { e.IssuedAt = time.Time{} }, "issued_at timestamp"},
		{"missing tenant", func(e *Envelope) { e.Invocation.TenantID = "" }, "tenant_id is invalid"},
		{"missing subject", func(e *Envelope) { e.Invocation.Subject = "" }, "subject is invalid"},
		{"missing connector", func(e *Envelope) { e.Invocation.ConnectorID = "" }, "connector_id is invalid"},
		{"missing capability", func(e *Envelope) { e.Invocation.Capability = "" }, "capability is invalid"},
		{"missing trace", func(e *Envelope) { e.Invocation.TraceID = "" }, "trace_id is invalid"},
		{"missing resource", func(e *Envelope) { e.Invocation.Resource = "" }, "resource is required"},
		{"invalid action", func(e *Envelope) { e.Invocation.Action = "sudo" }, "invalid action"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := validEnvelope()
			tc.mutate(&env)
			err := env.Validate()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to contain %q", err, tc.want)
			}
		})
	}
}

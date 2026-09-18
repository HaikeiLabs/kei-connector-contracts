package connectors

import "time"

// Client is the Kei-only connector client. It is the single governed entry
// point for a connector call: it mints the typed envelope, runs the fail-closed
// policy decision, and builds the audit attribution. It lives only inside Kei's
// control plane — connectors and third-party harnesses never link it — so an
// envelope is always minted by the platform, never self-asserted.
type Client struct {
	now func() time.Time
}

// InvokeOutcome is the complete governed result of one invocation. The caller
// persists the audit record and, on allow, proceeds to the provider using the
// connector's credential reference — never inline secret material.
type InvokeOutcome struct {
	Envelope    Envelope
	Decision    PolicyDecision
	AuditRecord AuditRecord
}

// NewClient returns a connector client minting envelopes with the real clock.
func NewClient() *Client {
	return &Client{now: time.Now}
}

// Invoke mints the governed envelope for an invocation and evaluates it
// fail-closed. It performs no external I/O: the caller resolves connector
// metadata from its store, invokes the client, and persists the audit record.
func (c *Client) Invoke(m Metadata, in Invocation) InvokeOutcome {
	envelope := Envelope{
		Version:    EnvelopeVersion1,
		Invocation: in,
		MintedBy:   MintedByControlPlane,
		IssuedAt:   c.now().UTC(),
	}
	decision := Decide(m, envelope)
	return InvokeOutcome{
		Envelope:    envelope,
		Decision:    decision,
		AuditRecord: BuildAuditRecord(envelope, decision, envelope.IssuedAt),
	}
}

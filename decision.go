package connectors

import "fmt"

// Decision is the outcome of a governed connector invocation. It is
// deliberately two-state for now; the NIST four-state
// (permit/deny/not_applicable/indeterminate) upgrade is a separate track
// (tool-call-tagging build spec Phase 6).
type Decision string

const (
	DecisionAllow Decision = "allow"
	DecisionDeny  Decision = "deny"
)

// PolicyDecision is the fail-closed outcome of evaluating an envelope against
// a connector's metadata. On an allow it carries the matched capability; on a
// deny it carries a reason. It never carries an error: every invalid boundary
// is a deny.
type PolicyDecision struct {
	Decision   Decision   `json:"decision"`
	Reason     string     `json:"reason,omitempty"`
	Capability Capability `json:"capability,omitempty"`
}

// Decide evaluates a governed envelope against connector metadata. It is
// fail-closed: malformed metadata, an invalid envelope, an inactive connector,
// a cross-tenant/cross-workspace reference, an undefined capability, or an
// undeclared capability are all explicit denies. Only an explicit allow is an
// allow. Data-access and approval policy are decided by the ABAC policy layer,
// not here.
func Decide(m Metadata, env Envelope) PolicyDecision {
	if err := m.Validate(); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: "invalid connector metadata: " + err.Error()}
	}
	if err := env.Validate(); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: "invalid invocation envelope: " + err.Error()}
	}
	if m.Status != StatusActive {
		return PolicyDecision{Decision: DecisionDeny, Reason: "connector is not active"}
	}
	if env.Invocation.TenantID != m.TenantID ||
		env.Invocation.WorkspaceID != m.WorkspaceID ||
		env.Invocation.ConnectorID != m.ID {
		return PolicyDecision{Decision: DecisionDeny, Reason: "connector not found"}
	}
	// The capability must be both defined for the provider and declared on the
	// connector; either boundary failing is a deny.
	if _, ok := LookupCapability(m.Provider, env.Invocation.Capability); !ok {
		return PolicyDecision{Decision: DecisionDeny, Reason: fmt.Sprintf("capability %q is not defined for provider %q", env.Invocation.Capability, m.Provider)}
	}
	if err := ValidateCall(m, env.Invocation); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: err.Error()}
	}
	cap, _ := LookupCapability(m.Provider, env.Invocation.Capability)
	return PolicyDecision{Decision: DecisionAllow, Capability: cap}
}

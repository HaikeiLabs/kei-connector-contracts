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
	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/envelope"
)

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
	Decision   Decision            `json:"decision"`
	Reason     string              `json:"reason,omitempty"`
	Capability contract.Capability `json:"capability,omitempty"`
}

// Decide evaluates a governed envelope against connector metadata. It is
// fail-closed: malformed metadata, an invalid envelope, an inactive connector,
// a cross-tenant/cross-workspace reference, or a malformed or undeclared
// capability are all explicit denies. Only an explicit allow is an allow, and
// it carries the capability the connector declared. Decide is the control
// plane's structural decision: it keeps no capability list
// (contract.AuthorizeCall), and data access is decided by the ABAC policy
// layer, not here. The tenant runtime separately refuses a capability it cannot execute.
func Decide(m contract.Metadata, env envelope.Envelope) PolicyDecision {
	if err := m.Validate(); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: "invalid connector metadata: " + err.Error()}
	}
	if err := env.Validate(); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: "invalid invocation envelope: " + err.Error()}
	}
	if m.Status != contract.StatusActive {
		return PolicyDecision{Decision: DecisionDeny, Reason: "connector is not active"}
	}
	if env.Invocation.TenantID != m.TenantID ||
		env.Invocation.WorkspaceID != m.WorkspaceID ||
		env.Invocation.ConnectorID != m.ID {
		return PolicyDecision{Decision: DecisionDeny, Reason: "connector not found"}
	}
	if err := contract.AuthorizeCall(m, env.Invocation); err != nil {
		return PolicyDecision{Decision: DecisionDeny, Reason: err.Error()}
	}
	// AuthorizeCall has confirmed the connector declares the capability.
	for _, c := range m.Capabilities {
		if c.Name == env.Invocation.Capability {
			return PolicyDecision{Decision: DecisionAllow, Capability: c}
		}
	}
	return PolicyDecision{Decision: DecisionDeny, Reason: "capability is not allowed"}
}

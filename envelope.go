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
	"errors"
	"fmt"
	"time"
)

// EnvelopeVersion is the schema version of the governed invocation envelope.
// Consumers must fail closed on a version they do not recognize.
type EnvelopeVersion string

const (
	// EnvelopeVersion1 is the initial envelope schema. Bump it on a breaking
	// change; an unknown version must never be treated as valid.
	EnvelopeVersion1 EnvelopeVersion = "1.0"

	// MintedByControlPlane is the only authority allowed to mint a governed
	// envelope. Connectors and harnesses can never self-assert one (NIST
	// no-self-assertion): the platform stamps tenant/workspace/agent from the
	// harness key, never from the caller.
	MintedByControlPlane = "kei-control-plane"
)

// Envelope is the typed, governed invocation envelope for a connector call.
// It wraps the contract Invocation with minting provenance, so every call
// carries the full propagation chain — tenant, workspace, subject, agent,
// connector, capability, resource, trace id, and idempotency key — plus who
// minted it and when. The envelope is the unit of authorization and audit.
type Envelope struct {
	Version    EnvelopeVersion `json:"version"`
	Invocation Invocation      `json:"invocation"`
	MintedBy   string          `json:"minted_by"`
	IssuedAt   time.Time       `json:"issued_at"`
}

// Validate is fail-closed: every boundary is explicit and a missing or
// malformed field is an error, never a partial pass.
func (e Envelope) Validate() error {
	if e.Version != EnvelopeVersion1 {
		return fmt.Errorf("unsupported envelope version %q", e.Version)
	}
	if e.MintedBy != MintedByControlPlane {
		return errors.New("envelope must be minted by the control plane")
	}
	if e.IssuedAt.IsZero() {
		return errors.New("envelope is missing an issued_at timestamp")
	}
	return validateInvocation(e.Invocation)
}

// validateInvocation enforces the propagation fields shared by the envelope
// and the contract's ValidateCall, so a malformed envelope is rejected before
// policy evaluation.
func validateInvocation(in Invocation) error {
	for name, value := range map[string]string{
		"tenant_id":    in.TenantID,
		"workspace_id": in.WorkspaceID,
		"subject":      in.Subject,
		"agent_id":     in.AgentID,
		"connector_id": in.ConnectorID,
		"capability":   in.Capability,
		"trace_id":     in.TraceID,
	} {
		if err := validID(name, value); err != nil {
			return err
		}
	}
	if !validAction(in.Action) {
		return fmt.Errorf("invalid action %q", in.Action)
	}
	if in.Resource == "" {
		return errors.New("resource is required")
	}
	if in.IdempotencyKey != "" {
		if err := validID("idempotency_key", in.IdempotencyKey); err != nil {
			return err
		}
	}
	if in.ApprovalID != "" {
		if err := validID("approval_id", in.ApprovalID); err != nil {
			return err
		}
	}
	return nil
}

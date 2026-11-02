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
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/envelope"
	"github.com/HaikeiLabs/kei-connector-contracts/governance"
)

// AuditRecord is the audit-attribution row for a governed connector
// invocation. Its JSON shape mirrors the shared DecisionEvent/AuditRecord wire
// contract (ADR-010) so a record can be stored by the same pipeline; it adds
// workspace_id, which the abac.audit_records table already carries. A record
// never carries credential material — only references and digests.
type AuditRecord struct {
	SpanID           string    `json:"span_id"`
	InvokedAt        time.Time `json:"invoked_at"`
	InvokingSubject  string    `json:"invoking_subject"`
	ParentSpan       *string   `json:"parent_span,omitempty"`
	DelegationDepth  int       `json:"delegation_depth"`
	AgentID          string    `json:"agent_id"`
	AgentVersion     *string   `json:"agent_version,omitempty"`
	Framework        string    `json:"framework"`
	ToolName         string    `json:"tool_name"`
	ToolArgsDigest   *string   `json:"tool_args_digest,omitempty"`
	ResourcesTouched []string  `json:"resources_touched"`
	Decision         string    `json:"decision"`
	PolicyID         *string   `json:"policy_id,omitempty"`
	PolicyReason     *string   `json:"policy_reason,omitempty"`
	OrgID            *string   `json:"org_id,omitempty"`
	WorkspaceID      *string   `json:"workspace_id,omitempty"`
	// CredentialKind is set only on a record of a governed call that used a
	// connector credential: which kind of credential it used. Empty (omitted)
	// on every other record. It never carries the credential, a reference to
	// it, or any status of the agent session that holds it.
	CredentialKind CredentialKind `json:"credential_kind,omitempty"`
}

// CredentialKind is which kind of credential a governed call used.
type CredentialKind string

const (
	// CredentialKindSecretStore is a credential read from the organization's
	// secret store (an opaque_ref connector).
	CredentialKindSecretStore CredentialKind = "secret_store"
	// CredentialKindOAuth is an OAuth credential (an oauth connector).
	CredentialKindOAuth CredentialKind = "oauth"
)

// Valid reports whether k is one of the defined credential kinds.
func (k CredentialKind) Valid() bool {
	return k == CredentialKindSecretStore || k == CredentialKindOAuth
}

// CredentialKindFor maps a connector's credential source to the credential
// kind an audit record carries: opaque_ref (and the empty legacy value) is
// secret_store, oauth is oauth. ok is false for any other source.
func CredentialKindFor(source contract.CredentialSource) (CredentialKind, bool) {
	switch source {
	case "", contract.CredentialSourceOpaqueRef:
		return CredentialKindSecretStore, true
	case contract.CredentialSourceOAuth:
		return CredentialKindOAuth, true
	}
	return "", false
}

// BuildAuditRecord maps an envelope and its fail-closed decision onto the
// canonical audit shape. The span_id is the invocation's dedupe key: the
// idempotency key when present, otherwise the trace id, so a replayed
// invocation never duplicates an audit row. The tool_name encodes the
// connector and capability so connector calls are queryable as a class
// (tool_name LIKE 'connector/%').
func BuildAuditRecord(env envelope.Envelope, decision governance.PolicyDecision, invokedAt time.Time) AuditRecord {
	spanID := env.Invocation.IdempotencyKey
	if spanID == "" {
		spanID = env.Invocation.TraceID
	}
	rec := AuditRecord{
		SpanID:           spanID,
		InvokedAt:        invokedAt,
		InvokingSubject:  env.Invocation.Subject,
		AgentID:          env.Invocation.AgentID,
		Framework:        "connector",
		ToolName:         "connector/" + env.Invocation.ConnectorID + "/" + env.Invocation.Capability,
		ResourcesTouched: []string{env.Invocation.Resource},
		Decision:         string(decision.Decision),
	}
	if env.Invocation.TenantID != "" {
		org := env.Invocation.TenantID
		rec.OrgID = &org
	}
	if env.Invocation.WorkspaceID != "" {
		workspace := env.Invocation.WorkspaceID
		rec.WorkspaceID = &workspace
	}
	if decision.Reason != "" {
		reason := decision.Reason
		rec.PolicyReason = &reason
	}
	return rec
}

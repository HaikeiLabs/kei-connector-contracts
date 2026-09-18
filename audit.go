package connectors

import "time"

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
}

// BuildAuditRecord maps an envelope and its fail-closed decision onto the
// canonical audit shape. The span_id is the invocation's dedupe key: the
// idempotency key when present, otherwise the trace id, so a replayed
// invocation never duplicates an audit row. The tool_name encodes the
// connector and capability so connector calls are queryable as a class
// (tool_name LIKE 'connector/%').
func BuildAuditRecord(env Envelope, decision PolicyDecision, invokedAt time.Time) AuditRecord {
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

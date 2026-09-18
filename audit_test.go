package connectors

import (
	"testing"
	"time"
)

func TestBuildAuditRecordAttribution(t *testing.T) {
	env := decisionEnvelope()
	decision := PolicyDecision{Decision: DecisionAllow, Capability: Capability{Name: "lead.read", Action: ActionRead}}
	invokedAt := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	rec := BuildAuditRecord(env, decision, invokedAt)

	if rec.SpanID != "idem-1" {
		t.Fatalf("span_id = %q, want the idempotency key", rec.SpanID)
	}
	if rec.InvokingSubject != "u-1" {
		t.Fatalf("invoking_subject = %q, want u-1", rec.InvokingSubject)
	}
	if rec.AgentID != "a-1" {
		t.Fatalf("agent_id = %q, want a-1", rec.AgentID)
	}
	if rec.Framework != "connector" {
		t.Fatalf("framework = %q, want connector", rec.Framework)
	}
	if want := "connector/c-1/lead.read"; rec.ToolName != want {
		t.Fatalf("tool_name = %q, want %q", rec.ToolName, want)
	}
	if len(rec.ResourcesTouched) != 1 || rec.ResourcesTouched[0] != "leads/42" {
		t.Fatalf("resources_touched = %v", rec.ResourcesTouched)
	}
	if rec.Decision != string(DecisionAllow) {
		t.Fatalf("decision = %q", rec.Decision)
	}
	if rec.OrgID == nil || *rec.OrgID != "t-1" {
		t.Fatalf("org_id = %v, want t-1", rec.OrgID)
	}
	if rec.WorkspaceID == nil || *rec.WorkspaceID != "w-1" {
		t.Fatalf("workspace_id = %v, want w-1", rec.WorkspaceID)
	}
}

func TestBuildAuditRecordFallsBackToTraceAndRecordsDenyReason(t *testing.T) {
	env := decisionEnvelope()
	env.Invocation.IdempotencyKey = ""
	env.Invocation.TraceID = "trace-9"
	decision := PolicyDecision{Decision: DecisionDeny, Reason: "approval required"}

	rec := BuildAuditRecord(env, decision, time.Now())
	if rec.SpanID != "trace-9" {
		t.Fatalf("span_id = %q, want the trace id when no idempotency key", rec.SpanID)
	}
	if rec.Decision != string(DecisionDeny) {
		t.Fatalf("decision = %q, want deny", rec.Decision)
	}
	if rec.PolicyReason == nil || *rec.PolicyReason != "approval required" {
		t.Fatalf("policy_reason = %v, want the deny reason", rec.PolicyReason)
	}
}

func TestBuildAuditRecordDenyCarriesNoCapability(t *testing.T) {
	rec := BuildAuditRecord(decisionEnvelope(), PolicyDecision{Decision: DecisionDeny, Reason: "connector is not active"}, time.Now())
	if rec.ToolName != "connector/c-1/lead.read" {
		t.Fatalf("tool_name = %q; a deny must still be attributable", rec.ToolName)
	}
}

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
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/governance"
)

func TestBuildAuditRecordAttribution(t *testing.T) {
	env := decisionEnvelope()
	decision := governance.PolicyDecision{Decision: governance.DecisionAllow, Capability: contract.Capability{Name: "lead.read", Action: contract.ActionRead}}
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
	if rec.Decision != string(governance.DecisionAllow) {
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
	decision := governance.PolicyDecision{Decision: governance.DecisionDeny, Reason: "no matching policy"}

	rec := BuildAuditRecord(env, decision, time.Now())
	if rec.SpanID != "trace-9" {
		t.Fatalf("span_id = %q, want the trace id when no idempotency key", rec.SpanID)
	}
	if rec.Decision != string(governance.DecisionDeny) {
		t.Fatalf("decision = %q, want deny", rec.Decision)
	}
	if rec.PolicyReason == nil || *rec.PolicyReason != "no matching policy" {
		t.Fatalf("policy_reason = %v, want the deny reason", rec.PolicyReason)
	}
}

func TestBuildAuditRecordPreservesCanonicalPlainTextValues(t *testing.T) {
	env := decisionEnvelope()
	env.Invocation.Subject = "user:U-17"
	env.Invocation.AgentID = "agent:ReviewBot"
	env.Invocation.ConnectorID = "inst-7"
	env.Invocation.Capability = "pull_request.read"
	env.Invocation.Resource = "repos/Acme/Kei/pulls/17"

	rec := BuildAuditRecord(env, governance.PolicyDecision{Decision: governance.DecisionAllow}, time.Now())
	if rec.InvokingSubject != "user:U-17" || rec.AgentID != "agent:ReviewBot" {
		t.Fatalf("principal values changed: subject=%q agent=%q", rec.InvokingSubject, rec.AgentID)
	}
	if want := "connector/inst-7/pull_request.read"; rec.ToolName != want {
		t.Fatalf("tool_name = %q, want %q", rec.ToolName, want)
	}
	if len(rec.ResourcesTouched) != 1 || rec.ResourcesTouched[0] != "repos/Acme/Kei/pulls/17" {
		t.Fatalf("resources_touched = %v; resource identifier must remain verbatim", rec.ResourcesTouched)
	}
}

func TestBuildAuditRecordDenyCarriesNoCapability(t *testing.T) {
	rec := BuildAuditRecord(decisionEnvelope(), governance.PolicyDecision{Decision: governance.DecisionDeny, Reason: "connector is not active"}, time.Now())
	if rec.ToolName != "connector/c-1/lead.read" {
		t.Fatalf("tool_name = %q; a deny must still be attributable", rec.ToolName)
	}
}

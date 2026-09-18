package connectors

import (
	"strings"
	"testing"
)

// s3Fixture returns an active S3 connector scoped to a single billing bucket
// prefix, mirroring the governed-object-access matrix in
// tests/foundry/fixtures/08-governed-connectors.yaml.
func s3Fixture() Metadata {
	return Metadata{
		ID:            "c-s3-1",
		TenantID:      "t-1",
		WorkspaceID:   "w-1",
		Name:          "billing-s3",
		Provider:      ProviderS3,
		Status:        StatusActive,
		CredentialRef: "vault/tenant/t-1/s3-billing",
		Scopes:        []string{"s3:read"},
		Resources:     []string{"s3://billing-prod/invoices"},
		Policy: PolicyAttributes{
			AllowedActions:   []Action{ActionRead},
			AllowedResources: []string{"s3://billing-prod/invoices"},
		},
		Capabilities: []Capability{
			{Name: "object.read", Action: ActionRead},
			{Name: "object.list", Action: ActionRead},
		},
		CreatedBy: "u-1",
	}
}

// validInvocation is a well-formed governed call for the CRM fixture.
func validInvocation() Invocation {
	return Invocation{
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
}

// TestValidateCallTenantWorkspaceIsolation is the fail-closed guard: a call
// may only ever reach a connector that is in the caller's own tenant,
// workspace, and connector identity. Every boundary is an explicit check.
func TestValidateCallTenantWorkspaceIsolation(t *testing.T) {
	m := fixture()
	in := validInvocation()

	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("valid call rejected: %v", err)
	}

	tests := []struct {
		name string
		mut  func(*Invocation)
	}{
		{"cross-tenant", func(i *Invocation) { i.TenantID = "t-other" }},
		{"cross-workspace", func(i *Invocation) { i.WorkspaceID = "w-other" }},
		{"cross-connector", func(i *Invocation) { i.ConnectorID = "c-other" }},
		{"missing tenant", func(i *Invocation) { i.TenantID = "" }},
		{"missing workspace", func(i *Invocation) { i.WorkspaceID = "" }},
		{"missing connector", func(i *Invocation) { i.ConnectorID = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := validInvocation()
			tt.mut(&i)
			err := ValidateCall(m, i)
			if err == nil {
				t.Fatal("boundary violation was allowed")
			}
			// A missing id is reported as an invalid id; a wrong id is reported
			// as connector not found. Both are hard denials.
			if !strings.Contains(err.Error(), "connector not found") && !strings.Contains(err.Error(), "invalid") {
				t.Fatalf("error = %q, want a hard denial", err.Error())
			}
		})
	}
}

// TestValidateCallRevokedConnectors walks every non-active lifecycle status.
// The contract is fail-closed: only StatusActive may be invoked, so a revoked
// or otherwise disabled connector denies before any policy check runs.
func TestValidateCallRevokedConnectors(t *testing.T) {
	in := validInvocation()

	for _, status := range []Status{StatusPending, StatusSuspended, StatusRevoked, StatusFailed} {
		t.Run(string(status), func(t *testing.T) {
			m := fixture()
			m.Status = status
			if err := ValidateCall(m, in); err == nil || err.Error() != "connector is not active" {
				t.Fatalf("status %q error = %v, want connector is not active", status, err)
			}
		})
	}

	// Recovery control: the same invocation is denied while revoked and allowed
	// the moment the connector is active again.
	m := fixture()
	m.Status = StatusRevoked
	if err := ValidateCall(m, in); err == nil {
		t.Fatal("revoked connector was invocable")
	}
	m.Status = StatusActive
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("reactivated connector rejected: %v", err)
	}
}

// TestValidateCallDoesNotEnforceResourcePolicy documents the decoupling:
// resource/action/prefix policy moved out of the connector contract to the
// ABAC policy layer (governance pivot). ValidateCall still requires a resource
// and rejects an empty one, but no longer consults PolicyAttributes.
func TestValidateCallDoesNotEnforceResourcePolicy(t *testing.T) {
	m := fixture()
	m.Policy.AllowedResources = nil

	i := validInvocation()
	i.Resource = "billing/42"
	if err := ValidateCall(m, i); err != nil {
		t.Fatalf("resource outside the dormant policy was rejected: %v", err)
	}

	i.Resource = ""
	if err := ValidateCall(m, i); err == nil || err.Error() != "resource is required" {
		t.Fatalf("empty resource error = %v, want resource is required", err)
	}
}

// TestConnectorPolicyFieldsAreDormant documents the governance pivot: the
// connector's PolicyAttributes are no longer enforced by the contract; data
// access and tool-call authorization are decided by ABAC policies.
func TestConnectorPolicyFieldsAreDormant(t *testing.T) {
	m := s3Fixture()
	m.Policy.AllowedPrefixes = nil
	m.Policy.AllowedResources = nil
	m.Policy.AllowedActions = nil

	i := validInvocation()
	i.ConnectorID = m.ID
	i.Capability = "object.read"
	i.Action = ActionRead
	i.Resource = "s3://other-bucket/secret.json"

	if err := ValidateCall(m, i); err != nil {
		t.Fatalf("dormant policy enforced by the contract: %v", err)
	}
}

// TestValidateCredentialRefRejectsMaterial is the redaction boundary at the
// contract level: a credential_ref must be an opaque reference. URLs and
// inline secret-bearing strings are rejected so provider credentials stay in
// the credential/secret backend.
func TestValidateCredentialRefRejectsMaterial(t *testing.T) {
	rejected := []string{
		"",
		"https://example.com/token",
		"postgres://user:pass@db.example.com/prod",
		"vault\nsecret",
		"AKIA\tIOSFODNN7EXAMPLE",
	}
	for _, ref := range rejected {
		t.Run(ref, func(t *testing.T) {
			if err := ValidateCredentialRef(ref); err == nil {
				t.Errorf("accepted unsafe credential_ref %q", ref)
			}
		})
	}

	accepted := []string{
		"vault/tenant/t-1/crm",
		"secret/kei/connectors/s3-billing",
		"kv-v2/data/kei/github",
	}
	for _, ref := range accepted {
		t.Run(ref, func(t *testing.T) {
			if err := ValidateCredentialRef(ref); err != nil {
				t.Errorf("rejected opaque credential_ref %q: %v", ref, err)
			}
		})
	}

	t.Run("overlong reference rejected", func(t *testing.T) {
		ref := "vault/" + strings.Repeat("a", 512)
		if err := ValidateCredentialRef(ref); err == nil {
			t.Error("accepted credential_ref longer than 512 bytes")
		}
	})
}

func TestCredentialRefRejectsSecretShapedMaterial(t *testing.T) {
	for _, tc := range []struct {
		ref     string
		wantErr bool
	}{
		{"AKIAIOSFODNN7EXAMPLE", false},
		{"ghp_12345678901234567890", true},
	} {
		if err := ValidateCredentialRef(tc.ref); (err != nil) != tc.wantErr {
			t.Errorf("ValidateCredentialRef(%q) error = %v, wantErr %v", tc.ref, err, tc.wantErr)
		}
	}
}

// TestValidateCallDelegationAttribution requires every governed call to be
// attributable: subject, agent, and trace are mandatory identity fields.
func TestValidateCallDelegationAttribution(t *testing.T) {
	missing := []struct {
		name string
		mut  func(*Invocation)
		want string
	}{
		{"subject", func(i *Invocation) { i.Subject = "" }, "subject is invalid"},
		{"agent", func(i *Invocation) { i.AgentID = "" }, "agent_id is invalid"},
		{"trace", func(i *Invocation) { i.TraceID = "" }, "trace_id is invalid"},
		{"connector", func(i *Invocation) { i.ConnectorID = "" }, "connector_id is invalid"},
	}
	for _, tt := range missing {
		t.Run(tt.name, func(t *testing.T) {
			i := validInvocation()
			tt.mut(&i)
			if err := ValidateCall(fixture(), i); err == nil || err.Error() != tt.want {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}

	t.Run("metadata requires a creator", func(t *testing.T) {
		m := fixture()
		m.CreatedBy = ""
		if err := m.Validate(); err == nil {
			t.Fatal("connector metadata without created_by was accepted")
		}
	})
}

// TestValidateCallIdempotencyKeyStable documents the contract-level idempotency
// posture: ValidateCall is a pure function, so an invocation replayed with the
// same idempotency key produces the identical decision. Replay deduplication
// itself must be implemented by the invoking harness or the control plane; the
// contract has no side effects to dedupe (see docs/connectors-governance-report.md).
func TestValidateCallIdempotencyKeyStable(t *testing.T) {
	m := fixture()

	first := validInvocation()
	first.IdempotencyKey = "idem-replay-1"
	second := first

	if err := ValidateCall(m, first); err != nil {
		t.Fatalf("first invocation rejected: %v", err)
	}
	if err := ValidateCall(m, second); err != nil {
		t.Fatalf("replayed invocation rejected: %v", err)
	}

	t.Run("idempotency key not required", func(t *testing.T) {
		i := validInvocation()
		i.IdempotencyKey = ""
		if err := ValidateCall(m, i); err != nil {
			t.Fatalf("invocation without idempotency key rejected: %v", err)
		}
	})
}

// TestValidateCallDiscordSourcedDenialAndRecovery models the governed denial
// and recovery paths for a Discord-sourced agent (harness key bound to the
// discord service). The denial comes from connector lifecycle and tenant
// boundary enforcement, not from any Discord-specific field.
func TestValidateCallDiscordSourcedDenialAndRecovery(t *testing.T) {
	in := validInvocation()

	t.Run("revoked connector denies discord agent", func(t *testing.T) {
		m := fixture()
		m.Status = StatusRevoked
		if err := ValidateCall(m, in); err == nil || err.Error() != "connector is not active" {
			t.Fatalf("error = %v, want connector is not active", err)
		}
	})

	t.Run("cross-tenant denies discord agent", func(t *testing.T) {
		m := fixture()
		i := in
		i.TenantID = "t-other"
		if err := ValidateCall(m, i); err == nil || err.Error() != "connector not found" {
			t.Fatalf("error = %v, want connector not found", err)
		}
	})

	t.Run("reactivated connector recovers", func(t *testing.T) {
		m := fixture()
		m.Status = StatusSuspended
		if err := ValidateCall(m, in); err == nil {
			t.Fatal("suspended connector was invocable")
		}
		m.Status = StatusActive
		if err := ValidateCall(m, in); err != nil {
			t.Fatalf("reactivated connector rejected: %v", err)
		}
	})
}

// TestValidateCallDoesNotConsultPrefixes documents that the connector's
// PolicyAttributes.AllowedPrefixes field is fully dormant on the contract: it
// is not consulted by ValidateCall. Prefix boundary enforcement belongs to the
// ABAC policy layer (governance pivot), which supersedes GAP-CONN-002.
func TestValidateCallDoesNotConsultPrefixes(t *testing.T) {
	m := s3Fixture()
	m.Policy.AllowedPrefixes = []string{"s3://billing-prod/invoices"}

	i := validInvocation()
	i.ConnectorID = m.ID
	i.Capability = "object.read"
	i.Resource = "s3://billing-prod/invoices/2026/01/invoice-1001.json"

	if err := ValidateCall(m, i); err != nil {
		t.Fatalf("dormant AllowedPrefixes enforced by the contract: %v", err)
	}
}

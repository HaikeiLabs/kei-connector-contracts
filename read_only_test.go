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
	"testing"
	"time"
)

// The connector catalog is a credential surface, not an authorization grant
// (docs/connector-contract-freeze.md). It deliberately names provider mutation
// capabilities so a connector's credential scope can be described, while ADR-011
// puts the live allow/deny decision in the tenant-side runtime PEP.
//
// These tests pin the operational consequence of that split: a connector is
// read-only unless its own metadata explicitly opts in to a mutation. A
// connector that declares only reads can never be talked into a write, no
// matter what the caller asks for or which approval id it presents.
//
// This is the fail-closed guarantee that keeps "the catalog lists writes" from
// silently becoming "connectors can write".
var catalogMutations = []struct {
	provider   Provider
	capability string
	action     Action
}{
	{ProviderGitHub, "issue.create", ActionCreate},
	{ProviderGitHub, "issue.update", ActionUpdate},
	{ProviderGitHub, "pull_request.create", ActionCreate},
	{ProviderGitHub, "pull_request.update", ActionUpdate},
	{ProviderGitHub, "issue.comment", ActionComment},
	{ProviderCRM, "lead.create", ActionCreate},
	{ProviderCRM, "lead.update", ActionUpdate},
	{ProviderLinear, "issue.create", ActionCreate},
	{ProviderLinear, "issue.update", ActionUpdate},
}

// readOnlyConnector returns an active connector whose metadata declares only a
// read capability, plus that capability and an in-scope resource.
func readOnlyConnector(t *testing.T, provider Provider) (Metadata, string, string) {
	t.Helper()
	m := fixture()
	m.Provider = provider
	switch provider {
	case ProviderCRM:
		m.Capabilities = []Capability{{Name: "lead.read", Action: ActionRead}}
		m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"leads"}}
		return m, "lead.read", "leads/42"
	case ProviderGitHub:
		m.Capabilities = []Capability{{Name: "issue.read", Action: ActionRead}}
		m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"repos/acme/kei"}}
		return m, "issue.read", "repos/acme/kei/issues/7"
	case ProviderLinear:
		m.Capabilities = []Capability{{Name: "issue.read", Action: ActionRead}}
		m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"linear/issue"}}
		return m, "issue.read", "linear/issue/KEI-42"
	case ProviderFreshBooks:
		m.Capabilities = []Capability{{Name: "invoice.read", Action: ActionRead}}
		m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"accounts/acct-1"}}
		return m, "invoice.read", "accounts/acct-1/invoices/inv-1"
	case ProviderMercury:
		m.Capabilities = []Capability{{Name: "account.read", Action: ActionRead}}
		m.Policy = PolicyAttributes{AllowedActions: []Action{ActionRead}, AllowedResources: []string{"accounts/acct-1"}}
		return m, "account.read", "accounts/acct-1"
	default:
		t.Fatalf("unexpected provider %q", provider)
		return Metadata{}, "", ""
	}
}

// TestReadOnlyConnectorRejectsCatalogMutations proves the opt-in boundary: a
// connector declaring only a read capability rejects every mutation the
// provider catalog names, even though those names resolve in the catalog.
func TestReadOnlyConnectorRejectsCatalogMutations(t *testing.T) {
	for _, tc := range catalogMutations {
		t.Run(string(tc.provider)+"/"+tc.capability, func(t *testing.T) {
			// The capability exists in the catalog: this test is about the
			// connector's declared surface, not about the name being unknown.
			if _, ok := LookupCapability(tc.provider, tc.capability); !ok {
				t.Fatalf("precondition: %q should be in the %q catalog", tc.capability, tc.provider)
			}

			meta, _, resource := readOnlyConnector(t, tc.provider)
			in := Invocation{
				TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
				AgentID: "a-1", ConnectorID: meta.ID, Capability: tc.capability,
				Action: tc.action, Resource: resource, TraceID: "trace-1",
			}
			if err := ValidateCall(meta, in); err == nil {
				t.Fatalf("read-only connector accepted mutation %q", tc.capability)
			}
		})
	}
}

// TestApprovalIDCannotUnlockUndeclaredMutation proves an approval id is not a
// bypass: presenting one, with a policy permissive enough to admit the action,
// still fails because the connector never declared the capability.
func TestApprovalIDCannotUnlockUndeclaredMutation(t *testing.T) {
	for _, tc := range catalogMutations {
		t.Run(string(tc.provider)+"/"+tc.capability, func(t *testing.T) {
			meta, _, resource := readOnlyConnector(t, tc.provider)
			meta.Policy.AllowedActions = []Action{ActionRead, ActionCreate, ActionUpdate, ActionComment}

			in := Invocation{
				TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
				AgentID: "a-1", ConnectorID: meta.ID, Capability: tc.capability,
				Action: tc.action, Resource: resource, TraceID: "trace-1",
				ApprovalID: "approval-1",
			}
			if err := ValidateCall(meta, in); err == nil {
				t.Fatalf("mutation %q was authorized with an approval id", tc.capability)
			}
		})
	}
}

// TestDecideDeniesUndeclaredMutation is the end-to-end fail-closed proof
// through the policy decision: a mutation envelope against a read-only
// connector is always a deny, never an allow, regardless of approval.
func TestDecideDeniesUndeclaredMutation(t *testing.T) {
	for _, tc := range catalogMutations {
		t.Run(string(tc.provider)+"/"+tc.capability, func(t *testing.T) {
			meta, _, resource := readOnlyConnector(t, tc.provider)
			meta.Policy.AllowedActions = []Action{ActionRead, ActionCreate, ActionUpdate, ActionComment}

			env := Envelope{
				Version: EnvelopeVersion1,
				Invocation: Invocation{
					TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
					AgentID: "a-1", ConnectorID: meta.ID, Capability: tc.capability,
					Action: tc.action, Resource: resource, TraceID: "trace-1",
					ApprovalID: "approval-1",
				},
				MintedBy: MintedByControlPlane,
				IssuedAt: time.Now().UTC(),
			}
			if d := Decide(meta, env); d.Decision != DecisionDeny {
				t.Fatalf("mutation %q decision = %s (%s), want deny", tc.capability, d.Decision, d.Reason)
			}
		})
	}
}

// TestReadCallNeedsNoApprovalID proves the read path stays usable: a governed
// read against a read-only connector is valid with an empty approval id.
func TestReadCallNeedsNoApprovalID(t *testing.T) {
	for _, provider := range []Provider{ProviderCRM, ProviderGitHub, ProviderLinear} {
		t.Run(string(provider), func(t *testing.T) {
			meta, readCap, resource := readOnlyConnector(t, provider)
			in := Invocation{
				TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
				AgentID: "a-1", ConnectorID: meta.ID, Capability: readCap,
				Action: ActionRead, Resource: resource, TraceID: "trace-1",
			}
			if err := ValidateCall(meta, in); err != nil {
				t.Fatalf("read call without an approval id rejected: %v", err)
			}
		})
	}
}

// TestDeclaredMutationRemainsReachable is the counterpart guarantee: opting in
// is what makes a write reachable. A connector that declares the mutation and
// whose policy admits the action passes validation. Without this, the tests
// above could pass because writes are broken rather than because they are
// gated.
func TestDeclaredMutationRemainsReachable(t *testing.T) {
	for _, tc := range catalogMutations {
		t.Run(string(tc.provider)+"/"+tc.capability, func(t *testing.T) {
			meta, _, resource := readOnlyConnector(t, tc.provider)
			meta.Capabilities = []Capability{{Name: tc.capability, Action: tc.action}}
			meta.Policy.AllowedActions = []Action{tc.action}

			in := Invocation{
				TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
				AgentID: "a-1", ConnectorID: meta.ID, Capability: tc.capability,
				Action: tc.action, Resource: resource, TraceID: "trace-1",
				ApprovalID: "approval-1",
			}
			if err := ValidateCall(meta, in); err != nil {
				t.Fatalf("declared mutation %q rejected: %v", tc.capability, err)
			}
		})
	}
}

// financeProviders are the providers whose catalogs deliberately contain no
// mutation at all.
var financeProviders = []Provider{ProviderFreshBooks, ProviderMercury}

// financeMutations are operations a caller might plausibly ask a bookkeeping
// or banking agent to perform. None is defined for either provider, which is
// the point: the guarantee for money movement is stronger than the opt-in
// default the tests above pin. For GitHub or Linear a connector may opt in to
// a write; for FreshBooks and Mercury there is no capability to opt in to.
var financeMutations = []struct {
	capability string
	action     Action
}{
	{"payment.create", ActionCreate},
	{"invoice.create", ActionCreate},
	{"invoice.update", ActionUpdate},
	{"expense.update", ActionUpdate},
	{"transfer.create", ActionCreate},
	{"transaction.create", ActionCreate},
	{"recipient.create", ActionCreate},
	{"account.update", ActionUpdate},
}

// TestFinanceCatalogsDefineNoMutation proves money movement is absent from the
// catalog rather than merely disabled by policy. A capability that is not
// defined cannot be declared by a connector, so no configuration mistake can
// make it reachable.
func TestFinanceCatalogsDefineNoMutation(t *testing.T) {
	for _, provider := range financeProviders {
		for _, capability := range CapabilitiesFor(provider) {
			if capability.Action != ActionRead {
				t.Errorf("%s capability %q has action %q; finance providers are read-only",
					provider, capability.Name, capability.Action)
			}
		}
		for _, tc := range financeMutations {
			if _, ok := LookupCapability(provider, tc.capability); ok {
				t.Errorf("%s defines mutation capability %q; finance providers are read-only",
					provider, tc.capability)
			}
		}
	}
}

// TestFinanceMutationsRejectedEvenWhenDeclared is the belt-and-braces case: a
// connector that forges a mutation capability in its own metadata, with a
// permissive policy and an approval id, is still refused. ValidateCall checks
// the capability against the provider catalog, so metadata cannot invent a
// payment.
func TestFinanceMutationsRejectedEvenWhenDeclared(t *testing.T) {
	for _, provider := range financeProviders {
		for _, tc := range financeMutations {
			t.Run(string(provider)+"/"+tc.capability, func(t *testing.T) {
				meta, _, resource := readOnlyConnector(t, provider)
				meta.Capabilities = append(meta.Capabilities, Capability{Name: tc.capability, Action: tc.action})
				meta.Policy.AllowedActions = []Action{ActionRead, ActionCreate, ActionUpdate, ActionDelete}
				meta.Policy.DestructiveEnabled = true

				in := Invocation{
					TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
					AgentID: "a-1", ConnectorID: meta.ID, Capability: tc.capability,
					Action: tc.action, Resource: resource, TraceID: "trace-1",
					ApprovalID: "approval-1",
				}
				if err := ValidateCall(meta, in); err == nil {
					t.Fatalf("%s accepted forged mutation %q", provider, tc.capability)
				}
			})
		}
	}
}

// TestFinanceReadsRemainReachable guards against the tests above passing
// because finance connectors are broken rather than because writes are
// absent. Every declared read capability must still validate.
func TestFinanceReadsRemainReachable(t *testing.T) {
	for _, provider := range financeProviders {
		meta, capability, resource := readOnlyConnector(t, provider)
		in := Invocation{
			TenantID: meta.TenantID, WorkspaceID: meta.WorkspaceID, Subject: "u-1",
			AgentID: "a-1", ConnectorID: meta.ID, Capability: capability,
			Action: ActionRead, Resource: resource, TraceID: "trace-1",
		}
		if err := ValidateCall(meta, in); err != nil {
			t.Fatalf("%s read %q rejected: %v", provider, capability, err)
		}
	}
}

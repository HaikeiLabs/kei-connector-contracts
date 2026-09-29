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
)

// These mirror the governance package's test fixtures for tests here.
func fixture() contract.Metadata {
	return contract.Metadata{ID: "c-1", TenantID: "t-1", WorkspaceID: "w-1", Name: "CRM", Provider: contract.ProviderCRM, Status: contract.StatusActive, CredentialRef: "vault/tenant/t-1/crm", Scopes: []string{"leads:read"}, Resources: []string{"leads/*"}, Policy: contract.PolicyAttributes{AllowedActions: []contract.Action{contract.ActionRead, contract.ActionUpdate}, AllowedResources: []string{"leads"}}, Capabilities: []contract.Capability{{Name: "lead.read", Action: contract.ActionRead}}, CreatedBy: "u-1"}
}

func decisionFixture() contract.Metadata {
	m := fixture()
	m.Capabilities = []contract.Capability{
		{Name: "lead.read", Action: contract.ActionRead},
		{Name: "lead.update", Action: contract.ActionUpdate},
	}
	return m
}
func decisionEnvelope() envelope.Envelope {
	return envelope.Envelope{
		Version: envelope.EnvelopeVersion1,
		Invocation: contract.Invocation{
			TenantID:       "t-1",
			WorkspaceID:    "w-1",
			Subject:        "u-1",
			AgentID:        "a-1",
			ConnectorID:    "c-1",
			Capability:     "lead.read",
			Action:         contract.ActionRead,
			Resource:       "leads/42",
			TraceID:        "trace-1",
			IdempotencyKey: "idem-1",
		},
		MintedBy: envelope.MintedByControlPlane,
		IssuedAt: time.Now().UTC(),
	}
}

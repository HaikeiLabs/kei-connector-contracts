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

import "github.com/HaikeiLabs/kei-connector-contracts/contract"

// fixture mirrors contract/contract_test.go's fixture for tests in this package.
func fixture() contract.Metadata {
	return contract.Metadata{ID: "c-1", TenantID: "t-1", WorkspaceID: "w-1", Name: "CRM", Provider: contract.ProviderCRM, Status: contract.StatusActive, CredentialRef: "vault/tenant/t-1/crm", Scopes: []string{"leads:read"}, Resources: []string{"leads/*"}, Policy: contract.PolicyAttributes{AllowedActions: []contract.Action{contract.ActionRead, contract.ActionUpdate}, AllowedResources: []string{"leads"}}, Capabilities: []contract.Capability{{Name: "lead.read", Action: contract.ActionRead}}, CreatedBy: "u-1"}
}

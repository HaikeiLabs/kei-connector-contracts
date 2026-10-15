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

package contract

import "testing"

func TestOnBehalfOfValidate(t *testing.T) {
	valid := OnBehalfOf{
		TenantID:        "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		ClientID:        "11111111-2222-3333-4444-555555555555",
		Resource:        "https://analysis.windows.net/powerbi/api",
		Scopes:          []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
		ClientSecretRef: "kei://connectors/c-http/client_secret",
	}
	if err := valid.Validate("c-http"); err != nil {
		t.Fatalf("valid OnBehalfOf rejected: %v", err)
	}

	for name, tc := range map[string]struct {
		obo OnBehalfOf
		id  string
	}{
		"tenant_id not a guid": {
			obo: OnBehalfOf{
				TenantID: "not-a-guid", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"client_id not a guid": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "not-a-guid",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"resource is http": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "http://analysis.windows.net/powerbi/api", Scopes: []string{"http://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"scopes is empty": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: nil,
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"scope uses .default": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/.default"},
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"scope under different resource": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://graph.microsoft.com/User.Read"},
				ClientSecretRef: "kei://connectors/c-http/client_secret",
			},
			id: "c-http",
		},
		"wrong connector id in ref": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "kei://connectors/other-connector/client_secret",
			},
			id: "c-http",
		},
		"wrong field in ref": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "kei://connectors/c-http/api_key",
			},
			id: "c-http",
		},
		"literal secret value": {
			obo: OnBehalfOf{
				TenantID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ClientID: "11111111-2222-3333-4444-555555555555",
				Resource: "https://analysis.windows.net/powerbi/api", Scopes: []string{"https://analysis.windows.net/powerbi/api/Dataset.Read.All"},
				ClientSecretRef: "sk_literal_secret_value",
			},
			id: "c-http",
		},
	} {
		if err := tc.obo.Validate(tc.id); err == nil {
			t.Errorf("%s: was accepted", name)
		}
	}
}

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

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"
)

var updateSetupGolden = flag.Bool("update-setup", false, "rewrite schemas/connector-setup.v1.json")

// Only the connectors we build are configurable: s3, notion, and the finance
// providers have no setup schema.
func TestSetupSchemasCoverExactlyTheSupportedProviders(t *testing.T) {
	want := []Provider{ProviderGmail, ProviderGoogle, ProviderLinear, ProviderGitHub, ProviderTito, ProviderCRM, ProviderHTTPAPI}
	got := SetupSchemas()
	if len(got) != len(want) {
		t.Fatalf("SetupSchemas() has %d providers, want %d", len(got), len(want))
	}
	for i, p := range want {
		if got[i].Provider != p || got[i].Schema != SetupSchemaVersion {
			t.Fatalf("schema %d = %s/%s, want %s/%s", i, got[i].Provider, got[i].Schema, p, SetupSchemaVersion)
		}
	}
	for _, p := range []Provider{ProviderS3, ProviderNotion, ProviderFreshBooks, ProviderMercury} {
		if _, ok := SetupSchemaFor(p); ok {
			t.Errorf("%s has a setup schema", p)
		}
	}
}

// The OAuth choices a setup screen offers must be exactly the account models
// the contract allows for that provider, with per_user as the default.
func TestSetupAuthOptionsMatchAccountModels(t *testing.T) {
	for _, s := range SetupSchemas() {
		var offered []AccountModel
		for _, a := range s.Auth {
			offered = append(offered, a.AccountModels...)
		}
		allowed := AccountModelsFor(s.Provider)
		if len(offered) != len(allowed) {
			t.Fatalf("%s offers %v, contract allows %v", s.Provider, offered, allowed)
		}
		for i := range allowed {
			if offered[i] != allowed[i] {
				t.Fatalf("%s offers %v, contract allows %v", s.Provider, offered, allowed)
			}
		}
		if len(allowed) > 0 && s.Auth[0].DefaultAccountModel != AccountModelPerUser {
			t.Errorf("%s default account model = %q, want per_user", s.Provider, s.Auth[0].DefaultAccountModel)
		}
	}
}

// Each connector has at most one secret field per account model, because a
// connector has one credential_ref.
func TestSetupSchemasDeclareAtMostOneSecretPerModel(t *testing.T) {
	for _, s := range SetupSchemas() {
		models := AccountModelsFor(s.Provider)
		if len(models) == 0 {
			models = []AccountModel{""}
		}
		for _, model := range models {
			secrets := 0
			for _, f := range s.FieldsFor(model) {
				if f.Secret {
					secrets++
				}
			}
			if secrets > 1 {
				t.Errorf("%s/%s declares %d secret fields", s.Provider, model, secrets)
			}
		}
	}
}

func TestValidateConfigAcceptsValidNonSecretFields(t *testing.T) {
	cases := []struct {
		provider Provider
		model    AccountModel
		config   map[string]any
	}{
		{ProviderTito, "", map[string]any{"account_slug": "acme"}},
		{ProviderGmail, AccountModelPerUser, map[string]any{"include_body": true}},
		{ProviderGmail, AccountModelPerUser, map[string]any{}},
		{ProviderGmail, AccountModelDomainDelegation, map[string]any{"impersonate_email": "events@example.com"}},
		{ProviderGoogle, AccountModelShared, map[string]any{"drive_id": "0AFx-drive_1"}},
		{ProviderCRM, "", map[string]any{"base_url": "https://haikeilabs.com", "assertion_audience": "kei-crm", "allowed_resources": []any{"leads"}}},
		{ProviderLinear, AccountModelShared, nil},
	}
	for _, tc := range cases {
		if err := ValidateConfig(tc.provider, tc.model, tc.config); err != nil {
			t.Errorf("%s/%s config rejected: %v", tc.provider, tc.model, err)
		}
	}
}

func TestValidateConfigRejectsInvalidConfigWithoutEchoingValues(t *testing.T) {
	cases := []struct {
		name     string
		provider Provider
		model    AccountModel
		config   map[string]any
	}{
		{"missing required", ProviderTito, "", map[string]any{}},
		{"bad pattern", ProviderTito, "", map[string]any{"account_slug": "has spaces-VALUE"}},
		{"secret field in config", ProviderTito, "", map[string]any{"account_slug": "acme", "api_token": "tok-VALUE-123456789012"}},
		{"unknown field", ProviderTito, "", map[string]any{"account_slug": "acme", "color": "VALUE"}},
		{"wrong type", ProviderGmail, AccountModelPerUser, map[string]any{"include_body": "VALUE"}},
		{"field for another model", ProviderGmail, AccountModelPerUser, map[string]any{"impersonate_email": "VALUE@example.com"}},
		{"not an email", ProviderGmail, AccountModelDomainDelegation, map[string]any{"impersonate_email": "Name <VALUE@example.com>"}},
		{"http base url", ProviderCRM, "", map[string]any{"base_url": "http://VALUE.example.com", "assertion_audience": "kei-crm"}},
		{"url with userinfo", ProviderCRM, "", map[string]any{"base_url": "https://u:VALUE@example.com", "assertion_audience": "kei-crm"}},
		{"model not allowed", ProviderGitHub, AccountModelDomainDelegation, map[string]any{}},
		{"no schema", ProviderS3, "", map[string]any{}},
	}
	for _, tc := range cases {
		err := ValidateConfig(tc.provider, tc.model, tc.config)
		if err == nil {
			t.Errorf("%s: accepted", tc.name)
			continue
		}
		if strings.Contains(err.Error(), "VALUE") {
			t.Errorf("%s: error echoes the submitted value: %v", tc.name, err)
		}
	}
}

// schemas/connector-setup.v1.json is the export the web UI and CLI render. It
// must always equal the Go data; regenerate with -update-setup.
func TestSetupSchemaJSONExportIsCurrent(t *testing.T) {
	want, err := json.MarshalIndent(SetupSchemas(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want = append(want, '\n')
	const path = "../schemas/connector-setup.v1.json"
	if *updateSetupGolden {
		if err := os.WriteFile(path, want, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (run go test -run TestSetupSchemaJSONExportIsCurrent -update-setup)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is stale; run go test -run TestSetupSchemaJSONExportIsCurrent -update-setup", path)
	}
}

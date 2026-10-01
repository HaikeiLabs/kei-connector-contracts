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

package setup

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

var updateSetupGolden = flag.Bool("update-setup", false, "rewrite schemas/connector-setup.v1.json")

// Only the connectors we build are configurable: s3 and the finance
// providers have no setup schema.
func TestSetupSchemasCoverExactlyTheSupportedProviders(t *testing.T) {
	want := []contract.Provider{contract.ProviderGmail, contract.ProviderGoogle, contract.ProviderLinear, contract.ProviderGitHub, contract.ProviderTito, contract.ProviderNotion, contract.ProviderDiscord, contract.ProviderGrafana, contract.ProviderCRM, contract.ProviderHTTPAPI}
	got := SetupSchemas()
	if len(got) != len(want) {
		t.Fatalf("SetupSchemas() has %d providers, want %d", len(got), len(want))
	}
	for i, p := range want {
		if got[i].Provider != p || got[i].Schema != SetupSchemaVersion {
			t.Fatalf("schema %d = %s/%s, want %s/%s", i, got[i].Provider, got[i].Schema, p, SetupSchemaVersion)
		}
	}
	for _, p := range []contract.Provider{contract.ProviderS3, contract.ProviderFreshBooks, contract.ProviderMercury} {
		if _, ok := SetupSchemaFor(p); ok {
			t.Errorf("%s has a setup schema", p)
		}
	}
}

// The OAuth choices a setup screen offers must be exactly the account models
// the contract allows for that provider, with per_user as the default.
func TestSetupAuthOptionsMatchAccountModels(t *testing.T) {
	for _, s := range SetupSchemas() {
		var offered []contract.AccountModel
		for _, a := range s.Auth {
			offered = append(offered, a.AccountModels...)
		}
		allowed := contract.AccountModelsFor(s.Provider)
		if len(offered) != len(allowed) {
			t.Fatalf("%s offers %v, contract allows %v", s.Provider, offered, allowed)
		}
		for i := range allowed {
			if offered[i] != allowed[i] {
				t.Fatalf("%s offers %v, contract allows %v", s.Provider, offered, allowed)
			}
		}
		if len(allowed) > 0 && s.Auth[0].DefaultAccountModel != contract.AccountModelPerUser {
			t.Errorf("%s default account model = %q, want per_user", s.Provider, s.Auth[0].DefaultAccountModel)
		}
	}
}

// Each connector has at most one secret field per account model, because a
// connector has one credential_ref.
func TestSetupSchemasDeclareAtMostOneSecretPerModel(t *testing.T) {
	for _, s := range SetupSchemas() {
		models := contract.AccountModelsFor(s.Provider)
		if len(models) == 0 {
			models = []contract.AccountModel{""}
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

// Notion is set up with an internal integration token and nothing else: the
// pages the integration can read are the ones shared with it in Notion.
func TestNotionSetupIsOneIntegrationToken(t *testing.T) {
	s, ok := SetupSchemaFor(contract.ProviderNotion)
	if !ok || len(s.Auth) != 1 || s.Auth[0].CredentialSource != contract.CredentialSourceOpaqueRef || len(s.Auth[0].AccountModels) != 0 {
		t.Fatalf("notion schema = %+v", s)
	}
	if len(s.Fields) != 1 || s.Fields[0].Name != "api_token" || !s.Fields[0].Secret || !s.Fields[0].Required || s.Fields[0].Location != SetupLocationCredential {
		t.Fatalf("notion fields = %+v", s.Fields)
	}
	token := s.Fields[0]
	for _, v := range []string{"ntn_" + strings.Repeat("a1B2", 11), "secret_" + strings.Repeat("Zz9", 14)} {
		if !validSetupString(token, v) {
			t.Errorf("token %q rejected", v)
		}
	}
	for _, v := range []string{"", "ntn_", "Bearer ntn_" + strings.Repeat("a", 40), "ntn_" + strings.Repeat("a", 20) + " x", "sk_" + strings.Repeat("a", 40), "ntn_" + strings.Repeat("a", 300)} {
		if validSetupString(token, v) {
			t.Errorf("token %q accepted", v)
		}
	}
}

func TestValidateConfigAcceptsValidNonSecretFields(t *testing.T) {
	cases := []struct {
		provider contract.Provider
		model    contract.AccountModel
		config   map[string]any
	}{
		{contract.ProviderTito, "", map[string]any{"account_slug": "acme"}},
		{contract.ProviderNotion, "", nil},
		{contract.ProviderNotion, "", map[string]any{}},
		{contract.ProviderGmail, contract.AccountModelPerUser, map[string]any{"include_body": true}},
		{contract.ProviderGmail, contract.AccountModelPerUser, map[string]any{}},
		{contract.ProviderGmail, contract.AccountModelDomainDelegation, map[string]any{"impersonate_email": "events@example.com"}},
		{contract.ProviderGoogle, contract.AccountModelShared, map[string]any{"drive_id": "0AFx-drive_1"}},
		{contract.ProviderCRM, "", map[string]any{"base_url": "https://haikeilabs.com", "assertion_audience": "kei-crm", "allowed_resources": []any{"leads"}}},
		{contract.ProviderLinear, contract.AccountModelShared, nil},
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
		provider contract.Provider
		model    contract.AccountModel
		config   map[string]any
	}{
		{"missing required", contract.ProviderTito, "", map[string]any{}},
		{"bad pattern", contract.ProviderTito, "", map[string]any{"account_slug": "has spaces-VALUE"}},
		{"secret field in config", contract.ProviderTito, "", map[string]any{"account_slug": "acme", "api_token": "tok-VALUE-123456789012"}},
		{"unknown field", contract.ProviderTito, "", map[string]any{"account_slug": "acme", "color": "VALUE"}},
		{"wrong type", contract.ProviderGmail, contract.AccountModelPerUser, map[string]any{"include_body": "VALUE"}},
		{"field for another model", contract.ProviderGmail, contract.AccountModelPerUser, map[string]any{"impersonate_email": "VALUE@example.com"}},
		{"not an email", contract.ProviderGmail, contract.AccountModelDomainDelegation, map[string]any{"impersonate_email": "Name <VALUE@example.com>"}},
		{"http base url", contract.ProviderCRM, "", map[string]any{"base_url": "http://VALUE.example.com", "assertion_audience": "kei-crm"}},
		{"url with userinfo", contract.ProviderCRM, "", map[string]any{"base_url": "https://u:VALUE@example.com", "assertion_audience": "kei-crm"}},
		{"model not allowed", contract.ProviderGitHub, contract.AccountModelDomainDelegation, map[string]any{}},
		{"no schema", contract.ProviderS3, "", map[string]any{}},
		{"notion token in config", contract.ProviderNotion, "", map[string]any{"api_token": "ntn_VALUE0000000000000000"}},
		{"notion has no config fields", contract.ProviderNotion, "", map[string]any{"workspace_id": "VALUE"}},
		{"notion has no account models", contract.ProviderNotion, contract.AccountModelShared, map[string]any{}},
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

func connectorMetadata(provider contract.Provider, source contract.CredentialSource, model contract.AccountModel, config map[string]any) contract.Metadata {
	return contract.Metadata{ID: "c-1", TenantID: "t-1", WorkspaceID: "w-1", Name: "conn", Provider: provider, Status: contract.StatusActive,
		CredentialSource: source, CredentialRef: "vault/tenant/t-1/conn", Scopes: []string{"read"},
		Capabilities: contract.CapabilitiesFor(provider)[:1], CreatedBy: "u-1", AccountModel: model, Config: config}
}

// ValidateMetadata is the full connector check: the contract's structural
// and account-model rules first, then the setup schema's config rules.
func TestValidateMetadataChecksStructureThenConfig(t *testing.T) {
	delegation := connectorMetadata(contract.ProviderGmail, contract.CredentialSourceOpaqueRef, contract.AccountModelDomainDelegation, map[string]any{"impersonate_email": "events@example.com"})
	if err := ValidateMetadata(delegation); err != nil {
		t.Fatalf("valid delegation connector rejected: %v", err)
	}
	missing := delegation
	missing.Config = map[string]any{}
	if err := ValidateMetadata(missing); err == nil {
		t.Error("delegation connector without impersonate_email was accepted")
	}
	legacy := connectorMetadata(contract.ProviderTito, contract.CredentialSourceOpaqueRef, "", nil)
	if err := ValidateMetadata(legacy); err != nil {
		t.Errorf("legacy connector without config rejected: %v", err)
	}
	legacy.Config = map[string]any{"account_slug": "has spaces"}
	if err := ValidateMetadata(legacy); err == nil {
		t.Error("legacy connector with invalid config was accepted")
	}
	structural := delegation
	structural.Provider = contract.ProviderGitHub
	if err := ValidateMetadata(structural); err == nil {
		t.Error("structurally invalid connector was accepted")
	}
}

// TestValidateMetadataRequiresExecutableCapabilities: ValidateMetadata is the
// runtime's check, so it refuses a capability the runtime has no code for,
// even though the control plane's structural Metadata.Validate accepts it.
func TestValidateMetadataRequiresExecutableCapabilities(t *testing.T) {
	m := connectorMetadata(contract.ProviderCRM, contract.CredentialSourceOpaqueRef, "", nil)
	m.Capabilities = []contract.Capability{{Name: "deal.read", Action: contract.ActionRead}}
	if err := m.Validate(); err != nil {
		t.Fatalf("structural check rejected a declared capability: %v", err)
	}
	if err := ValidateMetadata(m); err == nil || err.Error() != `capability "deal.read" is not defined for provider "crm"` {
		t.Fatalf("ValidateMetadata error = %v, want the executable-capability error", err)
	}
}

// TestCRMAllowedResourcesIncludeInvestors (HAI-210): the crm contract defines
// investor.list and investor.read, so the setup schema must let a connector
// scope itself to investors alongside leads. The default stays leads and
// customers, and resources outside the crm surface are still refused.
func TestCRMAllowedResourcesIncludeInvestors(t *testing.T) {
	base := func(resources ...any) map[string]any {
		return map[string]any{"base_url": "https://haikeilabs.com", "assertion_audience": "kei-crm", "allowed_resources": resources}
	}
	for name, config := range map[string]map[string]any{
		"investors with leads": base("leads", "investors"),
		"investors wildcard":   base("investors/*"),
	} {
		if err := ValidateConfig(contract.ProviderCRM, "", config); err != nil {
			t.Errorf("%s rejected: %v", name, err)
		}
	}
	for name, config := range map[string]map[string]any{
		"unknown resource":   base("leads", "invoices"),
		"investor singular":  base("investor"),
		"investors sub-path": base("investors/42"),
	} {
		if err := ValidateConfig(contract.ProviderCRM, "", config); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	schema, _ := SetupSchemaFor(contract.ProviderCRM)
	for _, f := range schema.Fields {
		if f.Name == "allowed_resources" {
			if got := fmt.Sprint(f.Default); got != "[leads customers]" {
				t.Errorf("allowed_resources default = %s, want [leads customers]", got)
			}
		}
	}
}

// HAI-309: Discord is set up with a bot token and nothing else; the guilds it
// can read are the ones the bot was invited to.
func TestDiscordSetupIsOneBotToken(t *testing.T) {
	s, ok := SetupSchemaFor(contract.ProviderDiscord)
	if !ok || len(s.Auth) != 1 || s.Auth[0].CredentialSource != contract.CredentialSourceOpaqueRef || len(s.Auth[0].AccountModels) != 0 {
		t.Fatalf("discord schema = %+v", s)
	}
	if len(s.Fields) != 1 || s.Fields[0].Name != "bot_token" || !s.Fields[0].Secret || !s.Fields[0].Required || s.Fields[0].Location != SetupLocationCredential {
		t.Fatalf("discord fields = %+v", s.Fields)
	}
	token := s.Fields[0]
	valid := "MTA4" + strings.Repeat("x", 20) + ".GhIj9k." + strings.Repeat("Ab_-", 8)
	if !validSetupString(token, valid) {
		t.Errorf("token %q rejected", valid)
	}
	for _, v := range []string{"", "short", "Bot " + valid, valid + " x", strings.Repeat("a", 300)} {
		if validSetupString(token, v) {
			t.Errorf("token %q accepted", v)
		}
	}
	if err := ValidateConfig(contract.ProviderDiscord, "", map[string]any{"bot_token": "VALUE"}); err == nil {
		t.Error("discord bot_token in config was accepted")
	}
}

// HAI-309: Grafana takes its instance URL (stored in Metadata.Grafana, never
// in config) and a service-account token behind the connector's opaque_ref.
func TestGrafanaSetupIsBaseURLAndServiceAccountToken(t *testing.T) {
	s, ok := SetupSchemaFor(contract.ProviderGrafana)
	if !ok || len(s.Auth) != 1 || s.Auth[0].CredentialSource != contract.CredentialSourceOpaqueRef || len(s.Auth[0].AccountModels) != 0 {
		t.Fatalf("grafana schema = %+v", s)
	}
	if len(s.Fields) != 2 {
		t.Fatalf("grafana fields = %+v", s.Fields)
	}
	base, token := s.Fields[0], s.Fields[1]
	if base.Name != "base_url" || base.Secret || !base.Required || base.Type != SetupFieldHTTPSURL || base.Location != SetupLocationGrafana {
		t.Errorf("grafana base_url field = %+v", base)
	}
	if token.Name != "service_account_token" || !token.Secret || !token.Required || token.Location != SetupLocationCredential {
		t.Errorf("grafana token field = %+v", token)
	}
	valid := "glsa_" + strings.Repeat("aB3", 10) + "_1a2b3c4d"
	if !validSetupString(token, valid) {
		t.Errorf("token %q rejected", valid)
	}
	for _, v := range []string{"", "glsa_", "Bearer " + valid, "eyJrIjoi" + strings.Repeat("a", 40), valid + " x"} {
		if validSetupString(token, v) {
			t.Errorf("token %q accepted", v)
		}
	}
	if err := ValidateConfig(contract.ProviderGrafana, "", map[string]any{"base_url": "https://VALUE.example.com"}); err == nil {
		t.Error("grafana base_url in config was accepted")
	}
}

// ValidateMetadata enforces the Grafana base URL rules on the connector, and
// that discord and grafana connectors are executable.
func TestValidateMetadataChecksDiscordAndGrafana(t *testing.T) {
	discord := connectorMetadata(contract.ProviderDiscord, contract.CredentialSourceOpaqueRef, "", nil)
	discord.Capabilities = contract.CapabilitiesFor(contract.ProviderDiscord)
	if err := ValidateMetadata(discord); err != nil {
		t.Errorf("discord connector rejected: %v", err)
	}
	grafana := connectorMetadata(contract.ProviderGrafana, contract.CredentialSourceOpaqueRef, "", nil)
	grafana.Capabilities = contract.CapabilitiesFor(contract.ProviderGrafana)
	grafana.Grafana = &contract.GrafanaConfig{BaseURL: "https://grafana.internal:3000"}
	if err := ValidateMetadata(grafana); err != nil {
		t.Errorf("self-hosted grafana connector rejected: %v", err)
	}
	for _, u := range []string{"http://grafana.example.com", "https://169.254.169.254", "https://metadata.google.internal", "https://grafana.example.com/api"} {
		bad := grafana
		bad.Grafana = &contract.GrafanaConfig{BaseURL: u}
		if err := ValidateMetadata(bad); err == nil {
			t.Errorf("grafana connector with base_url %q accepted", u)
		}
	}
	missing := grafana
	missing.Grafana = nil
	if err := ValidateMetadata(missing); err == nil {
		t.Error("grafana connector without a base url was accepted")
	}
}

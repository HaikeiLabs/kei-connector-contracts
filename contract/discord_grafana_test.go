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

func grafanaMetadata(baseURL string) Metadata {
	m := frozenMetadata(ProviderGrafana, CapabilitiesFor(ProviderGrafana))
	m.Grafana = &GrafanaConfig{BaseURL: baseURL}
	return m
}

// HAI-309: Discord and Grafana are read-only shared-secret providers. A bot
// token or service-account token sits behind an opaque_ref; neither has OAuth.
func TestDiscordAndGrafanaAreReadOnlyWithOpaqueCredential(t *testing.T) {
	for _, m := range []Metadata{frozenMetadata(ProviderDiscord, CapabilitiesFor(ProviderDiscord)), grafanaMetadata("https://grafana.example.com")} {
		if !ValidProvider(m.Provider) {
			t.Fatalf("%s is not a valid provider", m.Provider)
		}
		if err := m.Validate(); err != nil {
			t.Fatalf("%s metadata rejected: %v", m.Provider, err)
		}
		if len(CapabilitiesFor(m.Provider)) == 0 {
			t.Fatalf("%s declares no capabilities", m.Provider)
		}
		for _, c := range CapabilitiesFor(m.Provider) {
			if c.Action != ActionRead {
				t.Errorf("%s capability %q has action %q, want read", m.Provider, c.Name, c.Action)
			}
		}
		oauth := m
		oauth.CredentialSource, oauth.Subject = CredentialSourceOAuth, "subject-1"
		if err := oauth.Validate(); err == nil {
			t.Errorf("%s connector with an oauth credential source was accepted", m.Provider)
		}
		if len(AccountModelsFor(m.Provider)) != 0 {
			t.Errorf("%s has account models", m.Provider)
		}
	}
	for _, write := range []Capability{
		{Name: "message.create", Action: ActionCreate},
		{Name: "member.list", Action: ActionRead},
		{Name: "member.read", Action: ActionRead},
		{Name: "channel.read", Action: ActionUpdate},
	} {
		if err := frozenMetadata(ProviderDiscord, []Capability{write}).ValidateExecutable(); err == nil {
			t.Errorf("discord capability %q/%q was accepted", write.Name, write.Action)
		}
	}
	for _, write := range []Capability{
		{Name: "dashboard.create", Action: ActionCreate},
		{Name: "alert_rule.update", Action: ActionUpdate},
		{Name: "annotation.create", Action: ActionCreate},
	} {
		if err := grafanaMetadata("https://grafana.example.com").ValidateExecutable(); err != nil {
			t.Fatalf("grafana metadata not executable: %v", err)
		}
		m := grafanaMetadata("https://grafana.example.com")
		m.Capabilities = []Capability{write}
		if err := m.ValidateExecutable(); err == nil {
			t.Errorf("grafana capability %q/%q was accepted", write.Name, write.Action)
		}
	}
}

// The Grafana base URL is required on a grafana connector and refused on any
// other provider, the same shape as the http_api block.
func TestGrafanaConfigIsRequiredOnlyForGrafana(t *testing.T) {
	missing := frozenMetadata(ProviderGrafana, CapabilitiesFor(ProviderGrafana))
	if err := missing.Validate(); err == nil {
		t.Error("grafana connector without grafana metadata was accepted")
	}
	other := frozenMetadata(ProviderDiscord, CapabilitiesFor(ProviderDiscord))
	other.Grafana = &GrafanaConfig{BaseURL: "https://grafana.example.com"}
	if err := other.Validate(); err == nil {
		t.Error("grafana metadata on a discord connector was accepted")
	}
}

// Grafana is often self-hosted, so private and internal hosts are allowed.
// Link-local and cloud-metadata destinations never are.
func TestGrafanaBaseURL(t *testing.T) {
	for _, u := range []string{
		"https://grafana.example.com",
		"https://acme.grafana.net",
		"https://grafana.example.com/",
		"https://grafana.example.com:3000",
		"https://10.0.4.12",
		"https://192.168.1.20:3000",
		"https://grafana.internal",
		"https://grafana.corp.local",
		"https://localhost:3000",
		"https://127.0.0.1:3000",
		"https://[fd12:3456::1]",
		"https://grafana.cafe",
		"https://0x.example.com",
	} {
		if err := (GrafanaConfig{BaseURL: u}).Validate(); err != nil {
			t.Errorf("base_url %q rejected: %v", u, err)
		}
		if err := grafanaMetadata(u).Validate(); err != nil {
			t.Errorf("metadata with base_url %q rejected: %v", u, err)
		}
	}
	for _, u := range []string{
		"",
		"http://grafana.example.com",
		"grafana.example.com",
		"https://",
		"https://grafana.example.com/grafana",
		"https://grafana.example.com/?orgId=1",
		"https://grafana.example.com#x",
		"https://admin:pw@grafana.example.com",
		"https://169.254.169.254",
		"https://169.254.10.1:3000",
		"https://[fe80::1]",
		"https://[fd00:ec2::254]",
		"https://0.0.0.0",
		"https://metadata",
		"https://metadata.google.internal",
		"https://METADATA.google.internal",
		"https://instance-data",
		"https://instance-data.ec2.internal",
		"https://metadata.azure.com",
		"https://[::ffff:169.254.169.254]",
		"https://2852039166",
		"https://0xa9fea9fe",
		"https://0251.0376.0251.0376",
		"https://169.254.169.254.",
	} {
		if err := (GrafanaConfig{BaseURL: u}).Validate(); err == nil {
			t.Errorf("base_url %q accepted", u)
		}
		if err := grafanaMetadata(u).Validate(); err == nil {
			t.Errorf("metadata with base_url %q accepted", u)
		}
	}
}

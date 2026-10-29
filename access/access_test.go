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

package access

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

var updateAccessGolden = flag.Bool("update-access", false, "rewrite schemas/connector-access.v1.json")

// adapterProviders are the providers with a kei-connector-runtime adapter
// (internal/connector/invoke/execute.go). Only these get an api mode.
var adapterProviders = []contract.Provider{
	contract.ProviderGmail, contract.ProviderGoogle, contract.ProviderLinear, contract.ProviderGitHub,
	contract.ProviderTito, contract.ProviderNotion, contract.ProviderDiscord, contract.ProviderGrafana,
	contract.ProviderHTTPAPI, contract.ProviderSkillsRepo,
}

// HAI-319 freeze: every adapter provider has exactly one api mode, and no
// provider has a cli, mcp, or webhook mode until the owner decides one.
func TestAccessFreezeIsAPIOnlyForAdapterProviders(t *testing.T) {
	got := Access()
	if len(got) != len(adapterProviders) {
		t.Fatalf("Access() has %d providers, want %d", len(got), len(adapterProviders))
	}
	for i, p := range adapterProviders {
		a := got[i]
		if a.Provider != p || a.Schema != AccessSchemaVersion {
			t.Fatalf("access %d = %s/%s, want %s/%s", i, a.Provider, a.Schema, p, AccessSchemaVersion)
		}
		if len(a.Access) != 1 || a.Access[0].Kind != KindAPI || a.Access[0].API == nil || a.Access[0].API.Adapter != p {
			t.Fatalf("%s access = %+v, want one api mode", p, a.Access)
		}
		if err := a.Validate(); err != nil {
			t.Errorf("%s access invalid: %v", p, err)
		}
	}
	for _, p := range []contract.Provider{contract.ProviderS3, contract.ProviderFreshBooks, contract.ProviderMercury} {
		if _, ok := AccessFor(p); ok {
			t.Errorf("%s has access modes but no runtime adapter", p)
		}
	}
}

// The api mode reads every read capability and declares only the writes the
// runtime executes: Linear issue.create and comment.create, and http_api
// http.post and http.patch.
func TestAPIModeReadsAndWrites(t *testing.T) {
	for _, p := range adapterProviders {
		a, _ := AccessFor(p)
		mode := a.Access[0]
		var reads, writes []string
		for _, c := range contract.CapabilitiesFor(p) {
			if c.Action == contract.ActionRead {
				reads = append(reads, c.Name)
			} else if APIWrite(p, c.Name) {
				writes = append(writes, c.Name)
			}
		}
		if !reflect.DeepEqual(sorted(mode.Reads), sorted(reads)) {
			t.Errorf("%s reads = %v, want %v", p, mode.Reads, reads)
		}
		if !reflect.DeepEqual(sorted(mode.Writes), sorted(writes)) {
			t.Errorf("%s writes = %v, want %v", p, mode.Writes, writes)
		}
		if len(mode.Credentials) != 0 {
			t.Errorf("%s api mode wires credentials", p)
		}
		if len(mode.Setup) == 0 || len(mode.Verify) == 0 {
			t.Errorf("%s api mode has no setup or verify steps", p)
		}
	}
	if _, ok := AccessFor("crm"); ok {
		t.Error("crm has access modes; there is no crm provider (HAI-384)")
	}
	linear, _ := AccessFor(contract.ProviderLinear)
	if want := []string{"comment.create", "issue.create"}; !reflect.DeepEqual(sorted(linear.Access[0].Writes), want) {
		t.Errorf("linear writes = %v, want %v", linear.Access[0].Writes, want)
	}
	httpAPI, _ := AccessFor(contract.ProviderHTTPAPI)
	if want := []string{"http.patch", "http.post"}; !reflect.DeepEqual(sorted(httpAPI.Access[0].Writes), want) {
		t.Errorf("http_api writes = %v, want %v", httpAPI.Access[0].Writes, want)
	}
}

func TestAccessForReturnsACopy(t *testing.T) {
	a, _ := AccessFor(contract.ProviderGitHub)
	a.Access[0].Reads[0] = "issue.delete"
	b, _ := AccessFor(contract.ProviderGitHub)
	if b.Access[0].Reads[0] == "issue.delete" {
		t.Fatal("AccessFor shares its data with the caller")
	}
}

func TestCredentialReference(t *testing.T) {
	if got := CredentialReference("conn-1", "token"); got != "kei://connectors/conn-1/token" {
		t.Fatalf("CredentialReference = %q", got)
	}
}

func validCLI() ProviderAccess {
	return ProviderAccess{Schema: AccessSchemaVersion, Provider: contract.ProviderGitHub, Access: []AccessMode{{
		Kind:        KindCLI,
		CLI:         &CLIMode{Binary: "gh", MinVersion: "2.40.0", Install: "brew install gh"},
		Credentials: []CredentialEnv{{Env: "GH_TOKEN", Field: "token"}},
		Setup:       []Step{{Name: "install", Description: "Install the GitHub CLI."}},
		Verify:      []Step{{Name: "version", Description: "Check the version.", Run: []string{"gh", "--version"}}},
		Reads:       []string{"repository.read", "issue.read"},
		Writes:      []string{"issue.create"},
	}}}
}

// The non-api kinds are schema only: the shapes validate, though no
// provider declares one yet.
func TestValidateAcceptsEveryKind(t *testing.T) {
	sse := validCLI()
	sse.Access[0] = AccessMode{Kind: KindMCP, MCP: &MCPMode{Transport: TransportSSE, URL: "https://mcp.example.com/sse"},
		Credentials: []CredentialEnv{{Env: "GITHUB_TOKEN", Field: "token"}}, Reads: []string{"issue.read"}}
	remote := validCLI()
	remote.Access[0] = AccessMode{Kind: KindMCP, MCP: &MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com/mcp"}, Reads: []string{"issue.read"}}
	hook := validCLI()
	hook.Access[0] = AccessMode{Kind: KindWebhook, Webhook: &WebhookMode{Direction: DirectionInbound, SecretField: "webhook_secret"}, Reads: []string{"issue.read"}}
	all := validCLI()
	all.Access = append(all.Access, sse.Access[0], hook.Access[0], AccessMode{Kind: KindAPI, API: &APIMode{Adapter: contract.ProviderGitHub}, Reads: []string{"issue.read"}})
	for name, a := range map[string]ProviderAccess{"cli": validCLI(), "mcp sse": sse, "mcp http": remote, "webhook": hook, "all kinds": all} {
		if err := a.Validate(); err != nil {
			t.Errorf("%s rejected: %v", name, err)
		}
	}
}

func TestValidateRejectsMalformedModes(t *testing.T) {
	cases := map[string]func(*ProviderAccess){
		"schema":   func(a *ProviderAccess) { a.Schema = "kei.connector-access/v0" },
		"provider": func(a *ProviderAccess) { a.Provider = "atlassian" },
		"kind":     func(a *ProviderAccess) { a.Access[0].Kind = "ssh" },
		"block for other kind": func(a *ProviderAccess) {
			a.Access[0].MCP = &MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com"}
		},
		"missing block":  func(a *ProviderAccess) { a.Access[0].CLI = nil },
		"duplicate kind": func(a *ProviderAccess) { a.Access = append(a.Access, a.Access[0]) },
		"binary path":    func(a *ProviderAccess) { a.Access[0].CLI.Binary = "/usr/bin/gh" },
		"min version":    func(a *ProviderAccess) { a.Access[0].CLI.MinVersion = "latest" },
		"install hint":   func(a *ProviderAccess) { a.Access[0].CLI.Install = "" },
		"unknown read":   func(a *ProviderAccess) { a.Access[0].Reads = []string{"repo.clone"} },
		"write as read":  func(a *ProviderAccess) { a.Access[0].Reads = []string{"issue.create"} },
		"read as write":  func(a *ProviderAccess) { a.Access[0].Writes = []string{"issue.read"} },
		"duplicate read": func(a *ProviderAccess) { a.Access[0].Reads = []string{"issue.read", "issue.read"} },
		"env name":       func(a *ProviderAccess) { a.Access[0].Credentials[0].Env = "gh-token" },
		"field name":     func(a *ProviderAccess) { a.Access[0].Credentials[0].Field = "kei://connectors/c/token" },
		"duplicate env": func(a *ProviderAccess) {
			a.Access[0].Credentials = append(a.Access[0].Credentials, a.Access[0].Credentials[0])
		},
		"step without name": func(a *ProviderAccess) { a.Access[0].Setup[0].Name = "" },
		"step without text": func(a *ProviderAccess) { a.Access[0].Verify[0].Description = "" },
		"duplicate step":    func(a *ProviderAccess) { a.Access[0].Setup = append(a.Access[0].Setup, a.Access[0].Setup[0]) },
		"secret in step": func(a *ProviderAccess) {
			a.Access[0].Verify[0].Run = []string{"gh", "auth", "login", "--with-token=ghp_abc"}
		},
		"secret in description": func(a *ProviderAccess) { a.Access[0].Setup[0].Description = "export GH_TOKEN=ghp_abc" },
		"api with credentials": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindAPI, API: &APIMode{Adapter: contract.ProviderGitHub}, Credentials: []CredentialEnv{{Env: "GH_TOKEN", Field: "token"}}}
		},
		"api adapter for another provider": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindAPI, API: &APIMode{Adapter: contract.ProviderLinear}}
		},
		"webhook with credentials": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindWebhook, Webhook: &WebhookMode{Direction: DirectionInbound, SecretField: "s"}, Credentials: a.Access[0].Credentials}
		},
		"webhook direction": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindWebhook, Webhook: &WebhookMode{Direction: "both", SecretField: "s"}}
		},
		"webhook secret field": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindWebhook, Webhook: &WebhookMode{Direction: DirectionOutbound}}
		},
		"mcp transport": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindMCP, MCP: &MCPMode{Transport: "websocket", URL: "https://mcp.example.com"}}
		},
		"mcp without url": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindMCP, MCP: &MCPMode{Transport: TransportStreamableHTTP}}
		},
		"mcp without https url": func(a *ProviderAccess) {
			a.Access[0] = AccessMode{Kind: KindMCP, MCP: &MCPMode{Transport: TransportStreamableHTTP, URL: "http://mcp.example.com"}}
		},
	}
	for name, mutate := range cases {
		a := validCLI()
		mutate(&a)
		if err := a.Validate(); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

// HAI-484: MCP connectors are remote only (streamable_http or sse over
// https). stdio is gone, and the url must be a clean https endpoint.
func TestMCPModeIsRemoteOnly(t *testing.T) {
	doc := func(m *MCPMode) ProviderAccess {
		a := validCLI()
		a.Access[0] = AccessMode{Kind: KindMCP, MCP: m, Reads: []string{"issue.read"}}
		return a
	}
	t.Run("stdio transport is rejected", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: "stdio", URL: "https://mcp.example.com"}).Validate(); err == nil {
			t.Fatal("stdio transport accepted, want rejected")
		}
	})
	t.Run("http url is rejected", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportStreamableHTTP, URL: "http://mcp.example.com"}).Validate(); err == nil {
			t.Fatal("http url accepted, want rejected")
		}
	})
	t.Run("userinfo url is rejected", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportStreamableHTTP, URL: "https://user:pass@mcp.example.com"}).Validate(); err == nil {
			t.Fatal("userinfo url accepted, want rejected")
		}
	})
	t.Run("query url is rejected", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com/mcp?x=1"}).Validate(); err == nil {
			t.Fatal("query url accepted, want rejected")
		}
	})
	t.Run("fragment url is rejected", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com/mcp#f"}).Validate(); err == nil {
			t.Fatal("fragment url accepted, want rejected")
		}
	})
	t.Run("sse https is accepted", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportSSE, URL: "https://mcp.example.com/sse"}).Validate(); err != nil {
			t.Fatalf("sse https rejected: %v", err)
		}
	})
	t.Run("streamable_http https is accepted", func(t *testing.T) {
		if err := doc(&MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com/mcp"}).Validate(); err != nil {
			t.Fatalf("streamable_http https rejected: %v", err)
		}
	})
	t.Run("resource defaults to url", func(t *testing.T) {
		m := &MCPMode{Transport: TransportStreamableHTTP, URL: "https://mcp.example.com/mcp"}
		if err := doc(m).Validate(); err != nil {
			t.Fatalf("validate: %v", err)
		}
		if m.Resource != m.URL {
			t.Fatalf("resource = %q, want %q", m.Resource, m.URL)
		}
	})
}

// schemas/connector-access.v1.json is the export the console, CLI, and
// skills read. It must equal the Go data; regenerate with -update-access.
func TestAccessJSONExportIsCurrent(t *testing.T) {
	want, err := json.MarshalIndent(Access(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want = append(want, '\n')
	const path = "../schemas/connector-access.v1.json"
	if *updateAccessGolden {
		if err := os.WriteFile(path, want, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (run go test ./access -run TestAccessJSONExportIsCurrent -update-access)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is stale; run go test ./access -run TestAccessJSONExportIsCurrent -update-access", path)
	}
}

// The JSON Schema's enums and field names must match the Go types.
func TestAccessJSONSchemaMatchesGoTypes(t *testing.T) {
	raw, err := os.ReadFile("../schemas/connector-access.v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Defs map[string]struct {
			Properties map[string]struct {
				Const string   `json:"const"`
				Enum  []string `json:"enum"`
			} `json:"properties"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("schema is not JSON: %v", err)
	}
	enum := func(def, prop string) []string {
		return sorted(schema.Defs[def].Properties[prop].Enum)
	}
	if got := schema.Defs["provider_access"].Properties["schema"].Const; got != AccessSchemaVersion {
		t.Errorf("schema const = %q, want %q", got, AccessSchemaVersion)
	}
	for _, tc := range []struct {
		def, prop string
		want      []string
	}{
		{"access_mode", "kind", []string{string(KindAPI), string(KindCLI), string(KindMCP), string(KindWebhook)}},
		{"mcp", "transport", []string{string(TransportStreamableHTTP), string(TransportSSE)}},
		{"webhook", "direction", []string{string(DirectionInbound), string(DirectionOutbound)}},
	} {
		if got := enum(tc.def, tc.prop); !reflect.DeepEqual(got, sorted(tc.want)) {
			t.Errorf("%s.%s enum = %v, want %v", tc.def, tc.prop, got, tc.want)
		}
	}
	for def, typ := range map[string]any{
		"provider_access": ProviderAccess{}, "access_mode": AccessMode{}, "api": APIMode{}, "cli": CLIMode{},
		"mcp": MCPMode{}, "webhook": WebhookMode{}, "credential": CredentialEnv{}, "step": Step{},
	} {
		var fields []string
		rt := reflect.TypeOf(typ)
		for i := 0; i < rt.NumField(); i++ {
			fields = append(fields, strings.Split(rt.Field(i).Tag.Get("json"), ",")[0])
		}
		var props []string
		for name := range schema.Defs[def].Properties {
			props = append(props, name)
		}
		if !reflect.DeepEqual(sorted(props), sorted(fields)) {
			t.Errorf("%s properties = %v, Go fields = %v", def, sorted(props), sorted(fields))
		}
	}
}

func sorted(values []string) []string {
	out := append([]string{}, values...)
	sort.Strings(out)
	if out == nil {
		return []string{}
	}
	return out
}

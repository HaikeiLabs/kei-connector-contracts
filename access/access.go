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

// Package access describes how each provider can be reached: through the
// runtime adapter (api), a provider CLI, an MCP server, or a webhook. It is
// metadata only and grants nothing. ABAC policy still decides every call, a
// connector's capabilities still come from contract definitions, and the
// owner decides which modes a connector actually uses. schemas/
// connector-access.v1.json is its JSON export and
// schemas/connector-access.v1.schema.json its JSON Schema.
package access

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// AccessSchemaVersion identifies the provider access schema.
const AccessSchemaVersion = "kei.connector-access/v1"

// Kind is how a mode reaches the provider.
type Kind string

const (
	// KindAPI: the kei-connector-runtime adapter, run by
	// `kei-proxy connector invoke`. The runtime resolves the connector's
	// credential_ref itself.
	KindAPI Kind = "api"
	// KindCLI: the provider's own CLI, run under `kei-proxy run --`
	// (HAI-305), which injects credential references into the child only.
	KindCLI Kind = "cli"
	// KindMCP: an MCP server, started by the MCP launcher, which injects
	// credential references into the server only.
	KindMCP Kind = "mcp"
	// KindWebhook: provider events delivered to, or sent from, Kei.
	KindWebhook Kind = "webhook"
)

// Transport is an MCP server transport.
type Transport string

const (
	TransportStdio          Transport = "stdio"
	TransportStreamableHTTP Transport = "streamable_http"
)

// Direction is which way a webhook flows.
type Direction string

const (
	// DirectionInbound: the provider calls Kei.
	DirectionInbound Direction = "inbound"
	// DirectionOutbound: Kei calls the provider's endpoint.
	DirectionOutbound Direction = "outbound"
)

// ProviderAccess lists a provider's access modes.
type ProviderAccess struct {
	Schema   string            `json:"schema"`
	Provider contract.Provider `json:"provider"`
	Access   []AccessMode      `json:"access"`
}

// AccessMode is one way to reach a provider. Exactly the block that matches
// Kind is set.
type AccessMode struct {
	Kind    Kind         `json:"kind"`
	API     *APIMode     `json:"api,omitempty"`
	CLI     *CLIMode     `json:"cli,omitempty"`
	MCP     *MCPMode     `json:"mcp,omitempty"`
	Webhook *WebhookMode `json:"webhook,omitempty"`
	// Credentials wire credential fields into a cli or mcp process's
	// environment as references (see CredentialReference), never values.
	Credentials []CredentialEnv `json:"credentials,omitempty"`
	// Setup and Verify are ordered.
	Setup  []Step `json:"setup,omitempty"`
	Verify []Step `json:"verify,omitempty"`
	// Reads are the read capabilities this mode serves.
	Reads []string `json:"reads"`
	// Writes are the write capabilities this mode serves. Writes are made
	// by agent action tools (ADR-028), never implied by a read.
	Writes []string `json:"writes,omitempty"`
}

// APIMode runs the provider's runtime adapter.
type APIMode struct {
	Adapter contract.Provider `json:"adapter"`
}

// CLIMode runs a provider CLI.
type CLIMode struct {
	// Binary is a bare executable name, never a path.
	Binary     string `json:"binary"`
	MinVersion string `json:"min_version"`
	// Install is a human-readable install hint.
	Install string `json:"install"`
}

// MCPMode runs or connects to an MCP server.
type MCPMode struct {
	Package   string    `json:"package,omitempty"`
	Transport Transport `json:"transport"`
	// Command is the stdio server's argv.
	Command []string `json:"command,omitempty"`
	// URL is the streamable_http endpoint (https).
	URL string `json:"url,omitempty"`
}

// WebhookMode is a provider webhook.
type WebhookMode struct {
	Direction Direction `json:"direction"`
	// SecretField names the credential field holding the signing secret.
	SecretField string `json:"secret_field"`
}

// CredentialEnv sets Env to CredentialReference(<connector id>, Field) in a
// cli or mcp process.
type CredentialEnv struct {
	Env   string `json:"env"`
	Field string `json:"field"`
}

// Step is one ordered setup or verify step.
type Step struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Run         []string `json:"run,omitempty"`
}

// CredentialReference is the env value `kei-proxy run` and the MCP launcher
// resolve into the child process: kei://connectors/<connector_id>/<field>.
func CredentialReference(connectorID, field string) string {
	return "kei://connectors/" + connectorID + "/" + field
}

// apiWrites are the write capabilities the runtime executes through its
// adapter, per provider. crm executes every write it defines; Linear executes
// issue.create and comment.create. Other providers' write capabilities, and
// Linear issue.update, are not served by api mode.
var apiWrites = map[contract.Provider]map[string]bool{
	contract.ProviderLinear: {"issue.create": true, "comment.create": true},
}

// APIWrite reports whether the runtime adapter executes the named write
// capability for provider.
func APIWrite(provider contract.Provider, capability string) bool {
	if provider == contract.ProviderCRM {
		return true
	}
	return apiWrites[provider][capability]
}

func apiMode(provider contract.Provider) AccessMode {
	mode := AccessMode{
		Kind: KindAPI,
		API:  &APIMode{Adapter: provider},
		Setup: []Step{
			{Name: "create_connector", Description: "Create the connector in the Kei console with the provider's setup fields (schemas/connector-setup.v1.json). Secrets go to the console-connected secret manager; the connector keeps only its credential_ref."},
			{Name: "declare_capabilities", Description: "Declare the capabilities the connector exposes. ABAC policy decides who may use them."},
		},
		Reads: []string{},
	}
	for _, c := range contract.CapabilitiesFor(provider) {
		if c.Action == contract.ActionRead {
			mode.Reads = append(mode.Reads, c.Name)
		} else if APIWrite(provider, c.Name) {
			mode.Writes = append(mode.Writes, c.Name)
		}
	}
	mode.Verify = []Step{{
		Name:        "invoke_read",
		Description: "Run one read through the runtime adapter.",
		Run:         []string{"kei-proxy", "connector", "invoke", "--connector", "<connector_id>", "--capability", mode.Reads[0], "--action", "read", "--resource", "<resource>"},
	}}
	return mode
}

func apiOnly(provider contract.Provider) ProviderAccess {
	return ProviderAccess{Schema: AccessSchemaVersion, Provider: provider, Access: []AccessMode{apiMode(provider)}}
}

// providerAccess lists, in display order, the providers with a runtime
// adapter (kei-connector-runtime internal/connector/invoke/execute.go). Each
// has only its api mode: cli, mcp, and webhook modes are the owner's call.
var providerAccess = []ProviderAccess{
	apiOnly(contract.ProviderGmail),
	apiOnly(contract.ProviderGoogle),
	apiOnly(contract.ProviderLinear),
	apiOnly(contract.ProviderGitHub),
	apiOnly(contract.ProviderTito),
	apiOnly(contract.ProviderNotion),
	apiOnly(contract.ProviderDiscord),
	apiOnly(contract.ProviderGrafana),
	apiOnly(contract.ProviderCRM),
	apiOnly(contract.ProviderHTTPAPI),
}

// Access returns every provider's access modes in display order.
func Access() []ProviderAccess {
	out := make([]ProviderAccess, len(providerAccess))
	for i, a := range providerAccess {
		out[i] = a.clone()
	}
	return out
}

// AccessFor returns a provider's access modes, or false when it has none.
func AccessFor(provider contract.Provider) (ProviderAccess, bool) {
	for _, a := range providerAccess {
		if a.Provider == provider {
			return a.clone(), true
		}
	}
	return ProviderAccess{}, false
}

func (a ProviderAccess) clone() ProviderAccess {
	out := a
	out.Access = make([]AccessMode, len(a.Access))
	for i, m := range a.Access {
		c := m
		if m.API != nil {
			api := *m.API
			c.API = &api
		}
		if m.CLI != nil {
			cli := *m.CLI
			c.CLI = &cli
		}
		if m.MCP != nil {
			mcp := *m.MCP
			mcp.Command = cloneStrings(m.MCP.Command)
			c.MCP = &mcp
		}
		if m.Webhook != nil {
			hook := *m.Webhook
			c.Webhook = &hook
		}
		c.Credentials = append([]CredentialEnv(nil), m.Credentials...)
		c.Setup, c.Verify = cloneSteps(m.Setup), cloneSteps(m.Verify)
		c.Reads, c.Writes = cloneStrings(m.Reads), cloneStrings(m.Writes)
		out.Access[i] = c
	}
	return out
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

func cloneSteps(steps []Step) []Step {
	if steps == nil {
		return nil
	}
	out := make([]Step, len(steps))
	for i, s := range steps {
		out[i] = Step{Name: s.Name, Description: s.Description, Run: cloneStrings(s.Run)}
	}
	return out
}

var (
	binaryName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)
	minVersion = regexp.MustCompile(`^[0-9]+(\.[0-9]+){0,2}$`)
	envName    = regexp.MustCompile(`^[A-Z_][A-Z0-9_]{0,127}$`)
	fieldName  = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
)

// Validate checks the access document's shape: the schema version, a known
// provider, one well-formed mode per kind, reads and writes that are the
// provider's read and write capabilities, credential wiring only on cli and
// mcp modes, named ordered steps, and no secret material anywhere.
func (a ProviderAccess) Validate() error {
	if a.Schema != AccessSchemaVersion {
		return fmt.Errorf("access schema must be %q", AccessSchemaVersion)
	}
	if !contract.ValidProvider(a.Provider) {
		return fmt.Errorf("unsupported provider %q", a.Provider)
	}
	kinds := map[Kind]bool{}
	for _, m := range a.Access {
		if kinds[m.Kind] {
			return fmt.Errorf("access mode %q is declared twice", m.Kind)
		}
		kinds[m.Kind] = true
		if err := m.validate(a.Provider); err != nil {
			return fmt.Errorf("access mode %q: %w", m.Kind, err)
		}
	}
	return nil
}

func (m AccessMode) validate(provider contract.Provider) error {
	blocks := 0
	for _, set := range []bool{m.API != nil, m.CLI != nil, m.MCP != nil, m.Webhook != nil} {
		if set {
			blocks++
		}
	}
	if blocks != 1 {
		return errors.New("exactly one mode block must be set")
	}
	switch m.Kind {
	case KindAPI:
		if m.API == nil || m.API.Adapter != provider {
			return errors.New("api adapter must be the provider")
		}
	case KindCLI:
		if m.CLI == nil || !binaryName.MatchString(m.CLI.Binary) || !minVersion.MatchString(m.CLI.MinVersion) || strings.TrimSpace(m.CLI.Install) == "" {
			return errors.New("cli needs a bare binary name, a numeric min_version, and an install hint")
		}
	case KindMCP:
		if err := m.MCP.validate(); err != nil {
			return err
		}
	case KindWebhook:
		if m.Webhook == nil || (m.Webhook.Direction != DirectionInbound && m.Webhook.Direction != DirectionOutbound) || !fieldName.MatchString(m.Webhook.SecretField) {
			return errors.New("webhook needs a direction and a secret_field")
		}
	default:
		return errors.New("unknown access kind")
	}
	if len(m.Credentials) > 0 && m.Kind != KindCLI && m.Kind != KindMCP {
		return errors.New("only cli and mcp modes wire credentials")
	}
	envs := map[string]bool{}
	for _, c := range m.Credentials {
		if !envName.MatchString(c.Env) || !fieldName.MatchString(c.Field) || envs[c.Env] {
			return errors.New("credentials need unique upper-case env names and credential field names")
		}
		envs[c.Env] = true
	}
	for _, steps := range [][]Step{m.Setup, m.Verify} {
		names := map[string]bool{}
		for _, s := range steps {
			if !fieldName.MatchString(s.Name) || strings.TrimSpace(s.Description) == "" || names[s.Name] {
				return errors.New("steps need unique names and a description")
			}
			names[s.Name] = true
			if containsSecretMaterial(append([]string{s.Description}, s.Run...)...) {
				return errors.New("steps must not contain secret material")
			}
		}
	}
	if m.CLI != nil && containsSecretMaterial(m.CLI.Install) {
		return errors.New("install hint must not contain secret material")
	}
	if m.MCP != nil && containsSecretMaterial(m.MCP.Command...) {
		return errors.New("mcp command must not contain secret material")
	}
	if err := validateCapabilities(provider, m.Reads, true); err != nil {
		return err
	}
	return validateCapabilities(provider, m.Writes, false)
}

func (m *MCPMode) validate() error {
	if m == nil {
		return errors.New("mcp block is required")
	}
	switch m.Transport {
	case TransportStdio:
		if len(m.Command) == 0 || m.Command[0] == "" || m.URL != "" {
			return errors.New("stdio mcp needs a command and no url")
		}
	case TransportStreamableHTTP:
		u, err := url.Parse(m.URL)
		if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || len(m.Command) > 0 {
			return errors.New("streamable_http mcp needs an https url and no command")
		}
	default:
		return fmt.Errorf("unknown mcp transport %q", m.Transport)
	}
	return nil
}

func validateCapabilities(provider contract.Provider, names []string, read bool) error {
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			return fmt.Errorf("capability %q is listed twice", name)
		}
		seen[name] = true
		c, ok := lookup(provider, name)
		if !ok {
			return fmt.Errorf("capability %q is not defined for provider %q", name, provider)
		}
		if (c.Action == contract.ActionRead) != read {
			if read {
				return fmt.Errorf("capability %q is not a read", name)
			}
			return fmt.Errorf("capability %q is not a write", name)
		}
	}
	return nil
}

func lookup(provider contract.Provider, name string) (contract.Capability, bool) {
	for _, c := range contract.CapabilitiesFor(provider) {
		if c.Name == name {
			return c, true
		}
	}
	return contract.Capability{}, false
}

// containsSecretMaterial mirrors the contract's secret markers: access
// metadata carries references and field names, never credential values.
func containsSecretMaterial(values ...string) bool {
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, marker := range []string{"bearer ", "basic ", "api_key=", "api-key=", "secret=", "token=", "sk_", "ghp_", "gho_", "glsa_", "ntn_", "xoxb-"} {
			if strings.Contains(lower, marker) {
				return true
			}
		}
	}
	return false
}

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

package bundle

import (
	"encoding/json"
	"net"
	"net/url"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// Connector is a policy_set.connectors entry: a workspace connector the bundle
// knows, with its provider and declared capabilities. In a v2 bundle the entry
// is also the connector's non-secret descriptor (ADR-034, HAI-378 amendment):
// where its credential lives and its non-secret config, never a credential.
// v1 entries carry only id, provider and capabilities.
type Connector struct {
	ID           string                `json:"id"`
	Name         string                `json:"name,omitempty"`
	Provider     string                `json:"provider"`
	Capabilities []ConnectorCapability `json:"capabilities"`
	// CredentialSource is oauth or opaque_ref; always set on a v2 descriptor.
	CredentialSource string `json:"credential_source,omitempty"`
	AccountModel     string `json:"account_model,omitempty"`
	// CredentialRef is the runtime-local kei://connectors/<id>/<field>
	// reference of an opaque_ref connector, resolved inside the daemon.
	CredentialRef string `json:"credential_ref,omitempty"`
	// CredentialBackend is the kind of the organization's credential store.
	CredentialBackend string                  `json:"credential_backend,omitempty"`
	Config            map[string]any          `json:"config,omitempty"`
	HTTPAPI           *contract.HTTPAPI       `json:"http_api,omitempty"`
	Grafana           *contract.GrafanaConfig `json:"grafana,omitempty"`
	// ADR-036 §4.3 descriptor fields (catalog HAI-570). Optional in the closed
	// schema until every runtime accepts them; mcp requires slug and mcp.
	Kind string         `json:"kind,omitempty"`
	Slug string         `json:"slug,omitempty"`
	MCP  *MCPDescriptor `json:"mcp,omitempty"`
}

// MCPDescriptor is the mcp block of a v2 connector descriptor (ADR-036 §4.3).
type MCPDescriptor struct {
	URL                  string               `json:"url"`
	Transport            string               `json:"transport"`
	AuthMethod           string               `json:"auth_method"`
	ApprovedManifestHash string               `json:"approved_manifest_hash,omitempty"`
	ApprovedTools        []MCPApprovedTool    `json:"approved_tools,omitempty"`
	CredentialFields     []MCPCredentialField `json:"credential_fields,omitempty"`
}

// MCPCredentialField is one credential_fields entry of a secret_ref mcp
// descriptor (contract "mcp.credential_fields exact schema"): the secret at
// kei://connectors/<connector id>/<field> is sent in header, after prefix.
// It never carries a value.
type MCPCredentialField struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	Prefix string `json:"prefix,omitempty"`
}

// MCPApprovedTool is one approved tool of an mcp descriptor (ADR-036 §4.3).
// approved_tools is sorted by name.
type MCPApprovedTool struct {
	Name     string `json:"name"`
	ToolHash string `json:"tool_hash"`
}

// ConnectorCapability is one declared capability of a bundle connector.
type ConnectorCapability struct {
	Name     string `json:"name"`
	Mutating bool   `json:"mutating"`
}

// Connector returns the bundle's connector entry with id.
func (v *Verified) Connector(id string) (Connector, bool) {
	set, err := v.decodePolicySet()
	if err != nil {
		return Connector{}, false
	}
	for _, c := range set.Connectors {
		if c.ID == id {
			return c, true
		}
	}
	return Connector{}, false
}

// CarriesConnectorDescriptors reports whether the bundle is v2 and its policy
// set has a connectors section. Without one (an older catalog) the runtime
// keeps its pre-descriptor behavior; with one, every routed connector must
// have an entry.
func (v *Verified) CarriesConnectorDescriptors() bool {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" {
		return false
	}
	var set struct {
		Connectors *[]json.RawMessage `json:"connectors"`
	}
	return json.Unmarshal(v.Payload.PolicySet, &set) == nil && set.Connectors != nil
}

// MCPConnector returns the bundle's mcp connector descriptor whose slug is
// slug (ADR-036 §4.3, §4.5 step 5). The closed v2 schema guarantees an mcp
// descriptor has its mcp block.
func (v *Verified) MCPConnector(slug string) (Connector, bool) {
	set, err := v.decodePolicySet()
	if err != nil || slug == "" {
		return Connector{}, false
	}
	for _, c := range set.Connectors {
		if c.Kind == "mcp" && c.Slug == slug && c.MCP != nil {
			return c, true
		}
	}
	return Connector{}, false
}

// MCPHosts returns host:port of every mcp descriptor URL in the bundle, the
// only upstream destinations the MCP route may dial (ADR-036 §4.10).
func (v *Verified) MCPHosts() []string {
	set, err := v.decodePolicySet()
	if err != nil {
		return nil
	}
	var hosts []string
	for _, c := range set.Connectors {
		if c.Kind != "mcp" || c.MCP == nil {
			continue
		}
		u, err := url.Parse(c.MCP.URL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" {
			continue
		}
		port := u.Port()
		if port == "" {
			port = "443"
		}
		hosts = append(hosts, net.JoinHostPort(strings.ToLower(u.Hostname()), port))
	}
	return hosts
}

// AgentBoundToConnector reports whether agent_connector_bindings binds agentID
// to connectorID, whatever the capabilities (ADR-036 §4.5 step 6).
func (v *Verified) AgentBoundToConnector(agentID, connectorID string) bool {
	if v == nil || agentID == "" || connectorID == "" {
		return false
	}
	return v.agentBound(agentID, connectorID, nil)
}

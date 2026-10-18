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

// Package skill generates a per-connector SKILL.md from a v1 tool-definition
// file. The output references the harness-level credential pass-through
// contract but never carries a token, secret, or client/DVL name.
package skill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ToolDef is a single entry in a .tool-definition.v1.json array.
type ToolDef struct {
	Schema      string          `json:"schema"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permission  string          `json:"permission"`
	Category    string          `json:"category"`
	Service     string          `json:"service,omitempty"`
	Version     string          `json:"version,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	Binding     json.RawMessage `json:"binding,omitempty"`
}

// Param is a structured parameter from the tool-definition array form.
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Enum        []any  `json:"enum,omitempty"`
	Default     any    `json:"default,omitempty"`
}

// connectorInfo holds metadata about a connector for human-readable output.
type connectorInfo struct {
	ID        string // the connector id used in filenames, e.g. "linear", "google_drive"
	HumanName string // "Linear", "Google Drive"
	Service   string // the service field common to tools
}

// knownConnectors maps connector IDs to display info.
// TODO: this could be driven from a registry; keep it inline for now.
var knownConnectors = map[string]connectorInfo{
	"linear":       {ID: "linear", HumanName: "Linear", Service: "linear"},
	"github":       {ID: "github", HumanName: "GitHub", Service: "github"},
	"google_drive": {ID: "google_drive", HumanName: "Google Drive", Service: "drive"},
	"http_api":     {ID: "http_api", HumanName: "HTTP API", Service: "http_api"},
}

// Options controls the generator output.
type Options struct {
	// ConnectorID is the programmatic connector name, e.g. "linear", "google_drive".
	// When empty it is derived from the service field of the first tool.
	ConnectorID string
}

// Generate reads a tool-definition array from r and returns the SKILL.md bytes.
func Generate(r io.Reader, opts Options) ([]byte, error) {
	var tools []ToolDef
	dec := json.NewDecoder(r)
	if err := dec.Decode(&tools); err != nil {
		return nil, fmt.Errorf("decoding tool definitions: %w", err)
	}
	if len(tools) == 0 {
		return nil, fmt.Errorf("no tool definitions found")
	}

	info := resolveConnector(tools, opts)

	var buf bytes.Buffer

	// YAML frontmatter.
	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "name: %s-connector\n", strings.ReplaceAll(info.ID, "_", "-"))
	fmt.Fprintf(&buf, "description: Governed %s connector — tool reference for harness use. Use when an agent needs to interact with %s through Kei-governed tools. Prefer this skill over generic %s knowledge.\n", info.HumanName, info.HumanName, info.HumanName)
	buf.WriteString("---\n\n")
	fmt.Fprintf(&buf, "# %s connector — governed tool reference\n\n", info.HumanName)

	fmt.Fprintf(&buf, "Use this skill when an agent needs to interact with %s through the governed connector. ", info.HumanName)
	buf.WriteString("Every governed tool call goes through `kei-proxy authorize` first. ")
	buf.WriteString("After an allow decision the harness obtains the credential via `POST /v1/credentials:release` ")
	buf.WriteString("over the Unix socket. The agent never holds, asks for, or prints a raw credential.\n\n")

	// Tool reference overview table.
	buf.WriteString("## Tool reference\n\n")
	buf.WriteString("| Tool | Description | Permission |\n")
	buf.WriteString("| --- | --- | --- |\n")
	for _, t := range tools {
		// Escape pipe characters in description.
		desc := strings.ReplaceAll(t.Description, "|", "\\|")
		buf.WriteString(fmt.Sprintf("| `%s` | %s | `%s` |\n", t.Name, desc, t.Permission))
	}
	buf.WriteString("\n")

	// Per-tool details.
	buf.WriteString("### Tool details\n\n")
	for _, t := range tools {
		fmt.Fprintf(&buf, "#### `%s`\n\n", t.Name)
		fmt.Fprintf(&buf, "- **Description**: %s\n", t.Description)
		fmt.Fprintf(&buf, "- **Permission**: `%s`\n", t.Permission)

		// Parse and render parameters.
		params, err := parseParams(t.Parameters)
		hasParams := err == nil && len(params) > 0
		if hasParams {
			buf.WriteString("- **Parameters**:\n\n")
			buf.WriteString("  | Name | Type | Required | Description |\n")
			buf.WriteString("  | --- | --- | --- | --- |\n")
			for _, p := range params {
				typ := p.Type
				if typ == "" {
					typ = "string"
				}
				req := "no"
				if p.Required {
					req = "yes"
				}
				pdesc := strings.ReplaceAll(p.Description, "|", "\\|")
				buf.WriteString(fmt.Sprintf("  | `%s` | %s | %s | %s |\n", p.Name, typ, req, pdesc))
			}
			buf.WriteString("\n")
		} else {
			buf.WriteString("- **Parameters**: none\n\n")
		}

		if len(t.Tags) > 0 {
			buf.WriteString(fmt.Sprintf("- **Tags**: `%s`\n\n", strings.Join(t.Tags, "`, `")))
		}
	}

	// Credential pass-through section.
	buf.WriteString("## Credential pass-through\n\n")
	buf.WriteString("Kei supplies the credential at run time. After `kei-proxy authorize` returns an allow ")
	buf.WriteString("decision, the harness calls `POST /v1/credentials:release` over the Unix socket\n")
	buf.WriteString("(`KEI_RUNTIME_SOCKET_PATH`, default `/run/kei-proxy/runtime.sock`) to obtain the ")
	buf.WriteString("short-lived access token.\n\n")
	buf.WriteString("The request identifies the tool by name and carries the session in `X-Kei-Session`. ")
	buf.WriteString("On allow the response carries `credential.access_token` and `credential.expires_at`. ")
	buf.WriteString("See `docs/credential-passthrough.md` in the kei-connector-contracts repo for the ")
	buf.WriteString("full contract.\n\n")
	buf.WriteString("The agent never holds, asks for, or prints a raw credential. The credential value ")
	buf.WriteString("crosses the socket only at call time and is never written to a config file, ")
	buf.WriteString("environment variable, disk cache, or log line.\n\n")

	// Connect section.
	buf.WriteString("## Connect\n\n")
	buf.WriteString("An owner connects this provider with:\n\n")
	buf.WriteString("```\n")
	fmt.Fprintf(&buf, "kei connectors create --provider %s --workspace W\n", info.ID)
	buf.WriteString("```\n\n")
	buf.WriteString("Check status with `kei connectors get <id>` and re-authorize with ")
	buf.WriteString("`kei connectors reconnect <id>`.\n")

	return buf.Bytes(), nil
}

// resolveConnector picks the connector info from the tools and options.
func resolveConnector(tools []ToolDef, opts Options) connectorInfo {
	if opts.ConnectorID != "" {
		if info, ok := knownConnectors[opts.ConnectorID]; ok {
			return info
		}
		return connectorInfo{ID: opts.ConnectorID, HumanName: opts.ConnectorID}
	}
	// Derive from the first tool's service or category.
	svc := tools[0].Service
	if svc == "" {
		svc = tools[0].Category
	}
	if info, ok := knownConnectors[svc]; ok {
		return info
	}
	return connectorInfo{ID: svc, HumanName: svc}
}

// parseParams attempts to parse parameters as an array of Param.
func parseParams(raw json.RawMessage) ([]Param, error) {
	if raw == nil || len(raw) == 0 {
		return nil, nil
	}
	// Try array form first.
	var arr []Param
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}
	// Object form (open-ended JSON Schema) — not detailed here.
	return nil, fmt.Errorf("unsupported parameter format")
}

// FormatConnectorName returns a human-readable name for a connector id.
// Exported for use by the cmd or tests.
func FormatConnectorName(id string) string {
	if info, ok := knownConnectors[id]; ok {
		return info.HumanName
	}
	return id
}

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

// Package manifest builds kei.connector-manifest/v1 (proposal 0002, section
// 4.6): the connectors, their connection kinds and setup fields, every
// operation as a capability, and the one-time legacy capability mapping
// (section 8.1). The policy catalog takes it in at build time.
package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
)

// Schema is the manifest's schema id.
const Schema = "kei.connector-manifest/v1"

// Manifest is the whole document. Every list is sorted (connectors and
// capabilities by id, legacy rows by connector and capability); connections,
// scopes and setup fields keep their declared order.
type Manifest struct {
	Schema             string      `json:"schema"`
	ContractsVersion   string      `json:"contracts_version"`
	Connectors         []Connector `json:"connectors"`
	LegacyCapabilities []LegacyRow `json:"legacy_capabilities"`
}

// Connector is one connector.
type Connector struct {
	ID           string       `json:"id"`
	BaseURL      string       `json:"base_url,omitempty"`
	Connections  []Connection `json:"connections"`
	Capabilities []Capability `json:"capabilities"`
}

// Connection is one connection kind of a connector.
type Connection struct {
	Name         string       `json:"name"`
	Kind         string       `json:"kind"`
	AuthMethod   string       `json:"auth_method"`
	AccountModel string       `json:"account_model"`
	BaseURL      string       `json:"base_url,omitempty"`
	BaseURLField string       `json:"base_url_field,omitempty"`
	HeaderName   string       `json:"header_name,omitempty"`
	Engine       string       `json:"engine,omitempty"`
	Scopes       []string     `json:"scopes"`
	SetupFields  []SetupField `json:"setup_fields"`
}

// SetupField is a value an admin fills in.
type SetupField struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Pattern string `json:"pattern,omitempty"`
	Secret  bool   `json:"secret"`
}

// Capability is one operation: its id is the capability policies grant.
type Capability struct {
	ID               string    `json:"id"`
	Class            string    `json:"class"`
	Destructive      bool      `json:"destructive"`
	Deprecated       bool      `json:"deprecated"`
	ReplacedBy       string    `json:"replaced_by"`
	Scopes           []string  `json:"scopes"`
	Connections      []string  `json:"connections"`
	ResourceTemplate string    `json:"resource_template"`
	Semantics        Semantics `json:"semantics"`
	Request          string    `json:"request"`
	Response         string    `json:"response"`
}

// Semantics is an operation's versioned semantic mapping (section 8.3).
type Semantics struct {
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	Version      uint32 `json:"version"`
}

// LegacyRow maps one coarse catalog capability to the operation ids it
// covers (section 8.1).
type LegacyRow struct {
	Connector  string   `json:"connector" yaml:"connector"`
	Capability string   `json:"capability" yaml:"capability"`
	Operations []string `json:"operations" yaml:"operations"`
}

// ParseLegacy reads legacy_capabilities.yaml: a list of rows. Unknown keys
// are errors.
func ParseLegacy(b []byte) ([]LegacyRow, error) {
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	var rows []LegacyRow
	if err := dec.Decode(&rows); err != nil {
		return nil, fmt.Errorf("legacy capabilities: %w", err)
	}
	return rows, nil
}

// CoarseClass is the class a coarse capability grants: "read" for a
// capability whose verb is read (issue.read), "write" for every other verb
// (issue.create, issue.comment, pull_request.update).
func CoarseClass(capability string) string {
	if capability[strings.LastIndex(capability, ".")+1:] == "read" {
		return "read"
	}
	return "write"
}

// LegacyProvider is the provider name a connector has in the frozen
// connector-setup and connector-access exports.
func LegacyProvider(connectorID string) string {
	if connectorID == "database" {
		return "postgres"
	}
	return connectorID
}

// ClassName is "read" or "write".
func ClassName(c contracts.OperationClass) string {
	return enumName(c.String(), "OPERATION_CLASS_")
}

func enumName(s, prefix string) string { return strings.ToLower(strings.TrimPrefix(s, prefix)) }

// Build makes the manifest from a validated set. It fails when a legacy row
// names an unknown connector or operation, an operation of another
// connector, an operation of another class than the coarse capability's, or
// a destructive operation: the mapping never widens a grant.
func Build(set *model.Set, version string, legacy []LegacyRow) (*Manifest, error) {
	m := &Manifest{Schema: Schema, ContractsVersion: version, Connectors: []Connector{}, LegacyCapabilities: []LegacyRow{}}
	for _, c := range set.Connectors {
		mc := Connector{ID: c.ID, BaseURL: c.Opt.GetBaseUrl(), Connections: []Connection{}, Capabilities: []Capability{}}
		for _, conn := range c.Opt.GetConnections() {
			mconn := Connection{
				Name:         conn.GetName(),
				Kind:         enumName(conn.GetKind().String(), "CONNECTION_KIND_"),
				AuthMethod:   conn.GetAuthMethod(),
				AccountModel: enumName(conn.GetAccountModel().String(), "ACCOUNT_MODEL_"),
				BaseURL:      conn.GetBaseUrl(),
				BaseURLField: conn.GetBaseUrlField(),
				HeaderName:   conn.GetHeaderName(),
				Engine:       conn.GetEngine(),
				Scopes:       nonNil(conn.GetScopes()),
				SetupFields:  []SetupField{},
			}
			for _, f := range conn.GetSetupFields() {
				mconn.SetupFields = append(mconn.SetupFields, SetupField{Name: f.GetName(), Type: f.GetType(), Pattern: f.GetPattern(), Secret: f.GetSecret()})
			}
			mc.Connections = append(mc.Connections, mconn)
		}
		for _, op := range c.Ops {
			opt := op.Opt
			mc.Capabilities = append(mc.Capabilities, Capability{
				ID:               op.ID,
				Class:            ClassName(opt.GetClass()),
				Destructive:      opt.GetDestructive(),
				Deprecated:       opt.GetDeprecated(),
				ReplacedBy:       opt.GetReplacedBy(),
				Scopes:           nonNil(opt.GetScopes()),
				Connections:      nonNil(opt.GetConnections()),
				ResourceTemplate: opt.GetResourceTemplate(),
				Semantics: Semantics{
					Action:       enumName(opt.GetSemantics().GetAction().String(), "SEMANTIC_ACTION_"),
					ResourceType: opt.GetSemantics().GetResourceType(),
					Version:      opt.GetSemantics().GetVersion(),
				},
				Request:  string(op.Request.Desc.FullName()),
				Response: string(op.Response.Desc.FullName()),
			})
		}
		sort.Slice(mc.Capabilities, func(i, j int) bool { return mc.Capabilities[i].ID < mc.Capabilities[j].ID })
		m.Connectors = append(m.Connectors, mc)
	}
	sort.Slice(m.Connectors, func(i, j int) bool { return m.Connectors[i].ID < m.Connectors[j].ID })

	rows, err := checkLegacy(set, legacy)
	if err != nil {
		return nil, err
	}
	m.LegacyCapabilities = rows
	return m, nil
}

func checkLegacy(set *model.Set, legacy []LegacyRow) ([]LegacyRow, error) {
	connectors := map[string]bool{}
	for _, c := range set.Connectors {
		connectors[c.ID] = true
	}
	ops := map[string]*model.Operation{}
	for _, op := range set.Ops {
		ops[op.ID] = op
	}
	seen := map[string]bool{}
	out := []LegacyRow{}
	for _, r := range legacy {
		key := r.Connector + " " + r.Capability
		switch {
		case !connectors[r.Connector]:
			return nil, fmt.Errorf("legacy %s: unknown connector %q", key, r.Connector)
		case r.Capability == "":
			return nil, fmt.Errorf("legacy row for %s: empty capability", r.Connector)
		case len(r.Operations) == 0:
			return nil, fmt.Errorf("legacy %s: maps to no operation; leave a capability with no operation out of the table", key)
		case seen[key]:
			return nil, fmt.Errorf("legacy %s: listed twice", key)
		}
		seen[key] = true
		want := CoarseClass(r.Capability)
		for _, id := range r.Operations {
			op := ops[id]
			switch {
			case op == nil:
				return nil, fmt.Errorf("legacy %s: unknown operation %s", key, id)
			case op.Opt.GetConnector() != r.Connector:
				return nil, fmt.Errorf("legacy %s: %s belongs to connector %s", key, id, op.Opt.GetConnector())
			case ClassName(op.Opt.GetClass()) != want:
				return nil, fmt.Errorf("legacy %s: %s is %s, the capability grants %s (the mapping never widens)", key, id, ClassName(op.Opt.GetClass()), want)
			case op.Opt.GetDestructive():
				return nil, fmt.Errorf("legacy %s: %s is destructive (the mapping never widens)", key, id)
			}
		}
		sorted := slices.Clone(r.Operations)
		sort.Strings(sorted)
		if len(slices.Compact(slices.Clone(sorted))) != len(sorted) {
			return nil, fmt.Errorf("legacy %s: an operation is listed twice", key)
		}
		out = append(out, LegacyRow{Connector: r.Connector, Capability: r.Capability, Operations: sorted})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Connector != out[j].Connector {
			return out[i].Connector < out[j].Connector
		}
		return out[i].Capability < out[j].Capability
	})
	return out, nil
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// Marshal renders the manifest: two-space indent, no HTML escaping, a
// trailing newline.
func Marshal(m *Manifest) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

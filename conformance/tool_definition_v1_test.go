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

package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// toolDefRef holds key fields we check in every tool-definition entry.
type toolDefRef struct {
	Schema      string          `json:"schema"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permission  string          `json:"permission"`
	Category    string          `json:"category"`
	Tags        []string        `json:"tags,omitempty"`
	Service     string          `json:"service,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	Binding     json.RawMessage `json:"binding,omitempty"`
}

func TestToolDefinitionV1SchemaExists(t *testing.T) {
	schemaPath := filepath.Join("..", "schemas", "tool-definition.v1.schema.json")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Fatal("tool-definition.v1.schema.json must exist")
	}
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc document
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("schema parse: %v", err)
	}
	if doc["$id"] == nil || !strings.Contains(doc["$id"].(string), "tool-definition.v1.schema") {
		t.Fatalf("schema $id missing or wrong: %v", doc["$id"])
	}
}

func discoverToolDefDataFiles(t *testing.T) []string {
	t.Helper()
	pattern := filepath.Join("..", "schemas", "*.tool-definition.v1.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("no .tool-definition.v1.json data files found in schemas/")
	}
	return matches
}

func loadToolDefFile(t *testing.T, path string) []toolDefRef {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var tools []toolDefRef
	if err := json.Unmarshal(b, &tools); err != nil {
		t.Fatalf("%s: %v", filepath.Base(path), err)
	}
	if len(tools) == 0 {
		t.Fatalf("%s: must contain at least one tool definition", filepath.Base(path))
	}
	return tools
}

func TestToolDefinitionV1SchemaConformance(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)

			// Verify each entry has the tool-definition/v1 schema marker.
			for i, tool := range tools {
				if tool.Schema != "kei.tool-definition/v1" {
					t.Errorf("%s[%d]: schema must be kei.tool-definition/v1, got %q", name, i, tool.Schema)
				}
			}
		})
	}
}

func TestToolDefinitionV1RequiredFields(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)
			for i, tool := range tools {
				if tool.Name == "" {
					t.Errorf("%s[%d]: name is required", name, i)
				}
				if tool.Description == "" {
					t.Errorf("%s[%d]: description is required", name, i)
				}
				if tool.Permission == "" {
					t.Errorf("%s[%d]: permission is required", name, i)
				}
				if tool.Category == "" {
					t.Errorf("%s[%d]: category is required", name, i)
				}
			}
		})
	}
}

func TestToolDefinitionV1ConnectorReadHasService(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)
			for i, tool := range tools {
				// Connector read tools (those with a binding) must declare a service.
				if tool.Binding != nil && tool.Service == "" {
					t.Errorf("%s[%d] %s: governed connector read tools must declare a service", name, i, tool.Name)
				}
			}
		})
	}
}

var validPermissions = map[string]bool{
	"linear.read":   true,
	"linear.write":  true,
	"file.write":    true,
	"github.read":   true,
	"github.write":  true,
	"drive.read":    true,
	"drive.write":   true,
	"http_api.read": true,
}

var validCategories = map[string]bool{
	"linear":   true,
	"github":   true,
	"drive":    true,
	"http_api": true,
}

func TestToolDefinitionV1PermissionEnum(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)
			for i, tool := range tools {
				if !validPermissions[tool.Permission] {
					t.Errorf("%s[%d] %s: invalid permission %q", name, i, tool.Name, tool.Permission)
				}
			}
		})
	}
}

func TestToolDefinitionV1CategoryEnum(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)
			for i, tool := range tools {
				if !validCategories[tool.Category] {
					t.Errorf("%s[%d] %s: invalid category %q", name, i, tool.Name, tool.Category)
				}
			}
		})
	}
}

func TestToolDefinitionV1UniqueNames(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			tools := loadToolDefFile(t, path)
			seen := make(map[string]int)
			for i, tool := range tools {
				if tool.Name == "" {
					continue
				}
				if j, exists := seen[tool.Name]; exists {
					t.Errorf("%s[%d] %s: duplicate tool name (also at index %d)", name, i, tool.Name, j)
				}
				seen[tool.Name] = i
			}
		})
	}
}

func TestToolDefinitionV1ProviderFileCount(t *testing.T) {
	expected := map[string]int{
		"linear":       4,  // linear.list_issues, linear.get_issue, linear.list_projects, linear.create_issue
		"google_drive": 5,  // drive.list_files, drive.get_file, docs.get_document, drive_create_doc, drive_update_doc
		"github":       10, // github.get_repository, github.get_issue, github.get_pull_request, list_prs, list_issues, create_issue, get_workflow_status, create_pull_request, github.create_fix_pr, github.file_write
		"http_api":     2,  // http_api.list_records, http_api.get_record
	}
	for provider, want := range expected {
		dataPath := filepath.Join("..", "schemas", provider+".tool-definition.v1.json")
		_, err := os.Stat(dataPath)
		if os.IsNotExist(err) {
			t.Errorf("missing data file for provider %q: %s", provider, dataPath)
			continue
		}
		tools := loadToolDefFile(t, dataPath)
		if len(tools) != want {
			t.Errorf("%s: expected %d tool definitions, got %d", provider, want, len(tools))
		}
	}
}

// HAI-723 (owner decision, option A, 2026-10-20): per ADR-039 the harness
// builds the provider call itself, so tool_binding must carry non-secret
// transport metadata: method, url_template, auth placement (no secret
// values) and arg mapping (tool param -> {in, name}).

func loadToolDefSchema(t *testing.T) document {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "schemas", "tool-definition.v1.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc document
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("schema parse: %v", err)
	}
	return doc
}

func schemaTransport(t *testing.T) document {
	t.Helper()
	schema := loadToolDefSchema(t)
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("schema has no $defs")
	}
	binding, ok := defs["tool_binding"].(map[string]any)
	if !ok {
		t.Fatal("schema has no $defs.tool_binding")
	}
	props, ok := binding["properties"].(map[string]any)
	if !ok {
		t.Fatal("$defs.tool_binding has no properties")
	}
	transport, ok := props["transport"].(map[string]any)
	if !ok {
		t.Fatal("$defs.tool_binding.properties.transport missing: harness cannot build provider calls (HAI-723)")
	}
	return transport
}

func enumSet(v any) map[string]bool {
	set := map[string]bool{}
	if list, ok := v.([]any); ok {
		for _, e := range list {
			if s, ok := e.(string); ok {
				set[s] = true
			}
		}
	}
	return set
}

func requireEnum(t *testing.T, got map[string]bool, want []string, label string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: enum %v, want exactly %v", label, got, want)
		return
	}
	for _, w := range want {
		if !got[w] {
			t.Errorf("%s: enum %v missing %q", label, got, w)
		}
	}
}

// asDocument accepts both the named document type and a bare map[string]any
// (JSON unmarshal produces the latter).
func asDocument(v any) (document, bool) {
	if d, ok := v.(document); ok {
		return d, true
	}
	if m, ok := v.(map[string]any); ok {
		return document(m), true
	}
	return nil, false
}

// validateTransport enforces the constraints the schema declares for a
// transport block: closed property sets (additionalProperties: false, so no
// secret values can slip in) and enum membership.
func validateTransport(transport, tschema document) []string {
	var errs []string
	props, _ := asDocument(tschema["properties"])
	if tschema["additionalProperties"] != false {
		errs = append(errs, "transport schema must set additionalProperties: false")
	}
	for key := range transport {
		if _, ok := props[key]; !ok {
			errs = append(errs, fmt.Sprintf("transport.%s: not declared by schema", key))
		}
	}
	if m, ok := transport["method"]; ok {
		methodSchema, _ := asDocument(props["method"])
		enum, _ := methodSchema["enum"].([]any)
		if !enumSet(enum)[fmt.Sprintf("%v", m)] {
			errs = append(errs, fmt.Sprintf("transport.method %v not in enum %v", m, enum))
		}
	}
	if a, ok := asDocument(transport["auth"]); ok {
		authSchema, _ := asDocument(props["auth"])
		authProps, _ := asDocument(authSchema["properties"])
		if authSchema["additionalProperties"] != false {
			errs = append(errs, "transport.auth schema must set additionalProperties: false (no secret values)")
		}
		for key := range a {
			if _, ok := authProps[key]; !ok {
				errs = append(errs, fmt.Sprintf("transport.auth.%s: not declared by schema (secret values are forbidden)", key))
			}
		}
		if p, ok := a["placement"]; ok {
			placementSchema, _ := asDocument(authProps["placement"])
			enum, _ := placementSchema["enum"].([]any)
			if !enumSet(enum)[fmt.Sprintf("%v", p)] {
				errs = append(errs, fmt.Sprintf("transport.auth.placement %v not in enum %v", p, enum))
			}
		}
	}
	if args, ok := asDocument(transport["args"]); ok {
		argsSchema, _ := asDocument(props["args"])
		valSchema, _ := asDocument(argsSchema["additionalProperties"])
		valProps, _ := asDocument(valSchema["properties"])
		for name, raw := range args {
			v, ok := asDocument(raw)
			if !ok {
				errs = append(errs, fmt.Sprintf("transport.args.%s: must be an object", name))
				continue
			}
			for key := range v {
				if _, ok := valProps[key]; !ok {
					errs = append(errs, fmt.Sprintf("transport.args.%s.%s: not declared by schema", name, key))
				}
			}
			if in, ok := v["in"]; ok {
				inSchema, _ := asDocument(valProps["in"])
				enum, _ := inSchema["enum"].([]any)
				if !enumSet(enum)[fmt.Sprintf("%v", in)] {
					errs = append(errs, fmt.Sprintf("transport.args.%s.in %v not in enum %v", name, in, enum))
				}
			}
		}
	}
	return errs
}

func TestToolDefinitionV1TransportBlock(t *testing.T) {
	tschema := schemaTransport(t)
	props := object(tschema["properties"])

	requireEnum(t, enumSet(object(props["method"])["enum"]), []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, "transport.method")
	requireEnum(t, enumSet(object(object(object(props["auth"])["properties"])["placement"])["enum"]), []string{"header", "query", "none"}, "transport.auth.placement")
	requireEnum(t, enumSet(object(object(object(object(props["args"])["additionalProperties"])["properties"])["in"])["enum"]), []string{"path", "query", "body", "header"}, "transport.args.*.in")

	valid := document{
		"method":       "GET",
		"url_template": "/v1/issues?state={state}",
		"auth":         document{"placement": "header", "name": "Authorization", "scheme": "Bearer"},
		"args":         document{"state": document{"in": "query", "name": "state"}},
	}
	if errs := validateTransport(valid, tschema); len(errs) != 0 {
		t.Errorf("valid transport block must validate: %v", errs)
	}

	unknownMethod := document{
		"method":       "FETCH",
		"url_template": "/v1/issues",
	}
	if errs := validateTransport(unknownMethod, tschema); len(errs) == 0 {
		t.Error("transport with unknown method FETCH must fail")
	}

	secretField := document{
		"method":       "GET",
		"url_template": "/v1/issues",
		"auth":         document{"placement": "header", "name": "Authorization", "scheme": "Bearer", "token": "ghp_0123456789"},
	}
	if errs := validateTransport(secretField, tschema); len(errs) == 0 {
		t.Error("transport.auth with a secret-looking field (token) must fail")
	}
}

func TestConnectorDefinitionsHaveTransportForHTTPTools(t *testing.T) {
	for _, connector := range []string{"linear", "github", "google_drive", "http_api"} {
		path := filepath.Join("..", "schemas", connector+".tool-definition.v1.json")
		tools := loadToolDefFile(t, path)
		for _, tool := range tools {
			isMCP := false
			for _, tag := range tool.Tags {
				// harness-action tools (HAI-856) run a bounded harness action,
				// not a single provider HTTP call, so they carry no transport.
				if strings.EqualFold(tag, "mcp") || strings.EqualFold(tag, "harness-action") {
					isMCP = true
				}
			}
			if isMCP {
				continue
			}
			t.Run(tool.Name, func(t *testing.T) {
				if tool.Binding == nil {
					t.Fatalf("non-MCP tool %s must declare binding.transport", tool.Name)
				}
				var binding map[string]json.RawMessage
				if err := json.Unmarshal(tool.Binding, &binding); err != nil {
					t.Fatal(err)
				}
				var transport struct {
					Method      string                     `json:"method"`
					URLTemplate string                     `json:"url_template"`
					Args        map[string]json.RawMessage `json:"args"`
				}
				if raw := binding["transport"]; len(raw) == 0 {
					t.Fatalf("non-MCP tool %s must declare binding.transport", tool.Name)
				} else if err := json.Unmarshal(raw, &transport); err != nil {
					t.Fatal(err)
				}
				if transport.Method == "" || transport.URLTemplate == "" {
					t.Fatalf("transport for %s must include method and url_template", tool.Name)
				}
				for _, placeholder := range regexp.MustCompile(`\{([^{}]+)\}`).FindAllStringSubmatch(transport.URLTemplate, -1) {
					if _, ok := transport.Args[placeholder[1]]; !ok {
						t.Errorf("URL placeholder {%s} must have an entry in transport.args", placeholder[1])
					}
				}
			})
		}
	}
}

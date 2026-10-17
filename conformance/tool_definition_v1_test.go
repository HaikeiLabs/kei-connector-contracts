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
	"os"
	"path/filepath"
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
		"linear":       5, // linear.list_issues, linear.get_issue, linear.list_projects, linear.create_issue, linear.create_followup_task
		"google_drive": 5, // drive.list_files, drive.get_file, docs.get_document, drive_create_doc, drive_update_doc
		"github":       8, // github.get_repository, github.get_issue, github.get_pull_request, list_prs, list_issues, create_issue, get_workflow_status, create_pull_request
		"http_api":     2, // http_api.list_records, http_api.get_record
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

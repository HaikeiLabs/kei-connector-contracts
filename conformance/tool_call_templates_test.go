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
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/calltemplate"
)

// HAI-856: every connector tool-definition file must pass the tool call
// template profile (closed schema, placeholders only, no secrets), so a
// harness importing a pinned version can trust what it reads.
func TestConnectorToolDefinitionsAreValidCallTemplates(t *testing.T) {
	for _, path := range discoverToolDefDataFiles(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			if _, err := calltemplate.LoadFile(path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// templateParity pins the schema and permission of the first harness tools to
// the hand-written definitions they replace: github.create_fix_pr
// (kei-agents#45, src/agents/github/tools.py), the governed file_write
// (chat harness agent.py file_write_tool, config/tools-template.yaml) and the
// Linear write tools. A change here is a change to what policy sees.
type templateParity struct {
	file       string
	name       string
	permission string
	required   []string
	optional   []string
	defaults   map[string]any
	hasHTTP    bool
}

var firstTemplates = []templateParity{
	{
		file: "github", name: "github.create_fix_pr", permission: "github.write",
		required: []string{"body", "title", "unified_diff"},
	},
	{
		file: "github", name: "github.file_write", permission: "file.write",
		required: []string{"content", "path"},
		optional: []string{"branch", "commit_message", "owner", "repo"},
		defaults: map[string]any{"branch": "main", "commit_message": "", "owner": "soypete", "repo": ""},
		hasHTTP:  true,
	},
	{
		file: "linear", name: "linear.create_issue", permission: "linear.write",
		required: []string{"team_key", "title"},
		optional: []string{"description", "labels", "priority"},
		hasHTTP:  true,
	},
	{
		file: "linear", name: "linear.create_followup_task", permission: "linear.write",
		required: []string{"idempotency_key", "lead_id", "summary", "title"},
		optional: []string{"approval_id"},
		hasHTTP:  true,
	},
}

func loadTemplates(t *testing.T, connector string) []calltemplate.Template {
	t.Helper()
	tpls, err := calltemplate.LoadFile(filepath.Join("..", "schemas", connector+".tool-definition.v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	return tpls
}

func findTemplate(t *testing.T, connector, name string) calltemplate.Template {
	t.Helper()
	for _, tpl := range loadTemplates(t, connector) {
		if tpl.Name == name {
			return tpl
		}
	}
	t.Fatalf("template %s missing from %s", name, connector)
	return calltemplate.Template{}
}

func TestFirstToolCallTemplatesMatchHandWrittenDefinitions(t *testing.T) {
	for _, want := range firstTemplates {
		t.Run(want.name, func(t *testing.T) {
			got := findTemplate(t, want.file, want.name)
			if got.Permission != want.permission {
				t.Errorf("permission = %q, want %q", got.Permission, want.permission)
			}
			required, optional, defaults := parameterShape(t, got.Parameters)
			if !reflect.DeepEqual(required, want.required) {
				t.Errorf("required parameters = %v, want %v", required, want.required)
			}
			if !reflect.DeepEqual(optional, nonNil(want.optional)) {
				t.Errorf("optional parameters = %v, want %v", optional, want.optional)
			}
			for k, v := range want.defaults {
				if defaults[k] != v {
					t.Errorf("default for %s = %v, want %v", k, defaults[k], v)
				}
			}
			hasHTTP := got.Binding != nil && got.Binding.Transport != nil
			if hasHTTP != want.hasHTTP {
				t.Errorf("has transport = %v, want %v", hasHTTP, want.hasHTTP)
			}
		})
	}
}

func TestCreateFixPRSchemaIsClosedAndBounded(t *testing.T) {
	tpl := findTemplate(t, "github", "github.create_fix_pr")
	var schema struct {
		AdditionalProperties bool `json:"additionalProperties"`
		Properties           map[string]struct {
			Type      string `json:"type"`
			MaxLength int    `json:"maxLength"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(tpl.Parameters, &schema); err != nil {
		t.Fatal(err)
	}
	if schema.AdditionalProperties {
		t.Error("create_fix_pr must set additionalProperties: false")
	}
	for _, p := range []string{"title", "body", "unified_diff"} {
		if schema.Properties[p].Type != "string" {
			t.Errorf("%s must be a string", p)
		}
	}
	if got := schema.Properties["unified_diff"].MaxLength; got != 32768 {
		t.Errorf("unified_diff maxLength = %d, want 32768 (32 KiB)", got)
	}
}

func TestFileWriteContentIsBase64EncodedForTheContentsAPI(t *testing.T) {
	tr := findTemplate(t, "github", "github.file_write").Binding.Transport
	if tr.Method != "PUT" {
		t.Errorf("file_write method = %s, want PUT", tr.Method)
	}
	if got := tr.Args["content"]; got.In != "body" || got.Encoding != "base64" {
		t.Errorf("content arg = %+v, want body + base64", got)
	}
}

// parameterShape reads required/optional names and defaults from either
// parameter form (list of parameters, or JSON Schema object).
func parameterShape(t *testing.T, raw json.RawMessage) (required, optional []string, defaults map[string]any) {
	t.Helper()
	defaults = map[string]any{}
	var list []struct {
		Name     string `json:"name"`
		Required bool   `json:"required"`
		Default  any    `json:"default"`
	}
	if err := json.Unmarshal(raw, &list); err == nil {
		for _, p := range list {
			if p.Required {
				required = append(required, p.Name)
			} else {
				optional = append(optional, p.Name)
			}
			if p.Default != nil {
				defaults[p.Name] = p.Default
			}
		}
	} else {
		var obj struct {
			Required   []string                  `json:"required"`
			Properties map[string]map[string]any `json:"properties"`
		}
		if err := json.Unmarshal(raw, &obj); err != nil {
			t.Fatal(err)
		}
		isReq := map[string]bool{}
		for _, r := range obj.Required {
			isReq[r] = true
		}
		for name, prop := range obj.Properties {
			if isReq[name] {
				required = append(required, name)
			} else {
				optional = append(optional, name)
			}
			if d, ok := prop["default"]; ok {
				defaults[name] = d
			}
		}
	}
	sort.Strings(required)
	sort.Strings(optional)
	return required, nonNil(optional), defaults
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

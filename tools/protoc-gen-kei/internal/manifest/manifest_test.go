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

package manifest_test

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/plugintest"
)

const testVersion = "v0.0.0-test"

func repoRoot() string { return filepath.Join(plugintest.RepoProto(), "..") }

// connectorProtos lists every .proto under proto/kei/connectors.
func connectorProtos(t *testing.T) []string {
	t.Helper()
	var files []string
	root := plugintest.RepoProto()
	err := filepath.WalkDir(filepath.Join(root, "kei", "connectors"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".proto") {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	return files
}

func loadSet(t *testing.T) *model.Set {
	t.Helper()
	gen := plugintest.Plugin(t, plugintest.RepoProto(), connectorProtos(t), "", nil)
	set, err := model.Load(gen.Files)
	if err != nil {
		t.Fatalf("model.Load: %v", err)
	}
	return set
}

func loadLegacy(t *testing.T) []manifest.LegacyRow {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(plugintest.RepoProto(), "kei", "contracts", "v1", "legacy_capabilities.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	rows, err := manifest.ParseLegacy(b)
	if err != nil {
		t.Fatalf("ParseLegacy: %v", err)
	}
	return rows
}

func build(t *testing.T) []byte {
	t.Helper()
	m, err := manifest.Build(loadSet(t), testVersion, loadLegacy(t))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, err := manifest.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	return b
}

func TestManifest_ValidatesAgainstSchema(t *testing.T) {
	schemaPath := filepath.Join(repoRoot(), "schemas", "connector-manifest.v1.schema.json")
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaPath)
	if err != nil {
		t.Fatalf("compiling the schema: %v", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(build(t)))
	if err != nil {
		t.Fatal(err)
	}
	if err := sch.Validate(doc); err != nil {
		t.Fatalf("manifest does not validate: %v", err)
	}
	// The schema is not vacuous: a manifest with an unknown class fails.
	bad := strings.Replace(string(build(t)), `"class": "read"`, `"class": "admin"`, 1)
	badDoc, _ := jsonschema.UnmarshalJSON(strings.NewReader(bad))
	if err := sch.Validate(badDoc); err == nil {
		t.Fatalf("a capability with class admin validated")
	}
}

func TestManifest_Deterministic(t *testing.T) {
	a, b := build(t), build(t)
	if !bytes.Equal(a, b) {
		t.Fatalf("two builds differ")
	}
	if !bytes.HasSuffix(a, []byte("}\n")) {
		t.Fatalf("manifest must end with a newline")
	}
}

// TestManifest_Shape pins the section 4.6 shape for github.issues.get and
// the sorted order of every list.
func TestManifest_Shape(t *testing.T) {
	var m struct {
		Schema           string `json:"schema"`
		ContractsVersion string `json:"contracts_version"`
		Connectors       []struct {
			ID          string `json:"id"`
			Connections []struct {
				Name        string `json:"name"`
				Kind        string `json:"kind"`
				AuthMethod  string `json:"auth_method"`
				AccountMode string `json:"account_model"`
				SetupFields []struct {
					Name   string `json:"name"`
					Secret bool   `json:"secret"`
				} `json:"setup_fields"`
			} `json:"connections"`
			Capabilities []map[string]any `json:"capabilities"`
		} `json:"connectors"`
		Legacy []manifest.LegacyRow `json:"legacy_capabilities"`
	}
	if err := json.Unmarshal(build(t), &m); err != nil {
		t.Fatal(err)
	}
	if m.Schema != "kei.connector-manifest/v1" || m.ContractsVersion != testVersion {
		t.Fatalf("header = %q %q", m.Schema, m.ContractsVersion)
	}
	var ids []string
	for _, c := range m.Connectors {
		ids = append(ids, c.ID)
		var capIDs []string
		for _, cap := range c.Capabilities {
			capIDs = append(capIDs, cap["id"].(string))
		}
		if !sort.StringsAreSorted(capIDs) {
			t.Errorf("%s capabilities not sorted: %v", c.ID, capIDs)
		}
		if c.ID != "github" {
			continue
		}
		if len(c.Connections) != 2 || c.Connections[1].Name != "pat" || c.Connections[1].Kind != "http_bearer" ||
			c.Connections[1].AccountMode != "shared" || !c.Connections[1].SetupFields[0].Secret {
			t.Errorf("github connections = %+v", c.Connections)
		}
		for _, cap := range c.Capabilities {
			if cap["id"] != "github.issues.get" {
				continue
			}
			want := map[string]any{
				"id": "github.issues.get", "class": "read", "destructive": false,
				"deprecated": false, "replaced_by": "",
				"scopes": []any{"issues:read"}, "connections": []any{"user_oauth", "pat"},
				"resource_template": "repos/{owner}/{repository}",
				"semantics":         map[string]any{"action": "read", "resource_type": "issue", "version": float64(1)},
				"request":           "kei.connectors.github.v1.GetIssueRequest",
				"response":          "kei.connectors.github.v1.GetIssueResponse",
			}
			for k, v := range want {
				got, _ := json.Marshal(cap[k])
				exp, _ := json.Marshal(v)
				if string(got) != string(exp) {
					t.Errorf("github.issues.get %s = %s, want %s", k, got, exp)
				}
			}
		}
	}
	if !sort.StringsAreSorted(ids) || len(ids) < 4 {
		t.Fatalf("connectors = %v", ids)
	}
}

func TestManifest_LegacyRowsReferenceKnownOps(t *testing.T) {
	set := loadSet(t)
	if _, err := manifest.Build(set, testVersion, loadLegacy(t)); err != nil {
		t.Fatalf("the committed mapping: %v", err)
	}
	rows := []manifest.LegacyRow{{Connector: "github", Capability: "issue.read", Operations: []string{"github.issues.get", "github.issues.nope"}}}
	if _, err := manifest.Build(set, testVersion, rows); err == nil || !strings.Contains(err.Error(), "github.issues.nope") {
		t.Fatalf("Build with an unknown op = %v, want an error naming it", err)
	}
	rows = []manifest.LegacyRow{{Connector: "linear", Capability: "issue.read", Operations: []string{"github.issues.get"}}}
	if _, err := manifest.Build(set, testVersion, rows); err == nil {
		t.Fatalf("Build with another connector's op: want an error")
	}
}

// TestLegacyRows_Expected pins the section 8.1 rows whose operations exist,
// so a dropped, widened or narrowed row shows up as a test diff.
func TestLegacyRows_Expected(t *testing.T) {
	want := map[string]string{
		"github issue.read":     "github.issues.get github.issues.list",
		"linear comment.create": "linear.comments.create",
		"linear issue.create":   "linear.issues.create",
		"linear issue.read":     "linear.issues.get linear.issues.search",
		"linear team.read":      "linear.teams.get",
	}
	m, err := manifest.Build(loadSet(t), testVersion, loadLegacy(t))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, r := range m.LegacyCapabilities {
		got[r.Connector+" "+r.Capability] = strings.Join(r.Operations, " ")
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("legacy %s = %q, want %q", k, got[k], v)
		}
	}
	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf("unexpected legacy row %s", k)
		}
	}
}

// TestManifest_NoWidening: a coarse capability maps only to operations of
// its own class (".read" capabilities to read operations, the rest to
// write operations), and never to a destructive one.
func TestManifest_NoWidening(t *testing.T) {
	set := loadSet(t)
	byID := map[string]*model.Operation{}
	for _, op := range set.Ops {
		byID[op.ID] = op
	}
	for _, row := range loadLegacy(t) {
		want := manifest.CoarseClass(row.Capability)
		for _, id := range row.Operations {
			op := byID[id]
			if op == nil {
				t.Fatalf("%s %s maps to unknown %s", row.Connector, row.Capability, id)
			}
			if got := manifest.ClassName(op.Opt.GetClass()); got != want || op.Opt.GetDestructive() {
				t.Errorf("%s %s (%s) maps to %s (%s, destructive=%v)", row.Connector, row.Capability, want, id, got, op.Opt.GetDestructive())
			}
		}
	}
	widening := []manifest.LegacyRow{{Connector: "github", Capability: "issue.read", Operations: []string{"github.contents.put"}}}
	if _, err := manifest.Build(set, testVersion, widening); err == nil {
		t.Fatalf("mapping issue.read to a write: want an error")
	}
}

// TestLegacyExports_Consistent replaces the byte-match of section 7: the
// frozen exports stay frozen until phase 5, and only the facts the manifest
// derives are checked against them.
func TestLegacyExports_Consistent(t *testing.T) {
	type setupField struct {
		Name   string `json:"name"`
		Secret bool   `json:"secret"`
	}
	var setup []struct {
		Provider string       `json:"provider"`
		Fields   []setupField `json:"fields"`
	}
	var access []struct {
		Provider string `json:"provider"`
	}
	for name, v := range map[string]any{"connector-setup.v1.json": &setup, "connector-access.v1.json": &access} {
		b, err := os.ReadFile(filepath.Join(repoRoot(), "schemas", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b, v); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	setupBy := map[string][]setupField{}
	for _, p := range setup {
		setupBy[p.Provider] = p.Fields
	}
	accessBy := map[string]bool{}
	for _, p := range access {
		accessBy[p.Provider] = true
	}
	legacy := 0
	for _, c := range loadSet(t).Connectors {
		provider := manifest.LegacyProvider(c.ID)
		fields, inSetup := setupBy[provider]
		// A connector added after the exports were frozen is in neither,
		// and that is fine; one export without the other is not.
		if !inSetup && !accessBy[provider] {
			continue
		}
		legacy++
		if !inSetup {
			t.Errorf("connector %s: %q is in connector-access.v1.json but not connector-setup.v1.json", c.ID, provider)
		}
		if !accessBy[provider] {
			t.Errorf("connector %s: %q is in connector-setup.v1.json but not connector-access.v1.json", c.ID, provider)
		}
		for _, conn := range c.Opt.GetConnections() {
			for _, sf := range conn.GetSetupFields() {
				for _, lf := range fields {
					if lf.Name == sf.GetName() && lf.Secret != sf.GetSecret() {
						t.Errorf("%s field %s: secret=%v in the manifest, %v in the legacy export", c.ID, sf.GetName(), sf.GetSecret(), lf.Secret)
					}
				}
			}
		}
	}
	// The four connectors the exports already covered must all be found.
	if legacy < 4 {
		t.Errorf("only %d manifest connectors matched a legacy export entry, want at least github, linear, http_api and database", legacy)
	}
	if manifest.LegacyProvider("database") != "postgres" {
		t.Fatalf("database must alias postgres")
	}
}

// TestManifest_MatchesCommitted: manifest/connector-manifest.v1.json is the
// plugin's output for the committed protos (only the version may differ).
func TestManifest_MatchesCommitted(t *testing.T) {
	committed, err := os.ReadFile(filepath.Join(repoRoot(), "manifest", "connector-manifest.v1.json"))
	if err != nil {
		t.Fatalf("reading the committed manifest (run buf generate): %v", err)
	}
	var head struct {
		Version string `json:"contracts_version"`
	}
	if err := json.Unmarshal(committed, &head); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Build(loadSet(t), head.Version, loadLegacy(t))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := manifest.Marshal(m)
	if !bytes.Equal(got, committed) {
		t.Fatalf("manifest/connector-manifest.v1.json is stale; run buf generate")
	}
}

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

package render_test

import (
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/plugintest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/render"
)

// TestRustNestedMessageBody: a repeated-message body param (the shape of
// power_bi.datasets.execute_queries) renders with the canonical body rules:
// only populated fields, proto names unless json_name differs, and nested
// messages encoded field by field.
func TestRustNestedMessageBody(t *testing.T) {
	var opts render.Options
	gen := plugintest.Plugin(t, "testdata/nested", []string{"kei/testnested/v1/nested.proto"}, "paths=source_relative,lang=rust", opts.Set)
	if err := render.Generate(gen, opts); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	resp := gen.Response()
	if resp.Error != nil {
		t.Fatalf("plugin error: %s", resp.GetError())
	}
	var out string
	for _, f := range resp.File {
		out += f.GetContent()
	}
	for _, want := range []string{
		`if !req.queries.is_empty() {`,
		`"queries".to_string()`,
		`"query".to_string()`,
		`"max_rows".to_string()`,
		`"title".to_string()`,
		`"include_nulls".to_string()`,
		`"tags".to_string()`,
		`if let Some(v1) = &v0.label {`,
		`if let Some(v0) = &req.settings {`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated Rust lacks %s", want)
		}
	}
	if strings.Contains(out, `"maxRows"`) {
		t.Errorf("a REST body must use the proto name max_rows, not maxRows")
	}
}

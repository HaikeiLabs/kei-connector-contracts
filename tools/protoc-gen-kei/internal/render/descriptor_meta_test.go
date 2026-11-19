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
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/plugintest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/render"
)

// TestOpDescriptorMeta checks that every language's generated Op descriptor
// carries the operation metadata authorization and audit need beyond the
// capability id: the resource template, the method name, the connection
// names, and the transport kind. github.issues.get is the reference op.
func TestOpDescriptorMeta(t *testing.T) {
	const wantRT = "repos/{owner}/{repository}"
	for _, tc := range []struct {
		lang   string
		file   string
		want   []string
		absent []string
	}{
		{
			lang: "go", file: "kei/connectors/github/v1/githubops/ops.go",
			want: []string{
				`ID: "github.issues.get"`,
				`ResourceTemplate: "repos/{owner}/{repository}"`,
				`MethodName: "GetIssue"`,
				`Connections: []string{"user_oauth", "pat"}`,
				`Transport: connection.TransportHTTP`,
			},
			absent: []string{"TransportHTTPTransport"},
		},
		{
			lang: "python", file: "kei/connectors/github/v1/ops.py",
			want: []string{
				`id="github.issues.get"`,
				`resource_template="repos/{owner}/{repository}"`,
				`method_name="GetIssue"`,
				`connections=("user_oauth", "pat")`,
				`transport=Transport.HTTP`,
			},
		},
		{
			lang: "ts", file: "kei/connectors/github/v1/ops.ts",
			want: []string{
				`id: "github.issues.get"`,
				`resourceTemplate: "repos/{owner}/{repository}"`,
				`methodName: "GetIssue"`,
				`connections: ["user_oauth","pat"]`,
				`transport: "http"`,
			},
		},
		{
			lang: "rust", file: "kei/connectors/github/v1/ops.rs",
			want: []string{
				`id: "github.issues.get"`,
				`resource_template: "repos/{owner}/{repository}"`,
				`method_name: "GetIssue"`,
				`connections: &["user_oauth", "pat"]`,
				`transport: Transport::Http`,
			},
		},
	} {
		t.Run(tc.lang, func(t *testing.T) {
			var opts render.Options
			gen := plugintest.Plugin(t, "../../../proto",
				[]string{"kei/connectors/github/v1/connector.proto", "kei/connectors/github/v1/issues.proto"},
				fmt.Sprintf("paths=source_relative,lang=%s", tc.lang), opts.Set)
			if err := render.Generate(gen, opts); err != nil {
				t.Fatalf("Generate: %v", err)
			}
			resp := gen.Response()
			if resp.Error != nil {
				t.Fatalf("plugin error: %s", resp.GetError())
			}
			var content string
			for _, f := range resp.File {
				if f.GetName() == tc.file {
					content = f.GetContent()
				}
			}
			if content == "" {
				t.Fatalf("no output %s; got %v", tc.file, keys(resp.File))
			}
			for _, w := range tc.want {
				if !strings.Contains(content, w) {
					t.Errorf("ops output missing %q", w)
				}
			}
			for _, a := range tc.absent {
				if strings.Contains(content, a) {
					t.Errorf("ops output must not contain %q", a)
				}
			}
		})
	}
}

func keys(files []*pluginpb.CodeGeneratorResponse_File) []string {
	var ks []string
	for _, f := range files {
		ks = append(ks, f.GetName())
	}
	return ks
}

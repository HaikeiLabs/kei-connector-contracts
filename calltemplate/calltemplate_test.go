// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package calltemplate

import (
	"strings"
	"testing"
)

const validTemplate = `[{
  "schema": "kei.tool-definition/v1",
  "name": "github.file_write",
  "description": "Write a file",
  "parameters": {
    "type": "object",
    "additionalProperties": false,
    "required": ["path", "content"],
    "properties": {
      "path": {"type": "string"},
      "content": {"type": "string"}
    }
  },
  "permission": "github.write",
  "category": "github",
  "service": "github",
  "binding": {
    "connector_id": "conn_github_1",
    "transport": {
      "method": "PUT",
      "url_template": "https://api.github.com/repos/{owner}/{repo}/contents/{path}",
      "auth": {"placement": "header", "name": "Authorization", "scheme": "Bearer"},
      "args": {
        "owner": {"in": "path"},
        "repo": {"in": "path"},
        "path": {"in": "path"},
        "content": {"in": "body"}
      }
    }
  }
}]`

func TestParseAcceptsValidTemplate(t *testing.T) {
	tpls, err := Parse([]byte(validTemplate))
	if err != nil {
		t.Fatalf("valid template rejected: %v", err)
	}
	if len(tpls) != 1 || tpls[0].Name != "github.file_write" {
		t.Fatalf("unexpected templates: %+v", tpls)
	}
}

func TestParseRejectsMalformedTemplates(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(string) string
		wantErr string
	}{
		{"not an array", func(s string) string { return strings.Trim(s, "[]") }, "array"},
		{"empty array", func(string) string { return "[]" }, "at least one"},
		{"unknown field", func(s string) string { return strings.Replace(s, `"service"`, `"secret_note": "x", "service"`, 1) }, "unknown field"},
		{"bad schema marker", func(s string) string {
			return strings.Replace(s, "kei.tool-definition/v1", "kei.tool-definition/v9", 1)
		}, "schema"},
		{"bad name", func(s string) string { return strings.Replace(s, `"github.file_write"`, `"GitHub File"`, 1) }, "name"},
		{"missing description", func(s string) string { return strings.Replace(s, `"description": "Write a file",`, "", 1) }, "description"},
		{"missing permission", func(s string) string { return strings.Replace(s, `"permission": "github.write",`, "", 1) }, "permission"},
		{"tool without a schema", func(s string) string {
			i := strings.Index(s, `"parameters"`)
			j := strings.Index(s, `"permission"`)
			return s[:i] + s[j:]
		}, "parameters"},
		{"open parameter schema", func(s string) string { return strings.Replace(s, `"additionalProperties": false,`, "", 1) }, "additionalProperties"},
		{"required not in properties", func(s string) string { return strings.Replace(s, `["path", "content"]`, `["path", "nope"]`, 1) }, "nope"},
		{"unknown method", func(s string) string { return strings.Replace(s, `"PUT"`, `"FETCH"`, 1) }, "method"},
		{"missing url_template", func(s string) string {
			return strings.Replace(s, `"url_template": "https://api.github.com/repos/{owner}/{repo}/contents/{path}",`, "", 1)
		}, "url_template"},
		{"placeholder without path arg", func(s string) string { return strings.Replace(s, `"repo": {"in": "path"},`, "", 1) }, "{repo}"},
		{"bad arg location", func(s string) string {
			return strings.Replace(s, `"content": {"in": "body"}`, `"content": {"in": "cookie"}`, 1)
		}, "cookie"},
		{"secret field in auth", func(s string) string {
			return strings.Replace(s, `"scheme": "Bearer"`, `"scheme": "Bearer", "token": "abc"`, 1)
		}, "token"},
		{"github token in url", func(s string) string {
			return strings.Replace(s, "api.github.com", "api.github.com/ghp_0123456789abcdefghij0123456789abcd", 1)
		}, "secret"},
		{"literal bearer value", func(s string) string {
			return strings.Replace(s, `"description": "Write a file"`, `"description": "send Authorization: Bearer abcdef0123456789"`, 1)
		}, "secret"},
		{"userinfo in url", func(s string) string {
			return strings.Replace(s, "https://api.github.com", "https://bob:hunter2@api.github.com", 1)
		}, "secret"},
		{"literal credential in query", func(s string) string {
			return strings.Replace(s, "contents/{path}", "contents/{path}?access_token=abc123", 1)
		}, "secret"},
		{"secret-named key in config", func(s string) string {
			return strings.Replace(s, `"connector_id"`, `"config": {"api_key": "abc123"}, "connector_id"`, 1)
		}, "secret"},
		{"duplicate name", func(s string) string { return "[" + strings.Trim(s, "[]") + "," + strings.Trim(s, "[]") + "]" }, "duplicate"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.mutate(validTemplate)))
			if err == nil {
				t.Fatalf("expected an error containing %q, got none", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestParseAcceptsPlaceholdersAndRegistrationOnlyTemplates(t *testing.T) {
	// A bearer scheme with a placeholder value is not a secret, and an action
	// tool may omit binding.transport (the harness runs a bounded action).
	src := `[{"schema":"kei.tool-definition/v1","name":"github.create_fix_pr","description":"d",
	  "parameters":[{"name":"title","required":true}],"permission":"github.write","category":"github",
	  "binding":{"connector_id":"c","delegated_context":["repository"]}}]`
	if _, err := Parse([]byte(src)); err != nil {
		t.Fatalf("registration-only template rejected: %v", err)
	}
}

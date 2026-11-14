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

package contracts_test

import (
	"regexp"
	"testing"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// mcpDottedID is the same dotted-id shape every operation id must pass
// (see registry_test.go): lowercase [a-z0-9_] segments, three or more.
var mcpDottedID = regexp.MustCompile(`^[a-z0-9_]+(\.[a-z0-9_]+){2,}$`)

// TestOperationIDForMCP pins the id rule for tools discovered from a remote
// MCP server: mcp.<server>.<tool>, with server and tool normalized so every
// character outside [a-z0-9_] becomes an underscore. An id from a remote MCP
// server enters the same policy id space as declared operations, so it must
// pass the same dotted-id shape.
func TestOperationIDForMCP(t *testing.T) {
	tests := []struct {
		name   string
		server string
		tool   string
		want   string
	}{
		{"plain names pass through", "grafana", "render_dashboard", "mcp.grafana.render_dashboard"},
		{"uppercase is lowercased", "GitHub", "CreateIssue", "mcp.github.createissue"},
		{"spaces become underscores", "my server", "create issue", "mcp.my_server.create_issue"},
		{"dashes and dots become underscores", "my-server", "tool.v1", "mcp.my_server.tool_v1"},
		{"every other rune becomes an underscore", "Küche", "tool", "mcp.k_che.tool"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := contracts.OperationIDForMCP(tt.server, tt.tool)
			if err != nil {
				t.Fatalf("OperationIDForMCP(%q, %q) returned error: %v", tt.server, tt.tool, err)
			}
			if got != tt.want {
				t.Errorf("OperationIDForMCP(%q, %q) = %q, want %q", tt.server, tt.tool, got, tt.want)
			}
			if !mcpDottedID.MatchString(got) {
				t.Errorf("OperationIDForMCP(%q, %q) = %q, which does not match the dotted-id pattern", tt.server, tt.tool, got)
			}
		})
	}
}

// TestOperationIDForMCPEmptyPartsRejected pins that an empty server or tool
// produces an error, never an id.
func TestOperationIDForMCPEmptyPartsRejected(t *testing.T) {
	tests := []struct {
		name   string
		server string
		tool   string
	}{
		{"empty server", "", "render_dashboard"},
		{"empty tool", "grafana", ""},
		{"both empty", "", ""},
		{"whitespace-only server", "  ", "render_dashboard"},
		{"whitespace-only tool", "grafana", "\t"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := contracts.OperationIDForMCP(tt.server, tt.tool)
			if err == nil {
				t.Fatalf("OperationIDForMCP(%q, %q) = %q, want an error for an empty part", tt.server, tt.tool, got)
			}
			if got != "" {
				t.Errorf("OperationIDForMCP(%q, %q) returned id %q alongside an error, want an empty id", tt.server, tt.tool, got)
			}
		})
	}
}

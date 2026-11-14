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

package contracts

import (
	"errors"
	"strings"
)

// Operation ids for tools discovered from a remote MCP server.
//
// A remote MCP server exposes tools that are not declared as fixed
// operations in this module. The harness assigns them ids in the same space
// as declared operations — mcp.<server>.<tool> — so policy and audit see one
// id space. Assignment happens in the harness: kei-proxy stays
// decision-only (ADR-039) and does no routing.
//
// The id rule is normalization, not composition of raw names: every
// character outside [a-z0-9_] in the server or tool name becomes an
// underscore, so an id always matches the dotted-id shape policy matches on.

// ErrEmptyMCPServerSegment and ErrEmptyMCPToolSegment are returned when a
// part needed for an MCP operation id is empty.
var (
	ErrEmptyMCPServerSegment = errors.New("mcp operation id: server segment is required")
	ErrEmptyMCPToolSegment   = errors.New("mcp operation id: tool segment is required")
)

// OperationIDForMCP returns the operation id for a tool call against a
// remote MCP server: mcp.<server>.<tool>, with server and tool normalized so
// that every character outside [a-z0-9_] becomes an underscore. The result
// therefore always matches the dotted-id pattern (lowercase [a-z0-9_]
// segments, three or more).
//
// An empty server or tool (after trimming spaces) is an error; whitespace is
// not a name.
func OperationIDForMCP(server, tool string) (string, error) {
	server = strings.TrimSpace(server)
	if server == "" {
		return "", ErrEmptyMCPServerSegment
	}
	tool = strings.TrimSpace(tool)
	if tool == "" {
		return "", ErrEmptyMCPToolSegment
	}
	return "mcp." + normalizeMCPSegment(server) + "." + normalizeMCPSegment(tool), nil
}

// normalizeMCPSegment maps every character outside [a-z0-9_] to an
// underscore, lowercasing as it goes. The result is never empty: the caller
// has rejected empty input, and any remaining rune maps to "_".
func normalizeMCPSegment(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

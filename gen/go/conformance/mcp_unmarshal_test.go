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

// MCP contract conformance (HAI-1001, plan task 2.3): the tool.call fixtures
// unmarshal into the generated Go types, JSON nulls inside the Struct fields
// survive, and the unknown extra fields each fixture carries are detected by
// strict unmarshal (the drift gate).
package conformance

import (
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	mcpv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/mcp/v1"
)

func TestMCPFixturesUnmarshal(t *testing.T) {
	reqData := readFixture(t, "mcp", "tool.call.request.json")
	var call mcpv1.ToolCall
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(reqData, &call); err != nil {
		t.Fatalf("request fixture did not unmarshal: %v", err)
	}
	if got := call.GetServer(); got != "metrics" {
		t.Errorf("server = %q, want metrics", got)
	}
	if got := call.GetTool(); got != "render_dashboard" {
		t.Errorf("tool = %q, want render_dashboard", got)
	}
	args := call.GetArguments()
	if got := args.GetFields()["dashboard_uid"].GetStringValue(); got != "synthetic-dash-1" {
		t.Errorf("arguments.dashboard_uid = %q, want synthetic-dash-1", got)
	}
	if got := args.GetFields()["panel_filter"]; got.GetNullValue() != structpb.NullValue_NULL_VALUE {
		t.Errorf("arguments.panel_filter = %v, want JSON null preserved", got)
	}
	if got := args.GetFields()["variables"].GetStructValue().GetFields()["region"]; got.GetNullValue() != structpb.NullValue_NULL_VALUE {
		t.Errorf("arguments.variables.region = %v, want JSON null preserved", got)
	}

	respData := readFixture(t, "mcp", "tool.call.response.json")
	var result mcpv1.ToolResult
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(respData, &result); err != nil {
		t.Fatalf("response fixture did not unmarshal: %v", err)
	}
	if got := len(result.GetContent()); got != 2 {
		t.Fatalf("len(content) = %d, want 2", got)
	}
	if got := result.GetContent()[0].GetFields()["type"].GetStringValue(); got != "text" {
		t.Errorf("content[0].type = %q, want text", got)
	}
	if got := result.GetContent()[1].GetFields()["resource"].GetStructValue().GetFields()["mimeType"]; got.GetNullValue() != structpb.NullValue_NULL_VALUE {
		t.Errorf("content[1].resource.mimeType = %v, want JSON null preserved", got)
	}
	if result.GetIsError() {
		t.Errorf("is_error = true, want false")
	}
}

// TestMCPFixturesStrictDetectsUnknown proves the drift gate on both MCP
// fixtures: strict unmarshal fails because each carries one unknown extra
// field, so a provider payload that no longer fits the declared shape is
// caught.
func TestMCPFixturesStrictDetectsUnknown(t *testing.T) {
	for name, msg := range map[string]proto.Message{
		"tool.call.request.json":  &mcpv1.ToolCall{},
		"tool.call.response.json": &mcpv1.ToolResult{},
	} {
		data := readFixture(t, "mcp", name)
		if err := (protojson.UnmarshalOptions{}).Unmarshal(data, msg); err == nil {
			t.Errorf("strict unmarshal of %s succeeded; the unknown extra field must stay unknown so the drift gate detects provider additions", name)
		}
	}
}

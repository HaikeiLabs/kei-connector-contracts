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

//! MCP contract conformance (HAI-1001, task 2.3): the tool.call fixtures
//! unmarshal into the generated prost types (serde impls from prost-serde
//! with ignore_unknown_fields, so the unknown extra fields each fixture
//! carries are skipped). Spot asserts prove fields and JSON nulls survive
//! with their synthetic values intact.

use haikeilabs_connector_contracts::kei::connectors::mcp::v1::{ToolCall, ToolResult};
use pbjson_types::value::Kind;

fn fixture(name: &str) -> String {
    let path = format!(
        "{}/../../fixtures/operations/mcp/{}",
        env!("CARGO_MANIFEST_DIR"),
        name
    );
    std::fs::read_to_string(&path).expect("read fixture")
}

#[test]
fn tool_call_unmarshals() {
    let call: ToolCall = serde_json::from_str(&fixture("tool.call.request.json"))
        .expect("tool.call request unmarshals");
    assert_eq!(call.server, "metrics");
    assert_eq!(call.tool, "render_dashboard");

    let arguments = call.arguments.expect("arguments present");
    assert_eq!(
        arguments
            .fields
            .get("dashboard_uid")
            .and_then(|v| v.kind.as_ref()),
        Some(&Kind::StringValue("synthetic-dash-1".to_string()))
    );

    // JSON nulls survive inside the Struct: pbjson-types maps JSON null to a
    // Value whose kind is None.
    match arguments
        .fields
        .get("panel_filter")
        .map(|v| v.kind.is_none())
    {
        Some(true) => {}
        other => panic!("panel_filter is not null: {:?}", other),
    }
    let variables = arguments
        .fields
        .get("variables")
        .and_then(|v| v.kind.as_ref())
        .expect("variables present");
    match variables {
        Kind::StructValue(s) => match s.fields.get("region").map(|v| v.kind.is_none()) {
            Some(true) => {}
            other => panic!("variables.region is not null: {:?}", other),
        },
        other => panic!("variables is not a struct: {:?}", other),
    }
}

#[test]
fn tool_result_unmarshals() {
    let result: ToolResult = serde_json::from_str(&fixture("tool.call.response.json"))
        .expect("tool.call response unmarshals");
    assert_eq!(result.content.len(), 2);
    assert_eq!(
        result.content[0]
            .fields
            .get("type")
            .and_then(|v| v.kind.as_ref()),
        Some(&Kind::StringValue("text".to_string()))
    );

    // JSON nulls survive inside the nested Struct of a content block.
    let resource = result.content[1]
        .fields
        .get("resource")
        .and_then(|v| v.kind.as_ref())
        .expect("resource present");
    match resource {
        Kind::StructValue(s) => match s.fields.get("mimeType").map(|v| v.kind.is_none()) {
            Some(true) => {}
            other => panic!("resource.mimeType is not null: {:?}", other),
        },
        other => panic!("resource is not a struct: {:?}", other),
    }
    assert!(!result.is_error);
}

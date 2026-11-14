# Copyright 2026 Haikei Labs
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""MCP contract conformance (HAI-1001, task 2.3): the tool.call fixtures
unmarshal into the generated python types.

Lenient parse (ignore_unknown_fields=True) proves the happy path, JSON nulls
inside the Struct fields survive, and one spot assert per message proves a
field survived unmarshal with its synthetic value intact. The strict parse of
both fixtures fails because each carries one unknown extra field — the same
drift gate the pilot operations run.
"""

from pathlib import Path

import pytest
from google.protobuf import json_format
from google.protobuf.struct_pb2 import NULL_VALUE

from kei.connectors.mcp.v1 import tool_pb2

FIXTURES = Path(__file__).resolve().parents[3] / "fixtures" / "operations" / "mcp"


def _fixture(name: str) -> str:
    return (FIXTURES / name).read_text()


def test_tool_call_unmarshals():
    call = json_format.Parse(
        _fixture("tool.call.request.json"),
        tool_pb2.ToolCall(),
        ignore_unknown_fields=True,
    )

    assert call.server == "metrics"
    assert call.tool == "render_dashboard"
    arguments = call.arguments
    assert arguments.fields["dashboard_uid"].string_value == "synthetic-dash-1"
    # JSON nulls survive as NullValue inside the Struct.
    assert arguments.fields["panel_filter"].null_value == NULL_VALUE
    assert arguments.fields["variables"].struct_value.fields["region"].null_value == NULL_VALUE
    assert [
        value.string_value for value in arguments.fields["tags"].list_value.values
    ] == [
        "latency",
        "p99",
    ]


def test_tool_result_unmarshals():
    result = json_format.Parse(
        _fixture("tool.call.response.json"),
        tool_pb2.ToolResult(),
        ignore_unknown_fields=True,
    )

    assert len(result.content) == 2
    assert result.content[0].fields["type"].string_value == "text"
    # JSON nulls survive inside the nested Struct of a content block.
    resource = result.content[1].fields["resource"].struct_value
    assert resource.fields["mimeType"].null_value == NULL_VALUE
    assert result.is_error is False


@pytest.mark.parametrize(
    "name,message",
    [
        ("tool.call.request.json", tool_pb2.ToolCall),
        ("tool.call.response.json", tool_pb2.ToolResult),
    ],
)
def test_strict_parse_detects_unknown_field(name, message):
    """Each MCP fixture carries one unknown extra field, so the strict parse
    (the drift gate) must fail instead of silently ignoring it."""
    with pytest.raises(json_format.ParseError):
        json_format.Parse(_fixture(name), message())

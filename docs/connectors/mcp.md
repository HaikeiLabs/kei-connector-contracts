# mcp connector

The `mcp` connector covers tool calls against remote MCP servers. Unlike
`http_api` and `database`, it declares no fixed operations: a remote server
exposes its own tools, and the set changes per server. This contract fixes
two things instead: how such a tool gets an operation id, so policy and
audit see one id space with the declared operations, and the payload shape
of a call and its result. The general rules are in
[operation-contracts.md](../operation-contracts.md).

Source: `proto/kei/connectors/mcp/v1/tool.proto`; the id rule is
`OperationIDForMCP` in `gen/go/kei/contracts/v1/mcp.go`.

## Operation ids

A tool on a remote MCP server has the operation id `mcp.<server>.<tool>`.

- `<server>` and `<tool>` are normalized: the name is lowercased and every
  character outside `[a-z0-9_]` becomes `_`. Server `my-server`, tool
  `Create Issue` is `mcp.my_server.create_issue`. The id therefore always
  matches the dotted-id shape policies match on.
- An empty server or tool name (after trimming spaces) is an error
  (`ErrEmptyMCPServerSegment`, `ErrEmptyMCPToolSegment`). The harness must
  reject the tool rather than mint an id.
- The harness assigns the id when it registers the tool. Policies target it
  like any other id, for example `connector:mcp.my_server.*`.

There is no class, HTTP method, path template, scope list or credential
placement in this file: the messages carry no `operation` option, because
the id is per tool, not per message.

## No runtime routing

The harness calls the MCP server. kei-proxy stays decision-only (ADR-039):
it decides whether `mcp.<server>.<tool>` is allowed and does no routing,
and nothing in this repository generates a client that relays the call.

## Messages

| Message | Field | Type |
| --- | --- | --- |
| `ToolCall` | `server` | `string` (MCP-facing name, not normalized) |
| | `tool` | `string` (MCP-facing name, not normalized) |
| | `arguments` | `google.protobuf.Struct` |
| `ToolResult` | `content` | `repeated google.protobuf.Struct` |
| | `is_error` | `bool` (the MCP `isError` flag; false means success) |

`arguments` and each `content` block are `google.protobuf.Struct`: tool
schemas vary per server, so the input and the returned content blocks stay
free-form JSON. JSON `null` inside them survives as a null `Value`.

## Unmarshal provider JSON

Unmarshal leniently at run time. The snippets follow the repository's MCP
conformance tests.

Go:

```go
import (
	"google.golang.org/protobuf/encoding/protojson"

	mcpv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/mcp/v1"
)

var result mcpv1.ToolResult
if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(body, &result); err != nil {
	return err
}
```

Python (package root `gen/python/src`):

```python
from google.protobuf import json_format
from kei.connectors.mcp.v1 import tool_pb2

result = json_format.Parse(body, tool_pb2.ToolResult(), ignore_unknown_fields=True)
```

TypeScript (protobuf-es v2, generated into `gen/ts/src`):

```ts
import { fromJsonString } from "@bufbuild/protobuf";
import { ToolResultSchema } from "./kei/connectors/mcp/v1/tool_pb";

const result = fromJsonString(ToolResultSchema, body, { ignoreUnknownFields: true });
```

Rust (`ignore_unknown_fields` is set at generation time in `buf.gen.yaml`):

```rust
use haikeilabs_connector_contracts::kei::connectors::mcp::v1::ToolResult;

let result: ToolResult = serde_json::from_str(&body)?;
```

The TypeScript snippet uses the protobuf-es v2 API; the repository has no
TypeScript test yet.

## Fixtures

`fixtures/operations/mcp/`:

- `tool.call.request.json`: a synthetic `ToolCall` with nested objects,
  arrays and `null` values in `arguments`.
- `tool.call.response.json`: a `ToolResult` with a text block and a
  resource block.

Each carries one unknown top-level field (`session_extra`, `trace_extra`)
so strict unmarshal can prove it detects a provider addition.

## Testing

On every pull request:

- Go: `gen/go/conformance/mcp_unmarshal_test.go` (lenient unmarshal, and
  strict unmarshal must fail on both fixtures) and
  `gen/go/kei/contracts/v1/mcp_test.go` (the id rule).
- Python: `gen/python/tests/test_mcp_unmarshal.py`.
- Rust: `gen/rust/tests/mcp_unmarshal.rs`.

There is no live MCP server in CI: the harness owns the call, so the
contract is proven on fixtures only.

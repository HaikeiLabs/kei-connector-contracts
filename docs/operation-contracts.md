# Operation contracts

Accepted 2026-11-11 via
[proposal 0001](proposals/0001-operation-contracts.md) with the owner
amendment recorded there. This document is the rules an operation contract
follows. The contract stays metadata only: it describes API calls and never
carries provider payloads, results, customer content, or credential values,
and the root Go module stays stdlib-only.

## Operation ids

An operation id is a dotted name: `connector.resource.verb`, for example
`github.issues.get`.

- Segments are lowercase snake_case, and there are at least three segments.
- A tool discovered from a remote MCP server is `mcp.<server>.<tool>`, and a
  harness built-in is `builtin.<name>`.
- The `<server>` and `<tool>` segments of an MCP id are normalized: every
  character outside `[a-z0-9_]` becomes `_`, so an id always matches the
  dotted-id shape (for example server `my-server`, tool `Create Issue` is
  `mcp.my_server.create_issue`). An empty server or tool is an error — the
  harness must reject it rather than mint an id. Assignment happens in the
  harness when the tool is registered; kei-proxy stays decision-only
  (ADR-039) and does no routing. MCP payloads are typed in
  `proto/kei/connectors/mcp/v1/tool.proto` (`ToolCall`/`ToolResult`, free
  form `google.protobuf.Struct` arguments and content blocks).
- Ids are policy targets under the existing dot-segment matching, with `*`
  as a segment wildcard: `connector:github.*` matches `github.issues.get`.
  No policy schema change is needed.

## Layout

- `proto/kei/contracts/v1/options.proto` holds the operation message option
  (the id, scopes, and `read`/`write` class fields).
- Each connector is its own file:
  `proto/kei/connectors/<connector>/v1/<resource>.proto`, containing
  `<Op>Request` and `<Op>Response` message pairs for that resource.

buf generates plain Go structs, Python, TS and Rust types from the same
files. proto is the contract language only: there are no gRPC services, no
grpc-gateway, and no generated HTTP clients that relay calls. The harness
makes every call; kei-proxy stays decision-only (ADR-039).

## json_name

Set `json_name` on every field whose provider JSON key differs from
protojson's lowerCamel default. One generated type then unmarshals the
provider's real payload without a second mapping layer.

## Value encoding

Fields typed `google.protobuf.Value` or `google.protobuf.ListValue` carry raw
JSON, so a fixture has to use a JSON form that survives the round trip:

- `NULL` is JSON `null`.
- A 64-bit integer above 2^53 is a JSON string (`"9007199254740993"`), never a
  JSON number: JSON numbers parse as `float64` and lose the low digits.
- A `NUMERIC`/`DECIMAL` value is a JSON string
  (`"12345678901234567890.123456789012345678"`) so its precision and scale are
  preserved.

`RunQueryResponse.rows` are positional: each row lists values in `columns`
order.

## Unknown fields

Runtime unmarshal discards unknown fields, so a provider adding a key does
not break a release. The nightly drift check runs the same unmarshal strict
instead, and files a ticket when a provider payload no longer fits the
declared shape.

## Versioning

Within v1, changes are additive only, enforced mechanically by
`buf breaking`. A breaking change takes a v2, not an edit in place — the
same rule the frozen Go contract follows
([CONTRIBUTING.md](../CONTRIBUTING.md)).

## Fixtures

Fixtures are synthetic or sanitized only. Never check in real customer data,
and never a credential value — an API key, token, password, or connection
string is not a fixture. Fixtures exist so generated types can be proven to
unmarshal real provider shapes; the HAI-990 gate is 10 pilot operations whose
fixtures unmarshal in Go, Python and Rust.

## Scope

The existing `kei.tool-definition/v1` schemas, the call template profile
(`docs/tool-call-templates.md`, `calltemplate`), and `providers/` are
unchanged by this layer. Operation contracts add a source of truth; they do
not replace or deprecate anything.

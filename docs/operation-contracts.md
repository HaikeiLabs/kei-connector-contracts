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

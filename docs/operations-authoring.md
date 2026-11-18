# Add an operation

The proto is the source for operation and connection metadata. See
[proposal 0002](proposals/0002-connector-operations.md), including its owner
decisions, for the design. Some generation and release gates are still being
implemented.

Connector pages list their connections and operations:
[Airtable](connectors/airtable.md), [Atlassian](connectors/atlassian.md),
[database](connectors/database.md), [GitHub](connectors/github.md),
[Gmail](connectors/gmail.md), [Google Sheets](connectors/google_sheets.md),
[HTTP API](connectors/http_api.md), [HubSpot](connectors/hubspot.md),
[Linear](connectors/linear.md), [Notion](connectors/notion.md),
[Omni](connectors/omni.md), [Power BI](connectors/power_bi.md), and
[Sentry](connectors/sentry.md).

## Declare the contract

1. Add request and response messages under
   `proto/kei/connectors/<connector>/v1/`. Use `json_name` when a provider key
   differs from the default. Declare a unique `connector.resource.verb` id;
   it is also the capability id. Writes need a scope; mark deletes or
   overwrites `destructive`.
2. Declare the connector once with the `connector` file option in
   `connector.proto`: `id`, default `base_url`, and `connections`. Each
   connection has a stable `name`, `kind`, `auth_method`, `account_model`,
   optional `base_url` or `base_url_field`, `setup_fields` and `scopes`.
   Setup fields have `name`, `type`, optional `pattern` and `secret`.
   `engine` limits SQL DSNs; `header_name` names a custom credential header.
   Describe secrets, never put their values in proto.
3. Attach the `operation` option to the request. This example is copied from `proto/kei/connectors/github/v1/issues.proto`:

```proto
message GetIssueRequest {
  string owner = 1 [(kei.contracts.v1.param) = {in: PARAM_IN_PATH}];
  string repository = 2 [(kei.contracts.v1.param) = {in: PARAM_IN_PATH}];
  int64 issue_number = 3 [(kei.contracts.v1.param) = {in: PARAM_IN_PATH}];

  option (kei.contracts.v1.operation) = {
    id: "github.issues.get"
    connector: "github"
    class: OPERATION_CLASS_READ
    http_method: "GET"
    path_template: "/repos/{owner}/{repository}/issues/{issue_number}"
    scopes: ["issues:read"]
    credential_placement: CREDENTIAL_PLACEMENT_HEADER_BEARER
    method_name: "GetIssue"
    connections: ["user_oauth", "pat"]
    resource_template: "repos/{owner}/{repository}"
    semantics: {action: SEMANTIC_ACTION_READ resource_type: "issue" version: 1}
    errors: [
      {status: 401 kind: ERROR_KIND_UNAUTHORIZED},
      {status: 403 kind: ERROR_KIND_FORBIDDEN},
      {status: 404 kind: ERROR_KIND_NOT_FOUND},
      {status: 429 kind: ERROR_KIND_RATE_LIMITED retryable: true}
    ]
  };
}
```

The options include `method_name` (default: request name without `Request`),
accepted `connections`, `resource_template` (policy resource with request
field variables), `errors` (status, kind, retryable), and `semantics` (action,
resource type, version). Bump the semantics version when its mapping changes;
it is metadata, not a grant. `credential_placement` remains for compatibility;
the connection kind determines placement.

Use the `param` field option to override placement with `in` (path, query,
body or header) and optional wire `name`. By default, template fields go in
the path, other fields in the query for GET and the body otherwise. Set
`multi_segment` only for path values containing slashes: each segment is
escaped separately, slashes stay, and `.` and `..` segments are rejected.

For pagination, set `cursor_field` on the request and `next_cursor_field` on
the response. Each call returns one page; callers loop until the next cursor
is empty. For GraphQL, set `graphql_document` to one reviewed fixed document
with typed request fields as variables; check it against the pinned provider
schema. Set `positional_helper` to add a positional helper beside the default
message-based method; `signature=message|positional` selects the package form.
The database operation uses a request dialect and a SQL connection.

## Fixtures and checks

Add synthetic request and response JSON under `fixtures/operations/<connector>/`.
Include `null` where applicable and an extra provider field to exercise
unknown-field handling. Do not include credentials or customer data. For raw
JSON values, encode integers above 2^53 as strings to preserve precision;
encode exact decimal values as strings too. Add or update fixture unmarshal
and replay coverage for the operation, including each supported language.
Update its connector page and any applicable legacy capability mapping.

`buf generate` runs the message plugins and `protoc-gen-kei-ops` with all
protos for each configured language, validating operations and emitting
interfaces, implementations and descriptors. `protoc-gen-kei-manifest` emits
the connector and capability manifest using the legacy mapping file. Commit
generated changes. The current configuration generates operations for Go,
Python and TypeScript; Rust operation generation and four-language replay
parity are planned. Rust message types are already generated.

Before pushing, run `buf lint`, `buf breaking --against '.git#branch=main'`,
`buf generate`, and verify the generated diff. Run the generator tests with
`go -C tools/protoc-gen-kei test ./...`, Go conformance tests with
`go -C gen/go test ./...`, Python tests with
`PYTHONPATH=gen/python/src python -m pytest -q gen/python/tests`, Rust tests
with `cargo test --manifest-path gen/rust/Cargo.toml`, and TypeScript tests
and build with `npm --prefix gen/ts test` and `npm --prefix gen/ts run build`.
See [CI](../.github/workflows/ci.yaml) for the full checks.

## Version and deprecate

Capability ids are permanent: never reuse an id for different behavior.
Additive changes fit within the current proto version; breaking message changes
need a new version. To rename an operation, add the new id, mark the old one
`deprecated`, set `replaced_by` to the new id, and add a capability mapping row
from old to new without widening grants. Keep the old operation callable for
at least two minor releases; remove it only in a major release after the
catalog reports no policy still names it. Release language packages together.

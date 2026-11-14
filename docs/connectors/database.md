# database connector

The `database` connector runs one read-only SQL statement against a
connected database and returns the result set. It is read-only: no DDL and
no writes. The request names the SQL dialect so the runtime binds parameters
correctly for that engine. The general rules are in
[operation-contracts.md](../operation-contracts.md).

Source: `proto/kei/connectors/database/v1/query.proto`.

## Operations

| Id | Class | HTTP method | Path template | Scopes | Credential placement |
| --- | --- | --- | --- | --- | --- |
| `database.query.run` | read | `POST` | `/query` | `query:read` | none |

The credential placement is `CREDENTIAL_PLACEMENT_NONE`: the connection is
not an HTTP request with a header credential. No connection string or
password is ever part of the contract or a fixture.

## Messages

| Message | Field | Type |
| --- | --- | --- |
| `RunQueryRequest` | `dialect` | `Dialect` enum |
| | `sql` | `string` |
| | `params` | `repeated google.protobuf.Value` |
| | `max_rows` | `optional int64` |
| `Column` | `name` | `string` |
| | `type` | `string` (the engine's type name) |
| `RunQueryResponse` | `columns` | `repeated Column` |
| | `rows` | `repeated google.protobuf.ListValue` |
| | `truncated` | `bool` |

`Dialect` values: `DIALECT_POSTGRES`, `DIALECT_MYSQL`, `DIALECT_MSSQL`,
`DIALECT_SNOWFLAKE`, `DIALECT_SUPABASE` (zero value `DIALECT_UNSPECIFIED`).

`params` are positional bind parameters, each a `google.protobuf.Value`.
`rows` are positional too: each row is a `google.protobuf.ListValue` with
one value per entry in `columns`, in column order. `truncated` is true when
`max_rows` cut the result set short.

Because cells are raw JSON values, follow the
[value encoding](../operation-contracts.md#value-encoding) rules:

- SQL `NULL` is JSON `null`.
- An integer above 2^53 is a JSON string (`"9007199254740993"`).
- A `NUMERIC`/`DECIMAL` is a JSON string
  (`"12345678901234567890.123456789012345678"`).

## Unmarshal provider JSON

Unmarshal leniently at run time. The snippets follow the repository's
conformance tests.

Go:

```go
import (
	"google.golang.org/protobuf/encoding/protojson"

	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
)

var resp databasev1.RunQueryResponse
if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(body, &resp); err != nil {
	return err
}
```

Python (package root `gen/python/src`):

```python
from google.protobuf import json_format
from kei.connectors.database.v1 import query_pb2

resp = json_format.Parse(body, query_pb2.RunQueryResponse(), ignore_unknown_fields=True)
```

TypeScript (protobuf-es v2, generated into `gen/ts/src`):

```ts
import { fromJsonString } from "@bufbuild/protobuf";
import { RunQueryResponseSchema } from "./kei/connectors/database/v1/query_pb";

const resp = fromJsonString(RunQueryResponseSchema, body, { ignoreUnknownFields: true });
```

Rust (`ignore_unknown_fields` is set at generation time in `buf.gen.yaml`):

```rust
use haikeilabs_connector_contracts::kei::connectors::database::v1::RunQueryResponse;

let resp: RunQueryResponse = serde_json::from_str(&body)?;
```

The TypeScript snippet uses the protobuf-es v2 API; the repository has no
TypeScript test yet.

## Fixtures

`fixtures/operations/database/`:

- `query.run.request.json`: a `DIALECT_POSTGRES` query with one bind
  parameter above 2^53, passed as a string.
- `query.run.response.json`: four columns and two rows covering a bigint id
  above 2^53, unicode text, a `null` cell and a high-precision `numeric`. It
  carries one unknown extra field (`execution_time_ms`) for the strict drift
  check.

## Testing

Fixture unmarshal runs on every pull request: Go
(`gen/go/conformance/unmarshal_test.go`, lenient and strict), Python
(`gen/python/tests/test_unmarshal.py`) and Rust
(`gen/rust/tests/unmarshal.rs`).

The contract is also proven against real databases. The `Database testenv`
job in `.github/workflows/ci.yaml` runs on every pull request to `main`, on
pushes to `main`, and nightly. It starts the `testenv/compose.yaml` stack
with four services — Postgres, Supabase (a second Postgres standing in for
Supabase-local), MySQL and MSSQL — waits for them to be healthy and for the
MSSQL seed sidecar, then runs:

```sh
go -C gen/go test -count=1 -tags testenv ./conformance -run TestEnvDatabaseQueryRun -v
```

`TestEnvDatabaseQueryRun` (`gen/go/conformance/testenv_test.go`) runs the
same three queries on each service — one with NULLs, one with a bigint above
2^53, one with NUMERIC precision — through the client inside each container
(`psql`, `mysql`, `sqlcmd`), converts the result to the `RunQueryResponse`
JSON shape and unmarshals it into the generated type. The contract modules
gain no database driver, and the file only builds with the `testenv` build
tag, so ordinary `go test` runs leave it out.

The seed data (`testenv/seed/`) is synthetic: 101 accounts and 301 issues
per database, including NULLs, unicode text and the id `9007199254740993`.
The stack's passwords are local placeholders only.

To run it locally, see [testenv/README.md](../../testenv/README.md). MSSQL is
behind the `mssql` compose profile and its image is amd64-only, so locally it
needs an amd64 host; CI runs on an amd64 runner.

# http_api connector

The `http_api` connector fronts a generic REST API whose collections and
record shapes belong to the provider, not to this repository. The contract
declares how a record is listed or read and how the call is governed; the
record body itself stays open (`google.protobuf.Struct`), so one generated
type unmarshals any collection while policy and audit still see a typed
call. The general rules are in [operation-contracts.md](../operation-contracts.md).

Source: `proto/kei/connectors/http_api/v1/records.proto`.

## Operations

| Id | Class | HTTP method | Path template | Scopes | Credential placement |
| --- | --- | --- | --- | --- | --- |
| `http_api.records.list` | read | `GET` | `/{collection}` | `records:read` | header bearer |
| `http_api.records.get` | read | `GET` | `/{collection}/{record_id}` | `records:read` | header bearer |

Both operations are reads. The credential is placed as a bearer header by
the connector runtime at invocation time; it is never part of the contract
or a fixture.

## Messages

| Message | Field | Type |
| --- | --- | --- |
| `ListRecordsRequest` | `collection` | `string` |
| | `limit` | `optional int32` |
| | `cursor` | `optional string` |
| `ListRecordsResponse` | `records` | `repeated google.protobuf.Struct` |
| | `next_cursor` | `optional string` |
| `GetRecordRequest` | `collection` | `string` |
| | `record_id` | `string` |
| `GetRecordResponse` | `record` | `google.protobuf.Struct` |

`records` and `record` are `google.protobuf.Struct`: free-form JSON objects.
A JSON `null` inside a record survives as a null `Value`. Follow the
[value encoding](../operation-contracts.md#value-encoding) rules for large
integers and decimals inside a record.

## Unmarshal provider JSON

Unmarshal leniently at run time so a provider adding a key does not break a
release. The snippets follow the repository's conformance tests.

Go (`google.golang.org/protobuf/encoding/protojson`):

```go
import (
	"google.golang.org/protobuf/encoding/protojson"

	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
)

var resp httpapiv1.ListRecordsResponse
if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(body, &resp); err != nil {
	return err
}
```

Python (`google.protobuf.json_format`, package root `gen/python/src`):

```python
from google.protobuf import json_format
from kei.connectors.http_api.v1 import records_pb2

resp = json_format.Parse(body, records_pb2.ListRecordsResponse(), ignore_unknown_fields=True)
```

TypeScript (protobuf-es v2, generated into `gen/ts/src`):

```ts
import { fromJsonString } from "@bufbuild/protobuf";
import { ListRecordsResponseSchema } from "./kei/connectors/http_api/v1/records_pb";

const resp = fromJsonString(ListRecordsResponseSchema, body, { ignoreUnknownFields: true });
```

Rust (prost types with pbjson serde impls; `ignore_unknown_fields` is set at
generation time in `buf.gen.yaml`, so plain `serde_json` is lenient):

```rust
use haikeilabs_connector_contracts::kei::connectors::http_api::v1::ListRecordsResponse;

let resp: ListRecordsResponse = serde_json::from_str(&body)?;
```

The TypeScript snippet uses the protobuf-es v2 API; the repository has no
TypeScript test yet. Go, Python and Rust are covered by the tests below.

## Fixtures

`fixtures/operations/http_api/`:

- `records.list.request.json`, `records.list.response.json`
- `records.get.request.json`, `records.get.response.json`

Fixtures are synthetic. Each response fixture carries one unknown extra
field (`total_count`, `etag`) so the strict drift check can prove it detects
a provider addition.

## Testing

- Go: `gen/go/conformance/unmarshal_test.go` unmarshals every fixture with
  `DiscardUnknown`, and strict unmarshal must fail on the response fixtures.
- Python: `gen/python/tests/test_unmarshal.py` (`ignore_unknown_fields=True`).
- Rust: `gen/rust/tests/unmarshal.rs`.

CI runs them on every pull request: Python in the `Build and test` job, Go
(`go test ./...` in `gen/go`) and Rust (`cargo test` in `gen/rust`) in the
`Proto toolchain` job.

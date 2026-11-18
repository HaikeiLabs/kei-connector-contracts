# omni connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/omni/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `api_key` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `base_url` (url), `api_key` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `omni.documents.list` | `GET` | `/api/v1/documents` | `read` | false | `api_key` | `documents` |
| `omni.queries.run` | `POST` | `/api/v1/query/run` | `read` | false | `api_key` | `models/{model_id}` |

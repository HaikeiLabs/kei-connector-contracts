# airtable connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/airtable/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `pat` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `token` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `airtable.records.get` | `GET` | `/v0/{base_id}/{table_id_or_name}/{record_id}` | `read` | false | `pat` | `bases/{base_id}` |
| `airtable.records.list` | `GET` | `/v0/{base_id}/{table_id_or_name}` | `read` | false | `pat` | `bases/{base_id}` |

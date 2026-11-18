# google_sheets connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/google_sheets/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |
| `domain_delegation` | `CONNECTION_KIND_HTTP_BEARER` | `service_principal_key` | `ACCOUNT_MODEL_SHARED` | `service_account_json` (string, secret), `subject_email` (string) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `google_sheets.spreadsheets.get` | `GET` | `/v4/spreadsheets/{spreadsheet_id}` | `read` | false | `user_oauth`, `domain_delegation` | `spreadsheets/{spreadsheet_id}` |
| `google_sheets.values.get` | `GET` | `/v4/spreadsheets/{spreadsheet_id}/values/{range}` | `read` | false | `user_oauth`, `domain_delegation` | `spreadsheets/{spreadsheet_id}` |

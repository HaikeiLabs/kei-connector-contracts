# Notion connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/notion/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |
| `integration_token` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `token` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `notion.databases.query` | `POST` | `/v1/databases/{database_id}/query` | `read` | false | `user_oauth`, `integration_token` | `databases/{database_id}` |
| `notion.pages.get` | `GET` | `/v1/pages/{page_id}` | `read` | false | `user_oauth`, `integration_token` | `pages/{page_id}` |
| `notion.search` | `POST` | `/v1/search` | `read` | false | `user_oauth`, `integration_token` | `search` |

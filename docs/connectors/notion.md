# Notion connector

## Operations

| Operation id | Method | Path | Class | Connections |
| --- | --- | --- | --- | --- |
| `notion.databases.query` | `POST` | `/v1/databases/{database_id}/query` | `read` | `user_oauth`, `integration_token` |
| `notion.pages.get` | `GET` | `/v1/pages/{page_id}` | `read` | `user_oauth`, `integration_token` |

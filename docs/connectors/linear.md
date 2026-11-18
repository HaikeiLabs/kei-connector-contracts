# linear connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/linear/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |
| `api_key` | `CONNECTION_KIND_HTTP_HEADER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `api_key` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `linear.comments.create` | `POST` | `/graphql` (GraphQL) | `write` | false | `user_oauth`, `api_key` | `issues/{issue_id}` |
| `linear.issues.create` | `POST` | `/graphql` (GraphQL) | `write` | false | `user_oauth`, `api_key` | `teams/{team_id}` |
| `linear.issues.get` | `POST` | `/graphql` (GraphQL) | `read` | false | `user_oauth`, `api_key` | `issues/{id}` |
| `linear.issues.search` | `POST` | `/graphql` (GraphQL) | `read` | false | `user_oauth`, `api_key` | `issues` |
| `linear.teams.get` | `POST` | `/graphql` (GraphQL) | `read` | false | `user_oauth`, `api_key` | `teams/{team_id}` |

# gmail connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/gmail/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |
| `domain_delegation` | `CONNECTION_KIND_HTTP_BEARER` | `service_principal_key` | `ACCOUNT_MODEL_SHARED` | `service_account_json` (string, secret), `subject_email` (string) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `gmail.messages.get` | `GET` | `/gmail/v1/users/{user_id}/messages/{message_id}` | `read` | false | `user_oauth`, `domain_delegation` | `users/{user_id}` |
| `gmail.messages.list` | `GET` | `/gmail/v1/users/{user_id}/messages` | `read` | false | `user_oauth`, `domain_delegation` | `users/{user_id}` |

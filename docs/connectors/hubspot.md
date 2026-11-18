# hubspot connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/hubspot/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |
| `private_app` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `token` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `hubspot.contacts.get` | `GET` | `/crm/v3/objects/contacts/{contact_id}` | `read` | false | `user_oauth`, `private_app` | `contacts/{contact_id}` |
| `hubspot.contacts.list` | `GET` | `/crm/v3/objects/contacts` | `read` | false | `user_oauth`, `private_app` | `contacts` |
| `hubspot.deals.list` | `GET` | `/crm/v3/objects/deals` | `read` | false | `user_oauth`, `private_app` | `deals` |

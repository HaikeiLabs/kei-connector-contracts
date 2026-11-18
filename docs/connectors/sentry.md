# Sentry connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/sentry/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `auth_token` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `token` (string, secret), `organization_slug` (string) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `sentry.issues.get` | `GET` | `/api/0/organizations/{organization_slug}/issues/{issue_id}/` | `read` | false | `auth_token` | `organizations/{organization_slug}` |
| `sentry.issues.list` | `GET` | `/api/0/organizations/{organization_slug}/issues/` | `read` | false | `auth_token` | `organizations/{organization_slug}` |

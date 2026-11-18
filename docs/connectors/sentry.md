# Sentry connector

## Operations

| Operation id | Method | Path | Class | Connections |
| --- | --- | --- | --- | --- |
| `sentry.issues.get` | `GET` | `/api/0/organizations/{organization_slug}/issues/{issue_id}/` | `read` | `auth_token` |
| `sentry.issues.list` | `GET` | `/api/0/organizations/{organization_slug}/issues/` | `read` | `auth_token` |

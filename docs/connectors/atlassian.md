# atlassian connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/atlassian/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | `site_url` (url) |
| `api_token` | `CONNECTION_KIND_HTTP_BASIC` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `site_url` (url), `email` (string), `token` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `atlassian.confluence.pages.get` | `GET` | `/wiki/api/v2/pages/{page_id}` | `read` | false | `user_oauth`, `api_token` | `confluence/pages/{page_id}` |
| `atlassian.confluence.pages.list` | `GET` | `/wiki/api/v2/spaces/{space_id}/pages` | `read` | false | `user_oauth`, `api_token` | `confluence/pages` |
| `atlassian.jira.issues.get` | `GET` | `/rest/api/3/issue/{issue_key}` | `read` | false | `user_oauth`, `api_token` | `jira/issues/{issue_key}` |
| `atlassian.jira.issues.search` | `POST` | `/rest/api/3/search/jql` | `read` | false | `user_oauth`, `api_token` | `jira/issues` |

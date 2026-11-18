# github connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/github/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `sso_reuse` | `ACCOUNT_MODEL_PER_USER` | — |
| `pat` | `CONNECTION_KIND_HTTP_BEARER` | `secret_ref` | `ACCOUNT_MODEL_SHARED` | `token` (string, secret) |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `github.actions.runs.list` | `GET` | `/repos/{owner}/{repository}/actions/runs` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.contents.put` | `PUT` | `/repos/{owner}/{repository}/contents/{path}` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.blobs.create` | `POST` | `/repos/{owner}/{repository}/git/blobs` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.commits.create` | `POST` | `/repos/{owner}/{repository}/git/commits` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.refs.create` | `POST` | `/repos/{owner}/{repository}/git/refs` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.refs.get` | `GET` | `/repos/{owner}/{repository}/git/ref/{ref}` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.refs.update` | `PATCH` | `/repos/{owner}/{repository}/git/refs/{ref}` | `write` | true | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.git.trees.create` | `POST` | `/repos/{owner}/{repository}/git/trees` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.issues.comments.create` | `POST` | `/repos/{owner}/{repository}/issues/{issue_number}/comments` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.issues.create` | `POST` | `/repos/{owner}/{repository}/issues` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.issues.get` | `GET` | `/repos/{owner}/{repository}/issues/{issue_number}` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.issues.list` | `GET` | `/repos/{owner}/{repository}/issues` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.create` | `POST` | `/repos/{owner}/{repository}/pulls` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.get` | `GET` | `/repos/{owner}/{repository}/pulls/{pull_number}` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.list` | `GET` | `/repos/{owner}/{repository}/pulls` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.reviews.create` | `POST` | `/repos/{owner}/{repository}/pulls/{pull_number}/reviews` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.reviews.list` | `GET` | `/repos/{owner}/{repository}/pulls/{pull_number}/reviews` | `read` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |
| `github.pulls.update` | `PATCH` | `/repos/{owner}/{repository}/pulls/{pull_number}` | `write` | false | `user_oauth`, `pat` | `repos/{owner}/{repository}` |

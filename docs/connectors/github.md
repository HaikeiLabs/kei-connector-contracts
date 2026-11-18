# github connector

## Operations

| Operation id | Method | Path | Class | Connections |
| --- | --- | --- | --- | --- |
| `github.contents.put` | `PUT` | `/repos/{owner}/{repository}/contents/{path}` | `write` | `user_oauth`, `pat` |
| `github.issues.get` | `GET` | `/repos/{owner}/{repository}/issues/{issue_number}` | `read` | `user_oauth`, `pat` |
| `github.issues.list` | `GET` | `/repos/{owner}/{repository}/issues` | `read` | `user_oauth`, `pat` |
| `github.pulls.get` | `GET` | `/repos/{owner}/{repository}/pulls/{pull_number}` | `read` | `user_oauth`, `pat` |

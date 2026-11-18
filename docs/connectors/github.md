# github connector

## Operations on main

The table follows the proto options on `main`; `—` means no connection is declared there.

| Operation id | Method | Path | Class | Connections |
| --- | --- | --- | --- | --- |
| `github.contents.put` | `PUT` | `/repos/{owner}/{repository}/contents/{path}` | `write` | — |
| `github.issues.get` | `GET` | `/repos/{owner}/{repository}/issues/{issue_number}` | `read` | — |
| `github.issues.list` | `GET` | `/repos/{owner}/{repository}/issues` | `read` | — |
| `github.pulls.get` | `GET` | `/repos/{owner}/{repository}/pulls/{pull_number}` | `read` | — |

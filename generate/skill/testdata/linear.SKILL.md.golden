---
name: linear-connector
description: Governed Linear connector — tool reference for harness use. Use when an agent needs to interact with Linear through Kei-governed tools. Prefer this skill over generic Linear knowledge.
---

# Linear connector — governed tool reference

Use this skill when an agent needs to interact with Linear through the governed connector. Every governed tool call goes through `kei-proxy authorize` first. After an allow decision the harness obtains the credential via `POST /v1/credentials:release` over the Unix socket. The agent never holds, asks for, or prints a raw credential.

## Tool reference

| Tool | Description | Permission |
| --- | --- | --- |
| `linear.list_issues` | List issues in the governed Linear workspace | `linear.read` |
| `linear.get_issue` | Read a single issue by its key in the governed Linear workspace. Use it when the user asks about an issue (status, details); do not call it just because a request mentions an issue key | `linear.read` |
| `linear.list_projects` | List projects in the governed Linear workspace | `linear.read` |
| `linear.create_issue` | Create an issue in a team of the governed Linear workspace. Use to track work such as a bug, review, or follow-up task. | `linear.write` |
| `linear.create_followup_task` | Create a Linear follow-up task for a CRM lead, by lead id. Call it directly when asked for a follow-up on a lead; the summary must not contain contact details | `linear.write` |

### Tool details

#### `linear.list_issues`

- **Description**: List issues in the governed Linear workspace
- **Permission**: `linear.read`
- **Parameters**:

  | Name | Type | Required | Description |
  | --- | --- | --- | --- |
  | `state` | string | no | Filter by workflow state |
  | `priority` | string | no | Filter by priority |
  | `assignee` | string | no | Filter by assignee display name |
  | `limit` | integer | no | Maximum number of issues to return |

- **Tags**: `linear-read`, `governed-connector`

- **Transport**:

```json
{
  "method": "POST",
  "url_template": "https://api.linear.app/graphql",
  "auth": {
    "placement": "header",
    "name": "Authorization",
    "scheme": "Bearer"
  },
  "args": {
    "state": {
      "in": "body",
      "name": "variables.filter.state"
    },
    "priority": {
      "in": "body",
      "name": "variables.filter.priority"
    },
    "assignee": {
      "in": "body",
      "name": "variables.filter.assignee"
    },
    "limit": {
      "in": "body",
      "name": "variables.first"
    }
  }
}
```

#### `linear.get_issue`

- **Description**: Read a single issue by its key in the governed Linear workspace. Use it when the user asks about an issue (status, details); do not call it just because a request mentions an issue key
- **Permission**: `linear.read`
- **Parameters**:

  | Name | Type | Required | Description |
  | --- | --- | --- | --- |
  | `issue_key` | string | yes | Issue key, e.g. ENG-123 |

- **Tags**: `linear-read`, `governed-connector`

- **Transport**:

```json
{
  "method": "POST",
  "url_template": "https://api.linear.app/graphql",
  "auth": {
    "placement": "header",
    "name": "Authorization",
    "scheme": "Bearer"
  },
  "args": {
    "issue_key": {
      "in": "body",
      "name": "variables.id"
    }
  }
}
```

#### `linear.list_projects`

- **Description**: List projects in the governed Linear workspace
- **Permission**: `linear.read`
- **Parameters**:

  | Name | Type | Required | Description |
  | --- | --- | --- | --- |
  | `state` | string | no | Filter by project state |
  | `limit` | integer | no | Maximum number of projects to return |

- **Tags**: `linear-read`, `governed-connector`

- **Transport**:

```json
{
  "method": "POST",
  "url_template": "https://api.linear.app/graphql",
  "auth": {
    "placement": "header",
    "name": "Authorization",
    "scheme": "Bearer"
  },
  "args": {
    "state": {
      "in": "body",
      "name": "variables.filter.state"
    },
    "limit": {
      "in": "body",
      "name": "variables.first"
    }
  }
}
```

#### `linear.create_issue`

- **Description**: Create an issue in a team of the governed Linear workspace. Use to track work such as a bug, review, or follow-up task.
- **Permission**: `linear.write`
- **Parameters**:

  | Name | Type | Required | Description |
  | --- | --- | --- | --- |
  | `team_key` | string | yes | Key of the Linear team that owns the issue (the prefix of its issue keys, e.g. KEI for KEI-123) |
  | `title` | string | yes | Issue title |
  | `description` | string | no | Issue body in Markdown |
  | `priority` | string | no | Issue priority |
  | `labels` | string | no | Comma-separated label names |

- **Tags**: `linear-write`, `governed-action`

- **Transport**:

```json
{
  "method": "POST",
  "url_template": "https://api.linear.app/graphql",
  "auth": {
    "placement": "header",
    "name": "Authorization",
    "scheme": "Bearer"
  },
  "args": {
    "team_key": {
      "in": "body",
      "name": "variables.input.teamId"
    },
    "title": {
      "in": "body",
      "name": "variables.input.title"
    },
    "description": {
      "in": "body",
      "name": "variables.input.description"
    },
    "priority": {
      "in": "body",
      "name": "variables.input.priority"
    },
    "labels": {
      "in": "body",
      "name": "variables.input.labelIds"
    }
  }
}
```

#### `linear.create_followup_task`

- **Description**: Create a Linear follow-up task for a CRM lead, by lead id. Call it directly when asked for a follow-up on a lead; the summary must not contain contact details
- **Permission**: `linear.write`
- **Parameters**:

  | Name | Type | Required | Description |
  | --- | --- | --- | --- |
  | `lead_id` | string | yes | Governed CRM lead identifier |
  | `title` | string | yes | Follow-up task title |
  | `summary` | string | yes | Redacted CRM context approved for sharing |
  | `approval_id` | string | no | Approval reference for the write |
  | `idempotency_key` | string | yes | Stable replay key for this task creation |

- **Tags**: `linear-write`, `crm-followup`, `governed-action`, `harness-neutral`

- **Transport**:

```json
{
  "method": "POST",
  "url_template": "https://api.linear.app/graphql",
  "auth": {
    "placement": "header",
    "name": "Authorization",
    "scheme": "Bearer"
  },
  "args": {
    "lead_id": {
      "in": "body",
      "name": "variables.input.teamId"
    },
    "title": {
      "in": "body",
      "name": "variables.input.title"
    },
    "summary": {
      "in": "body",
      "name": "variables.input.description"
    },
    "approval_id": {
      "in": "header",
      "name": "X-Approval-Id"
    },
    "idempotency_key": {
      "in": "body",
      "name": "variables.input.id"
    }
  }
}
```

## Credential pass-through

Kei supplies the credential at run time. After `kei-proxy authorize` returns an allow decision, the harness calls `POST /v1/credentials:release` over the Unix socket
(`KEI_RUNTIME_SOCKET_PATH`, default `/run/kei-proxy/runtime.sock`) to obtain the short-lived access token.

The request identifies the tool by name and carries the session in `X-Kei-Session`. On allow the response carries `credential.access_token` and `credential.expires_at`. See `docs/credential-passthrough.md` in the kei-connector-contracts repo for the full contract.

The agent never holds, asks for, or prints a raw credential. The credential value crosses the socket only at call time and is never written to a config file, environment variable, disk cache, or log line.

## Connect

An owner connects this provider with:

```
kei connectors create --provider linear --workspace W
```

Check status with `kei connectors get <id>` and re-authorize with `kei connectors reconnect <id>`.

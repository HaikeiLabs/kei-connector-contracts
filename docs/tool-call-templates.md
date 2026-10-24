# Tool call templates (HAI-856)

Decision record, 2026-10-24 (owner approved). Builds on ADR-039 (connectors as
data, HAI-719) and the `binding.transport` fields from HAI-723.

## Decision

A **tool call template** is a `kei.tool-definition/v1` entry
(`schemas/tool-definition.v1.schema.json`) read under a stricter profile
(`schemas/tool-call-template.v1.schema.json`). There is no second format: the
entries in `schemas/*.tool-definition.v1.json` are the templates.

- Contracts are the catalog of tool-call **options**. A harness imports only
  the tools for the functions it runs.
- The harness pins a connector-contracts version into its image at build time
  (like `kei-proxy.version`). Nothing is fetched at runtime.
- At startup the harness reads and validates the entries it needs and registers
  tool names and schemas through agentware in **one** sync.
- The same entries drive the call after an allow, so what policy sees is what
  runs.
- Templates hold **placeholders only**. The credential is injected by the
  harness after the allow (ADR-039), at the place named by
  `binding.transport.auth`. A template never holds a secret or a resolved
  credential.
- A template that fails validation is a startup error. The harness must not
  register a partial set.

A runtime-mounted ConfigMap is a possible later addition, not the default.

## The profile

Everything in the tool-definition schema applies. On top of it a template must:

1. Carry a tool schema in `parameters`: a list of named parameters, or a JSON
   Schema object with `type: object`, `additionalProperties: false`, and every
   `required` name present in `properties`. A tool without a schema is rejected.
2. Have `permission`, `description`, a unique `name`, and no unknown fields.
3. If `binding.transport` is present: a valid `method`, a `url_template`, valid
   `auth.placement`, every `args.*.in` in `path|query|body|header`, and an
   `in: path` entry in `args` for every `{placeholder}` in `url_template`.
4. Contain no secret: no token-shaped values (GitHub, Slack, AWS, JWT, `sk-`,
   PEM keys), no literal `Bearer`/`Basic` credential, no `user:pass@` in a URL,
   no literal credential in a query string, and no secret-named key (`token`,
   `secret`, `password`, `api_key`, `credential`, ...) holding a literal value.

`binding.transport` is optional. A template without it describes an action the
harness runs itself (tag `harness-action`), for example `github.create_fix_pr`,
where the server picks the repository, branch and draft status. It still
registers name, schema and permission.

`transport.args.<param>.encoding` (`base64`) is the one addition to the
transport fields: the GitHub contents API takes file content base64 encoded.

## First templates

| Tool | File | Permission | Notes |
| --- | --- | --- | --- |
| `linear.create_issue` | `schemas/linear.tool-definition.v1.json` | `linear.write` | GraphQL `POST` |
| `linear.create_followup_task` | `schemas/linear.tool-definition.v1.json` | `linear.write` | GraphQL `POST` |
| `github.create_fix_pr` | `schemas/github.tool-definition.v1.json` | `github.write` | closed schema: `title`, `body`, `unified_diff` (max 32 KiB); no transport (kei-agents#45) |
| `github.file_write` | `schemas/github.tool-definition.v1.json` | `file.write` | `PUT` contents API; governed file_write (HAI-855) |

`conformance/tool_call_templates_test.go` pins their schemas and permissions
against the hand-written definitions they replace (kei-agents `GITHUB_FIX_ACTION_TOOL`,
the chat harness `file_write_tool` and `config/tools-template.yaml`).

Known gap: the chat harness `agent.py` hand-written `github.get_issue` transport
uses `{repository}`, the contracts entry uses `{repo}`. The harness must follow
the template when it migrates.

## How a harness imports and validates templates (for the agentware loader)

1. Build time: pin the contracts version (a git tag) in the image, for example
   a `kei-connector-contracts.version` file next to `kei-proxy.version`, and
   copy `schemas/` from that tag into the image.
2. Startup: for each function the harness runs, name the tools it needs
   (`github.create_fix_pr`, ...). Read the matching
   `schemas/<connector>.tool-definition.v1.json` files and select entries by
   `name`. A missing name is a startup error.
3. Validate the selected entries with the rules above. Go harnesses call
   `calltemplate.LoadFile(path)` or `calltemplate.Parse(data)`, which return
   every problem in one joined error. Other languages must implement the same
   rules; `calltemplate/calltemplate_test.go` is the reference case list.
4. Register: build the agentware manifest from the same entries (`name`,
   `description`, `parameters`, `permission`) and sync **once**. Any validation
   or sync error aborts startup loudly (non-zero exit, error in the log).
5. After an allow: look up the entry by tool name, fill `url_template`
   placeholders and `args` from the tool arguments, apply `encoding`, then
   inject the released credential at `auth.placement` / `auth.name` /
   `auth.scheme`. Templates never contain it.

The loader itself (agentware, catalog, harness) is out of scope for this change.

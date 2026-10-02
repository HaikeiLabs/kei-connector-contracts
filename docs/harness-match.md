# Harness Match — `kei.harness-match/v1`

This package implements the **kei.harness-match/v1** dialect for matching a
pre-built harness's tool call against a set of policy entries.  It is the
shared, dependency-free reference implementation published in
`kei-connector-contracts` so that every Kei service and third-party connector
evaluates harness rules identically.

The dialect is defined by **ADR-029**; see that document for the design
rationale, bundle format, and provenance chain.

## Semantics

A harness call is described by its kind (`claude_code`, `codex`, `opencode`, or
`custom`), the tool name, and optional payload fields (shell argv, skill name,
path, home, MCP server/tool, agent ID, harness ID).  Policies are evaluated in
order; the first enabled policy whose source **and** destination both match
produces the outcome.

### Outcomes

| Outcome | Meaning |
|---------|---------|
| `permit` | A matching permit policy was found. |
| `deny` | A matching deny policy was found. |
| `unmatched` | No enabled policy matched the call. |
| `not_renderable` | The call cannot be rendered into a policy entry (see reason codes). |

### Reason codes

| Reason | Trigger |
|--------|---------|
| `unsupported_kind` | The harness kind is not in the known set. |
| `path_not_absolute` | The supplied path is relative or not clean. |
| `evaluate_via_pdp` | The `custom` kind must be evaluated by the PDP, not the harness renderer. |
| `human_source` | The only matching source is a human-only pattern (`user:`, `email:`, `group:`, `org:`). |
| `dst_not_renderable` | The destination pattern cannot be evaluated for this call (e.g. `~/` glob with empty home). |
| `bundle_invalid` | (reserved) The policy bundle is malformed. |

### Source matching (`src`)

| Pattern | Matches |
|---------|---------|
| `*` | Any harness or agent. |
| `harness:*` | Any harness. |
| `harness:<kind>` | Harness of the given kind (`claude_code`, `codex`, `opencode`). |
| `harness:<harness_id>` | Harness with the given registered ID (UUID). |
| `agent:*` | Any agent (when the call has an agent ID). |
| `agent:<id>` | Specific agent by ID. |
| `user:`, `email:`, `group:`, `org:` | Human-only patterns — never matched by a harness renderer. |

### Destination matching (`dst`)

| Scheme | Format | Matches |
|--------|--------|---------|
| `shell:` | `shell:<argv-prefix>` | Token-wise prefix of the call's shell argv. `shell:*` matches any shell command. |
| `skill:` | `skill:<name>` | Exact skill name. `skill:*` matches any skill. |
| `path:` | `path:<glob>` | Per-segment glob with `*` (within a segment) and `**` (zero or more segments). `~/` expands to the call's home directory. |
| `mcp:` | `mcp:<server>` or `mcp:<server>.<tool>` | MCP server name, or server + tool joined by `.`. `mcp:*` matches any MCP server. |
| `tool:` | `tool:<name>` | Exact tool registry name. `tool:*` matches any tool. |
| `*` | (bare) | Matches any destination. |

The shell-permit rule (ADR-029 §9.4): a permit on a shell command counts only
when paired with a `shell:` argv-prefix destination.  `shell:*` permits are
rejected (they must use a specific argv prefix).  Deny policies on `shell:*`
still apply.

### Scope

A policy may carry an optional `scope_agent_id`.  When set, the policy matches
only calls whose `AgentID` equals the scoped value.  This is checked before
source or destination matching.

## Fixtures

The file `harnessmatch/testdata/cases.v1.json` contains the shared fixture set
(40+ cases) that every consumer should run to ensure consistent behaviour.
Load them with `harnessmatch.LoadFixtures()`.

## Usage

```go
import "github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"

call := harnessmatch.Call{
    Kind: "claude_code",
    Tool: "claude_code.edit",
    Argv: []string{"git", "commit", "-m", "msg"},
}

policies := []harnessmatch.Policy{
    {ID: "p1", Src: "*", Dst: "shell:git commit -m", Action: "permit", Enabled: true},
}

result := harnessmatch.Evaluate(call, policies)
// result.Outcome == harnessmatch.OutcomePermit
// result.PolicyID == "p1"
```

# Connector policy field contract

This document defines the canonical meanings of plain-text values shared by
connector policy producers, evaluators, and audit consumers. It does not define
a new policy schema or change the existing flat pattern matcher.

## Dimensions

- **Source/provider** is the canonical `Metadata.Provider` value (for example,
  `github` or `google_drive`). It identifies the integration, not the acting
  principal. Resolve it from trusted connector metadata, never from a
  caller-provided alias.
- **Connector instance** is the exact `Metadata.ID` / invocation
  `connector_id`. It distinguishes instances of the same provider.
- **Operation/capability** is the exact invocation `capability` (for example,
  `pull_request.read`). `action` (`read`, `create`, `update`, `delete`, or
  `comment`) is a separate operation classification. Neither is the policy
  effect.
- **Resource type** is the kind of object, such as `pull_request`; it is not
  present as a dedicated field in this contract version. Do not infer it from
  the identifier.
- **Resource identifier** is invocation `resource`, retained verbatim. It is
  provider-defined plain text, not a display name or a normalized path.
- **Principal** is the authenticated invocation `subject` and `agent_id`.
  Provider and connector identity are not principals.
- **Scope** is the authenticated `tenant_id` and `workspace_id`; both bound
  the invocation to its connector metadata. Scope is not part of a resource
  string.
- **Policy effect** in a policy catalog is `permit` or `deny`. The connector
  decision wire value is currently `allow` or `deny` (`DecisionAllow` /
  `DecisionDeny`); do not silently rewrite that existing enum. A catalog must
  document its mapping (`permit` -> `allow`, `deny` -> `deny`) when projecting
  policy effect into this decision field.

All these values are opaque canonical strings: preserve exact bytes, case, and
punctuation. Do not add friendly-name aliases, lowercase, parse/reformat, or
resolve identifiers through display metadata. Provider and capability names
must use the values declared in this module.

## Selector compatibility and matching boundary

This module does not own the policy catalog's pattern matcher. Existing
catalog behavior remains authoritative and must not be changed by interpreting
these strings differently:

- Exact values mean exact equality. The contract does not imply prefix,
  hierarchy, or partial-glob matching.
- A flat legacy destination selector that is tested against several candidate
  strings remains an OR of those candidates. Separate legacy selectors must
  not be assumed to compose as a conjunction.
- Do not reinterpret legacy `action:resource` strings as a new field encoding.
  Preserve existing bare capability/resource selectors and their established
  compatibility behavior.
- A structured conjunction across provider, connector, capability, resource
  type, and resource identifier requires a separately versioned policy
  contract. This contract does not concatenate dimensions into strings.
- Wildcard syntax and principal-selector grammar belong to the policy catalog
  matcher. This contract does not expand partial `*`, introduce hierarchy
  wildcards, or change legacy wildcard behavior.

## Audit projection

The current `AuditRecord` exposes subject (`invoking_subject`), agent
(`agent_id`), scope (`org_id`, `workspace_id`), connector plus capability in
`tool_name` (`connector/<connector_id>/<capability>`), exact resource in
`resources_touched`, and outcome in `decision`. `policy_id` and
`policy_reason` are optional. This audit shape does **not** have dedicated
provider or resource-type fields; do not parse them from `tool_name` or
`resources_touched`.

Example using available contract fields (illustrative values):

```json
{
  "invoking_subject": "user:U-17",
  "agent_id": "agent:ReviewBot",
  "org_id": "org-4",
  "workspace_id": "workspace-9",
  "tool_name": "connector/inst-7/pull_request.read",
  "resources_touched": ["repos/Acme/Kei/pulls/17"],
  "decision": "allow",
  "policy_id": "policy-12"
}
```

The provider (`github`) and resource type (`pull_request`) are intentionally
not represented in that record: adding either to the shared audit wire shape
requires an explicit additive contract decision and producer support. The
sample retains the identifier's original case and punctuation. A catalog's
policy effect `permit` corresponds to the recorded decision `allow`; a deny
is recorded as `deny`.

## Ownership boundary

This shared module owns connector metadata, invocation, decision, and audit
wire contracts. The policy catalog owns selector parsing, matching, priority,
and persisted policy effect. It must keep server-side documentation and tests
for its legacy matcher semantics; this document is not a matcher
implementation or a migration instruction.

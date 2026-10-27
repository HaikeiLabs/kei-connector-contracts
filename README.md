# kei-connector-contracts

The versioned connector contract shared by Kei services: connector metadata,
capabilities, invocation envelopes, decisions, audit shapes, and provider
definitions.

Extracted from `kei` under
[ADR-017](https://github.com/HaikeiLabs/kei/blob/main/docs/adr/017-service-repository-boundaries-and-naming.md),
which names this package "the first extraction seam". Previously it lived at
`contracts/connectors` in `kei` and was consumed through a filesystem `replace`
directive, so it could not be used by any repository outside that checkout.

## Invariants

This module is **metadata only**. Per
[ADR-011](https://github.com/HaikeiLabs/kei/blob/main/docs/adr/011-tenant-side-proxy-data-plane.md)
it must never carry provider payloads, provider results, customer content, or
credential values — `credential_ref` is an opaque pointer, never key material.
The one request payload it can describe is an http_api `http.post` or
`http.patch` body (`HTTPInvocation.Body`: a JSON object, `application/json`
only, capped at the connector's `request_max_bytes`). That request stays in the
tenant runtime: the control plane authorizes the call on the capability's
method and the resource path (`AuthorizeCall`) and never needs the body.

It is also deliberately **stdlib-only**. Every service imports it, so a
dependency added here propagates to all of them. CI fails if `go.mod` gains a
`require` block.

## Authorization and approvals

Connector contracts describe the operation, capability, resource, and
metadata needed for runtime execution. They do not carry a one-time human
approval for an invocation. Access is granted to users and groups through
workspace permissions before execution; the tenant-side runtime executes
operations covered by the workspace-bound identity. Contract PR #6 removes
the exported `Invocation.ApprovalID`, `PolicyAttributes.DestructiveEnabled`,
and `providers.DestructiveAllowed` fields. See the proposed
[ADR-027: Approvals grant access; there are no per-call approvals](https://github.com/HaikeiLabs/kei/pull/671)
for the decision and migration inventory; runtime consumers still need to
remove their remaining approval transport fields.

The contract is frozen: see `contract_freeze_test.go` and
`providers/contract_freeze_test.go`. Breaking changes require a new contract or
envelope version rather than an edit in place.

## Runtime policy bundle health (v1)

The shared `policybundlehealth` package and
`schemas/runtime-policy-bundle-health.v1.schema.json` own the runtime heartbeat
report and catalog read projection. Canonical examples live under
`schemas/examples/runtime-policy-bundle-health/`. The contract is metadata
only; the policy bundle itself remains local to the tenant runtime, and only
its wire types are shared (see
[Shared policy bundle](#shared-policy-bundle-hai-609-unreleased)). See
`Health.Validate` and `Health.Decode` for the closed state/reason and
field-invariant validation.

## Authorize request (v1)

The `authorizerequest` package and `schemas/authorize-request.v1.schema.json`
own the `kei.authorize-request/v1` request a harness sends to the tenant-local
`kei-proxy` daemon. The optional `caller` object (`harness`, `origin`
`interactive` or `automation`, and an opaque `session_ref`) lets the audit
trail say who called. A request without it stays valid and its caller is
unknown. Canonical examples live under `schemas/examples/authorize-request/`.
The contract is metadata only: `session_ref` is a pointer and must never carry
message content.

## Consumers

- `kei-connector-runtime` — tenant-side evaluation and provider execution
- `kei-policy-catalog` — connector metadata and registrations
- `kei-cli` — reads and inspects a policy bundle
- `kei` — integration tests

## Packages

Dependencies point one way: `contract` ← `capability`, `envelope`, `setup`;
`envelope` ← `governance` ← `audit` ← `invoke`; `providers` imports
`contract`; `bundle` imports `contract`, `harnessmatch`, `policy` and
`policybundlehealth`.

| Import path | Contents |
| --- | --- |
| `github.com/HaikeiLabs/kei-connector-contracts/contract` | Provider, Status, Action, Capability, and the capability definitions; `Metadata`, `Invocation`, the `http_api` types, and `GrafanaConfig`; `Validate`, `ValidateCall`, `ValidateInvocation`, `ValidateCredentialRef`; credential source and account models (`accountmodel.go`) |
| `.../setup` | the setup schema (`SetupSchema`, `SetupSchemas`, `SetupSchemaFor`), `ValidateConfig`, and `ValidateMetadata` (the complete connector check: `Metadata.Validate`, then config) |
| `.../access` | per-provider access modes (`ProviderAccess`, `AccessMode`, `Access`, `AccessFor`, `CredentialReference`); metadata only |
| `.../capability` | `LookupCapability`, `CapabilityFor`, `CapabilitiesForProvider`, `Providers` |
| `.../envelope` | `Envelope`, `EnvelopeVersion1`, `MintedByControlPlane` |
| `.../governance` | `Decision`, `PolicyDecision`, `Decide` |
| `.../audit` | `AuditRecord`, `BuildAuditRecord` |
| `.../invoke` | `Client`, `NewClient`, `InvokeOutcome` (envelope + decision + audit record) |
| `.../providers` | provider payloads, results, and read-only clients |
| `.../canonical` | `Resource` (the canonical resource URI of an MCP server, ADR-036 §3.4), `ErrInvalidResource` |
| `.../canonical/canonicaltest` | `ResourceCases`, the shared golden cases for `canonical.Resource`, also published as `canonical/canonicaltest/resource.json` |
| `.../policy` | `Principal`/`ParsePrincipal`, `Target`/`ParseTarget`, `MatchPrincipal`, `Rule`/`Request`/`Decide`/`Precedence`, `Render` (a thin entry point over `harnessmatch`), and the conformance vectors in `policy/testdata/vectors/` |
| `.../bundle` | the policy bundle wire types (`Payload`, `Audience`, `Identity`, `Harness`, `AgentModelProfile`, `InstallationSubject`, `Verified`), the bundle validators (`Validate`, `ValidateV2`, `RejectDuplicateJSON`), the decisions a bundle yields (`Evaluate`, `EvaluateConnector`, `EvaluateSession`, `EvaluateV2HarnessTool`, `EvaluateV2Plan`, `EvaluateV2Model`, `EvaluateHarness`) and the validator hooks a caller installs (`SetSchemaValidator`, `SetV2SchemaValidator`, `SetPlanValidator`) |

Account models are a file in `contract`, not a package: `Metadata.AccountModel`
is a contract field and `Metadata.Validate` enforces the account-model rules,
so the code must live with the `Metadata` type. The setup schema is its own
package; because `setup` imports `contract`, `Metadata.Validate` does not
check config fields. Validate a connector with `setup.ValidateMetadata`,
which runs `Metadata.Validate` and then the config rules.
`schemas/connector-setup.v1.json` stays at the repository root.

The root import path `github.com/HaikeiLabs/kei-connector-contracts` has no Go
package from v0.2.0 on. v0.1.0 consumers change `connectors.X` to the package
above that now holds `X`.

## Providers as a harness library

The `providers` package is a stable, harness-consumable library
(ADR-036 section 9.2). A harness — the agentware SDK, a chat harness, or a
custom SDK agent — imports it to read and write provider data, supplying its
own credential and HTTP client. Each provider constructor takes a credential
(a `CredentialResolver`) and an `*http.Client`, and returns a client bound to
the provider's fixed base endpoint. The package never imports
`kei-connector-runtime`; `providers/boundary_test.go`
(`TestProvidersImportBoundary`) enforces that, and breaking changes follow the
module's semantic versioning.

```go
import (
	"context"
	"errors"
	"net/http"

	"github.com/HaikeiLabs/kei-connector-contracts/providers"
)

// resolver is the harness's credential seam: it resolves an opaque kei://
// reference to a token and must never log the secret.
type resolver struct{}

func (resolver) Resolve(ctx context.Context, tenantID, workspaceID, subject, ref string) (string, error) {
	return "", errors.New("credential store not wired")
}

func main() {
	// A credential and an *http.Client in; no runtime import.
	github := providers.NewGitHubHTTP(&http.Client{}, resolver{})
	_ = github
}
```

## Credential source, account models, and setup (v0.2.0)

- `Metadata.CredentialSource` is `oauth` or `opaque_ref`. Empty means
  `opaque_ref`, so v0.1.0-shaped metadata stays valid.
- `Metadata.AccountModel` says whose account an OAuth-style connector acts as:
  `per_user` (the invoking user's token; no subject), `shared` (one account
  for the connector; subject `connector:<id>`), or `domain_delegation`
  (Google Workspace service account over `opaque_ref`, with an
  `impersonate_email` config field). `AccountModelsFor` lists what each
  provider allows. `CredentialSubject` names whose credential an invocation
  uses, for audit.
- `SetupSchemaFor` / `SetupSchemas` describe each buildable connector's setup
  fields (type, pattern, which are secret, which account models they apply
  to), and `ValidateConfig` checks `Metadata.Config`. The JSON export for the
  web UI and CLI is `schemas/connector-setup.v1.json`; a test keeps it in
  sync with the Go data.
- New providers `gmail` and `tito` (read-only), moved from
  `kei-policy-catalog`'s in-tree fork, along with the legacy `kei-oauth`
  reference rejection.

Contract: `docs/connector-execution-contract.md` in `kei-connector-runtime`.

## Discord and Grafana (HAI-309, unreleased)

Two read-only, shared-secret providers. Both use `credential_source`
`opaque_ref` only (no OAuth, no account models), and every capability is
`read`. The contract carries the payloads (`providers/discord.go`,
`providers/grafana.go`, each with `Validate`) and the resource grammar
(`providers.ValidateDiscordResource`, `providers.ValidateGrafanaResource`);
there is no backend seam, and the runtime executes the reads.

| Provider | Capability | Resource | Input |
| --- | --- | --- | --- |
| `discord` | `guild.list` | `guilds` | |
| `discord` | `guild.read` | `guilds/<snowflake>` | |
| `discord` | `channel.list` | `guilds/<snowflake>/channels` | |
| `discord` | `channel.read` | `channels/<snowflake>` | |
| `discord` | `thread.list` | `channels/<snowflake>/threads` | |
| `discord` | `message.list` | `channels/<snowflake>/messages` | `limit` 1-100, `before` or `after` snowflake |
| `discord` | `message.read` | `channels/<snowflake>/messages/<snowflake>` | |
| `discord` | `role.list` | `guilds/<snowflake>/roles` | |
| `grafana` | `folder.list` | `folders` | |
| `grafana` | `folder.read` | `folders/<uid>` | |
| `grafana` | `dashboard.search` | `search` | `query`, `tags`, `limit` |
| `grafana` | `dashboard.read` | `dashboards/<uid>` | |
| `grafana` | `datasource.list` | `datasources` | |
| `grafana` | `datasource.read` | `datasources/<uid>` | |
| `grafana` | `alert_rule.list` | `alert_rules` | |
| `grafana` | `annotation.list` | `annotations` | `from`/`to` (at most 31 days), `dashboard_uid`, `tags`, `limit` |
| `grafana` | `datasource.query` | `datasources/<uid>/query` | `from`/`to` (at most 7 days), `max_rows` (at most 1000), `queries` (1-10) |

Discord has no `member.*` capability: member lists are personal data.

Setup: `discord` takes the secret `bot_token`. `grafana` takes `base_url`,
stored in the new `Metadata.Grafana` (`GrafanaConfig{BaseURL}`, setup location
`grafana`), and the secret `service_account_token`. `Metadata.Validate`
requires `Metadata.Grafana` on a grafana connector and refuses it elsewhere.
The base URL is an https origin with no path, query, fragment, or userinfo.
Private and internal hosts are allowed, because Grafana is often self-hosted.
Link-local addresses (169.254.0.0/16, fe80::/10), `fd00:ec2::254`, unspecified
addresses, numeric-encoded IPv4 hosts, and cloud-metadata hostnames
(`metadata`, `metadata.*`, `instance-data`, `instance-data.*`) are always
refused. The runtime must still check the resolved address when it dials.

## Skills repository (unreleased)

A dedicated, read-only provider for a workspace's agent skill source: a git
repository holding one folder per skill, each with a `SKILL.md`. It is its own
provider rather than a `github` capability, so a skill source never implies
access to a repository's issues or pull requests.

| Provider | Capability | Resource | Result |
| --- | --- | --- | --- |
| `skills_repo` | `skills.list` | `skills` | `SkillsRepoListing`: `ref`, the resolved `commit`, and per skill `name`, `version`, `description`, `harnesses`, `path` |
| `skills_repo` | `skills.fetch` | `skills/<name>` | `SkillsRepoArchive`: the skill entry, `commit`, `format` (`tar+gzip`), `digest` (`sha256:<hex>`), `size`, and the `archive` bytes |

- A skill name is one segment of letters, digits, `-` and `_` (at most 64), so
  a fetch cannot leave the skills folder. A folder whose `SKILL.md`
  frontmatter names another skill is left out of the listing and refused by
  `skills.fetch`.
- `version` and `harnesses` are read from the frontmatter, top level or under
  `metadata:`. Values that are not plain tokens are dropped.
- The archive is inline in the invoke result: the result shape has no blob
  reference, and a reference would need a second, ungoverned download. It is
  deterministic for a commit (path order, fixed timestamps, modes 0644/0755),
  bounded (500 files, 2 MiB per file, 8 MiB in total), and refuses symlinks
  and submodules.
- Setup: `repo_url` (`https://github.com/<owner>/<repo>`; GitHub is the only
  host today), `ref` (default `main`), and `subdirectory` (default `skills/`,
  `/` for the root). Auth is the invoking user's GitHub OAuth token
  (`per_user`) or a shared-secret `token` over `opaque_ref`. The new optional
  `SetupField.CredentialSources` limits the token field to the `opaque_ref`
  path.
- `providers.GitHubSkillsRepo` is the GitHub backend behind the
  `SkillsRepoBackend` seam. It resolves the ref to a commit and reads trees
  and blobs by object id, so one call reads one snapshot.

`skills_repo` is in the access manifest with an `api` mode that reads
`skills.list` and `skills.fetch`; kei-connector-runtime executes it (HAI-431).

## Provider access modes (HAI-319, unreleased)

`access.Access()` lists how each provider can be reached. The JSON export is
`schemas/connector-access.v1.json` and its JSON Schema is
`schemas/connector-access.v1.schema.json` (`kei.connector-access/v1`). This is
**metadata only and grants nothing**: ABAC policy still decides every call,
and the owner decides which modes each connector actually uses.

Each provider has an `access` list. Each mode has:

- **`kind`:** one of the following.
  - `api`: the runtime adapter, run by `kei-proxy connector invoke`.
  - `cli`: a provider CLI, given by `binary`, `min_version` and an `install` hint.
  - `mcp`: a remote MCP server, given by `transport` (`streamable_http` or `sse`) and an https `url` (no userinfo, query, or fragment). `resource` is the canonical resource URI (defaults to `url`); `issuer` is the optional expected authorization server issuer.
  - `webhook`: given by `direction` (`inbound` or `outbound`) and the `secret_field` of the signing secret.
- **`credentials`** (cli and mcp only): pairs of `env` and `field`. The env var
  holds the reference `kei://connectors/<connector_id>/<field>`, which
  `kei-proxy run --` (HAI-305) resolves into the child process only (cli), or
  the runtime's MCP route attaches to upstream requests (mcp).
- **`setup` and `verify`:** ordered steps.
- **`reads` and `writes`:** the provider's read and write capabilities. Writes
  are made by agent action tools (ADR-028).

Today only `api` is populated. It is set for the providers with a runtime
adapter: gmail, google_drive, linear, github, tito, notion, discord, grafana,
crm, http_api and skills_repo. Only crm, linear and http_api list writes, because the
runtime executes them: every crm write, Linear `issue.create` and
`comment.create`, and http_api `http.post` and `http.patch`. No provider has a
cli, mcp or webhook mode yet.

## Canonical resource URIs (HAI-566, unreleased)

`canonical.Resource` turns an MCP server URL into its canonical resource URI.
It is the one implementation that the OIDC bridge, the policy catalog and the
runtime share: a token's resource, a stored connector's `mcp_url` and the
route's audience check are compared through it, so two copies that differed
would make every connection read as changed.

- It lowercases the scheme and host, and drops the default port `:443`. Any
  other port is kept as written, so `:0443` stays.
- It keeps IPv6 literals as written, apart from case.
- It strips every trailing slash unless the path is `/`, and an empty path
  becomes `/`.
- It rejects any URL that is not absolute `https`, and any URL with user
  info, a query (even an empty one) or a fragment.
- It checks shape only. It does not resolve hosts or apply SSRF rules.
- Its error never contains the submitted URL, which may carry credentials.

Test against `canonicaltest.ResourceCases()`, or `resource.json` from another
language, rather than a local table. The cases pin current results,
including the odd ones marked with a `note`.

## Shared policy model (HAI-606, unreleased)

`policy` is the one implementation of policy matching and decisions, shared by
the catalog, the runtime and the CLI so they cannot drift apart. It replaces
the matcher copies those services each carried: principals (who), targets
(what), decisions, and rendering.

- `ParsePrincipal`/`ParseTarget` parse `<scheme>:<value>` patterns; `*` is
  any, a bare name is a group for principals and an error for targets, and an
  unknown scheme or an empty value is an error.
- `MatchPrincipal` matches a principal pattern against user/email/groups with
  the same semantics as the catalog's `MatchesPattern`. `MatchTarget` matches
  a target pattern against a request target: exact, `*`, `<scheme>:*`, a shell
  token prefix (`shell:git` matches `shell:git push`, and a trailing slash on
  the last token is a prefix within it), or a connector dot-segment pattern
  with `*` segments (`connector:github.*`). `path:`, `skill:`, `tool:` and
  `capability:` keep exact matching.
- `Decide` orders the matching enabled rules by `Precedence` (priority DESC,
  deny before permit, ID ASC); the first decides, or the outcome is
  `no_match`. The runtime calls it when it authorizes a tool call.
- `Render` turns the rules that apply to one subject into the native
  allow/deny entries of a harness kind (for example the lists `kei harness
  sync` writes for a coding harness). The harness is the rendering context, not
  a principal: matching rules become entries, non-matching ones are withheld,
  an unknown kind is unrenderable, and a permit that overlaps a deny that is
  not rendered is withheld (fail closed).

Conformance vectors live in `policy/testdata/vectors/` (`parse-principal`,
`parse-target`, `match-principal`, `match-target`, `decide`, `render`).

## Shared policy bundle (HAI-609, unreleased)

`bundle` is the shared implementation of the policy bundle contract: the
document the policy catalog issues to an installation, the checks it has to
pass, and the decisions every service derives from it. It is the one copy the
CLI (which reads a bundle) and the tenant runtime (which enforces one) share,
so the two cannot drift apart.

- `Payload` is the bundle envelope: `Audience`, the `Identities` table, the
  registered `Harnesses`, each agent's `AgentModelProfile`, the installation
  `Subject`, the refresh window and the raw `PolicySet`. `Verified` is a
  validated bundle together with its digest.
- `Validate` and `ValidateV2` check the envelope: identity, audience, version
  rollback, the issued-at clock-skew bound (`MaxIssuedAtClockSkew`), the
  validity window, the subject-resolution mode and the policy-set dialect.
  `RejectDuplicateJSON` rejects bundle bytes that parse two ways.
- `Evaluate`, `EvaluateConnector`, `EvaluateSession` and the v2 evaluators
  (`EvaluateV2HarnessTool`, `EvaluateV2Plan`, `EvaluateV2Model`,
  `EvaluateHarness`) decide a call against a validated bundle. `Decision`
  carries the outcome, the deciding policy and a reason code.
- `SetSchemaValidator`, `SetV2SchemaValidator` and `SetPlanValidator` install
  the checks this dependency-free module cannot make itself: the bundle's JSON
  Schema documents and a v2 tool's operation plan both need a JSON Schema
  implementation, so the enforcing service compiles them and passes them in. A
  validator that is not installed is not applied.

The bundle's persisted state — the copy an installation enforces — is not part
of this package: it is runtime state rather than a wire contract, and it is
platform-specific. The same goes for the bundle's canonical JSON Schema
documents and example payloads, which stay with the services that compile and
issue them.

## Versioning

Tagged with semver. Consumers require a released version; no `replace`
directives.

## Building a connector

The contract is public and Apache 2.0 licensed so that anyone can implement a
connector against it without going through Haikei. Model the provider's shapes
in `providers/`, declare capabilities, and the runtime enforces policy over them.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the rules that keep the contract
safe to depend on — chiefly that it stays metadata-only, stdlib-only, and that
breaking changes take a new version rather than an edit in place.

## Licence

[Apache License 2.0](LICENSE). Every Go file carries the header; CI enforces it
via `scripts/license-header.sh --check`, and running the script without
arguments adds it to new files.

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
only; the policy bundle remains local to the runtime. See `Health.Validate` and
`Health.Decode` for the closed state/reason and field-invariant validation.

## Consumers

- `kei-connector-runtime` — tenant-side evaluation and provider execution
- `kei-policy-catalog` — connector metadata and registrations
- `kei` — integration tests

## Packages

Dependencies point one way: `contract` ← `capability`, `envelope`, `setup`;
`envelope` ← `governance` ← `audit` ← `invoke`; `providers` imports
`contract`.

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
  - `mcp`: an MCP server, given by `package` and `transport` (`stdio` with `command`, or `streamable_http` with an https `url`).
  - `webhook`: given by `direction` (`inbound` or `outbound`) and the `secret_field` of the signing secret.
- **`credentials`** (cli and mcp only): pairs of `env` and `field`. The env var
  holds the reference `kei://connectors/<connector_id>/<field>`, which
  `kei-proxy run --` (HAI-305) or the MCP launcher resolves into the child
  process only.
- **`setup` and `verify`:** ordered steps.
- **`reads` and `writes`:** the provider's read and write capabilities. Writes
  are made by agent action tools (ADR-028).

Today only `api` is populated. It is set for the providers with a runtime
adapter: gmail, google_drive, linear, github, tito, notion, discord, grafana,
crm and http_api. Only crm lists writes, because the runtime executes crm
writes. No provider has a cli, mcp or webhook mode yet.

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

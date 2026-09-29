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
| `github.com/HaikeiLabs/kei-connector-contracts/contract` | Provider, Status, Action, Capability, and the capability definitions; `Metadata`, `Invocation`, and the `http_api` types; `Validate`, `ValidateCall`, `ValidateInvocation`, `ValidateCredentialRef`; credential source and account models (`accountmodel.go`) |
| `.../setup` | the setup schema (`SetupSchema`, `SetupSchemas`, `SetupSchemaFor`), `ValidateConfig`, and `ValidateMetadata` (the complete connector check: `Metadata.Validate`, then config) |
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

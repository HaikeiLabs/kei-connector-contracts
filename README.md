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

The contract is frozen: see `contract_freeze_test.go` and
`providers/contract_freeze_test.go`. Breaking changes require a new contract or
envelope version rather than an edit in place.

## Consumers

- `kei-connector-runtime` — tenant-side evaluation and provider execution
- `kei-policy-catalog` — connector metadata and registrations
- `kei` — integration tests

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

# Contributing

This repository holds the connector contract: the types, capabilities,
invocation envelopes, decisions and audit shapes that Kei services and
third-party connectors share.

It is public so that anyone can build a connector against it without going
through Haikei. Contributions are welcome, particularly new provider contracts.

## The contract is frozen

`contract_freeze_test.go` and `providers/contract_freeze_test.go` pin the
existing surface. That is deliberate: every Kei service and every third-party
connector compiles against these shapes, so a change here is a change to a
published API.

- **Additive changes** — a new provider, a new capability, a new optional field
  — are the normal case.
- **Breaking changes** require a new contract or envelope version. Do not edit
  an existing shape in place, and do not modify a freeze test to make a change
  pass. If a freeze test fails, that is the test doing its job.

## What must never enter this module

The contract carries **metadata only**. These rules are not style preferences;
they are the boundary the platform is built on, and CI enforces the last one.

- No credential values. `credential_ref` is an opaque pointer to a location in a
  secret backend — never an API key, token, password or connection string.
- No provider payloads or results. Customer content stays in the tenant runtime
  and must never be modelled here.
- No database, handler, or provider-execution dependencies.
- **No third-party dependencies at all.** This module is stdlib-only. Every
  service imports it, so a dependency added here propagates to all of them. CI
  fails the build if `go.mod` gains a `require` block.

`governance_test.go` and `read_only_test.go` assert several of these directly,
including that values which look like secrets are rejected rather than stored.

## Adding a provider

1. Add `providers/<name>.go` with the provider's types and capability
   definitions. Model the shapes the provider returns; do not implement calls to
   it.
2. Declare capabilities read-only unless writes are genuinely required —
   `read_only_test.go` pins that default deliberately.
3. Add tests alongside, including a freeze test entry for the new surface.
4. Run `go build ./... && go vet ./... && gofmt -l . && go test ./...`.

## Pull requests

Every path is covered by `CODEOWNERS`, so a contract owner reviews each change.
Explain in the description what consumers must do differently, if anything —
"no consumer action required" is a useful thing to state explicitly.

CI runs build, vet, gofmt, the full test suite, and the dependency check on
`ubuntu-latest`. All must pass.

## Licensing and ownership

Contributions are accepted under the [Apache License 2.0](LICENSE), which
includes a patent grant. You retain copyright in what you contribute; you are
licensing it, not assigning it.

Haikei Labs does not assume ownership of, or support obligations for,
connectors built against this contract in other repositories. Publishing a
connector here is not required in order to use Kei.

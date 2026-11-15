# Releasing generated contracts

1. Bump the same version in `gen/python/pyproject.toml`, `gen/ts/package.json`
   (and its lockfile), and `gen/rust/Cargo.toml` (and its lockfile if needed).
   Run `python3 scripts/check-package-versions.py`, then merge the changes to
   `main`.
2. Tag the merged commit `vX.Y.Z` and push the tag. The publish workflow checks
   that the tag is on main and matches all three package versions. To retry,
   run `publish.yml` with workflow dispatch and `version` set to `X.Y.Z`;
   the release tag must already exist.

The `prepare` job validates the release. `npm` tests, builds and publishes
`@haikeilabs/connector-contracts` with provenance. `pypi` builds and publishes
`haikeilabs-connector-contracts`. `crates` publishes the Rust crate
`haikeilabs-connector-contracts`. `go-tag` creates `gen/go/vX.Y.Z` on the
same commit, making `github.com/HaikeiLabs/kei-connector-contracts/gen/go@vX.Y.Z`
available as a Go module. Registry publishes are independent; a retry may need
individual already-published packages to be handled separately.

## One-time owner setup

Configure the GitHub `release` environment and register trusted publishers for
repository `HaikeiLabs/kei-connector-contracts`, workflow `publish.yml`, and
environment `release` on each registry:

- PyPI: create a pending publisher for `haikeilabs-connector-contracts`.
- npm: configure `@haikeilabs/connector-contracts`.
- crates.io: configure `haikeilabs-connector-contracts`.

npm and crates.io may require the package to exist before trusted publishing can
be configured. Verify on the registry's trusted-publishing page; the first
publish may need to be manual or use a short-lived token. The workflow itself
uses OIDC and stores no registry token.

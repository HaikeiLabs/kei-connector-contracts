# Kei policy/manifest conformance fixtures

These deterministic fixtures are drafted from the active HAI-280 policy-set v2 compiler/schema, HAI-392 manifest-v3 decoder/types, HAI-390 preflight route interface, ADR-035, and HAI-398 Discord capability/resource contract. They are intentionally confined to test fixtures; they do not define or change production implementations.

- `manifest-v3.json`: connector and harness-executor route union.
- `policy-set-v2.json`: exact tool-bound checks, capability deny, parent/resource selectors, and direct harness-native grant.
- `policy-bundle-v2.json`: fresh installation-audience wrapper around the policy set.
- `policy-bundle-v1-compat.json` and `manifest-v2-compat.json`: legacy migration inputs.

`go test ./conformance` exercises the shared cases. Bundle v2 is currently not consumed by the active HAI-390 runtime branch: its validator accepts `kei.policy-bundle/v1` and `kei.match/v1` only (see that branch's `docs/tool-call-preflight-integration.md`). Therefore this fixture must not be represented as an end-to-end runtime/Discord regression until the v2 PDP integration lands. HAI-392 SDK and HAI-398 owners should review the fixture paths and the Discord capability/resource choices before these examples are considered frozen.

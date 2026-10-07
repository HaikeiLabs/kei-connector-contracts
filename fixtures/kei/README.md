# Kei policy/manifest conformance fixtures

These deterministic fixtures are drafted from the active HAI-280 policy-set v2 compiler/schema, HAI-392 manifest-v3 decoder/types, HAI-390 preflight route interface, ADR-035, and HAI-398 Discord capability/resource contract. They are intentionally confined to test fixtures; they do not define or change production implementations.

- `manifest-v3.json`: connector and harness-executor route union. Connector routes carry both explicit `agent_id` and `connector_id`; harness-executor routes omit connector-only fields.
- `policy-set-v2.json`: exact tool-bound checks, capability deny, parent/resource selectors, and direct harness-native grant.
- `policy-bundle-v2.json`: fresh installation-audience wrapper around the policy set.
- `policy-bundle-v1-compat.json` and `manifest-v2-compat.json`: legacy migration inputs.

`go test ./conformance` exercises the shared cases, including missing-agent and mismatched-route-identity rejection. Connector routing is authoritative: a local harness handler cannot substitute for a connector-bound route. HAI-392 SDK #168, runtime #130, and catalog #184 have merged; v3 routes now require both `agent_id` and `connector_id` and exact catalog binding validation.

Policy-set v2 and policy-bundle v1 compatibility remain distinct fixtures. Bundle v2 is still not served by the active HAI-390 runtime branch: its validator accepts `kei.policy-bundle/v1` and `kei.match/v1` only (see that branch's `docs/tool-call-preflight-integration.md`). The Discord end-to-end regression is not complete. Runtime preflight is blocked on the owner decision whether to provide a freshness-bounded local subject snapshot with invalidation or exclude subject-dependent v2 calls. Do not add a per-call online subject resolution, substitute caller identity, or claim Discord E2E completion before the runtime gate is resolved and exercised.

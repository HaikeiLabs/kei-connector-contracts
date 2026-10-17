# Credential pass-through over the Unix socket (HAI-722)

This contract describes how a harness obtains a connector credential from the
`kei-proxy` runtime **at call time, after an allow decision**, over the
owner-only local Unix socket. It is the owner decision of 2026-10-17 on HAI-722 /
ADR-039: connector credentials live in the runtime's **session memory** (never
a harness config file), the harness pulls one over the socket when it makes its
provider call, and the runtime **flushes** them when the session ends.

The trust boundary is the owner-only Unix socket
(`KEI_RUNTIME_SOCKET_PATH`, default `/run/kei-proxy/runtime.sock`), the same
listener that serves `/v1/authorize`. The credential value is the one place key
material leaves the daemon: it crosses that socket and nothing else. It is never
written to a harness config file, an environment variable, a disk cache, or a
log line.

## The route

`POST /v1/credentials:release` — a **sibling** of `/v1/authorize`, not a change
to it. It reuses the socket's existing checks (runtime token in `Authorization`,
session in the `X-Kei-Session` header) and the same identifying shape as
`/v1/authorize`. The route is fail-closed: it re-runs the same decision as
`/v1/authorize` for the request and releases a credential **only when that
decision is allow**. It does not trust a client-supplied "I was already
allowed" claim, and it never falls back to a config token, a cached value, or
another session's credential.

- Request schema: `kei.credential-release-request/v1`
- Response schema: `kei.credential-release-response/v1`
- Served only when session credential release is enabled in the daemon; without
  it the route is not served (404), like the other credentialed routes.

## Request

The request carries the session id (header), and exactly one of `tool` or
`connector{connector_id, capability}`, plus the `resource` the provider call
targets — the same fields `/v1/authorize` uses to decide.

| Field | Where | Meaning |
| --- | --- | --- |
| `X-Kei-Session` | header | The named session that holds the credential. |
| `schema` | body | `kei.credential-release-request/v1`. |
| `tool` | body | The harness tool name. Exactly one of `tool` / `connector`. |
| `connector.connector_id` | body | The connector, when the call is connector-scoped. |
| `connector.capability` | body | The capability on that connector. |
| `resource` | body | The resource the provider call targets. |
| `span_id`, `parent_span`, `delegation_depth`, `framework`, `agent_id`, `tool_args_digest` | body | Optional audit context, as in `/v1/authorize`. |

```json
{
  "schema": "kei.credential-release-request/v1",
  "connector": { "connector_id": "github", "capability": "read" },
  "resource": "repo:owner/name",
  "span_id": "0af76b1f1c1c1c1c"
}
```

## Response

### On allow — the credential is released

HTTP 200, `Cache-Control: no-store`. The `credential` object carries the
connector's access token **value** (not an opaque reference — there is no second
exchange or handle store) and its `expires_at`. The value is short-lived and
bounded by the bridge-issued token's lifetime; the runtime does not extend it.

```json
{
  "schema": "kei.credential-release-response/v1",
  "decision": "allow",
  "credential": {
    "connector_id": "github",
    "issuer": "https://auth.example.com",
    "resource": "repo:owner/name",
    "account_model": "per_user",
    "access_token": "<access-token>",
    "expires_at": "2026-10-17T10:00:00Z"
  },
  "policy_id": "pol_01J...",
  "policy_bundle_version": 42,
  "policy_bundle_digest": "sha256:0123...",
  "policy_revision": 7,
  "bundle_state": "ready"
}
```

### On deny / missing / expired credential

HTTP 200 with `decision: "deny"`, a typed `reason_code`, and **no** `credential`
field. There is never a fallback: a missing, expired, or rejected credential is
a deny, not a substitute value.

| `reason_code` | Meaning |
| --- | --- |
| `session_not_found` | The named session does not exist or has ended (already flushed). |
| `no_session` | The session holds no linked credential (never loaded, or the load failed). |
| `credential_not_connected` | The connector is not connected for this user. |
| `credential_reconnect_required` | The connector must be reconnected before it can be used. |
| `credential_max_lifetime_exceeded` | The credential's maximum lifetime was exceeded. |
| `credential_audience_mismatch` | The held credential's audience does not cover the requested resource. |
| `credential_provider_unavailable` | The bridge/provider could not be reached. |
| `credential_provider_unauthorized` | The provider rejected the re-fetched credential. |
| *(a `/v1/authorize` policy reason)* | The re-decision itself denied (e.g. `policy_bundle_missing`, `agent_mismatch`). |

```json
{
  "schema": "kei.credential-release-response/v1",
  "decision": "deny",
  "reason_code": "credential_not_connected",
  "reason": "the connector is not connected for this user",
  "bundle_state": "ready"
}
```

Transport-level failures (missing/invalid runtime token, runtime not
configured) are HTTP errors, as on the other socket routes.

## Refresh during a long session

The session loads its credentials once, at session start (one bridge call). The
release route serves the held credential with no new bridge call. If the
provider rejects a released credential (upstream 401), the harness calls the
route again; the runtime re-fetches that **one** connector credential through
the bridge and releases the new one — exactly one re-fetch, one retry. If the
re-fetch or the retry fails, the route denies with
`credential_provider_unauthorized` and the session must reconnect. There is no
further re-fetch. This is the same single-retry behaviour the daemon already
uses for in-daemon credential use.

## Revocation

When the installation is revoked, or a bundle refresh revokes the session, the
runtime flushes the session's credentials (the existing revocation path). After
that, the release route denies with `no_session`. A single connector credential
is also forgotten when the upstream rejects it (the existing drop path); the
next release re-resolves it once.

## Session-end flush

`DELETE /v1/session` — or the idle-TTL sweep, or daemon shutdown — flushes the
session's credential store: every held token is zeroed and dropped. After the
flush, the release route denies with `no_session` (or `session_not_found` for a
named session that no longer exists). Credentials exist only in the runtime's
session memory for the life of the session.

## Audit

The route emits one audit event per call: the decision (allow/deny), the
`reason_code`, a `credential_released` boolean (true only on a successful
release), the `connector_id`, `resource`, `session_id`, span context, and the
policy metadata (`policy_id`, bundle version/digest/revision). It **never**
records the access token value, the request payload, or any secret. The token
crosses the socket and is not written to the audit log. This extends the
existing decision audit event with a `credential_released` field; the decision
and the release are both recorded, the material is not.

## agentware SDK setup

The agentware SDK already talks to the kei-proxy socket — Python
`KeiProxySocketAuthorizeClient` (`pedro_agentware.kei.runtime_socket`) and
TypeScript `KeiProxySocketAuthorizeClient` (`kei/socketAuthorizeClient.ts`) —
sending the runtime token in `Authorization` and the session in
`X-Kei-Session`. The one new SDK call reuses that same client and transport; it
issues `POST /v1/credentials:release` at call time, after the allow decision,
and hands the returned value straight to the provider call (the `Authorization`
header). It never writes the value to a config file, an environment variable,
or a log.

Per-language name and signature (one method on the existing socket client):

- **Python** — `KeiProxySocketAuthorizeClient.release_credential(session, *, tool="", connector_id="", capability="", resource="", span_id="", parent_span="", delegation_depth=0, agent_id="", framework="", tool_args_digest="") -> CredentialRelease`
- **TypeScript** — `KeiProxySocketAuthorizeClient.releaseCredential(session, input): Promise<CredentialRelease>`
- **Go** — `(*SocketClient).ReleaseCredential(ctx, session, req ReleaseRequest) (ReleaseResult, error)`; Go's socket client mirrors the Python/TypeScript one (Go currently ships the CLI client `CLIClient`; the socket client is the same Unix-socket transport).

A deny is returned as a typed result (the `reason_code`), not raised; a
transport failure (unreachable daemon, timeout, non-2xx) raises the client's
existing error and the harness fails closed.

Session start and end map to load and flush through the SDK's **existing**
session methods — no new session API is added:

- `open_session` / `openSession` (`PUT /v1/session`) is what makes the daemon
  **load** the session's credentials into session memory.
- `close_session` (`DELETE /v1/session`) is what makes the daemon **flush**
  them.
- `release_credential` is the call-time fetch in between.

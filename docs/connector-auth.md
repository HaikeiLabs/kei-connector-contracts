# Connector auth matrix

Which auth method each connector uses, where its credential is placed, whose
account it acts as, and whether it is available or planned. This document is
metadata: it names methods and placements, never values
([CONTRIBUTING.md](../CONTRIBUTING.md)). Tracking: HAI-991.

## Auth methods

| Method | Meaning | Status |
| --- | --- | --- |
| `sso_reuse` | The invoking user's existing SSO identity is reused; no connector-specific registration is stored | available |
| `oauth` | OAuth, in two tiers: tier 1 is a hard-coded SSO app (one Haikei-held registration); tier 2 is a per-org registration the customer brings | available |
| `secret_ref` | An `opaque_ref` pointer into the install's configured secret store; sealed delivery | available |
| `kei_token_exchange` | A short-lived token minted by exchanging the session identity | planned; runtime `docs/oidc-token-exchange.md` |
| `service_principal_key` | A service principal key held in the install's configured secret store; the key-pair JWT is minted harness-side | planned; runtime HAI-993 |
| `saml` | Identity sign-in only; never a connector credential | planned |

Credential placement follows the method. An `oauth`, `sso_reuse` or
`kei_token_exchange` credential lives in the runtime's session memory and is
released at call time over the owner-only Unix socket
([credential pass-through](credential-passthrough.md)). A `secret_ref` or
`service_principal_key` lives in the install's configured secret store —
one store per install, set in its config; for remote cloud deploys that is
AWS Secrets Manager, Azure Key Vault or GCP Secret Manager, with HashiCorp
Vault later — and is released per session. A local bring-your-own harness
(for example Claude Code) uses SSO identity and has no secret store. OAuth
scope handling is [OAuth capability scope contracts](oauth-scopes.md).

## Matrix

One row per connector. "Shared" means one account serves the connector;
"per_user" means the invoking user's own account. A row is "available" when
the connector works today; "planned, TBD" marks what is not decided yet.

| Connector | Auth method(s) | Credential placement | Account model | Status |
| --- | --- | --- | --- | --- |
| Snowflake | `service_principal_key` (key-pair JWT minted harness-side) | install's configured secret store, released per session | shared | planned |
| Postgres / generic database | `secret_ref` (DSN, read-only role) | install's configured secret store | shared | planned |
| generic MCP | `oauth` | runtime session memory | per_user | available |
| generic HTTP API | `secret_ref` | install's configured secret store | shared | available |
| Linear | `oauth` (tier 2, per-org registration) | runtime session memory | per_user | available |
| Notion | `oauth` or `secret_ref` | runtime session memory; install's configured secret store | per_user; shared | available |
| GitHub | `sso_reuse` or `secret_ref` (fine-grained PAT) | runtime session memory; install's configured secret store | per_user; shared | available |
| Omni | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| Power BI | `oauth` (Entra on-behalf-of), moving to `kei_token_exchange` | runtime session memory | per_user | planned, TBD |
| Gmail | `oauth` (tier 1, hard-coded SSO) | runtime session memory | per_user | available |
| Google Sheets | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| Airtable | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| Supabase | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| MySQL | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| MSSQL | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| HubSpot | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| Atlassian | planned, TBD | planned, TBD | planned, TBD | planned, TBD |
| Sentry | planned, TBD | planned, TBD | planned, TBD | planned, TBD |

## Rule

Operation contracts name the auth method and the credential placement, never
values. The runtime releases a credential per session and per call, after an
allow decision (ADR-039; kei-proxy stays decision-only), and the harness
places it on the provider call.

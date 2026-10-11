# OAuth capability scope contracts

`providers.RequiredOAuthScopes`, `ValidateGrantedOAuthScopes`, and
`OAuthReconsentScopes` are pure contract helpers. Consent/connect code should use
`RequiredOAuthScopes` to request the least-privilege union, compare it with the
stored grant, and use `OAuthReconsentScopes` when enabled capabilities change.
Runtime provider `Guard` calls `ValidateGrantedOAuthScopes` before backend
execution for OAuth credentials when `DeclaresOAuthScopes` holds (see below).
Neither path performs a control-plane call.

## When the preflight runs (interim rule)

Connector `scopes` are not always provider OAuth scopes: `kei connectors
create` writes capability names, and the runtime's session path writes a
`session` placeholder. Until the preflight checks the user's actual grant
(the scope the bridge stored on the token, passed through the catalog), the
owner decision of 2026-10-10 (option C) limits it. `Guard` runs the
preflight only when the provider has a scope map **and** the connector's
`scopes` name at least one of that provider's mapped OAuth scopes. Otherwise
it skips the preflight, as v0.7.0 did. That covers GitHub (no map),
capability-name scopes, and session metadata. When the preflight runs, it is
fail-closed.

## Current coverage and integration boundary

`setup.SetupSchemaFor` currently offers OAuth for Google Drive, Gmail, Linear,
and GitHub. This repository contains no connect/consent handler or OAuth
callback consumer: setup declares available auth modes, not requested scopes.
Consequently the helpers are not yet wired to initiate connect or re-consent.
Runtime preflight is wired and denies OAuth providers/capabilities without an
explicit mapping.

Mappings currently cover Google Drive (`drive.search`,
`drive.metadata.read`, `docs.read`), Gmail (`message.search`, `message.get`),
and Linear. GitHub OAuth has no scope contract in this repository: no
provider-specific OAuth registration/authorization configuration declares exact
scope strings or grant semantics. Do not guess these from capability names or
provider permissions. Add the map when the owning OAuth app contract is
available. Until then GitHub OAuth calls skip the preflight (interim rule
above). Opaque-ref credentials are unaffected.

### Linear

Linear scopes are the Linear OAuth app scopes
(<https://linear.app/developers/oauth-2-0-authentication>), requested
comma-separated. Each capability maps to the narrowest scope that covers it:

| Capability | Scope |
| --- | --- |
| `team.read`, `project.read`, `cycle.read`, `issue.read` | `read` |
| `issue.create` | `issues:create` |
| `comment.create` | `comments:create` |
| `issue.update` | `write` |

The connector's `scopes` must cover its declared capabilities, or `Guard`
denies before any network call. kei-oidc-bridge requests
`read,write,issues:create,comments:create` so one grant covers every Linear
capability.

Consent/connect and incremental re-consent remain blocked on integrating these
helpers in the external connect flow and establishing an exact GitHub OAuth
scope contract. HAI-398 should remain open until both are integrated and
verified end to end.

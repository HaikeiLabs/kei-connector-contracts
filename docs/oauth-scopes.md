# OAuth capability scope contracts

`providers.RequiredOAuthScopes`, `ValidateGrantedOAuthScopes`, and
`OAuthReconsentScopes` are pure contract helpers. Consent/connect code should use
`RequiredOAuthScopes` to request the least-privilege union, compare it with the
stored grant, and use `OAuthReconsentScopes` when enabled capabilities change.
Runtime provider `Guard` calls `ValidateGrantedOAuthScopes` before backend
execution for OAuth credentials. Neither path performs a control-plane call.

## Current coverage and integration boundary

`setup.SetupSchemaFor` currently offers OAuth for Google Drive, Gmail, Linear,
and GitHub. This repository contains no connect/consent handler or OAuth
callback consumer: setup declares available auth modes, not requested scopes.
Consequently the helpers are not yet wired to initiate connect or re-consent.
Runtime preflight is wired and denies OAuth providers/capabilities without an
explicit mapping.

Mappings currently cover Google Drive (`drive.search`,
`drive.metadata.read`, `docs.read`) and Gmail (`message.search`, `message.get`).
Linear and GitHub OAuth have no scope contracts in this repository: no
provider-specific OAuth registration/authorization configuration declares exact
scope strings or grant semantics. Do not guess these from capability names or
provider permissions. Add their maps when the owning OAuth app contract is
available. Until then OAuth runtime calls fail closed. Opaque-ref credentials
are unaffected.

Consent/connect and incremental re-consent remain blocked on integrating these
helpers in the external connect flow and establishing exact Linear/GitHub OAuth
scope contracts. HAI-398 should remain open until both are integrated and
verified end to end.

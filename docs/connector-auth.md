# Connector auth matrix

Which auth method each connector uses, where its credential is placed, whose
account it acts as, and whether it is available or planned. This document is
metadata: it names methods and placements, never values
([CONTRIBUTING.md](../CONTRIBUTING.md)).

## Auth methods

| Method | Meaning | Status |
| --- | --- | --- |
| `sso_reuse` | The invoking user's existing SSO identity is reused; no connector-specific registration is stored | available |
| `oauth` | OAuth, in two tiers: tier 1 is a hard-coded SSO app (one Haikei-held registration); tier 2 is a per-org registration the customer brings | available |
| `secret_ref` | An `opaque_ref` pointer into the install's configured secret store; sealed delivery | available |
| `kei_token_exchange` | A short-lived token minted by exchanging the session identity | planned; runtime `docs/oidc-token-exchange.md` |
| `service_principal_key` | A service principal key held in the install's configured secret store; the key-pair JWT is minted harness-side | planned |
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
the connector works today; "planned, TBD" marks what is not decided yet;
"proposed" marks a researched recommendation the owner has not yet decided
(see [Set B auth notes](#set-b-auth-notes)). Database connectors are shared
only for now; a per-user database identity may follow once enterprise SSO
identity is available.

| Connector | Auth method(s) | Credential placement | Account model | Status |
| --- | --- | --- | --- | --- |
| Snowflake | `service_principal_key` (key-pair JWT minted harness-side) | install's configured secret store, released per session | shared | planned |
| Postgres / generic database | `secret_ref` (DSN, read-only role) | install's configured secret store | shared | available |
| generic MCP | `oauth` | runtime session memory | per_user | available |
| generic HTTP API | `secret_ref` | install's configured secret store | shared | available |
| Linear | `oauth` (tier 2, per-org registration) | runtime session memory | per_user | available |
| Notion | `oauth` or `secret_ref` | runtime session memory; install's configured secret store | per_user; shared | available |
| GitHub | `sso_reuse` or `secret_ref` (fine-grained PAT) | runtime session memory; install's configured secret store | per_user; shared | available |
| Omni | proposed: `secret_ref` (personal access token of a dedicated least-privilege Omni user); alt: `secret_ref` organization API key | install's configured secret store | shared | proposed |
| Power BI | `oauth` (Entra on-behalf-of), moving to `kei_token_exchange` | runtime session memory | per_user | planned, TBD |
| Gmail | `oauth` (tier 1, hard-coded SSO) | runtime session memory | per_user | available |
| Google Sheets | proposed: `oauth` (tier 1, the Haikei-held Google app, `spreadsheets.readonly`); alt: `service_principal_key` (Google service account) | runtime session memory; alt: install's configured secret store | per_user; alt: shared | proposed |
| Airtable | proposed: `secret_ref` (personal access token, read scopes, limited to named bases); alt: `oauth` | install's configured secret store; alt: runtime session memory | shared; alt: per_user | proposed |
| Supabase | `secret_ref` (Postgres DSN, read-only role) | install's configured secret store | shared | planned |
| MySQL | `secret_ref` (DSN, read-only role) | install's configured secret store | shared | planned |
| MSSQL | `secret_ref` (DSN, read-only role) | install's configured secret store | shared | planned |
| HubSpot | proposed: `secret_ref` (service key or private-app token, read scopes); alt: `oauth` (tier 2) | install's configured secret store; alt: runtime session memory | shared | proposed |
| Atlassian (Jira + Confluence) | proposed: `secret_ref` (service account OAuth 2.0 client credentials, read scopes); alt: `oauth` (tier 2, 3LO) | install's configured secret store; alt: runtime session memory | shared; alt: per_user | proposed |
| Sentry | proposed: `secret_ref` (internal integration token, read permissions); alt: `oauth` (tier 2) | install's configured secret store; alt: runtime session memory | shared; alt: per_user | proposed |

## Set B auth notes

Every recommendation here is **proposed**; the owner decides. Each picks one
default (least privilege, read-only, installable by the customer's admin) and
names the alternative. Statements are from the vendor's official docs as read
on 2026-11-16; where those docs are silent, the note says so.

**Database connectors (Postgres, MySQL, MSSQL, Supabase).** Decided pattern:
`secret_ref` to a DSN for a read-only database role, held in the install's
configured secret store, shared account. Least privilege is enforced by the
database role's grants, not by Kei. For Supabase the connector uses the
project's Postgres connection string with its own read-only role, not the
project's API keys. Per-user database identity is out of scope until
enterprise SSO identity is available.

**Airtable.** Auth is a personal access token or an OAuth access token;
legacy API keys stopped working on 2024-02-01
([authentication](https://airtable.com/developers/web/api/authentication)).
A token carries scopes and is also limited to the bases or workspaces it is
granted; on Enterprise, a service account's token is created from the admin
panel. OAuth needs an integration registered at `/create/oauth`; access tokens
last 60 minutes and refresh tokens 60 days, refreshable indefinitely, with PKCE
([OAuth reference](https://airtable.com/developers/web/api/oauth-reference)).
Read scopes: `data.records:read`, `schema.bases:read`, and optionally
`data.recordComments:read` ([scopes](https://airtable.com/developers/web/api/scopes)).
Proposed default: `secret_ref` holding a token for a service account (or a
dedicated user) with those read scopes, granted only the named bases; shared.
Alternative: `oauth` per_user, which needs an OAuth integration registration
(Haikei-held or per-org); the docs read do not describe an admin approval step.

**HubSpot.** Options are OAuth apps (multi-account), static-auth private apps
(single account), and service keys
([auth overview](https://developers.hubspot.com/docs/guides/apps/authentication/intro-to-auth)).
Service keys are bearer tokens created by Super Admins or users with the
Developer tools permission, configured with object-level scopes, and do not
expire on a schedule; rotation (manual, or automatic every 3, 6 or 12 months)
replaces them ([service keys](https://developers.hubspot.com/docs/apps/developer-platform/build-apps/authentication/account-service-keys)).
OAuth needs an app registration; only a Super Admin or a user with Marketplace
Access can install it
([working with OAuth](https://developers.hubspot.com/docs/apps/legacy-apps/authentication/working-with-oauth)),
and access tokens last 30 minutes
([token management](https://developers.hubspot.com/docs/api-reference/latest/authentication/manage-oauth-tokens)).
Read scopes are per object, for example `crm.objects.contacts.read`,
`crm.objects.companies.read`, `crm.objects.deals.read`, and the matching
`crm.schemas.*.read` for property metadata
([scopes](https://developers.hubspot.com/docs/apps/legacy-apps/authentication/scopes)).
Proposed default: `secret_ref` holding a service key (or a legacy private-app
token) with only the needed `.read` scopes; shared. Alternative: `oauth`
tier 2. The docs read do not say whether an OAuth install acts with the
installing user's permissions or account-wide, so per_user is not assumed.

**Atlassian (Jira + Confluence, one connector).** Options are OAuth 2.0 (3LO)
apps, scoped API tokens, and service accounts. 3LO needs an app registered in
the developer console; refresh tokens rotate on every use and expire after 90
days of inactivity; `offline_access` is required for a refresh token; one app
covers both Jira and Confluence; org admins can see and remove it under
Connected apps, and the page read does not describe a pre-approval step
([3LO apps](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps/)).
API tokens can be scoped and expire after 1 to 365 days
([API tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/)).
A service account, created by an org admin in Atlassian Administration, can
hold an OAuth 2.0 client-credentials credential with scopes chosen at creation
(Jira, Confluence, Goals, Projects only); its access tokens last 60 minutes
([service account OAuth credential](https://support.atlassian.com/user-management/docs/create-oauth-2-0-credential-for-service-accounts/)).
Read scopes: Jira `read:jira-work`, `read:jira-user`; Confluence
`read:confluence-content.all`, `read:confluence-space.summary`,
`search:confluence` (granular equivalents include `read:page:confluence`,
`read:space:confluence`)
([Confluence scopes](https://developer.atlassian.com/cloud/confluence/scopes-for-oauth-2-3LO-and-forge-apps/)).
Proposed default: `secret_ref` holding the service account's client ID and
secret with read scopes only, the 60-minute token minted harness-side; shared,
with what it can see set by the service account's project and space
permissions. Alternative: `oauth` tier 2 (3LO) per_user, so results follow
each user's own Jira and Confluence permissions.

**Sentry.** Options are user auth tokens, organization tokens from an internal
integration, OAuth (authorization code with PKCE; access tokens expire after
30 days and refresh), DSN auth for a few endpoints, and legacy API keys, which
are disabled for new accounts ([API auth](https://docs.sentry.io/api/auth/)).
An internal integration issues an org-wide token that does not expire until
revoked ([internal integrations](https://docs.sentry.io/organization/integrations/integration-platform/internal-integration/)).
Read scopes: `org:read`, `project:read`, `team:read`, `member:read`,
`event:read` (issues and events)
([permissions](https://docs.sentry.io/api/permissions/)). The pages read do not
say which org role may create an internal integration. Proposed default:
`secret_ref` holding an internal integration token with read permissions only;
shared. Alternative: `oauth` tier 2 per_user.

**Google Sheets.** Read scope `spreadsheets.readonly` is classed sensitive and
needs Google's additional app verification; `drive.readonly` is restricted;
`drive.file` (per-file, non-sensitive) only covers files the app opened or the
user picked, which does not fit an agent reading arbitrary sheets
([Sheets scopes](https://developers.google.com/workspace/sheets/api/scopes)).
Service accounts sign a key-pair JWT for a one-hour access token; reading
Workspace users' data needs domain-wide delegation, which a Workspace super
admin grants and which impersonates users
([service accounts](https://developers.google.com/identity/protocols/oauth2/service-account)).
That page does not cover sharing a sheet directly with a service account.
Proposed default: `oauth` tier 1 per_user, reusing the Haikei-held Google app
that serves Gmail, with `spreadsheets.readonly` added (requires that app's
verification to cover the scope). Alternative: `service_principal_key` shared
with a Google service account; domain-wide delegation is not recommended
because it is broader than read access to sheets.

**Omni.** The API takes bearer tokens of two kinds
([authentication](https://docs.omni.co/api/authentication)): organization API
keys, created by Organization Admins, carrying Organization Admin permissions
and no expiry; and personal access tokens, which an Organization Admin must
enable, which a user with Restricted Querier or higher creates, which inherit
that user's permissions and do not expire. Some endpoints (document
export/import, SCIM) accept only organization keys. OAuth exists only for the
Omni CLI and MCP server, not for direct API calls
([MCP authentication](https://docs.omni.co/ai/mcp/authentication)). Omni has no
API scopes; access is the token owner's role. Proposed default: `secret_ref`
holding a personal access token of a dedicated Omni user with the least role
that can query; shared. Alternative: `secret_ref` organization API key, only
if admin-only endpoints are needed, since it carries admin permissions.

## Rule

Operation contracts name the auth method and the credential placement, never
values. The runtime releases a credential per session and per call, after an
allow decision (ADR-039; kei-proxy stays decision-only), and the harness
places it on the provider call.

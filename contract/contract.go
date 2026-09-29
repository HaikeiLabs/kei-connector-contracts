// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package contract contains the provider-neutral contract used by Kei's
// control plane and connector clients.  It intentionally contains references
// to credentials, never credential material.
package contract

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Provider string

const (
	ProviderCRM     Provider = "crm"
	ProviderLinear  Provider = "linear"
	ProviderGitHub  Provider = "github"
	ProviderGoogle  Provider = "google_drive"
	ProviderNotion  Provider = "notion"
	ProviderS3      Provider = "s3"
	ProviderHTTPAPI Provider = "http_api"

	// ProviderFreshBooks is the bookkeeping provider. Its capability surface
	// is read-only: invoices, expenses, payments and clients are read to
	// reconcile and report, never created or altered through Kei.
	ProviderFreshBooks Provider = "freshbooks"

	// ProviderMercury is the banking provider. Its capability surface is
	// read-only by construction, matching the provider's own read-only
	// token and MCP surface: no capability here can move money.
	ProviderMercury Provider = "mercury"

	// ProviderGmail is its own provider rather than a google_drive capability:
	// it needs a separate OAuth grant (gmail.readonly) and a separate
	// resource surface, and a Drive connector must never imply mail access.
	ProviderGmail Provider = "gmail"

	// ProviderTito is the Tito ticketing admin API. It has no OAuth: its API
	// token lives in the customer's secret backend behind an opaque_ref.
	ProviderTito Provider = "tito"
)

// CredentialSource is typed metadata for where a connector's credential
// comes from. It is a closed set, not a string namespace: the source must
// never be encoded into credential_ref itself.
type CredentialSource string

const (
	// CredentialSourceOAuth marks a Kei-managed connector OAuth credential:
	// the token is bridge-owned, and the control plane fetches it for one
	// allowed invocation. Which account the token belongs to is the
	// connector's AccountModel.
	CredentialSourceOAuth CredentialSource = "oauth"
	// CredentialSourceOpaqueRef is the shared-secret path: credential_ref
	// points at the customer's secret backend and the tenant runtime resolves
	// it. The value avoids the substring "secret" so API redaction guards that
	// reject credential-like terms do not flag this metadata field.
	CredentialSourceOpaqueRef CredentialSource = "opaque_ref"
)

// Definitions are the deliberately small initial provider surface. New
// operations must be added here before an agent can request them.
var definitions = map[Provider][]Capability{
	// Investor reads serve fundraising (HAI-210). Investor stage changes are
	// agent action tools, never CRM capabilities.
	ProviderCRM:     {{Name: "lead.read", Action: ActionRead}, {Name: "lead.create", Action: ActionCreate}, {Name: "lead.update", Action: ActionUpdate}, {Name: "investor.list", Action: ActionRead}, {Name: "investor.read", Action: ActionRead}},
	ProviderLinear:  {{Name: "team.read", Action: ActionRead}, {Name: "project.read", Action: ActionRead}, {Name: "cycle.read", Action: ActionRead}, {Name: "issue.read", Action: ActionRead}, {Name: "issue.create", Action: ActionCreate}, {Name: "issue.update", Action: ActionUpdate}},
	ProviderGitHub:  {{Name: "repository.read", Action: ActionRead}, {Name: "issue.read", Action: ActionRead}, {Name: "pull_request.read", Action: ActionRead}, {Name: "check.read", Action: ActionRead}, {Name: "workflow.read", Action: ActionRead}, {Name: "issue.create", Action: ActionCreate}, {Name: "issue.update", Action: ActionUpdate}, {Name: "pull_request.create", Action: ActionCreate}, {Name: "pull_request.update", Action: ActionUpdate}, {Name: "issue.comment", Action: ActionComment}},
	ProviderGoogle:  {{Name: "drive.search", Action: ActionRead}, {Name: "drive.metadata.read", Action: ActionRead}, {Name: "docs.read", Action: ActionRead}},
	ProviderNotion:  {{Name: "search", Action: ActionRead}, {Name: "page.read", Action: ActionRead}, {Name: "database.query", Action: ActionRead}},
	ProviderS3:      {{Name: "object.list", Action: ActionRead}, {Name: "object.read", Action: ActionRead}},
	ProviderHTTPAPI: {{Name: "http.get", Action: ActionRead}, {Name: "http.head", Action: ActionRead}},
	// Finance providers declare reads only. There is deliberately no
	// payment, transfer, invoice.create or expense.update capability: a
	// name absent from this catalog cannot be invoked at all, so the
	// read-only guarantee does not depend on policy being configured
	// correctly. See read_only_test.go.
	ProviderFreshBooks: {{Name: "invoice.read", Action: ActionRead}, {Name: "expense.read", Action: ActionRead}, {Name: "payment.read", Action: ActionRead}, {Name: "client.read", Action: ActionRead}},
	ProviderMercury:    {{Name: "account.read", Action: ActionRead}, {Name: "transaction.read", Action: ActionRead}, {Name: "balance.read", Action: ActionRead}},
	ProviderGmail:      {{Name: "message.search", Action: ActionRead}, {Name: "message.get", Action: ActionRead}},
	ProviderTito:       {{Name: "event.list", Action: ActionRead}, {Name: "event.get", Action: ActionRead}, {Name: "release.list", Action: ActionRead}, {Name: "ticket.summary", Action: ActionRead}},
}

func CapabilitiesFor(provider Provider) []Capability {
	capabilities := definitions[provider]
	result := make([]Capability, len(capabilities))
	copy(result, capabilities)
	return result
}

// PolicyAttributes is the connector's declared credential-surface policy. It
// is dormant on the connector: enforcement moved to the ABAC policy layer
// (governance pivot), so these fields round-trip for compatibility but are
// not consulted by ValidateCall or Decide.
type PolicyAttributes struct {
	AllowedActions     []Action `json:"allowed_actions"`
	AllowedResources   []string `json:"allowed_resources"`
	AllowedPrefixes    []string `json:"allowed_prefixes,omitempty"`
	DestructiveEnabled bool     `json:"destructive_enabled"`
	// GmailIncludeBody is the explicit opt-in for message.get to return the
	// message body. Unlike the fields above it is consulted, by the Gmail
	// client: without it, results carry metadata and snippet only.
	GmailIncludeBody bool `json:"gmail_include_body,omitempty"`
}

type Metadata struct {
	ID            string           `json:"id"`
	TenantID      string           `json:"tenant_id"`
	WorkspaceID   string           `json:"workspace_id"`
	Name          string           `json:"name"`
	Provider      Provider         `json:"provider"`
	Status        Status           `json:"status"`
	CredentialRef string           `json:"credential_ref"`
	Scopes        []string         `json:"scopes"`
	Resources     []string         `json:"resources"`
	Policy        PolicyAttributes `json:"policy"`
	Capabilities  []Capability     `json:"capabilities"`
	CreatedBy     string           `json:"created_by"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	HTTPAPI       *HTTPAPI         `json:"http_api,omitempty"`

	// CredentialSource is the typed credential source. Empty means
	// opaque_ref, so metadata built against v0.1.0 stays valid; see
	// EffectiveCredentialSource.
	CredentialSource CredentialSource `json:"credential_source,omitempty"`
	// Subject is the canonical subject the connector's credential is bound
	// to. A legacy oauth connector (no AccountModel) must declare it; a
	// per_user connector must not; a shared one uses SharedSubject(ID).
	Subject string `json:"subject,omitempty"`
	// AccountModel is whose account an OAuth-style connector acts as. Empty
	// keeps the legacy single-subject oauth rules.
	AccountModel AccountModel `json:"account_model,omitempty"`
	// Config holds the provider's non-secret setup fields (see SetupSchemaFor).
	// Secret values never appear here; they are behind CredentialRef.
	Config map[string]any `json:"config,omitempty"`
}

// EffectiveCredentialSource returns the connector's credential source, with an
// empty value meaning opaque_ref.
func (m Metadata) EffectiveCredentialSource() CredentialSource {
	if m.CredentialSource == "" {
		return CredentialSourceOpaqueRef
	}
	return m.CredentialSource
}

type Invocation struct {
	TenantID       string          `json:"tenant_id"`
	WorkspaceID    string          `json:"workspace_id"`
	Subject        string          `json:"subject"`
	AgentID        string          `json:"agent_id"`
	ConnectorID    string          `json:"connector_id"`
	Capability     string          `json:"capability"`
	Action         Action          `json:"action"`
	Resource       string          `json:"resource"`
	TraceID        string          `json:"trace_id"`
	IdempotencyKey string          `json:"idempotency_key,omitempty"`
	ApprovalID     string          `json:"approval_id,omitempty"`
	HTTP           *HTTPInvocation `json:"http,omitempty"`
}

var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,254}$`)

func validID(name, value string) error {
	if value == "" || !identifier.MatchString(value) {
		return fmt.Errorf("%s is invalid", name)
	}
	return nil
}

func (m Metadata) Validate() error {
	for name, value := range map[string]string{
		"id": m.ID, "tenant_id": m.TenantID, "workspace_id": m.WorkspaceID,
		"name": m.Name, "credential_ref": m.CredentialRef, "created_by": m.CreatedBy,
	} {
		if err := validID(name, value); err != nil {
			return err
		}
	}
	if !validProvider(m.Provider) {
		return fmt.Errorf("unsupported provider %q", m.Provider)
	}
	if !validStatus(m.Status) {
		return fmt.Errorf("invalid connector status %q", m.Status)
	}
	if err := m.validateCredentialBinding(); err != nil {
		return err
	}
	if m.Provider == ProviderHTTPAPI {
		if m.HTTPAPI == nil {
			return errors.New("http_api metadata is required")
		}
		if err := m.HTTPAPI.Validate(); err != nil {
			return err
		}
	} else if m.HTTPAPI != nil {
		return errors.New("http_api metadata is only valid for provider http_api")
	}
	if len(m.Scopes) == 0 {
		return errors.New("connector must declare at least one scope")
	}
	if len(m.Capabilities) == 0 {
		return errors.New("connector must declare capabilities")
	}
	for _, c := range m.Capabilities {
		if containsSecretMaterial(c.Description) {
			return errors.New("metadata must not contain secret material")
		}
		if err := validID("capability", c.Name); err != nil {
			return err
		}
		if !validAction(c.Action) {
			return fmt.Errorf("invalid capability action %q", c.Action)
		}
		if !definedCapability(m.Provider, c) {
			return fmt.Errorf("capability %q is not defined for provider %q", c.Name, m.Provider)
		}
	}
	return nil
}

func containsSecretMaterial(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"bearer ", "basic ", "api_key=", "api-key=", "secret=", "token=", "sk_", "ghp_", "gho_", "xoxb-"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func definedCapability(provider Provider, candidate Capability) bool {
	for _, capability := range definitions[provider] {
		if capability.Name == candidate.Name && capability.Action == candidate.Action {
			return true
		}
	}
	return false
}

func validProvider(p Provider) bool {
	return p == ProviderCRM || p == ProviderLinear || p == ProviderGitHub || p == ProviderGoogle || p == ProviderNotion || p == ProviderS3 || p == ProviderHTTPAPI ||
		p == ProviderFreshBooks || p == ProviderMercury || p == ProviderGmail || p == ProviderTito
}
func validStatus(s Status) bool {
	return s == StatusPending || s == StatusActive || s == StatusSuspended || s == StatusRevoked || s == StatusFailed
}
func validAction(a Action) bool {
	return a == ActionRead || a == ActionCreate || a == ActionUpdate || a == ActionDelete || a == ActionComment
}

// ValidateCall is fail-closed: every boundary is explicit, and a missing
// tenant/workspace, inactive connector, undeclared capability, or action
// mismatch is denied. Connector policy — allowed resources, actions, and
// approvals — is decided by the ABAC policy layer, not by the connector's
// credential surface.
func ValidateCall(m Metadata, in Invocation) error {
	if err := m.Validate(); err != nil {
		return err
	}
	for name, value := range map[string]string{"tenant_id": in.TenantID, "workspace_id": in.WorkspaceID, "subject": in.Subject, "agent_id": in.AgentID, "connector_id": in.ConnectorID, "capability": in.Capability, "trace_id": in.TraceID} {
		if err := validID(name, value); err != nil {
			return err
		}
	}
	if m.Status != StatusActive {
		return fmt.Errorf("connector is not active")
	}
	if in.TenantID != m.TenantID || in.WorkspaceID != m.WorkspaceID || in.ConnectorID != m.ID {
		return errors.New("connector not found")
	}
	if in.Resource == "" {
		return errors.New("resource is required")
	}
	// A legacy oauth connector's credential is bound to the subject it
	// declared; an invocation for any other subject would act under the
	// wrong identity. per_user and shared connectors are gated by policy.
	if m.AccountModel == "" && m.EffectiveCredentialSource() == CredentialSourceOAuth && in.Subject != m.Subject {
		return errors.New("invoking subject does not match the connector's declared subject")
	}
	if m.Provider == ProviderHTTPAPI {
		if in.HTTP == nil {
			return errors.New("http invocation is required")
		}
		if err := m.HTTPAPI.ValidateInvocation(*in.HTTP); err != nil {
			return err
		}
	}
	for _, c := range m.Capabilities {
		if c.Name == in.Capability {
			if c.Action != in.Action {
				return errors.New("capability action mismatch")
			}
			return nil
		}
	}
	return errors.New("capability is not allowed")
}

// ValidateCredentialRef accepts opaque references only. URLs and inline
// secrets are rejected so provider credentials remain in Kei's credential
// service/secret backend. Any scheme-bearing reference is rejected: the
// credential type/source is typed metadata (CredentialSource + Subject), not a
// string namespace. The legacy kei-oauth namespace, in both its colon and
// hyphen forms, is migration-only and is rejected.
func ValidateCredentialRef(ref string) error {
	if ref == "" || len(ref) > 512 {
		return errors.New("credential_ref is required")
	}
	if u, err := url.Parse(ref); err == nil && u.Scheme != "" {
		return errors.New("credential_ref must be an opaque reference")
	}
	if strings.ContainsAny(ref, "\r\n\t") {
		return errors.New("credential_ref contains invalid characters")
	}
	lower := strings.ToLower(ref)
	if strings.HasPrefix(lower, "kei-oauth:") || strings.HasPrefix(lower, "kei-oauth-") {
		return errors.New("credential_ref must not use the legacy kei-oauth namespace")
	}
	if strings.Contains(ref, "=") || strings.HasPrefix(lower, "bearer ") || strings.HasPrefix(lower, "sk_") || strings.HasPrefix(lower, "ghp_") || strings.HasPrefix(lower, "gho_") || strings.HasPrefix(lower, "xoxb-") {
		return errors.New("credential_ref must not contain credential material")
	}
	return nil
}

// ValidateInvocation checks an invocation's identifiers, action, and
// resource. Envelope validation applies it to the invocation it carries.
func ValidateInvocation(in Invocation) error {
	for name, value := range map[string]string{
		"tenant_id":    in.TenantID,
		"workspace_id": in.WorkspaceID,
		"subject":      in.Subject,
		"agent_id":     in.AgentID,
		"connector_id": in.ConnectorID,
		"capability":   in.Capability,
		"trace_id":     in.TraceID,
	} {
		if err := validID(name, value); err != nil {
			return err
		}
	}
	if !validAction(in.Action) {
		return fmt.Errorf("invalid action %q", in.Action)
	}
	if in.Resource == "" {
		return errors.New("resource is required")
	}
	if in.IdempotencyKey != "" {
		if err := validID("idempotency_key", in.IdempotencyKey); err != nil {
			return err
		}
	}
	if in.ApprovalID != "" {
		if err := validID("approval_id", in.ApprovalID); err != nil {
			return err
		}
	}
	return nil
}

// ValidProvider reports whether p is a provider the contract defines.
func ValidProvider(p Provider) bool { return validProvider(p) }

// DefinedProviders returns every provider with a capability definition, in no
// particular order.
func DefinedProviders() []Provider {
	out := make([]Provider, 0, len(definitions))
	for p := range definitions {
		out = append(out, p)
	}
	return out
}

// Package connectors contains the provider-neutral contract used by Kei's
// control plane and connector clients.  It intentionally contains references
// to credentials, never credential material.
package connectors

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
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusRevoked   Status = "revoked"
	StatusFailed    Status = "failed"
)

type Action string

const (
	ActionRead    Action = "read"
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionComment Action = "comment"
)

// Capability is the only provider operation a connector client may expose.
// It describes the credential surface: which operations a connector's
// credential can perform. Data-access and tool-call authorization — including
// which capabilities require an approval — are decided by ABAC policies, not
// by capability flags on the connector.
type Capability struct {
	Name        string `json:"name"`
	Action      Action `json:"action"`
	Description string `json:"description"`
}

// Definitions are the deliberately small initial provider surface. New
// operations must be added here before an agent can request them.
var definitions = map[Provider][]Capability{
	ProviderCRM:     {{Name: "lead.read", Action: ActionRead}, {Name: "lead.create", Action: ActionCreate}, {Name: "lead.update", Action: ActionUpdate}},
	ProviderLinear:  {{Name: "team.read", Action: ActionRead}, {Name: "project.read", Action: ActionRead}, {Name: "cycle.read", Action: ActionRead}, {Name: "issue.read", Action: ActionRead}, {Name: "issue.create", Action: ActionCreate}, {Name: "issue.update", Action: ActionUpdate}},
	ProviderGitHub:  {{Name: "repository.read", Action: ActionRead}, {Name: "issue.read", Action: ActionRead}, {Name: "pull_request.read", Action: ActionRead}, {Name: "check.read", Action: ActionRead}, {Name: "workflow.read", Action: ActionRead}, {Name: "issue.create", Action: ActionCreate}, {Name: "issue.update", Action: ActionUpdate}, {Name: "pull_request.create", Action: ActionCreate}, {Name: "pull_request.update", Action: ActionUpdate}, {Name: "issue.comment", Action: ActionComment}},
	ProviderGoogle:  {{Name: "drive.search", Action: ActionRead}, {Name: "drive.metadata.read", Action: ActionRead}, {Name: "docs.read", Action: ActionRead}},
	ProviderNotion:  {{Name: "search", Action: ActionRead}, {Name: "page.read", Action: ActionRead}, {Name: "database.query", Action: ActionRead}},
	ProviderS3:      {{Name: "object.list", Action: ActionRead}, {Name: "object.read", Action: ActionRead}},
	ProviderHTTPAPI: {{Name: "http.get", Action: ActionRead}, {Name: "http.head", Action: ActionRead}},
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
	return p == ProviderCRM || p == ProviderLinear || p == ProviderGitHub || p == ProviderGoogle || p == ProviderNotion || p == ProviderS3 || p == ProviderHTTPAPI
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
// service/secret backend.
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
	if strings.Contains(ref, "=") || strings.HasPrefix(lower, "bearer ") || strings.HasPrefix(lower, "sk_") || strings.HasPrefix(lower, "ghp_") || strings.HasPrefix(lower, "gho_") || strings.HasPrefix(lower, "xoxb-") {
		return errors.New("credential_ref must not contain credential material")
	}
	return nil
}

package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type LinearTeam struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type LinearProject struct {
	ID      string `json:"id"`
	TeamKey string `json:"team_key"`
	Name    string `json:"name"`
	State   string `json:"state"`
}

type LinearCycle struct {
	ID       string    `json:"id"`
	TeamKey  string    `json:"team_key"`
	Name     string    `json:"name"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type LinearIssue struct {
	ID         string `json:"id"`
	TeamKey    string `json:"team_key"`
	ProjectID  string `json:"project_id,omitempty"`
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	State      string `json:"state"`
}

type TeamReadPayload struct{}

func (TeamReadPayload) Capability() string { return "team.read" }

type ProjectReadPayload struct{}

func (ProjectReadPayload) Capability() string { return "project.read" }

type CycleReadPayload struct{}

func (CycleReadPayload) Capability() string { return "cycle.read" }

type LinearIssueReadPayload struct{}

func (LinearIssueReadPayload) Capability() string { return "issue.read" }

// LinearBackend is the seam for Linear reads.
type LinearBackend interface {
	Team(ctx context.Context, key string) (LinearTeam, error)
	Project(ctx context.Context, id string) (LinearProject, error)
	Cycle(ctx context.Context, id string) (LinearCycle, error)
	Issue(ctx context.Context, id string) (LinearIssue, error)
}

// LinearCredentialScope is the invocation binding supplied to credential
// resolution. CredentialRef remains an opaque location; token material never
// enters the governed metadata or invocation shapes.
type LinearCredentialScope struct {
	TenantID    string
	WorkspaceID string
	Subject     string
}

// LinearCredentialResolver resolves a connector credential for one governed
// invocation. Implementations must enforce the binding represented by Scope.
type LinearCredentialResolver interface {
	Resolve(context.Context, string, LinearCredentialScope) (string, error)
}

// LinearGraphQLTransport is the authenticated GraphQL transport seam. The
// endpoint is supplied by the client and is always the fixed Linear endpoint.
// Implementations must not log or include Token in returned errors.
type LinearGraphQLTransport interface {
	Query(context.Context, string, string, string, map[string]any) (json.RawMessage, error)
}

// MemoryLinear remains an explicit offline test backend. It is never used by
// runtime construction.
type MemoryLinear struct {
	Teams    map[string]LinearTeam
	Projects map[string]LinearProject
	Cycles   map[string]LinearCycle
	Issues   map[string]LinearIssue
}

func (m MemoryLinear) Team(_ context.Context, key string) (LinearTeam, error) {
	t, ok := m.Teams[key]
	if !ok {
		return LinearTeam{}, errors.New("team not found")
	}
	return t, nil
}

func (m MemoryLinear) Project(_ context.Context, id string) (LinearProject, error) {
	p, ok := m.Projects[id]
	if !ok {
		return LinearProject{}, errors.New("project not found")
	}
	return p, nil
}

func (m MemoryLinear) Cycle(_ context.Context, id string) (LinearCycle, error) {
	c, ok := m.Cycles[id]
	if !ok {
		return LinearCycle{}, errors.New("cycle not found")
	}
	return c, nil
}

func (m MemoryLinear) Issue(_ context.Context, id string) (LinearIssue, error) {
	i, ok := m.Issues[id]
	if !ok {
		return LinearIssue{}, errors.New("issue not found")
	}
	return i, nil
}

// LinearClient serves team.read, project.read, cycle.read, and issue.read.
// Resources are linear/team/<key>, linear/project/<id>, linear/cycle/<id>,
// and linear/issue/<id>.
type LinearClient struct {
	backend     LinearBackend
	transport   LinearGraphQLTransport
	credentials LinearCredentialResolver
}

func NewLinear(backend LinearBackend) *LinearClient { return &LinearClient{backend: backend} }

// NewLinearRuntime constructs the authenticated, read-only Linear provider.
// Both dependencies are required at invocation time; a nil dependency fails
// closed without attempting a network call or exposing credential material.
func NewLinearRuntime(transport LinearGraphQLTransport, credentials LinearCredentialResolver) *LinearClient {
	return &LinearClient{transport: transport, credentials: credentials}
}

func (c *LinearClient) Provider() connectors.Provider { return connectors.ProviderLinear }

func (c *LinearClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	kind, id, ok := linearResource(inv.Resource)
	if !ok {
		return Result{}, errors.New("resource must be linear/<team|project|cycle|issue>/<id>")
	}
	if c.backend == nil && (c.transport == nil || c.credentials == nil) {
		return Result{}, errors.New("linear runtime dependencies are unavailable")
	}
	if c.backend == nil {
		token, err := c.resolveToken(ctx, meta, inv)
		if err != nil {
			return Result{}, err
		}
		backend := linearGraphQLBackend{transport: c.transport, token: token}
		return c.invokeBackend(ctx, inv, payload, kind, id, backend)
	}
	return c.invokeBackend(ctx, inv, payload, kind, id, c.backend)
}

func (c *LinearClient) resolveToken(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation) (string, error) {
	if err := connectors.ValidateCredentialRef(meta.CredentialRef); err != nil {
		return "", errors.New("linear credential is unavailable")
	}
	token, err := c.credentials.Resolve(ctx, meta.CredentialRef, LinearCredentialScope{
		TenantID: inv.TenantID, WorkspaceID: inv.WorkspaceID, Subject: inv.Subject,
	})
	if err != nil || strings.TrimSpace(token) == "" {
		return "", errors.New("linear credential is unavailable")
	}
	return token, nil
}

func (c *LinearClient) invokeBackend(ctx context.Context, inv connectors.Invocation, payload Payload, kind, id string, backend LinearBackend) (Result, error) {
	switch payload.(type) {
	case TeamReadPayload:
		if kind != "team" {
			return Result{}, errors.New("team.read requires a linear/team/<key> resource")
		}
		data, err := backend.Team(ctx, id)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case ProjectReadPayload:
		if kind != "project" {
			return Result{}, errors.New("project.read requires a linear/project/<id> resource")
		}
		data, err := backend.Project(ctx, id)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case CycleReadPayload:
		if kind != "cycle" {
			return Result{}, errors.New("cycle.read requires a linear/cycle/<id> resource")
		}
		data, err := backend.Cycle(ctx, id)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case LinearIssueReadPayload:
		if kind != "issue" {
			return Result{}, errors.New("issue.read requires a linear/issue/<id> resource")
		}
		data, err := backend.Issue(ctx, id)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func linearResource(resource string) (kind, id string, ok bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 3 || parts[0] != "linear" || parts[2] == "" {
		return "", "", false
	}
	switch parts[1] {
	case "team", "project", "cycle", "issue":
		return parts[1], parts[2], true
	default:
		return "", "", false
	}
}

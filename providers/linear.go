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

package providers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
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

// Linear write bounds. Title and body limits keep a write inside the
// runtime's 16 KiB input cap; priority is Linear's 0 (none) to 4 (low).
const (
	MaxLinearTitleRunes       = 256
	MaxLinearDescriptionBytes = 12 << 10
	MaxLinearCommentBytes     = 12 << 10
	MaxLinearPriority         = 4
)

// LinearIssueCreatePayload requests issue.create. The team is the issue's
// parent (Haikei semantics rule), so TeamID is required and must equal the
// team named by the linear/team/<team> resource. TeamID is the team key
// (for example "ENG") or the team's Linear id. The workspace is never an
// argument: it is the connector's OAuth grant.
type LinearIssueCreatePayload struct {
	TeamID      string `json:"team_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Priority    *int   `json:"priority,omitempty"`
}

func (LinearIssueCreatePayload) Capability() string { return "issue.create" }

// Validate checks the payload against the issue.create bounds.
func (p LinearIssueCreatePayload) Validate() error {
	switch {
	case !linearID.MatchString(p.TeamID):
		return errors.New("team_id is invalid")
	case strings.TrimSpace(p.Title) == "" || utf8.RuneCountInString(p.Title) > MaxLinearTitleRunes || strings.ContainsAny(p.Title, "\r\n"):
		return errors.New("title is invalid")
	case len(p.Description) > MaxLinearDescriptionBytes || !utf8.ValidString(p.Description):
		return errors.New("description is invalid")
	case p.Priority != nil && (*p.Priority < 0 || *p.Priority > MaxLinearPriority):
		return errors.New("priority is out of range")
	}
	return nil
}

// LinearCommentCreatePayload requests comment.create. The issue is the
// comment's parent, so IssueID is required and must equal the issue named by
// the linear/issue/<issue> resource. IssueID is the identifier ("ENG-42") or
// the issue's Linear id.
type LinearCommentCreatePayload struct {
	IssueID string `json:"issue_id"`
	Body    string `json:"body"`
}

func (LinearCommentCreatePayload) Capability() string { return "comment.create" }

// Validate checks the payload against the comment.create bounds.
func (p LinearCommentCreatePayload) Validate() error {
	switch {
	case !linearID.MatchString(p.IssueID):
		return errors.New("issue_id is invalid")
	case strings.TrimSpace(p.Body) == "" || len(p.Body) > MaxLinearCommentBytes || !utf8.ValidString(p.Body):
		return errors.New("body is invalid")
	}
	return nil
}

// linearID is a Linear team key, issue identifier, or entity id.
var linearID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

// LinearIssueCreated is the issue.create result. It names the new issue and
// carries none of the submitted title or description.
type LinearIssueCreated struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	TeamKey    string `json:"team_key"`
	URL        string `json:"url"`
}

// LinearCommentCreated is the comment.create result. It carries no body.
type LinearCommentCreated struct {
	ID      string `json:"id"`
	IssueID string `json:"issue_id"`
	URL     string `json:"url"`
}

// LinearWriteBackend is the seam for Linear writes. ClientID is a
// deterministic UUID derived from the invocation's idempotency key: Linear
// accepts it as the new entity's id, so a replayed write cannot create a
// second issue or comment.
type LinearWriteBackend interface {
	CreateIssue(ctx context.Context, clientID string, in LinearIssueCreatePayload) (LinearIssueCreated, error)
	CreateComment(ctx context.Context, clientID string, in LinearCommentCreatePayload) (LinearCommentCreated, error)
}

// LinearWriteID derives the Linear entity id for one write from the
// connector, the invoking subject, the capability, and the idempotency key.
// It is a version 4 shaped UUID so Linear accepts it as an input id.
func LinearWriteID(inv contract.Invocation) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{inv.TenantID, inv.WorkspaceID, inv.ConnectorID, inv.Subject, inv.Capability, inv.IdempotencyKey}, "\x00")))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

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
// runtime construction. Writes are recorded in Issues and Comments, keyed by
// the client id, so a replay returns the first result.
type MemoryLinear struct {
	Teams    map[string]LinearTeam
	Projects map[string]LinearProject
	Cycles   map[string]LinearCycle
	Issues   map[string]LinearIssue
	Comments map[string]LinearCommentCreated
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

func (m MemoryLinear) CreateIssue(_ context.Context, clientID string, in LinearIssueCreatePayload) (LinearIssueCreated, error) {
	team, ok := m.Teams[in.TeamID]
	if !ok || m.Issues == nil {
		return LinearIssueCreated{}, errors.New("team not found")
	}
	if existing, ok := m.Issues[clientID]; ok {
		return LinearIssueCreated{ID: existing.ID, Identifier: existing.Identifier, TeamKey: existing.TeamKey}, nil
	}
	issue := LinearIssue{ID: clientID, TeamKey: team.Key, Identifier: fmt.Sprintf("%s-%d", team.Key, len(m.Issues)+1), Title: in.Title, State: "Triage"}
	m.Issues[clientID] = issue
	return LinearIssueCreated{ID: issue.ID, Identifier: issue.Identifier, TeamKey: issue.TeamKey}, nil
}

func (m MemoryLinear) CreateComment(_ context.Context, clientID string, in LinearCommentCreatePayload) (LinearCommentCreated, error) {
	if m.Comments == nil {
		return LinearCommentCreated{}, errors.New("comments are unavailable")
	}
	if existing, ok := m.Comments[clientID]; ok {
		return existing, nil
	}
	for _, issue := range m.Issues {
		if issue.ID == in.IssueID || issue.Identifier == in.IssueID {
			comment := LinearCommentCreated{ID: clientID, IssueID: issue.ID}
			m.Comments[clientID] = comment
			return comment, nil
		}
	}
	return LinearCommentCreated{}, errors.New("issue not found")
}

// LinearClient serves team.read, project.read, cycle.read, issue.read,
// issue.create, and comment.create. Read resources are linear/team/<key>,
// linear/project/<id>, linear/cycle/<id>, and linear/issue/<id>. A write's
// resource is its parent: issue.create on linear/team/<team>, comment.create
// on linear/issue/<issue>.
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

func (c *LinearClient) Provider() contract.Provider { return contract.ProviderLinear }

func (c *LinearClient) Invoke(ctx context.Context, meta contract.Metadata, inv contract.Invocation, payload Payload) (Result, error) {
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

// invokeWrite runs issue.create or comment.create. Writes require an
// idempotency key (connector execution contract, writes) and a payload whose
// parent identifier equals the resource. The client never retries a write.
func (c *LinearClient) invokeWrite(ctx context.Context, inv contract.Invocation, payload Payload, kind, id string, backend LinearBackend) (Result, error) {
	writer, ok := backend.(LinearWriteBackend)
	if !ok {
		return Result{}, errors.New("linear writes are unavailable")
	}
	if inv.IdempotencyKey == "" {
		return Result{}, errors.New("linear writes require an idempotency key")
	}
	clientID := LinearWriteID(inv)
	switch p := payload.(type) {
	case LinearIssueCreatePayload:
		if err := p.Validate(); err != nil {
			return Result{}, err
		}
		if kind != "team" || p.TeamID != id {
			return Result{}, errors.New("issue.create requires a linear/team/<team> resource matching team_id")
		}
		data, err := writer.CreateIssue(ctx, clientID, p)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case LinearCommentCreatePayload:
		if err := p.Validate(); err != nil {
			return Result{}, err
		}
		if kind != "issue" || p.IssueID != id {
			return Result{}, errors.New("comment.create requires a linear/issue/<issue> resource matching issue_id")
		}
		data, err := writer.CreateComment(ctx, clientID, p)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func (c *LinearClient) resolveToken(ctx context.Context, meta contract.Metadata, inv contract.Invocation) (string, error) {
	if err := contract.ValidateCredentialRef(meta.CredentialRef); err != nil {
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

func (c *LinearClient) invokeBackend(ctx context.Context, inv contract.Invocation, payload Payload, kind, id string, backend LinearBackend) (Result, error) {
	switch payload.(type) {
	case LinearIssueCreatePayload, LinearCommentCreatePayload:
		return c.invokeWrite(ctx, inv, payload, kind, id, backend)
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

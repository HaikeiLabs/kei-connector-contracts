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
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

type CRMLead struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Stage     string    `json:"stage"`
	CreatedAt time.Time `json:"created_at"`
}

// LeadReadPayload requests lead.read. An empty ID lists all leads in scope;
// a non-empty ID must agree with the invocation resource.
type LeadReadPayload struct {
	ID string `json:"id,omitempty"`
}

func (LeadReadPayload) Capability() string { return "lead.read" }

// CRMBackend is the seam for CRM reads. authRef is the connector's opaque
// credential reference; backends resolve it against the credential service
// and must reject calls without a valid authenticated session.
type CRMBackend interface {
	ListLeads(ctx context.Context, authRef string) ([]CRMLead, error)
	Lead(ctx context.Context, authRef, id string) (CRMLead, error)
}

// CRMInvestor is one fundraising investor. The fields are the complete
// allowlist a CRM may return for an investor: there is no free text.
type CRMInvestor struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Firm      string    `json:"firm"`
	Email     string    `json:"email"`
	Stage     string    `json:"stage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InvestorStages are the fundraising pipeline stages, in order. They mirror
// kei-agents InvestorStage.
var InvestorStages = []string{"prospect", "contacted", "meeting", "diligence", "committed", "passed"}

// Investor list bounds.
const (
	MaxInvestorQuery     = 200
	MaxInvestorLimit     = 100
	DefaultInvestorLimit = 20
	MaxInvestorCursor    = 128
)

// InvestorListPayload requests investor.list on the investors collection.
// Every field is optional: Stage filters by pipeline stage, Query matches
// name, firm or email literally, Limit bounds the page (default 20) and
// Cursor is the opaque next_cursor of a previous page.
type InvestorListPayload struct {
	Stage  string `json:"stage,omitempty"`
	Query  string `json:"q,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

func (InvestorListPayload) Capability() string { return "investor.list" }

// Validate checks the payload against the investor.list bounds.
func (p InvestorListPayload) Validate() error {
	switch {
	case p.Stage != "" && !slices.Contains(InvestorStages, p.Stage):
		return errors.New("unknown investor stage")
	case len(p.Query) > MaxInvestorQuery:
		return errors.New("investor query is too long")
	case p.Limit < 0 || p.Limit > MaxInvestorLimit:
		return errors.New("investor limit is out of range")
	case len(p.Cursor) > MaxInvestorCursor:
		return errors.New("investor cursor is too long")
	}
	return nil
}

// InvestorReadPayload requests investor.read. It carries nothing: the
// resource investors/<id> names the investor.
type InvestorReadPayload struct{}

func (InvestorReadPayload) Capability() string { return "investor.read" }

// InvestorPage is the investor.list result. NextCursor is null on the last
// page.
type InvestorPage struct {
	Records    []CRMInvestor `json:"records"`
	NextCursor *string       `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

// CRMInvestorBackend is the seam for investor reads. It is separate from
// CRMBackend so existing lead backends keep compiling; a backend without it
// fails investor reads closed.
type CRMInvestorBackend interface {
	ListInvestors(ctx context.Context, authRef string, filter InvestorListPayload) (InvestorPage, error)
	Investor(ctx context.Context, authRef, id string) (CRMInvestor, error)
}

// MemoryCRM is the in-memory CRM backend used until a real one exists. When
// AuthRef is set, only that opaque reference is accepted, simulating an
// authenticated session.
type MemoryCRM struct {
	AuthRef   string
	Leads     map[string]CRMLead
	Investors map[string]CRMInvestor
}

var errUnauthenticated = errors.New("crm session is not authenticated")

func (m MemoryCRM) ListLeads(_ context.Context, authRef string) ([]CRMLead, error) {
	if m.AuthRef != "" && authRef != m.AuthRef {
		return nil, errUnauthenticated
	}
	var out []CRMLead
	for _, l := range m.Leads {
		out = append(out, l)
	}
	return out, nil
}

func (m MemoryCRM) Lead(_ context.Context, authRef, id string) (CRMLead, error) {
	if m.AuthRef != "" && authRef != m.AuthRef {
		return CRMLead{}, errUnauthenticated
	}
	l, ok := m.Leads[id]
	if !ok {
		return CRMLead{}, errors.New("lead not found")
	}
	return l, nil
}

// ListInvestors pages investors in descending id order; the cursor is the
// last id of the previous page.
func (m MemoryCRM) ListInvestors(_ context.Context, authRef string, f InvestorListPayload) (InvestorPage, error) {
	if m.AuthRef != "" && authRef != m.AuthRef {
		return InvestorPage{}, errUnauthenticated
	}
	limit := f.Limit
	if limit == 0 {
		limit = DefaultInvestorLimit
	}
	query := strings.ToLower(f.Query)
	var matched []CRMInvestor
	for _, i := range m.Investors {
		if f.Stage != "" && i.Stage != f.Stage {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(i.Name+"\x00"+i.Firm+"\x00"+i.Email), query) {
			continue
		}
		if f.Cursor != "" && i.ID >= f.Cursor {
			continue
		}
		matched = append(matched, i)
	}
	sort.Slice(matched, func(a, b int) bool { return matched[a].ID > matched[b].ID })
	page := InvestorPage{Records: matched}
	if len(matched) > limit {
		page.Records, page.HasMore = matched[:limit], true
		next := matched[limit-1].ID
		page.NextCursor = &next
	}
	if page.Records == nil {
		page.Records = []CRMInvestor{}
	}
	return page, nil
}

func (m MemoryCRM) Investor(_ context.Context, authRef, id string) (CRMInvestor, error) {
	if m.AuthRef != "" && authRef != m.AuthRef {
		return CRMInvestor{}, errUnauthenticated
	}
	i, ok := m.Investors[id]
	if !ok {
		return CRMInvestor{}, errors.New("investor not found")
	}
	return i, nil
}

// CRMClient serves lead.read, investor.list and investor.read. Lead
// resources are leads (list) or leads/<id> (single lead); investor resources
// are investors (investor.list) and investors/<id> (investor.read). Customer
// reads are not exposed because the contract defines no customer
// capabilities.
type CRMClient struct {
	backend CRMBackend
}

func NewCRM(backend CRMBackend) *CRMClient { return &CRMClient{backend: backend} }

func (c *CRMClient) Provider() contract.Provider { return contract.ProviderCRM }

func (c *CRMClient) Invoke(ctx context.Context, meta contract.Metadata, inv contract.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	switch p := payload.(type) {
	case InvestorListPayload:
		return c.listInvestors(ctx, meta, inv, p)
	case InvestorReadPayload:
		return c.readInvestor(ctx, meta, inv)
	}
	p, ok := payload.(LeadReadPayload)
	if !ok {
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
	id := ""
	if inv.Resource != "leads" {
		parts := strings.Split(inv.Resource, "/")
		if len(parts) != 2 || parts[0] != "leads" || parts[1] == "" {
			return Result{}, errors.New("resource must be leads or leads/<id>")
		}
		id = parts[1]
	}
	if p.ID != id {
		return Result{}, errors.New("payload id does not match the invocation resource")
	}
	if id == "" {
		leads, err := c.backend.ListLeads(ctx, meta.CredentialRef)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: leads}, nil
	}
	lead, err := c.backend.Lead(ctx, meta.CredentialRef, id)
	if err != nil {
		return Result{}, err
	}
	return Result{Capability: inv.Capability, Data: lead}, nil
}

var errNoInvestorBackend = errors.New("crm backend does not serve investors")

func (c *CRMClient) listInvestors(ctx context.Context, meta contract.Metadata, inv contract.Invocation, p InvestorListPayload) (Result, error) {
	if inv.Resource != "investors" {
		return Result{}, errors.New("investor.list resource must be investors")
	}
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	backend, ok := c.backend.(CRMInvestorBackend)
	if !ok {
		return Result{}, errNoInvestorBackend
	}
	page, err := backend.ListInvestors(ctx, meta.CredentialRef, p)
	if err != nil {
		return Result{}, err
	}
	return Result{Capability: inv.Capability, Data: page}, nil
}

func (c *CRMClient) readInvestor(ctx context.Context, meta contract.Metadata, inv contract.Invocation) (Result, error) {
	id, ok := strings.CutPrefix(inv.Resource, "investors/")
	if !ok || id == "" || strings.Contains(id, "/") {
		return Result{}, errors.New("investor.read resource must be investors/<id>")
	}
	backend, ok := c.backend.(CRMInvestorBackend)
	if !ok {
		return Result{}, errNoInvestorBackend
	}
	investor, err := backend.Investor(ctx, meta.CredentialRef, id)
	if err != nil {
		return Result{}, err
	}
	return Result{Capability: inv.Capability, Data: investor}, nil
}

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
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
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

// MemoryCRM is the in-memory CRM backend used until a real one exists. When
// AuthRef is set, only that opaque reference is accepted, simulating an
// authenticated session.
type MemoryCRM struct {
	AuthRef string
	Leads   map[string]CRMLead
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

// CRMClient serves lead.read. Resources are leads (list) or leads/<id>
// (single lead). Customer reads are not exposed because the contract defines
// no customer capabilities.
type CRMClient struct {
	backend CRMBackend
}

func NewCRM(backend CRMBackend) *CRMClient { return &CRMClient{backend: backend} }

func (c *CRMClient) Provider() connectors.Provider { return connectors.ProviderCRM }

func (c *CRMClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
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

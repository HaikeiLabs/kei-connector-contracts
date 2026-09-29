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
	"regexp"
	"sort"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// TitoEvent is the read-only event record exposed by the client. Status is
// "live" or "draft".
type TitoEvent struct {
	Account   string `json:"account_slug"`
	Slug      string `json:"event_slug"`
	Title     string `json:"title"`
	URL       string `json:"url,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Status    string `json:"status"`
}

// TitoEventPage is the event.list result.
type TitoEventPage struct {
	Events   []TitoEvent `json:"events"`
	NextPage int         `json:"next_page,omitempty"`
}

// TitoRelease is a ticket type on an event. It carries counts, never
// attendees.
type TitoRelease struct {
	ID           int64  `json:"release_id"`
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Quantity     int    `json:"quantity,omitempty"`
	TicketsCount int    `json:"tickets_count"`
	Price        string `json:"price,omitempty"`
}

// TitoTicket is the only per-ticket data a backend reads: state and release.
// Attendee names, emails, and answers are never decoded.
type TitoTicket struct {
	State        string
	ReleaseTitle string
}

// TitoTicketSummary is the ticket.summary result: aggregate counts only.
type TitoTicketSummary struct {
	Account   string         `json:"account_slug"`
	Event     string         `json:"event_slug"`
	Total     int            `json:"total"`
	ByState   map[string]int `json:"by_state"`
	ByRelease map[string]int `json:"by_release"`
}

func summarizeTickets(account, event string, tickets []TitoTicket) TitoTicketSummary {
	summary := TitoTicketSummary{Account: account, Event: event, ByState: map[string]int{}, ByRelease: map[string]int{}}
	for _, ticket := range tickets {
		summary.Total++
		summary.ByState[ticket.State]++
		summary.ByRelease[ticket.ReleaseTitle]++
	}
	return summary
}

// EventListPayload requests event.list over tito/<account>/events. Status is
// upcoming (default), past, or archived; Page starts at 1.
type EventListPayload struct {
	Status string `json:"status,omitempty"`
	Page   int    `json:"page,omitempty"`
}

func (EventListPayload) Capability() string { return "event.list" }

// EventGetPayload requests event.get over tito/<account>/events/<event>.
type EventGetPayload struct{}

func (EventGetPayload) Capability() string { return "event.get" }

// ReleaseListPayload requests release.list over
// tito/<account>/events/<event>/releases.
type ReleaseListPayload struct{}

func (ReleaseListPayload) Capability() string { return "release.list" }

// TicketSummaryPayload requests ticket.summary over
// tito/<account>/events/<event>/ticket-summary.
type TicketSummaryPayload struct{}

func (TicketSummaryPayload) Capability() string { return "ticket.summary" }

// TitoBackend is the seam for the Tito admin API. Implementations receive
// validated account and event slugs only.
type TitoBackend interface {
	Events(ctx context.Context, account, status string, page int) ([]TitoEvent, int, error)
	Event(ctx context.Context, account, event string) (TitoEvent, error)
	Releases(ctx context.Context, account, event string) ([]TitoRelease, error)
	Tickets(ctx context.Context, account, event string) ([]TitoTicket, error)
}

// MemoryTito is the in-memory Tito backend used by tests and fixtures. Maps
// are keyed by "<account>/<event>".
type MemoryTito struct {
	Events   map[string]TitoEvent
	Releases map[string][]TitoRelease
	Tickets  map[string][]TitoTicket
}

// memoryTitoBackend adapts MemoryTito to TitoBackend; the adapter exists
// because MemoryTito's fields share the backend's method names.
type memoryTitoBackend struct{ MemoryTito }

func (m memoryTitoBackend) Events(_ context.Context, account, _ string, _ int) ([]TitoEvent, int, error) {
	keys := make([]string, 0, len(m.MemoryTito.Events))
	for key := range m.MemoryTito.Events {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var out []TitoEvent
	for _, key := range keys {
		if strings.HasPrefix(key, account+"/") {
			out = append(out, m.MemoryTito.Events[key])
		}
	}
	return out, 0, nil
}

func (m memoryTitoBackend) Event(_ context.Context, account, event string) (TitoEvent, error) {
	e, ok := m.MemoryTito.Events[account+"/"+event]
	if !ok {
		return TitoEvent{}, errors.New("event not found")
	}
	return e, nil
}

func (m memoryTitoBackend) Releases(_ context.Context, account, event string) ([]TitoRelease, error) {
	return m.MemoryTito.Releases[account+"/"+event], nil
}

func (m memoryTitoBackend) Tickets(_ context.Context, account, event string) ([]TitoTicket, error) {
	return m.MemoryTito.Tickets[account+"/"+event], nil
}

// TitoClient serves event.list, event.get, release.list, and ticket.summary.
// Tickets are only ever aggregated: no capability returns attendees.
type TitoClient struct {
	backend TitoBackend
}

// NewTito returns a client over the in-memory backend.
func NewTito(store MemoryTito) *TitoClient { return NewTitoBackend(memoryTitoBackend{store}) }

func NewTitoBackend(backend TitoBackend) *TitoClient { return &TitoClient{backend: backend} }

func (c *TitoClient) Provider() contract.Provider { return contract.ProviderTito }

func (c *TitoClient) Invoke(ctx context.Context, meta contract.Metadata, inv contract.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	ctx = withInvocationScope(ctx, meta, inv)
	account, event, leaf, ok := titoResource(inv.Resource)
	if !ok {
		return Result{}, errors.New("resource must be tito/<account>/events[/<event>[/releases|/ticket-summary]]")
	}
	switch p := payload.(type) {
	case EventListPayload:
		if event != "" {
			return Result{}, errors.New("resource must be tito/<account>/events")
		}
		status := p.Status
		if status == "" {
			status = "upcoming"
		}
		if status != "upcoming" && status != "past" && status != "archived" {
			return Result{}, errors.New("status must be upcoming, past, or archived")
		}
		if p.Page < 0 {
			return Result{}, errors.New("page must be positive")
		}
		page := p.Page
		if page == 0 {
			page = 1
		}
		events, next, err := c.backend.Events(ctx, account, status, page)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: TitoEventPage{Events: events, NextPage: next}}, nil
	case EventGetPayload:
		if event == "" || leaf != "" {
			return Result{}, errors.New("resource must be tito/<account>/events/<event>")
		}
		e, err := c.backend.Event(ctx, account, event)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: e}, nil
	case ReleaseListPayload:
		if event == "" || leaf != "releases" {
			return Result{}, errors.New("resource must be tito/<account>/events/<event>/releases")
		}
		releases, err := c.backend.Releases(ctx, account, event)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: releases}, nil
	case TicketSummaryPayload:
		if event == "" || leaf != "ticket-summary" {
			return Result{}, errors.New("resource must be tito/<account>/events/<event>/ticket-summary")
		}
		tickets, err := c.backend.Tickets(ctx, account, event)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: summarizeTickets(account, event, tickets)}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

var titoSlug = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$`)

// titoResource parses tito/<account>/events[/<event>[/<leaf>]]. Slugs are
// restricted so a resource can never reshape the fixed API path.
func titoResource(resource string) (account, event, leaf string, ok bool) {
	parts := strings.Split(resource, "/")
	if len(parts) < 3 || len(parts) > 5 || parts[0] != "tito" || parts[2] != "events" || !titoSlug.MatchString(parts[1]) {
		return "", "", "", false
	}
	account = parts[1]
	if len(parts) >= 4 {
		if !titoSlug.MatchString(parts[3]) {
			return "", "", "", false
		}
		event = parts[3]
	}
	if len(parts) == 5 {
		leaf = parts[4]
		if leaf != "releases" && leaf != "ticket-summary" {
			return "", "", "", false
		}
	}
	return account, event, leaf, true
}

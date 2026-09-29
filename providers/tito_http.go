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
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	titoAPIBase = "https://api.tito.io/v3"
	// titoTicketPageSize and titoMaxTicketPages bound ticket.summary. Past
	// the bound the summary fails closed instead of returning partial counts.
	titoTicketPageSize = 100
	titoMaxTicketPages = 50
)

// AuthenticatedTito is the production Tito admin API backend. It calls only
// GET endpoints on a fixed base; the API token is resolved per invocation
// from the connector's opaque credential reference.
type AuthenticatedTito struct {
	httpClient  *http.Client
	credentials CredentialResolver
}

func NewAuthenticatedTito(config RuntimeConfig) *TitoClient {
	config = normalizeRuntimeConfig(config)
	return NewTitoBackend(&AuthenticatedTito{httpClient: config.HTTPClient, credentials: config.Credentials})
}

func NewTitoHTTP(httpClient *http.Client, credentials CredentialResolver) *TitoClient {
	return NewAuthenticatedTito(RuntimeConfig{HTTPClient: httpClient, Credentials: credentials})
}

func (b *AuthenticatedTito) Events(ctx context.Context, account, status string, page int) ([]TitoEvent, int, error) {
	endpoint := titoAPIBase + "/" + url.PathEscape(account) + "/events"
	if status != "upcoming" {
		endpoint += "/" + url.PathEscape(status)
	}
	var response struct {
		Events []titoEventResponse `json:"events"`
		Meta   titoPageMeta        `json:"meta"`
	}
	if err := b.getJSON(ctx, endpoint+"?"+url.Values{"page[number]": {strconv.Itoa(page)}}.Encode(), &response); err != nil {
		return nil, 0, err
	}
	events := make([]TitoEvent, 0, len(response.Events))
	for _, event := range response.Events {
		events = append(events, event.model(account))
	}
	return events, response.Meta.NextPage, nil
}

func (b *AuthenticatedTito) Event(ctx context.Context, account, event string) (TitoEvent, error) {
	var response struct {
		Event titoEventResponse `json:"event"`
	}
	if err := b.getJSON(ctx, titoAPIBase+"/"+url.PathEscape(account)+"/"+url.PathEscape(event), &response); err != nil {
		return TitoEvent{}, err
	}
	return response.Event.model(account), nil
}

func (b *AuthenticatedTito) Releases(ctx context.Context, account, event string) ([]TitoRelease, error) {
	var response struct {
		Releases []struct {
			ID           int64       `json:"id"`
			Slug         string      `json:"slug"`
			Title        string      `json:"title"`
			Quantity     int         `json:"quantity"`
			TicketsCount int         `json:"tickets_count"`
			Price        json.Number `json:"price"`
		} `json:"releases"`
	}
	if err := b.getJSON(ctx, titoAPIBase+"/"+url.PathEscape(account)+"/"+url.PathEscape(event)+"/releases", &response); err != nil {
		return nil, err
	}
	releases := make([]TitoRelease, 0, len(response.Releases))
	for _, r := range response.Releases {
		releases = append(releases, TitoRelease{ID: r.ID, Slug: r.Slug, Title: r.Title, Quantity: r.Quantity, TicketsCount: r.TicketsCount, Price: r.Price.String()})
	}
	return releases, nil
}

// Tickets decodes only state and release_title from each ticket, so attendee
// PII in the provider response is discarded at the decode boundary.
func (b *AuthenticatedTito) Tickets(ctx context.Context, account, event string) ([]TitoTicket, error) {
	endpoint := titoAPIBase + "/" + url.PathEscape(account) + "/" + url.PathEscape(event) + "/tickets"
	var tickets []TitoTicket
	for page := 1; ; {
		if page > titoMaxTicketPages {
			return nil, errors.New("too many tickets to summarize")
		}
		var response struct {
			Tickets []struct {
				State        string `json:"state"`
				ReleaseTitle string `json:"release_title"`
			} `json:"tickets"`
			Meta titoPageMeta `json:"meta"`
		}
		params := url.Values{"page[number]": {strconv.Itoa(page)}, "page[size]": {strconv.Itoa(titoTicketPageSize)}}
		if err := b.getJSON(ctx, endpoint+"?"+params.Encode(), &response); err != nil {
			return nil, err
		}
		for _, t := range response.Tickets {
			tickets = append(tickets, TitoTicket{State: t.State, ReleaseTitle: t.ReleaseTitle})
		}
		if response.Meta.NextPage == 0 {
			return tickets, nil
		}
		if response.Meta.NextPage <= page {
			return nil, errors.New("provider response was invalid")
		}
		page = response.Meta.NextPage
	}
}

type titoPageMeta struct {
	NextPage int `json:"next_page"`
}

type titoEventResponse struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Live      bool   `json:"live"`
}

func (e titoEventResponse) model(account string) TitoEvent {
	status := "draft"
	if e.Live {
		status = "live"
	}
	return TitoEvent{Account: account, Slug: e.Slug, Title: e.Title, URL: e.URL, StartDate: e.StartDate, EndDate: e.EndDate, Status: status}
}

func (b *AuthenticatedTito) getJSON(ctx context.Context, endpoint string, out any) error {
	scope, err := scopeFromContext(ctx)
	if err != nil {
		return err
	}
	token, err := b.credentials.Resolve(ctx, scope.tenantID, scope.workspaceID, scope.subject, scope.credential)
	if err != nil || token == "" {
		return errors.New("credential resolution failed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return errors.New("provider request could not be created")
	}
	req.Header.Set("Authorization", "Token token="+token)
	req.Header.Set("Accept", "application/json")
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return errors.New("provider request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("provider request was rejected")
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errors.New("provider response was invalid")
	}
	return nil
}

var _ TitoBackend = (*AuthenticatedTito)(nil)

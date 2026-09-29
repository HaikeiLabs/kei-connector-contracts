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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// Attendee PII markers planted in fixtures; no result may carry them.
var titoPIIMarkers = []string{"FAKE-ATTENDEE", "attendee@example.test", `"email"`, `"name"`}

func titoStore() MemoryTito {
	return MemoryTito{
		Events: map[string]TitoEvent{
			"acme/conf-2026": {Account: "acme", Slug: "conf-2026", Title: "Conf 2026", URL: "https://ti.to/acme/conf-2026", StartDate: "2026-10-01", EndDate: "2026-10-02", Status: "live"},
		},
		Releases: map[string][]TitoRelease{
			"acme/conf-2026": {{ID: 11, Slug: "ga", Title: "General Admission", Quantity: 100, TicketsCount: 3}},
		},
		Tickets: map[string][]TitoTicket{
			"acme/conf-2026": {
				{State: "complete", ReleaseTitle: "General Admission"},
				{State: "complete", ReleaseTitle: "General Admission"},
				{State: "incomplete", ReleaseTitle: "Speaker"},
			},
		},
	}
}

func titoMeta(t *testing.T) connectors.Metadata {
	t.Helper()
	return metaFor(t, connectors.ProviderTito, []string{"tito/acme/events"}, nil)
}

func TestTitoReadsReturnEventReleaseAndAggregateData(t *testing.T) {
	client := NewTito(titoStore())
	meta := titoMeta(t)
	ctx := context.Background()

	result, err := client.Invoke(ctx, meta, invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{})
	if err != nil {
		t.Fatalf("event.list: %v", err)
	}
	if page := result.Data.(TitoEventPage); len(page.Events) != 1 || page.Events[0].Slug != "conf-2026" || page.Events[0].Title != "Conf 2026" {
		t.Fatalf("events = %+v", page)
	}

	result, err = client.Invoke(ctx, meta, invocation("event.get", connectors.ActionRead, "tito/acme/events/conf-2026"), EventGetPayload{})
	if err != nil {
		t.Fatalf("event.get: %v", err)
	}
	if event := result.Data.(TitoEvent); event.URL != "https://ti.to/acme/conf-2026" || event.StartDate != "2026-10-01" || event.Status != "live" {
		t.Fatalf("event = %+v", event)
	}

	result, err = client.Invoke(ctx, meta, invocation("release.list", connectors.ActionRead, "tito/acme/events/conf-2026/releases"), ReleaseListPayload{})
	if err != nil {
		t.Fatalf("release.list: %v", err)
	}
	if releases := result.Data.([]TitoRelease); len(releases) != 1 || releases[0].Title != "General Admission" || releases[0].TicketsCount != 3 {
		t.Fatalf("releases = %+v", releases)
	}

	result, err = client.Invoke(ctx, meta, invocation("ticket.summary", connectors.ActionRead, "tito/acme/events/conf-2026/ticket-summary"), TicketSummaryPayload{})
	if err != nil {
		t.Fatalf("ticket.summary: %v", err)
	}
	summary := result.Data.(TitoTicketSummary)
	if summary.Total != 3 || summary.ByState["complete"] != 2 || summary.ByState["incomplete"] != 1 || summary.ByRelease["General Admission"] != 2 || summary.ByRelease["Speaker"] != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	assertNoTitoPII(t, result)
}

func TestTitoRejectsWritesAndMismatchedRequests(t *testing.T) {
	client := NewTito(titoStore())
	meta := titoMeta(t)
	for name, tc := range map[string]struct {
		inv     connectors.Invocation
		payload Payload
	}{
		"ticket create is not a capability": {invocation("ticket.create", connectors.ActionCreate, "tito/acme/events/conf-2026"), EventGetPayload{}},
		"read capability as a write":        {invocation("event.get", connectors.ActionUpdate, "tito/acme/events/conf-2026"), EventGetPayload{}},
		"delete is locked":                  {invocation("event.get", connectors.ActionDelete, "tito/acme/events/conf-2026"), EventGetPayload{}},
		"payload replay":                    {invocation("event.get", connectors.ActionRead, "tito/acme/events/conf-2026"), TicketSummaryPayload{}},
		"list on an event resource":         {invocation("event.list", connectors.ActionRead, "tito/acme/events/conf-2026"), EventListPayload{}},
		"get on the collection":             {invocation("event.get", connectors.ActionRead, "tito/acme/events"), EventGetPayload{}},
		"summary on the event":              {invocation("ticket.summary", connectors.ActionRead, "tito/acme/events/conf-2026"), TicketSummaryPayload{}},
		"releases without an event":         {invocation("release.list", connectors.ActionRead, "tito/acme/events//releases"), ReleaseListPayload{}},
		"raw ticket listing":                {invocation("ticket.summary", connectors.ActionRead, "tito/acme/events/conf-2026/tickets"), TicketSummaryPayload{}},
		"slug with a query":                 {invocation("event.get", connectors.ActionRead, "tito/acme/events/conf?x=1"), EventGetPayload{}},
		"encoded slash in a slug":           {invocation("event.get", connectors.ActionRead, "tito/acme/events/conf%2Freg"), EventGetPayload{}},
		"traversal":                         {invocation("event.get", connectors.ActionRead, "tito/acme/events/../conf-2026"), EventGetPayload{}},
		"unknown event status filter":       {invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{Status: "deleted"}},
		"negative page":                     {invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{Page: -1}},
	} {
		if _, err := client.Invoke(context.Background(), meta, tc.inv, tc.payload); err == nil {
			t.Errorf("%s: invocation was allowed", name)
		}
	}
}

func TestAuthenticatedTitoContract(t *testing.T) {
	resolver := &recordingResolver{token: "fake-tito-token"}
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Token token=fake-tito-token" {
			t.Errorf("authorization = %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("accept = %q", got)
		}
		paths = append(paths, r.URL.Path+"?"+r.URL.RawQuery)
		switch r.URL.Path {
		case "/v3/acme/events", "/v3/acme/events/past":
			_, _ = io.WriteString(w, `{"events":[{"slug":"conf-2026","title":"Conf 2026","url":"https://ti.to/acme/conf-2026","start_date":"2026-10-01","end_date":"2026-10-02","live":true,"test_mode":false}],"meta":{"next_page":2}}`)
		case "/v3/acme/conf-2026":
			_, _ = io.WriteString(w, `{"event":{"slug":"conf-2026","title":"Conf 2026","url":"https://ti.to/acme/conf-2026","start_date":"2026-10-01","end_date":"2026-10-02","live":false,"test_mode":false}}`)
		case "/v3/acme/conf-2026/releases":
			_, _ = io.WriteString(w, `{"releases":[{"id":11,"slug":"ga","title":"General Admission","quantity":100,"tickets_count":3,"price":"50.0","secret":false}]}`)
		case "/v3/acme/conf-2026/tickets":
			if r.URL.Query().Get("page[number]") == "2" {
				_, _ = io.WriteString(w, `{"tickets":[{"state":"incomplete","release_title":"Speaker","name":"FAKE-ATTENDEE Three","first_name":"FAKE-ATTENDEE","email":"attendee@example.test","phone_number":"FAKE-ATTENDEE-PHONE"}],"meta":{"next_page":null}}`)
				return
			}
			_, _ = io.WriteString(w, `{"tickets":[{"state":"complete","release_title":"General Admission","name":"FAKE-ATTENDEE One","email":"attendee@example.test"},{"state":"complete","release_title":"General Admission","name":"FAKE-ATTENDEE Two","email":"attendee@example.test"}],"meta":{"next_page":2}}`)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := NewTitoHTTP(testHTTPClient(server, func(_ *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "api.tito.io" || req.Method != http.MethodGet {
			t.Errorf("request = %s %s", req.Method, req.URL.String())
		}
	}), resolver)
	meta := titoMeta(t)
	ctx := context.Background()

	result, err := client.Invoke(ctx, meta, invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{Status: "past"})
	if err != nil {
		t.Fatalf("event.list: %v", err)
	}
	page := result.Data.(TitoEventPage)
	if len(page.Events) != 1 || page.Events[0].Account != "acme" || page.Events[0].Status != "live" || page.NextPage != 2 {
		t.Fatalf("events = %+v", page)
	}
	result, err = client.Invoke(ctx, meta, invocation("event.get", connectors.ActionRead, "tito/acme/events/conf-2026"), EventGetPayload{})
	if err != nil || result.Data.(TitoEvent).Status != "draft" {
		t.Fatalf("event.get = %+v, %v", result, err)
	}
	result, err = client.Invoke(ctx, meta, invocation("release.list", connectors.ActionRead, "tito/acme/events/conf-2026/releases"), ReleaseListPayload{})
	if err != nil || len(result.Data.([]TitoRelease)) != 1 || result.Data.([]TitoRelease)[0].Price != "50.0" {
		t.Fatalf("release.list = %+v, %v", result, err)
	}
	result, err = client.Invoke(ctx, meta, invocation("ticket.summary", connectors.ActionRead, "tito/acme/events/conf-2026/ticket-summary"), TicketSummaryPayload{})
	if err != nil {
		t.Fatalf("ticket.summary: %v", err)
	}
	summary := result.Data.(TitoTicketSummary)
	if summary.Total != 3 || summary.ByState["complete"] != 2 || summary.ByRelease["Speaker"] != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	assertNoTitoPII(t, result)

	want := []string{
		"/v3/acme/events/past?page%5Bnumber%5D=1",
		"/v3/acme/conf-2026?",
		"/v3/acme/conf-2026/releases?",
		"/v3/acme/conf-2026/tickets?page%5Bnumber%5D=1&page%5Bsize%5D=100",
		"/v3/acme/conf-2026/tickets?page%5Bnumber%5D=2&page%5Bsize%5D=100",
	}
	if strings.Join(paths, " ") != strings.Join(want, " ") {
		t.Fatalf("requests = %v\nwant %v", paths, want)
	}
	if resolver.ten != "t-1" || resolver.workspace != "w-1" || resolver.subject != "u-1" || resolver.ref != meta.CredentialRef {
		t.Fatalf("scope = %+v", resolver)
	}
}

func TestAuthenticatedTitoTicketSummaryFailsClosedPastPageLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page[number]"))
		_, _ = fmt.Fprintf(w, `{"tickets":[{"state":"complete","release_title":"GA"}],"meta":{"next_page":%d}}`, page+1)
	}))
	defer server.Close()
	client := NewTitoHTTP(testHTTPClient(server, func(*testing.T, *http.Request) {}), &recordingResolver{token: "fake-tito-token"})
	_, err := client.Invoke(context.Background(), titoMeta(t), invocation("ticket.summary", connectors.ActionRead, "tito/acme/events/conf-2026/ticket-summary"), TicketSummaryPayload{})
	if err == nil || !strings.Contains(err.Error(), "too many tickets") {
		t.Fatalf("error = %v, want a fail-closed page limit", err)
	}
}

func TestAuthenticatedTitoFailsClosedAndRedactsErrors(t *testing.T) {
	client, err := New(connectors.ProviderTito)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Invoke(context.Background(), titoMeta(t), invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{})
	if err == nil || !strings.Contains(err.Error(), "credential") {
		t.Fatalf("unconfigured resolver error = %v", err)
	}

	noNetwork := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("must not send") })}
	_, err = NewTitoHTTP(noNetwork, credentialErrorResolver{}).Invoke(context.Background(), titoMeta(t), invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{})
	if err == nil || err.Error() != "credential resolution failed" {
		t.Fatalf("resolver error = %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"FAKE-ATTENDEE fake-tito-token"}`, http.StatusUnauthorized)
	}))
	defer server.Close()
	_, err = NewTitoHTTP(testHTTPClient(server, func(*testing.T, *http.Request) {}), &recordingResolver{token: "fake-tito-token"}).Invoke(context.Background(), titoMeta(t), invocation("event.list", connectors.ActionRead, "tito/acme/events"), EventListPayload{})
	if err == nil || strings.Contains(err.Error(), "fake-tito-token") || strings.Contains(err.Error(), "FAKE-ATTENDEE") {
		t.Fatalf("provider rejection error = %v", err)
	}
}

func assertNoTitoPII(t *testing.T, result Result) {
	t.Helper()
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range titoPIIMarkers {
		if strings.Contains(string(encoded), marker) {
			t.Fatalf("result carries attendee PII %q: %s", marker, encoded)
		}
	}
}

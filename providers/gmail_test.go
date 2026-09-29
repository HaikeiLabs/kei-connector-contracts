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
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

const gmailFixtureBody = "FAKE-BODY-MARKER quarterly numbers inside"

func gmailStore() MemoryGmail {
	return MemoryGmail{Messages: map[string]GmailMessage{
		"m-1": {ID: "m-1", ThreadID: "th-1", InternalDate: time.Unix(1700000000, 0).UTC(), From: "a@example.test", To: "b@example.test", Subject: "Quarterly update", Labels: []string{"INBOX"}, Snippet: "numbers inside", Body: gmailFixtureBody},
		"m-2": {ID: "m-2", ThreadID: "th-2", Subject: "Lunch", Snippet: "tacos?", Body: "FAKE-BODY-MARKER tacos"},
	}}
}

func gmailMeta(t *testing.T, includeBody bool) connectors.Metadata {
	t.Helper()
	m := metaFor(t, connectors.ProviderGmail, []string{"gmail/messages"}, nil)
	m.Policy.GmailIncludeBody = includeBody
	return m
}

func TestGmailReadonlyScopeIsPublished(t *testing.T) {
	if GmailReadonlyScope != "https://www.googleapis.com/auth/gmail.readonly" {
		t.Fatalf("GmailReadonlyScope = %q", GmailReadonlyScope)
	}
}

func TestGmailSearchReturnsMetadataAndSnippetOnly(t *testing.T) {
	client := NewGmail(gmailStore())
	result, err := client.Invoke(context.Background(), gmailMeta(t, true), invocation("message.search", connectors.ActionRead, "gmail/messages"), MessageSearchPayload{Query: "quarterly"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	page := result.Data.(GmailMessagePage)
	if len(page.Messages) != 1 || page.Messages[0].ID != "m-1" || page.Messages[0].Snippet != "numbers inside" || page.Messages[0].Subject != "Quarterly update" {
		t.Fatalf("page = %+v", page)
	}
	// Even with the body opt-in, search never returns bodies.
	assertNoGmailBody(t, result)
}

func TestGmailGetExcludesBodyByDefault(t *testing.T) {
	client := NewGmail(gmailStore())
	result, err := client.Invoke(context.Background(), gmailMeta(t, false), invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{MessageID: "m-1"})
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	msg := result.Data.(GmailMessage)
	if msg.ID != "m-1" || msg.ThreadID != "th-1" || msg.From != "a@example.test" || msg.Snippet != "numbers inside" || len(msg.Labels) != 1 {
		t.Fatalf("message = %+v", msg)
	}
	assertNoGmailBody(t, result)
}

func TestGmailGetIncludesBodyOnlyWithPolicyOptIn(t *testing.T) {
	client := NewGmail(gmailStore())
	result, err := client.Invoke(context.Background(), gmailMeta(t, true), invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{})
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got := result.Data.(GmailMessage).Body; got != gmailFixtureBody {
		t.Fatalf("body = %q, want the fixture body", got)
	}
}

// A backend that ignores the includeBody flag must still not leak a body.
func TestGmailClientStripsBodyEvenIfBackendReturnsIt(t *testing.T) {
	client := NewGmail(leakyGmail{})
	for _, tc := range []struct {
		capability, resource string
		payload              Payload
	}{
		{"message.search", "gmail/messages", MessageSearchPayload{}},
		{"message.get", "gmail/messages/m-1", MessageGetPayload{}},
	} {
		result, err := client.Invoke(context.Background(), gmailMeta(t, false), invocation(tc.capability, connectors.ActionRead, tc.resource), tc.payload)
		if err != nil {
			t.Fatalf("%s: %v", tc.capability, err)
		}
		assertNoGmailBody(t, result)
	}
}

func TestGmailRejectsWritesAndMismatchedRequests(t *testing.T) {
	client := NewGmail(gmailStore())
	meta := gmailMeta(t, false)
	for name, tc := range map[string]struct {
		inv     connectors.Invocation
		payload Payload
	}{
		"send is not a capability":     {invocation("message.send", connectors.ActionCreate, "gmail/messages"), MessageSearchPayload{}},
		"read capability as a write":   {invocation("message.get", connectors.ActionUpdate, "gmail/messages/m-1"), MessageGetPayload{}},
		"delete is locked":             {invocation("message.get", connectors.ActionDelete, "gmail/messages/m-1"), MessageGetPayload{}},
		"payload replay":               {invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageSearchPayload{}},
		"search on a message resource": {invocation("message.search", connectors.ActionRead, "gmail/messages/m-1"), MessageSearchPayload{}},
		"get on the collection":        {invocation("message.get", connectors.ActionRead, "gmail/messages"), MessageGetPayload{}},
		"payload id differs":           {invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{MessageID: "m-2"}},
		"foreign resource":             {invocation("message.get", connectors.ActionRead, "drive/d-1/files/m-1"), MessageGetPayload{}},
		"traversal":                    {invocation("message.get", connectors.ActionRead, "gmail/messages/../m-1"), MessageGetPayload{}},
		"negative page size":           {invocation("message.search", connectors.ActionRead, "gmail/messages"), MessageSearchPayload{PageSize: -1}},
		"page size above the hard cap": {invocation("message.search", connectors.ActionRead, "gmail/messages"), MessageSearchPayload{PageSize: gmailMaxPageSize + 1}},
	} {
		if _, err := client.Invoke(context.Background(), meta, tc.inv, tc.payload); err == nil {
			t.Errorf("%s: invocation was allowed", name)
		}
	}
	drive := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1"}, nil)
	if _, err := client.Invoke(context.Background(), drive, invocation("drive.search", connectors.ActionRead, "drive/d-1"), DriveSearchPayload{}); err == nil {
		t.Error("a google_drive connector was accepted by the Gmail client")
	}
}

func TestAuthenticatedGmailSearchRequestsMetadataOnly(t *testing.T) {
	resolver := &recordingResolver{token: "fake-gmail-token"}
	var sawFormats []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer fake-gmail-token" {
			t.Errorf("authorization = %q", got)
		}
		switch r.URL.Path {
		case "/gmail/v1/users/me/messages":
			q := r.URL.Query()
			if q.Get("q") != "from:alice" || q.Get("maxResults") != "2" || q.Get("pageToken") != "p-1" {
				t.Errorf("list query = %v", q)
			}
			_, _ = io.WriteString(w, `{"messages":[{"id":"m-1","threadId":"th-1"}],"nextPageToken":"p-2"}`)
		case "/gmail/v1/users/me/messages/m-1":
			sawFormats = append(sawFormats, r.URL.Query().Get("format"))
			if got := r.URL.Query()["metadataHeaders"]; strings.Join(got, ",") != "From,To,Subject" {
				t.Errorf("metadataHeaders = %v", got)
			}
			_, _ = io.WriteString(w, gmailMetadataFixture)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := NewGmailHTTP(testHTTPClient(server, assertGmailEndpoint(t)), resolver)
	meta := gmailMeta(t, true)
	result, err := client.Invoke(context.Background(), meta, invocation("message.search", connectors.ActionRead, "gmail/messages"), MessageSearchPayload{Query: "from:alice", PageSize: 2, PageToken: "p-1"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	page := result.Data.(GmailMessagePage)
	want := GmailMessage{ID: "m-1", ThreadID: "th-1", InternalDate: time.UnixMilli(1700000000123).UTC(), From: "Alice <alice@example.test>", To: "bob@example.test", Subject: "Hello", Labels: []string{"INBOX", "UNREAD"}, Snippet: "hi bob"}
	if len(page.Messages) != 1 || page.NextPageToken != "p-2" || !gmailMessagesEqual(page.Messages[0], want) {
		t.Fatalf("page = %+v", page)
	}
	if strings.Join(sawFormats, ",") != "metadata" {
		t.Fatalf("formats requested = %v, want metadata only", sawFormats)
	}
	if resolver.ten != "t-1" || resolver.workspace != "w-1" || resolver.subject != "u-1" || resolver.ref != meta.CredentialRef {
		t.Fatalf("scope = %+v", resolver)
	}
	assertNoGmailBody(t, result)
}

func TestAuthenticatedGmailGetFetchesBodyOnlyWithOptIn(t *testing.T) {
	for _, includeBody := range []bool{false, true} {
		var format string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/gmail/v1/users/me/messages/m-1" {
				t.Errorf("path = %q", r.URL.Path)
			}
			format = r.URL.Query().Get("format")
			if format == "full" {
				_, _ = io.WriteString(w, gmailFullFixture)
				return
			}
			_, _ = io.WriteString(w, gmailMetadataFixture)
		}))
		client := NewGmailHTTP(testHTTPClient(server, assertGmailEndpoint(t)), &recordingResolver{token: "fake-gmail-token"})
		result, err := client.Invoke(context.Background(), gmailMeta(t, includeBody), invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{MessageID: "m-1"})
		server.Close()
		if err != nil {
			t.Fatalf("includeBody=%v: %v", includeBody, err)
		}
		msg := result.Data.(GmailMessage)
		if !includeBody {
			if format != "metadata" {
				t.Errorf("default get requested format %q, want metadata", format)
			}
			assertNoGmailBody(t, result)
			continue
		}
		if format != "full" {
			t.Errorf("opt-in get requested format %q, want full", format)
		}
		if msg.Body != gmailFixtureBody || msg.Subject != "Hello" {
			t.Errorf("opt-in message = %+v", msg)
		}
	}
}

func TestAuthenticatedGmailFailsClosedAndRedactsErrors(t *testing.T) {
	client, err := New(connectors.ProviderGmail)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Invoke(context.Background(), gmailMeta(t, false), invocation("message.search", connectors.ActionRead, "gmail/messages"), MessageSearchPayload{})
	if err == nil || !strings.Contains(err.Error(), "credential") {
		t.Fatalf("unconfigured resolver error = %v", err)
	}

	noNetwork := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("must not send") })}
	_, err = NewGmailHTTP(noNetwork, credentialErrorResolver{}).Invoke(context.Background(), gmailMeta(t, false), invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{})
	if err == nil || err.Error() != "credential resolution failed" {
		t.Fatalf("resolver error = %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"FAKE-BODY-MARKER fake-gmail-token"}`, http.StatusForbidden)
	}))
	defer server.Close()
	_, err = NewGmailHTTP(testHTTPClient(server, assertGmailEndpoint(t)), &recordingResolver{token: "fake-gmail-token"}).Invoke(context.Background(), gmailMeta(t, false), invocation("message.get", connectors.ActionRead, "gmail/messages/m-1"), MessageGetPayload{})
	if err == nil || strings.Contains(err.Error(), "fake-gmail-token") || strings.Contains(err.Error(), "FAKE-BODY-MARKER") {
		t.Fatalf("provider rejection error = %v", err)
	}
}

const gmailMetadataFixture = `{"id":"m-1","threadId":"th-1","labelIds":["INBOX","UNREAD"],"snippet":"hi bob","internalDate":"1700000000123",
 "payload":{"headers":[{"name":"From","value":"Alice <alice@example.test>"},{"name":"To","value":"bob@example.test"},{"name":"Subject","value":"Hello"}]}}`

var gmailFullFixture = `{"id":"m-1","threadId":"th-1","labelIds":["INBOX"],"snippet":"hi bob","internalDate":"1700000000123",
 "payload":{"mimeType":"multipart/alternative","headers":[{"name":"Subject","value":"Hello"}],"parts":[
  {"mimeType":"text/html","body":{"data":"` + base64.URLEncoding.EncodeToString([]byte("<p>html</p>")) + `"}},
  {"mimeType":"text/plain","body":{"data":"` + base64.RawURLEncoding.EncodeToString([]byte(gmailFixtureBody)) + `"}}]}}`

func assertGmailEndpoint(t *testing.T) func(*testing.T, *http.Request) {
	return func(_ *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "gmail.googleapis.com" || req.Method != http.MethodGet {
			t.Errorf("request = %s %s", req.Method, req.URL.String())
		}
	}
}

func assertNoGmailBody(t *testing.T, result Result) {
	t.Helper()
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "FAKE-BODY-MARKER") || strings.Contains(string(encoded), `"body"`) {
		t.Fatalf("result carries a message body: %s", encoded)
	}
}

func gmailMessagesEqual(a, b GmailMessage) bool {
	return a.ID == b.ID && a.ThreadID == b.ThreadID && a.InternalDate.Equal(b.InternalDate) && a.From == b.From && a.To == b.To && a.Subject == b.Subject && strings.Join(a.Labels, ",") == strings.Join(b.Labels, ",") && a.Snippet == b.Snippet && a.Body == b.Body
}

type leakyGmail struct{}

func (leakyGmail) SearchMessages(context.Context, string, int, string) ([]GmailMessage, string, error) {
	return []GmailMessage{{ID: "m-1", Snippet: "s", Body: "FAKE-BODY-MARKER leaked"}}, "", nil
}

func (leakyGmail) Message(context.Context, string, bool) (GmailMessage, error) {
	return GmailMessage{ID: "m-1", Snippet: "s", Body: "FAKE-BODY-MARKER leaked"}, nil
}

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

// Package mockapi serves a provider-shaped HTTP API generated from the
// repository's contract metadata: one route per manifest operation, fixture
// responses, seeded synthetic pages, and failure injection. The tests run
// the generated Go operations against it.
package mockapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1/githubops"
	notionv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/notion/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/notion/v1/notionops"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// repoRoot is the repository root seen from this package's directory.
const repoRoot = "../.."

const (
	manifestPath = repoRoot + "/manifest/connector-manifest.v1.json"
	fixturesPath = repoRoot + "/fixtures/operations"
)

func newTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()
	s, err := New(Options{
		ManifestPath: manifestPath,
		FixturesDir:  fixturesPath,
		Seed:         7,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ts := httptest.NewServer(s)
	t.Cleanup(ts.Close)
	return s, ts
}

// bearerConn builds a connection.HTTP against the mock, optionally wrapping
// the transport to inject an X-Mock-Fail directive.
func bearerConn(t *testing.T, base string, fail string) connection.HTTP {
	t.Helper()
	client := http.DefaultClient
	if fail != "" {
		client = &http.Client{Transport: &failTransport{next: http.DefaultTransport, fail: fail}}
	}
	conn, err := connection.NewHTTP(base, &contracts.Connection{Name: "pat", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER}, "mock-token", client)
	if err != nil {
		t.Fatalf("NewHTTP: %v", err)
	}
	return conn
}

type failTransport struct {
	next http.RoundTripper
	fail string
}

func (f *failTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Mock-Fail", f.fail)
	return f.next.RoundTrip(req)
}

// TestEveryManifestOperationHasRoute: every manifest capability either has a
// route in the server's route table (HTTP and GraphQL operations) or is a
// SQL operation with deliberately no HTTP route. Both directions are
// asserted: no uncovered operation, no route without an operation.
func TestEveryManifestOperationHasRoute(t *testing.T) {
	s, _ := newTestServer(t)
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest struct {
		Connectors []struct {
			ID           string `json:"id"`
			Capabilities []struct {
				ID string `json:"id"`
			} `json:"capabilities"`
		} `json:"connectors"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}
	routes := map[string]string{} // op id -> "METHOD /template"
	for _, r := range s.Routes() {
		routes[r.OpID] = r.Method + " " + r.Template
	}
	covered := map[string]bool{}
	for _, c := range manifest.Connectors {
		for _, cap := range c.Capabilities {
			covered[cap.ID] = true
			op, ok := contracts.Lookup(cap.ID)
			if !ok {
				t.Fatalf("%s: not registered in the contracts registry (import its generated package)", cap.ID)
			}
			route, hasRoute := routes[cap.ID]
			switch {
			case op.GetHttpMethod() == "" && op.GetGraphqlDocument() == "":
				if hasRoute {
					t.Errorf("%s is SQL yet has route %s", cap.ID, route)
				}
			case !hasRoute:
				t.Errorf("%s: no route (want %s %s)", cap.ID, op.GetHttpMethod(), op.GetPathTemplate())
			}
		}
	}
	for id, route := range routes {
		if !covered[id] {
			t.Errorf("route %s (%s) serves no manifest operation", route, id)
		}
	}
}

// TestServesFixture_GitHubIssuesGet: the generated GitHub client decodes the
// served fixture losslessly, including the id above 2^53.
func TestServesFixture_GitHubIssuesGet(t *testing.T) {
	_, ts := newTestServer(t)
	fixture := readFixtureFile(t, "github", "issues.get.response.json")
	var want struct {
		ID     int64  `json:"id"`
		Number int64  `json:"number"`
		Title  string `json:"title"`
	}
	if err := json.Unmarshal(fixture, &want); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	resp, err := githubops.New(bearerConn(t, ts.URL, "")).GetIssue(context.Background(), &githubv1.GetIssueRequest{Owner: "acme", Repository: "widgets", IssueNumber: want.Number})
	if err != nil {
		t.Fatalf("GetIssue: %v", err)
	}
	if resp.GetId() != want.ID {
		t.Fatalf("id = %d, want %d (fixture value, lossless above 2^53)", resp.GetId(), want.ID)
	}
	if resp.GetTitle() != want.Title {
		t.Fatalf("title = %q, want %q", resp.GetTitle(), want.Title)
	}
}

// TestCursorPaging_GitHubPullsList: page one carries next_page = 2; the
// cursor page returns a varied list of a different length and no cursor.
func TestCursorPaging_GitHubPullsList(t *testing.T) {
	_, ts := newTestServer(t)
	fixture := readFixtureFile(t, "github", "pulls.list.response.json")
	var first struct {
		Pulls []map[string]any `json:"pulls"`
	}
	if err := json.Unmarshal(fixture, &first); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	conn := githubops.New(bearerConn(t, ts.URL, ""))
	ctx := context.Background()
	page1, err := conn.ListPulls(ctx, &githubv1.ListPullsRequest{Owner: "acme", Repository: "widgets"})
	if err != nil {
		t.Fatalf("ListPulls page 1: %v", err)
	}
	if got := len(page1.GetPulls()); got != len(first.Pulls) {
		t.Fatalf("page 1 pulls = %d, want fixture length %d", got, len(first.Pulls))
	}
	if page1.GetNextPage() != 2 {
		t.Fatalf("page 1 next_page = %d, want 2", page1.GetNextPage())
	}
	page2, err := conn.ListPulls(ctx, &githubv1.ListPullsRequest{Owner: "acme", Repository: "widgets", Page: proto.Int32(2)})
	if err != nil {
		t.Fatalf("ListPulls page 2: %v", err)
	}
	if got := len(page2.GetPulls()); got == len(first.Pulls) {
		t.Fatalf("page 2 list length %d must differ from page 1", got)
	}
	if page2.GetNextPage() != 0 {
		t.Fatalf("page 2 next_page = %d, want 0 (last page)", page2.GetNextPage())
	}
	if page2.GetPulls()[0].GetId() == first.Pulls[0]["id"] {
		t.Fatalf("page 2 ids are not varied")
	}
}

// TestCursorPaging_NotionDatabasesQuery: the string cursor rides the POST
// body; the second page is the last.
func TestCursorPaging_NotionDatabasesQuery(t *testing.T) {
	_, ts := newTestServer(t)
	conn := notionops.New(bearerConn(t, ts.URL, ""))
	ctx := context.Background()
	page1, err := conn.QueryDatabase(ctx, &notionv1.QueryDatabaseRequest{DatabaseId: "b2c3d4e5-0000-4000-8000-000000000001"})
	if err != nil {
		t.Fatalf("QueryDatabases page 1: %v", err)
	}
	if page1.GetNextCursor() == "" {
		t.Fatalf("page 1 next_cursor is empty, want the mock's token")
	}
	page2, err := conn.QueryDatabase(ctx, &notionv1.QueryDatabaseRequest{DatabaseId: "b2c3d4e5-0000-4000-8000-000000000001", StartCursor: proto.String(page1.GetNextCursor())})
	if err != nil {
		t.Fatalf("QueryDatabases page 2: %v", err)
	}
	if page2.GetNextCursor() != "" {
		t.Fatalf("page 2 next_cursor = %q, want empty (last page)", page2.GetNextCursor())
	}
}

// TestFailureInjection_MapsGeneratedErrors: X-Mock-Fail statuses surface as
// the generated client's mapped error kinds.
func TestFailureInjection_MapsGeneratedErrors(t *testing.T) {
	cases := []struct {
		fail      string
		wantKind  contracts.ErrorKind
		retryable bool
	}{
		{"401", contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false},
		{"404", contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false},
		{"429", contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true},
		{"500", contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, false},
	}
	for _, tc := range cases {
		t.Run(tc.fail, func(t *testing.T) {
			_, ts := newTestServer(t)
			_, err := githubops.New(bearerConn(t, ts.URL, tc.fail)).GetIssue(context.Background(), &githubv1.GetIssueRequest{Owner: "acme", Repository: "widgets", IssueNumber: 1})
			var e *connection.Error
			if !errAs(err, &e) {
				t.Fatalf("GetIssue error = %v, want *connection.Error", err)
			}
			if e.Kind != tc.wantKind || e.Retryable != tc.retryable {
				t.Fatalf("error = %+v, want kind %v retryable %v", e, tc.wantKind, tc.retryable)
			}
		})
	}
}

// TestRejectsWrongCredentialPlacement: a request without the credential in
// its declared place is a 400 with a reason, and the generated client maps
// it to an invalid-argument kind.
func TestRejectsWrongCredentialPlacement(t *testing.T) {
	_, ts := newTestServer(t)
	resp, err := http.Get(ts.URL + "/repos/acme/widgets/issues/1?token=mock-token")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "credential") {
		t.Fatalf("body = %s, want a reason mentioning the credential", body)
	}
}

// TestRejectsUnescapedPath: a path value with raw reserved characters is a
// 400 with a reason (the connection runtime always escapes them).
func TestRejectsUnescapedPath(t *testing.T) {
	_, ts := newTestServer(t)
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/repos/acme(ops)/widgets/issues/1", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer mock-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "escaped") {
		t.Fatalf("body = %s, want a reason about escaping", body)
	}
}

func readFixtureFile(t *testing.T, connector, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesPath, connector, name))
	if err != nil {
		t.Fatalf("read fixture %s/%s: %v", connector, name, err)
	}
	return data
}

func errAs(err error, target any) bool { return errors.As(err, target) }

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

// Replay tests: each generated operation runs against a local test server
// that checks the rendered method, path, query, body and credential
// placement, then returns the operation's response fixture.

package conformance

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1/databaseops"
	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1/githubops"
	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1/httpapiops"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1/linearops"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// replayToken is a synthetic credential.
const replayToken = "kei-replay-token-0001"

// recorded is what the test server saw.
type recorded struct {
	method, escapedPath string
	query               url.Values
	auth                string
	body                []byte
}

// replayServer serves body for every request and records it.
func replayServer(t *testing.T, status int, header http.Header, body []byte) (*httptest.Server, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method, rec.escapedPath, rec.query = r.Method, r.URL.EscapedPath(), r.URL.Query()
		rec.auth = r.Header.Get("Authorization")
		rec.body, _ = io.ReadAll(r.Body)
		for k, v := range header {
			w.Header()[k] = v
		}
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv, rec
}

func bearer(t *testing.T, base string) connection.HTTP {
	t.Helper()
	conn, err := connection.NewHTTP(base, &contracts.Connection{Name: "pat", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER}, replayToken, http.DefaultClient)
	if err != nil {
		t.Fatalf("NewHTTP: %v", err)
	}
	return conn
}

func loadRequest(t *testing.T, connector, op string, msg proto.Message) {
	t.Helper()
	if err := protojson.Unmarshal(readFixture(t, connector, opFixtureName(op, "request")), msg); err != nil {
		t.Fatalf("request fixture %s: %v", op, err)
	}
}

func TestGeneratedGetIssue_Replay(t *testing.T) {
	srv, rec := replayServer(t, http.StatusOK, nil, readFixture(t, "github", "issues.get.response.json"))
	q := githubops.New(bearer(t, srv.URL))
	resp, err := q.GetIssue(context.Background(), &githubv1.GetIssueRequest{Owner: "acme", Repository: "widgets", IssueNumber: 1347})
	if err != nil {
		t.Fatalf("GetIssue: %v", err)
	}
	if rec.method != http.MethodGet || rec.escapedPath != "/repos/acme/widgets/issues/1347" {
		t.Fatalf("request = %s %s", rec.method, rec.escapedPath)
	}
	if rec.auth != "Bearer "+replayToken {
		t.Fatalf("Authorization = %q", rec.auth)
	}
	if len(rec.query) != 0 || len(rec.body) != 0 {
		t.Fatalf("unexpected query %v or body %q", rec.query, rec.body)
	}
	if resp.GetNumber() != 1347 {
		t.Fatalf("Number = %d, want 1347", resp.GetNumber())
	}
	// Review Focus 2: the fixture's id is 2^53+1 and it has an unknown key.
	if resp.GetId() != 9007199254740993 {
		t.Fatalf("Id = %d, want 9007199254740993 (lossless above 2^53)", resp.GetId())
	}
	var _ githubops.Operations = q
	if githubops.OpGetIssue.ID != "github.issues.get" || githubops.OpGetIssue.Class != connection.Read {
		t.Fatalf("OpGetIssue = %+v", githubops.OpGetIssue)
	}
}

// TestGeneratedHTTP_Replay replays every HTTP pilot operation from its
// request and response fixtures.
func TestGeneratedHTTP_Replay(t *testing.T) {
	cases := []struct {
		op        string
		connector string
		// serve rewrites the response fixture into the provider's wire shape.
		serve  func(t *testing.T, fixture []byte) []byte
		method string
		path   string
		query  url.Values
		body   string // canonical JSON, empty for none
		call   func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error)
		check  func(t *testing.T, resp proto.Message)
	}{
		{
			op: "github.issues.list", connector: "github",
			serve:  unwrapArray("issues"), // GitHub returns a bare JSON array
			method: "GET", path: "/repos/example-org/example-repo/issues",
			query: url.Values{"state": {"open"}, "per_page": {"2"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListIssuesRequest{}
				loadRequest(t, "github", "github.issues.list", req)
				return githubops.New(conn).ListIssues(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*githubv1.ListIssuesResponse).GetIssues()); n != 2 {
					t.Fatalf("issues = %d, want 2", n)
				}
			},
		},
		{
			op: "github.pulls.get", connector: "github",
			method: "GET", path: "/repos/example-org/example-repo/pulls/42",
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.GetPullRequestRequest{}
				loadRequest(t, "github", "github.pulls.get", req)
				return githubops.New(conn).GetPullRequest(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if m.(*githubv1.GetPullRequestResponse).GetNumber() == 0 {
					t.Fatalf("number not decoded")
				}
			},
		},
		{
			op: "github.contents.put", connector: "github",
			// path is multi_segment: its "/" is kept, each segment escaped.
			method: "PUT", path: "/repos/example-org/example-repo/contents/docs/notes.txt",
			body: `{"branch":"main","content":"c3ludGhldGljIGZpbGUgY29udGVudA==","message":"Synthetic pilot commit"}`,
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.PutContentsRequest{}
				loadRequest(t, "github", "github.contents.put", req)
				return githubops.New(conn).PutContents(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if m.(*githubv1.PutContentsResponse).GetCommit().GetSha() == "" {
					t.Fatalf("commit sha not decoded")
				}
			},
		},
		{
			op: "github.pulls.list", connector: "github",
			serve:  unwrapArray("pulls"), // GitHub returns a bare JSON array
			method: "GET", path: "/repos/example-org/example-repo/pulls",
			query: url.Values{"state": {"open"}, "per_page": {"30"}, "page": {"1"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListPullsRequest{}
				loadRequest(t, "github", "github.pulls.list", req)
				return githubops.New(conn).ListPulls(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*githubv1.ListPullsResponse).GetPulls()); n != 1 {
					t.Fatalf("pulls = %d, want 1", n)
				}
			},
		},
		{
			op: "github.issues.create", connector: "github",
			method: "POST", path: "/repos/example-org/example-repo/issues",
			body: `{"body":"Synthetic body text","labels":["synthetic-label"],"title":"Synthetic new issue"}`,
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.CreateIssueRequest{}
				loadRequest(t, "github", "github.issues.create", req)
				return githubops.New(conn).CreateIssue(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if got := m.(*githubv1.CreateIssueResponse).GetNumber(); got != 42 {
					t.Fatalf("number = %d, want 42", got)
				}
			},
		},
		{
			op: "github.actions.workflows.runs.list", connector: "github",
			serve:  unwrapArray("workflow_runs"), // GitHub returns a bare JSON array
			method: "GET", path: "/repos/example-org/example-repo/actions/workflows/ci.yml/runs",
			query: url.Values{"branch": {"main"}, "event": {"push"}, "per_page": {"30"}, "page": {"1"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListWorkflowRunsByWorkflowRequest{}
				loadRequest(t, "github", "github.actions.workflows.runs.list", req)
				return githubops.New(conn).ListWorkflowRunsByWorkflow(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				runs := m.(*githubv1.ListWorkflowRunsByWorkflowResponse).GetWorkflowRuns()
				if len(runs) != 1 {
					t.Fatalf("workflow_runs = %d, want 1", len(runs))
				}
				if got := runs[0].GetActor().GetLogin(); got != "example-user" {
					t.Fatalf("actor.login = %q, want example-user", got)
				}
				if got := runs[0].GetHeadSha(); got != "synthetic-head-sha-1" {
					t.Fatalf("head_sha = %q, want the synthetic sha", got)
				}
			},
		},
		{
			op: "github.repos.contents.get", connector: "github",
			// path is multi_segment: its "/" is kept, each segment escaped.
			method: "GET", path: "/repos/example-org/example-repo/contents/docs/notes.md",
			query: url.Values{"ref": {"main"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.GetContentsRequest{}
				loadRequest(t, "github", "github.repos.contents.get", req)
				return githubops.New(conn).GetContents(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if got := m.(*githubv1.GetContentsResponse).GetEncoding(); got != "base64" {
					t.Fatalf("encoding = %q, want base64", got)
				}
			},
		},
		{
			op: "github.repos.contents.list", connector: "github",
			serve:  unwrapArray("entries"), // GitHub returns a bare JSON array
			method: "GET", path: "/repos/example-org/example-repo/contents/docs",
			query: url.Values{"ref": {"main"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListContentsRequest{}
				loadRequest(t, "github", "github.repos.contents.list", req)
				return githubops.New(conn).ListContents(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*githubv1.ListContentsResponse).GetEntries()); n != 2 {
					t.Fatalf("entries = %d, want 2", n)
				}
			},
		},
		{
			op: "github.repos.commits.list", connector: "github",
			serve:  unwrapArray("commits"), // GitHub returns a bare JSON array
			method: "GET", path: "/repos/example-org/example-repo/commits",
			query: url.Values{"sha": {"main"}, "per_page": {"30"}, "page": {"1"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListCommitsRequest{}
				loadRequest(t, "github", "github.repos.commits.list", req)
				return githubops.New(conn).ListCommits(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				commits := m.(*githubv1.ListCommitsResponse).GetCommits()
				if len(commits) != 1 {
					t.Fatalf("commits = %d, want 1", len(commits))
				}
				if got := commits[0].GetAuthor().GetLogin(); got != "example-user" {
					t.Fatalf("commits[0].author.login = %q, want example-user", got)
				}
			},
		},
		{
			op: "github.repos.commits.get", connector: "github",
			// ref is multi_segment: a sha or branch name.
			method: "GET", path: "/repos/example-org/example-repo/commits/synthetic-commit-sha-3",
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.GetCommitRequest{}
				loadRequest(t, "github", "github.repos.commits.get", req)
				return githubops.New(conn).GetCommit(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if got := m.(*githubv1.GetCommitResponse).GetCommit().GetMessage(); got != "Synthetic commit message" {
					t.Fatalf("commit.message = %q, want Synthetic commit message", got)
				}
			},
		},
		{
			op: "github.actions.runs.attempts.get", connector: "github",
			method: "GET", path: "/repos/example-org/example-repo/actions/runs/9007199254741001/attempts/2",
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.GetRunAttemptRequest{}
				loadRequest(t, "github", "github.actions.runs.attempts.get", req)
				return githubops.New(conn).GetRunAttempt(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if got := m.(*githubv1.GetRunAttemptResponse).GetRunAttempt(); got != 2 {
					t.Fatalf("run_attempt = %d, want 2", got)
				}
			},
		},
		{
			op: "github.actions.runs.jobs.list", connector: "github",
			method: "GET", path: "/repos/example-org/example-repo/actions/runs/9007199254741001/jobs",
			query: url.Values{"per_page": {"30"}, "page": {"1"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &githubv1.ListRunJobsRequest{}
				loadRequest(t, "github", "github.actions.runs.jobs.list", req)
				return githubops.New(conn).ListRunJobs(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*githubv1.ListRunJobsResponse).GetJobs()); n != 1 {
					t.Fatalf("jobs = %d, want 1", n)
				}
			},
		},
		{
			op: "http_api.records.list", connector: "http_api",
			method: "GET", path: "/contacts", query: url.Values{"limit": {"2"}},
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &httpapiv1.ListRecordsRequest{}
				loadRequest(t, "http_api", "http_api.records.list", req)
				return httpapiops.New(conn).ListRecords(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*httpapiv1.ListRecordsResponse).GetRecords()); n != 2 {
					t.Fatalf("records = %d, want 2", n)
				}
			},
		},
		{
			op: "http_api.records.get", connector: "http_api",
			method: "GET", path: "/contacts/rec-1",
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &httpapiv1.GetRecordRequest{}
				loadRequest(t, "http_api", "http_api.records.get", req)
				return httpapiops.New(conn).GetRecord(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if m.(*httpapiv1.GetRecordResponse).GetRecord() == nil {
					t.Fatalf("record not decoded")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.op, func(t *testing.T) {
			fixture := readFixture(t, tc.connector, opFixtureName(tc.op, "response"))
			if tc.serve != nil {
				fixture = tc.serve(t, fixture)
			}
			srv, rec := replayServer(t, http.StatusOK, nil, fixture)
			resp, err := tc.call(context.Background(), bearer(t, srv.URL), t)
			if err != nil {
				t.Fatalf("call: %v", err)
			}
			if rec.method != tc.method || rec.escapedPath != tc.path {
				t.Fatalf("request = %s %s, want %s %s", rec.method, rec.escapedPath, tc.method, tc.path)
			}
			if rec.auth != "Bearer "+replayToken {
				t.Fatalf("Authorization = %q", rec.auth)
			}
			wantQuery := tc.query
			if wantQuery == nil {
				wantQuery = url.Values{}
			}
			if !reflect.DeepEqual(rec.query, wantQuery) {
				t.Fatalf("query = %v, want %v", rec.query, wantQuery)
			}
			if string(rec.body) != tc.body {
				t.Fatalf("body = %s, want %s", rec.body, tc.body)
			}
			tc.check(t, resp)
		})
	}
}

// TestGeneratedGraphQL_Replay replays the Linear operations: the request is
// {"query": <the operation's document>, "variables": <request fields by
// JSON name>}, and the response fixture is served inside GraphQL's data
// envelope (under the root field when the fixture is that field's value).
func TestGeneratedGraphQL_Replay(t *testing.T) {
	cases := []struct {
		op        string
		root      string // the data key wrapping the fixture; empty when the fixture is data itself
		variables string
		call      func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error)
		check     func(t *testing.T, resp proto.Message)
	}{
		{
			op:        "linear.issues.get",
			variables: `{"id":"00000000-0000-0000-0000-000000000101"}`,
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &linearv1.GetIssueRequest{}
				loadRequest(t, "linear", "linear.issues.get", req)
				return linearops.New(conn).GetIssue(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if got := m.(*linearv1.GetIssueResponse).GetIssue().GetIdentifier(); got != "ENG-101" {
					t.Fatalf("identifier = %q", got)
				}
			},
		},
		{
			op:        "linear.issues.create",
			root:      "issueCreate",
			variables: `{"priority":3,"teamId":"00000000-0000-0000-0000-000000000002","title":"Synthetic pilot issue"}`,
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &linearv1.CreateIssueRequest{}
				loadRequest(t, "linear", "linear.issues.create", req)
				return linearops.New(conn).CreateIssue(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				r := m.(*linearv1.CreateIssueResponse)
				if !r.GetSuccess() || r.GetIssue().GetIdentifier() != "ENG-102" {
					t.Fatalf("resp = %v", r)
				}
			},
		},
		{
			op:        "linear.issues.search",
			root:      "searchIssues",
			variables: `{"first":10,"term":"synthetic"}`,
			call: func(ctx context.Context, conn connection.HTTP, t *testing.T) (proto.Message, error) {
				req := &linearv1.SearchIssuesRequest{}
				loadRequest(t, "linear", "linear.issues.search", req)
				return linearops.New(conn).SearchIssues(ctx, req)
			},
			check: func(t *testing.T, m proto.Message) {
				if n := len(m.(*linearv1.SearchIssuesResponse).GetNodes()); n == 0 {
					t.Fatalf("no nodes decoded")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.op, func(t *testing.T) {
			fixture := readFixture(t, "linear", opFixtureName(tc.op, "response"))
			data := json.RawMessage(fixture)
			if tc.root != "" {
				data, _ = json.Marshal(map[string]json.RawMessage{tc.root: fixture})
			}
			envelope, _ := json.Marshal(map[string]json.RawMessage{"data": data})
			srv, rec := replayServer(t, http.StatusOK, nil, envelope)
			resp, err := tc.call(context.Background(), bearer(t, srv.URL), t)
			if err != nil {
				t.Fatalf("call: %v", err)
			}
			if rec.method != http.MethodPost || rec.escapedPath != "/graphql" {
				t.Fatalf("request = %s %s", rec.method, rec.escapedPath)
			}
			var body struct {
				Query     string          `json:"query"`
				Variables json.RawMessage `json:"variables"`
			}
			if err := json.Unmarshal(rec.body, &body); err != nil {
				t.Fatalf("body %s: %v", rec.body, err)
			}
			op, ok := contracts.Lookup(tc.op)
			if !ok {
				t.Fatalf("%s is not registered", tc.op)
			}
			if body.Query != op.GetGraphqlDocument() {
				t.Fatalf("query = %q, want the operation's document", body.Query)
			}
			if string(body.Variables) != tc.variables {
				t.Fatalf("variables = %s, want %s", body.Variables, tc.variables)
			}
			tc.check(t, resp)
		})
	}
}

// TestGeneratedHTTP_ErrorMapping checks the declared error table and that no
// credential or body reaches the error (Review Focus 3).
func TestGeneratedHTTP_ErrorMapping(t *testing.T) {
	srv, _ := replayServer(t, http.StatusTooManyRequests, nil, []byte(`{"message":"slow down `+replayToken+`"}`))
	_, err := githubops.New(bearer(t, srv.URL)).GetIssue(context.Background(), &githubv1.GetIssueRequest{Owner: "acme", Repository: "widgets", IssueNumber: 1})
	var e *connection.Error
	if !errors.As(err, &e) {
		t.Fatalf("GetIssue error = %v, want *connection.Error", err)
	}
	if e.Kind != contracts.ErrorKind_ERROR_KIND_RATE_LIMITED || !e.Retryable || e.OpID != "github.issues.get" {
		t.Fatalf("error = %+v", e)
	}
	if strings.Contains(err.Error(), replayToken) || strings.Contains(err.Error(), "slow down") {
		t.Fatalf("error leaks the body or credential: %v", err)
	}
}

// TestGeneratedHTTP_PathEscaping (Review Focus 1): reserved characters in
// path values never change the path shape.
func TestGeneratedHTTP_PathEscaping(t *testing.T) {
	srv, rec := replayServer(t, http.StatusOK, nil, readFixture(t, "github", "issues.get.response.json"))
	_, err := githubops.New(bearer(t, srv.URL)).GetIssue(context.Background(), &githubv1.GetIssueRequest{Owner: "a/b", Repository: "..", IssueNumber: 7})
	if err != nil {
		t.Fatalf("GetIssue: %v", err)
	}
	if rec.escapedPath != "/repos/a%2Fb/%2E%2E/issues/7" {
		t.Fatalf("path = %s", rec.escapedPath)
	}
}

// unwrapArray turns the fixture's {"<field>": [...]} into the bare array a
// provider sends.
func unwrapArray(field string) func(t *testing.T, fixture []byte) []byte {
	return func(t *testing.T, fixture []byte) []byte {
		t.Helper()
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(fixture, &obj); err != nil {
			t.Fatalf("fixture: %v", err)
		}
		return obj[field]
	}
}

// replayDriver is a database/sql driver that serves fixed rows.
type replayDriver struct {
	mu    sync.Mutex
	query string
	args  []driver.NamedValue
	cols  []string
	rows  [][]driver.Value
}

func (d *replayDriver) Open(string) (driver.Conn, error) { return replayConn{d}, nil }

type replayConn struct{ d *replayDriver }

func (c replayConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c replayConn) Close() error                        { return nil }
func (c replayConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (c replayConn) QueryContext(_ context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	c.d.mu.Lock()
	defer c.d.mu.Unlock()
	c.d.query, c.d.args = q, args
	return &replayRows{cols: c.d.cols, rows: c.d.rows}, nil
}

type replayRows struct {
	cols []string
	rows [][]driver.Value
	i    int
}

func (r *replayRows) Columns() []string { return r.cols }
func (r *replayRows) Close() error      { return nil }
func (r *replayRows) Next(dest []driver.Value) error {
	if r.i >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.i])
	r.i++
	return nil
}

var registerReplayDriver sync.Once

// TestGeneratedRunQuery_Replay runs database.query.run over a fake driver
// whose rows reproduce the response fixture, and compares the result with
// the fixture.
func TestGeneratedRunQuery_Replay(t *testing.T) {
	d := &replayDriver{
		cols: []string{"id", "name", "note", "balance"},
		rows: [][]driver.Value{
			{int64(9007199254740993), "Big-ID Account — ünïcodé ✓", nil, []byte("12345678901234567890.123456789012345678")},
			{int64(1), "Synthetic One", "Synthetic note", []byte("0.000000000000000001")},
		},
	}
	registerReplayDriver.Do(func() { sql.Register("kei-replay", d) })
	db, err := sql.Open("kei-replay", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	req := &databasev1.RunQueryRequest{}
	loadRequest(t, "database", "database.query.run", req)
	resp, err := databaseops.New(connection.NewSQL(db)).RunQuery(context.Background(), req)
	if err != nil {
		t.Fatalf("RunQuery: %v", err)
	}
	if d.query != req.GetSql() || len(d.args) != 1 || d.args[0].Value != "9007199254740993" {
		t.Fatalf("driver got %q %v", d.query, d.args)
	}
	rows := resp.GetRows()
	if len(rows) != 2 || resp.GetTruncated() {
		t.Fatalf("rows = %d truncated = %v", len(rows), resp.GetTruncated())
	}
	first := rows[0].GetValues()
	if first[0].GetStringValue() != "9007199254740993" {
		t.Fatalf("big id = %v, want a lossless string", first[0])
	}
	if first[3].GetStringValue() != "12345678901234567890.123456789012345678" {
		t.Fatalf("numeric = %v, want a precision-preserving string", first[3])
	}
	if _, ok := first[2].GetKind().(*structpb.Value_NullValue); !ok {
		t.Fatalf("null cell = %v, want null", first[2])
	}
	if rows[1].GetValues()[0].GetNumberValue() != 1 {
		t.Fatalf("small id = %v, want the number 1", rows[1].GetValues()[0])
	}
	if names := []string{resp.GetColumns()[0].GetName(), resp.GetColumns()[3].GetName()}; names[0] != "id" || names[1] != "balance" {
		t.Fatalf("columns = %v", resp.GetColumns())
	}
}

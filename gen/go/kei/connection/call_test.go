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

package connection_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

type seen struct {
	method, path string
	query        url.Values
	body         []byte
}

func callServer(t *testing.T, header http.Header, status int, body string) (connection.HTTP, *seen) {
	t.Helper()
	s := &seen{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.method, s.path, s.query = r.Method, r.URL.EscapedPath(), r.URL.Query()
		s.body, _ = io.ReadAll(r.Body)
		for k, v := range header {
			w.Header()[k] = v
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	conn, err := connection.NewHTTP(srv.URL, bearerConn(), syntheticToken, srv.Client())
	if err != nil {
		t.Fatalf("NewHTTP: %v", err)
	}
	return conn, s
}

func jsonEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("body %q: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want %q: %v", want, err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

var opRead = connection.Op{ID: "acme.things.get", Class: connection.Read}

func TestCallHTTP_PlacementAndBigIntBody(t *testing.T) {
	conn, s := callServer(t, nil, 200, `{}`)
	spec := &connection.HTTPSpec{
		Op: opRead, Method: "POST", PathTemplate: "/repos/{owner}/{repository}",
		Params: []connection.Param{
			{Field: "owner", In: connection.InPath, Wire: "owner"},
			{Field: "repository", In: connection.InPath, Wire: "repository"},
			{Field: "issue_number", In: connection.InBody, Wire: "number"},
		},
	}
	req := &githubv1.GetIssueRequest{Owner: "acme", Repository: "w", IssueNumber: 9007199254740993}
	if err := connection.CallHTTP(context.Background(), conn, spec, req, &githubv1.GetIssueResponse{}); err != nil {
		t.Fatalf("CallHTTP: %v", err)
	}
	if s.path != "/repos/acme/w" {
		t.Fatalf("path = %s", s.path)
	}
	// An int64 goes in the body as a JSON number, exactly, not as a string.
	if string(s.body) != `{"number":9007199254740993}` {
		t.Fatalf("body = %s", s.body)
	}
}

func TestCallHTTP_MultiSegmentPath(t *testing.T) {
	spec := &connection.HTTPSpec{
		Op: opRead, Method: "PUT", PathTemplate: "/repos/{owner}/{repository}/contents/{path}",
		Params: []connection.Param{
			{Field: "owner", In: connection.InPath, Wire: "owner"},
			{Field: "repository", In: connection.InPath, Wire: "repository"},
			{Field: "path", In: connection.InPath, Wire: "path", MultiSegment: true},
		},
	}
	conn, s := callServer(t, nil, 200, `{}`)
	req := &githubv1.PutContentsRequest{Owner: "a/b", Repository: "w", Path: "docs/my notes.txt"}
	if err := connection.CallHTTP(context.Background(), conn, spec, req, &githubv1.PutContentsResponse{}); err != nil {
		t.Fatalf("CallHTTP: %v", err)
	}
	if s.path != "/repos/a%2Fb/w/contents/docs/my%20notes.txt" {
		t.Fatalf("path = %s", s.path)
	}
	req.Path = "docs/../secrets"
	if err := connection.CallHTTP(context.Background(), conn, spec, req, &githubv1.PutContentsResponse{}); err == nil {
		t.Fatalf("a .. segment: want an error")
	}
}

func TestCallHTTP_QueryOmitsUnsetFields(t *testing.T) {
	conn, s := callServer(t, nil, 200, `[]`)
	spec := &connection.HTTPSpec{
		Op: opRead, Method: "GET", PathTemplate: "/repos/{owner}/{repository}/issues",
		Params: []connection.Param{
			{Field: "owner", In: connection.InPath, Wire: "owner"},
			{Field: "repository", In: connection.InPath, Wire: "repository"},
			{Field: "state", In: connection.InQuery, Wire: "state"},
			{Field: "per_page", In: connection.InQuery, Wire: "per_page"},
		},
	}
	req := &githubv1.ListIssuesRequest{Owner: "acme", Repository: "w", PerPage: proto.Int32(0)}
	resp := &githubv1.ListIssuesResponse{}
	if err := connection.CallHTTP(context.Background(), conn, spec, req, resp); err != nil {
		t.Fatalf("CallHTTP: %v", err)
	}
	// state is unset and omitted; per_page is explicitly 0 and sent.
	if !reflect.DeepEqual(s.query, url.Values{"per_page": {"0"}}) {
		t.Fatalf("query = %v", s.query)
	}
	if len(s.body) != 0 {
		t.Fatalf("GET sent a body: %q", s.body)
	}
}

func TestCallHTTP_LinkPagination(t *testing.T) {
	spec := &connection.HTTPSpec{
		Op: opRead, Method: "GET", PathTemplate: "/{collection}",
		Params: []connection.Param{
			{Field: "collection", In: connection.InPath, Wire: "collection"},
			{Field: "limit", In: connection.InQuery, Wire: "limit"},
			{Field: "cursor", In: connection.InQuery, Wire: "cursor"},
		},
		CursorField: "cursor", NextCursorField: "next_cursor",
	}
	link := http.Header{"Link": {`<https://api.example.test/contacts?limit=2&cursor=page-3>; rel="next", <https://api.example.test/contacts?cursor=page-9>; rel="last"`}}

	t.Run("from Link header", func(t *testing.T) {
		conn, _ := callServer(t, link, 200, `{"records":[]}`)
		resp := &httpapiv1.ListRecordsResponse{}
		if err := connection.CallHTTP(context.Background(), conn, spec, &httpapiv1.ListRecordsRequest{Collection: "contacts"}, resp); err != nil {
			t.Fatalf("CallHTTP: %v", err)
		}
		if resp.GetNextCursor() != "page-3" {
			t.Fatalf("next cursor = %q, want page-3", resp.GetNextCursor())
		}
	})
	t.Run("body wins", func(t *testing.T) {
		conn, _ := callServer(t, link, 200, `{"records":[],"next_cursor":"from-body"}`)
		resp := &httpapiv1.ListRecordsResponse{}
		if err := connection.CallHTTP(context.Background(), conn, spec, &httpapiv1.ListRecordsRequest{Collection: "contacts"}, resp); err != nil {
			t.Fatalf("CallHTTP: %v", err)
		}
		if resp.GetNextCursor() != "from-body" {
			t.Fatalf("next cursor = %q", resp.GetNextCursor())
		}
	})
	t.Run("last page", func(t *testing.T) {
		conn, _ := callServer(t, nil, 200, `{"records":[]}`)
		resp := &httpapiv1.ListRecordsResponse{}
		if err := connection.CallHTTP(context.Background(), conn, spec, &httpapiv1.ListRecordsRequest{Collection: "contacts"}, resp); err != nil {
			t.Fatalf("CallHTTP: %v", err)
		}
		if resp.NextCursor != nil {
			t.Fatalf("next cursor = %q, want unset on the last page", resp.GetNextCursor())
		}
	})
}

const searchDoc = "query SearchIssues($term: String!, $first: Int) { searchIssues(term: $term, first: $first) { nodes { id } pageInfo { hasNextPage endCursor } } }"

func graphQLSpec(doc string) *connection.HTTPSpec {
	return &connection.HTTPSpec{
		Op: opRead, Method: "POST", PathTemplate: "/graphql", GraphQLDocument: doc,
		Params: []connection.Param{
			{Field: "term", In: connection.InBody, Wire: "term"},
			{Field: "first", In: connection.InBody, Wire: "first"},
			{Field: "after", In: connection.InBody, Wire: "after"},
		},
	}
}

func TestCallHTTP_GraphQL(t *testing.T) {
	t.Run("root field unwrapped", func(t *testing.T) {
		conn, s := callServer(t, nil, 200, `{"data":{"searchIssues":{"nodes":[{"id":"i-1","extra":1}],"pageInfo":{"hasNextPage":true,"endCursor":"c2"}}}}`)
		resp := &linearv1.SearchIssuesResponse{}
		err := connection.CallHTTP(context.Background(), conn, graphQLSpec(searchDoc), &linearv1.SearchIssuesRequest{Term: "synthetic", First: proto.Int32(10)}, resp)
		if err != nil {
			t.Fatalf("CallHTTP: %v", err)
		}
		if s.method != "POST" || s.path != "/graphql" {
			t.Fatalf("request = %s %s", s.method, s.path)
		}
		jsonEqual(t, s.body, `{"query":`+mustJSON(searchDoc)+`,"variables":{"term":"synthetic","first":10}}`)
		if len(resp.GetNodes()) != 1 || resp.GetNodes()[0].GetId() != "i-1" || resp.GetPageInfo().GetEndCursor() != "c2" {
			t.Fatalf("resp = %v", resp)
		}
	})
	t.Run("data matches the response", func(t *testing.T) {
		conn, _ := callServer(t, nil, 200, `{"data":{"issue":{"id":"i-9","identifier":"ENG-9"}}}`)
		spec := &connection.HTTPSpec{Op: opRead, Method: "POST", PathTemplate: "/graphql", GraphQLDocument: "query I($id: String!) { issue(id: $id) { id identifier } }",
			Params: []connection.Param{{Field: "id", In: connection.InBody, Wire: "id"}}}
		resp := &linearv1.GetIssueResponse{}
		if err := connection.CallHTTP(context.Background(), conn, spec, &linearv1.GetIssueRequest{Id: "i-9"}, resp); err != nil {
			t.Fatalf("CallHTTP: %v", err)
		}
		if resp.GetIssue().GetIdentifier() != "ENG-9" {
			t.Fatalf("resp = %v", resp)
		}
	})
}

// TestCallHTTP_NestedNames: nested message fields use JSON names
// (camelCase) inside GraphQL variables and wire names in a REST body.
func TestCallHTTP_NestedNames(t *testing.T) {
	req := &githubv1.PutContentsResponse{Content: &githubv1.ContentFile{HtmlUrl: "u"}}
	params := []connection.Param{{Field: "content", In: connection.InBody, Wire: "content"}}

	conn, s := callServer(t, nil, 200, `{"data":{}}`)
	gql := &connection.HTTPSpec{Op: opRead, Method: "POST", PathTemplate: "/graphql", GraphQLDocument: "query Q { q }", Params: params}
	if err := connection.CallHTTP(context.Background(), conn, gql, req, &githubv1.GetIssueResponse{}); err != nil {
		t.Fatalf("CallHTTP: %v", err)
	}
	jsonEqual(t, s.body, `{"query":"query Q { q }","variables":{"content":{"htmlUrl":"u"}}}`)

	conn, s = callServer(t, nil, 200, `{}`)
	rest := &connection.HTTPSpec{Op: opRead, Method: "POST", PathTemplate: "/x", Params: params}
	if err := connection.CallHTTP(context.Background(), conn, rest, req, &githubv1.GetIssueResponse{}); err != nil {
		t.Fatalf("CallHTTP: %v", err)
	}
	jsonEqual(t, s.body, `{"content":{"html_url":"u"}}`)
}

func TestCallHTTP_GraphQLErrors(t *testing.T) {
	cases := []struct {
		code      string
		kind      contracts.ErrorKind
		retryable bool
	}{
		{"NOT_FOUND", contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false},
		{"UNAUTHENTICATED", contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false},
		{"FORBIDDEN", contracts.ErrorKind_ERROR_KIND_FORBIDDEN, false},
		{"RATELIMITED", contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true},
		{"SOMETHING_ELSE", contracts.ErrorKind_ERROR_KIND_INVALID, false},
		{"", contracts.ErrorKind_ERROR_KIND_INVALID, false},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			body := `{"data":null,"errors":[{"message":"bad ` + syntheticToken + `","extensions":{"code":"` + tc.code + `"}}]}`
			conn, _ := callServer(t, nil, 200, body)
			err := connection.CallHTTP(context.Background(), conn, graphQLSpec(searchDoc), &linearv1.SearchIssuesRequest{Term: "x"}, &linearv1.SearchIssuesResponse{})
			var e *connection.Error
			if !errors.As(err, &e) {
				t.Fatalf("err = %v, want *connection.Error", err)
			}
			if e.Kind != tc.kind || e.Retryable != tc.retryable || e.Status != 200 {
				t.Fatalf("error = %+v", e)
			}
			if strings.Contains(err.Error(), syntheticToken) || strings.Contains(err.Error(), "bad") {
				t.Fatalf("error leaks the provider message: %v", err)
			}
		})
	}
}

func TestCallHTTP_DecodeErrorHasNoBody(t *testing.T) {
	conn, _ := callServer(t, nil, 200, `{"number": "`+syntheticToken+`"`)
	err := connection.CallHTTP(context.Background(), conn, &connection.HTTPSpec{Op: opRead, Method: "GET", PathTemplate: "/x"}, &githubv1.GetIssueRequest{}, &githubv1.GetIssueResponse{})
	if err == nil {
		t.Fatalf("want a decode error")
	}
	if strings.Contains(err.Error(), syntheticToken) {
		t.Fatalf("decode error leaks the body: %v", err)
	}
}

func TestCallHTTP_StatusErrorUsesTable(t *testing.T) {
	conn, _ := callServer(t, nil, 404, `{"message":"Not Found"}`)
	spec := &connection.HTTPSpec{Op: opRead, Method: "GET", PathTemplate: "/x", Errors: []*contracts.ErrorMapping{{Status: 404, Kind: contracts.ErrorKind_ERROR_KIND_NOT_FOUND}}}
	err := connection.CallHTTP(context.Background(), conn, spec, &githubv1.GetIssueRequest{}, &githubv1.GetIssueResponse{})
	var e *connection.Error
	if !errors.As(err, &e) || e.Kind != contracts.ErrorKind_ERROR_KIND_NOT_FOUND || e.OpID != opRead.ID {
		t.Fatalf("err = %v", err)
	}
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

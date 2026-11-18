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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// syntheticToken is a made-up credential; tests assert it never leaks.
const syntheticToken = "kei-synthetic-token-7f3a9c"

func TestRenderPath_EscapesSegments(t *testing.T) {
	cases := []struct {
		name     string
		template string
		values   map[string]string
		want     string
	}{
		{"plain", "/repos/{owner}/{repository}", map[string]string{"owner": "acme", "repository": "widgets"}, "/repos/acme/widgets"},
		{"slash and space", "/repos/{owner}/{repository}", map[string]string{"owner": "a/b", "repository": "x y"}, "/repos/a%2Fb/x%20y"},
		{"percent", "/repos/{owner}/{repository}", map[string]string{"owner": "100%", "repository": "w"}, "/repos/100%25/w"},
		{"dot dot", "/repos/{owner}/{repository}", map[string]string{"owner": "..", "repository": "."}, "/repos/%2E%2E/%2E"},
		{"query chars", "/repos/{owner}/{repository}", map[string]string{"owner": "a?b", "repository": "c#d"}, "/repos/a%3Fb/c%23d"},
		{"partial segment", "/files/{name}.json", map[string]string{"name": "a/b"}, "/files/a%2Fb.json"},
		{"number", "/repos/{owner}/{repository}/issues/{issue_number}", map[string]string{"owner": "acme", "repository": "widgets", "issue_number": "1347"}, "/repos/acme/widgets/issues/1347"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := connection.RenderPath(tc.template, tc.values)
			if err != nil {
				t.Fatalf("RenderPath: %v", err)
			}
			if got != tc.want {
				t.Fatalf("RenderPath = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRenderPath_MultiSegment: a multi-segment variable keeps "/" between
// segments but escapes each segment, and rejects ".", ".." and empty
// segments; other variables stay fully escaped.
func TestRenderPath_MultiSegment(t *testing.T) {
	const tmpl = "/repos/{owner}/{repository}/contents/{path}"
	render := func(path string) (string, error) {
		return connection.RenderPath(tmpl, map[string]string{"owner": "a/b", "repository": "w", "path": path}, "path")
	}
	good := map[string]string{
		"docs/notes.txt":  "/repos/a%2Fb/w/contents/docs/notes.txt",
		"notes.txt":       "/repos/a%2Fb/w/contents/notes.txt",
		"a b/c%d/e?f#g":   "/repos/a%2Fb/w/contents/a%20b/c%25d/e%3Ff%23g",
		"heads/feature/x": "/repos/a%2Fb/w/contents/heads/feature/x",
	}
	for in, want := range good {
		got, err := render(in)
		if err != nil {
			t.Fatalf("RenderPath(%q): %v", in, err)
		}
		if got != want {
			t.Fatalf("RenderPath(%q) = %q, want %q", in, got, want)
		}
	}
	for _, in := range []string{"a/../b", "a//b", "../x", "a/..", "./a", "a/./b", "/a", "a/", ""} {
		if got, err := render(in); err == nil {
			t.Errorf("RenderPath(%q) = %q, want an error", in, got)
		}
	}
	// A variable not named multi-segment still escapes "/".
	got, err := connection.RenderPath(tmpl, map[string]string{"owner": "o", "repository": "r", "path": "docs/notes.txt"})
	if err != nil || got != "/repos/o/r/contents/docs%2Fnotes.txt" {
		t.Fatalf("RenderPath without multi-segment = %q, %v", got, err)
	}
}

func TestRenderPath_Errors(t *testing.T) {
	cases := []struct {
		name     string
		template string
		values   map[string]string
	}{
		{"missing variable", "/repos/{owner}/{repository}", map[string]string{"owner": "acme"}},
		{"empty value", "/repos/{owner}/{repository}", map[string]string{"owner": "", "repository": "w"}},
		{"unterminated", "/repos/{owner", map[string]string{"owner": "acme"}},
		{"empty name", "/repos/{}", map[string]string{}},
		{"stray close", "/repos/owner}", map[string]string{}},
		{"not absolute", "repos/{owner}", map[string]string{"owner": "acme"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := connection.RenderPath(tc.template, tc.values); err == nil {
				t.Fatalf("RenderPath = %q, want an error", got)
			}
		})
	}
}

func TestHTTPBearerSetsHeader(t *testing.T) {
	var got *http.Request
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Link", `<https://api.example.test/x?page=2>; rel="next"`)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()

	conn, err := connection.NewHTTP(srv.URL+"/api", bearerConn(), syntheticToken, srv.Client())
	if err != nil {
		t.Fatalf("NewHTTP: %v", err)
	}
	status, header, body, err := conn.Do(context.Background(), http.MethodPost, "/repos/a%2Fb/x%20y", url.Values{"per_page": {"10"}}, []byte(`{"title":"t"}`))
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if status != http.StatusOK || string(body) != `{"ok":true}` {
		t.Fatalf("Do = %d %q", status, body)
	}
	if header.Get("Link") == "" {
		t.Fatalf("response header Link missing")
	}
	if a := got.Header.Get("Authorization"); a != "Bearer "+syntheticToken {
		t.Fatalf("Authorization = %q", a)
	}
	if got.Method != http.MethodPost {
		t.Fatalf("method = %s", got.Method)
	}
	if p := got.URL.EscapedPath(); p != "/api/repos/a%2Fb/x%20y" {
		t.Fatalf("escaped path = %q", p)
	}
	if q := got.URL.RawQuery; q != "per_page=10" {
		t.Fatalf("query = %q", q)
	}
	if ct := got.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if gotBody != `{"title":"t"}` {
		t.Fatalf("body = %q", gotBody)
	}
}

func TestHTTPHeaderAndBasic(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	t.Run("header default name", func(t *testing.T) {
		c := &contracts.Connection{Name: "api_key", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER}
		conn, err := connection.NewHTTP(srv.URL, c, syntheticToken, srv.Client())
		if err != nil {
			t.Fatalf("NewHTTP: %v", err)
		}
		if _, _, _, err := conn.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
			t.Fatalf("Do: %v", err)
		}
		if v := got.Get("X-Api-Key"); v != syntheticToken {
			t.Fatalf("X-Api-Key = %q", v)
		}
		if v := got.Get("Authorization"); v != "" {
			t.Fatalf("Authorization set for header kind: %q", v)
		}
	})

	t.Run("header from connection", func(t *testing.T) {
		c := &contracts.Connection{Name: "api_key", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER, HeaderName: "X-Service-Token"}
		conn, err := connection.NewHTTP(srv.URL, c, syntheticToken, srv.Client())
		if err != nil {
			t.Fatalf("NewHTTP: %v", err)
		}
		if _, _, _, err := conn.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
			t.Fatalf("Do: %v", err)
		}
		if v := got.Get("X-Service-Token"); v != syntheticToken {
			t.Fatalf("X-Service-Token = %q", v)
		}
		if v := got.Get("X-Api-Key"); v != "" {
			t.Fatalf("default header also set: %q", v)
		}
	})

	t.Run("basic", func(t *testing.T) {
		c := &contracts.Connection{Name: "basic", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BASIC}
		conn, err := connection.NewHTTP(srv.URL, c, "svc-user:"+syntheticToken, srv.Client())
		if err != nil {
			t.Fatalf("NewHTTP: %v", err)
		}
		if _, _, _, err := conn.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
			t.Fatalf("Do: %v", err)
		}
		req := &http.Request{Header: got}
		user, pass, ok := req.BasicAuth()
		if !ok || user != "svc-user" || pass != syntheticToken {
			t.Fatalf("BasicAuth = %q %q %v", user, pass, ok)
		}
	})
}

func TestNewHTTP_Rejects(t *testing.T) {
	client := http.DefaultClient
	cases := []struct {
		name       string
		base       string
		conn       *contracts.Connection
		credential string
		client     *http.Client
	}{
		{"nil connection", "https://api.example.test", nil, syntheticToken, client},
		{"sql kind", "https://api.example.test", &contracts.Connection{Kind: contracts.ConnectionKind_CONNECTION_KIND_SQL_DSN}, syntheticToken, client},
		{"unspecified kind", "https://api.example.test", &contracts.Connection{}, syntheticToken, client},
		{"empty credential", "https://api.example.test", bearerConn(), "", client},
		{"newline in credential", "https://api.example.test", bearerConn(), syntheticToken + "\r\nX: y", client},
		{"bad header name", "https://api.example.test", &contracts.Connection{Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER, HeaderName: "X Bad:"}, syntheticToken, client},
		{"basic without colon", "https://api.example.test", &contracts.Connection{Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BASIC}, syntheticToken, client},
		{"nil client", "https://api.example.test", bearerConn(), syntheticToken, nil},
		{"relative base", "/api", bearerConn(), syntheticToken, client},
		{"bad scheme", "ftp://api.example.test", bearerConn(), syntheticToken, client},
		{"userinfo in base", (&url.URL{Scheme: "https", User: url.UserPassword("u", "p"), Host: "api.example.test"}).String(), bearerConn(), syntheticToken, client},
		{"query in base", "https://api.example.test/?a=b", bearerConn(), syntheticToken, client},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := connection.NewHTTP(tc.base, tc.conn, tc.credential, tc.client)
			if err == nil {
				t.Fatalf("NewHTTP: want an error")
			}
			if strings.Contains(err.Error(), syntheticToken) {
				t.Fatalf("error leaks the credential: %v", err)
			}
		})
	}
}

func TestHTTPDo_PathCannotChangeOrigin(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()
	conn, err := connection.NewHTTP(srv.URL, bearerConn(), syntheticToken, srv.Client())
	if err != nil {
		t.Fatalf("NewHTTP: %v", err)
	}
	for _, p := range []string{"//evil.example.test/x", "x", "", "http://evil.example.test/", "/a?b=c", "/a#b"} {
		if _, _, _, err := conn.Do(context.Background(), http.MethodGet, p, nil, nil); err == nil {
			t.Errorf("Do(%q): want an error", p)
		}
	}
	if hits != 0 {
		t.Fatalf("server got %d requests, want 0", hits)
	}
}

func TestConnection_NoCredentialInErrors(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusInternalServerError} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// A hostile or careless provider echoes the credential back.
				w.WriteHeader(status)
				_, _ = io.WriteString(w, `{"message":"bad token `+r.Header.Get("Authorization")+`"}`)
			}))
			defer srv.Close()

			conn, err := connection.NewHTTP(srv.URL, bearerConn(), syntheticToken, srv.Client())
			if err != nil {
				t.Fatalf("NewHTTP: %v", err)
			}
			op := connection.Op{ID: "github.issues.get", Class: connection.Read}
			got, _, _, err := conn.Do(context.Background(), http.MethodGet, "/x", nil, nil)
			if err != nil {
				t.Fatalf("Do: %v", err)
			}
			if got != status {
				t.Fatalf("status = %d", got)
			}
			mapped := connection.MapError(op, got, []*contracts.ErrorMapping{
				{Status: 401, Kind: contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED},
			})
			if mapped == nil {
				t.Fatalf("MapError(%d) = nil", status)
			}
			for _, s := range []string{
				mapped.Error(),
				fmt.Sprintf("%v", mapped), fmt.Sprintf("%+v", mapped), fmt.Sprintf("%#v", mapped),
				fmt.Sprintf("%v", conn), fmt.Sprintf("%+v", conn), fmt.Sprintf("%#v", conn), fmt.Sprintf("%s", conn),
			} {
				if strings.Contains(s, syntheticToken) {
					t.Fatalf("credential leaked: %s", s)
				}
				if strings.Contains(s, "bad token") {
					t.Fatalf("response body leaked: %s", s)
				}
			}
		})
	}

	t.Run("transport error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		base := srv.URL
		srv.Close() // connection refused from now on
		conn, err := connection.NewHTTP(base, bearerConn(), syntheticToken, http.DefaultClient)
		if err != nil {
			t.Fatalf("NewHTTP: %v", err)
		}
		_, _, _, err = conn.Do(context.Background(), http.MethodGet, "/x", nil, nil)
		if err == nil {
			t.Fatalf("Do: want an error")
		}
		if strings.Contains(err.Error(), syntheticToken) {
			t.Fatalf("credential leaked: %v", err)
		}
	})
}

func TestMapError_Retryable(t *testing.T) {
	op := connection.Op{ID: "github.issues.get", Class: connection.Read}
	table := []*contracts.ErrorMapping{
		{Status: 404, Kind: contracts.ErrorKind_ERROR_KIND_NOT_FOUND},
		{Status: 429, Kind: contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, Retryable: true},
	}
	err := connection.MapError(op, 429, table)
	var e *connection.Error
	if !errors.As(err, &e) {
		t.Fatalf("MapError = %T, want *connection.Error", err)
	}
	if e.Kind != contracts.ErrorKind_ERROR_KIND_RATE_LIMITED || !e.Retryable || e.Status != 429 || e.OpID != "github.issues.get" {
		t.Fatalf("MapError = %+v", e)
	}
	if !strings.Contains(e.Error(), "github.issues.get") || !strings.Contains(e.Error(), "429") {
		t.Fatalf("Error() = %q, want op id and status", e.Error())
	}
}

func TestMapError_Defaults(t *testing.T) {
	op := connection.Op{ID: "github.issues.get", Class: connection.Read}
	for _, ok := range []int{200, 201, 204, 299} {
		if err := connection.MapError(op, ok, nil); err != nil {
			t.Errorf("MapError(%d) = %v, want nil", ok, err)
		}
	}
	cases := []struct {
		status    int
		kind      contracts.ErrorKind
		retryable bool
	}{
		{400, contracts.ErrorKind_ERROR_KIND_INVALID, false},
		{401, contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false},
		{403, contracts.ErrorKind_ERROR_KIND_FORBIDDEN, false},
		{404, contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false},
		{409, contracts.ErrorKind_ERROR_KIND_CONFLICT, false},
		{422, contracts.ErrorKind_ERROR_KIND_INVALID, false},
		{429, contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true},
		{500, contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, false},
		{503, contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, true},
		{302, contracts.ErrorKind_ERROR_KIND_UNSPECIFIED, false},
	}
	for _, tc := range cases {
		var e *connection.Error
		if !errors.As(connection.MapError(op, tc.status, nil), &e) {
			t.Fatalf("MapError(%d): not a *connection.Error", tc.status)
		}
		if e.Kind != tc.kind || e.Retryable != tc.retryable {
			t.Errorf("MapError(%d) = %v retryable=%v, want %v retryable=%v", tc.status, e.Kind, e.Retryable, tc.kind, tc.retryable)
		}
	}
	// The declared table wins over the defaults.
	var e *connection.Error
	errors.As(connection.MapError(op, 500, []*contracts.ErrorMapping{{Status: 500, Kind: contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, Retryable: true}}), &e)
	if !e.Retryable {
		t.Fatalf("table mapping ignored: %+v", e)
	}
}

func TestOpClassString(t *testing.T) {
	if connection.Read.String() != "read" || connection.Write.String() != "write" {
		t.Fatalf("Class.String = %q %q", connection.Read, connection.Write)
	}
	var zero connection.Class
	if zero == connection.Read || zero == connection.Write {
		t.Fatalf("the zero Class must be neither read nor write")
	}
}

func bearerConn() *contracts.Connection {
	return &contracts.Connection{Name: "pat", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER}
}

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

package mockapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Options configures a Server.
type Options struct {
	// ManifestPath is the connector manifest (manifest/connector-manifest.v1.json).
	ManifestPath string
	// FixturesDir is fixtures/operations.
	FixturesDir string
	// Seed offsets the synthetic ids of generated pages; 0 means 1.
	Seed int
	// FailDelay is how long the "slow" failure directive sleeps, in seconds.
	// 0 means 2.
	FailDelay int
}

// Server is an http.Handler serving the contract-derived provider mock.
type Server struct {
	routes  []*route
	graphQL map[string]*route // GraphQL root field
	opts    Options
}

// Routes returns the public view of the route table.
func (s *Server) Routes() []Route { return makeRoutes(s.routes) }

// New builds a Server from the manifest and the linked descriptors. It
// fails on contracts drift so a stale manifest or a missing generated
// import surfaces at construction, not as a mystery 404.
func New(opts Options) (*Server, error) {
	if opts.Seed == 0 {
		opts.Seed = 1
	}
	if opts.FailDelay == 0 {
		opts.FailDelay = 2
	}
	raw, err := os.ReadFile(opts.ManifestPath)
	if err != nil {
		return nil, fmt.Errorf("mockapi: read manifest: %w", err)
	}
	routes, err := buildRoutes(raw)
	if err != nil {
		return nil, err
	}
	s := &Server{routes: routes, graphQL: map[string]*route{}, opts: opts}
	for _, r := range routes {
		if r.RootField != "" {
			s.graphQL[r.RootField] = r
		}
	}
	return s, nil
}

// reqState is what ServeHTTP learns about one request.
type reqState struct {
	route *route
	body  []byte                     // raw request body (may be empty)
	vars  map[string]json.RawMessage // GraphQL variables
}

// ServeHTTP matches the request to a route, validates it, applies failure
// injection and serves the fixture.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	st := &reqState{body: body}
	if reason := s.matchRoute(r, st); reason != "" {
		if st.route == nil && reason == "method" {
			mockFail(w, http.StatusMethodNotAllowed, "mockapi: "+r.URL.EscapedPath()+" does not accept "+r.Method)
			return
		}
		mockFail(w, http.StatusNotFound, reason)
		return
	}
	route := st.route
	if reason := validateCredential(route, r); reason != "" {
		mockFail(w, http.StatusBadRequest, reason)
		return
	}
	if route.RootField == "" {
		if reason := validatePath(route, r); reason != "" {
			mockFail(w, http.StatusBadRequest, reason)
			return
		}
	}
	if fail := r.Header.Get("X-Mock-Fail"); fail != "" {
		if s.applyFailure(w, fail) {
			return
		}
	}
	s.serve(w, r, st)
}

// matchRoute resolves the route. Empty reason = matched; "method" = the
// path matches but the method does not; otherwise the reason is a 404 body.
func (s *Server) matchRoute(r *http.Request, st *reqState) string {
	if r.URL.Path == "/graphql" {
		if r.Method != http.MethodPost {
			return "method"
		}
		var body struct {
			Query     string                     `json:"query"`
			Variables map[string]json.RawMessage `json:"variables"`
		}
		if err := json.Unmarshal(st.body, &body); err != nil || body.Query == "" {
			return "mockapi: request body must be {\"query\": <document>, \"variables\": {...}}"
		}
		st.vars = body.Variables
		m := rootFieldRe.FindStringSubmatch(body.Query)
		if m == nil {
			return "mockapi: the GraphQL document has no root field"
		}
		route := s.graphQL[m[1]]
		if route == nil {
			return "mockapi: no route for the GraphQL root field " + m[1]
		}
		st.route = route
		return ""
	}
	route, methodOnly := s.matchREST(r)
	if route != nil {
		st.route = route
		return ""
	}
	if methodOnly {
		return "method"
	}
	return "mockapi: no route for " + r.Method + " " + r.URL.EscapedPath()
}

// matchREST finds the REST route whose compiled template matches the path;
// (nil, true) means some route's template matched but none with the right
// method, so the caller reports 405.
func (s *Server) matchREST(r *http.Request) (*route, bool) {
	escaped := r.URL.EscapedPath()
	methodOnly := false
	for _, route := range s.routes {
		if route.RootField != "" {
			continue
		}
		if !route.re.MatchString(escaped) {
			continue
		}
		if route.Method != r.Method {
			methodOnly = true
			continue
		}
		return route, false
	}
	return nil, methodOnly
}

// serve writes the operation's response: the fixture verbatim for a plain
// HTTP route, wrapped in {"data": …} for GraphQL, with seeded pagination
// when the operation has cursor fields.
func (s *Server) serve(w http.ResponseWriter, r *http.Request, st *reqState) {
	route := st.route
	fixture, err := os.ReadFile(s.fixtureFile(route))
	if err != nil {
		mockFail(w, http.StatusInternalServerError, "mockapi: fixture for "+route.OpID+": "+err.Error())
		return
	}
	body := fixture
	if route.cursorWire != "" {
		page, reason := s.cursorPage(route, r, st)
		if reason != "" {
			mockFail(w, http.StatusBadRequest, reason)
			return
		}
		if page == 0 {
			body = firstPage(fixture, route)
		} else {
			body = lastPage(fixture, route, s.opts.Seed)
		}
	}
	if strings.EqualFold(r.Header.Get("X-Mock-Fail"), "truncate") {
		body = body[:len(body)*3/5]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if route.RootField != "" {
		data, _ := json.Marshal(map[string]json.RawMessage{"data": json.RawMessage(body)})
		_, _ = w.Write(data)
		return
	}
	_, _ = w.Write(body)
}

// fixtureFile maps an operation id to its response fixture:
// github.issues.get -> fixtures/operations/github/issues.get.response.json.
func (s *Server) fixtureFile(route *route) string {
	suffix := strings.TrimPrefix(route.OpID, route.connector.GetId()+".")
	return s.opts.FixturesDir + "/" + route.connector.GetId() + "/" + suffix + ".response.json"
}

// cursorPage extracts the page cursor from the request and maps it to a
// page number: 0 = first page (no cursor), anything else = the last page.
func (s *Server) cursorPage(route *route, r *http.Request, st *reqState) (int, string) {
	raw := ""
	switch {
	case route.RootField != "":
		if v, ok := st.vars[route.cursorWire]; ok && string(v) != "null" {
			raw = strings.Trim(string(v), `"`)
		}
	case route.cursorInBody:
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(st.body, &obj); err != nil {
			return 0, "mockapi: request body is not a JSON object"
		}
		if v, ok := obj[route.cursorWire]; ok && string(v) != "null" {
			raw = strings.Trim(string(v), `"`)
		}
	default:
		raw = r.URL.Query().Get(route.cursorWire)
	}
	if raw == "" {
		return 0, ""
	}
	if route.cursorIsInt {
		if _, err := strconv.Atoi(raw); err != nil {
			return 0, "mockapi: cursor " + raw + " is malformed: want a page number"
		}
		return 1, ""
	}
	if !strings.HasPrefix(raw, "mock-page-") {
		return 0, "mockapi: cursor " + raw + " is malformed: want a mock-page-<n> token"
	}
	return 1, ""
}

// firstPage rewrites the fixture's next cursor to the mock's token.
func firstPage(fixture []byte, route *route) []byte {
	var obj map[string]any
	if err := json.Unmarshal(fixture, &obj); err != nil {
		return fixture
	}
	if route.cursorIsInt {
		obj[route.nextWire] = 2
	} else {
		obj[route.nextWire] = "mock-page-2"
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return fixture
	}
	return out
}

// lastPage serves a varied, cursor-free page: the seeded variation of the
// fixture with the next-cursor key removed and the list one item longer.
func lastPage(fixture []byte, route *route, seed int) []byte {
	var doc any
	if err := json.Unmarshal(fixture, &doc); err != nil {
		return fixture
	}
	v := vary(growList(doc), seed*1000+1)
	if obj, ok := v.(map[string]any); ok {
		delete(obj, route.nextWire)
		fixture, err := json.Marshal(obj)
		if err != nil {
			return fixture
		}
		return fixture
	}
	out, err := json.Marshal(v)
	if err != nil {
		return fixture
	}
	return out
}

// applyFailure handles the X-Mock-Fail directive; true = responded.
func (s *Server) applyFailure(w http.ResponseWriter, fail string) bool {
	switch strings.ToLower(fail) {
	case "401", "403", "404", "429", "500":
		status, _ := strconv.Atoi(fail)
		if fail == "429" {
			w.Header().Set("Retry-After", "1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{}`))
		return true
	case "slow":
		sleepFor(s.opts.FailDelay)
	case "truncate":
		// Handled in serve: shorten the body.
	}
	return false
}

// mockFail writes a JSON error body with a reason.
func mockFail(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(map[string]string{"error": reason})
	_, _ = w.Write(b)
}

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

// Cross-language parity (proposal 0002, section 11): every HTTP and GraphQL
// operation in the manifest goes through the generated operations of each
// language against one recording server. The rendered requests must be
// identical and the decoded results, as canonical JSON, equal.
//
// Each language has a runner behind scripts/parity/run.sh. A runner that is
// not available prints "SKIP: <reason>" and exits 3: that fails the test for
// a language in requiredParityLanguages and is logged for the others.

package conformance

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// parityLanguages are every language with generated operations.
var parityLanguages = []string{"go", "ts", "python", "rust"}

// requiredParityLanguages must run; an absent runner fails the test. All
// four are listed, so no language can silently skip.
var requiredParityLanguages = []string{"go", "ts", "python", "rust"}

const parityToken = "kei-parity-token-0001"

var (
	manifestPath = filepath.Join("..", "..", "..", "manifest", "connector-manifest.v1.json")
	runScript    = filepath.Join("..", "..", "..", "scripts", "parity", "run.sh")
)

// parityCase is one operation call, as the runners read it.
type parityCase struct {
	Op           string          `json:"op"`
	Connector    string          `json:"connector"`
	Method       string          `json:"method"`
	RequestType  string          `json:"request_type"`
	ResponseType string          `json:"response_type"`
	Request      json.RawMessage `json:"request"`
}

// parityResult is one line a runner prints.
type parityResult struct {
	Op       string          `json:"op"`
	Response json.RawMessage `json:"response,omitempty"`
	Error    string          `json:"error,omitempty"`
}

// parityRequest is what the server saw.
type parityRequest struct {
	Method, Path, Body, Authorization, ContentType string
	Query                                          url.Values
}

type parityServer struct {
	mu       sync.Mutex
	seen     map[string]parityRequest // lang + " " + op
	response map[string][]byte        // op -> body to serve
}

func (s *parityServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/", 3)
	if len(parts) < 2 {
		http.Error(w, "want /<lang>/<op>/...", http.StatusBadRequest)
		return
	}
	lang, op := parts[0], parts[1]
	rest := "/"
	if len(parts) == 3 {
		rest += parts[2]
	}
	body, _ := io.ReadAll(r.Body)
	query, _ := url.ParseQuery(r.URL.RawQuery)
	s.mu.Lock()
	s.seen[lang+" "+op] = parityRequest{
		Method: r.Method, Path: rest, Body: string(body), Query: query,
		Authorization: r.Header.Get("Authorization"), ContentType: r.Header.Get("Content-Type"),
	}
	resp, ok := s.response[op]
	s.mu.Unlock()
	if !ok {
		http.Error(w, "unknown operation", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func TestParity_AllLanguages(t *testing.T) {
	if os.Getenv("KEI_PARITY") == "" {
		t.Skip("set KEI_PARITY=1 to run the cross-language parity test (the parity CI job does)")
	}
	cases, served := parityCases(t)
	srv := &parityServer{seen: map[string]parityRequest{}, response: served}
	ts := httptest.NewServer(srv)
	defer ts.Close()

	casesPath := filepath.Join(t.TempDir(), "cases.json")
	b, err := json.Marshal(cases)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(casesPath, b, 0o600); err != nil {
		t.Fatal(err)
	}

	results := map[string]map[string]parityResult{}
	for _, lang := range parityLanguages {
		out, skipped, err := runParity(lang, casesPath, ts.URL+"/"+lang)
		required := contains(requiredParityLanguages, lang)
		switch {
		case err != nil:
			t.Fatalf("%s runner: %v", lang, err)
		case skipped != "" && required:
			t.Errorf("%s runner is required but skipped: %s", lang, skipped)
			continue
		case skipped != "":
			t.Logf("%s runner skipped: %s", lang, skipped)
			continue
		}
		results[lang] = out
	}
	ref, ok := results["go"]
	if !ok {
		t.Fatalf("the Go reference runner did not run")
	}
	for _, c := range cases {
		r := ref[c.Op]
		if r.Error != "" || r.Response == nil {
			t.Errorf("go %s: %s", c.Op, r.Error)
		}
	}
	// Review Focus 2: an id above 2^53 survives, as protobuf JSON's string.
	if !strings.Contains(string(ref["github.issues.get"].Response), `"id":"9007199254740993"`) {
		t.Errorf("go github.issues.get = %s, want the 2^53+1 id", ref["github.issues.get"].Response)
	}

	langs := make([]string, 0, len(results))
	for lang := range results {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	for _, lang := range langs {
		if lang == "go" {
			continue
		}
		for _, c := range cases {
			want, got := srv.seen["go "+c.Op], srv.seen[lang+" "+c.Op]
			if !reflect.DeepEqual(want, got) {
				t.Errorf("%s %s request differs from go:\n go: %+v\n%3s: %+v", lang, c.Op, want, lang, got)
			}
			r, ok := results[lang][c.Op]
			switch {
			case !ok:
				t.Errorf("%s runner has no result for %s", lang, c.Op)
			case r.Error != "":
				t.Errorf("%s %s: %s", lang, c.Op, r.Error)
			case canonicalJSON(t, r.Response) != canonicalJSON(t, ref[c.Op].Response):
				t.Errorf("%s %s result differs from go:\n go: %s\n%3s: %s", lang, c.Op,
					canonicalJSON(t, ref[c.Op].Response), lang, canonicalJSON(t, r.Response))
			}
		}
	}
}

// parityCases builds one case per HTTP or GraphQL capability in the
// manifest, and the provider-shaped response the server returns for it.
func parityCases(t *testing.T) ([]parityCase, map[string][]byte) {
	t.Helper()
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("reading the manifest: %v", err)
	}
	var m struct {
		Connectors []struct {
			ID          string `json:"id"`
			Connections []struct {
				Kind string `json:"kind"`
			} `json:"connections"`
			Capabilities []struct {
				ID       string `json:"id"`
				Request  string `json:"request"`
				Response string `json:"response"`
			} `json:"capabilities"`
		} `json:"connectors"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	var cases []parityCase
	served := map[string][]byte{}
	for _, c := range m.Connectors {
		if len(c.Connections) > 0 && c.Connections[0].Kind == "sql_dsn" {
			continue // SQL has no cross-language test server; replay covers it per language
		}
		for _, capability := range c.Capabilities {
			op, ok := contracts.Lookup(capability.ID)
			if !ok {
				t.Fatalf("%s is in the manifest but not registered", capability.ID)
			}
			method := op.GetMethodName()
			if method == "" {
				name := capability.Request[strings.LastIndex(capability.Request, ".")+1:]
				method = strings.TrimSuffix(name, "Request")
			}
			cases = append(cases, parityCase{
				Op: capability.ID, Connector: c.ID, Method: method,
				RequestType: capability.Request, ResponseType: capability.Response,
				Request: readFixture(t, c.ID, opFixtureName(capability.ID, "request")),
			})
			served[capability.ID] = providerShape(t, c.ID, op, capability.Response,
				readFixture(t, c.ID, opFixtureName(capability.ID, "response")))
		}
	}
	if len(cases) == 0 {
		t.Fatal("no parity cases")
	}
	return cases, served
}

var graphQLRootRe = regexp.MustCompile(`^\s*(?:query|mutation)\b[^{]*\{\s*([A-Za-z_][A-Za-z0-9_]*)`)

// providerShape turns a response fixture (the response message's JSON)
// into what the provider sends: GraphQL's data envelope (under the root
// field when the fixture is that field's value), or a bare array for a
// GitHub list whose response has one repeated field.
func providerShape(t *testing.T, connector string, op *contracts.Operation, responseType string, fixture []byte) []byte {
	t.Helper()
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(fixture, &obj); err != nil {
		t.Fatalf("%s response fixture: %v", op.GetId(), err)
	}
	if doc := op.GetGraphqlDocument(); doc != "" {
		m := graphQLRootRe.FindStringSubmatch(doc)
		if m == nil {
			t.Fatalf("%s: cannot find the GraphQL root field", op.GetId())
		}
		data := json.RawMessage(fixture)
		if _, ok := obj[m[1]]; !ok {
			data, _ = json.Marshal(map[string]json.RawMessage{m[1]: fixture})
		}
		out, _ := json.Marshal(map[string]json.RawMessage{"data": data})
		return out
	}
	if connector == "github" {
		mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(responseType))
		if err != nil {
			t.Fatalf("%s: %v", responseType, err)
		}
		var lists []protoreflect.FieldDescriptor
		fields := mt.Descriptor().Fields()
		for i := 0; i < fields.Len(); i++ {
			if fields.Get(i).IsList() {
				lists = append(lists, fields.Get(i))
			}
		}
		if len(lists) == 1 {
			if v, ok := obj[string(lists[0].Name())]; ok {
				return v
			}
			if v, ok := obj[lists[0].JSONName()]; ok {
				return v
			}
		}
	}
	return fixture
}

// runParity runs one language's runner. skipped is the reason when the
// runner is not available (exit status 3).
func runParity(lang, casesPath, baseURL string) (map[string]parityResult, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", runScript, lang, casesPath, baseURL, parityToken)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 3 {
		return nil, strings.TrimSpace(strings.TrimPrefix(stdout.String(), "SKIP:")), nil
	}
	if err != nil {
		return nil, "", errors.New(err.Error() + ": " + stderr.String())
	}
	out := map[string]parityResult{}
	sc := bufio.NewScanner(&stdout)
	sc.Buffer(make([]byte, 0, 1<<20), 16<<20)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var r parityResult
		if err := json.Unmarshal(line, &r); err != nil {
			return nil, "", errors.New("bad runner line " + strconv.Quote(string(line)) + ": " + err.Error())
		}
		out[r.Op] = r
	}
	return out, "", sc.Err()
}

// canonicalJSON re-encodes a protobuf JSON result with sorted keys and
// numbers normalized, so formatting differences between languages' JSON
// encoders do not count.
func canonicalJSON(t *testing.T, raw json.RawMessage) string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return "invalid JSON: " + string(raw)
	}
	b, err := json.Marshal(normalizeNumbers(v))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func normalizeNumbers(v any) any {
	switch x := v.(type) {
	case json.Number:
		if f, err := strconv.ParseFloat(x.String(), 64); err == nil {
			return json.Number(strconv.FormatFloat(f, 'g', -1, 64))
		}
		return x
	case []any:
		for i := range x {
			x[i] = normalizeNumbers(x[i])
		}
	case map[string]any:
		for k := range x {
			x[k] = normalizeNumbers(x[k])
		}
	}
	return v
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

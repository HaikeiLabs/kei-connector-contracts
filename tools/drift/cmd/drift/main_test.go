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

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testManifest = `{
  "schema": "kei.connector-manifest/v1",
  "contracts_version": "v0.19.0",
  "connectors": [
    {
      "id": "github",
      "base_url": "https://api.github.com",
      "connections": [
        {
          "name": "pat",
          "kind": "http_bearer",
          "auth_method": "secret_ref",
          "account_model": "shared",
          "scopes": [],
          "setup_fields": [
            {"name": "token", "type": "string", "pattern": "^\\S+$", "secret": true}
          ]
        }
      ],
      "capabilities": [
        {"id": "github.issues.list", "class": "read"}
      ]
    }
  ]
}`

const testRequests = `{
  "requests": [
    {"op": "github.issues.list", "args": {"owner": "HaikeiLabs", "repository": "kei-connector-contracts", "state": "open", "per_page": 5}}
  ]
}`

const testToken = "supersecret-sandbox-token-value"

func writeTestFiles(t *testing.T) (manifestPath, requestsDir string) {
	t.Helper()
	dir := t.TempDir()
	manifestPath = filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(testManifest), 0o600); err != nil {
		t.Fatal(err)
	}
	requestsDir = filepath.Join(dir, "drift")
	if err := os.MkdirAll(requestsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(requestsDir, "github.json"), []byte(testRequests), 0o600); err != nil {
		t.Fatal(err)
	}
	return manifestPath, requestsDir
}

// githubListIssuesFixture is the canonical fixture body for github.issues.list.
const githubListIssuesFixture = `{
  "issues": [
    {"id": 5001, "number": 12, "title": "Synthetic first issue", "state": "open", "body": null, "html_url": "https://example.invalid/example-org/example-repo/issues/12"}
  ],
  "incomplete_results": false
}`

func newFixtureServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestSkipsWithoutCreds(t *testing.T) {
	manifestPath, requestsDir := writeTestFiles(t)
	var out strings.Builder
	code := Run(t.Context(), Options{
		ManifestPath: manifestPath,
		RequestsDir:  requestsDir,
		Env:          map[string]string{}, // no KEI_SANDBOX_GITHUB_TOKEN
		Out:          &out,
	})
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output:\n%s", code, out.String())
	}
	if got := out.String(); !strings.Contains(got, "SKIP github: no sandbox credentials") {
		t.Fatalf("output missing skip line:\n%s", got)
	}
	// The output must never carry internal ticket ids (public repo).
	if got := out.String(); strings.Contains(got, "HAI-") {
		t.Fatalf("output carries an internal ticket id:\n%s", got)
	}
}

func TestReportsUnknownFieldsAsDrift(t *testing.T) {
	manifestPath, requestsDir := writeTestFiles(t)
	// The canonical fixture plus one more unknown field: strict unmarshal
	// must flag drift. The canonical fixture already carries its own drift
	// gate (incomplete_results is not modeled by the contract), so the first
	// reported unknown field is deterministic.
	drifted := `{
  "issues": [
    {"id": 5001, "number": 12, "title": "Synthetic first issue", "state": "open", "body": null, "html_url": "https://example.invalid/example-org/example-repo/issues/12"}
  ],
  "incomplete_results": false,
  "extra_field": "unexpected"
}`
	srv := newFixtureServer(t, drifted)
	var out strings.Builder
	code := Run(t.Context(), Options{
		ManifestPath: manifestPath,
		RequestsDir:  requestsDir,
		Env:          map[string]string{"KEI_SANDBOX_GITHUB_TOKEN": testToken},
		BaseURLs:     map[string]string{"github": srv.URL},
		Out:          &out,
	})
	if code != 1 {
		t.Fatalf("exit = %d, want 1; output:\n%s", code, out.String())
	}
	if got := out.String(); !strings.Contains(got, "DRIFT github.issues.list") || !strings.Contains(got, `unknown field "incomplete_results"`) {
		t.Fatalf("output missing drift line:\n%s", got)
	}
}

func TestNeverPrintsCredential(t *testing.T) {
	manifestPath, requestsDir := writeTestFiles(t)
	// The provider answers 401 and echoes the (invalid) token in the error
	// body: no matter what the tool does internally, the token value must
	// never reach stdout/stderr.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message": "Bad credentials for ` + testToken + `"}`))
	}))
	t.Cleanup(srv.Close)
	var out strings.Builder
	code := Run(t.Context(), Options{
		ManifestPath: manifestPath,
		RequestsDir:  requestsDir,
		Env:          map[string]string{"KEI_SANDBOX_GITHUB_TOKEN": testToken},
		BaseURLs:     map[string]string{"github": srv.URL},
		Out:          &out,
	})
	if code != 1 {
		t.Fatalf("exit = %d, want 1; output:\n%s", code, out.String())
	}
	if got := out.String(); strings.Contains(got, testToken) {
		t.Fatalf("output leaks the credential:\n%s", got)
	}
}

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

// Result-required checks: a write operation whose result_required paths are
// not all held by a 2xx response fails with ERROR_KIND_CONTRACT_VIOLATION,
// not with a silent success. Mirrors the contract in
// docs/operation-contracts.md.

package conformance

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1/linearops"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// badResultFixtures are the response fixtures that each miss one
// result_required path of linear.issues.create.
var badResultFixtures = []struct {
	name    string
	fixture string
}{
	{"success false", "issues.create.response.success-false.json"},
	{"issue missing", "issues.create.response.no-issue.json"},
	{"empty identifier", "issues.create.response.empty-identifier.json"},
}

// createIssue sends the operation's request fixture.
func createIssue(t *testing.T, conn connection.HTTP) (*linearv1.CreateIssueResponse, error) {
	t.Helper()
	req := &linearv1.CreateIssueRequest{}
	loadRequest(t, "linear", "linear.issues.create", req)
	return linearops.New(conn).CreateIssue(context.Background(), req)
}

// graphqlEnvelope wraps a linear response fixture in GraphQL's data
// envelope, under the issueCreate root field.
func graphqlEnvelope(t *testing.T, root string, fixture []byte) []byte {
	t.Helper()
	data, err := json.Marshal(map[string]json.RawMessage{root: fixture})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := json.Marshal(map[string]json.RawMessage{"data": data})
	if err != nil {
		t.Fatal(err)
	}
	return envelope
}

// TestGeneratedGraphQL_ResultRequired rejects each 2xx response that misses
// a result_required path with ERROR_KIND_CONTRACT_VIOLATION.
func TestGeneratedGraphQL_ResultRequired(t *testing.T) {
	for _, tc := range badResultFixtures {
		t.Run(tc.name, func(t *testing.T) {
			envelope := graphqlEnvelope(t, "issueCreate", readFixture(t, "linear", tc.fixture))
			srv, _ := replayServer(t, http.StatusOK, nil, envelope)
			_, err := createIssue(t, bearer(t, srv.URL))
			var e *connection.Error
			if !errors.As(err, &e) {
				t.Fatalf("CreateIssue error = %v, want *connection.Error", err)
			}
			if e.Kind != contracts.ErrorKind_ERROR_KIND_CONTRACT_VIOLATION {
				t.Fatalf("kind = %s, want ERROR_KIND_CONTRACT_VIOLATION", e.Kind)
			}
			if e.OpID != "linear.issues.create" || e.Status != http.StatusOK {
				t.Fatalf("error = %+v, want op linear.issues.create at status 200", e)
			}
		})
	}
}

// TestGeneratedGraphQL_ResultRequiredGood keeps the good fixture passing:
// the check only fires when a path is actually missing.
func TestGeneratedGraphQL_ResultRequiredGood(t *testing.T) {
	envelope := graphqlEnvelope(t, "issueCreate", readFixture(t, "linear", opFixtureName("linear.issues.create", "response")))
	srv, _ := replayServer(t, http.StatusOK, nil, envelope)
	resp, err := createIssue(t, bearer(t, srv.URL))
	if err != nil {
		t.Fatalf("CreateIssue: %v", err)
	}
	if !resp.GetSuccess() || resp.GetIssue().GetIdentifier() != "ENG-102" {
		t.Fatalf("resp = %+v", resp)
	}
}

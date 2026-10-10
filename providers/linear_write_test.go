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
	"regexp"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// linearScriptTransport answers each GraphQL operation by its operation name
// and records every request in order.
type linearScriptTransport struct {
	answers  map[string]json.RawMessage
	requests []linearScriptRequest
}

type linearScriptRequest struct {
	operation string
	vars      map[string]any
}

var linearOperation = regexp.MustCompile(`^(?:query|mutation) (\w+)`)

func (t *linearScriptTransport) Query(_ context.Context, endpoint, token, query string, variables map[string]any) (json.RawMessage, error) {
	if endpoint != linearGraphQLEndpoint || token != "linear-access-token" {
		return nil, errors.New("bad request")
	}
	op := linearOperation.FindStringSubmatch(query)[1]
	t.requests = append(t.requests, linearScriptRequest{operation: op, vars: variables})
	if data, ok := t.answers[op]; ok {
		return data, nil
	}
	return nil, errors.New("Linear GraphQL request failed")
}

func (t *linearScriptTransport) operations() []string {
	out := make([]string, 0, len(t.requests))
	for _, r := range t.requests {
		out = append(out, r.operation)
	}
	return out
}

const teamUUIDFixture = "5a3b982d-8f04-4971-956c-fbcb2c68642a"

func linearWriteMeta(t *testing.T) contract.Metadata {
	t.Helper()
	m := metaFor(t, contract.ProviderLinear, []string{"linear/*"}, nil)
	m.CredentialSource = contract.CredentialSourceOAuth
	m.AccountModel = contract.AccountModelPerUser
	m.Scopes = []string{LinearReadScope, LinearWriteScope, LinearIssuesCreateScope, LinearCommentsCreateScope}
	if err := m.Validate(); err != nil {
		t.Fatalf("fixture metadata is invalid: %v", err)
	}
	return m
}

func linearWriteInvocation(capability, resource string) contract.Invocation {
	inv := invocation(capability, contract.ActionCreate, resource)
	inv.IdempotencyKey = "idem-1"
	return inv
}

func linearWriteClient(answers map[string]json.RawMessage) (*LinearClient, *linearScriptTransport) {
	transport := &linearScriptTransport{answers: answers}
	return NewLinearRuntime(transport, &linearTestResolver{token: "linear-access-token"}), transport
}

func TestLinearIssueCreateResolvesTeamKeyAndSendsClientID(t *testing.T) {
	inv := linearWriteInvocation("issue.create", "linear/team/HAI")
	clientID := LinearWriteID(inv)
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"TeamID":      json.RawMessage(`{"teams":{"nodes":[{"id":"` + teamUUIDFixture + `","key":"HAI"}]}}`),
		"IssueCreate": json.RawMessage(`{"issueCreate":{"success":true,"issue":{"id":"` + clientID + `","identifier":"HAI-7","url":"https://linear.app/haikei/issue/HAI-7","team":{"key":"HAI"}}}}`),
	})
	priority := 2
	payload := LinearIssueCreatePayload{TeamID: "HAI", Title: "Crash on login", Description: "private-issue-body", Priority: &priority}

	result, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, payload)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	got := result.Data.(LinearIssueCreated)
	if got != (LinearIssueCreated{ID: clientID, Identifier: "HAI-7", TeamKey: "HAI", URL: "https://linear.app/haikei/issue/HAI-7"}) {
		t.Fatalf("result = %+v", got)
	}
	if ops := strings.Join(transport.operations(), ","); ops != "TeamID,IssueCreate" {
		t.Fatalf("operations = %s", ops)
	}
	input := transport.requests[1].vars["input"].(map[string]any)
	if input["id"] != clientID || input["teamId"] != teamUUIDFixture || input["title"] != "Crash on login" || input["priority"] != 2 || input["description"] != "private-issue-body" {
		t.Fatalf("issueCreate input = %#v", input)
	}
	encoded, _ := json.Marshal(result)
	if strings.Contains(string(encoded), "private-issue-body") || strings.Contains(string(encoded), "Crash on login") {
		t.Fatalf("result echoes submitted content: %s", encoded)
	}
}

func TestLinearIssueCreateWithTeamUUIDSkipsLookup(t *testing.T) {
	inv := linearWriteInvocation("issue.create", "linear/team/"+teamUUIDFixture)
	clientID := LinearWriteID(inv)
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"IssueCreate": json.RawMessage(`{"issueCreate":{"success":true,"issue":{"id":"` + clientID + `","identifier":"HAI-8","url":"u","team":{"key":"HAI"}}}}`),
	})
	if _, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, LinearIssueCreatePayload{TeamID: teamUUIDFixture, Title: "t"}); err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if ops := strings.Join(transport.operations(), ","); ops != "IssueCreate" {
		t.Fatalf("operations = %s", ops)
	}
}

func TestLinearIssueCreateReplayReturnsExistingIssueWithoutRetry(t *testing.T) {
	inv := linearWriteInvocation("issue.create", "linear/team/HAI")
	clientID := LinearWriteID(inv)
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"TeamID":   json.RawMessage(`{"teams":{"nodes":[{"id":"` + teamUUIDFixture + `","key":"HAI"}]}}`),
		"IssueRef": json.RawMessage(`{"issue":{"id":"` + clientID + `","identifier":"HAI-7","url":"u","team":{"key":"HAI"}}}`),
	})
	result, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, LinearIssueCreatePayload{TeamID: "HAI", Title: "t"})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if result.Data.(LinearIssueCreated).Identifier != "HAI-7" {
		t.Fatalf("result = %+v", result.Data)
	}
	if ops := strings.Join(transport.operations(), ","); ops != "TeamID,IssueCreate,IssueRef" {
		t.Fatalf("operations = %s (the mutation must run once)", ops)
	}
}

func TestLinearIssueCreateFailureIsGeneric(t *testing.T) {
	inv := linearWriteInvocation("issue.create", "linear/team/HAI")
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"TeamID": json.RawMessage(`{"teams":{"nodes":[{"id":"` + teamUUIDFixture + `","key":"HAI"}]}}`),
	})
	_, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, LinearIssueCreatePayload{TeamID: "HAI", Title: "t", Description: "private-issue-body"})
	if err == nil || err.Error() != "Linear issue create failed" {
		t.Fatalf("error = %v", err)
	}
	if ops := strings.Join(transport.operations(), ","); ops != "TeamID,IssueCreate,IssueRef" {
		t.Fatalf("operations = %s", ops)
	}
}

func TestLinearIssueCreateRejectsAmbiguousTeam(t *testing.T) {
	inv := linearWriteInvocation("issue.create", "linear/team/HAI")
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"TeamID": json.RawMessage(`{"teams":{"nodes":[]}}`),
	})
	if _, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, LinearIssueCreatePayload{TeamID: "HAI", Title: "t"}); err == nil {
		t.Fatal("unknown team accepted")
	}
	if ops := strings.Join(transport.operations(), ","); ops != "TeamID" {
		t.Fatalf("operations = %s", ops)
	}
}

func TestLinearCommentCreateSendsClientID(t *testing.T) {
	inv := linearWriteInvocation("comment.create", "linear/issue/HAI-7")
	clientID := LinearWriteID(inv)
	client, transport := linearWriteClient(map[string]json.RawMessage{
		"CommentCreate": json.RawMessage(`{"commentCreate":{"success":true,"comment":{"id":"` + clientID + `","url":"https://linear.app/c","issue":{"id":"issue-uuid"}}}}`),
	})
	result, err := client.Invoke(context.Background(), linearWriteMeta(t), inv, LinearCommentCreatePayload{IssueID: "HAI-7", Body: "private-comment-body"})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if got := result.Data.(LinearCommentCreated); got != (LinearCommentCreated{ID: clientID, IssueID: "issue-uuid", URL: "https://linear.app/c"}) {
		t.Fatalf("result = %+v", got)
	}
	input := transport.requests[0].vars["input"].(map[string]any)
	if input["id"] != clientID || input["issueId"] != "HAI-7" || input["body"] != "private-comment-body" {
		t.Fatalf("commentCreate input = %#v", input)
	}
}

func TestLinearWritesFailClosedBeforeAnyNetworkCall(t *testing.T) {
	priority := 9
	for name, tc := range map[string]struct {
		inv     contract.Invocation
		payload Payload
	}{
		"no idempotency key": {func() contract.Invocation {
			inv := linearWriteInvocation("issue.create", "linear/team/HAI")
			inv.IdempotencyKey = ""
			return inv
		}(), LinearIssueCreatePayload{TeamID: "HAI", Title: "t"}},
		"team_id differs from resource":  {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "ENG", Title: "t"}},
		"missing team_id":                {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{Title: "t"}},
		"issue resource for create":      {linearWriteInvocation("issue.create", "linear/issue/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "t"}},
		"empty title":                    {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "  "}},
		"multi-line title":               {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "a\nb"}},
		"priority out of range":          {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "t", Priority: &priority}},
		"description too large":          {linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "t", Description: strings.Repeat("x", MaxLinearDescriptionBytes+1)}},
		"issue_id differs":               {linearWriteInvocation("comment.create", "linear/issue/HAI-7"), LinearCommentCreatePayload{IssueID: "HAI-8", Body: "b"}},
		"team resource for comment":      {linearWriteInvocation("comment.create", "linear/team/HAI"), LinearCommentCreatePayload{IssueID: "HAI", Body: "b"}},
		"empty body":                     {linearWriteInvocation("comment.create", "linear/issue/HAI-7"), LinearCommentCreatePayload{IssueID: "HAI-7", Body: ""}},
		"payload for another capability": {linearWriteInvocation("comment.create", "linear/issue/HAI-7"), LinearIssueCreatePayload{TeamID: "HAI", Title: "t"}},
	} {
		t.Run(name, func(t *testing.T) {
			client, transport := linearWriteClient(nil)
			if _, err := client.Invoke(context.Background(), linearWriteMeta(t), tc.inv, tc.payload); err == nil {
				t.Fatal("write accepted")
			}
			if len(transport.requests) != 0 {
				t.Fatalf("network calls = %v", transport.operations())
			}
		})
	}
}

func TestLinearWriteRequiresGrantedScope(t *testing.T) {
	meta := linearWriteMeta(t)
	meta.Scopes = []string{LinearReadScope}
	client, transport := linearWriteClient(nil)
	_, err := client.Invoke(context.Background(), meta, linearWriteInvocation("issue.create", "linear/team/HAI"), LinearIssueCreatePayload{TeamID: "HAI", Title: "t"})
	if err == nil || !strings.Contains(err.Error(), "OAuth preflight") {
		t.Fatalf("error = %v", err)
	}
	if len(transport.requests) != 0 {
		t.Fatalf("network calls = %v", transport.operations())
	}
}

func TestLinearReadOnlyConnectorNeedsOnlyReadScope(t *testing.T) {
	meta := linearWriteMeta(t)
	meta.Capabilities = []contract.Capability{{Name: "team.read", Action: contract.ActionRead}, {Name: "issue.read", Action: contract.ActionRead}}
	meta.Scopes = []string{LinearReadScope}
	transport := &linearTestTransport{data: json.RawMessage(`{"team":{"key":"KEI","name":"Kei Engineering"}}`)}
	client := NewLinearRuntime(transport, &linearTestResolver{token: "linear-access-token"})
	if _, err := client.Invoke(context.Background(), meta, invocation("team.read", contract.ActionRead, "linear/team/KEI"), TeamReadPayload{}); err != nil {
		t.Fatalf("OAuth read rejected: %v", err)
	}
}

func TestLinearWriteIDIsStableAndScoped(t *testing.T) {
	a := linearWriteInvocation("issue.create", "linear/team/HAI")
	if LinearWriteID(a) != LinearWriteID(a) {
		t.Fatal("write id is not deterministic")
	}
	if !linearUUID.MatchString(LinearWriteID(a)) || LinearWriteID(a)[14] != '4' {
		t.Fatalf("write id %q is not a v4-shaped UUID", LinearWriteID(a))
	}
	for name, mutate := range map[string]func(*contract.Invocation){
		"subject":    func(in *contract.Invocation) { in.Subject = "u-2" },
		"key":        func(in *contract.Invocation) { in.IdempotencyKey = "idem-2" },
		"connector":  func(in *contract.Invocation) { in.ConnectorID = "c-2" },
		"capability": func(in *contract.Invocation) { in.Capability = "comment.create" },
	} {
		b := a
		mutate(&b)
		if LinearWriteID(a) == LinearWriteID(b) {
			t.Errorf("write id ignores %s", name)
		}
	}
}

func TestMemoryLinearWritesReplay(t *testing.T) {
	store := MemoryLinear{Teams: map[string]LinearTeam{"HAI": {Key: "HAI", Name: "Haikei"}}, Issues: map[string]LinearIssue{}, Comments: map[string]LinearCommentCreated{}}
	client := NewLinear(store)
	meta := metaFor(t, contract.ProviderLinear, []string{"linear/*"}, nil)
	inv := linearWriteInvocation("issue.create", "linear/team/HAI")
	first, err := client.Invoke(context.Background(), meta, inv, LinearIssueCreatePayload{TeamID: "HAI", Title: "t"})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	second, err := client.Invoke(context.Background(), meta, inv, LinearIssueCreatePayload{TeamID: "HAI", Title: "t"})
	if err != nil || first.Data != second.Data || len(store.Issues) != 1 {
		t.Fatalf("replay = %+v, %v (issues %d)", second.Data, err, len(store.Issues))
	}
	issue := first.Data.(LinearIssueCreated)
	comment := linearWriteInvocation("comment.create", "linear/issue/"+issue.Identifier)
	if _, err := client.Invoke(context.Background(), meta, comment, LinearCommentCreatePayload{IssueID: issue.Identifier, Body: "b"}); err != nil {
		t.Fatalf("comment: %v", err)
	}
}

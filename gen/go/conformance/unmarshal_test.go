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

// Package conformance is the HAI-990 gate: every pilot operation registered
// in kei.contracts.v1 must have synthetic request and response fixtures, and
// those fixtures must unmarshal into the generated Go types.
//
// Lenient unmarshal (DiscardUnknown) proves the happy path: a provider may
// add keys without breaking the contract. Strict unmarshal proves the drift
// gate: every response fixture carries one unknown extra field, so a provider
// payload that no longer fits the declared shape is detected. A separate
// request-side drift fixture proves a renamed request key (repo instead of
// repository) is caught instead of silently unmarshaling an empty field.
package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	airtablev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/airtable/v1"
	atlassianv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/atlassian/v1"
	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	hubspotv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/hubspot/v1"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	notionv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/notion/v1"
	sentryv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/sentry/v1"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

const fixturesDir = "../../../fixtures/operations"

// fixtureCase pairs one operation id with the generated request/response
// messages its fixtures must unmarshal into, plus one spot assertion proving
// a field survived unmarshal with its synthetic value intact.
type fixtureCase struct {
	op        string
	connector string
	request   func() proto.Message
	response  func() proto.Message
	spotCheck func(t *testing.T, resp proto.Message)
}

var fixtureCases = []fixtureCase{
	{
		op:        "github.issues.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetIssueRequest{} },
		response:  func() proto.Message { return &githubv1.GetIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issue := resp.(*githubv1.GetIssueResponse)
			if issue.GetNumber() != 1347 {
				t.Errorf("number = %d, want 1347", issue.GetNumber())
			}
			if issue.GetId() != 9007199254740993 {
				t.Errorf("id = %d, want 9007199254740993 (int64 must not go through float64)", issue.GetId())
			}
		},
	},
	{
		op:        "github.issues.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListIssuesRequest{} },
		response:  func() proto.Message { return &githubv1.ListIssuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := len(resp.(*githubv1.ListIssuesResponse).GetIssues()); got != 2 {
				t.Errorf("len(issues) = %d, want 2", got)
			}
		},
	},
	{
		op:        "github.contents.put",
		connector: "github",
		request:   func() proto.Message { return &githubv1.PutContentsRequest{} },
		response:  func() proto.Message { return &githubv1.PutContentsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := resp.(*githubv1.PutContentsResponse).GetContent().GetName(); got != "notes.txt" {
				t.Errorf("content.name = %q, want notes.txt", got)
			}
		},
	},
	{
		op:        "github.pulls.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetPullRequestRequest{} },
		response:  func() proto.Message { return &githubv1.GetPullRequestResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := resp.(*githubv1.GetPullRequestResponse).GetNumber(); got != 42 {
				t.Errorf("number = %d, want 42", got)
			}
		},
	},
	{
		op:        "linear.issues.get",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.GetIssueRequest{} },
		response:  func() proto.Message { return &linearv1.GetIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := resp.(*linearv1.GetIssueResponse).GetIssue().GetIdentifier(); got != "ENG-101" {
				t.Errorf("issue.identifier = %q, want ENG-101", got)
			}
			if got := len(resp.(*linearv1.GetIssueResponse).GetIssue().GetLabels()); got != 1 {
				t.Errorf("len(issue.labels) = %d, want 1", got)
			}
		},
	},
	{
		op:        "linear.issues.create",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.CreateIssueRequest{} },
		response:  func() proto.Message { return &linearv1.CreateIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := resp.(*linearv1.CreateIssueResponse).GetSuccess(); !got {
				t.Errorf("success = %v, want true", got)
			}
		},
	},
	{
		op:        "linear.issues.search",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.SearchIssuesRequest{} },
		response:  func() proto.Message { return &linearv1.SearchIssuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := len(resp.(*linearv1.SearchIssuesResponse).GetNodes()); got != 2 {
				t.Errorf("len(nodes) = %d, want 2", got)
			}
		},
	},
	{
		op:        "linear.teams.get",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.GetTeamRequest{} },
		response:  func() proto.Message { return &linearv1.GetTeamResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := resp.(*linearv1.GetTeamResponse).GetTeam().GetKey(); got != "ENG" {
				t.Errorf("team.key = %q, want ENG", got)
			}
		},
	},
	{
		op:        "linear.comments.create",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.CreateCommentRequest{} },
		response:  func() proto.Message { return &linearv1.CreateCommentResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			created := resp.(*linearv1.CreateCommentResponse)
			if !created.GetSuccess() {
				t.Errorf("success = %v, want true", created.GetSuccess())
			}
			if got := created.GetComment().GetId(); got == "" {
				t.Errorf("comment.id is empty, want a synthetic id")
			}
		},
	},
	{
		op:        "atlassian.jira.issues.get",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.GetJiraIssueRequest{} },
		response:  func() proto.Message { return &atlassianv1.GetJiraIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.GetJiraIssueResponse)
			if got.GetKey() != "SYN-101" {
				t.Errorf("key = %q, want SYN-101", got.GetKey())
			}
			if got.GetFields().GetSummary() != "Synthetic Jira issue" {
				t.Errorf("fields.summary = %q, want Synthetic Jira issue", got.GetFields().GetSummary())
			}
			if got.GetFields().Description != nil {
				t.Errorf("fields.description = %v, want nil (provider null)", got.GetFields().Description)
			}
		},
	},
	{
		op:        "atlassian.jira.issues.search",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.SearchJiraIssuesRequest{} },
		response:  func() proto.Message { return &atlassianv1.SearchJiraIssuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.SearchJiraIssuesResponse)
			if len(got.GetIssues()) != 2 {
				t.Errorf("len(issues) = %d, want 2", len(got.GetIssues()))
			}
			if got.GetNextPageToken() != "synthetic-next-token" {
				t.Errorf("next_page_token = %q, want synthetic-next-token", got.GetNextPageToken())
			}
			if got.GetIssues()[1].GetFields().Description != nil {
				t.Errorf("issues[1].fields.description = %v, want nil (provider null)", got.GetIssues()[1].GetFields().Description)
			}
		},
	},
	{
		op:        "atlassian.confluence.pages.get",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.GetConfluencePageRequest{} },
		response:  func() proto.Message { return &atlassianv1.GetConfluencePageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.GetConfluencePageResponse)
			if got.GetId() != 78643201 {
				t.Errorf("id = %d, want 78643201", got.GetId())
			}
			if got.GetTitle() != "Synthetic Confluence page" {
				t.Errorf("title = %q, want Synthetic Confluence page", got.GetTitle())
			}
			if got.ParentId != nil {
				t.Errorf("parent_id = %v, want nil (provider null)", got.ParentId)
			}
		},
	},
	{
		op:        "atlassian.confluence.pages.list",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.ListConfluencePagesRequest{} },
		response:  func() proto.Message { return &atlassianv1.ListConfluencePagesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.ListConfluencePagesResponse)
			if len(got.GetResults()) != 2 {
				t.Errorf("len(results) = %d, want 2", len(got.GetResults()))
			}
			if got.GetResults()[0].GetTitle() != "Synthetic first page" {
				t.Errorf("results[0].title = %q, want Synthetic first page", got.GetResults()[0].GetTitle())
			}
		},
	},
	{
		op:        "airtable.records.get",
		connector: "airtable",
		request:   func() proto.Message { return &airtablev1.GetRecordRequest{} },
		response:  func() proto.Message { return &airtablev1.GetRecordResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			rec := resp.(*airtablev1.GetRecordResponse)
			if rec.GetId() != "recSynthetic00003" {
				t.Errorf("id = %q, want recSynthetic00003", rec.GetId())
			}
			if got := rec.GetFields().GetFields()["Name"].GetStringValue(); got != "Synthetic third task" {
				t.Errorf("fields.Name = %q, want Synthetic third task", got)
			}
			if got := rec.GetFields().GetFields()["Priority"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("fields.Priority = %v, want null", got)
			}
		},
	},
	{
		op:        "airtable.records.list",
		connector: "airtable",
		request:   func() proto.Message { return &airtablev1.ListRecordsRequest{} },
		response:  func() proto.Message { return &airtablev1.ListRecordsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			list := resp.(*airtablev1.ListRecordsResponse)
			if got := len(list.GetRecords()); got != 2 {
				t.Errorf("len(records) = %d, want 2", got)
			}
			if list.GetOffset() != "synthetic-next-offset" {
				t.Errorf("offset = %q, want synthetic-next-offset", list.GetOffset())
			}
			if got := list.GetRecords()[1].GetFields().GetFields()["Priority"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("records[1].fields.Priority = %v, want null", got)
			}
		},
	},
	{
		op:        "airtable.records.create",
		connector: "airtable",
		request:   func() proto.Message { return &airtablev1.CreateRecordRequest{} },
		response:  func() proto.Message { return &airtablev1.CreateRecordResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			rec := resp.(*airtablev1.CreateRecordResponse)
			if rec.GetId() != "recSynthetic00004" {
				t.Errorf("id = %q, want recSynthetic00004", rec.GetId())
			}
			if got := rec.GetFields().GetFields()["Status"].GetStringValue(); got != "Todo" {
				t.Errorf("fields.Status = %q, want Todo", got)
			}
		},
	},
	{
		op:        "airtable.records.update",
		connector: "airtable",
		request:   func() proto.Message { return &airtablev1.UpdateRecordRequest{} },
		response:  func() proto.Message { return &airtablev1.UpdateRecordResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			rec := resp.(*airtablev1.UpdateRecordResponse)
			if got := rec.GetFields().GetFields()["Status"].GetStringValue(); got != "Done" {
				t.Errorf("fields.Status = %q, want Done", got)
			}
			if got := rec.GetFields().GetFields()["Priority"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("fields.Priority = %v, want null", got)
			}
		},
	},
	{
		op:        "hubspot.contacts.get",
		connector: "hubspot",
		request:   func() proto.Message { return &hubspotv1.GetContactRequest{} },
		response:  func() proto.Message { return &hubspotv1.GetContactResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			contact := resp.(*hubspotv1.GetContactResponse)
			if contact.GetId() != "701" {
				t.Errorf("id = %q, want 701", contact.GetId())
			}
			if got := contact.GetProperties().GetFields()["email"].GetStringValue(); got != "synthetic.contact@example.invalid" {
				t.Errorf("properties.email = %q, want the synthetic email", got)
			}
			if got := contact.GetProperties().GetFields()["lastname"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("properties.lastname = %v, want null", got)
			}
		},
	},
	{
		op:        "hubspot.contacts.create",
		connector: "hubspot",
		request:   func() proto.Message { return &hubspotv1.CreateContactRequest{} },
		response:  func() proto.Message { return &hubspotv1.CreateContactResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			contact := resp.(*hubspotv1.CreateContactResponse)
			if contact.GetId() != "701550001" {
				t.Errorf("id = %q, want 701550001", contact.GetId())
			}
			if got := contact.GetProperties().GetFields()["email"].GetStringValue(); got != "synthetic.new@example.invalid" {
				t.Errorf("properties.email = %q, want the synthetic email", got)
			}
		},
	},
	{
		op:        "hubspot.contacts.update",
		connector: "hubspot",
		request:   func() proto.Message { return &hubspotv1.UpdateContactRequest{} },
		response:  func() proto.Message { return &hubspotv1.UpdateContactResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			contact := resp.(*hubspotv1.UpdateContactResponse)
			if got := contact.GetProperties().GetFields()["lifecyclestage"].GetStringValue(); got != "customer" {
				t.Errorf("properties.lifecyclestage = %q, want customer", got)
			}
		},
	},
	{
		op:        "hubspot.contacts.list",
		connector: "hubspot",
		request:   func() proto.Message { return &hubspotv1.ListContactsRequest{} },
		response:  func() proto.Message { return &hubspotv1.ListContactsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			list := resp.(*hubspotv1.ListContactsResponse)
			if got := len(list.GetResults()); got != 2 {
				t.Errorf("len(results) = %d, want 2", got)
			}
			if list.GetNextAfter() != "synthetic-next-cursor" {
				t.Errorf("next_after = %q, want synthetic-next-cursor", list.GetNextAfter())
			}
			if got := list.GetResults()[1].GetProperties().GetFields()["lastname"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("results[1].properties.lastname = %v, want null", got)
			}
		},
	},
	{
		op:        "hubspot.deals.list",
		connector: "hubspot",
		request:   func() proto.Message { return &hubspotv1.ListDealsRequest{} },
		response:  func() proto.Message { return &hubspotv1.ListDealsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			list := resp.(*hubspotv1.ListDealsResponse)
			if got := len(list.GetResults()); got != 2 {
				t.Errorf("len(results) = %d, want 2", got)
			}
			if got := list.GetResults()[0].GetProperties().GetFields()["dealname"].GetStringValue(); got != "Synthetic first deal" {
				t.Errorf("results[0].properties.dealname = %q, want Synthetic first deal", got)
			}
			if got := list.GetResults()[0].GetProperties().GetFields()["description"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("results[0].properties.description = %v, want null", got)
			}
		},
	},
	{
		op:        "notion.pages.get",
		connector: "notion",
		request:   func() proto.Message { return &notionv1.GetPageRequest{} },
		response:  func() proto.Message { return &notionv1.GetPageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			page := resp.(*notionv1.GetPageResponse)
			if page.GetId() != "5f8b7e4a-0000-4000-8000-000000000001" {
				t.Errorf("id = %q, want the synthetic page id", page.GetId())
			}
			if page.PublicUrl != nil {
				t.Error("public_url is set, want nil (null in the fixture)")
			}
		},
	},
	{
		op:        "notion.databases.query",
		connector: "notion",
		request:   func() proto.Message { return &notionv1.QueryDatabaseRequest{} },
		response:  func() proto.Message { return &notionv1.QueryDatabaseResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			q := resp.(*notionv1.QueryDatabaseResponse)
			if got := len(q.GetResults()); got != 2 {
				t.Errorf("len(results) = %d, want 2", got)
			}
			if q.GetNextCursor() != "synthetic-next-cursor" {
				t.Errorf("next_cursor = %q, want synthetic-next-cursor", q.GetNextCursor())
			}
			if !q.GetHasMore() {
				t.Error("has_more = false, want true")
			}
			if q.GetResults()[1].PublicUrl != nil {
				t.Error("results[1].public_url is set, want nil (null in the fixture)")
			}
		},
	},
	{
		op:        "sentry.issues.get",
		connector: "sentry",
		request:   func() proto.Message { return &sentryv1.GetIssueRequest{} },
		response:  func() proto.Message { return &sentryv1.GetIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issue := resp.(*sentryv1.GetIssueResponse)
			if issue.GetShortId() != "PROJ-103" {
				t.Errorf("short_id = %q, want PROJ-103", issue.GetShortId())
			}
			if issue.ShareId != nil {
				t.Error("share_id is set, want nil (null in the fixture)")
			}
		},
	},
	{
		op:        "sentry.issues.list",
		connector: "sentry",
		request:   func() proto.Message { return &sentryv1.ListIssuesRequest{} },
		response:  func() proto.Message { return &sentryv1.ListIssuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issues := resp.(*sentryv1.ListIssuesResponse).GetIssues()
			if got := len(issues); got != 2 {
				t.Errorf("len(issues) = %d, want 2", got)
			}
			if issues[1].ShareId == nil || issues[1].GetShareId() != "synthetic-share-id-2" {
				t.Errorf("issues[1].share_id = %v, want synthetic-share-id-2", issues[1].ShareId)
			}
		},
	},
	{
		op:        "http_api.records.list",
		connector: "http_api",
		request:   func() proto.Message { return &httpapiv1.ListRecordsRequest{} },
		response:  func() proto.Message { return &httpapiv1.ListRecordsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			if got := len(resp.(*httpapiv1.ListRecordsResponse).GetRecords()); got != 2 {
				t.Errorf("len(records) = %d, want 2", got)
			}
		},
	},
	{
		op:        "http_api.records.get",
		connector: "http_api",
		request:   func() proto.Message { return &httpapiv1.GetRecordRequest{} },
		response:  func() proto.Message { return &httpapiv1.GetRecordResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			record := resp.(*httpapiv1.GetRecordResponse).GetRecord()
			if got := record.GetFields()["id"].GetStringValue(); got != "rec-1" {
				t.Errorf("record.id = %q, want rec-1", got)
			}
		},
	},
	{
		op:        "database.query.run",
		connector: "database",
		request:   func() proto.Message { return &databasev1.RunQueryRequest{} },
		response:  func() proto.Message { return &databasev1.RunQueryResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			result := resp.(*databasev1.RunQueryResponse)
			if got := len(result.GetColumns()); got != 4 {
				t.Fatalf("len(columns) = %d, want 4", got)
			}
			if got := result.GetColumns()[0].GetName(); got != "id" {
				t.Errorf("columns[0].name = %q, want id", got)
			}
			if got := result.GetColumns()[0].GetType(); got != "int8" {
				t.Errorf("columns[0].type = %q, want int8", got)
			}
			if got := result.GetColumns()[3].GetType(); got != "numeric" {
				t.Errorf("columns[3].type = %q, want numeric", got)
			}
			rows := result.GetRows()
			if got := len(rows); got != 2 {
				t.Fatalf("len(rows) = %d, want 2", got)
			}
			big := rows[0].GetValues()
			if got := big[0].GetStringValue(); got != "9007199254740993" {
				t.Errorf("rows[0][0] = %q, want the >2^53 id as a string (JSON numbers are float64)", got)
			}
			if got := big[1].GetStringValue(); got != "Big-ID Account — ünïcodé ✓" {
				t.Errorf("rows[0][1] = %q, want the unicode account name", got)
			}
			if got := big[2].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("rows[0][2] = %v, want null", got)
			}
			if got := big[3].GetStringValue(); got != "12345678901234567890.123456789012345678" {
				t.Errorf("rows[0][3] = %q, want NUMERIC carried as a precision-preserving string", got)
			}
			small := rows[1].GetValues()
			if got := small[3].GetStringValue(); got != "0.000000000000000001" {
				t.Errorf("rows[1][3] = %q, want numeric fraction preserved", got)
			}
			if result.GetTruncated() {
				t.Error("truncated = true, want false")
			}
		},
	},
}

// readFixture reads one synthetic fixture. Paths are relative to this test's
// package directory: the fixtures live at the repository root.
func readFixture(t *testing.T, connector, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesDir, connector, name))
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	return data
}

// opFixtureName maps an operation id to its fixture file name:
// github.issues.get + "request" becomes issues.get.request.json.
func opFixtureName(op, kind string) string {
	return op[strings.Index(op, ".")+1:] + "." + kind + ".json"
}

func TestFixturesUnmarshal(t *testing.T) {
	registered := contracts.All()
	if len(registered) != len(fixtureCases) {
		t.Fatalf("contracts.All() registers %d operations but the fixture table covers %d; every registered operation needs a fixture", len(registered), len(fixtureCases))
	}
	covered := make(map[string]bool, len(fixtureCases))
	for _, tt := range fixtureCases {
		covered[tt.op] = true
	}
	for _, op := range registered {
		if !covered[op.GetId()] {
			t.Errorf("operation %q is registered but has no fixture table entry", op.GetId())
		}
	}

	for _, tt := range fixtureCases {
		t.Run(tt.op, func(t *testing.T) {
			reqData := readFixture(t, tt.connector, opFixtureName(tt.op, "request"))
			respData := readFixture(t, tt.connector, opFixtureName(tt.op, "response"))

			req := tt.request()
			if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(reqData, req); err != nil {
				t.Fatalf("request fixture did not unmarshal: %v", err)
			}
			resp := tt.response()
			if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(respData, resp); err != nil {
				t.Fatalf("response fixture did not unmarshal: %v", err)
			}
			tt.spotCheck(t, resp)
		})
	}
}

func TestDatabaseQueryRunRequestFixture(t *testing.T) {
	data := readFixture(t, "database", "query.run.request.json")
	var req databasev1.RunQueryRequest
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, &req); err != nil {
		t.Fatalf("request fixture did not unmarshal: %v", err)
	}
	if got := req.GetDialect(); got != databasev1.Dialect_DIALECT_POSTGRES {
		t.Errorf("dialect = %v, want DIALECT_POSTGRES", got)
	}
	if got := req.GetSql(); !strings.Contains(got, "accounts") {
		t.Errorf("sql = %q, want the synthetic accounts query", got)
	}
	params := req.GetParams()
	if len(params) != 1 {
		t.Fatalf("len(params) = %d, want 1", len(params))
	}
	if got := params[0].GetStringValue(); got != "9007199254740993" {
		t.Errorf("params[0] = %q, want the >2^53 bind value as a string", got)
	}
	if got := req.GetMaxRows(); got != 100 {
		t.Errorf("max_rows = %d, want 100", got)
	}
}

func TestFixturesStrictDetectsUnknown(t *testing.T) {
	for _, tt := range fixtureCases {
		t.Run(tt.op, func(t *testing.T) {
			respData := readFixture(t, tt.connector, opFixtureName(tt.op, "response"))
			err := (protojson.UnmarshalOptions{}).Unmarshal(respData, tt.response())
			if err == nil {
				t.Errorf("strict unmarshal of the %s response fixture succeeded; every response fixture must carry one unknown extra field so the drift gate detects provider additions", tt.op)
			}
		})
	}
}

func TestGitHubRepoDriftCaught(t *testing.T) {
	data := readFixture(t, "github", "issues.get.request_repo_drift.json")
	err := (protojson.UnmarshalOptions{}).Unmarshal(data, &githubv1.GetIssueRequest{})
	if err == nil {
		t.Fatal("strict unmarshal of issues.get.request_repo_drift.json succeeded; the repo key must stay unknown so a rename of repository is caught instead of silently leaving the field empty")
	}
}

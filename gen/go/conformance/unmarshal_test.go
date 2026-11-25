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

// Package conformance is the gate: every pilot operation registered
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
	gmailv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/gmail/v1"
	googlecalendarv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/google_calendar/v1"
	googlesheetsv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/google_sheets/v1"
	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	hubspotv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/hubspot/v1"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	notionv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/notion/v1"
	omniv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/omni/v1"
	powerbiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/power_bi/v1"
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
			pull := resp.(*githubv1.GetPullRequestResponse)
			if got := pull.GetNumber(); got != 42 {
				t.Errorf("number = %d, want 42", got)
			}
			if got := pull.GetDiffUrl(); got != "https://example.invalid/example-org/example-repo/pull/42.diff" {
				t.Errorf("diff_url = %q, want the .diff URL", got)
			}
			if got := pull.GetPatchUrl(); got != "https://example.invalid/example-org/example-repo/pull/42.patch" {
				t.Errorf("patch_url = %q, want the .patch URL", got)
			}
		},
	},
	{
		op:        "github.pulls.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListPullsRequest{} },
		response:  func() proto.Message { return &githubv1.ListPullsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			pulls := resp.(*githubv1.ListPullsResponse)
			if got := pulls.GetPulls()[0].GetHead().GetRef(); got != "feature/x" {
				t.Errorf("pulls[0].head.ref = %q, want feature/x", got)
			}
			if got := pulls.GetPulls()[0].GetId(); got != 9007199254740999 {
				t.Errorf("pulls[0].id = %d, want 9007199254740999 (int64 must not go through float64)", got)
			}
		},
	},
	{
		op:        "github.actions.runs.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListWorkflowRunsRequest{} },
		response:  func() proto.Message { return &githubv1.ListWorkflowRunsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			runs := resp.(*githubv1.ListWorkflowRunsResponse)
			if got := runs.GetWorkflowRuns()[0].GetHeadBranch(); got != "feature/x" {
				t.Errorf("workflow_runs[0].head_branch = %q, want feature/x", got)
			}
			if got := runs.GetWorkflowRuns()[0].GetId(); got != 9007199254741001 {
				t.Errorf("workflow_runs[0].id = %d, want 9007199254741001 (int64 must not go through float64)", got)
			}
			if got := runs.GetWorkflowRuns()[0].GetActor().GetLogin(); got != "example-user" {
				t.Errorf("workflow_runs[0].actor.login = %q, want example-user", got)
			}
			if got := runs.GetWorkflowRuns()[0].GetHeadSha(); got != "synthetic-head-sha-1" {
				t.Errorf("workflow_runs[0].head_sha = %q, want the synthetic sha", got)
			}
		},
	},
	{
		op:        "github.actions.workflows.runs.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListWorkflowRunsByWorkflowRequest{} },
		response:  func() proto.Message { return &githubv1.ListWorkflowRunsByWorkflowResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			runs := resp.(*githubv1.ListWorkflowRunsByWorkflowResponse)
			if got := len(runs.GetWorkflowRuns()); got != 1 {
				t.Errorf("len(workflow_runs) = %d, want 1", got)
			}
			if got := runs.GetWorkflowRuns()[0].GetActor().GetLogin(); got != "example-user" {
				t.Errorf("workflow_runs[0].actor.login = %q, want example-user", got)
			}
			if got := runs.GetWorkflowRuns()[0].GetHeadSha(); got != "synthetic-head-sha-1" {
				t.Errorf("workflow_runs[0].head_sha = %q, want the synthetic sha", got)
			}
		},
	},
	{
		op:        "github.repos.contents.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetContentsRequest{} },
		response:  func() proto.Message { return &githubv1.GetContentsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			file := resp.(*githubv1.GetContentsResponse)
			if got := file.GetEncoding(); got != "base64" {
				t.Errorf("encoding = %q, want base64", got)
			}
			if got := file.GetContent(); got != "c3ludGhldGljLWNvbnRlbnQ=" {
				t.Errorf("content = %q, want the synthetic base64", got)
			}
		},
	},
	{
		op:        "github.repos.contents.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListContentsRequest{} },
		response:  func() proto.Message { return &githubv1.ListContentsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			entries := resp.(*githubv1.ListContentsResponse)
			if got := len(entries.GetEntries()); got != 2 {
				t.Errorf("len(entries) = %d, want 2", got)
			}
			if got := entries.GetEntries()[0].GetType(); got != "file" {
				t.Errorf("entries[0].type = %q, want file", got)
			}
			if got := entries.GetEntries()[1].GetType(); got != "dir" {
				t.Errorf("entries[1].type = %q, want dir", got)
			}
		},
	},
	{
		op:        "github.repos.commits.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListCommitsRequest{} },
		response:  func() proto.Message { return &githubv1.ListCommitsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			commits := resp.(*githubv1.ListCommitsResponse)
			if got := commits.GetCommits()[0].GetSha(); got != "synthetic-commit-sha-3" {
				t.Errorf("commits[0].sha = %q, want synthetic-commit-sha-3", got)
			}
			if got := commits.GetCommits()[0].GetAuthor().GetLogin(); got != "example-user" {
				t.Errorf("commits[0].author.login = %q, want example-user", got)
			}
		},
	},
	{
		op:        "github.repos.commits.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetCommitRequest{} },
		response:  func() proto.Message { return &githubv1.GetCommitResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			commit := resp.(*githubv1.GetCommitResponse)
			if got := commit.GetCommit().GetMessage(); got != "Synthetic commit message" {
				t.Errorf("commit.message = %q, want Synthetic commit message", got)
			}
			if got := commit.GetAuthor().GetLogin(); got != "example-user" {
				t.Errorf("author.login = %q, want example-user", got)
			}
		},
	},
	{
		op:        "github.actions.runs.attempts.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetRunAttemptRequest{} },
		response:  func() proto.Message { return &githubv1.GetRunAttemptResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			run := resp.(*githubv1.GetRunAttemptResponse)
			if got := run.GetRunAttempt(); got != 2 {
				t.Errorf("run_attempt = %d, want 2", got)
			}
			if got := run.GetHeadSha(); got != "synthetic-head-sha-1" {
				t.Errorf("head_sha = %q, want the synthetic sha", got)
			}
		},
	},
	{
		op:        "github.actions.runs.jobs.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListRunJobsRequest{} },
		response:  func() proto.Message { return &githubv1.ListRunJobsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			jobs := resp.(*githubv1.ListRunJobsResponse)
			if got := jobs.GetJobs()[0].GetId(); got != 9007199254741005 {
				t.Errorf("jobs[0].id = %d, want 9007199254741005 (int64 must not go through float64)", got)
			}
			if got := jobs.GetJobs()[0].GetRunnerName(); got != "synthetic-runner-1" {
				t.Errorf("jobs[0].runner_name = %q, want synthetic-runner-1", got)
			}
		},
	},
	{
		op:        "github.issues.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateIssueRequest{} },
		response:  func() proto.Message { return &githubv1.CreateIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issue := resp.(*githubv1.CreateIssueResponse)
			if got := issue.GetNumber(); got != 42 {
				t.Errorf("number = %d, want 42", got)
			}
			if got := issue.GetAuthor().GetLogin(); got != "example-user" {
				t.Errorf("author.login = %q, want example-user (json_name user)", got)
			}
		},
	},
	{
		op:        "github.issues.comments.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateIssueCommentRequest{} },
		response:  func() proto.Message { return &githubv1.CreateIssueCommentResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			comment := resp.(*githubv1.CreateIssueCommentResponse)
			if got := comment.GetBody(); got != "Synthetic comment" {
				t.Errorf("body = %q, want Synthetic comment", got)
			}
			if got := comment.GetAuthor().GetLogin(); got != "example-user" {
				t.Errorf("author.login = %q, want example-user (json_name user)", got)
			}
		},
	},
	{
		op:        "github.repos.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetRepositoryRequest{} },
		response:  func() proto.Message { return &githubv1.GetRepositoryResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			repo := resp.(*githubv1.GetRepositoryResponse)
			if got := repo.GetId(); got != 9007199254741011 {
				t.Errorf("id = %d, want 9007199254741011 (int64 must not go through float64)", got)
			}
			if got := repo.GetFullName(); got != "example-org/example-repo" {
				t.Errorf("full_name = %q, want example-org/example-repo", got)
			}
			if got := repo.GetOwner().GetLogin(); got != "example-user" {
				t.Errorf("owner.login = %q, want example-user (json_name user)", got)
			}
		},
	},
	{
		op:        "github.checks.runs.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListCheckRunsRequest{} },
		response:  func() proto.Message { return &githubv1.ListCheckRunsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			runs := resp.(*githubv1.ListCheckRunsResponse)
			if got := runs.GetCheckRuns()[0].GetId(); got != 9007199254741012 {
				t.Errorf("check_runs[0].id = %d, want 9007199254741012 (int64 must not go through float64)", got)
			}
			if got := runs.GetCheckRuns()[0].GetConclusion(); got != "success" {
				t.Errorf("check_runs[0].conclusion = %q, want success", got)
			}
			if got := runs.GetTotalCount(); got != 2 {
				t.Errorf("total_count = %d, want 2", got)
			}
		},
	},
	{
		op:        "github.issues.update",
		connector: "github",
		request:   func() proto.Message { return &githubv1.UpdateIssueRequest{} },
		response:  func() proto.Message { return &githubv1.UpdateIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issue := resp.(*githubv1.UpdateIssueResponse)
			if got := issue.GetNumber(); got != 42 {
				t.Errorf("number = %d, want 42", got)
			}
			if got := issue.GetState(); got != "closed" {
				t.Errorf("state = %q, want closed", got)
			}
			if got := issue.GetAuthor().GetLogin(); got != "example-user" {
				t.Errorf("author.login = %q, want example-user (json_name user)", got)
			}
		},
	},
	{
		op:        "github.pulls.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreatePullRequestRequest{} },
		response:  func() proto.Message { return &githubv1.CreatePullRequestResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			pull := resp.(*githubv1.CreatePullRequestResponse)
			if got := pull.GetNumber(); got != 43 {
				t.Errorf("number = %d, want 43", got)
			}
			if !pull.GetDraft() {
				t.Error("draft = false, want true")
			}
		},
	},
	{
		op:        "github.pulls.update",
		connector: "github",
		request:   func() proto.Message { return &githubv1.UpdatePullRequestRequest{} },
		response:  func() proto.Message { return &githubv1.UpdatePullRequestResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			pull := resp.(*githubv1.UpdatePullRequestResponse)
			if got := pull.GetTitle(); got != "Synthetic updated title" {
				t.Errorf("title = %q, want Synthetic updated title", got)
			}
			if got := pull.GetState(); got != "open" {
				t.Errorf("state = %q, want open", got)
			}
		},
	},
	{
		op:        "github.pulls.reviews.list",
		connector: "github",
		request:   func() proto.Message { return &githubv1.ListPullReviewsRequest{} },
		response:  func() proto.Message { return &githubv1.ListPullReviewsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			reviews := resp.(*githubv1.ListPullReviewsResponse)
			if got := reviews.GetReviews()[0].GetState(); got != "APPROVED" {
				t.Errorf("reviews[0].state = %q, want APPROVED", got)
			}
			if got := reviews.GetReviews()[0].GetAuthor().GetLogin(); got != "example-reviewer" {
				t.Errorf("reviews[0].author.login = %q, want example-reviewer (json_name user)", got)
			}
		},
	},
	{
		op:        "github.pulls.reviews.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreatePullReviewRequest{} },
		response:  func() proto.Message { return &githubv1.CreatePullReviewResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			review := resp.(*githubv1.CreatePullReviewResponse)
			if got := review.GetState(); got != "APPROVED" {
				t.Errorf("state = %q, want APPROVED", got)
			}
			if got := review.GetId(); got != 9007199254741006 {
				t.Errorf("id = %d, want 9007199254741006 (int64 must not go through float64)", got)
			}
		},
	},
	{
		op:        "github.git.refs.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetRefRequest{} },
		response:  func() proto.Message { return &githubv1.GetRefResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			ref := resp.(*githubv1.GetRefResponse)
			if got := ref.GetRef(); got != "refs/heads/feature/x" {
				t.Errorf("ref = %q, want refs/heads/feature/x", got)
			}
			if got := ref.GetObject().GetSha(); got != "synthetic-ref-sha-1" {
				t.Errorf("object.sha = %q, want the synthetic ref sha", got)
			}
		},
	},
	{
		op:        "github.git.refs.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateRefRequest{} },
		response:  func() proto.Message { return &githubv1.CreateRefResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			ref := resp.(*githubv1.CreateRefResponse)
			if got := ref.GetRef(); got != "refs/heads/feature/y" {
				t.Errorf("ref = %q, want refs/heads/feature/y", got)
			}
			if got := ref.GetObject().GetType(); got != "commit" {
				t.Errorf("object.type = %q, want commit", got)
			}
		},
	},
	{
		op:        "github.git.refs.update",
		connector: "github",
		request:   func() proto.Message { return &githubv1.UpdateRefRequest{} },
		response:  func() proto.Message { return &githubv1.UpdateRefResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			ref := resp.(*githubv1.UpdateRefResponse)
			if got := ref.GetObject().GetSha(); got != "synthetic-ref-sha-2" {
				t.Errorf("object.sha = %q, want the moved synthetic sha", got)
			}
		},
	},
	{
		op:        "github.git.blobs.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateBlobRequest{} },
		response:  func() proto.Message { return &githubv1.CreateBlobResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			blob := resp.(*githubv1.CreateBlobResponse)
			if got := blob.GetSha(); got != "synthetic-blob-sha-1" {
				t.Errorf("sha = %q, want the synthetic blob sha", got)
			}
		},
	},
	{
		op:        "github.git.trees.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateTreeRequest{} },
		response:  func() proto.Message { return &githubv1.CreateTreeResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			tree := resp.(*githubv1.CreateTreeResponse)
			if got := tree.GetSha(); got != "synthetic-tree-sha-1" {
				t.Errorf("sha = %q, want the synthetic tree sha", got)
			}
			if got := tree.GetTree()[0].GetPath(); got != "notes.txt" {
				t.Errorf("tree[0].path = %q, want notes.txt", got)
			}
		},
	},
	{
		op:        "github.git.commits.create",
		connector: "github",
		request:   func() proto.Message { return &githubv1.CreateCommitRequest{} },
		response:  func() proto.Message { return &githubv1.CreateCommitResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			commit := resp.(*githubv1.CreateCommitResponse)
			if got := commit.GetSha(); got != "synthetic-commit-sha-2" {
				t.Errorf("sha = %q, want the synthetic commit sha", got)
			}
			if got := len(commit.GetParents()); got != 1 {
				t.Errorf("len(parents) = %d, want 1", got)
			}
		},
	},
	{
		op:        "github.git.commits.get",
		connector: "github",
		request:   func() proto.Message { return &githubv1.GetGitCommitRequest{} },
		response:  func() proto.Message { return &githubv1.GetGitCommitResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			commit := resp.(*githubv1.GetGitCommitResponse)
			if got := commit.GetSha(); got != "synthetic-commit-sha-3" {
				t.Errorf("sha = %q, want the synthetic commit sha", got)
			}
			if got := commit.GetMessage(); got != "Synthetic commit message" {
				t.Errorf("message = %q, want the synthetic commit message", got)
			}
			if got := len(commit.GetParents()); got != 1 {
				t.Errorf("len(parents) = %d, want 1", got)
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
		op:        "linear.projects.list",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.ListProjectsRequest{} },
		response:  func() proto.Message { return &linearv1.ListProjectsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			projects := resp.(*linearv1.ListProjectsResponse)
			if got := projects.GetNodes()[0].GetName(); got != "Synthetic pilot project" {
				t.Errorf("nodes[0].name = %q, want Synthetic pilot project", got)
			}
			if !projects.GetPageInfo().GetHasNextPage() {
				t.Errorf("page_info.has_next_page = %v, want true", projects.GetPageInfo().GetHasNextPage())
			}
			if got := projects.GetPageInfo().GetEndCursor(); got != "synthetic-cursor-1" {
				t.Errorf("page_info.end_cursor = %q, want synthetic-cursor-1", got)
			}
		},
	},
	{
		op:        "linear.cycles.list",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.ListCyclesRequest{} },
		response:  func() proto.Message { return &linearv1.ListCyclesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			cycles := resp.(*linearv1.ListCyclesResponse)
			if got := cycles.GetNodes()[0].GetNumber(); got != 7 {
				t.Errorf("nodes[0].number = %d, want 7", got)
			}
			if !cycles.GetNodes()[0].GetIsActive() {
				t.Errorf("nodes[0].is_active = %v, want true", cycles.GetNodes()[0].GetIsActive())
			}
		},
	},
	{
		op:        "linear.issues.update",
		connector: "linear",
		request:   func() proto.Message { return &linearv1.UpdateIssueRequest{} },
		response:  func() proto.Message { return &linearv1.UpdateIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			updated := resp.(*linearv1.UpdateIssueResponse)
			if !updated.GetSuccess() {
				t.Errorf("success = %v, want true", updated.GetSuccess())
			}
			if got := updated.GetIssue().GetTitle(); got != "Synthetic updated issue" {
				t.Errorf("issue.title = %q, want Synthetic updated issue", got)
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
		op:        "atlassian.jira.issues.create",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.CreateJiraIssueRequest{} },
		response:  func() proto.Message { return &atlassianv1.CreateJiraIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.CreateJiraIssueResponse)
			if got.GetId() != "30001" {
				t.Errorf("id = %q, want 30001", got.GetId())
			}
			if got.GetKey() != "SYN-301" {
				t.Errorf("key = %q, want SYN-301", got.GetKey())
			}
		},
	},
	{
		op:        "atlassian.jira.issues.comments.create",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.CreateJiraCommentRequest{} },
		response:  func() proto.Message { return &atlassianv1.CreateJiraCommentResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.CreateJiraCommentResponse)
			if got.GetId() != "40001" {
				t.Errorf("id = %q, want 40001", got.GetId())
			}
			if got.GetBody() == nil || got.GetBody().GetFields()["type"].GetStringValue() != "doc" {
				t.Errorf("body.type = %v, want doc", got.GetBody())
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
		op:        "atlassian.confluence.pages.create",
		connector: "atlassian",
		request:   func() proto.Message { return &atlassianv1.CreateConfluencePageRequest{} },
		response:  func() proto.Message { return &atlassianv1.CreateConfluencePageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*atlassianv1.CreateConfluencePageResponse)
			if got.GetId() != 78643205 {
				t.Errorf("id = %d, want 78643205", got.GetId())
			}
			if got.GetTitle() != "Synthetic created page" {
				t.Errorf("title = %q, want Synthetic created page", got.GetTitle())
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
		op:        "notion.pages.create",
		connector: "notion",
		request:   func() proto.Message { return &notionv1.CreatePageRequest{} },
		response:  func() proto.Message { return &notionv1.CreatePageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			page := resp.(*notionv1.CreatePageResponse)
			if page.GetId() != "d4e5f6a7-0000-4000-8000-000000000005" {
				t.Errorf("id = %q, want the synthetic created page id", page.GetId())
			}
			if page.PublicUrl != nil {
				t.Error("public_url is set, want nil (null in the fixture)")
			}
		},
	},
	{
		op:        "notion.pages.update",
		connector: "notion",
		request:   func() proto.Message { return &notionv1.UpdatePageRequest{} },
		response:  func() proto.Message { return &notionv1.UpdatePageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			page := resp.(*notionv1.UpdatePageResponse)
			if page.GetLastEditedTime() != "2026-02-11T09:45:00.000Z" {
				t.Errorf("last_edited_time = %q, want the synthetic edit time", page.GetLastEditedTime())
			}
			if page.GetArchived() {
				t.Error("archived = true, want false")
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
		op:        "sentry.issues.update",
		connector: "sentry",
		request:   func() proto.Message { return &sentryv1.UpdateIssueRequest{} },
		response:  func() proto.Message { return &sentryv1.UpdateIssueResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			issue := resp.(*sentryv1.UpdateIssueResponse)
			if issue.GetStatus() != "resolved" {
				t.Errorf("status = %q, want resolved", issue.GetStatus())
			}
			if issue.ShareId != nil {
				t.Error("share_id is set, want nil (null in the fixture)")
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
	{
		op:        "gmail.messages.list",
		connector: "gmail",
		request:   func() proto.Message { return &gmailv1.ListMessagesRequest{} },
		response:  func() proto.Message { return &gmailv1.ListMessagesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*gmailv1.ListMessagesResponse)
			if len(got.GetMessages()) != 2 {
				t.Errorf("len(messages) = %d, want 2", len(got.GetMessages()))
			}
			if got.GetNextPageToken() != "synthetic-next-token" {
				t.Errorf("next_page_token = %q, want synthetic-next-token", got.GetNextPageToken())
			}
		},
	},
	{
		op:        "gmail.messages.get",
		connector: "gmail",
		request:   func() proto.Message { return &gmailv1.GetMessageRequest{} },
		response:  func() proto.Message { return &gmailv1.GetMessageResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*gmailv1.GetMessageResponse)
			if got.GetId() != "msg-1001" {
				t.Errorf("id = %q, want msg-1001", got.GetId())
			}
			if got.GetSnippet() != "Synthetic message preview text" {
				t.Errorf("snippet = %q, want Synthetic message preview text", got.GetSnippet())
			}
			if got.InternalDate != nil {
				t.Errorf("internal_date = %v, want nil (provider null)", got.InternalDate)
			}
		},
	},
	{
		op:        "gmail.drafts.create",
		connector: "gmail",
		request:   func() proto.Message { return &gmailv1.CreateDraftRequest{} },
		response:  func() proto.Message { return &gmailv1.CreateDraftResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*gmailv1.CreateDraftResponse)
			if got.GetId() != "draft-3001" {
				t.Errorf("id = %q, want draft-3001", got.GetId())
			}
			if got.GetMessage().GetId() != "msg-1002" {
				t.Errorf("message.id = %q, want msg-1002", got.GetMessage().GetId())
			}
			if got.GetMessage().GetThreadId() != "thr-2002" {
				t.Errorf("message.thread_id = %q, want thr-2002", got.GetMessage().GetThreadId())
			}
		},
	},
	{
		op:        "google_sheets.values.get",
		connector: "google_sheets",
		request:   func() proto.Message { return &googlesheetsv1.GetValuesRequest{} },
		response:  func() proto.Message { return &googlesheetsv1.GetValuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*googlesheetsv1.GetValuesResponse)
			if got.GetRange() != "Sheet1!A1:B2" {
				t.Errorf("range = %q, want Sheet1!A1:B2", got.GetRange())
			}
			if len(got.GetValues()) != 2 {
				t.Errorf("len(values) = %d, want 2", len(got.GetValues()))
			}
			if got.GetValues()[0].GetValues()[1].GetNullValue() != structpb.NullValue_NULL_VALUE {
				t.Errorf("values[0][1] = %v, want null (provider null)", got.GetValues()[0].GetValues()[1])
			}
		},
	},
	{
		op:        "google_sheets.values.update",
		connector: "google_sheets",
		request:   func() proto.Message { return &googlesheetsv1.UpdateValuesRequest{} },
		response:  func() proto.Message { return &googlesheetsv1.UpdateValuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*googlesheetsv1.UpdateValuesResponse)
			if got.GetUpdatedRange() != "Sheet1!A1:B2" {
				t.Errorf("updated_range = %q, want Sheet1!A1:B2", got.GetUpdatedRange())
			}
			if got.GetUpdatedCells() != 4 {
				t.Errorf("updated_cells = %d, want 4", got.GetUpdatedCells())
			}
		},
	},
	{
		op:        "google_sheets.spreadsheets.get",
		connector: "google_sheets",
		request:   func() proto.Message { return &googlesheetsv1.GetSpreadsheetRequest{} },
		response:  func() proto.Message { return &googlesheetsv1.GetSpreadsheetResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*googlesheetsv1.GetSpreadsheetResponse)
			if got.GetSpreadsheetId() != "synthetic-sheet-1" {
				t.Errorf("spreadsheet_id = %q, want synthetic-sheet-1", got.GetSpreadsheetId())
			}
			if got.GetProperties().GetTitle() != "Synthetic sheet" {
				t.Errorf("properties.title = %q, want Synthetic sheet", got.GetProperties().GetTitle())
			}
		},
	},
	{
		op:        "google_sheets.values.append",
		connector: "google_sheets",
		request:   func() proto.Message { return &googlesheetsv1.AppendValuesRequest{} },
		response:  func() proto.Message { return &googlesheetsv1.AppendValuesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*googlesheetsv1.AppendValuesResponse)
			if got.GetSpreadsheetId() != "synthetic-sheet-1" {
				t.Errorf("spreadsheet_id = %q, want synthetic-sheet-1", got.GetSpreadsheetId())
			}
			if got := got.GetTableRange(); got != "Sheet1!A1:B2" {
				t.Errorf("table_range = %q, want Sheet1!A1:B2", got)
			}
			if got := got.GetUpdates().GetUpdatedRange(); got != "Sheet1!A3:B4" {
				t.Errorf("updates.updated_range = %q, want Sheet1!A3:B4", got)
			}
			if got := got.GetUpdates().GetUpdatedCells(); got != 4 {
				t.Errorf("updates.updated_cells = %d, want 4", got)
			}
		},
	},
	{
		op:        "google_calendar.calendars.list",
		connector: "google_calendar",
		request:   func() proto.Message { return &googlecalendarv1.ListCalendarsRequest{} },
		response:  func() proto.Message { return &googlecalendarv1.ListCalendarsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*googlecalendarv1.ListCalendarsResponse)
			if len(got.GetItems()) != 2 {
				t.Errorf("len(items) = %d, want 2", len(got.GetItems()))
			}
			if !got.GetItems()[0].GetPrimary() {
				t.Error("items[0].primary = false, want true")
			}
			if got.GetItems()[1].Description != nil {
				t.Error("items[1].description is set, want nil (provider null)")
			}
			if got.GetNextPageToken() != "synthetic-next-token" {
				t.Errorf("next_page_token = %q, want synthetic-next-token", got.GetNextPageToken())
			}
		},
	},
	{
		op:        "google_calendar.events.list",
		connector: "google_calendar",
		request:   func() proto.Message { return &googlecalendarv1.ListEventsRequest{} },
		response:  func() proto.Message { return &googlecalendarv1.ListEventsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			events := resp.(*googlecalendarv1.ListEventsResponse)
			if len(events.GetItems()) != 2 {
				t.Errorf("len(items) = %d, want 2", len(events.GetItems()))
			}
			if got := events.GetItems()[0].GetId(); got != "9007199254740993" {
				t.Errorf("items[0].id = %q, want the >2^53 provider id as a string", got)
			}
			if got := events.GetItems()[0].GetStart().GetDateTime(); got != "2026-01-15T10:30:00+01:00" {
				t.Errorf("items[0].start.date_time = %q, want the synthetic start", got)
			}
			if got := events.GetItems()[1].GetStart().GetDate(); got != "2026-01-20" {
				t.Errorf("items[1].start.date = %q, want the all-day date", got)
			}
		},
	},
	{
		op:        "google_calendar.events.get",
		connector: "google_calendar",
		request:   func() proto.Message { return &googlecalendarv1.GetEventRequest{} },
		response:  func() proto.Message { return &googlecalendarv1.GetEventResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			event := resp.(*googlecalendarv1.GetEventResponse)
			if event.GetId() != "9007199254740993" {
				t.Errorf("id = %q, want the >2^53 provider id as a string", event.GetId())
			}
			if event.Description != nil {
				t.Error("description is set, want nil (provider null)")
			}
			if got := event.GetLocation(); got != "Synthetic Conference Room" {
				t.Errorf("location = %q, want Synthetic Conference Room", got)
			}
		},
	},
	{
		op:        "google_calendar.events.create",
		connector: "google_calendar",
		request:   func() proto.Message { return &googlecalendarv1.CreateEventRequest{} },
		response:  func() proto.Message { return &googlecalendarv1.CreateEventResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			event := resp.(*googlecalendarv1.CreateEventResponse)
			if event.GetId() != "evt-synthetic-created-1" {
				t.Errorf("id = %q, want evt-synthetic-created-1", event.GetId())
			}
			if got := event.GetStart().GetDateTime(); got != "2026-02-10T14:00:00+01:00" {
				t.Errorf("start.date_time = %q, want the synthetic start", got)
			}
		},
	},
	{
		op:        "google_calendar.events.update",
		connector: "google_calendar",
		request:   func() proto.Message { return &googlecalendarv1.UpdateEventRequest{} },
		response:  func() proto.Message { return &googlecalendarv1.UpdateEventResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			event := resp.(*googlecalendarv1.UpdateEventResponse)
			if got := event.GetSummary(); got != "Synthetic updated event" {
				t.Errorf("summary = %q, want Synthetic updated event", got)
			}
			if event.Location != nil {
				t.Error("location is set, want nil (provider null)")
			}
		},
	},
	{
		op:        "notion.search",
		connector: "notion",
		request:   func() proto.Message { return &notionv1.SearchRequest{} },
		response:  func() proto.Message { return &notionv1.SearchResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			got := resp.(*notionv1.SearchResponse)
			if len(got.GetResults()) != 2 {
				t.Errorf("len(results) = %d, want 2", len(got.GetResults()))
			}
			if got := got.GetResults()[0].GetFields()["id"].GetStringValue(); got != "9007199254740993" {
				t.Errorf("results[0].id = %q, want the >2^53 provider id as a string", got)
			}
			if got := got.GetResults()[0].GetFields()["cover"].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
				t.Errorf("results[0].cover = %v, want null", got)
			}
			if got := got.GetResults()[1].GetFields()["object"].GetStringValue(); got != "database" {
				t.Errorf("results[1].object = %q, want database", got)
			}
			if got.GetNextCursor() != "synthetic-next-token" {
				t.Errorf("next_cursor = %q, want synthetic-next-token", got.GetNextCursor())
			}
			if !got.GetHasMore() {
				t.Error("has_more = false, want true")
			}
		},
	},
	{
		op:        "omni.documents.list",
		connector: "omni",
		request:   func() proto.Message { return &omniv1.ListDocumentsRequest{} },
		response:  func() proto.Message { return &omniv1.ListDocumentsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			page := resp.(*omniv1.ListDocumentsResponse)
			if got := len(page.GetDocuments()); got != 2 {
				t.Errorf("len(documents) = %d, want 2", got)
			}
			if got := page.GetDocuments()[0].GetName(); got != "Q3 Revenue — ünïcodé ✓" {
				t.Errorf("documents[0].name = %q, want the unicode document name", got)
			}
			if got := page.GetNextCursor(); got != "synthetic-cursor-2" {
				t.Errorf("next_cursor = %q, want synthetic-cursor-2", got)
			}
		},
	},
	{
		op:        "omni.queries.run",
		connector: "omni",
		request:   func() proto.Message { return &omniv1.RunQueryRequest{} },
		response:  func() proto.Message { return &omniv1.RunQueryResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			result := resp.(*omniv1.RunQueryResponse)
			if got := len(result.GetRows()); got != 2 {
				t.Fatalf("len(rows) = %d, want 2", got)
			}
			if got := result.GetRows()[0].GetFields()["region"].GetStringValue(); got != "EMEA" {
				t.Errorf("rows[0].region = %q, want EMEA", got)
			}
			if got := result.GetRows()[0].GetFields()["revenue"].GetNumberValue(); got != 1234567.89 {
				t.Errorf("rows[0].revenue = %v, want 1234567.89", got)
			}
		},
	},
	{
		op:        "power_bi.datasets.list",
		connector: "power_bi",
		request:   func() proto.Message { return &powerbiv1.ListDatasetsRequest{} },
		response:  func() proto.Message { return &powerbiv1.ListDatasetsResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			page := resp.(*powerbiv1.ListDatasetsResponse)
			if got := len(page.GetValue()); got != 2 {
				t.Errorf("len(value) = %d, want 2", got)
			}
			dataset := page.GetValue()[0]
			if got := dataset.GetName(); got != "Sales Marketing" {
				t.Errorf("value[0].name = %q, want Sales Marketing", got)
			}
			if got := dataset.GetConfiguredBy(); got != "synthetic-admin@example.invalid" {
				t.Errorf("value[0].configured_by = %q, want the PascalCase wire key decoded", got)
			}
			if got := dataset.GetTargetStorageMode(); got != "Abf" {
				t.Errorf("value[0].target_storage_mode = %q, want Abf", got)
			}
		},
	},
	{
		op:        "power_bi.datasets.execute_queries",
		connector: "power_bi",
		request:   func() proto.Message { return &powerbiv1.ExecuteQueriesRequest{} },
		response:  func() proto.Message { return &powerbiv1.ExecuteQueriesResponse{} },
		spotCheck: func(t *testing.T, resp proto.Message) {
			result := resp.(*powerbiv1.ExecuteQueriesResponse)
			if got := len(result.GetResults()); got != 1 {
				t.Fatalf("len(results) = %d, want 1", got)
			}
			tables := result.GetResults()[0].GetTables()
			if got := len(tables); got != 1 {
				t.Fatalf("len(results[0].tables) = %d, want 1", got)
			}
			if got := len(tables[0].GetRows()); got != 2 {
				t.Fatalf("len(results[0].tables[0].rows) = %d, want 2", got)
			}
			if got := tables[0].GetRows()[0].GetFields()["Region"].GetStringValue(); got != "EMEA" {
				t.Errorf("rows[0].Region = %q, want EMEA", got)
			}
			if got := tables[0].GetRows()[0].GetFields()["Sales Amount"].GetNumberValue(); got != 1234567.89 {
				t.Errorf("rows[0].Sales Amount = %v, want 1234567.89", got)
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

// TestDatabaseDialectFixtureRows covers the per-dialect query.run fixtures:
// each engine's dialect enum on the request, one engine-specific typed value
// per dialect (jsonb/timestamptz, bigint unsigned/datetime(6), nvarchar
// unicode/datetimeoffset, NUMBER(38) and VARIANT as strings, uuid), the NULL
// cell, and the drift gate: strict unmarshal must fail because every fixture
// carries one unknown extra field.
func TestDatabaseDialectFixtureRows(t *testing.T) {
	type dialectCellCheck struct {
		row, col int
		want     string
		null     bool
	}
	for _, tc := range []struct {
		name     string
		dialect  databasev1.Dialect
		col2Type string
		checks   []dialectCellCheck
	}{
		{
			name:     "postgres",
			dialect:  databasev1.Dialect_DIALECT_POSTGRES,
			col2Type: "timestamptz",
			checks: []dialectCellCheck{
				{row: 0, col: 0, want: "9007199254740993"},
				{row: 0, col: 1, want: `{"owner": "synthetic", "tags": ["ops"]}`},
				{row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00"},
				{row: 0, col: 3, want: "12345678901234567890.123456789012345678"},
				{row: 1, col: 1, null: true},
			},
		},
		{
			name:     "mysql",
			dialect:  databasev1.Dialect_DIALECT_MYSQL,
			col2Type: "datetime(6)",
			checks: []dialectCellCheck{
				{row: 0, col: 0, want: "18446744073709551615"},
				{row: 0, col: 1, want: `{"owner": "synthetic", "tags": ["ops"]}`},
				{row: 0, col: 2, want: "2026-01-15 10:30:00.123456"},
				{row: 0, col: 3, want: "12345678901234567890.123456"},
				{row: 1, col: 1, null: true},
			},
		},
		{
			name:     "mssql",
			dialect:  databasev1.Dialect_DIALECT_MSSQL,
			col2Type: "datetimeoffset",
			checks: []dialectCellCheck{
				{row: 0, col: 0, want: "9007199254740993"},
				{row: 0, col: 1, want: "München — ünïcodé ✓"},
				{row: 0, col: 2, want: "2026-01-15T10:30:00.1234567+02:00"},
				{row: 0, col: 3, want: "12345678901234567890.1234567890"},
				{row: 1, col: 2, null: true},
			},
		},
		{
			name:     "snowflake",
			dialect:  databasev1.Dialect_DIALECT_SNOWFLAKE,
			col2Type: "NUMBER(38,18)",
			checks: []dialectCellCheck{
				{row: 0, col: 0, want: "99999999999999999999999999999999999999"},
				{row: 0, col: 1, want: "Big-ID Account — ünïcodé ✓"},
				{row: 0, col: 2, want: "12345678901234567890.123456789012345678"},
				{row: 0, col: 3, want: `{"owner": "synthetic", "tags": ["ops"]}`},
				{row: 1, col: 3, null: true},
			},
		},
		{
			name:     "supabase",
			dialect:  databasev1.Dialect_DIALECT_SUPABASE,
			col2Type: "timestamptz",
			checks: []dialectCellCheck{
				{row: 0, col: 0, want: "6f9619ff-8b86-d011-b42d-00c04fc964ff"},
				{row: 0, col: 1, want: `{"owner": "synthetic", "tags": ["ops"]}`},
				{row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00"},
				{row: 0, col: 3, want: "12345678901234567890.123456789012345678"},
				{row: 1, col: 1, null: true},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reqData := readFixture(t, "database", "query.run."+tc.name+".request.json")
			var req databasev1.RunQueryRequest
			if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(reqData, &req); err != nil {
				t.Fatalf("request fixture did not unmarshal: %v", err)
			}
			if got := req.GetDialect(); got != tc.dialect {
				t.Errorf("dialect = %v, want %v", got, tc.dialect)
			}

			respData := readFixture(t, "database", "query.run."+tc.name+".response.json")
			resp := &databasev1.RunQueryResponse{}
			if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(respData, resp); err != nil {
				t.Fatalf("response fixture did not unmarshal: %v", err)
			}
			if got := len(resp.GetColumns()); got != 4 {
				t.Fatalf("len(columns) = %d, want 4", got)
			}
			if got := resp.GetColumns()[2].GetType(); got != tc.col2Type {
				t.Errorf("columns[2].type = %q, want %q", got, tc.col2Type)
			}
			if got := len(resp.GetRows()); got != 2 {
				t.Fatalf("len(rows) = %d, want 2", got)
			}
			for _, c := range tc.checks {
				value := resp.GetRows()[c.row].GetValues()[c.col]
				if c.null {
					if got := value.GetNullValue(); got != structpb.NullValue_NULL_VALUE {
						t.Errorf("rows[%d][%d] = %v, want null", c.row, c.col, value)
					}
					continue
				}
				if got := value.GetStringValue(); got != c.want {
					t.Errorf("rows[%d][%d] = %q, want %q", c.row, c.col, got, c.want)
				}
			}
			if resp.GetTruncated() {
				t.Error("truncated = true, want false")
			}
			if err := (protojson.UnmarshalOptions{}).Unmarshal(respData, &databasev1.RunQueryResponse{}); err == nil {
				t.Errorf("strict unmarshal of query.run.%s.response.json succeeded; every per-dialect response fixture must carry one unknown extra field so the drift gate detects provider additions", tc.name)
			}
		})
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

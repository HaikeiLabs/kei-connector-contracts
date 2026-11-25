# Copyright 2026 Haikei Labs
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Gate: every pilot operation fixture unmarshals in Python.

Every registered pilot operation needs synthetic request and response
fixtures under fixtures/operations/<connector>/, and both must parse into
the generated python types with json_format.Parse(ignore_unknown_fields=True)
(the same lenient semantics as the Go DiscardUnknown gate). One spot assert
per op proves a field survived unmarshal with its synthetic value intact.
"""

from pathlib import Path

import pytest
from google.protobuf import json_format

from kei.connectors.atlassian.v1 import confluence_pb2, jira_pb2
from kei.connectors.airtable.v1 import records_pb2 as airtable_records_pb2
from kei.connectors.database.v1 import query_pb2
from kei.connectors.github.v1 import (
    actions_pb2,
    checks_pb2,
    commits_pb2,
    contents_pb2,
    git_pb2,
    issues_pb2,
    pulls_pb2,
    repos_pb2,
)
from kei.connectors.gmail.v1 import messages_pb2
from kei.connectors.google_calendar.v1 import calendars_pb2, events_pb2
from kei.connectors.google_sheets.v1 import spreadsheets_pb2
from kei.connectors.http_api.v1 import records_pb2
from kei.connectors.hubspot.v1 import contacts_pb2, deals_pb2
from kei.connectors.linear.v1 import (
    comments_pb2,
    cycles_pb2,
    issues_pb2 as linear_issues_pb2,
    projects_pb2,
    teams_pb2,
)
from kei.connectors.omni.v1 import documents_pb2, queries_pb2
from kei.connectors.power_bi.v1 import datasets_pb2
from kei.connectors.notion.v1 import databases_pb2, pages_pb2, search_pb2
from kei.connectors.sentry.v1 import issues_pb2 as sentry_issues_pb2

FIXTURES = Path(__file__).resolve().parents[3] / "fixtures" / "operations"

INT64_SPOT = 9007199254740993  # 2^53 + 1: must not go through float64


def _fixture(connector: str, name: str) -> str:
    return (FIXTURES / connector / name).read_text()


def _op_file(op: str, kind: str) -> str:
    return op.split(".", 1)[1] + "." + kind + ".json"


def _spot_github_issues_get(response):
    assert response.number == 1347
    assert response.id == INT64_SPOT


def _spot_github_issues_list(response):
    assert len(response.issues) == 2


def _spot_github_contents_put(response):
    assert response.content.name == "notes.txt"


def _spot_github_pulls_get(response):
    assert response.number == 42
    assert response.diff_url == "https://example.invalid/example-org/example-repo/pull/42.diff"
    assert response.patch_url == "https://example.invalid/example-org/example-repo/pull/42.patch"


def _spot_github_pulls_list(response):
    pull = response.pulls[0]
    assert pull.head.ref == "feature/x"
    assert pull.base.ref == "main"
    assert pull.author.login == "example-user"
    assert pull.labels[0].name == "synthetic-label"
    assert pull.id == 9007199254740999
    assert response.next_page == 2


def _spot_github_actions_runs_list(response):
    assert response.workflow_runs[0].head_branch == "feature/x"
    assert response.workflow_runs[0].id == 9007199254741001
    assert response.workflow_runs[0].actor.login == "example-user"
    assert response.workflow_runs[0].head_sha == "synthetic-head-sha-1"


def _spot_github_actions_workflows_runs_list(response):
    assert len(response.workflow_runs) == 1
    assert response.workflow_runs[0].actor.login == "example-user"
    assert response.workflow_runs[0].head_sha == "synthetic-head-sha-1"


def _spot_github_repos_contents_get(response):
    assert response.encoding == "base64"
    assert response.content == "c3ludGhldGljLWNvbnRlbnQ="


def _spot_github_repos_contents_list(response):
    assert len(response.entries) == 2
    assert response.entries[0].type == "file"
    assert response.entries[1].type == "dir"


def _spot_github_repos_commits_list(response):
    assert response.commits[0].sha == "synthetic-commit-sha-3"
    assert response.commits[0].author.login == "example-user"


def _spot_github_repos_commits_get(response):
    assert response.commit.message == "Synthetic commit message"
    assert response.author.login == "example-user"


def _spot_github_actions_runs_attempts_get(response):
    assert response.run_attempt == 2
    assert response.head_sha == "synthetic-head-sha-1"


def _spot_github_actions_runs_jobs_list(response):
    assert response.jobs[0].id == 9007199254741005
    assert response.jobs[0].runner_name == "synthetic-runner-1"


def _spot_github_issues_create(response):
    assert response.number == 42
    assert response.author.login == "example-user"
    assert response.labels[0].name == "synthetic-label"


def _spot_github_issues_comments_create(response):
    assert response.body == "Synthetic comment"
    assert response.author.login == "example-user"


def _spot_github_repos_get(response):
    assert response.id == 9007199254741011
    assert response.full_name == "example-org/example-repo"
    assert response.owner.login == "example-user"


def _spot_github_checks_runs_list(response):
    assert response.check_runs[0].id == 9007199254741012
    assert response.check_runs[0].conclusion == "success"
    assert response.total_count == 2


def _spot_github_issues_update(response):
    assert response.number == 42
    assert response.state == "closed"
    assert response.author.login == "example-user"


def _spot_github_pulls_create(response):
    assert response.number == 43
    assert response.draft is True
    assert response.head.ref == "feature/x"
    assert response.base.ref == "main"


def _spot_github_pulls_update(response):
    assert response.title == "Synthetic updated title"
    assert response.state == "open"


def _spot_github_pulls_reviews_list(response):
    assert response.reviews[0].state == "APPROVED"
    assert response.reviews[0].author.login == "example-reviewer"


def _spot_github_pulls_reviews_create(response):
    assert response.state == "APPROVED"
    assert response.id == 9007199254741006


def _spot_github_git_refs_get(response):
    assert response.ref == "refs/heads/feature/x"
    assert response.object.sha == "synthetic-ref-sha-1"


def _spot_github_git_refs_create(response):
    assert response.ref == "refs/heads/feature/y"
    assert response.object.type == "commit"


def _spot_github_git_refs_update(response):
    assert response.object.sha == "synthetic-ref-sha-2"


def _spot_github_git_blobs_create(response):
    assert response.sha == "synthetic-blob-sha-1"


def _spot_github_git_trees_create(response):
    assert response.sha == "synthetic-tree-sha-1"
    assert response.tree[0].path == "notes.txt"


def _spot_github_git_commits_create(response):
    assert response.sha == "synthetic-commit-sha-2"
    assert len(response.parents) == 1


def _spot_github_git_commits_get(response):
    assert response.sha == "synthetic-commit-sha-3"
    assert response.message == "Synthetic commit message"
    assert response.tree.sha == "synthetic-tree-sha-4"
    assert response.parents[0].sha == "synthetic-commit-sha-2"


def _spot_linear_issues_get(response):
    assert response.issue.identifier == "ENG-101"
    assert response.issue.labels[0].name == "synthetic-label"


def _spot_linear_issues_create(response):
    assert response.success is True


def _spot_linear_issues_search(response):
    assert len(response.nodes) == 2


def _spot_linear_teams_get(response):
    assert response.team.key == "ENG"


def _spot_linear_comments_create(response):
    assert response.success is True
    assert response.comment.body == "Synthetic comment"


def _spot_linear_projects_list(response):
    assert response.nodes[0].name == "Synthetic pilot project"
    assert response.page_info.has_next_page is True
    assert response.page_info.end_cursor == "synthetic-cursor-1"


def _spot_linear_cycles_list(response):
    assert response.nodes[0].number == 7
    assert response.nodes[0].is_active is True


def _spot_linear_issues_update(response):
    assert response.success is True
    assert response.issue.title == "Synthetic updated issue"


def _spot_gmail_messages_list(response):
    assert len(response.messages) == 2
    assert response.next_page_token == "synthetic-next-token"


def _spot_gmail_messages_get(response):
    assert response.id == "msg-1001"
    assert response.snippet == "Synthetic message preview text"
    assert not response.HasField("internal_date")


def _spot_gmail_drafts_create(response):
    assert response.id == "draft-3001"
    assert response.message.id == "msg-1002"
    assert response.message.thread_id == "thr-2002"


def _spot_google_sheets_values_get(response):
    assert response.range == "Sheet1!A1:B2"
    assert len(response.values) == 2
    assert response.values[0].values[1].WhichOneof("kind") == "null_value"


def _spot_google_sheets_values_update(response):
    assert response.updated_range == "Sheet1!A1:B2"
    assert response.updated_cells == 4


def _spot_google_sheets_spreadsheets_get(response):
    assert response.spreadsheet_id == "synthetic-sheet-1"
    assert response.properties.title == "Synthetic sheet"


def _spot_google_sheets_values_append(response):
    assert response.spreadsheet_id == "synthetic-sheet-1"
    assert response.table_range == "Sheet1!A1:B2"
    assert response.updates.updated_range == "Sheet1!A3:B4"
    assert response.updates.updated_cells == 4


def _spot_google_calendar_calendars_list(response):
    assert len(response.items) == 2
    assert response.items[0].primary is True
    assert not response.items[1].HasField("description")
    assert response.next_page_token == "synthetic-next-token"


def _spot_google_calendar_events_list(response):
    assert len(response.items) == 2
    assert response.items[0].id == "9007199254740993"
    assert response.items[0].start.date_time == "2026-01-15T10:30:00+01:00"
    assert response.items[1].start.date == "2026-01-20"
    assert not response.items[1].start.HasField("date_time")


def _spot_google_calendar_events_get(response):
    assert response.id == "9007199254740993"
    assert not response.HasField("description")
    assert response.location == "Synthetic Conference Room"


def _spot_google_calendar_events_create(response):
    assert response.id == "evt-synthetic-created-1"
    assert response.start.date_time == "2026-02-10T14:00:00+01:00"


def _spot_google_calendar_events_update(response):
    assert response.summary == "Synthetic updated event"
    assert not response.HasField("location")


def _spot_atlassian_jira_issues_get(response):
    assert response.key == "SYN-101"
    assert response.fields.summary == "Synthetic Jira issue"
    assert not response.fields.HasField("description")


def _spot_atlassian_jira_issues_create(response):
    assert response.id == "30001"
    assert response.key == "SYN-301"


def _spot_atlassian_jira_comments_create(response):
    assert response.id == "40001"
    assert response.body.fields["type"].string_value == "doc"
    assert response.created == "2026-02-11T12:00:00.000Z"


def _spot_atlassian_jira_issues_search(response):
    assert len(response.issues) == 2
    assert response.next_page_token == "synthetic-next-token"
    assert not response.issues[1].fields.HasField("description")


def _spot_atlassian_confluence_pages_get(response):
    assert response.id == 78643201
    assert response.title == "Synthetic Confluence page"
    assert not response.HasField("parent_id")


def _spot_atlassian_confluence_pages_create(response):
    assert response.id == 78643205
    assert response.title == "Synthetic created page"
    assert not response.HasField("parent_id")


def _spot_atlassian_confluence_pages_list(response):
    assert len(response.results) == 2
    assert response.results[0].title == "Synthetic first page"


def _spot_airtable_records_get(response):
    assert response.id == "recSynthetic00003"
    assert response.fields.fields["Name"].string_value == "Synthetic third task"
    assert response.fields.fields["Priority"].WhichOneof("kind") == "null_value"


def _spot_airtable_records_create(response):
    assert response.id == "recSynthetic00004"
    assert response.fields.fields["Status"].string_value == "Todo"


def _spot_airtable_records_update(response):
    assert response.fields.fields["Status"].string_value == "Done"
    assert response.fields.fields["Priority"].WhichOneof("kind") == "null_value"


def _spot_airtable_records_list(response):
    assert len(response.records) == 2
    assert response.offset == "synthetic-next-offset"
    assert response.records[1].fields.fields["Priority"].WhichOneof("kind") == "null_value"


def _spot_hubspot_contacts_get(response):
    assert response.id == "701"
    assert response.properties.fields["email"].string_value == "synthetic.contact@example.invalid"
    assert response.properties.fields["lastname"].WhichOneof("kind") == "null_value"


def _spot_hubspot_contacts_create(response):
    assert response.id == "701550001"
    assert response.properties.fields["email"].string_value == "synthetic.new@example.invalid"


def _spot_hubspot_contacts_update(response):
    assert response.properties.fields["lifecyclestage"].string_value == "customer"


def _spot_hubspot_contacts_list(response):
    assert len(response.results) == 2
    assert response.next_after == "synthetic-next-cursor"
    assert response.results[1].properties.fields["lastname"].WhichOneof("kind") == "null_value"


def _spot_hubspot_deals_list(response):
    assert len(response.results) == 2
    assert response.results[0].properties.fields["dealname"].string_value == "Synthetic first deal"
    assert response.results[0].properties.fields["description"].WhichOneof("kind") == "null_value"


def _spot_notion_pages_get(response):
    assert response.id == "5f8b7e4a-0000-4000-8000-000000000001"
    assert not response.HasField("public_url")


def _spot_notion_pages_create(response):
    assert response.id == "d4e5f6a7-0000-4000-8000-000000000005"
    assert not response.HasField("public_url")


def _spot_notion_pages_update(response):
    assert response.last_edited_time == "2026-02-11T09:45:00.000Z"
    assert response.archived is False


def _spot_notion_search(response):
    assert len(response.results) == 2
    assert response.results[0].fields["id"].string_value == "9007199254740993"
    assert response.results[0].fields["cover"].WhichOneof("kind") == "null_value"
    assert response.results[1].fields["object"].string_value == "database"
    assert response.next_cursor == "synthetic-next-token"
    assert response.has_more is True


def _spot_notion_databases_query(response):
    assert len(response.results) == 2
    assert response.next_cursor == "synthetic-next-cursor"
    assert response.has_more is True
    assert not response.results[1].HasField("public_url")


def _spot_sentry_issues_get(response):
    assert response.short_id == "PROJ-103"
    assert not response.HasField("share_id")


def _spot_sentry_issues_update(response):
    assert response.status == "resolved"
    assert response.short_id == "PROJ-103"
    assert not response.HasField("share_id")


def _spot_sentry_issues_list(response):
    assert len(response.issues) == 2
    assert response.issues[1].share_id == "synthetic-share-id-2"


def _spot_http_api_records_list(response):
    assert len(response.records) == 2


def _spot_http_api_records_get(response):
    assert response.record.fields["id"].string_value == "rec-1"


def _spot_database_query_run(response):
    assert len(response.columns) == 4
    assert response.columns[0].name == "id"
    assert response.columns[0].type == "int8"
    assert response.columns[3].type == "numeric"
    assert len(response.rows) == 2

    big = response.rows[0].values
    # JSON numbers are float64, so the >2^53 id and NUMERIC values travel as
    # strings; NULL travels as the protobuf null value.
    assert big[0].string_value == "9007199254740993"
    assert big[1].string_value == "Big-ID Account — ünïcodé ✓"
    assert big[2].WhichOneof("kind") == "null_value"
    assert big[3].string_value == "12345678901234567890.123456789012345678"

    small = response.rows[1].values
    assert small[3].string_value == "0.000000000000000001"
    assert response.truncated is False


def _spot_omni_documents_list(response):
    assert len(response.documents) == 2
    assert response.documents[0].name == "Q3 Revenue — ünïcodé ✓"
    assert response.next_cursor == "synthetic-cursor-2"


def _spot_omni_queries_run(response):
    assert len(response.rows) == 2
    assert response.rows[0].fields["region"].string_value == "EMEA"
    assert response.rows[0].fields["revenue"].number_value == 1234567.89


def _spot_power_bi_datasets_list(response):
    assert len(response.value) == 2
    dataset = response.value[0]
    assert dataset.name == "Sales Marketing"
    # ConfiguredBy and TargetStorageMode are PascalCase on the wire; the
    # fields set json_name to match.
    assert dataset.configured_by == "synthetic-admin@example.invalid"
    assert dataset.target_storage_mode == "Abf"


def _spot_power_bi_datasets_execute_queries(response):
    assert len(response.results) == 1
    assert len(response.results[0].tables) == 1
    rows = response.results[0].tables[0].rows
    assert len(rows) == 2
    assert rows[0].fields["Region"].string_value == "EMEA"
    assert rows[0].fields["Sales Amount"].number_value == 1234567.89


def _check_database_query_run_request(request):
    assert request.dialect == query_pb2.DIALECT_POSTGRES
    assert "accounts" in request.sql
    assert len(request.params) == 1
    assert request.params[0].string_value == "9007199254740993"
    assert request.max_rows == 100


# Per-dialect query.run fixtures: each engine's dialect enum, one
# engine-specific typed value (jsonb/timestamptz, bigint unsigned/datetime(6),
# nvarchar unicode/datetimeoffset, NUMBER(38) and VARIANT as strings, uuid),
# and the NULL cell. The >2^53 integers and precision-critical decimals travel
# as strings; NULL travels as the protobuf null value.
DATABASE_DIALECT_CASES = [
    (
        "postgres",
        query_pb2.DIALECT_POSTGRES,
        "timestamptz",
        [
            (0, 0, "9007199254740993", False),
            (0, 1, '{"owner": "synthetic", "tags": ["ops"]}', False),
            (0, 2, "2026-01-15T10:30:00.123456+00:00", False),
            (0, 3, "12345678901234567890.123456789012345678", False),
            (1, 1, None, True),
        ],
    ),
    (
        "mysql",
        query_pb2.DIALECT_MYSQL,
        "datetime(6)",
        [
            (0, 0, "18446744073709551615", False),
            (0, 1, '{"owner": "synthetic", "tags": ["ops"]}', False),
            (0, 2, "2026-01-15 10:30:00.123456", False),
            (0, 3, "12345678901234567890.123456", False),
            (1, 1, None, True),
        ],
    ),
    (
        "mssql",
        query_pb2.DIALECT_MSSQL,
        "datetimeoffset",
        [
            (0, 0, "9007199254740993", False),
            (0, 1, "München — ünïcodé ✓", False),
            (0, 2, "2026-01-15T10:30:00.1234567+02:00", False),
            (0, 3, "12345678901234567890.1234567890", False),
            (1, 2, None, True),
        ],
    ),
    (
        "snowflake",
        query_pb2.DIALECT_SNOWFLAKE,
        "NUMBER(38,18)",
        [
            (0, 0, "99999999999999999999999999999999999999", False),
            (0, 1, "Big-ID Account — ünïcodé ✓", False),
            (0, 2, "12345678901234567890.123456789012345678", False),
            (0, 3, '{"owner": "synthetic", "tags": ["ops"]}', False),
            (1, 3, None, True),
        ],
    ),
    (
        "supabase",
        query_pb2.DIALECT_SUPABASE,
        "timestamptz",
        [
            (0, 0, "6f9619ff-8b86-d011-b42d-00c04fc964ff", False),
            (0, 1, '{"owner": "synthetic", "tags": ["ops"]}', False),
            (0, 2, "2026-01-15T10:30:00.123456+00:00", False),
            (0, 3, "12345678901234567890.123456789012345678", False),
            (1, 1, None, True),
        ],
    ),
]


@pytest.mark.parametrize(
    "name,dialect,col2_type,checks",
    DATABASE_DIALECT_CASES,
    ids=[case[0] for case in DATABASE_DIALECT_CASES],
)
def test_database_dialect_fixture_rows(name, dialect, col2_type, checks):
    request = json_format.Parse(
        _fixture("database", "query.run.%s.request.json" % name),
        query_pb2.RunQueryRequest(),
        ignore_unknown_fields=True,
    )
    assert request.dialect == dialect

    response_file = "query.run.%s.response.json" % name
    response = json_format.Parse(
        _fixture("database", response_file),
        query_pb2.RunQueryResponse(),
        ignore_unknown_fields=True,
    )
    assert len(response.columns) == 4
    assert response.columns[2].type == col2_type
    assert len(response.rows) == 2
    for row, col, want, is_null in checks:
        value = response.rows[row].values[col]
        if is_null:
            assert value.WhichOneof("kind") == "null_value"
        else:
            assert value.string_value == want
    assert response.truncated is False

    # Drift gate: strict parsing must reject the unknown extra field each
    # per-dialect response fixture carries.
    with pytest.raises(json_format.ParseError):
        json_format.Parse(
            _fixture("database", response_file),
            query_pb2.RunQueryResponse(),
            ignore_unknown_fields=False,
        )


OPS = [
    ("github.issues.get", "github", issues_pb2.GetIssueRequest, issues_pb2.GetIssueResponse, _spot_github_issues_get),
    ("github.issues.list", "github", issues_pb2.ListIssuesRequest, issues_pb2.ListIssuesResponse, _spot_github_issues_list),
    ("github.contents.put", "github", contents_pb2.PutContentsRequest, contents_pb2.PutContentsResponse, _spot_github_contents_put),
    ("github.pulls.get", "github", pulls_pb2.GetPullRequestRequest, pulls_pb2.GetPullRequestResponse, _spot_github_pulls_get),
    ("github.pulls.list", "github", pulls_pb2.ListPullsRequest, pulls_pb2.ListPullsResponse, _spot_github_pulls_list),
    ("github.actions.runs.list", "github", actions_pb2.ListWorkflowRunsRequest, actions_pb2.ListWorkflowRunsResponse, _spot_github_actions_runs_list),
    ("github.actions.workflows.runs.list", "github", actions_pb2.ListWorkflowRunsByWorkflowRequest, actions_pb2.ListWorkflowRunsByWorkflowResponse, _spot_github_actions_workflows_runs_list),
    ("github.repos.contents.get", "github", contents_pb2.GetContentsRequest, contents_pb2.GetContentsResponse, _spot_github_repos_contents_get),
    ("github.repos.contents.list", "github", contents_pb2.ListContentsRequest, contents_pb2.ListContentsResponse, _spot_github_repos_contents_list),
    ("github.repos.commits.list", "github", commits_pb2.ListCommitsRequest, commits_pb2.ListCommitsResponse, _spot_github_repos_commits_list),
    ("github.repos.commits.get", "github", commits_pb2.GetCommitRequest, commits_pb2.GetCommitResponse, _spot_github_repos_commits_get),
    ("github.actions.runs.attempts.get", "github", actions_pb2.GetRunAttemptRequest, actions_pb2.GetRunAttemptResponse, _spot_github_actions_runs_attempts_get),
    ("github.actions.runs.jobs.list", "github", actions_pb2.ListRunJobsRequest, actions_pb2.ListRunJobsResponse, _spot_github_actions_runs_jobs_list),
    ("github.issues.create", "github", issues_pb2.CreateIssueRequest, issues_pb2.CreateIssueResponse, _spot_github_issues_create),
    ("github.issues.comments.create", "github", issues_pb2.CreateIssueCommentRequest, issues_pb2.CreateIssueCommentResponse, _spot_github_issues_comments_create),
    ("github.repos.get", "github", repos_pb2.GetRepositoryRequest, repos_pb2.GetRepositoryResponse, _spot_github_repos_get),
    ("github.checks.runs.list", "github", checks_pb2.ListCheckRunsRequest, checks_pb2.ListCheckRunsResponse, _spot_github_checks_runs_list),
    ("github.issues.update", "github", issues_pb2.UpdateIssueRequest, issues_pb2.UpdateIssueResponse, _spot_github_issues_update),
    ("github.pulls.create", "github", pulls_pb2.CreatePullRequestRequest, pulls_pb2.CreatePullRequestResponse, _spot_github_pulls_create),
    ("github.pulls.update", "github", pulls_pb2.UpdatePullRequestRequest, pulls_pb2.UpdatePullRequestResponse, _spot_github_pulls_update),
    ("github.pulls.reviews.list", "github", pulls_pb2.ListPullReviewsRequest, pulls_pb2.ListPullReviewsResponse, _spot_github_pulls_reviews_list),
    ("github.pulls.reviews.create", "github", pulls_pb2.CreatePullReviewRequest, pulls_pb2.CreatePullReviewResponse, _spot_github_pulls_reviews_create),
    ("github.git.refs.get", "github", git_pb2.GetRefRequest, git_pb2.GetRefResponse, _spot_github_git_refs_get),
    ("github.git.refs.create", "github", git_pb2.CreateRefRequest, git_pb2.CreateRefResponse, _spot_github_git_refs_create),
    ("github.git.refs.update", "github", git_pb2.UpdateRefRequest, git_pb2.UpdateRefResponse, _spot_github_git_refs_update),
    ("github.git.blobs.create", "github", git_pb2.CreateBlobRequest, git_pb2.CreateBlobResponse, _spot_github_git_blobs_create),
    ("github.git.trees.create", "github", git_pb2.CreateTreeRequest, git_pb2.CreateTreeResponse, _spot_github_git_trees_create),
    ("github.git.commits.create", "github", git_pb2.CreateCommitRequest, git_pb2.CreateCommitResponse, _spot_github_git_commits_create),
    ("github.git.commits.get", "github", git_pb2.GetGitCommitRequest, git_pb2.GetGitCommitResponse, _spot_github_git_commits_get),
    ("linear.issues.get", "linear", linear_issues_pb2.GetIssueRequest, linear_issues_pb2.GetIssueResponse, _spot_linear_issues_get),
    ("linear.issues.create", "linear", linear_issues_pb2.CreateIssueRequest, linear_issues_pb2.CreateIssueResponse, _spot_linear_issues_create),
    ("linear.issues.search", "linear", linear_issues_pb2.SearchIssuesRequest, linear_issues_pb2.SearchIssuesResponse, _spot_linear_issues_search),
    ("linear.teams.get", "linear", teams_pb2.GetTeamRequest, teams_pb2.GetTeamResponse, _spot_linear_teams_get),
    ("linear.comments.create", "linear", comments_pb2.CreateCommentRequest, comments_pb2.CreateCommentResponse, _spot_linear_comments_create),
    ("linear.projects.list", "linear", projects_pb2.ListProjectsRequest, projects_pb2.ListProjectsResponse, _spot_linear_projects_list),
    ("linear.cycles.list", "linear", cycles_pb2.ListCyclesRequest, cycles_pb2.ListCyclesResponse, _spot_linear_cycles_list),
    ("linear.issues.update", "linear", linear_issues_pb2.UpdateIssueRequest, linear_issues_pb2.UpdateIssueResponse, _spot_linear_issues_update),
    ("gmail.messages.list", "gmail", messages_pb2.ListMessagesRequest, messages_pb2.ListMessagesResponse, _spot_gmail_messages_list),
    ("gmail.messages.get", "gmail", messages_pb2.GetMessageRequest, messages_pb2.GetMessageResponse, _spot_gmail_messages_get),
    ("gmail.drafts.create", "gmail", messages_pb2.CreateDraftRequest, messages_pb2.CreateDraftResponse, _spot_gmail_drafts_create),
    ("google_sheets.values.get", "google_sheets", spreadsheets_pb2.GetValuesRequest, spreadsheets_pb2.GetValuesResponse, _spot_google_sheets_values_get),
    ("google_sheets.values.update", "google_sheets", spreadsheets_pb2.UpdateValuesRequest, spreadsheets_pb2.UpdateValuesResponse, _spot_google_sheets_values_update),
    ("google_sheets.spreadsheets.get", "google_sheets", spreadsheets_pb2.GetSpreadsheetRequest, spreadsheets_pb2.GetSpreadsheetResponse, _spot_google_sheets_spreadsheets_get),
    ("google_sheets.values.append", "google_sheets", spreadsheets_pb2.AppendValuesRequest, spreadsheets_pb2.AppendValuesResponse, _spot_google_sheets_values_append),
    ("google_calendar.calendars.list", "google_calendar", calendars_pb2.ListCalendarsRequest, calendars_pb2.ListCalendarsResponse, _spot_google_calendar_calendars_list),
    ("google_calendar.events.list", "google_calendar", events_pb2.ListEventsRequest, events_pb2.ListEventsResponse, _spot_google_calendar_events_list),
    ("google_calendar.events.get", "google_calendar", events_pb2.GetEventRequest, events_pb2.GetEventResponse, _spot_google_calendar_events_get),
    ("google_calendar.events.create", "google_calendar", events_pb2.CreateEventRequest, events_pb2.CreateEventResponse, _spot_google_calendar_events_create),
    ("google_calendar.events.update", "google_calendar", events_pb2.UpdateEventRequest, events_pb2.UpdateEventResponse, _spot_google_calendar_events_update),
    ("atlassian.jira.issues.get", "atlassian", jira_pb2.GetJiraIssueRequest, jira_pb2.GetJiraIssueResponse, _spot_atlassian_jira_issues_get),
    ("atlassian.jira.issues.create", "atlassian", jira_pb2.CreateJiraIssueRequest, jira_pb2.CreateJiraIssueResponse, _spot_atlassian_jira_issues_create),
    ("atlassian.jira.issues.comments.create", "atlassian", jira_pb2.CreateJiraCommentRequest, jira_pb2.CreateJiraCommentResponse, _spot_atlassian_jira_comments_create),
    ("atlassian.jira.issues.search", "atlassian", jira_pb2.SearchJiraIssuesRequest, jira_pb2.SearchJiraIssuesResponse, _spot_atlassian_jira_issues_search),
    ("atlassian.confluence.pages.get", "atlassian", confluence_pb2.GetConfluencePageRequest, confluence_pb2.GetConfluencePageResponse, _spot_atlassian_confluence_pages_get),
    ("atlassian.confluence.pages.create", "atlassian", confluence_pb2.CreateConfluencePageRequest, confluence_pb2.CreateConfluencePageResponse, _spot_atlassian_confluence_pages_create),
    ("atlassian.confluence.pages.list", "atlassian", confluence_pb2.ListConfluencePagesRequest, confluence_pb2.ListConfluencePagesResponse, _spot_atlassian_confluence_pages_list),
    ("airtable.records.get", "airtable", airtable_records_pb2.GetRecordRequest, airtable_records_pb2.GetRecordResponse, _spot_airtable_records_get),
    ("airtable.records.create", "airtable", airtable_records_pb2.CreateRecordRequest, airtable_records_pb2.CreateRecordResponse, _spot_airtable_records_create),
    ("airtable.records.update", "airtable", airtable_records_pb2.UpdateRecordRequest, airtable_records_pb2.UpdateRecordResponse, _spot_airtable_records_update),
    ("airtable.records.list", "airtable", airtable_records_pb2.ListRecordsRequest, airtable_records_pb2.ListRecordsResponse, _spot_airtable_records_list),
    ("hubspot.contacts.get", "hubspot", contacts_pb2.GetContactRequest, contacts_pb2.GetContactResponse, _spot_hubspot_contacts_get),
    ("hubspot.contacts.create", "hubspot", contacts_pb2.CreateContactRequest, contacts_pb2.CreateContactResponse, _spot_hubspot_contacts_create),
    ("hubspot.contacts.update", "hubspot", contacts_pb2.UpdateContactRequest, contacts_pb2.UpdateContactResponse, _spot_hubspot_contacts_update),
    ("hubspot.contacts.list", "hubspot", contacts_pb2.ListContactsRequest, contacts_pb2.ListContactsResponse, _spot_hubspot_contacts_list),
    ("hubspot.deals.list", "hubspot", deals_pb2.ListDealsRequest, deals_pb2.ListDealsResponse, _spot_hubspot_deals_list),
    ("notion.pages.get", "notion", pages_pb2.GetPageRequest, pages_pb2.GetPageResponse, _spot_notion_pages_get),
    ("notion.pages.create", "notion", pages_pb2.CreatePageRequest, pages_pb2.CreatePageResponse, _spot_notion_pages_create),
    ("notion.pages.update", "notion", pages_pb2.UpdatePageRequest, pages_pb2.UpdatePageResponse, _spot_notion_pages_update),
    ("notion.databases.query", "notion", databases_pb2.QueryDatabaseRequest, databases_pb2.QueryDatabaseResponse, _spot_notion_databases_query),
    ("notion.search", "notion", search_pb2.SearchRequest, search_pb2.SearchResponse, _spot_notion_search),
    ("sentry.issues.get", "sentry", sentry_issues_pb2.GetIssueRequest, sentry_issues_pb2.GetIssueResponse, _spot_sentry_issues_get),
    ("sentry.issues.update", "sentry", sentry_issues_pb2.UpdateIssueRequest, sentry_issues_pb2.UpdateIssueResponse, _spot_sentry_issues_update),
    ("sentry.issues.list", "sentry", sentry_issues_pb2.ListIssuesRequest, sentry_issues_pb2.ListIssuesResponse, _spot_sentry_issues_list),
    ("http_api.records.list", "http_api", records_pb2.ListRecordsRequest, records_pb2.ListRecordsResponse, _spot_http_api_records_list),
    ("http_api.records.get", "http_api", records_pb2.GetRecordRequest, records_pb2.GetRecordResponse, _spot_http_api_records_get),
    ("database.query.run", "database", query_pb2.RunQueryRequest, query_pb2.RunQueryResponse, _spot_database_query_run),
    ("omni.documents.list", "omni", documents_pb2.ListDocumentsRequest, documents_pb2.ListDocumentsResponse, _spot_omni_documents_list),
    ("omni.queries.run", "omni", queries_pb2.RunQueryRequest, queries_pb2.RunQueryResponse, _spot_omni_queries_run),
    ("power_bi.datasets.list", "power_bi", datasets_pb2.ListDatasetsRequest, datasets_pb2.ListDatasetsResponse, _spot_power_bi_datasets_list),
    ("power_bi.datasets.execute_queries", "power_bi", datasets_pb2.ExecuteQueriesRequest, datasets_pb2.ExecuteQueriesResponse, _spot_power_bi_datasets_execute_queries),
]


@pytest.mark.parametrize(
    "op,connector,request_cls,response_cls,spot",
    OPS,
    ids=[entry[0] for entry in OPS],
)
def test_fixtures_unmarshal(op, connector, request_cls, response_cls, spot):
    request = json_format.Parse(
        _fixture(connector, _op_file(op, "request")),
        request_cls(),
        ignore_unknown_fields=True,
    )
    assert request is not None
    if op == "database.query.run":
        _check_database_query_run_request(request)

    response = json_format.Parse(
        _fixture(connector, _op_file(op, "response")),
        response_cls(),
        ignore_unknown_fields=True,
    )
    assert response is not None

    spot(response)

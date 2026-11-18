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

"""HAI-990 gate: the pilot operation fixtures unmarshal in Python.

Every registered pilot operation needs synthetic request and response
fixtures under fixtures/operations/<connector>/, and both must parse into
the generated python types with json_format.Parse(ignore_unknown_fields=True)
(the same lenient semantics as the Go DiscardUnknown gate). One spot assert
per op proves a field survived unmarshal with its synthetic value intact.
"""

from pathlib import Path

import pytest
from google.protobuf import json_format

from kei.connectors.database.v1 import query_pb2
from kei.connectors.github.v1 import (
    actions_pb2,
    contents_pb2,
    git_pb2,
    issues_pb2,
    pulls_pb2,
)
from kei.connectors.http_api.v1 import records_pb2
from kei.connectors.linear.v1 import (
    comments_pb2,
    issues_pb2 as linear_issues_pb2,
    teams_pb2,
)

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


def _spot_github_issues_create(response):
    assert response.number == 42
    assert response.author.login == "example-user"
    assert response.labels[0].name == "synthetic-label"


def _spot_github_issues_comments_create(response):
    assert response.body == "Synthetic comment"
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


def _check_database_query_run_request(request):
    assert request.dialect == query_pb2.DIALECT_POSTGRES
    assert "accounts" in request.sql
    assert len(request.params) == 1
    assert request.params[0].string_value == "9007199254740993"
    assert request.max_rows == 100


OPS = [
    ("github.issues.get", "github", issues_pb2.GetIssueRequest, issues_pb2.GetIssueResponse, _spot_github_issues_get),
    ("github.issues.list", "github", issues_pb2.ListIssuesRequest, issues_pb2.ListIssuesResponse, _spot_github_issues_list),
    ("github.contents.put", "github", contents_pb2.PutContentsRequest, contents_pb2.PutContentsResponse, _spot_github_contents_put),
    ("github.pulls.get", "github", pulls_pb2.GetPullRequestRequest, pulls_pb2.GetPullRequestResponse, _spot_github_pulls_get),
    ("github.pulls.list", "github", pulls_pb2.ListPullsRequest, pulls_pb2.ListPullsResponse, _spot_github_pulls_list),
    ("github.actions.runs.list", "github", actions_pb2.ListWorkflowRunsRequest, actions_pb2.ListWorkflowRunsResponse, _spot_github_actions_runs_list),
    ("github.issues.create", "github", issues_pb2.CreateIssueRequest, issues_pb2.CreateIssueResponse, _spot_github_issues_create),
    ("github.issues.comments.create", "github", issues_pb2.CreateIssueCommentRequest, issues_pb2.CreateIssueCommentResponse, _spot_github_issues_comments_create),
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
    ("linear.issues.get", "linear", linear_issues_pb2.GetIssueRequest, linear_issues_pb2.GetIssueResponse, _spot_linear_issues_get),
    ("linear.issues.create", "linear", linear_issues_pb2.CreateIssueRequest, linear_issues_pb2.CreateIssueResponse, _spot_linear_issues_create),
    ("linear.issues.search", "linear", linear_issues_pb2.SearchIssuesRequest, linear_issues_pb2.SearchIssuesResponse, _spot_linear_issues_search),
    ("linear.teams.get", "linear", teams_pb2.GetTeamRequest, teams_pb2.GetTeamResponse, _spot_linear_teams_get),
    ("linear.comments.create", "linear", comments_pb2.CreateCommentRequest, comments_pb2.CreateCommentResponse, _spot_linear_comments_create),
    ("http_api.records.list", "http_api", records_pb2.ListRecordsRequest, records_pb2.ListRecordsResponse, _spot_http_api_records_list),
    ("http_api.records.get", "http_api", records_pb2.GetRecordRequest, records_pb2.GetRecordResponse, _spot_http_api_records_get),
    ("database.query.run", "database", query_pb2.RunQueryRequest, query_pb2.RunQueryResponse, _spot_database_query_run),
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

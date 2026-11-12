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

"""HAI-990 gate: the 10 pilot operation fixtures unmarshal in Python.

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
from kei.connectors.github.v1 import contents_pb2, issues_pb2, pulls_pb2
from kei.connectors.http_api.v1 import records_pb2
from kei.connectors.linear.v1 import issues_pb2 as linear_issues_pb2

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


def _spot_linear_issues_get(response):
    assert response.issue.identifier == "ENG-101"


def _spot_linear_issues_create(response):
    assert response.success is True


def _spot_linear_issues_search(response):
    assert len(response.nodes) == 2


def _spot_http_api_records_list(response):
    assert len(response.records) == 2


def _spot_http_api_records_get(response):
    assert response.record.fields["id"].string_value == "rec-1"


def _spot_database_query_run(response):
    assert len(response.rows) == 2


OPS = [
    ("github.issues.get", "github", issues_pb2.GetIssueRequest, issues_pb2.GetIssueResponse, _spot_github_issues_get),
    ("github.issues.list", "github", issues_pb2.ListIssuesRequest, issues_pb2.ListIssuesResponse, _spot_github_issues_list),
    ("github.contents.put", "github", contents_pb2.PutContentsRequest, contents_pb2.PutContentsResponse, _spot_github_contents_put),
    ("github.pulls.get", "github", pulls_pb2.GetPullRequestRequest, pulls_pb2.GetPullRequestResponse, _spot_github_pulls_get),
    ("linear.issues.get", "linear", linear_issues_pb2.GetIssueRequest, linear_issues_pb2.GetIssueResponse, _spot_linear_issues_get),
    ("linear.issues.create", "linear", linear_issues_pb2.CreateIssueRequest, linear_issues_pb2.CreateIssueResponse, _spot_linear_issues_create),
    ("linear.issues.search", "linear", linear_issues_pb2.SearchIssuesRequest, linear_issues_pb2.SearchIssuesResponse, _spot_linear_issues_search),
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

    response = json_format.Parse(
        _fixture(connector, _op_file(op, "response")),
        response_cls(),
        ignore_unknown_fields=True,
    )
    assert response is not None

    spot(response)

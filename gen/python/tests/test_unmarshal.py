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

"""HAI-990 gate: the 25 pilot operation fixtures unmarshal in Python.

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
from kei.connectors.github.v1 import contents_pb2, issues_pb2, pulls_pb2
from kei.connectors.http_api.v1 import records_pb2
from kei.connectors.hubspot.v1 import contacts_pb2, deals_pb2
from kei.connectors.linear.v1 import (
    comments_pb2,
    issues_pb2 as linear_issues_pb2,
    teams_pb2,
)
from kei.connectors.notion.v1 import databases_pb2, pages_pb2
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


def _spot_atlassian_jira_issues_get(response):
    assert response.key == "SYN-101"
    assert response.fields.summary == "Synthetic Jira issue"
    assert not response.fields.HasField("description")


def _spot_atlassian_jira_issues_search(response):
    assert len(response.issues) == 2
    assert response.next_page_token == "synthetic-next-token"
    assert not response.issues[1].fields.HasField("description")


def _spot_atlassian_confluence_pages_get(response):
    assert response.id == 78643201
    assert response.title == "Synthetic Confluence page"
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


def _spot_notion_databases_query(response):
    assert len(response.results) == 2
    assert response.next_cursor == "synthetic-next-cursor"
    assert response.has_more is True
    assert not response.results[1].HasField("public_url")


def _spot_sentry_issues_get(response):
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
    ("linear.issues.get", "linear", linear_issues_pb2.GetIssueRequest, linear_issues_pb2.GetIssueResponse, _spot_linear_issues_get),
    ("linear.issues.create", "linear", linear_issues_pb2.CreateIssueRequest, linear_issues_pb2.CreateIssueResponse, _spot_linear_issues_create),
    ("linear.issues.search", "linear", linear_issues_pb2.SearchIssuesRequest, linear_issues_pb2.SearchIssuesResponse, _spot_linear_issues_search),
    ("linear.teams.get", "linear", teams_pb2.GetTeamRequest, teams_pb2.GetTeamResponse, _spot_linear_teams_get),
    ("linear.comments.create", "linear", comments_pb2.CreateCommentRequest, comments_pb2.CreateCommentResponse, _spot_linear_comments_create),
    ("atlassian.jira.issues.get", "atlassian", jira_pb2.GetJiraIssueRequest, jira_pb2.GetJiraIssueResponse, _spot_atlassian_jira_issues_get),
    ("atlassian.jira.issues.search", "atlassian", jira_pb2.SearchJiraIssuesRequest, jira_pb2.SearchJiraIssuesResponse, _spot_atlassian_jira_issues_search),
    ("atlassian.confluence.pages.get", "atlassian", confluence_pb2.GetConfluencePageRequest, confluence_pb2.GetConfluencePageResponse, _spot_atlassian_confluence_pages_get),
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
    ("notion.databases.query", "notion", databases_pb2.QueryDatabaseRequest, databases_pb2.QueryDatabaseResponse, _spot_notion_databases_query),
    ("sentry.issues.get", "sentry", sentry_issues_pb2.GetIssueRequest, sentry_issues_pb2.GetIssueResponse, _spot_sentry_issues_get),
    ("sentry.issues.list", "sentry", sentry_issues_pb2.ListIssuesRequest, sentry_issues_pb2.ListIssuesResponse, _spot_sentry_issues_list),
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

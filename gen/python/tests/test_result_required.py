# Result-required checks: a write operation whose result_required paths are
# not all held by a 2xx response fails with a contract_violation error, not
# with a silent success. Mirrors gen/go/conformance/result_required_test.go.
import asyncio
import json

import pytest
from google.protobuf import json_format

from kei.connection import OperationError, http_connection
from kei.connectors.linear.v1 import issues_pb2 as linear_issues_pb2
from kei.connectors.linear.v1 import ops as linear_ops

from test_replay import REPLAY_TOKEN, _Client, _Server, read_fixture

BAD_FIXTURES = [
    ("success false", "issues.create.response.success-false.json"),
    ("issue missing", "issues.create.response.no-issue.json"),
    ("empty identifier", "issues.create.response.empty-identifier.json"),
]


def _server_for(root: str, fixture_name: str) -> _Server:
    fixture = json.loads(read_fixture("linear", fixture_name))
    return _Server(200, json.dumps({"data": {root: fixture}}).encode())


def _create(server: _Server):
    conn = http_connection(server.base, "bearer", REPLAY_TOKEN, _Client())
    queries = linear_ops.LinearQueries(conn)
    req = json_format.ParseDict(
        json.loads(read_fixture("linear", "issues.create.request.json")),
        linear_issues_pb2.CreateIssueRequest(),
    )
    return queries.create_issue(req)


@pytest.mark.parametrize("name,fixture_name", BAD_FIXTURES, ids=[n for n, _ in BAD_FIXTURES])
def test_result_required_rejects_bad_responses(name, fixture_name):
    server = _server_for("issueCreate", fixture_name)
    try:
        with pytest.raises(OperationError) as excinfo:
            asyncio.run(_create(server))
        error = excinfo.value
        assert error.kind == "contract_violation"
        assert error.op_id == "linear.issues.create"
        assert error.status == 200
    finally:
        server.close()


def test_result_required_keeps_good_response():
    server = _server_for("issueCreate", "issues.create.response.json")
    try:
        resp = asyncio.run(_create(server))
        assert resp.success is True
        assert resp.issue.identifier == "ENG-102"
    finally:
        server.close()

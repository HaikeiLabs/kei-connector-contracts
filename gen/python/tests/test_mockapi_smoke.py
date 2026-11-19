# Smoke tests: the generated Python operations run against the real
# contract-derived mock server (testenv/mockapi) — fixture responses,
# GraphQL root-field routing, and failure injection via the X-Mock-Fail
# header mapping to the generated error kinds.
#
# The mock is the Go binary from testenv/mockapi/cmd/mockapi: built once
# per module with the local Go toolchain (or taken from MOCKAPI_BIN) and
# started on a free port. Skipped where no Go toolchain is available.
import asyncio
import http.client
import os
import re
import shutil
import subprocess
import tempfile
import urllib.parse
from pathlib import Path

import pytest

from kei.connection import OperationError, http_connection
from kei.connectors.github.v1 import issues_pb2
from kei.connectors.github.v1 import ops as github_ops
from kei.connectors.linear.v1 import issues_pb2 as linear_issues_pb2
from kei.connectors.linear.v1 import ops as linear_ops

REPO_ROOT = Path(__file__).resolve().parents[3]
TOKEN = "smoke-token"

requires_go = pytest.mark.skipif(
    "MOCKAPI_BIN" not in os.environ and shutil.which("go") is None,
    reason="the Go toolchain (or MOCKAPI_BIN) is needed to build and run the mock",
)


class _Client:
    """An httpx-shaped client over stdlib http.client (the CI python
    environment has protobuf and pytest only), with room to inject extra
    request headers — that is how the failure-injection case reaches the
    mock without touching the generated code."""

    def __init__(self, extra_headers=None):
        self._extra = dict(extra_headers or {})

    class _Headers:
        def __init__(self, items):
            self._items = list(items)

        def multi_items(self):
            return self._items

    class _Response:
        def __init__(self, status_code, headers, content):
            self.status_code = status_code
            self.headers = headers
            self.content = content

    async def request(self, method, url, headers=None, content=None):
        parts = urllib.parse.urlsplit(url)
        conn = http.client.HTTPConnection(parts.hostname, parts.port, timeout=10)
        try:
            target = parts.path + ("?" + parts.query if parts.query else "")
            merged = dict(headers or {})
            merged.update(self._extra)
            conn.request(method, target, body=content, headers=merged)
            response = conn.getresponse()
            data = response.read()
            items = response.getheaders()
            status = response.status
        finally:
            conn.close()
        return _Client._Response(status, _Client._Headers(items), data)


@pytest.fixture(scope="module")
def mock_base():
    binary = os.environ.get("MOCKAPI_BIN")
    tmp = None
    if not binary:
        tmp = tempfile.mkdtemp(prefix="mockapi-smoke-")
        binary = os.path.join(tmp, "mockapi")
        subprocess.run(
            ["go", "-C", "testenv/mockapi", "build", "-o", binary, "./cmd/mockapi"],
            cwd=REPO_ROOT,
            check=True,
        )
    proc = subprocess.Popen(
        [
            binary,
            "-addr",
            "127.0.0.1:0",
            "-manifest",
            str(REPO_ROOT / "manifest" / "connector-manifest.v1.json"),
            "-fixtures",
            str(REPO_ROOT / "fixtures" / "operations"),
        ],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.PIPE,
        text=True,
    )
    try:
        # The server prints its bound address (127.0.0.1:0 resolves to a
        # free port) as its first stderr line, after net.Listen.
        line = proc.stderr.readline()
        if proc.poll() is not None or not line:
            raise RuntimeError("mockapi failed to start: " + proc.stderr.read())
        match = re.search(r"on (\S+)$", line.strip())
        if not match:
            raise RuntimeError("mockapi printed no address: " + line)
        yield "http://" + match.group(1)
    finally:
        proc.terminate()
        proc.wait(timeout=10)
        if tmp:
            shutil.rmtree(tmp, ignore_errors=True)


@requires_go
def test_get_issue_replays_fixture(mock_base):
    queries = github_ops.GithubQueries(http_connection(mock_base, "bearer", TOKEN, _Client()))
    resp = asyncio.run(
        queries.get_issue(issues_pb2.GetIssueRequest(owner="acme", repository="widgets", issue_number=1347))
    )
    assert resp.number == 1347
    # The fixture id is 2^53+1: it must survive the decode losslessly.
    assert resp.id == 9007199254740993
    assert resp.title == "Synthetic pilot issue"


@requires_go
def test_linear_graphql_replays_fixture(mock_base):
    queries = linear_ops.LinearQueries(http_connection(mock_base, "bearer", TOKEN, _Client()))
    resp = asyncio.run(
        queries.get_issue(linear_issues_pb2.GetIssueRequest(id="00000000-0000-0000-0000-000000000101"))
    )
    # GraphQL decoding unwraps the {"data": {"issue": ...}} envelope.
    assert resp.issue.identifier == "ENG-101"


@requires_go
def test_failure_injection_maps_unauthorized(mock_base):
    queries = github_ops.GithubQueries(
        http_connection(mock_base, "bearer", TOKEN, _Client(extra_headers={"X-Mock-Fail": "401"}))
    )
    with pytest.raises(OperationError) as exc:
        asyncio.run(
            queries.get_issue(issues_pb2.GetIssueRequest(owner="acme", repository="widgets", issue_number=1347))
        )
    assert exc.value.kind == "unauthorized"
    assert exc.value.status == 401
    assert exc.value.op_id == "github.issues.get"

# Replay tests: generated operations run against a local test server that
# checks the rendered method, path, query, body and credential placement,
# then returns the operation's response fixture. Mirrors
# gen/go/conformance/replay_test.go.
import asyncio
import http.client
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import urlsplit

import pytest
from google.protobuf import json_format

from kei.connection import Column, OperationError, QueryResult, http_connection
from kei.connectors.database.v1 import ops as database_ops
from kei.connectors.database.v1 import query_pb2
from kei.connectors.github.v1 import issues_pb2
from kei.connectors.github.v1 import ops as github_ops

REPLAY_TOKEN = "synthetic-token"
FIXTURES = Path(__file__).resolve().parents[3] / "fixtures" / "operations"


def read_fixture(connector: str, name: str) -> bytes:
    return (FIXTURES / connector / name).read_bytes()


class _Server:
    """Serves one fixed status and body for every request and records what
    it saw: the raw (still percent-encoded) target, method, Authorization
    and body."""

    def __init__(self, status: int, body: bytes):
        rec = {"method": "", "path": "", "auth": "", "body": b""}

        class Handler(BaseHTTPRequestHandler):
            def do_GET(self) -> None:
                self._handle()

            def do_POST(self) -> None:
                self._handle()

            def do_PUT(self) -> None:
                self._handle()

            def do_DELETE(self) -> None:
                self._handle()

            def _handle(self) -> None:
                rec["method"] = self.command
                rec["path"] = self.path
                rec["auth"] = self.headers.get("Authorization", "")
                length = int(self.headers.get("Content-Length") or 0)
                rec["body"] = self.rfile.read(length)
                self.send_response(status)
                self.send_header("Content-Type", "application/json")
                self.end_headers()
                self.wfile.write(body)

            def log_message(self, *args) -> None:
                pass

        self.rec = rec
        self._httpd = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        self._thread = threading.Thread(target=self._httpd.serve_forever, daemon=True)
        self._thread.start()

    @property
    def base(self) -> str:
        host, port = self._httpd.server_address[:2]
        return f"http://{host}:{port}"

    def close(self) -> None:
        self._httpd.shutdown()
        self._httpd.server_close()


class _Client:
    """An httpx-shaped client over stdlib http.client, so the replay runs in
    the CI python environment (protobuf and pytest only)."""

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
        parts = urlsplit(url)
        conn = http.client.HTTPConnection(parts.hostname, parts.port, timeout=10)
        try:
            target = parts.path + ("?" + parts.query if parts.query else "")
            conn.request(method, target, body=content, headers=dict(headers or {}))
            response = conn.getresponse()
            data = response.read()
            items = response.getheaders()
            status = response.status
        finally:
            conn.close()
        return _Client._Response(status, _Client._Headers(items), data)


def _github_queries(base: str):
    conn = http_connection(base, "bearer", REPLAY_TOKEN, _Client())
    return github_ops.GithubQueries(conn)


def test_get_issue_replay():
    server = _Server(200, read_fixture("github", "issues.get.response.json"))
    try:
        queries = _github_queries(server.base)
        resp = asyncio.run(
            queries.get_issue(
                issues_pb2.GetIssueRequest(owner="acme", repository="widgets", issue_number=1347)
            )
        )
        assert server.rec["method"] == "GET"
        assert server.rec["path"] == "/repos/acme/widgets/issues/1347"
        assert server.rec["auth"] == "Bearer " + REPLAY_TOKEN
        assert server.rec["body"] == b""
        assert resp.number == 1347
        # The fixture's id is 2^53+1 and it has an unknown key: both must
        # survive the decode losslessly.
        assert resp.id == 9007199254740993
        assert github_ops.OP_GET_ISSUE.id == "github.issues.get"
        assert github_ops.OP_GET_ISSUE.op_class.value == "read"
    finally:
        server.close()


def test_graphql_replay():
    # GraphQL decoding unwraps data: its single key is either a response
    # field (issue) or a root field wrapping the response (searchIssues).
    import json

    from kei.connectors.linear.v1 import issues_pb2 as linear_issues_pb2
    from kei.connectors.linear.v1 import ops as linear_ops

    cases = [
        ("issues.get", None, lambda q: q.get_issue(linear_issues_pb2.GetIssueRequest(id="i-1")),
         lambda r: r.issue.identifier == "ENG-101"),
        ("issues.search", "searchIssues", lambda q: q.search_issues(linear_issues_pb2.SearchIssuesRequest(term="x")),
         lambda r: len(r.nodes) > 0),
    ]
    for name, root, call, check in cases:
        fixture = json.loads(read_fixture("linear", f"{name}.response.json"))
        data = fixture if root is None else {root: fixture}
        server = _Server(200, json.dumps({"data": data}).encode())
        try:
            conn = http_connection(server.base, "bearer", REPLAY_TOKEN, _Client())
            resp = asyncio.run(call(linear_ops.LinearQueries(conn)))
            assert server.rec["method"] == "POST" and server.rec["path"] == "/graphql"
            assert check(resp), f"linear.{name}: {resp}"
        finally:
            server.close()


def test_path_escaping():
    server = _Server(200, read_fixture("github", "issues.get.response.json"))
    try:
        queries = _github_queries(server.base)
        asyncio.run(
            queries.get_issue(
                issues_pb2.GetIssueRequest(owner="a/b", repository="..", issue_number=7)
            )
        )
        # Reserved characters in path values never change the path shape.
        assert server.rec["path"] == "/repos/a%2Fb/%2E%2E/issues/7"
    finally:
        server.close()


def test_error_mapping_does_not_leak():
    server = _Server(429, b'{"message": "slow down ' + REPLAY_TOKEN.encode() + b'"}')
    try:
        queries = _github_queries(server.base)
        with pytest.raises(OperationError) as caught:
            asyncio.run(
                queries.get_issue(
                    issues_pb2.GetIssueRequest(owner="acme", repository="widgets", issue_number=1)
                )
            )
        err = caught.value
        assert err.kind == "rate_limited"
        assert err.retryable
        assert err.op_id == "github.issues.get"
        # Neither the credential nor the body may reach the error message.
        assert REPLAY_TOKEN not in str(err)
        assert "slow down" not in str(err)
    finally:
        server.close()


class _FakeSQL:
    """A SqlConnection whose rows reproduce the query.run response fixture."""

    def __init__(self):
        self.seen: tuple | None = None

    async def query(self, sql, params, max_rows):
        self.seen = (sql, list(params), max_rows)
        return QueryResult(
            columns=[
                Column(name="id", type="BIGINT"),
                Column(name="name", type="TEXT"),
                Column(name="note", type="TEXT"),
                Column(name="balance", type="NUMERIC"),
            ],
            rows=[
                [9007199254740993, "Big-ID Account — ünïcodé ✓", None, "12345678901234567890.123456789012345678"],
                [1, "Synthetic One", "Synthetic note", "0.000000000000000001"],
            ],
            truncated=False,
        )


def test_run_query_replay():
    fake = _FakeSQL()
    queries = database_ops.DatabaseQueries(fake)
    req = json_format.Parse(
        read_fixture("database", "query.run.request.json"), query_pb2.RunQueryRequest()
    )
    resp = asyncio.run(queries.run_query(req))
    sql, params, max_rows = fake.seen
    assert sql == req.sql
    # A 64-bit parameter crosses as the exact string, never a float.
    assert params == ["9007199254740993"]
    assert max_rows == 100
    assert len(resp.rows) == 2 and not resp.truncated
    first = resp.rows[0].values
    assert first[0].string_value == "9007199254740993"
    assert first[2].WhichOneof("kind") == "null_value"
    assert first[3].string_value == "12345678901234567890.123456789012345678"
    assert resp.rows[1].values[0].number_value == 1
    assert [c.name for c in resp.columns[:1]] == ["id"]
    assert resp.columns[3].name == "balance"

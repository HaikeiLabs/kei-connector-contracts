# The SQL runtime: the SqlConnection protocol generated query operations
# run on, and the request/response translation around it. Mirrors
# gen/go/kei/connection/{sql,sqlcall}.go: it takes the connection as an
# argument and never fetches, logs or prints credentials.
import math
from datetime import datetime, timezone
from typing import Any, NamedTuple, Protocol, Sequence

from google.protobuf import json_format, struct_pb2

from .http import _g_format, _is_set

# DefaultMaxRows bounds the rows a query returns when the request sets no
# max_rows.
DEFAULT_MAX_ROWS = 1000

# MaxExactInt is the largest integer a number cell represents exactly
# (2^53, the IEEE 754 double's safe integer range). A larger integer is
# delivered as a string instead.
MAX_EXACT_INT = 1 << 53


class Column(NamedTuple):
    name: str
    type: str


class QueryResult(NamedTuple):
    """One provider query result. rows cells use JSON types: None, bool,
    int, float, str; a cell the provider cannot represent as a number is a
    string."""

    columns: "list[Column]"
    rows: "list[list[Any]]"
    truncated: bool


class SqlConnection(Protocol):
    """An authenticated SQL connection to one provider database.

    query sends one statement with parameters and a row limit. max_rows
    bounds the rows returned; a provider that has more sets truncated. It
    raises on transport, auth and SQL errors."""

    async def query(self, sql: str, params: Sequence[Any], max_rows: int) -> QueryResult: ...


def _rfc3339(dt: datetime) -> str:
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    s = dt.astimezone(timezone.utc).isoformat()
    return s.replace("+00:00", "Z")


def _set_value(v, cell) -> None:
    """Put one cell into a struct_pb2.Value, mirroring ToValue."""
    if cell is None:
        v.null_value = struct_pb2.NULL_VALUE
    elif isinstance(cell, bool):
        v.bool_value = cell
    elif isinstance(cell, int):
        if cell > MAX_EXACT_INT or cell < -MAX_EXACT_INT:
            v.string_value = str(cell)
        else:
            v.number_value = float(cell)
    elif isinstance(cell, float):
        if math.isnan(cell) or math.isinf(cell):
            v.string_value = _g_format(cell)
        else:
            v.number_value = cell
    elif isinstance(cell, str):
        v.string_value = cell
    elif isinstance(cell, (bytes, bytearray)):
        v.string_value = bytes(cell).decode("utf-8", "replace")
    elif isinstance(cell, datetime):
        v.string_value = _rfc3339(cell)
    else:
        v.string_value = str(cell)


async def call_sql(conn, op, req, resp_type):
    """Run one SQL query operation: translate the request's sql, params
    (google.protobuf.Value) and max_rows into a query, then translate the
    result into a new resp_type message. A request without sql, params and
    max_rows fields raises ValueError."""
    op_id = getattr(op, "id", op)
    if conn is None:
        raise ValueError(f"connection: {op_id}: nil connection")
    fields = req.DESCRIPTOR.fields_by_name
    for name in ("sql", "params", "max_rows"):
        if name not in fields:
            raise ValueError(f"connection: {op_id}: request is not a query request")
    args = [json_format.MessageToDict(v) for v in req.params]
    max_rows = DEFAULT_MAX_ROWS
    if _is_set(req, fields["max_rows"]):
        if int(req.max_rows) > 0:
            max_rows = int(req.max_rows)
    try:
        result = await conn.query(req.sql, args, max_rows)
    except Exception as e:
        raise ValueError(f"connection: {op_id}: {e}") from e
    resp = resp_type()
    for c in result.columns:
        resp.columns.add(name=c.name, type=c.type.lower())
    for row in result.rows:
        values = resp.rows.add()
        for cell in row:
            _set_value(values.values.add(), cell)
    resp.truncated = result.truncated
    return resp

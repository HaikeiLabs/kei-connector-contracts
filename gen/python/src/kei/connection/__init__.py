# The kei.connection runtime: typed HTTP and SQL connections plus the
# request encoder and response decoder generated operations run on. It is
# hand-written, dependency-light and credential-safe: callers pass the
# credential in, nothing logs, prints or embeds it, and no error carries a
# credential or a response body.
#
# The HTTP connection ships an httpx adapter (http_connection, an optional
# dependency: pip install httpx); everything else is stdlib + protobuf.
from enum import Enum
from typing import NamedTuple

from .errors import OperationError, map_error
from .http import (
    DEFAULT_HEADER_NAME,
    MAX_RESPONSE_BYTES,
    ConnectionKind,
    HttpConnection,
    TransportError,
    call_http,
    http_connection,
)
from .path import render_path
from .sql import (
    DEFAULT_MAX_ROWS,
    MAX_EXACT_INT,
    Column,
    QueryResult,
    SqlConnection,
    call_sql,
)

__all__ = [
    "DEFAULT_HEADER_NAME",
    "DEFAULT_MAX_ROWS",
    "MAX_EXACT_INT",
    "MAX_RESPONSE_BYTES",
    "Column",
    "ConnectionKind",
    "HttpConnection",
    "Op",
    "OpClass",
    "OperationError",
    "QueryResult",
    "SqlConnection",
    "TransportError",
    "call_http",
    "call_sql",
    "http_connection",
    "map_error",
    "render_path",
]


class OpClass(str, Enum):
    """What an operation does to the provider: read, write, or
    unspecified."""

    UNSPECIFIED = "unspecified"
    READ = "read"
    WRITE = "write"


class Op(NamedTuple):
    """One operation a connector exposes: its capability id, class and
    whether it destroys or substantially changes provider state. The id is
    also the policy capability: connector:<id>."""

    id: str
    op_class: OpClass
    destructive: bool

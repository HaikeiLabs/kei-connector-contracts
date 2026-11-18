# The HTTP runtime: an authenticated connection over an httpx.AsyncClient,
# the canonical request encoder and the response decoder that generated
# operations run on. Mirrors gen/go/kei/connection/{http,call,encode}.go:
# it takes the credential as an argument and never fetches, logs or prints
# it, and no error it produces contains a credential or a response body.
import base64
import json
import math
import re
from decimal import Decimal
from typing import Any, Mapping, Protocol, Sequence, Union
from urllib.parse import parse_qs, quote_plus, urlencode, urlsplit

from google.protobuf import json_format

from .errors import (
    FORBIDDEN,
    NOT_FOUND,
    RATE_LIMITED,
    UNAUTHORIZED,
    OperationError,
    map_error,
)
from .path import render_path

# DefaultHeaderName is the header an HTTP_HEADER connection puts the
# credential in when the connection declares no header name.
DEFAULT_HEADER_NAME = "X-Api-Key"

# MaxResponseBytes bounds the response body request reads.
MAX_RESPONSE_BYTES = 32 << 20

_VALID_HEADER_NAME = re.compile("[!#$%&'*+.^_`|~\\w-]+")

_LINK_RE = re.compile(r"<([^>]*)>((?:\s*;\s*[^;,]*)*)")

# Query values: a scalar str, or a sequence for fields that repeat the key.
QueryValue = Union[str, Sequence[str]]


class ConnectionKind:
    """The HTTP connection kinds an http_connection accepts."""

    BEARER = "bearer"
    HEADER = "header"
    BASIC = "basic"


class TransportError(Exception):
    """A request failed at the transport layer. The message names the method
    and path only: the credential travels in headers, never the URL."""


class HttpConnection(Protocol):
    """An authenticated HTTP connection to one provider origin. The
    credential is already attached; callers pass a rendered, escaped path.

    request sends one request. path must be an absolute path (from
    render_path) with no query or fragment; it is appended to the base
    URL's path. query values are str, or a sequence of str for a field that
    repeats the key. A non-2xx status is not an error: map it with
    map_error. The returned header mapping has lower-case keys."""

    async def request(
        self,
        method: str,
        path: str,
        query: Mapping[str, QueryValue],
        body: "bytes | None",
    ) -> "tuple[int, Mapping[str, str], bytes]": ...


def http_connection(
    base_url: str,
    kind: str,
    credential: str,
    client: Any,
    header_name: str = "",
) -> HttpConnection:
    """Build an HTTP connection for kind:

      - "bearer": Authorization: Bearer <credential>;
      - "header": the credential in header_name (DefaultHeaderName when it
        is empty);
      - "basic": the credential is "user:password".

    base_url must be an absolute http or https URL with no user info, query
    or fragment. client is required: the connection makes no request except
    through it.
    """
    if client is None:
        raise ValueError("connection: http_connection: client is required")
    if credential == "":
        raise ValueError("connection: http_connection: empty credential")
    if any(c in credential for c in "\r\n\x00"):
        raise ValueError("connection: http_connection: credential has a control character")
    base = _parse_base(base_url)
    if kind == ConnectionKind.BEARER:
        headers = {"authorization": "Bearer " + credential}
    elif kind == ConnectionKind.HEADER:
        name = header_name or DEFAULT_HEADER_NAME
        if not _VALID_HEADER_NAME.fullmatch(name):
            raise ValueError(f"connection: http_connection: invalid header name {name!r}")
        headers = {name.lower(): credential}
    elif kind == ConnectionKind.BASIC:
        user, sep, password = credential.partition(":")
        if not sep or user == "":
            raise ValueError("connection: http_connection: a basic credential must be user:password")
        token = base64.b64encode(f"{user}:{password}".encode()).decode("ascii")
        headers = {"authorization": "Basic " + token}
    else:
        raise ValueError(f"connection: http_connection: kind {kind!r} is not an HTTP connection")
    return _HttpxConnection(base, headers, client)


def _parse_base(raw: str):
    try:
        u = urlsplit(raw)
    except ValueError as e:
        raise ValueError(f"connection: invalid base URL: {e}") from None
    if u.username is not None or u.password is not None:
        raise ValueError("connection: base URL must not carry user info")
    if u.scheme not in ("http", "https"):
        raise ValueError(f"connection: base URL {raw!r} must be http or https")
    if not u.hostname:
        raise ValueError(f"connection: base URL {raw!r} has no host")
    if u.query or u.fragment:
        raise ValueError(f"connection: base URL {raw!r} must not have a query or fragment")
    return u


class _HttpxConnection:
    """Holds the credential in a private mapping; str and repr never print
    it."""

    def __init__(self, base, headers: Mapping[str, str], client):
        self._base = base
        self._headers = dict(headers)
        self._client = client

    def __str__(self) -> str:
        return f"connection.HTTP(base={self._base.scheme}://{self._base.netloc})"

    def __repr__(self) -> str:
        return str(self)

    async def request(
        self,
        method: str,
        path: str,
        query: Mapping[str, QueryValue],
        body: "bytes | None",
    ) -> "tuple[int, Mapping[str, str], bytes]":
        if not path.startswith("/") or path.startswith("//") or "?" in path or "#" in path:
            raise ValueError(
                f"connection: path {path!r} must be an absolute path with no query or fragment"
            )
        target = self._base.scheme + "://" + self._base.netloc + self._base.path.rstrip("/") + path
        pairs: list[tuple[str, str]] = []
        for k in sorted(query):
            v = query[k]
            if isinstance(v, str):
                pairs.append((k, v))
            else:
                pairs.extend((k, item) for item in v)
        qs = urlencode(pairs, safe="-._~", quote_via=quote_plus)
        if qs:
            target += "?" + qs
        headers = {"accept": "application/json"}
        headers.update(self._headers)
        if body is not None:
            headers["content-type"] = "application/json"
        try:
            resp = await self._client.request(method, target, headers=headers, content=body)
        except Exception as e:
            raise TransportError(f"connection: {method} {path}: {e}") from e
        if len(resp.content) > MAX_RESPONSE_BYTES:
            raise TransportError(
                f"connection: {method} {path}: response larger than {MAX_RESPONSE_BYTES} bytes"
            )
        out: dict[str, str] = {}
        for name, value in resp.headers.multi_items():
            k = name.lower()
            out[k] = value if k not in out else out[k] + ", " + value
        return resp.status_code, out, resp.content


# --- Canonical request JSON -------------------------------------------------
#
# Shared by every language's runtime so the cross-language parity test can
# compare request bodies byte for byte:
#   - object keys sorted by byte order, no insignificant whitespace;
#   - integers (including 64-bit) as exact JSON numbers, never strings;
#   - strings with only the escapes JSON requires (no HTML escaping, UTF-8
#     kept as is);
#   - enums by name, bytes as standard base64;
#   - only populated fields; well-known types in their protobuf JSON form.

_STRING_ESCAPES = {'"': '\\"', "\\": "\\\\", "\n": "\\n", "\r": "\\r", "\t": "\\t"}


def _json_string(s: str) -> str:
    out = ['"']
    for ch in s:
        esc = _STRING_ESCAPES.get(ch)
        if esc is not None:
            out.append(esc)
        elif ch < " ":
            out.append("\\u%04x" % ord(ch))
        else:
            out.append(ch)
    out.append('"')
    return "".join(out)


def _json_float(f: float) -> str:
    # Go's encoding/json: shortest round-trip digits; exponent form only
    # when abs < 1e-6 or abs >= 1e21, with a two-digit exponent.
    s = repr(f)
    if s.endswith(".0"):
        s = s[:-2]
    if "e" in s:
        mant, _, exp = s.partition("e")
        e = int(exp)
        if 1e-6 <= abs(f) < 1e21:
            s = format(Decimal(mant + "e" + str(e)), "f")
        else:
            sign = "+" if e >= 0 else "-"
            s = f"{mant}e{sign}{abs(e):02d}"
    return s


def _dumps(v: Any) -> str:
    if v is None:
        return "null"
    if v is True:
        return "true"
    if v is False:
        return "false"
    if isinstance(v, str):
        return _json_string(v)
    if isinstance(v, int):
        return str(v)
    if isinstance(v, float):
        return _json_float(v)
    if isinstance(v, list):
        return "[" + ",".join(_dumps(e) for e in v) + "]"
    if isinstance(v, dict):
        return "{" + ",".join(_json_string(k) + ":" + _dumps(v[k]) for k in sorted(v)) + "}"
    raise TypeError(f"connection: {type(v).__name__} is not valid JSON")


def _default_json_name(name: str) -> str:
    out = []
    upper = False
    for ch in name:
        if ch == "_":
            upper = True
        elif upper and "a" <= ch <= "z":
            out.append(ch.upper())
            upper = False
        else:
            out.append(ch)
            upper = False
    return "".join(out)


def _wire_name(fd) -> str:
    """The key a field uses on the wire: its json_name when that differs
    from the protobuf default (lowerCamelCase), otherwise the proto field
    name."""
    if fd.json_name != _default_json_name(fd.name):
        return fd.json_name
    return fd.name


def _is_map(fd) -> bool:
    return (
        fd.message_type is not None
        and fd.message_type.GetOptions().map_entry
    )


def _is_set(msg, fd) -> bool:
    """Whether msg has fd populated, mirroring Go's Message.Has."""
    if fd.is_repeated:
        return len(getattr(msg, fd.name)) > 0
    if fd.message_type is not None or fd.has_presence:
        try:
            return msg.HasField(fd.name)
        except ValueError:
            return False
    return getattr(msg, fd.name) != fd.default_value


def _g_format(f: float) -> str:
    if math.isnan(f):
        return "NaN"
    if math.isinf(f):
        return "+Inf" if f > 0 else "-Inf"
    s = repr(f)
    if s.endswith(".0"):
        s = s[:-2]
    return s


def _scalar_string(fd, value) -> str:
    """Format a scalar for a path, a query value or a JSON literal."""
    if fd.enum_type is not None:
        try:
            return fd.enum_type.values_by_number[value].name
        except KeyError:
            return str(value)
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, int):
        return str(value)
    if isinstance(value, float):
        return _g_format(value)
    if isinstance(value, str):
        return value
    if isinstance(value, (bytes, bytearray)):
        return base64.b64encode(bytes(value)).decode("ascii")
    raise ValueError(f"connection: field {fd.name}: a {fd.type} cannot be a path or query value")


def _encode_message(msg, graphql: bool) -> dict:
    out: dict[str, Any] = {}
    for fd in msg.DESCRIPTOR.fields:
        if not _is_set(msg, fd):
            continue
        name = fd.json_name if graphql else _wire_name(fd)
        out[name] = _encode_field(msg, fd, graphql)
    return out


def _encode_field(msg, fd, graphql: bool) -> Any:
    if fd.is_repeated:
        if _is_map(fd):
            value_fd = fd.message_type.fields_by_name["value"]
            return {
                str(k): _encode_single(value_fd, v, graphql)
                for k, v in getattr(msg, fd.name).items()
            }
        return [_encode_single(fd, v, graphql) for v in getattr(msg, fd.name)]
    return _encode_single(fd, getattr(msg, fd.name), graphql)


def _encode_single(fd, v, graphql: bool) -> Any:
    if fd.message_type is not None:
        if fd.message_type.full_name.startswith("google.protobuf."):
            return json_format.MessageToDict(v)
        return _encode_message(v, graphql)
    if fd.enum_type is not None:
        try:
            return fd.enum_type.values_by_number[v].name
        except KeyError:
            return str(v)
    if isinstance(v, (bytes, bytearray)):
        return base64.b64encode(bytes(v)).decode("ascii")
    if isinstance(v, float) and (math.isnan(v) or math.isinf(v)):
        raise ValueError(f"connection: field {fd.name}: {v!r} is not valid JSON")
    return v


def _encode_request(spec: Mapping[str, Any], req):
    """Render the path, query and canonical JSON body of one request."""
    fields = req.DESCRIPTOR.fields_by_name
    graphql = bool(spec["graphql_document"])
    path_values: dict[str, str] = {}
    multi: list[str] = []
    query: dict[str, list[str]] = {}
    body: dict[str, Any] = {}
    has_body = graphql
    for p in spec["params"]:
        fd = fields.get(p["field"])
        if fd is None:
            raise ValueError(f"connection: spec names unknown request field {p['field']!r}")
        where = p["in"]
        if where == "path":
            if fd.is_repeated:
                raise ValueError(f"connection: field {fd.name}: a list cannot be a path value")
            path_values[p["wire"]] = _scalar_string(fd, getattr(req, fd.name))
            if p["multi_segment"]:
                multi.append(p["wire"])
        elif where == "query":
            if not _is_set(req, fd):
                continue
            if _is_map(fd):
                raise ValueError(f"connection: field {fd.name}: a map cannot be a query value")
            if fd.is_repeated:
                query.setdefault(p["wire"], []).extend(
                    _scalar_string(fd, v) for v in getattr(req, fd.name)
                )
            else:
                query[p["wire"]] = [_scalar_string(fd, getattr(req, fd.name))]
        elif where == "body":
            has_body = True
            if not _is_set(req, fd):
                continue
            body[p["wire"]] = _encode_field(req, fd, graphql)
        else:
            raise ValueError(f"connection: field {fd.name}: unsupported placement {where!r}")
    path = render_path(spec["path_template"], path_values, multi)
    if not has_body:
        return path, query, None
    if graphql:
        payload: Any = {"query": spec["graphql_document"], "variables": body}
    else:
        payload = body
    return path, query, _dumps(payload).encode("utf-8")


# --- Response decoding ------------------------------------------------------


def _decode_error(spec) -> ValueError:
    return ValueError(
        f"connection: {spec['op'].id}: response does not match the contract"
    )


def _sole_repeated_field(md, next_cursor: str):
    """md's only list field, not counting the next-cursor field, or None."""
    found = None
    for fd in md.fields:
        if not fd.is_repeated or _is_map(fd) or fd.name == next_cursor:
            continue
        if found is not None:
            return None
        found = fd
    return found


def _decode_json(spec, body: bytes, resp) -> None:
    body = body.strip()
    if not body:
        return
    try:
        value = json.loads(body)
    except ValueError:
        raise _decode_error(spec) from None
    _decode_value(spec, value, resp)


def _decode_value(spec, value, resp) -> None:
    if isinstance(value, list):
        # A JSON array fills the response's only repeated field.
        fd = _sole_repeated_field(resp.DESCRIPTOR, spec["next_cursor_field"])
        if fd is None:
            raise _decode_error(spec)
        value = {fd.json_name: value}
    try:
        json_format.ParseDict(value, resp, ignore_unknown_fields=True)
    except json_format.ParseError:
        raise _decode_error(spec) from None


def _graphql_kind(code: str) -> "tuple[str, bool]":
    if code == "NOT_FOUND":
        return NOT_FOUND, False
    if code == "UNAUTHENTICATED":
        return UNAUTHORIZED, False
    if code == "FORBIDDEN":
        return FORBIDDEN, False
    if code == "RATELIMITED":
        return RATE_LIMITED, True
    return "invalid", False


def _decode_graphql(spec, status: int, body: bytes, resp) -> None:
    try:
        env = json.loads(body)
    except ValueError:
        raise _decode_error(spec) from None
    if not isinstance(env, dict):
        raise _decode_error(spec)
    errors = env.get("errors")
    if errors:
        code = ""
        first = errors[0]
        if isinstance(first, dict):
            ext = first.get("extensions")
            if isinstance(ext, dict):
                code = ext.get("code") or ""
        kind, retryable = _graphql_kind(str(code))
        raise OperationError(spec["op"].id, kind, status, retryable)
    data = env.get("data")
    if data is None:
        return
    if isinstance(data, dict) and len(data) == 1:
        key = next(iter(data))
        md = resp.DESCRIPTOR
        if key not in md.fields_by_name and key not in md.fields_by_json_name:
            # GraphQL: data with one key that is not a response field is
            # unwrapped.
            data = data[key]
    _decode_value(spec, data, resp)


def _next_link(headers: Mapping[str, str]) -> str:
    """The rel=\"next\" URL of a Link header, or \"\"."""
    for m in _LINK_RE.finditer(headers.get("link", "")):
        for param in m.group(2).split(";"):
            k, sep, val = param.strip().partition("=")
            if not sep or k.strip().lower() != "rel":
                continue
            for rel in val.strip().strip('"').split():
                if rel.lower() == "next":
                    return m.group(1)
    return ""


def _fill_next_cursor(spec, headers: Mapping[str, str], resp) -> None:
    """Set an unset next-cursor field from the Link header. A value that
    does not fit the field's type is ignored."""
    cursor_field = spec.get("cursor_field", "")
    next_cursor_field = spec.get("next_cursor_field", "")
    if not cursor_field or not next_cursor_field:
        return
    fd = resp.DESCRIPTOR.fields_by_name.get(next_cursor_field)
    if fd is None or fd.is_repeated or _is_set(resp, fd):
        return
    link = _next_link(headers)
    if not link:
        return
    wire = cursor_field
    for p in spec["params"]:
        if p["field"] == cursor_field:
            wire = p["wire"]
    parsed = urlsplit(link)
    raw = ""
    for k, vs in parse_qs(parsed.query).items():
        if k == wire and vs:
            raw = vs[0]
            break
    if not raw:
        return
    try:
        if fd.type == fd.TYPE_STRING:
            setattr(resp, next_cursor_field, raw)
        elif fd.type in (
            fd.TYPE_INT32,
            fd.TYPE_SINT32,
            fd.TYPE_SFIXED32,
            fd.TYPE_INT64,
            fd.TYPE_SINT64,
            fd.TYPE_SFIXED64,
            fd.TYPE_UINT32,
            fd.TYPE_FIXED32,
            fd.TYPE_UINT64,
            fd.TYPE_FIXED64,
        ):
            setattr(resp, next_cursor_field, int(raw))
    except ValueError:
        pass


# --- The generated-operation entry point ------------------------------------


async def call_http(conn, spec: Mapping[str, Any], req, resp_type):
    """Run one HTTP operation: render the path, place the request fields,
    send the request, map a non-2xx status through spec's error table and
    decode the JSON response into a new resp_type message, discarding
    unknown fields."""
    op_id = spec["op"].id
    if conn is None:
        raise ValueError(f"connection: {op_id}: nil connection")
    try:
        path, query, body = _encode_request(spec, req)
        status, headers, resp_body = await conn.request(spec["method"], path, query, body)
    except OperationError:
        raise
    except Exception as e:
        raise ValueError(f"connection: {op_id}: {e}") from e
    err = map_error(op_id, status, spec["errors"], resp_body)
    if err is not None:
        raise err
    resp = resp_type()
    if spec["graphql_document"]:
        _decode_graphql(spec, status, resp_body, resp)
    else:
        _decode_json(spec, resp_body, resp)
    _fill_next_cursor(spec, headers, resp)
    return resp

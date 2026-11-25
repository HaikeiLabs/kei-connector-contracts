# Provider error mapping. An OperationError carries only the operation id,
# the status and the mapped kind: never the response body, a header or a
# credential. Mirrors gen/go/kei/connection/errors.go.
from typing import Iterable, Mapping, Union

# Error kinds, as they appear on the wire and in generated specs. The values
# are the Go ErrorKind names minus the ERROR_KIND_ prefix, lowercased.
UNSPECIFIED = "unspecified"
INVALID = "invalid"
UNAUTHORIZED = "unauthorized"
FORBIDDEN = "forbidden"
NOT_FOUND = "not_found"
CONFLICT = "conflict"
RATE_LIMITED = "rate_limited"
UNAVAILABLE = "unavailable"
# Raised only by the runtime's response check: a 2xx response missed a
# result_required path, so the contract cannot confirm the write happened.
# Never mapped from a provider status.
CONTRACT_VIOLATION = "contract_violation"


class OperationError(Exception):
    """A provider error returned by a generated operation."""

    def __init__(self, op_id: str, kind: str, status: int, retryable: bool = False):
        self.op_id = op_id
        self.kind = kind
        self.status = status
        self.retryable = retryable
        retry = ", retryable" if retryable else ""
        # The message holds only the op id, kind and status: never a body,
        # message or credential.
        super().__init__(f"connection: {op_id}: {kind} (status {status}{retry})")

    def __repr__(self) -> str:
        return str(self)


ErrorTable = Iterable[Mapping[str, Union[int, str, bool]]]


def _default_kind(status: int) -> tuple[str, bool]:
    if status in (400, 422):
        return INVALID, False
    if status == 401:
        return UNAUTHORIZED, False
    if status == 403:
        return FORBIDDEN, False
    if status == 404:
        return NOT_FOUND, False
    if status == 409:
        return CONFLICT, False
    if status == 429:
        return RATE_LIMITED, True
    if status in (502, 503, 504):
        return UNAVAILABLE, True
    if 500 <= status < 600:
        return UNAVAILABLE, False
    return UNSPECIFIED, False


def map_error(
    op_id: str,
    status: int,
    table: ErrorTable,
    body: Union[bytes, str] = b"",
) -> "OperationError | None":
    """Map a provider HTTP status to an OperationError.

    A 2xx status returns None. Otherwise the operation's declared table
    decides; a status the table does not list falls back to the defaults:
    400 and 422 invalid, 401 unauthorized, 403 forbidden, 404 not found,
    409 conflict, 429 rate limited (retryable), 5xx unavailable (502, 503
    and 504 retryable), anything else unspecified.

    body is the response body, accepted only so callers can hand over
    everything they have: it is deliberately never read, so no credential
    or body text can reach the error.
    """
    _ = body
    if 200 <= status < 300:
        return None
    for m in table:
        if m["status"] == status:
            return OperationError(op_id, str(m["kind"]), status, bool(m.get("retryable", False)))
    kind, retryable = _default_kind(status)
    return OperationError(op_id, kind, status, retryable)

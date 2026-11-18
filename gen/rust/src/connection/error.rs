// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//! Provider error mapping, mirroring the Go runtime's `errors.go` exactly.
//! Errors carry only the operation id, the status and the mapped kind:
//! never the response body, a header or a credential.

use std::fmt;

use super::Op;
use crate::kei::contracts::v1::{ErrorKind, ErrorMapping};

/// OperationError is a provider or runtime error returned by a generated
/// operation. With no `detail` it describes a mapped provider error; with a
/// `detail` it wraps a lower-level failure (path, transport, decode). A
/// detail never contains a response body, a header value or a credential.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct OperationError {
    pub kind: ErrorKind,
    pub status: u16,
    pub retryable: bool,
    pub op_id: String,
    pub detail: Option<String>,
}

impl OperationError {
    /// wrap builds an error for a lower-level failure under one operation,
    /// like Go's `fmt.Errorf("connection: %s: %w", op, err)`.
    pub fn wrap(op_id: &str, detail: impl Into<String>) -> Self {
        OperationError {
            kind: ErrorKind::Unspecified,
            status: 0,
            retryable: false,
            op_id: op_id.to_string(),
            detail: Some(detail.into()),
        }
    }
}

impl fmt::Display for OperationError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "connection: {}", self.op_id)?;
        match &self.detail {
            Some(detail) => write!(f, ": {detail}"),
            None => write!(
                f,
                ": {} (status {}{})",
                kind_name(self.kind),
                self.status,
                if self.retryable { ", retryable" } else { "" },
            ),
        }
    }
}

impl std::error::Error for OperationError {}

/// MapError maps a provider HTTP status to an error. A 2xx status returns
/// None. Otherwise the operation's declared table decides; a status the
/// table does not list falls back to the defaults: 400 and 422 invalid, 401
/// unauthorized, 403 forbidden, 404 not found, 409 conflict, 429 rate
/// limited (retryable), 5xx unavailable (502, 503 and 504 retryable),
/// anything else unspecified.
pub fn map_error(op: &Op, status: u16, table: &[ErrorMapping]) -> Option<OperationError> {
    if (200..300).contains(&status) {
        return None;
    }
    for m in table {
        if m.status == status as i32 {
            return Some(OperationError {
                kind: ErrorKind::try_from(m.kind).unwrap_or(ErrorKind::Unspecified),
                status,
                retryable: m.retryable,
                op_id: op.id.clone(),
                detail: None,
            });
        }
    }
    let (kind, retryable) = default_kind(status);
    Some(OperationError {
        kind,
        status,
        retryable,
        op_id: op.id.clone(),
        detail: None,
    })
}

fn default_kind(status: u16) -> (ErrorKind, bool) {
    match status {
        400 | 422 => (ErrorKind::Invalid, false),
        401 => (ErrorKind::Unauthorized, false),
        403 => (ErrorKind::Forbidden, false),
        404 => (ErrorKind::NotFound, false),
        409 => (ErrorKind::Conflict, false),
        429 => (ErrorKind::RateLimited, true),
        502 | 503 | 504 => (ErrorKind::Unavailable, true),
        500..=599 => (ErrorKind::Unavailable, false),
        _ => (ErrorKind::Unspecified, false),
    }
}

/// kindName turns ERROR_KIND_NOT_FOUND into not_found.
fn kind_name(kind: ErrorKind) -> String {
    kind.as_str_name()
        .trim_start_matches("ERROR_KIND_")
        .to_ascii_lowercase()
}

/// graphQLKind maps a GraphQL errors[0].extensions.code to a kind:
/// NOT_FOUND, UNAUTHENTICATED (as unauthorized), FORBIDDEN, RATELIMITED (as
/// rate_limited, retryable); anything else is invalid.
pub(crate) fn graph_ql_kind(code: &str) -> (ErrorKind, bool) {
    match code {
        "NOT_FOUND" => (ErrorKind::NotFound, false),
        "UNAUTHENTICATED" => (ErrorKind::Unauthorized, false),
        "FORBIDDEN" => (ErrorKind::Forbidden, false),
        "RATELIMITED" => (ErrorKind::RateLimited, true),
        _ => (ErrorKind::Invalid, false),
    }
}

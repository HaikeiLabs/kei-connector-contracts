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

//! Package connection is the hand-written runtime that generated connector
//! operations run on (proposal 0002, section 5.2).
//!
//! It supplies the authenticated connection type (HTTP), the operation
//! descriptor, the path renderer, canonical request/response JSON helpers
//! and the error mapping. It takes credentials as arguments and never
//! fetches them: it has no kei-proxy client, no secret-store client, no
//! environment reads and no logging, and it reaches the network only through
//! the `reqwest::Client` the caller passes in. No error or debug value it
//! produces contains a credential or a provider response body.
//!
//! SQL connectors have no database-driver dependency: SQL operations run
//! through the caller-supplied [`Sql`] transport, like HTTP operations run
//! through the [`Http`] transport over `reqwest`.

pub mod error;
pub mod http;
pub mod path;

pub use error::{map_error, OperationError};
pub use http::{
    ConnectionError, DoResponse, HttpConnection, DEFAULT_HEADER_NAME, MAX_RESPONSE_BYTES,
};
pub use path::{render_path, PathError};

use base64::Engine as _;
use serde::de::DeserializeOwned;
use serde_json::Value;

/// Class is an operation's class. Unspecified is neither Read nor Write, so
/// an unset descriptor is never mistaken for a read.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Class {
    Unspecified,
    Read,
    Write,
}

impl Class {
    pub fn as_str(self) -> &'static str {
        match self {
            Class::Read => "read",
            Class::Write => "write",
            Class::Unspecified => "unspecified",
        }
    }
}

/// Op describes one generated operation, for authorization and audit. ID is
/// the operation id, which is also the capability and the policy target
/// (connector:<ID>).
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Op {
    pub id: &'static str,
    pub class: Class,
    pub destructive: bool,
}

/// Http is the transport generated HTTP operations run on: one authenticated
/// connection, with the credential already attached. Callers pass a
/// rendered, escaped path. A non-2xx status is not an error: map it with
/// `map_error`.
pub trait Http {
    fn do_request(
        &self,
        method: &str,
        path: &str,
        query: &[(String, String)],
        body: Option<&[u8]>,
    ) -> impl std::future::Future<Output = Result<DoResponse, ConnectionError>> + Send;
}

impl Http for HttpConnection {
    fn do_request(
        &self,
        method: &str,
        path: &str,
        query: &[(String, String)],
        body: Option<&[u8]>,
    ) -> impl std::future::Future<Output = Result<DoResponse, ConnectionError>> + Send {
        // The inherent method takes precedence; this does not recurse.
        self.do_request(method, path, query, body)
    }
}

/// encodeBytes base64-encodes a bytes field as proto JSON does (standard
/// alphabet, padded).
pub fn encode_bytes(b: &[u8]) -> String {
    base64::engine::general_purpose::STANDARD.encode(b)
}

/// SqlColumn is one column of a SQL result set.
#[derive(Debug, Clone)]
pub struct SqlColumn {
    pub name: String,
    pub ty: String,
}

/// SqlResult is one SQL query result: the columns, the rows (one JSON value
/// per cell) and whether the row cap truncated it.
#[derive(Debug, Clone)]
pub struct SqlResult {
    pub columns: Vec<SqlColumn>,
    pub rows: Vec<Vec<Value>>,
    pub truncated: bool,
}

/// Sql is the caller-supplied SQL transport. SQL connectors have no driver
/// dependency until a Rust consumer needs one.
pub trait Sql {
    fn query(
        &self,
        sql: &str,
        params: &[Value],
        max_rows: i64,
    ) -> impl std::future::Future<Output = Result<SqlResult, ConnectionError>> + Send;
}

/// callSql runs one SQL operation: it converts the request's
/// google.protobuf.Value params to JSON, applies the default row cap
/// (1000), and decodes the result into the response type. Large integers
/// become strings so no id loses precision (see sql_cell).
pub async fn call_sql<T: Sql, R: DeserializeOwned + Default>(
    conn: &T,
    op: &Op,
    sql: &str,
    params: &[pbjson_types::Value],
    max_rows: i64,
) -> Result<R, OperationError> {
    let args: Vec<Value> = params
        .iter()
        .map(|p| serde_json::to_value(p).unwrap_or(Value::Null))
        .collect();
    let max = if max_rows > 0 { max_rows } else { 1000 };
    let result = conn
        .query(sql, &args, max)
        .await
        .map_err(|e| OperationError::wrap(op.id, e.to_string()))?;
    let columns: Vec<Value> = result
        .columns
        .iter()
        .map(|c| serde_json::json!({"name": c.name, "type": c.ty.to_ascii_lowercase()}))
        .collect();
    let rows: Vec<Value> = result
        .rows
        .iter()
        .map(|row| Value::Array(row.iter().cloned().map(sql_cell).collect()))
        .collect();
    let json = serde_json::json!({"columns": columns, "rows": rows, "truncated": result.truncated});
    serde_json::from_value(json).map_err(|_| decode_error(&op.id))
}

/// sqlCell canonicalizes one SQL cell the way the Go runtime does: integers
/// beyond 2^53 become strings so no id loses precision; everything else
/// passes through.
fn sql_cell(v: Value) -> Value {
    match v {
        Value::Number(n) => {
            if let Some(i) = n.as_i64() {
                if i.abs() > 9_007_199_254_740_992 {
                    Value::String(i.to_string())
                } else {
                    Value::Number(n)
                }
            } else if let Some(u) = n.as_u64() {
                if u > 9_007_199_254_740_992 {
                    Value::String(u.to_string())
                } else {
                    Value::Number(n)
                }
            } else {
                Value::Number(n)
            }
        }
        other => other,
    }
}

/// decodeResponse decodes a JSON response into T: unknown fields are
/// discarded (pbjson's ignore_unknown_fields), and a top-level JSON array
/// fills T's only repeated field, named by `array_field`'s JSON name.
/// An empty body yields T::default(). Decoding errors never include the
/// response body.
pub fn decode_response<T: DeserializeOwned + Default>(
    op_id: &str,
    body: &[u8],
    array_field: Option<&str>,
) -> Result<T, OperationError> {
    let body: &[u8] = trim(body);
    if body.is_empty() {
        return Ok(T::default());
    }
    let value: Value = serde_json::from_slice(body).map_err(|_| decode_error(op_id))?;
    decode_value(op_id, value, array_field)
}

/// decodeGraphQLResponse decodes a GraphQL response: an "errors" array maps
/// extensions.code to an error kind (see graph_ql_kind); when data has one
/// key that is not a field of the response (named by `response_fields`, as
/// (proto name, JSON name) pairs), that key's value is decoded; otherwise
/// data itself is.
pub fn decode_graphql_response<T: DeserializeOwned + Default>(
    op_id: &str,
    status: u16,
    body: &[u8],
    response_fields: &[(&str, &str)],
    array_field: Option<&str>,
) -> Result<T, OperationError> {
    let env: Value = serde_json::from_slice(body).map_err(|_| decode_error(op_id))?;
    if let Some(errors) = env.get("errors").and_then(|v| v.as_array()) {
        if !errors.is_empty() {
            let code = errors[0]
                .get("extensions")
                .and_then(|e| e.get("code"))
                .and_then(|c| c.as_str())
                .unwrap_or("");
            let (kind, retryable) = error::graph_ql_kind(code);
            return Err(OperationError {
                kind,
                status,
                retryable,
                op_id: op_id.to_string(),
                detail: None,
            });
        }
    }
    let mut data = env.get("data").cloned().unwrap_or(Value::Null);
    if data.is_null() {
        return Ok(T::default());
    }
    if let Value::Object(obj) = &data {
        if obj.len() == 1 {
            for (k, v) in obj {
                if !response_fields
                    .iter()
                    .any(|(name, json)| *name == k.as_str() || *json == k.as_str())
                {
                    data = v.clone();
                    break;
                }
            }
        }
    }
    decode_value(op_id, data, array_field)
}

fn decode_value<T: DeserializeOwned + Default>(
    op_id: &str,
    value: Value,
    array_field: Option<&str>,
) -> Result<T, OperationError> {
    let value = match value {
        Value::Array(items) => match array_field {
            Some(name) => {
                let mut obj = serde_json::Map::new();
                obj.insert(name.to_string(), Value::Array(items));
                Value::Object(obj)
            }
            None => return Err(decode_error(op_id)),
        },
        v => v,
    };
    serde_json::from_value(value).map_err(|_| decode_error(op_id))
}

fn decode_error(op_id: &str) -> OperationError {
    OperationError::wrap(op_id, "response does not match the contract")
}

/// nextCursorFromLink returns the cursor-field query parameter of the
/// Link header's rel="next" URL, or None.
pub fn next_cursor_from_link(
    headers: &reqwest::header::HeaderMap,
    cursor_wire: &str,
) -> Option<String> {
    let link = next_link(headers)?;
    let url = match url::Url::parse(&link) {
        Ok(u) => u,
        Err(_) => {
            // A relative reference is valid in a Link header; parse it
            // against a placeholder base to read its query.
            let base = url::Url::parse("http://link.invalid/").expect("placeholder base");
            url::Url::options()
                .base_url(Some(&base))
                .parse(&link)
                .ok()?
        }
    };
    url.query_pairs()
        .find(|(k, _)| k == cursor_wire)
        .map(|(_, v)| v.into_owned())
}

/// nextLink returns the rel="next" URL of a Link header, or None.
fn next_link(headers: &reqwest::header::HeaderMap) -> Option<String> {
    for value in headers.get_all(reqwest::header::LINK) {
        let v = value.to_str().ok()?;
        let mut rest = v;
        while let Some(start) = rest.find('<') {
            let end = rest[start..].find('>')? + start;
            let link = rest[start + 1..end].to_string();
            let params = &rest[end + 1..];
            let params = match params.find('<') {
                Some(i) => &params[..i],
                None => params,
            };
            let is_next = params.split(';').any(|p| {
                let p = p.trim();
                match p.split_once('=') {
                    Some((k, val)) => {
                        k.trim().eq_ignore_ascii_case("rel")
                            && val
                                .trim()
                                .trim_matches('"')
                                .split_whitespace()
                                .any(|w| w.eq_ignore_ascii_case("next"))
                    }
                    None => false,
                }
            });
            if is_next {
                return Some(link);
            }
            rest = &rest[end + 1..];
        }
    }
    None
}

fn trim(body: &[u8]) -> &[u8] {
    let start = body.iter().position(|b| !b.is_ascii_whitespace());
    match start {
        Some(s) => {
            let end = body
                .iter()
                .rposition(|b| !b.is_ascii_whitespace())
                .unwrap_or(s);
            &body[s..=end]
        }
        None => &[],
    }
}

/// Render integral protobuf JSON numbers without a fractional suffix, as Go's
/// JSON encoder does. Restrict conversion to exactly representable integers;
/// preserve negative zero and all fractional values.
pub fn normalize_protojson_numbers(value: Value) -> Value {
    match value {
        Value::Array(items) => {
            Value::Array(items.into_iter().map(normalize_protojson_numbers).collect())
        }
        Value::Object(items) => Value::Object(
            items
                .into_iter()
                .map(|(k, v)| (k, normalize_protojson_numbers(v)))
                .collect(),
        ),
        Value::Number(ref number) if number.is_f64() => {
            let n = number.as_f64().expect("JSON number is finite");
            if n.fract() == 0.0
                && n.abs() <= 9_007_199_254_740_992.0
                && !(n == 0.0 && n.is_sign_negative())
            {
                Value::Number((n as i64).into())
            } else {
                value
            }
        }
        _ => value,
    }
}

#[cfg(test)]
mod protojson_number_tests {
    #[test]
    fn normalizes_integral_numbers_inside_structs_and_lists() {
        let input = serde_json::json!({"body": {"version": 1.0, "items": [2.0, 1.5, null, {}]}});
        let got = super::normalize_protojson_numbers(input);
        assert_eq!(
            serde_json::to_string(&got).unwrap(),
            r#"{"body":{"items":[2,1.5,null,{}],"version":1}}"#
        );
    }
}

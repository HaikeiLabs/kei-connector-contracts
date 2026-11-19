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

//! The authenticated HTTP connection, on `reqwest`, mirroring the Go
//! runtime's `http.go`. The credential is held in private fields and never
//! printed: Display and Debug are hand-written.

use std::fmt;

use base64::Engine as _;
use reqwest::header::{HeaderMap, HeaderName, HeaderValue, ACCEPT, AUTHORIZATION, CONTENT_TYPE};
use url::Url;

use crate::kei::contracts::v1::{Connection, ConnectionKind};

/// DefaultHeaderName is the header an HTTP_HEADER connection puts the
/// credential in when the connection declares no header name.
pub const DEFAULT_HEADER_NAME: &str = "X-Api-Key";

/// MaxResponseBytes bounds the response body `do_request` reads.
pub const MAX_RESPONSE_BYTES: usize = 32 << 20;

/// DoResponse is what one request returned. A non-2xx status is not an
/// error: map it with `map_error`.
#[derive(Debug, Clone)]
pub struct DoResponse {
    pub status: u16,
    pub headers: HeaderMap,
    pub body: Vec<u8>,
}

/// ConnectionError is a connection setup or transport error. Its Display
/// carries URLs and method names, never a credential or a response body.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ConnectionError(pub String);

impl fmt::Display for ConnectionError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "connection: {}", self.0)
    }
}

impl std::error::Error for ConnectionError {}

/// HttpConnection is an authenticated HTTP connection to one provider
/// origin. The credential is already attached; callers pass a rendered,
/// escaped path.
pub struct HttpConnection {
    base: Url,
    client: reqwest::Client,
    name: String,
    kind: ConnectionKind,
    header: Option<(HeaderName, String)>,
    basic: Option<(String, String)>,
}

impl HttpConnection {
    /// new builds an HTTP connection for conn's kind:
    ///
    /// - CONNECTION_KIND_HTTP_BEARER: `Authorization: Bearer <credential>`;
    /// - CONNECTION_KIND_HTTP_HEADER: the credential in the connection's
    ///   header_name (DEFAULT_HEADER_NAME when it is empty);
    /// - CONNECTION_KIND_HTTP_BASIC: the credential is "user:password".
    ///
    /// base_url must be an absolute http or https URL with no user info,
    /// query or fragment.
    pub fn new(
        base_url: &str,
        conn: &Connection,
        credential: &str,
        client: reqwest::Client,
    ) -> Result<Self, ConnectionError> {
        let base = parse_base(base_url)?;
        if credential.is_empty() {
            return Err(ConnectionError(format!(
                "NewHTTP {:?}: empty credential",
                conn.name
            )));
        }
        if credential
            .bytes()
            .any(|b| b == b'\r' || b == b'\n' || b == b'\0')
        {
            return Err(ConnectionError(format!(
                "NewHTTP {:?}: credential has a control character",
                conn.name
            )));
        }
        let kind = ConnectionKind::try_from(conn.kind).unwrap_or(ConnectionKind::Unspecified);
        let mut h = HttpConnection {
            base,
            client,
            name: conn.name.clone(),
            kind,
            header: None,
            basic: None,
        };
        match kind {
            ConnectionKind::HttpBearer => {
                h.header = Some((
                    HeaderName::from_static("authorization"),
                    format!("Bearer {credential}"),
                ));
            }
            ConnectionKind::HttpHeader => {
                let name = if conn.header_name.is_empty() {
                    DEFAULT_HEADER_NAME.to_string()
                } else {
                    conn.header_name.clone()
                };
                if !valid_header_name(&name) {
                    return Err(ConnectionError(format!(
                        "NewHTTP {:?}: invalid header name {:?}",
                        conn.name, name
                    )));
                }
                let name = HeaderName::from_bytes(name.as_bytes())
                    .map_err(|e| ConnectionError(format!("NewHTTP {:?}: {e}", conn.name)))?;
                h.header = Some((name, credential.to_string()));
            }
            ConnectionKind::HttpBasic => {
                let (user, pass) = match credential.split_once(':') {
                    Some((u, p)) => (u, p),
                    None => {
                        return Err(ConnectionError(format!(
                            "NewHTTP {:?}: a basic credential must be user:password",
                            conn.name
                        )))
                    }
                };
                if user.is_empty() {
                    return Err(ConnectionError(format!(
                        "NewHTTP {:?}: a basic credential must be user:password",
                        conn.name
                    )));
                }
                h.basic = Some((user.to_string(), pass.to_string()));
            }
            _ => {
                return Err(ConnectionError(format!(
                    "NewHTTP {:?}: kind {} is not an HTTP connection",
                    conn.name,
                    kind.as_str_name()
                )))
            }
        }
        Ok(h)
    }

    /// do_request sends one request. path must be an absolute path (from
    /// render_path) with no query or fragment; it is appended to the base
    /// URL's path. The query is a list of (name, value) pairs, encoded like
    /// Go's url.Values.Encode (keys sorted, space as '+').
    pub async fn do_request(
        &self,
        method: &str,
        path: &str,
        query: &[(String, String)],
        body: Option<&[u8]>,
    ) -> Result<DoResponse, ConnectionError> {
        let mut target = self.resolve(path)?;
        let encoded = encode_query(query);
        if !encoded.is_empty() {
            target.set_query(Some(&encoded));
        }
        let method = reqwest::Method::from_bytes(method.as_bytes())
            .map_err(|e| ConnectionError(format!("{method} {path}: {e}")))?;
        let mut req = self.client.request(method.clone(), target.clone());
        req = req.header(ACCEPT, HeaderValue::from_static("application/json"));
        if let Some(b) = body {
            req = req
                .header(CONTENT_TYPE, HeaderValue::from_static("application/json"))
                .body(b.to_vec());
        }
        match &self.header {
            Some((name, value)) => req = req.header(name.clone(), value.clone()),
            None => {
                if let Some((user, pass)) = &self.basic {
                    let value = format!(
                        "Basic {}",
                        base64::engine::general_purpose::STANDARD.encode(format!("{user}:{pass}"))
                    );
                    req = req.header(AUTHORIZATION, value);
                }
            }
        }
        let resp = req
            .send()
            .await
            .map_err(|e| ConnectionError(format!("{method} {path}: {e}")))?;
        let status = resp.status().as_u16();
        let headers = resp.headers().clone();
        let mut resp = resp;
        let mut body_out = Vec::new();
        loop {
            let chunk = resp
                .chunk()
                .await
                .map_err(|e| ConnectionError(format!("{method} {path}: reading response: {e}")))?;
            let Some(chunk) = chunk else { break };
            if body_out.len() + chunk.len() > MAX_RESPONSE_BYTES {
                return Err(ConnectionError(format!(
                    "{method} {path}: response larger than {MAX_RESPONSE_BYTES} bytes",
                )));
            }
            body_out.extend_from_slice(&chunk);
        }
        Ok(DoResponse {
            status,
            headers,
            body: body_out,
        })
    }

    /// resolve appends an escaped absolute path to the base URL's path and
    /// refuses anything that could change the origin.
    fn resolve(&self, path: &str) -> Result<Url, ConnectionError> {
        if !path.starts_with('/')
            || path.starts_with("//")
            || path.contains('?')
            || path.contains('#')
        {
            return Err(ConnectionError(format!(
                "path {path:?} must be an absolute path with no query or fragment",
            )));
        }
        let base_path = self.base.path().trim_end_matches('/');
        let raw = format!(
            "{}{}{}",
            self.base.origin().ascii_serialization(),
            base_path,
            path
        );
        let u =
            Url::parse(&raw).map_err(|e| ConnectionError(format!("invalid path {path:?}: {e}")))?;
        if u.scheme() != self.base.scheme()
            || u.host_str() != self.base.host_str()
            || u.port() != self.base.port()
        {
            return Err(ConnectionError(
                format!("path {path:?} changes the origin",),
            ));
        }
        Ok(u)
    }
}

impl fmt::Display for HttpConnection {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "connection.HttpConnection{{name={} kind={} base={}}}",
            self.name,
            self.kind.as_str_name(),
            self.base.as_str(),
        )
    }
}

// Debug is hand-written: the derived form would print the credential.
impl fmt::Debug for HttpConnection {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("HttpConnection")
            .field("name", &self.name)
            .field("kind", &self.kind.as_str_name())
            .field("base", &self.base.as_str())
            .finish()
    }
}

/// validHeaderName reports whether name is an RFC 9110 token.
fn valid_header_name(name: &str) -> bool {
    !name.is_empty()
        && name
            .bytes()
            .all(|b| b.is_ascii_alphanumeric() || b"!#$%&'*+-.^_`|~".contains(&b))
}

fn parse_base(raw: &str) -> Result<Url, ConnectionError> {
    let u = Url::parse(raw).map_err(|e| ConnectionError(format!("invalid base URL: {e}")))?;
    match () {
        _ if u.scheme() != "http" && u.scheme() != "https" => Err(ConnectionError(format!(
            "base URL {:?} must be http or https",
            redacted(&u)
        ))),
        _ if u.host_str().is_none() => Err(ConnectionError(format!(
            "base URL {:?} has no host",
            redacted(&u)
        ))),
        _ if !u.username().is_empty() || u.password().is_some() => Err(ConnectionError(
            "base URL must not carry user info".to_string(),
        )),
        _ if u.query().is_some() || u.fragment().is_some() => Err(ConnectionError(format!(
            "base URL {:?} must not have a query or fragment",
            redacted(&u)
        ))),
        _ => Ok(u),
    }
}

fn redacted(u: &Url) -> String {
    // The base URL never carries user info (checked above), so it is safe to
    // print as is.
    u.as_str().to_string()
}

/// encodeQuery encodes query pairs the way Go's url.Values.Encode does:
/// keys sorted, each key's values in insertion order, percent-encoded with
/// Go's QueryEscape set (space becomes '+').
pub(crate) fn encode_query(query: &[(String, String)]) -> String {
    let mut by_key: std::collections::BTreeMap<&str, Vec<&str>> = Default::default();
    for (k, v) in query {
        by_key.entry(k.as_str()).or_default().push(v.as_str());
    }
    let mut out = String::new();
    for (k, vals) in by_key {
        for v in vals {
            if !out.is_empty() {
                out.push('&');
            }
            out.push_str(&query_escape(k));
            out.push('=');
            out.push_str(&query_escape(v));
        }
    }
    out
}

/// queryEscape percent-encodes s with exactly Go's `url.QueryEscape` set:
/// A-Za-z0-9 and `-_.~` stay, space becomes '+', every other byte becomes an
/// uppercase %XX escape.
fn query_escape(s: &str) -> String {
    let mut out = String::with_capacity(s.len());
    for b in s.as_bytes() {
        match b {
            b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'-' | b'_' | b'.' | b'~' => {
                out.push(*b as char)
            }
            b' ' => out.push('+'),
            _ => {
                out.push('%');
                out.push_str(&format!("{b:02X}"));
            }
        }
    }
    out
}

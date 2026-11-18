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

//! Path-template rendering, mirroring the Go runtime's `path.go` exactly:
//! every value is one segment, escaped like Go's `url.PathEscape`, and no
//! value can change the shape of the path.

use std::fmt;

/// PathError is a path-template or path-value error. Its Display carries the
/// template or variable name, never a value (a value may be a credential).
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PathError(pub String);

impl fmt::Display for PathError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "connection: {}", self.0)
    }
}

impl std::error::Error for PathError {}

/// RenderPath fills a path template such as "/repos/{owner}/{repository}"
/// from `values`.
///
/// Every value is percent-encoded as a single path segment, so no value can
/// change the shape of the path: "/" becomes %2F, "%" becomes %25, "?" and
/// "#" are encoded, and a value made only of dots ("." or "..") is encoded
/// as %2E so it is never read as a relative segment. A missing or empty
/// value, an unterminated or empty variable, or a template that does not
/// start with "/" is an error.
///
/// A variable named in `multi_segment` (a param with multi_segment, such as
/// a file path or a Git ref) keeps "/" between its segments and escapes each
/// segment; an empty, "." or ".." segment is an error.
pub fn render_path(
    template: &str,
    values: &[(&str, &str)],
    multi_segment: &[&str],
) -> Result<String, PathError> {
    if !template.starts_with('/') {
        return Err(PathError(format!(
            "path template {template:?} must start with /"
        )));
    }
    let mut out = String::with_capacity(template.len());
    let mut rest = template;
    loop {
        let open = match rest.find('{') {
            Some(i) => i,
            None => {
                if rest.contains('}') {
                    return Err(PathError(format!(
                        "path template {template:?} has an unmatched }}",
                    )));
                }
                out.push_str(rest);
                return Ok(out);
            }
        };
        if let Some(close) = rest.find('}') {
            if close < open {
                return Err(PathError(format!(
                    "path template {template:?} has an unmatched }}",
                )));
            }
        }
        let end = rest[open..].find('}').ok_or_else(|| {
            PathError(format!(
                "path template {template:?} has an unterminated variable",
            ))
        })?;
        let name = &rest[open + 1..open + end];
        if name.is_empty() || name.contains('{') || name.contains('/') {
            return Err(PathError(format!(
                "path template {template:?} has an invalid variable",
            )));
        }
        let v = values
            .iter()
            .find(|(n, _)| *n == name)
            .map(|(_, v)| *v)
            .ok_or_else(|| PathError(format!("path variable {name:?} has no value")))?;
        if v.is_empty() {
            return Err(PathError(format!("path variable {name:?} is empty")));
        }
        out.push_str(&rest[..open]);
        if multi_segment.iter().any(|m| *m == name) {
            let mut segs = Vec::new();
            for s in v.split('/') {
                if s.is_empty() || s == "." || s == ".." {
                    return Err(PathError(format!(
                        "path variable {name:?} has an empty, . or .. segment",
                    )));
                }
                segs.push(escape_segment(s));
            }
            out.push_str(&segs.join("/"));
        } else {
            out.push_str(&escape_segment(v));
        }
        rest = &rest[open + end + 1..];
    }
}

/// escapeSegment percent-encodes v as one path segment.
fn escape_segment(v: &str) -> String {
    if !v.is_empty() && v.bytes().all(|b| b == b'.') {
        return "%2E".repeat(v.len());
    }
    path_escape(v)
}

/// pathEscape percent-encodes s with exactly Go's `url.PathEscape` set:
/// A-Za-z0-9 and `-_.~$&+:=@` stay; every other byte (UTF-8), including
/// `!'()*,;`, becomes an uppercase %XX escape.
fn path_escape(s: &str) -> String {
    let mut out = String::with_capacity(s.len());
    for b in s.as_bytes() {
        match b {
            b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'-' | b'_' | b'.' | b'~' | b'$'
            | b'&' | b'+' | b'=' | b':' | b'@' => out.push(*b as char),
            _ => {
                out.push('%');
                out.push_str(&format!("{b:02X}"));
            }
        }
    }
    out
}

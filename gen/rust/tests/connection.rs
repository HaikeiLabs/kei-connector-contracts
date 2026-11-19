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

//! Runtime tests for the hand-written connection module, mirroring the Go
//! runtime's `TestRenderPath_EscapesSegments`, `TestRenderPath_MultiSegment`
//! and `TestConnection_NoCredentialInErrors`.

use haikeilabs_connector_contracts::connection::{
    map_error, render_path, Class, HttpConnection, Op, Transport,
};
use haikeilabs_connector_contracts::kei::contracts::v1::{Connection, ConnectionKind};
use wiremock::matchers::{any, header};
use wiremock::{Mock, MockServer, ResponseTemplate};

/// A synthetic credential; tests assert it never leaks.
const SYNTHETIC_TOKEN: &str = "tok_SECRET";

#[test]
fn render_path_escapes_segments() {
    // The orchestrator's case: reserved characters stay inside one segment.
    let got = render_path(
        "/repos/{owner}/{repository}",
        &[("owner", "a/b"), ("repository", "x y")],
        &[],
    )
    .expect("render_path");
    assert_eq!(got, "/repos/a%2Fb/x%20y");

    // Go's url.PathEscape is the reference: keep $ & + : = @, escape ! ' ( ) * , ;
    let got =
        render_path("/x/{v}", &[("v", "a!b'c(d)e*f,g;h$i&j+k:l=m@n")], &[]).expect("render_path");
    assert_eq!(got, "/x/a%21b%27c%28d%29e%2Af%2Cg%3Bh$i&j+k:l=m@n");
    let got = render_path("/x/{v}", &[("v", "d/a!b")], &["v"]).expect("render_path");
    assert_eq!(got, "/x/d/a%21b");

    // Review Focus 1: %, dot-only values, query characters and a value that
    // ends inside a larger segment all keep the path shape.
    for (template, values, want) in [
        (
            "/repos/{owner}/{repository}",
            vec![("owner", "100%"), ("repository", "w")],
            "/repos/100%25/w",
        ),
        (
            "/repos/{owner}/{repository}",
            vec![("owner", ".."), ("repository", ".")],
            "/repos/%2E%2E/%2E",
        ),
        (
            "/repos/{owner}/{repository}",
            vec![("owner", "a?b"), ("repository", "c#d")],
            "/repos/a%3Fb/c%23d",
        ),
        (
            "/repos/{owner}/{repository}/issues/{issue_number}",
            vec![
                ("owner", "acme"),
                ("repository", "widgets"),
                ("issue_number", "1347"),
            ],
            "/repos/acme/widgets/issues/1347",
        ),
    ] {
        let got = render_path(template, &values, &[])
            .unwrap_or_else(|e| panic!("render_path {values:?}: {e}"));
        assert_eq!(got, want, "values {values:?}");
    }
    let got = render_path("/files/{name}.json", &[("name", "a/b")], &[]).expect("render_path");
    assert_eq!(got, "/files/a%2Fb.json");
}

#[test]
fn render_path_multi_segment() {
    const TMPL: &str = "/repos/{owner}/{repository}/contents/{path}";
    // A multi-segment variable keeps "/" between segments but escapes each
    // segment; other variables stay fully escaped.
    let got = render_path(
        TMPL,
        &[
            ("owner", "a/b"),
            ("repository", "w"),
            ("path", "docs/notes.txt"),
        ],
        &["path"],
    )
    .expect("render_path");
    assert_eq!(got, "/repos/a%2Fb/w/contents/docs/notes.txt");

    let got = render_path(
        TMPL,
        &[
            ("owner", "a/b"),
            ("repository", "w"),
            ("path", "a b/c%d/e?f#g"),
        ],
        &["path"],
    )
    .expect("render_path");
    assert_eq!(got, "/repos/a%2Fb/w/contents/a%20b/c%25d/e%3Ff%23g");

    // Empty, "." and ".." segments are errors.
    for bad in [
        "a/../b", "a//b", "../x", "a/..", "./a", "a/./b", "/a", "a/", "", "..",
    ] {
        assert!(
            render_path(
                TMPL,
                &[("owner", "o"), ("repository", "r"), ("path", bad)],
                &["path"]
            )
            .is_err(),
            "RenderPath({bad:?}) should be an error"
        );
    }

    // A variable not named multi-segment still escapes "/".
    let got = render_path(
        TMPL,
        &[
            ("owner", "o"),
            ("repository", "r"),
            ("path", "docs/notes.txt"),
        ],
        &[],
    )
    .expect("render_path");
    assert_eq!(got, "/repos/o/r/contents/docs%2Fnotes.txt");
}

#[tokio::test]
async fn connection_sends_auth_headers() {
    // Each kind sends its credential in the header the contract names; a
    // request without the expected header matches no mock and gets a 404.
    let srv = MockServer::start().await;
    Mock::given(header("authorization", "Bearer tok_SECRET"))
        .respond_with(ResponseTemplate::new(200))
        .mount(&srv)
        .await;
    Mock::given(header("X-Api-Key", "tok_SECRET"))
        .respond_with(ResponseTemplate::new(200))
        .mount(&srv)
        .await;
    Mock::given(header("authorization", "Basic dXNlcjpwYXNz"))
        .respond_with(ResponseTemplate::new(200))
        .mount(&srv)
        .await;

    for (kind, credential, header_name) in [
        (ConnectionKind::HttpBearer, SYNTHETIC_TOKEN, ""),
        (ConnectionKind::HttpHeader, SYNTHETIC_TOKEN, "X-Api-Key"),
        (ConnectionKind::HttpBasic, "user:pass", ""),
    ] {
        let conn = HttpConnection::new(
            &srv.uri(),
            &Connection {
                name: "pat".to_string(),
                kind: kind as i32,
                header_name: header_name.to_string(),
                ..Default::default()
            },
            credential,
            reqwest::Client::new(),
        )
        .expect("HttpConnection::new");
        let resp = conn
            .do_request("GET", "/x", &[], None)
            .await
            .expect("do_request");
        assert_eq!(
            resp.status, 200,
            "kind {:?} did not send its auth header",
            kind
        );
    }
    // Every mock was hit exactly once: the right header went out per kind.
    srv.verify().await;
}

#[tokio::test]
async fn no_credential_in_errors() {
    // A hostile or careless provider echoes the credential back.
    let srv = MockServer::start().await;
    srv.register(Mock::given(any()).respond_with(
        ResponseTemplate::new(401).set_body_string(format!("bad token {SYNTHETIC_TOKEN}")),
    ))
    .await;

    let conn = HttpConnection::new(
        &srv.uri(),
        &Connection {
            name: "pat".to_string(),
            kind: ConnectionKind::HttpBearer as i32,
            ..Default::default()
        },
        SYNTHETIC_TOKEN,
        reqwest::Client::new(),
    )
    .expect("HttpConnection::new");

    let resp = conn
        .do_request("GET", "/x", &[], None)
        .await
        .expect("do_request");
    assert_eq!(resp.status, 401);
    let err = map_error(
        &Op {
            id: "github.issues.get",
            class: Class::Read,
            destructive: false,
            resource_template: "repos/{owner}/{repository}",
            method_name: "GetIssue",
            connections: &["user_oauth", "pat"],
            transport: Transport::Http,
        },
        resp.status,
        &[],
    )
    .expect("a 401 maps to an error");
    let shown = format!("{err} {err:?}");
    assert!(
        !shown.contains(SYNTHETIC_TOKEN),
        "credential leaked: {shown}"
    );
    assert!(
        !shown.contains("bad token"),
        "response body leaked: {shown}"
    );

    // The connection's own Display/Debug never carries the credential either.
    let shown = format!("{conn} {conn:?}");
    assert!(
        !shown.contains(SYNTHETIC_TOKEN),
        "credential leaked: {shown}"
    );
}

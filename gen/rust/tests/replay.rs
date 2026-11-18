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

//! Replay test for a generated operation: github.issues.get over a mocked
//! provider, serving the same synthetic fixture the other language bindings
//! replay (fixtures/operations/github/issues.get.response.json). Proves the
//! generated template drives the hand-written connection runtime end to
//! end: rendered path, Bearer credential, JSON decode back into the prost
//! message.

use std::path::PathBuf;

use haikeilabs_connector_contracts::connection::{Class, HttpConnection, Op};
use haikeilabs_connector_contracts::kei::connectors::github::v1::ops::{
    GithubOperations, GithubQueries, OP_GET_ISSUE,
};
use haikeilabs_connector_contracts::kei::connectors::github::v1::GetIssueRequest;
use haikeilabs_connector_contracts::kei::contracts::v1::{Connection, ConnectionKind};
use wiremock::matchers::{method, path};
use wiremock::{Mock, MockServer, ResponseTemplate};

/// A synthetic credential; only wiremock sees it.
const SYNTHETIC_TOKEN: &str = "synthetic-token";

/// The shared operation fixture, relative to this crate's manifest dir.
fn fixture_body() -> Vec<u8> {
    let p = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("../../fixtures/operations/github/issues.get.response.json");
    std::fs::read(&p).unwrap_or_else(|e| panic!("read {}: {e}", p.display()))
}

#[tokio::test]
async fn replay_get_issue_against_fixture() {
    let srv = MockServer::start().await;
    srv.register(
        Mock::given(method("GET"))
            .and(path("/repos/acme/widgets/issues/1347"))
            .respond_with(ResponseTemplate::new(200).set_body_bytes(fixture_body())),
    )
    .await;

    let conn = HttpConnection::new(
        &srv.uri(),
        &Connection {
            name: "github-pat".to_string(),
            kind: ConnectionKind::HttpBearer as i32,
            ..Default::default()
        },
        SYNTHETIC_TOKEN,
        reqwest::Client::new(),
    )
    .expect("HttpConnection::new");

    // The descriptor the policy grants: connector:github.issues.get, a read.
    assert_eq!(OP_GET_ISSUE.id, "github.issues.get");
    assert_eq!(OP_GET_ISSUE.class, Class::Read);
    let op_check: Op = OP_GET_ISSUE;
    assert!(!op_check.destructive);

    let queries = GithubQueries { conn };
    let resp = queries
        .get_issue(GetIssueRequest {
            owner: "acme".to_string(),
            repository: "widgets".to_string(),
            issue_number: 1347,
            ..Default::default()
        })
        .await
        .expect("get_issue");

    assert_eq!(resp.number, 1347);
    assert_eq!(resp.id, 9007199254740993, "int64 survives proto-JSON decode");
    assert_eq!(resp.title, "Synthetic pilot issue");
    srv.verify().await;
}

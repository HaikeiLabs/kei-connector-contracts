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

//! Result-required checks: a write operation whose result_required paths
//! are not all held by a 2xx response fails with a contract violation, not
//! with a silent success. Mirrors gen/go/conformance/result_required_test.go.

use std::path::PathBuf;

use haikeilabs_connector_contracts::connection::{HttpConnection, OperationError};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::ops::{LinearOperations, LinearQueries};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::CreateIssueRequest;
use haikeilabs_connector_contracts::kei::contracts::v1::{Connection, ConnectionKind, ErrorKind};
use wiremock::matchers::{method, path};
use wiremock::{Mock, MockServer, ResponseTemplate};

/// A synthetic credential; only wiremock sees it.
const SYNTHETIC_TOKEN: &str = "synthetic-token";

fn fixture(name: &str) -> Vec<u8> {
    let p = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join(format!("../../fixtures/operations/linear/{name}.json"));
    std::fs::read(&p).unwrap_or_else(|e| panic!("read {}: {e}", p.display()))
}

fn request() -> CreateIssueRequest {
    CreateIssueRequest {
        title: "Synthetic pilot issue".to_string(),
        team_id: "00000000-0000-0000-0000-000000000002".to_string(),
        ..Default::default()
    }
}

async fn serve_envelope(fixture_name: &str) -> MockServer {
    let srv = MockServer::start().await;
    let envelope = serde_json::json!({
        "data": { "issueCreate": serde_json::from_slice::<serde_json::Value>(&fixture(fixture_name)).expect("fixture is JSON") }
    });
    srv.register(
        Mock::given(method("POST"))
            .and(path("/graphql"))
            .respond_with(ResponseTemplate::new(200).set_body_json(envelope)),
    )
    .await;
    srv
}

fn queries(srv: &MockServer) -> LinearQueries {
    let conn = HttpConnection::new(
        &srv.uri(),
        &Connection {
            name: "linear-pat".to_string(),
            kind: ConnectionKind::HttpBearer as i32,
            ..Default::default()
        },
        SYNTHETIC_TOKEN,
        reqwest::Client::new(),
    )
    .expect("HttpConnection::new");
    LinearQueries { conn }
}

#[tokio::test]
async fn rejects_a_response_missing_a_result_field() {
    for (name, fixture_name) in [
        ("success false", "issues.create.response.success-false"),
        ("issue missing", "issues.create.response.no-issue"),
        ("empty identifier", "issues.create.response.empty-identifier"),
    ] {
        let srv = serve_envelope(fixture_name).await;
        let err = queries(&srv)
            .create_issue(request())
            .await
            .expect_err(&format!("{name}: the call must fail"));
        let OperationError {
            kind,
            status,
            op_id,
            detail,
            ..
        } = &err;
        assert_eq!(
            *kind,
            ErrorKind::ContractViolation,
            "{name}: kind = {:?}",
            err
        );
        assert_eq!(*status, 200, "{name}: status = {:?}", err);
        assert_eq!(op_id, "linear.issues.create", "{name}: op id = {:?}", err);
        assert!(detail.is_none(), "{name}: detail = {:?}", err);
    }
}

#[tokio::test]
async fn keeps_the_good_fixture_passing() {
    let srv = serve_envelope("issues.create.response").await;
    let resp = queries(&srv).create_issue(request()).await.expect("create_issue");
    assert!(resp.success);
    assert_eq!(resp.issue.as_ref().expect("issue").identifier, "ENG-102");
    srv.verify().await;
}

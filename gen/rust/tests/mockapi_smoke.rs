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

//! Smoke tests: the generated Rust operations run against the real
//! contract-derived mock server (testenv/mockapi) — fixture responses,
//! GraphQL root-field routing, and failure injection via the X-Mock-Fail
//! header mapping to the generated error kinds.
//!
//! The mock is the Go binary from testenv/mockapi/cmd/mockapi, built once
//! per test with the local Go toolchain (or taken from MOCKAPI_BIN) and
//! started on a free port.

use std::io::Read;
use std::net::TcpStream;
use std::path::PathBuf;
use std::process::{Child, Command, Stdio};
use std::time::{Duration, Instant};

use haikeilabs_connector_contracts::connection::HttpConnection;
use haikeilabs_connector_contracts::kei::connectors::github::v1::ops::{GithubOperations, GithubQueries};
use haikeilabs_connector_contracts::kei::connectors::github::v1::GetIssueRequest;
use haikeilabs_connector_contracts::kei::connectors::linear::v1::ops::{LinearOperations, LinearQueries};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::GetIssueRequest as LinearGetIssueRequest;
use haikeilabs_connector_contracts::kei::contracts::v1::{Connection, ConnectionKind, ErrorKind};

/// Repository root, relative to this crate's manifest dir.
fn repo_root() -> PathBuf {
    PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../..")
}

struct Mock {
    child: Child,
    base: String,
}

impl Drop for Mock {
    fn drop(&mut self) {
        let _ = self.child.kill();
        let _ = self.child.wait();
    }
}

/// Start the mock on a free port: build (unless MOCKAPI_BIN is set), spawn,
/// and read the bound address from the server's first stderr line.
fn start_mock() -> Mock {
    let bin = match std::env::var("MOCKAPI_BIN") {
        Ok(bin) => PathBuf::from(bin),
        Err(_) => {
            let out = std::env::temp_dir().join(format!("mockapi-smoke-{}", std::process::id()));
            let status = Command::new("go")
                .args(["-C", "testenv/mockapi", "build", "-o", out.to_str().expect("temp path is utf-8"), "./cmd/mockapi"])
                .current_dir(repo_root())
                .status()
                .expect("go is needed to build the mock (or set MOCKAPI_BIN)");
            assert!(status.success(), "go build of the mock failed");
            out
        }
    };
    let mut child = Command::new(&bin)
        .args([
            "-addr",
            "127.0.0.1:0",
            "-manifest",
            &repo_root().join("manifest/connector-manifest.v1.json").display().to_string(),
            "-fixtures",
            &repo_root().join("fixtures/operations").display().to_string(),
        ])
        .stdout(Stdio::null())
        .stderr(Stdio::piped())
        .spawn()
        .expect("spawn mockapi");
    let mut line = String::new();
    let mut byte = [0u8; 1];
    let stderr = child.stderr.as_mut().expect("piped stderr");
    let deadline = Instant::now() + Duration::from_secs(30);
    loop {
        if Instant::now() > deadline {
            panic!("mockapi did not print a listening address");
        }
        match stderr.read(&mut byte) {
            Ok(0) => panic!("mockapi exited before listening"),
            Ok(_) => {
                if byte[0] == b'\n' {
                    break;
                }
                line.push(byte[0] as char);
            }
            Err(e) => panic!("read mockapi stderr: {e}"),
        }
    }
    let addr = line.trim().rsplit(' ').next().expect("the serving line ends with the address").to_string();
    for _ in 0..100 {
        if TcpStream::connect(&addr).is_ok() {
            return Mock { child, base: format!("http://{addr}") };
        }
        std::thread::sleep(Duration::from_millis(50));
    }
    panic!("mockapi never accepted connections on {addr}");
}

fn connection(base: &str, client: reqwest::Client) -> HttpConnection {
    HttpConnection::new(
        base,
        &Connection {
            name: "smoke".to_string(),
            kind: ConnectionKind::HttpBearer as i32,
            ..Default::default()
        },
        "smoke-token",
        client,
    )
    .expect("HttpConnection::new")
}

#[tokio::test]
async fn github_get_issue_replays_fixture() {
    let mock = start_mock();
    let queries = GithubQueries { conn: connection(&mock.base, reqwest::Client::new()) };
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
    // The fixture id is 2^53+1: it must survive the decode losslessly.
    assert_eq!(resp.id, 9007199254740993);
    assert_eq!(resp.title, "Synthetic pilot issue");
}

#[tokio::test]
async fn linear_get_issue_replays_graphql() {
    let mock = start_mock();
    let queries = LinearQueries { conn: connection(&mock.base, reqwest::Client::new()) };
    let resp = queries
        .get_issue(LinearGetIssueRequest { id: "00000000-0000-0000-0000-000000000101".to_string() })
        .await
        .expect("get_issue");
    // GraphQL decoding unwraps the {"data": {"issue": ...}} envelope.
    let issue = resp.issue.expect("issue");
    assert_eq!(issue.identifier, "ENG-101");
}

#[tokio::test]
async fn injected_401_maps_to_unauthorized() {
    let mock = start_mock();
    let mut headers = reqwest::header::HeaderMap::new();
    headers.insert("x-mock-fail", reqwest::header::HeaderValue::from_static("401"));
    let client = reqwest::Client::builder().default_headers(headers).build().expect("client");
    let queries = GithubQueries { conn: connection(&mock.base, client) };
    let err = queries
        .get_issue(GetIssueRequest {
            owner: "acme".to_string(),
            repository: "widgets".to_string(),
            issue_number: 1347,
            ..Default::default()
        })
        .await
        .expect_err("the injected 401 must surface as an OperationError");
    assert_eq!(err.kind, ErrorKind::Unauthorized);
    assert_eq!(err.status, 401);
    assert_eq!(err.op_id, "github.issues.get");
}

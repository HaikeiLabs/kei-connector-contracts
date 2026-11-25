// Copyright 2026 Haikei Labs
// SPDX-License-Identifier: Apache-2.0
//
// Rust runner of the cross-language parity test (protocol: see
// scripts/parity/run.sh). build.rs generates `dispatch` from the generated
// operations, so a new operation or connector needs no change here.
use haikeilabs_connector_contracts::connection::HttpConnection;
use haikeilabs_connector_contracts::kei::contracts::v1::{Connection, ConnectionKind};
use std::io::Write;

include!(concat!(env!("OUT_DIR"), "/dispatch.rs"));

#[tokio::main(flavor = "current_thread")]
async fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() != 4 {
        eprintln!("usage: kei-parity-runner <cases.json> <base url> <token>");
        std::process::exit(2);
    }
    let cases: Vec<serde_json::Value> = serde_json::from_str(
        &std::fs::read_to_string(&args[1]).expect("reading the cases"),
    )
    .expect("parsing the cases");
    let conn_spec = Connection {
        name: "parity".to_string(),
        kind: ConnectionKind::HttpBearer as i32,
        ..Default::default()
    };
    let client = reqwest::Client::new();
    let mut out = std::io::stdout().lock();
    for case in cases {
        let op = case["op"].as_str().unwrap_or_default().to_string();
        // A "@<variant>" suffix marks a deliberately broken response fixture
        // (parity's result_required cases); the dispatch still calls the
        // real operation.
        let op_id = op.split('@').next().unwrap_or(&op).to_string();
        let result = match HttpConnection::new(&format!("{}/{}", args[2], op), &conn_spec, &args[3], client.clone()) {
            Ok(conn) => dispatch(&op_id, conn, case["request"].clone()).await,
            Err(e) => Err(e.to_string()),
        };
        let line = match result {
            Ok(response) => serde_json::json!({ "op": op, "response": response }),
            Err(error) => serde_json::json!({ "op": op, "error": error }),
        };
        writeln!(out, "{line}").expect("writing a result");
    }
}

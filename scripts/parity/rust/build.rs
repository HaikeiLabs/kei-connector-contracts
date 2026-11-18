// Copyright 2026 Haikei Labs
// SPDX-License-Identifier: Apache-2.0
//
// Generates the runner's dispatch table from the generated Rust operations:
// every gen/rust/src/kei/connectors/<c>/v1/ops.rs declares a trait whose
// methods are documented "/// <fn> runs <op id> (...)" and declared
// "async fn <fn>(&self, req: <Req>) -> Result<<Resp>, OperationError>;".
// A new operation or connector therefore needs no change here. SQL
// connectors (a generic Queries<S>) are skipped: parity covers HTTP only.
use std::{env, fs, path::Path};

fn main() {
    let connectors = Path::new(env!("CARGO_MANIFEST_DIR")).join("../../../gen/rust/src/kei/connectors");
    println!("cargo:rerun-if-changed={}", connectors.display());
    let mut arms = String::new();
    let mut dirs: Vec<_> = fs::read_dir(&connectors).expect("gen/rust connectors").flatten().collect();
    dirs.sort_by_key(|d| d.file_name());
    for dir in dirs {
        let connector = dir.file_name().to_string_lossy().to_string();
        let ops = dir.path().join("v1/ops.rs");
        let Ok(src) = fs::read_to_string(&ops) else { continue };
        println!("cargo:rerun-if-changed={}", ops.display());
        let queries = src.lines().find_map(|l| {
            let l = l.trim();
            let name = l.strip_prefix("pub struct ")?.strip_suffix(" {")?;
            name.ends_with("Queries").then(|| name.to_string())
        });
        let Some(queries) = queries else { continue }; // SQL: Queries<S: Sql>
        let trait_name = src
            .lines()
            .find_map(|l| l.trim().strip_prefix("pub trait ")?.strip_suffix(" {").map(str::to_string))
            .expect("Operations trait");
        let module = format!("haikeilabs_connector_contracts::kei::connectors::{connector}::v1");
        let lines: Vec<&str> = src.lines().map(str::trim).collect();
        for w in lines.windows(2) {
            let (Some(doc), Some(sig)) = (w[0].strip_prefix("/// "), w[1].strip_prefix("async fn ")) else { continue };
            if !sig.ends_with(';') {
                continue; // the impl, not the trait declaration
            }
            let mut doc_words = doc.split_whitespace();
            let (Some(func), Some("runs"), Some(op)) = (doc_words.next(), doc_words.next(), doc_words.next()) else { continue };
            let req = sig.split("req: ").nth(1).and_then(|s| s.split(')').next()).expect("request type");
            let resp = sig.split("Result<").nth(1).and_then(|s| s.split(',').next()).expect("response type");
            arms.push_str(&format!(
                "        {op:?} => {{\n            use {module}::ops::{trait_name} as _;\n            let req: {module}::{req} = serde_json::from_value(request).map_err(|e| format!(\"request fixture: {{e}}\"))?;\n            let resp: {module}::{resp} = {module}::ops::{queries} {{ conn }}.{func}(req).await.map_err(|e| e.to_string())?;\n            serde_json::to_value(&resp).map_err(|e| e.to_string())\n        }}\n"
            ));
        }
    }
    let code = format!(
        "pub async fn dispatch(op: &str, conn: haikeilabs_connector_contracts::connection::HttpConnection, request: serde_json::Value) -> Result<serde_json::Value, String> {{\n    match op {{\n{arms}        other => Err(format!(\"no Rust operation {{other}}\")),\n    }}\n}}\n"
    );
    fs::write(Path::new(&env::var("OUT_DIR").unwrap()).join("dispatch.rs"), code).unwrap();
}

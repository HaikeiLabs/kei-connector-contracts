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

//! HAI-990 gate: every pilot operation response fixture unmarshals into the
//! generated prost types (serde impls from prost-serde with
//! ignore_unknown_fields, so the unknown extra fields each fixture carries
//! are skipped). One spot assert per op proves a field survived with its
//! synthetic value intact, and the >2^53 int64 id round-trips exactly.

use haikeilabs_connector_contracts::kei::connectors::database::v1::{
    RunQueryRequest, RunQueryResponse,
};
use haikeilabs_connector_contracts::kei::connectors::github::v1::{
    GetIssueResponse, GetPullRequestResponse, ListIssuesResponse, PutContentsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::http_api::v1::{
    GetRecordResponse, ListRecordsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::{
    CreateCommentResponse, CreateIssueResponse, GetIssueResponse as LinearGetIssueResponse,
    GetTeamResponse, SearchIssuesResponse,
};
use haikeilabs_connector_contracts::kei::connectors::omni::v1::{
    ListDocumentsResponse, RunQueryResponse as OmniRunQueryResponse,
};
use haikeilabs_connector_contracts::kei::connectors::power_bi::v1::{
    ExecuteQueriesResponse, ListDatasetsResponse,
};
use pbjson_types::value::Kind;

fn fixture(connector: &str, name: &str) -> String {
    let path = format!(
        "{}/../../fixtures/operations/{}/{}",
        env!("CARGO_MANIFEST_DIR"),
        connector,
        name
    );
    std::fs::read_to_string(&path).expect("read fixture")
}

#[test]
fn fixtures_unmarshal() {
    // github.issues.get — spot: number 1347 and the 2^53+1 id, exactly.
    let resp: GetIssueResponse =
        serde_json::from_str(&fixture("github", "issues.get.response.json"))
            .expect("github.issues.get response unmarshals");
    assert_eq!(resp.number, 1347);
    assert_eq!(resp.id, 9007199254740993, "int64 must not go through f64");

    // The >2^53 id round-trips exactly through serialize + reparse.
    let back = serde_json::to_string(&resp).expect("serialize github.issues.get response");
    let reparsed: GetIssueResponse =
        serde_json::from_str(&back).expect("reparse serialized github.issues.get response");
    assert_eq!(reparsed.id, 9007199254740993);
    assert_eq!(reparsed.number, 1347);

    // github.issues.list — spot: two entries.
    let resp: ListIssuesResponse =
        serde_json::from_str(&fixture("github", "issues.list.response.json"))
            .expect("github.issues.list response unmarshals");
    assert_eq!(resp.issues.len(), 2);

    // github.contents.put — spot: content.name.
    let resp: PutContentsResponse =
        serde_json::from_str(&fixture("github", "contents.put.response.json"))
            .expect("github.contents.put response unmarshals");
    assert_eq!(
        resp.content.as_ref().expect("content present").name,
        "notes.txt"
    );

    // github.pulls.get — spot: number 42.
    let resp: GetPullRequestResponse =
        serde_json::from_str(&fixture("github", "pulls.get.response.json"))
            .expect("github.pulls.get response unmarshals");
    assert_eq!(resp.number, 42);

    // linear.issues.get — spot: issue.identifier and labels[0].name.
    let resp: LinearGetIssueResponse =
        serde_json::from_str(&fixture("linear", "issues.get.response.json"))
            .expect("linear.issues.get response unmarshals");
    let issue = resp.issue.as_ref().expect("issue present");
    assert_eq!(issue.identifier, "ENG-101");
    assert_eq!(issue.labels[0].name, "synthetic-label");

    // linear.issues.create — spot: success.
    let resp: CreateIssueResponse =
        serde_json::from_str(&fixture("linear", "issues.create.response.json"))
            .expect("linear.issues.create response unmarshals");
    assert!(resp.success);

    // linear.issues.search — spot: two nodes.
    let resp: SearchIssuesResponse =
        serde_json::from_str(&fixture("linear", "issues.search.response.json"))
            .expect("linear.issues.search response unmarshals");
    assert_eq!(resp.nodes.len(), 2);

    // linear.teams.get — spot: team.key.
    let resp: GetTeamResponse =
        serde_json::from_str(&fixture("linear", "teams.get.response.json"))
            .expect("linear.teams.get response unmarshals");
    assert_eq!(resp.team.as_ref().expect("team present").key, "ENG");

    // linear.comments.create — spot: success and the created comment body.
    let resp: CreateCommentResponse =
        serde_json::from_str(&fixture("linear", "comments.create.response.json"))
            .expect("linear.comments.create response unmarshals");
    assert!(resp.success);
    assert_eq!(
        resp.comment.as_ref().expect("comment present").body,
        "Synthetic comment"
    );

    // http_api.records.list — spot: two Struct records.
    let resp: ListRecordsResponse =
        serde_json::from_str(&fixture("http_api", "records.list.response.json"))
            .expect("http_api.records.list response unmarshals");
    assert_eq!(resp.records.len(), 2);

    // http_api.records.get — spot: record.id inside the Struct body.
    let resp: GetRecordResponse =
        serde_json::from_str(&fixture("http_api", "records.get.response.json"))
            .expect("http_api.records.get response unmarshals");
    let record = resp.record.as_ref().expect("record present");
    match record.fields.get("id") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(id)),
        }) => assert_eq!(id, "rec-1"),
        other => panic!("record.id is not a string: {:?}", other),
    }

    // database.query.run — request: dialect, sql, positional params (the
    // >2^53 bind value travels as a string), max_rows.
    let req: RunQueryRequest = serde_json::from_str(&fixture("database", "query.run.request.json"))
        .expect("database.query.run request unmarshals");
    assert!(req.sql.contains("accounts"));
    assert_eq!(req.params.len(), 1);
    match req.params[0].kind.as_ref() {
        Some(Kind::StringValue(param)) => assert_eq!(param, "9007199254740993"),
        other => panic!("params[0] must be a string, got {:?}", other),
    }
    assert_eq!(req.max_rows, Some(100));

    // database.query.run — spot: columns, NULLs, bigint and NUMERIC strings.
    let resp: RunQueryResponse =
        serde_json::from_str(&fixture("database", "query.run.response.json"))
            .expect("database.query.run response unmarshals");
    assert_eq!(resp.columns.len(), 4);
    assert_eq!(resp.columns[0].name, "id");
    assert_eq!(resp.columns[0].r#type, "int8");
    assert_eq!(resp.columns[3].r#type, "numeric");
    assert_eq!(resp.rows.len(), 2);

    let big = &resp.rows[0].values;
    match big[0].kind.as_ref() {
        Some(Kind::StringValue(id)) => {
            assert_eq!(id, "9007199254740993", "int64 > 2^53 travels as a string")
        }
        other => panic!("rows[0][0] must be a string, got {:?}", other),
    }
    match big[1].kind.as_ref() {
        Some(Kind::StringValue(name)) => assert_eq!(name, "Big-ID Account — ünïcodé ✓"),
        other => panic!("rows[0][1] must be a string, got {:?}", other),
    }
    // prost-serde's pbjson_types::Value represents JSON null as kind None.
    match big[2].kind.as_ref() {
        None | Some(Kind::NullValue(_)) => {}
        other => panic!("rows[0][2] must be null, got {:?}", other),
    }
    match big[3].kind.as_ref() {
        Some(Kind::StringValue(balance)) => assert_eq!(
            balance, "12345678901234567890.123456789012345678",
            "NUMERIC keeps its precision as a string"
        ),
        other => panic!("rows[0][3] must be a string, got {:?}", other),
    }
    match resp.rows[1].values[3].kind.as_ref() {
        Some(Kind::StringValue(balance)) => assert_eq!(balance, "0.000000000000000001"),
        other => panic!("rows[1][3] must be a string, got {:?}", other),
    }
    assert!(!resp.truncated);

    // omni.documents.list — spot: two entries, the unicode name and the next
    // page cursor.
    let resp: ListDocumentsResponse =
        serde_json::from_str(&fixture("omni", "documents.list.response.json"))
            .expect("omni.documents.list response unmarshals");
    assert_eq!(resp.documents.len(), 2);
    assert_eq!(resp.documents[0].name, "Q3 Revenue — ünïcodé ✓");
    assert_eq!(resp.next_cursor.as_deref(), Some("synthetic-cursor-2"));

    // omni.queries.run — spot: two Struct rows with region and revenue.
    let resp: OmniRunQueryResponse =
        serde_json::from_str(&fixture("omni", "queries.run.response.json"))
            .expect("omni.queries.run response unmarshals");
    assert_eq!(resp.rows.len(), 2);
    match resp.rows[0].fields.get("region") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(region)),
        }) => assert_eq!(region, "EMEA"),
        other => panic!("rows[0].region is not a string: {:?}", other),
    }
    match resp.rows[0].fields.get("revenue") {
        Some(pbjson_types::Value {
            kind: Some(Kind::NumberValue(revenue)),
        }) => assert_eq!(*revenue, 1234567.89),
        other => panic!("rows[0].revenue is not a number: {:?}", other),
    }

    // power_bi.datasets.list — spot: two datasets under value; the PascalCase
    // wire keys ConfiguredBy and TargetStorageMode decode via json_name.
    let resp: ListDatasetsResponse =
        serde_json::from_str(&fixture("power_bi", "datasets.list.response.json"))
            .expect("power_bi.datasets.list response unmarshals");
    assert_eq!(resp.value.len(), 2);
    assert_eq!(resp.value[0].name, "Sales Marketing");
    assert_eq!(
        resp.value[0].configured_by.as_deref(),
        Some("synthetic-admin@example.invalid")
    );
    assert_eq!(resp.value[0].target_storage_mode.as_deref(), Some("Abf"));

    // power_bi.datasets.execute_queries — spot: one result, one table, two
    // Struct rows with Region and Sales Amount.
    let resp: ExecuteQueriesResponse = serde_json::from_str(&fixture(
        "power_bi",
        "datasets.execute_queries.response.json",
    ))
    .expect("power_bi.datasets.execute_queries response unmarshals");
    assert_eq!(resp.results.len(), 1);
    assert_eq!(resp.results[0].tables.len(), 1);
    let rows = &resp.results[0].tables[0].rows;
    assert_eq!(rows.len(), 2);
    match rows[0].fields.get("Region") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(region)),
        }) => assert_eq!(region, "EMEA"),
        other => panic!("rows[0].Region is not a string: {:?}", other),
    }
    match rows[0].fields.get("Sales Amount") {
        Some(pbjson_types::Value {
            kind: Some(Kind::NumberValue(amount)),
        }) => assert_eq!(*amount, 1234567.89),
        other => panic!("rows[0].Sales Amount is not a number: {:?}", other),
    }
}

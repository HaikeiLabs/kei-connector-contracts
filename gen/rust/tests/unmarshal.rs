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

use haikeilabs_connector_contracts::kei::connectors::airtable::v1::{
    GetRecordResponse as AirtableGetRecordResponse,
    ListRecordsResponse as AirtableListRecordsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::database::v1::{
    RunQueryRequest, RunQueryResponse,
};
use haikeilabs_connector_contracts::kei::connectors::github::v1::{
    GetIssueResponse, GetPullRequestResponse, ListIssuesResponse, PutContentsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::http_api::v1::{
    GetRecordResponse, ListRecordsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::hubspot::v1::{
    GetContactResponse, ListContactsResponse, ListDealsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::{
    CreateCommentResponse, CreateIssueResponse, GetIssueResponse as LinearGetIssueResponse,
    GetTeamResponse, SearchIssuesResponse,
};
use haikeilabs_connector_contracts::kei::connectors::notion::v1::{
    GetPageResponse, QueryDatabaseResponse,
};
use haikeilabs_connector_contracts::kei::connectors::sentry::v1::{
    GetIssueResponse as SentryGetIssueResponse, ListIssuesResponse as SentryListIssuesResponse,
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

    // airtable.records.get — spot: record id, the Name string, and the null
    // Priority field.
    let resp: AirtableGetRecordResponse =
        serde_json::from_str(&fixture("airtable", "records.get.response.json"))
            .expect("airtable.records.get response unmarshals");
    assert_eq!(resp.id, "recSynthetic00003");
    let fields = resp.fields.as_ref().expect("fields present");
    match fields.fields.get("Name") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(name)),
        }) => assert_eq!(name, "Synthetic third task"),
        other => panic!("fields.Name must be a string, got {:?}", other),
    }
    match fields.fields.get("Priority") {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("fields.Priority must be null, got {:?}", other),
    }

    // airtable.records.list — spot: two records, the offset cursor, and the
    // null Priority field on the second record.
    let resp: AirtableListRecordsResponse =
        serde_json::from_str(&fixture("airtable", "records.list.response.json"))
            .expect("airtable.records.list response unmarshals");
    assert_eq!(resp.records.len(), 2);
    assert_eq!(resp.offset.as_deref(), Some("synthetic-next-offset"));
    let second = resp.records[1].fields.as_ref().expect("fields present");
    match second.fields.get("Priority") {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("records[1].fields.Priority must be null, got {:?}", other),
    }

    // hubspot.contacts.get — spot: contact id, the email string, and the null
    // lastname property.
    let resp: GetContactResponse =
        serde_json::from_str(&fixture("hubspot", "contacts.get.response.json"))
            .expect("hubspot.contacts.get response unmarshals");
    assert_eq!(resp.id, "701");
    let props = resp.properties.as_ref().expect("properties present");
    match props.fields.get("email") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(email)),
        }) => assert_eq!(email, "synthetic.contact@example.invalid"),
        other => panic!("properties.email must be a string, got {:?}", other),
    }
    match props.fields.get("lastname") {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("properties.lastname must be null, got {:?}", other),
    }

    // hubspot.contacts.list — spot: two results and the next cursor.
    let resp: ListContactsResponse =
        serde_json::from_str(&fixture("hubspot", "contacts.list.response.json"))
            .expect("hubspot.contacts.list response unmarshals");
    assert_eq!(resp.results.len(), 2);
    assert_eq!(resp.next_after.as_deref(), Some("synthetic-next-cursor"));
    let second = resp.results[1].properties.as_ref().expect("properties present");
    match second.fields.get("lastname") {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("results[1].properties.lastname must be null, got {:?}", other),
    }

    // hubspot.deals.list — spot: two results and the null description.
    let resp: ListDealsResponse =
        serde_json::from_str(&fixture("hubspot", "deals.list.response.json"))
            .expect("hubspot.deals.list response unmarshals");
    assert_eq!(resp.results.len(), 2);
    let first = resp.results[0].properties.as_ref().expect("properties present");
    match first.fields.get("dealname") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(name)),
        }) => assert_eq!(name, "Synthetic first deal"),
        other => panic!("properties.dealname must be a string, got {:?}", other),
    }
    match first.fields.get("description") {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("properties.description must be null, got {:?}", other),
    }
    // notion.pages.get — spot: the synthetic page id and the null public_url.
    let resp: GetPageResponse =
        serde_json::from_str(&fixture("notion", "pages.get.response.json"))
            .expect("notion.pages.get response unmarshals");
    assert_eq!(resp.id, "5f8b7e4a-0000-4000-8000-000000000001");
    assert!(resp.public_url.is_none(), "public_url must stay null");

    // notion.databases.query — spot: two results, the next cursor, and the
    // null public_url on the second page.
    let resp: QueryDatabaseResponse =
        serde_json::from_str(&fixture("notion", "databases.query.response.json"))
            .expect("notion.databases.query response unmarshals");
    assert_eq!(resp.results.len(), 2);
    assert_eq!(resp.next_cursor.as_deref(), Some("synthetic-next-cursor"));
    assert!(resp.has_more);
    assert!(resp.results[1].public_url.is_none(), "public_url must stay null");

    // sentry.issues.get — spot: short_id and the null share_id.
    let resp: SentryGetIssueResponse =
        serde_json::from_str(&fixture("sentry", "issues.get.response.json"))
            .expect("sentry.issues.get response unmarshals");
    assert_eq!(resp.short_id, "PROJ-103");
    assert!(resp.share_id.is_none(), "share_id must stay null");

    // sentry.issues.list — spot: two issues and the share id on the second.
    let resp: SentryListIssuesResponse =
        serde_json::from_str(&fixture("sentry", "issues.list.response.json"))
            .expect("sentry.issues.list response unmarshals");
    assert_eq!(resp.issues.len(), 2);
    assert_eq!(
        resp.issues[1].share_id.as_deref(),
        Some("synthetic-share-id-2")
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
}

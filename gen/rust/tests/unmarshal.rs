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
use haikeilabs_connector_contracts::kei::connectors::atlassian::v1::{
    CreateConfluencePageResponse, CreateJiraCommentResponse, CreateJiraIssueResponse,
    GetConfluencePageResponse, GetJiraIssueResponse, ListConfluencePagesResponse,
    SearchJiraIssuesResponse,
};
use haikeilabs_connector_contracts::kei::connectors::database::v1::{
    Dialect, RunQueryRequest, RunQueryResponse,
};
use haikeilabs_connector_contracts::kei::connectors::github::v1::{
    CreateBlobResponse, CreateCommitResponse, CreateIssueCommentResponse,
    CreateIssueResponse as GhCreateIssueResponse, CreatePullRequestResponse,
    CreatePullReviewResponse, CreateRefResponse, CreateTreeResponse, GetIssueResponse,
    GetPullRequestResponse, GetRefResponse, ListIssuesResponse, ListPullReviewsResponse,
    ListPullsResponse, ListWorkflowRunsResponse, PutContentsResponse, UpdatePullRequestResponse,
    UpdateRefResponse,
};
use haikeilabs_connector_contracts::kei::connectors::gmail::v1::{
    GetMessageResponse, ListMessagesResponse,
};
use haikeilabs_connector_contracts::kei::connectors::google_sheets::v1::{
    GetSpreadsheetResponse, GetValuesResponse,
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
use haikeilabs_connector_contracts::kei::connectors::omni::v1::{
    ListDocumentsResponse, RunQueryResponse as OmniRunQueryResponse,
};
use haikeilabs_connector_contracts::kei::connectors::power_bi::v1::{
    ExecuteQueriesResponse, ListDatasetsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::notion::v1::{
    CreatePageResponse, GetPageResponse, QueryDatabaseResponse, UpdatePageResponse,
};
use haikeilabs_connector_contracts::kei::connectors::sentry::v1::{
    GetIssueResponse as SentryGetIssueResponse, ListIssuesResponse as SentryListIssuesResponse,
    UpdateIssueResponse as SentryUpdateIssueResponse,
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

    // github.pulls.list — spot: head/base refs, author login, labels, the
    // >2^53 id, and the next_page cursor.
    let resp: ListPullsResponse =
        serde_json::from_str(&fixture("github", "pulls.list.response.json"))
            .expect("github.pulls.list response unmarshals");
    let pull = &resp.pulls[0];
    assert_eq!(
        pull.head.as_ref().expect("head present").r#ref,
        "feature/x"
    );
    assert_eq!(pull.base.as_ref().expect("base present").r#ref, "main");
    assert_eq!(
        pull.author.as_ref().expect("author present").login,
        "example-user"
    );
    assert_eq!(pull.labels[0].name, "synthetic-label");
    assert_eq!(pull.id, 9007199254740999, "int64 must not go through f64");
    assert_eq!(resp.next_page, Some(2));

    // github.actions.runs.list — spot: head_branch and the >2^53 run id.
    let resp: ListWorkflowRunsResponse =
        serde_json::from_str(&fixture("github", "actions.runs.list.response.json"))
            .expect("github.actions.runs.list response unmarshals");
    assert_eq!(resp.workflow_runs[0].head_branch, "feature/x");
    assert_eq!(resp.workflow_runs[0].id, 9007199254741001);

    // github.issues.create — spot: number, author login, labels.
    let resp: GhCreateIssueResponse =
        serde_json::from_str(&fixture("github", "issues.create.response.json"))
            .expect("github.issues.create response unmarshals");
    assert_eq!(resp.number, 42);
    assert_eq!(resp.author.as_ref().expect("author present").login, "example-user");
    assert_eq!(resp.labels[0].name, "synthetic-label");

    // github.issues.comments.create — spot: body and author login.
    let resp: CreateIssueCommentResponse =
        serde_json::from_str(&fixture("github", "issues.comments.create.response.json"))
            .expect("github.issues.comments.create response unmarshals");
    assert_eq!(resp.body, "Synthetic comment");
    assert_eq!(resp.author.as_ref().expect("author present").login, "example-user");

    // github.pulls.create — spot: number, draft flag, head/base refs.
    let resp: CreatePullRequestResponse =
        serde_json::from_str(&fixture("github", "pulls.create.response.json"))
            .expect("github.pulls.create response unmarshals");
    assert_eq!(resp.number, 43);
    assert!(resp.draft);
    assert_eq!(resp.head.as_ref().expect("head present").r#ref, "feature/x");
    assert_eq!(resp.base.as_ref().expect("base present").r#ref, "main");

    // github.pulls.update — spot: title and state.
    let resp: UpdatePullRequestResponse =
        serde_json::from_str(&fixture("github", "pulls.update.response.json"))
            .expect("github.pulls.update response unmarshals");
    assert_eq!(resp.title, "Synthetic updated title");
    assert_eq!(resp.state, "open");

    // github.pulls.reviews.list — spot: review state and author login.
    let resp: ListPullReviewsResponse =
        serde_json::from_str(&fixture("github", "pulls.reviews.list.response.json"))
            .expect("github.pulls.reviews.list response unmarshals");
    assert_eq!(resp.reviews[0].state, "APPROVED");
    assert_eq!(resp.reviews[0].author.as_ref().expect("author present").login, "example-reviewer");

    // github.pulls.reviews.create — spot: state and the >2^53 review id.
    let resp: CreatePullReviewResponse =
        serde_json::from_str(&fixture("github", "pulls.reviews.create.response.json"))
            .expect("github.pulls.reviews.create response unmarshals");
    assert_eq!(resp.state, "APPROVED");
    assert_eq!(resp.id, 9007199254741006);

    // github.git.refs.get — spot: full ref and the commit sha it points at.
    let resp: GetRefResponse =
        serde_json::from_str(&fixture("github", "git.refs.get.response.json"))
            .expect("github.git.refs.get response unmarshals");
    assert_eq!(resp.r#ref, "refs/heads/feature/x");
    assert_eq!(resp.object.as_ref().expect("object present").sha, "synthetic-ref-sha-1");

    // github.git.refs.create — spot: ref and object type.
    let resp: CreateRefResponse =
        serde_json::from_str(&fixture("github", "git.refs.create.response.json"))
            .expect("github.git.refs.create response unmarshals");
    assert_eq!(resp.r#ref, "refs/heads/feature/y");
    assert_eq!(resp.object.as_ref().expect("object present").r#type, "commit");

    // github.git.refs.update — spot: the moved sha.
    let resp: UpdateRefResponse =
        serde_json::from_str(&fixture("github", "git.refs.update.response.json"))
            .expect("github.git.refs.update response unmarshals");
    assert_eq!(resp.object.as_ref().expect("object present").sha, "synthetic-ref-sha-2");

    // github.git.blobs.create — spot: blob sha.
    let resp: CreateBlobResponse =
        serde_json::from_str(&fixture("github", "git.blobs.create.response.json"))
            .expect("github.git.blobs.create response unmarshals");
    assert_eq!(resp.sha, "synthetic-blob-sha-1");

    // github.git.trees.create — spot: tree sha and entry path.
    let resp: CreateTreeResponse =
        serde_json::from_str(&fixture("github", "git.trees.create.response.json"))
            .expect("github.git.trees.create response unmarshals");
    assert_eq!(resp.sha, "synthetic-tree-sha-1");
    assert_eq!(resp.tree[0].path, "notes.txt");

    // github.git.commits.create — spot: commit sha and single parent.
    let resp: CreateCommitResponse =
        serde_json::from_str(&fixture("github", "git.commits.create.response.json"))
            .expect("github.git.commits.create response unmarshals");
    assert_eq!(resp.sha, "synthetic-commit-sha-2");
    assert_eq!(resp.parents.len(), 1);

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

    // gmail.messages.list — spot: two messages and the next cursor.
    let resp: ListMessagesResponse =
        serde_json::from_str(&fixture("gmail", "messages.list.response.json"))
            .expect("gmail.messages.list response unmarshals");
    assert_eq!(resp.messages.len(), 2);
    assert_eq!(resp.next_page_token.as_deref(), Some("synthetic-next-token"));

    // gmail.messages.get — spot: message id, the snippet, and the null
    // internal date.
    let resp: GetMessageResponse =
        serde_json::from_str(&fixture("gmail", "messages.get.response.json"))
            .expect("gmail.messages.get response unmarshals");
    assert_eq!(resp.id, "msg-1001");
    assert_eq!(resp.snippet, "Synthetic message preview text");
    assert!(resp.internal_date.is_none(), "internal_date must be null");

    // google_sheets.values.get — spot: the echoed range, two rows, and the
    // null cell in the first row.
    let resp: GetValuesResponse =
        serde_json::from_str(&fixture("google_sheets", "values.get.response.json"))
            .expect("google_sheets.values.get response unmarshals");
    assert_eq!(resp.range, "Sheet1!A1:B2");
    assert_eq!(resp.values.len(), 2);
    match resp.values[0].values.get(1) {
        Some(pbjson_types::Value {
            kind: None | Some(Kind::NullValue(_)),
        }) => {}
        other => panic!("values[0][1] must be null, got {:?}", other),
    }

    // google_sheets.spreadsheets.get — spot: spreadsheet id and the title.
    let resp: GetSpreadsheetResponse =
        serde_json::from_str(&fixture("google_sheets", "spreadsheets.get.response.json"))
            .expect("google_sheets.spreadsheets.get response unmarshals");
    assert_eq!(resp.spreadsheet_id, "synthetic-sheet-1");
    let properties = resp.properties.as_ref().expect("properties present");
    assert_eq!(properties.title, "Synthetic sheet");

    // atlassian.jira.issues.get — spot: issue key, the summary, and the null
    // description.
    let resp: GetJiraIssueResponse =
        serde_json::from_str(&fixture("atlassian", "jira.issues.get.response.json"))
            .expect("atlassian.jira.issues.get response unmarshals");
    assert_eq!(resp.key, "SYN-101");
    let fields = resp.fields.as_ref().expect("fields present");
    assert_eq!(fields.summary, "Synthetic Jira issue");
    assert!(fields.description.is_none(), "description must be null");

    // atlassian.jira.issues.create — spot: the created issue's id and key.
    let resp: CreateJiraIssueResponse =
        serde_json::from_str(&fixture("atlassian", "jira.issues.create.response.json"))
            .expect("atlassian.jira.issues.create response unmarshals");
    assert_eq!(resp.id, "30001");
    assert_eq!(resp.key, "SYN-301");

    // atlassian.jira.issues.comments.create — spot: the comment id, its ADF
    // body and the created timestamp.
    let resp: CreateJiraCommentResponse =
        serde_json::from_str(&fixture("atlassian", "jira.issues.comments.create.response.json"))
            .expect("atlassian.jira.issues.comments.create response unmarshals");
    assert_eq!(resp.id, "40001");
    let comment_body = resp.body.as_ref().expect("body present");
    match comment_body.fields.get("type") {
        Some(pbjson_types::Value {
            kind: Some(Kind::StringValue(kind)),
        }) => assert_eq!(kind, "doc"),
        other => panic!("body.type must be a string, got {:?}", other),
    }
    assert_eq!(resp.created, "2026-02-11T12:00:00.000Z");

    // atlassian.jira.issues.search — spot: two issues, the next cursor, and
    // the null description on the second issue.
    let resp: SearchJiraIssuesResponse =
        serde_json::from_str(&fixture("atlassian", "jira.issues.search.response.json"))
            .expect("atlassian.jira.issues.search response unmarshals");
    assert_eq!(resp.issues.len(), 2);
    assert_eq!(resp.next_page_token.as_deref(), Some("synthetic-next-token"));
    let second = resp.issues[1].fields.as_ref().expect("fields present");
    assert!(second.description.is_none(), "description must be null");

    // atlassian.confluence.pages.get — spot: page id, the title, and the null
    // parent_id of a root page.
    let resp: GetConfluencePageResponse =
        serde_json::from_str(&fixture("atlassian", "confluence.pages.get.response.json"))
            .expect("atlassian.confluence.pages.get response unmarshals");
    assert_eq!(resp.id, 78643201);
    assert_eq!(resp.title, "Synthetic Confluence page");
    assert!(resp.parent_id.is_none(), "parent_id must be null");

    // atlassian.confluence.pages.create — spot: the created page's id and
    // title, and the null parent_id of a root page.
    let resp: CreateConfluencePageResponse =
        serde_json::from_str(&fixture("atlassian", "confluence.pages.create.response.json"))
            .expect("atlassian.confluence.pages.create response unmarshals");
    assert_eq!(resp.id, 78643205);
    assert_eq!(resp.title, "Synthetic created page");
    assert!(resp.parent_id.is_none(), "parent_id must be null");

    // atlassian.confluence.pages.list — spot: two page entries.
    let resp: ListConfluencePagesResponse =
        serde_json::from_str(&fixture("atlassian", "confluence.pages.list.response.json"))
            .expect("atlassian.confluence.pages.list response unmarshals");
    assert_eq!(resp.results.len(), 2);
    assert_eq!(resp.results[0].title, "Synthetic first page");

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

    // notion.pages.create — spot: the synthetic created page id and the null
    // public_url.
    let resp: CreatePageResponse =
        serde_json::from_str(&fixture("notion", "pages.create.response.json"))
            .expect("notion.pages.create response unmarshals");
    assert_eq!(resp.id, "d4e5f6a7-0000-4000-8000-000000000005");
    assert!(resp.public_url.is_none(), "public_url must stay null");

    // notion.pages.update — spot: the synthetic edit time and archived false.
    let resp: UpdatePageResponse =
        serde_json::from_str(&fixture("notion", "pages.update.response.json"))
            .expect("notion.pages.update response unmarshals");
    assert_eq!(resp.last_edited_time, "2026-02-11T09:45:00.000Z");
    assert!(!resp.archived);

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

    // sentry.issues.update — spot: the resolved status and the null share_id.
    let resp: SentryUpdateIssueResponse =
        serde_json::from_str(&fixture("sentry", "issues.update.response.json"))
            .expect("sentry.issues.update response unmarshals");
    assert_eq!(resp.status, "resolved");
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

// database.query.run per-dialect fixtures: each engine's dialect enum, one
// engine-specific typed value (jsonb/timestamptz, bigint unsigned/datetime(6),
// nvarchar unicode/datetimeoffset, NUMBER(38) and VARIANT as strings, uuid),
// and the NULL cell. The >2^53 integers and precision-critical decimals travel
// as strings; prost-serde's pbjson_types::Value represents JSON null as kind
// None. pbjson deserializers ignore unknown fields, so the unknown extra field
// each fixture carries is not asserted here (the strict-mode drift gate lives
// in the Go, Python and TypeScript suites).
#[test]
fn database_dialect_fixture_rows() {
    struct DialectCell {
        row: usize,
        col: usize,
        want: &'static str,
        is_null: bool,
    }
    struct DialectCase {
        name: &'static str,
        dialect: i32,
        col2_type: &'static str,
        checks: &'static [DialectCell],
    }
    const CASES: &[DialectCase] = &[
        DialectCase {
            name: "postgres",
            dialect: Dialect::Postgres as i32,
            col2_type: "timestamptz",
            checks: &[
                DialectCell { row: 0, col: 0, want: "9007199254740993", is_null: false },
                DialectCell { row: 0, col: 1, want: "{\"owner\": \"synthetic\", \"tags\": [\"ops\"]}", is_null: false },
                DialectCell { row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00", is_null: false },
                DialectCell { row: 0, col: 3, want: "12345678901234567890.123456789012345678", is_null: false },
                DialectCell { row: 1, col: 1, want: "", is_null: true },
            ],
        },
        DialectCase {
            name: "mysql",
            dialect: Dialect::Mysql as i32,
            col2_type: "datetime(6)",
            checks: &[
                DialectCell { row: 0, col: 0, want: "18446744073709551615", is_null: false },
                DialectCell { row: 0, col: 1, want: "{\"owner\": \"synthetic\", \"tags\": [\"ops\"]}", is_null: false },
                DialectCell { row: 0, col: 2, want: "2026-01-15 10:30:00.123456", is_null: false },
                DialectCell { row: 0, col: 3, want: "12345678901234567890.123456", is_null: false },
                DialectCell { row: 1, col: 1, want: "", is_null: true },
            ],
        },
        DialectCase {
            name: "mssql",
            dialect: Dialect::Mssql as i32,
            col2_type: "datetimeoffset",
            checks: &[
                DialectCell { row: 0, col: 0, want: "9007199254740993", is_null: false },
                DialectCell { row: 0, col: 1, want: "München — ünïcodé ✓", is_null: false },
                DialectCell { row: 0, col: 2, want: "2026-01-15T10:30:00.1234567+02:00", is_null: false },
                DialectCell { row: 0, col: 3, want: "12345678901234567890.1234567890", is_null: false },
                DialectCell { row: 1, col: 2, want: "", is_null: true },
            ],
        },
        DialectCase {
            name: "snowflake",
            dialect: Dialect::Snowflake as i32,
            col2_type: "NUMBER(38,18)",
            checks: &[
                DialectCell { row: 0, col: 0, want: "99999999999999999999999999999999999999", is_null: false },
                DialectCell { row: 0, col: 1, want: "Big-ID Account — ünïcodé ✓", is_null: false },
                DialectCell { row: 0, col: 2, want: "12345678901234567890.123456789012345678", is_null: false },
                DialectCell { row: 0, col: 3, want: "{\"owner\": \"synthetic\", \"tags\": [\"ops\"]}", is_null: false },
                DialectCell { row: 1, col: 3, want: "", is_null: true },
            ],
        },
        DialectCase {
            name: "supabase",
            dialect: Dialect::Supabase as i32,
            col2_type: "timestamptz",
            checks: &[
                DialectCell { row: 0, col: 0, want: "6f9619ff-8b86-d011-b42d-00c04fc964ff", is_null: false },
                DialectCell { row: 0, col: 1, want: "{\"owner\": \"synthetic\", \"tags\": [\"ops\"]}", is_null: false },
                DialectCell { row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00", is_null: false },
                DialectCell { row: 0, col: 3, want: "12345678901234567890.123456789012345678", is_null: false },
                DialectCell { row: 1, col: 1, want: "", is_null: true },
            ],
        },
    ];

    for tc in CASES {
        let req: RunQueryRequest = serde_json::from_str(&fixture(
            "database",
            &format!("query.run.{}.request.json", tc.name),
        ))
        .unwrap_or_else(|err| panic!("query.run.{} request unmarshals: {err}", tc.name));
        assert_eq!(req.dialect, tc.dialect, "query.run.{} dialect", tc.name);

        let resp: RunQueryResponse = serde_json::from_str(&fixture(
            "database",
            &format!("query.run.{}.response.json", tc.name),
        ))
        .unwrap_or_else(|err| panic!("query.run.{} response unmarshals: {err}", tc.name));
        assert_eq!(resp.columns.len(), 4, "query.run.{} columns", tc.name);
        assert_eq!(resp.columns[2].r#type, tc.col2_type, "query.run.{} columns[2].type", tc.name);
        assert_eq!(resp.rows.len(), 2, "query.run.{} rows", tc.name);
        for cell in tc.checks {
            let value = &resp.rows[cell.row].values[cell.col];
            if cell.is_null {
                match value.kind.as_ref() {
                    None | Some(Kind::NullValue(_)) => {}
                    other => panic!("query.run.{} rows[{}][{}] must be null, got {:?}", tc.name, cell.row, cell.col, other),
                }
            } else {
                match value.kind.as_ref() {
                    Some(Kind::StringValue(got)) => assert_eq!(
                        got, cell.want,
                        "query.run.{} rows[{}][{}] keeps its synthetic value",
                        tc.name, cell.row, cell.col
                    ),
                    other => panic!("query.run.{} rows[{}][{}] must be a string, got {:?}", tc.name, cell.row, cell.col, other),
                }
            }
        }
        assert!(!resp.truncated, "query.run.{} truncated", tc.name);
    }
}

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
    CreateBlobResponse, CreateCommitResponse, CreateIssueCommentResponse,
    CreateIssueResponse as GhCreateIssueResponse, CreatePullRequestResponse,
    CreatePullReviewResponse, CreateRefResponse, CreateTreeResponse, GetIssueResponse,
    GetPullRequestResponse, GetRefResponse, ListIssuesResponse, ListPullReviewsResponse,
    ListPullsResponse, ListWorkflowRunsResponse, PutContentsResponse, UpdatePullRequestResponse,
    UpdateRefResponse,
};
use haikeilabs_connector_contracts::kei::connectors::http_api::v1::{
    GetRecordResponse, ListRecordsResponse,
};
use haikeilabs_connector_contracts::kei::connectors::linear::v1::{
    CreateCommentResponse, CreateIssueResponse, GetIssueResponse as LinearGetIssueResponse,
    GetTeamResponse, SearchIssuesResponse,
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

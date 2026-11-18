import { readFileSync } from "node:fs";
import { fromJson, type DescMessage, type JsonValue } from "@bufbuild/protobuf";
import { describe, expect, it } from "vitest";
import { isInteger, parse } from "lossless-json";
import * as gmailMessages from "../src/kei/connectors/gmail/v1/messages_pb.js";
import * as googleSheets from "../src/kei/connectors/google_sheets/v1/spreadsheets_pb.js";
import * as atlassianJira from "../src/kei/connectors/atlassian/v1/jira_pb.js";
import * as atlassianConfluence from "../src/kei/connectors/atlassian/v1/confluence_pb.js";
import * as ghIssues from "../src/kei/connectors/github/v1/issues_pb.js";
import * as ghContents from "../src/kei/connectors/github/v1/contents_pb.js";
import * as ghPulls from "../src/kei/connectors/github/v1/pulls_pb.js";
import * as ghActions from "../src/kei/connectors/github/v1/actions_pb.js";
import * as ghGit from "../src/kei/connectors/github/v1/git_pb.js";
import * as linear from "../src/kei/connectors/linear/v1/issues_pb.js";
import * as linearComments from "../src/kei/connectors/linear/v1/comments_pb.js";
import * as linearTeams from "../src/kei/connectors/linear/v1/teams_pb.js";
import * as airtableRecords from "../src/kei/connectors/airtable/v1/records_pb.js";
import * as hubspotContacts from "../src/kei/connectors/hubspot/v1/contacts_pb.js";
import * as hubspotDeals from "../src/kei/connectors/hubspot/v1/deals_pb.js";
import * as notionPages from "../src/kei/connectors/notion/v1/pages_pb.js";
import * as notionDatabases from "../src/kei/connectors/notion/v1/databases_pb.js";
import * as sentry from "../src/kei/connectors/sentry/v1/issues_pb.js";
import * as http from "../src/kei/connectors/http_api/v1/records_pb.js";
import * as db from "../src/kei/connectors/database/v1/query_pb.js";
import * as mcp from "../src/kei/connectors/mcp/v1/tool_pb.js";
import * as omniDocs from "../src/kei/connectors/omni/v1/documents_pb.js";
import * as omniQueries from "../src/kei/connectors/omni/v1/queries_pb.js";
import * as pbi from "../src/kei/connectors/power_bi/v1/datasets_pb.js";

function fixtureText(op: string, kind: string): string {
  const [connector, ...name] = op.split(".");
  return readFileSync(new URL(`../../../fixtures/operations/${connector}/${name.join(".")}.${kind}.json`, import.meta.url), "utf8");
}

function fixture(op: string, kind: string): JsonValue {
  return parse(fixtureText(op, kind), undefined, {
    parseNumber: (value) => {
      if (isInteger(value)) {
        const integer = BigInt(value);
        if (integer > BigInt(Number.MAX_SAFE_INTEGER) || integer < BigInt(Number.MIN_SAFE_INTEGER)) return value;
      }
      return Number(value);
    },
  }) as JsonValue;
}

// Each response includes an unknown field: lenient parsing accepts provider drift,
// while strict parsing detects it. Requests are also checked for compatibility.
const ops = [
  ["gmail.messages.list", gmailMessages.ListMessagesRequestSchema, gmailMessages.ListMessagesResponseSchema, (r: any) => { expect(r.messages).toHaveLength(2); expect(r.nextPageToken).toBe("synthetic-next-token"); }],
  ["gmail.messages.get", gmailMessages.GetMessageRequestSchema, gmailMessages.GetMessageResponseSchema, (r: any) => { expect(r.id).toBe("msg-1001"); expect(r.snippet).toBe("Synthetic message preview text"); expect(r.internalDate).toBeUndefined(); }],
  ["gmail.drafts.create", gmailMessages.CreateDraftRequestSchema, gmailMessages.CreateDraftResponseSchema, (r: any) => { expect(r.id).toBe("draft-3001"); expect(r.message.id).toBe("msg-1002"); expect(r.message.threadId).toBe("thr-2002"); }],
  ["google_sheets.values.get", googleSheets.GetValuesRequestSchema, googleSheets.GetValuesResponseSchema, (r: any) => { expect(r.range).toBe("Sheet1!A1:B2"); expect(r.values).toHaveLength(2); expect(r.values[0].values[1].kind.case).toBe("nullValue"); }],
  ["google_sheets.spreadsheets.get", googleSheets.GetSpreadsheetRequestSchema, googleSheets.GetSpreadsheetResponseSchema, (r: any) => { expect(r.spreadsheetId).toBe("synthetic-sheet-1"); expect(r.properties.title).toBe("Synthetic sheet"); }],
  ["google_sheets.values.update", googleSheets.UpdateValuesRequestSchema, googleSheets.UpdateValuesResponseSchema, (r: any) => { expect(r.updatedRange).toBe("Sheet1!A1:B2"); expect(r.updatedCells).toBe(4); }],
  ["atlassian.jira.issues.get", atlassianJira.GetJiraIssueRequestSchema, atlassianJira.GetJiraIssueResponseSchema, (r: any) => { expect(r.key).toBe("SYN-101"); expect(r.fields.summary).toBe("Synthetic Jira issue"); expect(r.fields.description).toBeUndefined(); }],
  ["atlassian.jira.issues.create", atlassianJira.CreateJiraIssueRequestSchema, atlassianJira.CreateJiraIssueResponseSchema, (r: any) => { expect(r.id).toBe("30001"); expect(r.key).toBe("SYN-301"); }],
  ["atlassian.jira.issues.comments.create", atlassianJira.CreateJiraCommentRequestSchema, atlassianJira.CreateJiraCommentResponseSchema, (r: any) => { expect(r.id).toBe("40001"); expect(r.body["type"]).toBe("doc"); expect(r.created).toBe("2026-02-11T12:00:00.000Z"); }],
  ["atlassian.jira.issues.search", atlassianJira.SearchJiraIssuesRequestSchema, atlassianJira.SearchJiraIssuesResponseSchema, (r: any) => { expect(r.issues).toHaveLength(2); expect(r.nextPageToken).toBe("synthetic-next-token"); expect(r.issues[1].fields.description).toBeUndefined(); }],
  ["atlassian.confluence.pages.get", atlassianConfluence.GetConfluencePageRequestSchema, atlassianConfluence.GetConfluencePageResponseSchema, (r: any) => { expect(r.id).toBe(78643201n); expect(r.title).toBe("Synthetic Confluence page"); expect(r.parentId).toBeUndefined(); }],
  ["atlassian.confluence.pages.create", atlassianConfluence.CreateConfluencePageRequestSchema, atlassianConfluence.CreateConfluencePageResponseSchema, (r: any) => { expect(r.id).toBe(78643205n); expect(r.title).toBe("Synthetic created page"); expect(r.parentId).toBeUndefined(); }],
  ["atlassian.confluence.pages.list", atlassianConfluence.ListConfluencePagesRequestSchema, atlassianConfluence.ListConfluencePagesResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.results[0].title).toBe("Synthetic first page"); }],
  ["github.issues.get", ghIssues.GetIssueRequestSchema, ghIssues.GetIssueResponseSchema, (r: any) => { expect(r.id).toBe(9007199254740993n); expect(r.number).toBe(1347n); }],
  ["github.issues.list", ghIssues.ListIssuesRequestSchema, ghIssues.ListIssuesResponseSchema, (r: any) => expect(r.issues).toHaveLength(2)],
  ["github.contents.put", ghContents.PutContentsRequestSchema, ghContents.PutContentsResponseSchema, (r: any) => expect(r.content.name).toBe("notes.txt")],
  ["github.pulls.get", ghPulls.GetPullRequestRequestSchema, ghPulls.GetPullRequestResponseSchema, (r: any) => expect(r.number).toBe(42n)],
  ["github.pulls.list", ghPulls.ListPullsRequestSchema, ghPulls.ListPullsResponseSchema, (r: any) => { expect(r.pulls[0].head.ref).toBe("feature/x"); expect(r.pulls[0].base.ref).toBe("main"); expect(r.pulls[0].author.login).toBe("example-user"); expect(r.pulls[0].labels[0].name).toBe("synthetic-label"); expect(r.pulls[0].id).toBe(9007199254740999n); expect(r.nextPage).toBe(2n); }],
  ["github.actions.runs.list", ghActions.ListWorkflowRunsRequestSchema, ghActions.ListWorkflowRunsResponseSchema, (r: any) => { expect(r.workflowRuns[0].headBranch).toBe("feature/x"); expect(r.workflowRuns[0].id).toBe(9007199254741001n); }],
  ["github.issues.create", ghIssues.CreateIssueRequestSchema, ghIssues.CreateIssueResponseSchema, (r: any) => { expect(r.number).toBe(42n); expect(r.author.login).toBe("example-user"); expect(r.labels[0].name).toBe("synthetic-label"); }],
  ["github.issues.comments.create", ghIssues.CreateIssueCommentRequestSchema, ghIssues.CreateIssueCommentResponseSchema, (r: any) => { expect(r.body).toBe("Synthetic comment"); expect(r.author.login).toBe("example-user"); }],
  ["github.pulls.create", ghPulls.CreatePullRequestRequestSchema, ghPulls.CreatePullRequestResponseSchema, (r: any) => { expect(r.number).toBe(43n); expect(r.draft).toBe(true); expect(r.head.ref).toBe("feature/x"); }],
  ["github.pulls.update", ghPulls.UpdatePullRequestRequestSchema, ghPulls.UpdatePullRequestResponseSchema, (r: any) => { expect(r.title).toBe("Synthetic updated title"); expect(r.state).toBe("open"); }],
  ["github.pulls.reviews.list", ghPulls.ListPullReviewsRequestSchema, ghPulls.ListPullReviewsResponseSchema, (r: any) => { expect(r.reviews[0].state).toBe("APPROVED"); expect(r.reviews[0].author.login).toBe("example-reviewer"); }],
  ["github.pulls.reviews.create", ghPulls.CreatePullReviewRequestSchema, ghPulls.CreatePullReviewResponseSchema, (r: any) => { expect(r.state).toBe("APPROVED"); expect(r.id).toBe(9007199254741006n); }],
  ["github.git.refs.get", ghGit.GetRefRequestSchema, ghGit.GetRefResponseSchema, (r: any) => { expect(r.ref).toBe("refs/heads/feature/x"); expect(r.object.sha).toBe("synthetic-ref-sha-1"); }],
  ["github.git.refs.create", ghGit.CreateRefRequestSchema, ghGit.CreateRefResponseSchema, (r: any) => { expect(r.ref).toBe("refs/heads/feature/y"); expect(r.object.type).toBe("commit"); }],
  ["github.git.refs.update", ghGit.UpdateRefRequestSchema, ghGit.UpdateRefResponseSchema, (r: any) => expect(r.object.sha).toBe("synthetic-ref-sha-2")],
  ["github.git.blobs.create", ghGit.CreateBlobRequestSchema, ghGit.CreateBlobResponseSchema, (r: any) => expect(r.sha).toBe("synthetic-blob-sha-1")],
  ["github.git.trees.create", ghGit.CreateTreeRequestSchema, ghGit.CreateTreeResponseSchema, (r: any) => { expect(r.sha).toBe("synthetic-tree-sha-1"); expect(r.tree[0].path).toBe("notes.txt"); }],
  ["github.git.commits.create", ghGit.CreateCommitRequestSchema, ghGit.CreateCommitResponseSchema, (r: any) => { expect(r.sha).toBe("synthetic-commit-sha-2"); expect(r.parents).toHaveLength(1); }],
  ["linear.issues.get", linear.GetIssueRequestSchema, linear.GetIssueResponseSchema, (r: any) => { expect(r.issue.identifier).toBe("ENG-101"); expect(r.issue.labels[0].name).toBe("synthetic-label"); }],
  ["linear.issues.create", linear.CreateIssueRequestSchema, linear.CreateIssueResponseSchema, (r: any) => expect(r.success).toBe(true)],
  ["linear.issues.search", linear.SearchIssuesRequestSchema, linear.SearchIssuesResponseSchema, (r: any) => expect(r.nodes).toHaveLength(2)],
  ["linear.teams.get", linearTeams.GetTeamRequestSchema, linearTeams.GetTeamResponseSchema, (r: any) => expect(r.team.key).toBe("ENG")],
  ["linear.comments.create", linearComments.CreateCommentRequestSchema, linearComments.CreateCommentResponseSchema, (r: any) => { expect(r.success).toBe(true); expect(r.comment.body).toBe("Synthetic comment"); }],
  ["airtable.records.get", airtableRecords.GetRecordRequestSchema, airtableRecords.GetRecordResponseSchema, (r: any) => { expect(r.id).toBe("recSynthetic00003"); expect(r.fields["Name"]).toBe("Synthetic third task"); expect(r.fields["Priority"]).toBeNull(); }],
  ["airtable.records.create", airtableRecords.CreateRecordRequestSchema, airtableRecords.CreateRecordResponseSchema, (r: any) => { expect(r.id).toBe("recSynthetic00004"); expect(r.fields["Status"]).toBe("Todo"); }],
  ["airtable.records.update", airtableRecords.UpdateRecordRequestSchema, airtableRecords.UpdateRecordResponseSchema, (r: any) => { expect(r.fields["Status"]).toBe("Done"); expect(r.fields["Priority"]).toBeNull(); }],
  ["airtable.records.list", airtableRecords.ListRecordsRequestSchema, airtableRecords.ListRecordsResponseSchema, (r: any) => { expect(r.records).toHaveLength(2); expect(r.offset).toBe("synthetic-next-offset"); expect(r.records[1].fields["Priority"]).toBeNull(); }],
  ["hubspot.contacts.get", hubspotContacts.GetContactRequestSchema, hubspotContacts.GetContactResponseSchema, (r: any) => { expect(r.id).toBe("701"); expect(r.properties.email).toBe("synthetic.contact@example.invalid"); expect(r.properties.lastname).toBeNull(); }],
  ["hubspot.contacts.create", hubspotContacts.CreateContactRequestSchema, hubspotContacts.CreateContactResponseSchema, (r: any) => { expect(r.id).toBe("701550001"); expect(r.properties.email).toBe("synthetic.new@example.invalid"); }],
  ["hubspot.contacts.update", hubspotContacts.UpdateContactRequestSchema, hubspotContacts.UpdateContactResponseSchema, (r: any) => { expect(r.properties.lifecyclestage).toBe("customer"); }],
  ["hubspot.contacts.list", hubspotContacts.ListContactsRequestSchema, hubspotContacts.ListContactsResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.nextAfter).toBe("synthetic-next-cursor"); }],
  ["hubspot.deals.list", hubspotDeals.ListDealsRequestSchema, hubspotDeals.ListDealsResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.results[0].properties.dealname).toBe("Synthetic first deal"); }],
  ["notion.pages.get", notionPages.GetPageRequestSchema, notionPages.GetPageResponseSchema, (r: any) => { expect(r.id).toBe("5f8b7e4a-0000-4000-8000-000000000001"); expect(r.publicUrl).toBeUndefined(); }],
  ["notion.pages.create", notionPages.CreatePageRequestSchema, notionPages.CreatePageResponseSchema, (r: any) => { expect(r.id).toBe("d4e5f6a7-0000-4000-8000-000000000005"); expect(r.publicUrl).toBeUndefined(); }],
  ["notion.pages.update", notionPages.UpdatePageRequestSchema, notionPages.UpdatePageResponseSchema, (r: any) => { expect(r.lastEditedTime).toBe("2026-02-11T09:45:00.000Z"); expect(r.archived).toBe(false); }],
  ["notion.databases.query", notionDatabases.QueryDatabaseRequestSchema, notionDatabases.QueryDatabaseResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.nextCursor).toBe("synthetic-next-cursor"); expect(r.results[1].publicUrl).toBeUndefined(); }],
  ["sentry.issues.get", sentry.GetIssueRequestSchema, sentry.GetIssueResponseSchema, (r: any) => { expect(r.shortId).toBe("PROJ-103"); expect(r.shareId).toBeUndefined(); }],
  ["sentry.issues.update", sentry.UpdateIssueRequestSchema, sentry.UpdateIssueResponseSchema, (r: any) => { expect(r.status).toBe("resolved"); expect(r.shortId).toBe("PROJ-103"); expect(r.shareId).toBeUndefined(); }],
  ["sentry.issues.list", sentry.ListIssuesRequestSchema, sentry.ListIssuesResponseSchema, (r: any) => { expect(r.issues).toHaveLength(2); expect(r.issues[1].shareId).toBe("synthetic-share-id-2"); }],
  ["http_api.records.list", http.ListRecordsRequestSchema, http.ListRecordsResponseSchema, (r: any) => expect(r.records).toHaveLength(2)],
  ["http_api.records.get", http.GetRecordRequestSchema, http.GetRecordResponseSchema, (r: any) => expect(r.record.id).toBe("rec-1")],
  ["database.query.run", db.RunQueryRequestSchema, db.RunQueryResponseSchema, (r: any) => expect(r.columns[0].name).toBe("id")],
  ["mcp.tool.call", mcp.ToolCallSchema, mcp.ToolResultSchema, (r: any) => expect(r.content).toHaveLength(2)],
  ["omni.documents.list", omniDocs.ListDocumentsRequestSchema, omniDocs.ListDocumentsResponseSchema, (r: any) => { expect(r.documents).toHaveLength(2); expect(r.documents[0].name).toBe("Q3 Revenue — ünïcodé ✓"); expect(r.nextCursor).toBe("synthetic-cursor-2"); }],
  ["omni.queries.run", omniQueries.RunQueryRequestSchema, omniQueries.RunQueryResponseSchema, (r: any) => { expect(r.rows).toHaveLength(2); expect(r.rows[0].region).toBe("EMEA"); expect(r.rows[0].revenue).toBe(1234567.89); }],
  ["power_bi.datasets.list", pbi.ListDatasetsRequestSchema, pbi.ListDatasetsResponseSchema, (r: any) => { expect(r.value).toHaveLength(2); expect(r.value[0].name).toBe("Sales Marketing"); expect(r.value[0].configuredBy).toBe("synthetic-admin@example.invalid"); expect(r.value[0].targetStorageMode).toBe("Abf"); }],
  ["power_bi.datasets.execute_queries", pbi.ExecuteQueriesRequestSchema, pbi.ExecuteQueriesResponseSchema, (r: any) => { expect(r.results).toHaveLength(1); expect(r.results[0].tables[0].rows).toHaveLength(2); expect(r.results[0].tables[0].rows[0].Region).toBe("EMEA"); expect(r.results[0].tables[0].rows[0]["Sales Amount"]).toBe(1234567.89); }],
] as const;

describe.each(ops)("%s", (op, request, response, spot) => {
  it("parses request and response with unknown fields ignored", () => {
    expect(fromJson(request as DescMessage, fixture(op, "request"), { ignoreUnknownFields: true })).toBeDefined();
    spot(fromJson(response as DescMessage, fixture(op, "response"), { ignoreUnknownFields: true }));
  });
  it("rejects unknown response fields in strict mode", () => {
    expect(() => fromJson(response as DescMessage, fixture(op, "response"))).toThrow();
  });
});

it("shows JSON.parse loses precision for the provider issue id", () => {
  const raw = JSON.parse(fixtureText("github.issues.get", "response"));
  expect(BigInt(raw.id)).toBe(9007199254740992n);
  expect(BigInt(raw.id)).not.toBe(9007199254740993n);
});

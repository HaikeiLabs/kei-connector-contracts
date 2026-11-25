import { readFileSync } from "node:fs";
import { fromJson, type DescMessage, type JsonValue } from "@bufbuild/protobuf";
import { describe, expect, it } from "vitest";
import { isInteger, parse } from "lossless-json";
import * as gmailMessages from "../src/kei/connectors/gmail/v1/messages_pb.js";
import * as gcCalendars from "../src/kei/connectors/google_calendar/v1/calendars_pb.js";
import * as gcEvents from "../src/kei/connectors/google_calendar/v1/events_pb.js";
import * as googleSheets from "../src/kei/connectors/google_sheets/v1/spreadsheets_pb.js";
import * as atlassianJira from "../src/kei/connectors/atlassian/v1/jira_pb.js";
import * as atlassianConfluence from "../src/kei/connectors/atlassian/v1/confluence_pb.js";
import * as ghIssues from "../src/kei/connectors/github/v1/issues_pb.js";
import * as ghContents from "../src/kei/connectors/github/v1/contents_pb.js";
import * as ghPulls from "../src/kei/connectors/github/v1/pulls_pb.js";
import * as ghActions from "../src/kei/connectors/github/v1/actions_pb.js";
import * as ghCommits from "../src/kei/connectors/github/v1/commits_pb.js";
import * as ghGit from "../src/kei/connectors/github/v1/git_pb.js";
import * as ghRepos from "../src/kei/connectors/github/v1/repos_pb.js";
import * as ghChecks from "../src/kei/connectors/github/v1/checks_pb.js";
import * as linear from "../src/kei/connectors/linear/v1/issues_pb.js";
import * as linearComments from "../src/kei/connectors/linear/v1/comments_pb.js";
import * as linearTeams from "../src/kei/connectors/linear/v1/teams_pb.js";
import * as linearProjects from "../src/kei/connectors/linear/v1/projects_pb.js";
import * as linearCycles from "../src/kei/connectors/linear/v1/cycles_pb.js";
import * as airtableRecords from "../src/kei/connectors/airtable/v1/records_pb.js";
import * as hubspotContacts from "../src/kei/connectors/hubspot/v1/contacts_pb.js";
import * as hubspotDeals from "../src/kei/connectors/hubspot/v1/deals_pb.js";
import * as notionPages from "../src/kei/connectors/notion/v1/pages_pb.js";
import * as notionDatabases from "../src/kei/connectors/notion/v1/databases_pb.js";
import * as notionSearch from "../src/kei/connectors/notion/v1/search_pb.js";
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
  ["google_sheets.values.append", googleSheets.AppendValuesRequestSchema, googleSheets.AppendValuesResponseSchema, (r: any) => { expect(r.spreadsheetId).toBe("synthetic-sheet-1"); expect(r.tableRange).toBe("Sheet1!A1:B2"); expect(r.updates.updatedRange).toBe("Sheet1!A3:B4"); expect(r.updates.updatedCells).toBe(4); }],
  ["google_calendar.calendars.list", gcCalendars.ListCalendarsRequestSchema, gcCalendars.ListCalendarsResponseSchema, (r: any) => { expect(r.items).toHaveLength(2); expect(r.items[0].primary).toBe(true); expect(r.items[1].description).toBeUndefined(); expect(r.nextPageToken).toBe("synthetic-next-token"); }],
  ["google_calendar.events.list", gcEvents.ListEventsRequestSchema, gcEvents.ListEventsResponseSchema, (r: any) => { expect(r.items).toHaveLength(2); expect(r.items[0].id).toBe("9007199254740993"); expect(r.items[0].start.dateTime).toBe("2026-01-15T10:30:00+01:00"); expect(r.items[1].start.date).toBe("2026-01-20"); expect(r.items[1].start.dateTime).toBeUndefined(); }],
  ["google_calendar.events.get", gcEvents.GetEventRequestSchema, gcEvents.GetEventResponseSchema, (r: any) => { expect(r.id).toBe("9007199254740993"); expect(r.description).toBeUndefined(); expect(r.location).toBe("Synthetic Conference Room"); }],
  ["google_calendar.events.create", gcEvents.CreateEventRequestSchema, gcEvents.CreateEventResponseSchema, (r: any) => { expect(r.id).toBe("evt-synthetic-created-1"); expect(r.start.dateTime).toBe("2026-02-10T14:00:00+01:00"); }],
  ["google_calendar.events.update", gcEvents.UpdateEventRequestSchema, gcEvents.UpdateEventResponseSchema, (r: any) => { expect(r.summary).toBe("Synthetic updated event"); expect(r.location).toBeUndefined(); }],
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
  ["github.pulls.get", ghPulls.GetPullRequestRequestSchema, ghPulls.GetPullRequestResponseSchema, (r: any) => { expect(r.number).toBe(42n); expect(r.diffUrl).toBe("https://example.invalid/example-org/example-repo/pull/42.diff"); expect(r.patchUrl).toBe("https://example.invalid/example-org/example-repo/pull/42.patch"); }],
  ["github.pulls.list", ghPulls.ListPullsRequestSchema, ghPulls.ListPullsResponseSchema, (r: any) => { expect(r.pulls[0].head.ref).toBe("feature/x"); expect(r.pulls[0].base.ref).toBe("main"); expect(r.pulls[0].author.login).toBe("example-user"); expect(r.pulls[0].labels[0].name).toBe("synthetic-label"); expect(r.pulls[0].id).toBe(9007199254740999n); expect(r.nextPage).toBe(2n); }],
  ["github.actions.runs.list", ghActions.ListWorkflowRunsRequestSchema, ghActions.ListWorkflowRunsResponseSchema, (r: any) => { expect(r.workflowRuns[0].headBranch).toBe("feature/x"); expect(r.workflowRuns[0].id).toBe(9007199254741001n); expect(r.workflowRuns[0].actor.login).toBe("example-user"); expect(r.workflowRuns[0].headSha).toBe("synthetic-head-sha-1"); }],
  ["github.actions.workflows.runs.list", ghActions.ListWorkflowRunsByWorkflowRequestSchema, ghActions.ListWorkflowRunsByWorkflowResponseSchema, (r: any) => { expect(r.workflowRuns).toHaveLength(1); expect(r.workflowRuns[0].actor.login).toBe("example-user"); expect(r.workflowRuns[0].headSha).toBe("synthetic-head-sha-1"); }],
  ["github.repos.contents.get", ghContents.GetContentsRequestSchema, ghContents.GetContentsResponseSchema, (r: any) => { expect(r.encoding).toBe("base64"); expect(r.content).toBe("c3ludGhldGljLWNvbnRlbnQ="); }],
  ["github.repos.contents.list", ghContents.ListContentsRequestSchema, ghContents.ListContentsResponseSchema, (r: any) => { expect(r.entries).toHaveLength(2); expect(r.entries[0].type).toBe("file"); expect(r.entries[1].type).toBe("dir"); }],
  ["github.repos.commits.list", ghCommits.ListCommitsRequestSchema, ghCommits.ListCommitsResponseSchema, (r: any) => { expect(r.commits[0].sha).toBe("synthetic-commit-sha-3"); expect(r.commits[0].author.login).toBe("example-user"); expect(r.nextPage).toBe(2n); }],
  ["github.repos.commits.get", ghCommits.GetCommitRequestSchema, ghCommits.GetCommitResponseSchema, (r: any) => { expect(r.commit.message).toBe("Synthetic commit message"); expect(r.author.login).toBe("example-user"); }],
  ["github.actions.runs.attempts.get", ghActions.GetRunAttemptRequestSchema, ghActions.GetRunAttemptResponseSchema, (r: any) => { expect(r.runAttempt).toBe(2); expect(r.headSha).toBe("synthetic-head-sha-1"); }],
  ["github.actions.runs.jobs.list", ghActions.ListRunJobsRequestSchema, ghActions.ListRunJobsResponseSchema, (r: any) => { expect(r.jobs[0].id).toBe(9007199254741005n); expect(r.jobs[0].runnerName).toBe("synthetic-runner-1"); }],
  ["github.issues.create", ghIssues.CreateIssueRequestSchema, ghIssues.CreateIssueResponseSchema, (r: any) => { expect(r.number).toBe(42n); expect(r.author.login).toBe("example-user"); expect(r.labels[0].name).toBe("synthetic-label"); }],
  ["github.issues.comments.create", ghIssues.CreateIssueCommentRequestSchema, ghIssues.CreateIssueCommentResponseSchema, (r: any) => { expect(r.body).toBe("Synthetic comment"); expect(r.author.login).toBe("example-user"); }],
  ["github.repos.get", ghRepos.GetRepositoryRequestSchema, ghRepos.GetRepositoryResponseSchema, (r: any) => { expect(r.id).toBe(9007199254741011n); expect(r.fullName).toBe("example-org/example-repo"); expect(r.owner.login).toBe("example-user"); }],
  ["github.checks.runs.list", ghChecks.ListCheckRunsRequestSchema, ghChecks.ListCheckRunsResponseSchema, (r: any) => { expect(r.checkRuns[0].id).toBe(9007199254741012n); expect(r.checkRuns[0].conclusion).toBe("success"); expect(r.totalCount).toBe(2n); }],
  ["github.issues.update", ghIssues.UpdateIssueRequestSchema, ghIssues.UpdateIssueResponseSchema, (r: any) => { expect(r.number).toBe(42n); expect(r.state).toBe("closed"); expect(r.author.login).toBe("example-user"); }],
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
  ["linear.projects.list", linearProjects.ListProjectsRequestSchema, linearProjects.ListProjectsResponseSchema, (r: any) => { expect(r.nodes[0].name).toBe("Synthetic pilot project"); expect(r.pageInfo.hasNextPage).toBe(true); expect(r.pageInfo.endCursor).toBe("synthetic-cursor-1"); }],
  ["linear.cycles.list", linearCycles.ListCyclesRequestSchema, linearCycles.ListCyclesResponseSchema, (r: any) => { expect(r.nodes[0].number).toBe(7); expect(r.nodes[0].isActive).toBe(true); }],
  ["linear.issues.update", linear.UpdateIssueRequestSchema, linear.UpdateIssueResponseSchema, (r: any) => { expect(r.success).toBe(true); expect(r.issue.title).toBe("Synthetic updated issue"); }],
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
  ["notion.search", notionSearch.SearchRequestSchema, notionSearch.SearchResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.results[0]["id"]).toBe("9007199254740993"); expect(r.results[0]["cover"]).toBeNull(); expect(r.results[1]["object"]).toBe("database"); expect(r.nextCursor).toBe("synthetic-next-token"); expect(r.hasMore).toBe(true); }],
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

// Per-dialect query.run fixtures: each engine's dialect enum, one
// engine-specific typed value (jsonb/timestamptz, bigint unsigned/datetime(6),
// nvarchar unicode/datetimeoffset, NUMBER(38) and VARIANT as strings, uuid),
// and the NULL cell. The >2^53 integers and precision-critical decimals travel
// as strings; NULL travels as the protobuf null kind.
const databaseDialects = [
  {
    name: "postgres",
    dialect: db.Dialect.POSTGRES,
    col2Type: "timestamptz",
    checks: [
      { row: 0, col: 0, want: "9007199254740993" },
      { row: 0, col: 1, want: '{"owner": "synthetic", "tags": ["ops"]}' },
      { row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00" },
      { row: 0, col: 3, want: "12345678901234567890.123456789012345678" },
      { row: 1, col: 1, want: null },
    ],
  },
  {
    name: "mysql",
    dialect: db.Dialect.MYSQL,
    col2Type: "datetime(6)",
    checks: [
      { row: 0, col: 0, want: "18446744073709551615" },
      { row: 0, col: 1, want: '{"owner": "synthetic", "tags": ["ops"]}' },
      { row: 0, col: 2, want: "2026-01-15 10:30:00.123456" },
      { row: 0, col: 3, want: "12345678901234567890.123456" },
      { row: 1, col: 1, want: null },
    ],
  },
  {
    name: "mssql",
    dialect: db.Dialect.MSSQL,
    col2Type: "datetimeoffset",
    checks: [
      { row: 0, col: 0, want: "9007199254740993" },
      { row: 0, col: 1, want: "München — ünïcodé ✓" },
      { row: 0, col: 2, want: "2026-01-15T10:30:00.1234567+02:00" },
      { row: 0, col: 3, want: "12345678901234567890.1234567890" },
      { row: 1, col: 2, want: null },
    ],
  },
  {
    name: "snowflake",
    dialect: db.Dialect.SNOWFLAKE,
    col2Type: "NUMBER(38,18)",
    checks: [
      { row: 0, col: 0, want: "99999999999999999999999999999999999999" },
      { row: 0, col: 1, want: "Big-ID Account — ünïcodé ✓" },
      { row: 0, col: 2, want: "12345678901234567890.123456789012345678" },
      { row: 0, col: 3, want: '{"owner": "synthetic", "tags": ["ops"]}' },
      { row: 1, col: 3, want: null },
    ],
  },
  {
    name: "supabase",
    dialect: db.Dialect.SUPABASE,
    col2Type: "timestamptz",
    checks: [
      { row: 0, col: 0, want: "6f9619ff-8b86-d011-b42d-00c04fc964ff" },
      { row: 0, col: 1, want: '{"owner": "synthetic", "tags": ["ops"]}' },
      { row: 0, col: 2, want: "2026-01-15T10:30:00.123456+00:00" },
      { row: 0, col: 3, want: "12345678901234567890.123456789012345678" },
      { row: 1, col: 1, want: null },
    ],
  },
] as const;

describe.each(databaseDialects)("database.query.run.$name", (tc) => {
  it("parses the request with the right dialect enum", () => {
    const req = fromJson(db.RunQueryRequestSchema, fixture(`database.query.run.${tc.name}`, "request"), { ignoreUnknownFields: true });
    expect(req.dialect).toBe(tc.dialect);
  });
  it("unmarshals one typed value per engine", () => {
    const resp = fromJson(db.RunQueryResponseSchema, fixture(`database.query.run.${tc.name}`, "response"), { ignoreUnknownFields: true });
    expect(resp.columns).toHaveLength(4);
    expect(resp.columns[2].type).toBe(tc.col2Type);
    expect(resp.rows).toHaveLength(2);
    for (const c of tc.checks) {
      const value = resp.rows[c.row].values[c.col];
      if (c.want === null) {
        expect(value.kind.case).toBe("nullValue");
      } else {
        expect(value.kind.case).toBe("stringValue");
        if (value.kind.case === "stringValue") expect(value.kind.value).toBe(c.want);
      }
    }
    expect(resp.truncated).toBe(false);
  });
  it("rejects the unknown extra field in strict mode", () => {
    expect(() => fromJson(db.RunQueryResponseSchema, fixture(`database.query.run.${tc.name}`, "response"))).toThrow();
  });
});

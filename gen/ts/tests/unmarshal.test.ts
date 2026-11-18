import { readFileSync } from "node:fs";
import { fromJson, type DescMessage, type JsonValue } from "@bufbuild/protobuf";
import { describe, expect, it } from "vitest";
import { isInteger, parse } from "lossless-json";
import * as atlassianJira from "../src/kei/connectors/atlassian/v1/jira_pb.js";
import * as atlassianConfluence from "../src/kei/connectors/atlassian/v1/confluence_pb.js";
import * as ghIssues from "../src/kei/connectors/github/v1/issues_pb.js";
import * as ghContents from "../src/kei/connectors/github/v1/contents_pb.js";
import * as ghPulls from "../src/kei/connectors/github/v1/pulls_pb.js";
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
  ["atlassian.jira.issues.get", atlassianJira.GetJiraIssueRequestSchema, atlassianJira.GetJiraIssueResponseSchema, (r: any) => { expect(r.key).toBe("SYN-101"); expect(r.fields.summary).toBe("Synthetic Jira issue"); expect(r.fields.description).toBeUndefined(); }],
  ["atlassian.jira.issues.search", atlassianJira.SearchJiraIssuesRequestSchema, atlassianJira.SearchJiraIssuesResponseSchema, (r: any) => { expect(r.issues).toHaveLength(2); expect(r.nextPageToken).toBe("synthetic-next-token"); expect(r.issues[1].fields.description).toBeUndefined(); }],
  ["atlassian.confluence.pages.get", atlassianConfluence.GetConfluencePageRequestSchema, atlassianConfluence.GetConfluencePageResponseSchema, (r: any) => { expect(r.id).toBe(78643201n); expect(r.title).toBe("Synthetic Confluence page"); expect(r.parentId).toBeUndefined(); }],
  ["atlassian.confluence.pages.list", atlassianConfluence.ListConfluencePagesRequestSchema, atlassianConfluence.ListConfluencePagesResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.results[0].title).toBe("Synthetic first page"); }],
  ["github.issues.get", ghIssues.GetIssueRequestSchema, ghIssues.GetIssueResponseSchema, (r: any) => { expect(r.id).toBe(9007199254740993n); expect(r.number).toBe(1347n); }],
  ["github.issues.list", ghIssues.ListIssuesRequestSchema, ghIssues.ListIssuesResponseSchema, (r: any) => expect(r.issues).toHaveLength(2)],
  ["github.contents.put", ghContents.PutContentsRequestSchema, ghContents.PutContentsResponseSchema, (r: any) => expect(r.content.name).toBe("notes.txt")],
  ["github.pulls.get", ghPulls.GetPullRequestRequestSchema, ghPulls.GetPullRequestResponseSchema, (r: any) => expect(r.number).toBe(42n)],
  ["linear.issues.get", linear.GetIssueRequestSchema, linear.GetIssueResponseSchema, (r: any) => { expect(r.issue.identifier).toBe("ENG-101"); expect(r.issue.labels[0].name).toBe("synthetic-label"); }],
  ["linear.issues.create", linear.CreateIssueRequestSchema, linear.CreateIssueResponseSchema, (r: any) => expect(r.success).toBe(true)],
  ["linear.issues.search", linear.SearchIssuesRequestSchema, linear.SearchIssuesResponseSchema, (r: any) => expect(r.nodes).toHaveLength(2)],
  ["linear.teams.get", linearTeams.GetTeamRequestSchema, linearTeams.GetTeamResponseSchema, (r: any) => expect(r.team.key).toBe("ENG")],
  ["linear.comments.create", linearComments.CreateCommentRequestSchema, linearComments.CreateCommentResponseSchema, (r: any) => { expect(r.success).toBe(true); expect(r.comment.body).toBe("Synthetic comment"); }],
  ["airtable.records.get", airtableRecords.GetRecordRequestSchema, airtableRecords.GetRecordResponseSchema, (r: any) => { expect(r.id).toBe("recSynthetic00003"); expect(r.fields["Name"]).toBe("Synthetic third task"); expect(r.fields["Priority"]).toBeNull(); }],
  ["airtable.records.list", airtableRecords.ListRecordsRequestSchema, airtableRecords.ListRecordsResponseSchema, (r: any) => { expect(r.records).toHaveLength(2); expect(r.offset).toBe("synthetic-next-offset"); expect(r.records[1].fields["Priority"]).toBeNull(); }],
  ["hubspot.contacts.get", hubspotContacts.GetContactRequestSchema, hubspotContacts.GetContactResponseSchema, (r: any) => { expect(r.id).toBe("701"); expect(r.properties.email).toBe("synthetic.contact@example.invalid"); expect(r.properties.lastname).toBeNull(); }],
  ["hubspot.contacts.list", hubspotContacts.ListContactsRequestSchema, hubspotContacts.ListContactsResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.nextAfter).toBe("synthetic-next-cursor"); }],
  ["hubspot.deals.list", hubspotDeals.ListDealsRequestSchema, hubspotDeals.ListDealsResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.results[0].properties.dealname).toBe("Synthetic first deal"); }],
  ["notion.pages.get", notionPages.GetPageRequestSchema, notionPages.GetPageResponseSchema, (r: any) => { expect(r.id).toBe("5f8b7e4a-0000-4000-8000-000000000001"); expect(r.publicUrl).toBeUndefined(); }],
  ["notion.databases.query", notionDatabases.QueryDatabaseRequestSchema, notionDatabases.QueryDatabaseResponseSchema, (r: any) => { expect(r.results).toHaveLength(2); expect(r.nextCursor).toBe("synthetic-next-cursor"); expect(r.results[1].publicUrl).toBeUndefined(); }],
  ["sentry.issues.get", sentry.GetIssueRequestSchema, sentry.GetIssueResponseSchema, (r: any) => { expect(r.shortId).toBe("PROJ-103"); expect(r.shareId).toBeUndefined(); }],
  ["sentry.issues.list", sentry.ListIssuesRequestSchema, sentry.ListIssuesResponseSchema, (r: any) => { expect(r.issues).toHaveLength(2); expect(r.issues[1].shareId).toBe("synthetic-share-id-2"); }],
  ["http_api.records.list", http.ListRecordsRequestSchema, http.ListRecordsResponseSchema, (r: any) => expect(r.records).toHaveLength(2)],
  ["http_api.records.get", http.GetRecordRequestSchema, http.GetRecordResponseSchema, (r: any) => expect(r.record.id).toBe("rec-1")],
  ["database.query.run", db.RunQueryRequestSchema, db.RunQueryResponseSchema, (r: any) => expect(r.columns[0].name).toBe("id")],
  ["mcp.tool.call", mcp.ToolCallSchema, mcp.ToolResultSchema, (r: any) => expect(r.content).toHaveLength(2)],
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

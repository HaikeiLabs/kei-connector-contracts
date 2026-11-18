import { readFileSync } from "node:fs";
import { fromJson, type DescMessage, type JsonValue } from "@bufbuild/protobuf";
import { describe, expect, it } from "vitest";
import { isInteger, parse } from "lossless-json";
import * as gmailMessages from "../src/kei/connectors/gmail/v1/messages_pb.js";
import * as googleSheets from "../src/kei/connectors/google_sheets/v1/spreadsheets_pb.js";
import * as ghIssues from "../src/kei/connectors/github/v1/issues_pb.js";
import * as ghContents from "../src/kei/connectors/github/v1/contents_pb.js";
import * as ghPulls from "../src/kei/connectors/github/v1/pulls_pb.js";
import * as linear from "../src/kei/connectors/linear/v1/issues_pb.js";
import * as linearComments from "../src/kei/connectors/linear/v1/comments_pb.js";
import * as linearTeams from "../src/kei/connectors/linear/v1/teams_pb.js";
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
  ["gmail.messages.list", gmailMessages.ListMessagesRequestSchema, gmailMessages.ListMessagesResponseSchema, (r: any) => { expect(r.messages).toHaveLength(2); expect(r.nextPageToken).toBe("synthetic-next-token"); }],
  ["gmail.messages.get", gmailMessages.GetMessageRequestSchema, gmailMessages.GetMessageResponseSchema, (r: any) => { expect(r.id).toBe("msg-1001"); expect(r.snippet).toBe("Synthetic message preview text"); expect(r.internalDate).toBeUndefined(); }],
  ["google_sheets.values.get", googleSheets.GetValuesRequestSchema, googleSheets.GetValuesResponseSchema, (r: any) => { expect(r.range).toBe("Sheet1!A1:B2"); expect(r.values).toHaveLength(2); expect(r.values[0].values[1].kind.case).toBe("nullValue"); }],
  ["google_sheets.spreadsheets.get", googleSheets.GetSpreadsheetRequestSchema, googleSheets.GetSpreadsheetResponseSchema, (r: any) => { expect(r.spreadsheetId).toBe("synthetic-sheet-1"); expect(r.properties.title).toBe("Synthetic sheet"); }],
  ["github.issues.get", ghIssues.GetIssueRequestSchema, ghIssues.GetIssueResponseSchema, (r: any) => { expect(r.id).toBe(9007199254740993n); expect(r.number).toBe(1347n); }],
  ["github.issues.list", ghIssues.ListIssuesRequestSchema, ghIssues.ListIssuesResponseSchema, (r: any) => expect(r.issues).toHaveLength(2)],
  ["github.contents.put", ghContents.PutContentsRequestSchema, ghContents.PutContentsResponseSchema, (r: any) => expect(r.content.name).toBe("notes.txt")],
  ["github.pulls.get", ghPulls.GetPullRequestRequestSchema, ghPulls.GetPullRequestResponseSchema, (r: any) => expect(r.number).toBe(42n)],
  ["linear.issues.get", linear.GetIssueRequestSchema, linear.GetIssueResponseSchema, (r: any) => { expect(r.issue.identifier).toBe("ENG-101"); expect(r.issue.labels[0].name).toBe("synthetic-label"); }],
  ["linear.issues.create", linear.CreateIssueRequestSchema, linear.CreateIssueResponseSchema, (r: any) => expect(r.success).toBe(true)],
  ["linear.issues.search", linear.SearchIssuesRequestSchema, linear.SearchIssuesResponseSchema, (r: any) => expect(r.nodes).toHaveLength(2)],
  ["linear.teams.get", linearTeams.GetTeamRequestSchema, linearTeams.GetTeamResponseSchema, (r: any) => expect(r.team.key).toBe("ENG")],
  ["linear.comments.create", linearComments.CreateCommentRequestSchema, linearComments.CreateCommentResponseSchema, (r: any) => { expect(r.success).toBe(true); expect(r.comment.body).toBe("Synthetic comment"); }],
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

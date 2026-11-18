import { createServer } from 'node:http';
import { readFileSync } from 'node:fs';
import { once } from 'node:events';
import { fromJson } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { newHTTP, type SQL } from '../src/kei/connection/index.js';
import { GithubQueries } from '../src/kei/connectors/github/v1/ops.js';
import { DatabaseQueries } from '../src/kei/connectors/database/v1/ops.js';
import { GetIssueRequestSchema } from '../src/kei/connectors/github/v1/issues_pb.js';
import { RunQueryRequestSchema } from '../src/kei/connectors/database/v1/query_pb.js';

function fixture(path: string): string {
  return readFileSync(new URL(`../../../fixtures/operations/${path}.json`, import.meta.url), 'utf8');
}

describe('generated operation replay', () => {
  it('replays github.issues.get and preserves the large id', async () => {
    let recorded: { method?: string; path?: string; authorization?: string; body: string } | undefined;
    const server = createServer((req, res) => {
      recorded = { method: req.method, path: req.url, authorization: req.headers.authorization, body: '' };
      req.on('data', chunk => { recorded!.body += chunk; });
      req.on('end', () => { res.setHeader('Content-Type', 'application/json'); res.end(fixture('github/issues.get.response')); });
    });
    server.listen(0, '127.0.0.1'); await once(server, 'listening');
    try {
      const address = server.address(); if (!address || typeof address === 'string') throw new Error('missing address');
      const request = fromJson(GetIssueRequestSchema, JSON.parse(fixture('github/issues.get.request')));
      const result = await new GithubQueries(newHTTP(`http://127.0.0.1:${address.port}`, 'bearer', 'synthetic-token')).getIssue(request);
      expect(recorded).toEqual({ method: 'GET', path: '/repos/example-org/example-repo/issues/1347', authorization: 'Bearer synthetic-token', body: '' });
      expect(result.id).toBe(9007199254740993n);
    } finally { server.close(); }
  });
  it('replays database.query.run with a fake SQL driver', async () => {
    let seen: unknown[] = [];
    const conn: SQL = { async query(sql, args, maxRows) {
      seen = [sql, args, maxRows];
      return { columns: [{ name: 'id', type: 'BIGINT' }, { name: 'name', type: 'TEXT' }, { name: 'note', type: 'TEXT' }, { name: 'balance', type: 'NUMERIC' }], rows: [[9007199254740993n, 'Big-ID Account — ünïcodé ✓', null, '12345678901234567890.123456789012345678'], [1n, 'Synthetic One', 'Synthetic note', '0.000000000000000001']], truncated: false };
    } };
    const request = fromJson(RunQueryRequestSchema, JSON.parse(fixture('database/query.run.request')));
    const result = await new DatabaseQueries(conn).runQuery(request);
    expect(seen).toEqual([request.sql, ['9007199254740993'], 100n]);
    expect(result.columns[0].name).toBe('id');
    expect(result.rows).toHaveLength(2);
    expect(result.rows[0].values[0].kind).toEqual({ case: 'stringValue', value: '9007199254740993' });
  });
});

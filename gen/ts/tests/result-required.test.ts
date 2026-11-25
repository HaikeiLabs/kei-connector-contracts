// Result-required checks: a write operation whose result_required paths are
// not all held by a 2xx response fails with a contract_violation error, not
// with a silent success. Mirrors gen/go/conformance/result_required_test.go.
import { createServer } from 'node:http';
import { readFileSync } from 'node:fs';
import { once } from 'node:events';
import { fromJson } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { ConnectionError, newHTTP } from '../src/kei/connection/index.js';
import { LinearQueries } from '../src/kei/connectors/linear/v1/ops.js';
import { CreateIssueRequestSchema } from '../src/kei/connectors/linear/v1/issues_pb.js';

function fixture(path: string): string {
  return readFileSync(new URL(`../../../fixtures/operations/${path}.json`, import.meta.url), 'utf8');
}

const badFixtures: [name: string, path: string][] = [
  ['success false', 'linear/issues.create.response.success-false'],
  ['issue missing', 'linear/issues.create.response.no-issue'],
  ['empty identifier', 'linear/issues.create.response.empty-identifier'],
];

describe('linear.issues.create result_required', () => {
  for (const [name, path] of badFixtures) {
    it(`rejects ${name} with a contract violation`, async () => {
      const envelope = JSON.stringify({ data: { issueCreate: JSON.parse(fixture(path)) } });
      const server = createServer((req, res) => {
        res.setHeader('Content-Type', 'application/json');
        res.end(envelope);
      });
      server.listen(0, '127.0.0.1');
      await once(server, 'listening');
      try {
        const address = server.address();
        if (!address || typeof address === 'string') throw new Error('missing address');
        const request = fromJson(CreateIssueRequestSchema, JSON.parse(fixture('linear/issues.create.request')));
        const queries = new LinearQueries(newHTTP(`http://127.0.0.1:${address.port}`, 'bearer', 'synthetic-token'));
        const error = await queries.createIssue(request).then(
          () => { throw new Error('createIssue succeeded'); },
          (e: unknown) => e,
        );
        expect(error).toBeInstanceOf(ConnectionError);
        const connectionError = error as ConnectionError;
        expect(connectionError.kind).toBe('contract_violation');
        expect(connectionError.opId).toBe('linear.issues.create');
        expect(connectionError.status).toBe(200);
      } finally { server.close(); }
    });
  }

  it('keeps the good fixture passing', async () => {
    const envelope = JSON.stringify({ data: { issueCreate: JSON.parse(fixture('linear/issues.create.response')) } });
    const server = createServer((req, res) => {
      res.setHeader('Content-Type', 'application/json');
      res.end(envelope);
    });
    server.listen(0, '127.0.0.1');
    await once(server, 'listening');
    try {
      const address = server.address();
      if (!address || typeof address === 'string') throw new Error('missing address');
      const request = fromJson(CreateIssueRequestSchema, JSON.parse(fixture('linear/issues.create.request')));
      const queries = new LinearQueries(newHTTP(`http://127.0.0.1:${address.port}`, 'bearer', 'synthetic-token'));
      const result = await queries.createIssue(request);
      expect(result.success).toBe(true);
      expect(result.issue?.identifier).toBe('ENG-102');
    } finally { server.close(); }
  });
});

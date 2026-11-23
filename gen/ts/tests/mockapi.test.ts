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

// Smoke tests: the generated TypeScript operations run against the real
// contract-derived mock server (testenv/mockapi) — fixture responses,
// GraphQL root-field routing, and failure injection via the X-Mock-Fail
// header mapping to the generated error kinds.
//
// The mock is the Go binary from testenv/mockapi/cmd/mockapi, built once
// with the local Go toolchain (or taken from MOCKAPI_BIN) and started on a
// free port. Skipped where no Go toolchain is available.
import { spawn, spawnSync, type ChildProcess } from 'node:child_process';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { fromJson } from '@bufbuild/protobuf';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { newHTTP, type HTTP, type HTTPResponse } from '../src/kei/connection/index.js';
import { GithubQueries } from '../src/kei/connectors/github/v1/ops.js';
import { LinearQueries } from '../src/kei/connectors/linear/v1/ops.js';
import { GetIssueRequestSchema } from '../src/kei/connectors/github/v1/issues_pb.js';
import { GetIssueRequestSchema as LinearGetIssueRequestSchema } from '../src/kei/connectors/linear/v1/issues_pb.js';

const repoRoot = fileURLToPath(new URL('../../../', import.meta.url));
const token = 'smoke-token';
const hasGo = Boolean(process.env.MOCKAPI_BIN) || spawnSync('go', ['version'], { stdio: 'ignore' }).status === 0;

async function startMock(): Promise<{ proc: ChildProcess; base: string }> {
  let dir = '';
  const bin =
    process.env.MOCKAPI_BIN ??
    (() => {
      dir = mkdtempSync(path.join(tmpdir(), 'mockapi-smoke-'));
      const binary = path.join(dir, 'mockapi');
      const build = spawnSync('go', ['-C', 'testenv/mockapi', 'build', '-o', binary, './cmd/mockapi'], {
        cwd: repoRoot,
        stdio: 'inherit',
      });
      if (build.status !== 0) throw new Error('go build of the mock failed');
      return binary;
    })();
  const proc = spawn(bin, [
    '-addr',
    '127.0.0.1:0',
    '-manifest',
    path.join(repoRoot, 'manifest/connector-manifest.v1.json'),
    '-fixtures',
    path.join(repoRoot, 'fixtures/operations'),
  ]);
  const address = new Promise<string>((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error('mockapi did not start')), 30000);
    proc.stderr!.on('data', (chunk: Buffer) => {
      const match = chunk.toString().match(/serving \d+ routes on (\S+)/);
      if (match) {
        clearTimeout(timer);
        resolve(match[1]);
      }
    });
    proc.on('exit', code => {
      clearTimeout(timer);
      reject(new Error(`mockapi exited before listening (code ${code})`));
    });
  });
  return { proc, base: `http://${await address}` };
}

// A second HTTP transport for the failure-injection case: newHTTP builds
// its headers internally, so the injected header needs its own transport.
// It still exercises the real client decode path end to end.
function injectingHTTP(base: string, extra: Record<string, string>): HTTP {
  return {
    async do(method, target, query, body) {
      const search = query.size ? `?${query}` : '';
      const res = await fetch(base + target + search, {
        method,
        headers: { Accept: 'application/json', Authorization: `Bearer ${token}`, ...extra },
        body,
      });
      return { status: res.status, headers: res.headers, body: await res.text() } satisfies HTTPResponse;
    },
  };
}

describe.skipIf(!hasGo)('generated operations against the contract mock', () => {
  let mock: { proc: ChildProcess; base: string } | undefined;
  let dir = '';

  // The go build of the mock can exceed vitest's 10s hook default on the
  // release runner; give it a generous budget.
  beforeAll(
    async () => {
      mock = await startMock();
      dir = path.dirname(mock.proc.spawnfile);
    },
    180000,
  );
  afterAll(() => {
    mock?.proc.kill();
    if (!process.env.MOCKAPI_BIN && dir.startsWith(tmpdir())) rmSync(dir, { recursive: true, force: true });
  });

  it('replays github.issues.get with the lossless id', async () => {
    const queries = new GithubQueries(newHTTP(mock!.base, 'bearer', token));
    const resp = await queries.getIssue(
      fromJson(GetIssueRequestSchema, { owner: 'acme', repository: 'widgets', issueNumber: 1347 }),
    );
    expect(resp.number).toBe(1347n);
    expect(resp.id).toBe(9007199254740993n);
    expect(resp.title).toBe('Synthetic pilot issue');
  });

  it('replays linear.issues.get over GraphQL', async () => {
    const queries = new LinearQueries(newHTTP(mock!.base, 'bearer', token));
    const resp = await queries.getIssue(
      fromJson(LinearGetIssueRequestSchema, { id: '00000000-0000-0000-0000-000000000101' }),
    );
    expect(resp.issue?.identifier).toBe('ENG-101');
  });

  it('maps an injected 401 to unauthorized', async () => {
    const queries = new GithubQueries(injectingHTTP(mock!.base, { 'X-Mock-Fail': '401' }));
    await expect(
      queries.getIssue(fromJson(GetIssueRequestSchema, { owner: 'acme', repository: 'widgets', issueNumber: 1347 })),
    ).rejects.toMatchObject({ opId: 'github.issues.get', kind: 'unauthorized', status: 401 });
  });
});

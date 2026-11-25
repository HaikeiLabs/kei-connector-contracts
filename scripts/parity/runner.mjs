// Copyright 2026 Haikei Labs
// SPDX-License-Identifier: Apache-2.0
//
// TypeScript runner of the cross-language parity test (protocol: see
// scripts/parity/run.sh). It runs the built package in gen/ts/dist and calls
// each operation by its generated method name, so a new operation or
// connector needs no change here.
import { readFileSync, readdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..', '..');
const dist = join(root, 'gen', 'ts', 'dist', 'kei');
const req = (p) => import(pathToFileURL(p).href);

const [casesPath, base, token] = process.argv.slice(2);
if (!token) {
  console.error('usage: runner.mjs <cases.json> <base url> <token>');
  process.exit(2);
}

// Resolve the protobuf runtime from gen/ts, where it is installed.
const { createRegistry, fromJson, toJson } = await req(join(root, 'gen', 'ts', 'node_modules', '@bufbuild', 'protobuf', 'dist', 'esm', 'index.js'));
const { isInteger, parse: losslessParseRaw } = await req(join(root, 'gen', 'ts', 'node_modules', 'lossless-json', 'lib', 'esm', 'index.js'));
const { newHTTP } = await req(join(dist, 'connection', 'index.js'));

// Parse with integers above the double-precision range kept as strings, so
// int64 values survive; fromJson accepts both shapes for 64-bit fields.
const losslessParse = (raw) => losslessParseRaw(raw, undefined, { parseNumber: (v) => isInteger(v) && (BigInt(v) > BigInt(Number.MAX_SAFE_INTEGER) || BigInt(v) < BigInt(Number.MIN_SAFE_INTEGER)) ? v : Number(v) });

const connectors = new Map();
async function connector(id) {
  if (connectors.has(id)) return connectors.get(id);
  const dir = join(dist, 'connectors', id, 'v1');
  const files = [];
  for (const name of readdirSync(dir).filter((n) => n.endsWith('_pb.js')).sort()) {
    for (const v of Object.values(await req(join(dir, name)))) if (v && v.kind === 'file') files.push(v);
  }
  const ops = await req(join(dir, 'ops.js'));
  const Queries = Object.entries(ops).find(([name, v]) => name.endsWith('Queries') && typeof v === 'function')?.[1];
  if (!Queries) throw new Error(`no Queries class for connector ${id}`);
  const entry = { registry: createRegistry(...files), Queries };
  connectors.set(id, entry);
  return entry;
}

const lcfirst = (s) => s.charAt(0).toLowerCase() + s.slice(1);

for (const c of losslessParse(readFileSync(casesPath, 'utf8'))) {
  let line;
  try {
    const { registry, Queries } = await connector(c.connector);
    const reqSchema = registry.getMessage(c.request_type);
    const respSchema = registry.getMessage(c.response_type);
    if (!reqSchema || !respSchema) throw new Error(`unknown message ${c.request_type} or ${c.response_type}`);
    const queries = new Queries(newHTTP(`${base}/${c.op}`, 'bearer', token));
    const method = queries[lcfirst(c.method)];
    if (typeof method !== 'function') throw new Error(`${c.connector} has no method ${lcfirst(c.method)}`);
    const resp = await method.call(queries, fromJson(reqSchema, c.request));
    line = { op: c.op, response: toJson(respSchema, resp) };
  } catch (e) {
    line = { op: c.op, error: String(e && e.message ? e.message : e) };
  }
  process.stdout.write(JSON.stringify(line) + '\n');
}

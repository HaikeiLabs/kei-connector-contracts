import { fromJson, toJson, type DescField, type DescMessage, type JsonValue, type MessageShape } from '@bufbuild/protobuf';
import { isInteger, LosslessNumber, parse, stringify } from 'lossless-json';

/** Render path variables without allowing a value to change path structure. */
export function renderPath(template: string, values: Record<string, string>, multiSegment: string[] = []): string {
  if (!template.startsWith('/')) throw new Error('connection: path must start with /');
  let result = '';
  let rest = template;
  while (rest.length) {
    const open = rest.indexOf('{');
    const close = rest.indexOf('}');
    if (open < 0) {
      if (close >= 0) throw new Error('connection: unmatched }');
      return result + rest;
    }
    if (close < open) throw new Error('connection: invalid path variable');
    const name = rest.slice(open + 1, close);
    if (!name || /[{/]/.test(name)) throw new Error('connection: invalid path variable');
    const value = values[name];
    if (!value) throw new Error(`connection: missing or empty path variable ${name}`);
    const escape = (v: string) => encodeURIComponent(v).replace(/%3A/gi, ':').replace(/%40/gi, '@').replace(/%24/gi, '$').replace(/%26/gi, '&').replace(/%2B/gi, '+').replace(/%2C/gi, ',').replace(/%3B/gi, ';').replace(/%3D/gi, '=');
    let encoded: string;
    if (multiSegment.includes(name)) {
      const parts = value.split('/');
      if (parts.some(part => !part || part === '.' || part === '..')) throw new Error(`connection: invalid path variable ${name}`);
      encoded = parts.map(escape).join('/');
    } else {
      encoded = /^\.+$/.test(value) ? '%2E'.repeat(value.length) : escape(value);
    }
    result += rest.slice(0, open) + encoded;
    rest = rest.slice(close + 1);
  }
  return result;
}

export interface ErrorMapping { status: number; kind: string; retryable?: boolean }
export class ConnectionError extends Error {
  constructor(readonly opId: string, readonly status: number, readonly kind: string, readonly retryable = false) {
    super(`connection: ${opId}: ${kind} (status ${status}${retryable ? ', retryable' : ''})`);
  }
}
export function mapError(opId: string, status: number, table: ErrorMapping[]): ConnectionError | undefined {
  if (status >= 200 && status < 300) return undefined;
  const mapping = table.find(entry => entry.status === status);
  if (mapping) return new ConnectionError(opId, status, mapping.kind, mapping.retryable);
  const kind = status === 400 || status === 422 ? 'invalid' : status === 401 ? 'unauthorized' : status === 403 ? 'forbidden' : status === 404 ? 'not_found' : status === 409 ? 'conflict' : status === 429 ? 'rate_limited' : status >= 500 && status < 600 ? 'unavailable' : 'unspecified';
  return new ConnectionError(opId, status, kind, [429, 502, 503, 504].includes(status));
}

export interface Op { id: string; class: 'read' | 'write'; destructive: boolean }
export interface Param { field: string; in: 'path' | 'query' | 'body'; wire: string; multiSegment: boolean }
export interface HTTPSpec { op: Op; method: string; pathTemplate: string; params: Param[]; errors: ErrorMapping[]; graphqlDocument: string; cursorField: string; nextCursorField: string }
export interface HTTPResponse { status: number; headers: Headers; body: string }
export interface HTTP { do(method: string, path: string, query: URLSearchParams, body?: string): Promise<HTTPResponse> }
export interface SQL { query(sql: string, args: JsonValue[], maxRows: bigint): Promise<{ columns: { name: string; type: string }[]; rows: unknown[][]; truncated: boolean }> }

export function newHTTP(baseURL: string, kind: 'bearer' | 'header' | 'basic', credential: string, client: typeof fetch = fetch, headerName = 'X-Api-Key'): HTTP {
  const base = new URL(baseURL);
  if (!['http:', 'https:'].includes(base.protocol) || base.username || base.password || base.search || base.hash) throw new Error('connection: invalid base URL');
  if (!credential || /[\r\n\0]/.test(credential)) throw new Error('connection: invalid credential');
  if (kind === 'basic' && (!credential.includes(':') || credential.startsWith(':'))) throw new Error('connection: invalid basic credential');
  if (kind === 'header' && !/^[!#$%&'*+.^_`|~\w-]+$/.test(headerName)) throw new Error('connection: invalid header name');
  const authorization = kind === 'basic' ? `Basic ${btoa(credential)}` : kind === 'bearer' ? `Bearer ${credential}` : credential;
  return { async do(method, path, query, body) {
    if (!path.startsWith('/') || path.startsWith('//') || /[?#]/.test(path)) throw new Error('connection: invalid path');
    const headers = new Headers({ Accept: 'application/json' });
    headers.set(kind === 'header' ? headerName : 'Authorization', authorization);
    if (body !== undefined) headers.set('Content-Type', 'application/json');
    const target = base.toString().replace(/\/$/, '') + path + (query.size ? `?${query}` : '');
    try {
      const response = await client(target, { method, headers, body, redirect: 'error' });
      const text = await response.text();
      if (new TextEncoder().encode(text).length > 32 * 1024 * 1024) throw new Error('too large');
      return { status: response.status, headers: response.headers, body: text };
    } catch { throw new Error('connection: transport failed'); }
  } };
}

function field(schema: DescMessage, name: string): DescField {
  const found = schema.fields.find(f => f.name === name);
  if (!found) throw new Error('connection: invalid operation specification');
  return found;
}
function canonical(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(canonical).join(',')}]`;
  if (value !== null && typeof value === 'object' && !(value instanceof LosslessNumber)) return `{${Object.entries(value).sort(([a], [b]) => a < b ? -1 : a > b ? 1 : 0).map(([k,v]) => `${JSON.stringify(k)}:${canonical(v)}`).join(',')}}`;
  return stringify(value) ?? 'null';
}
function wireObject(schema: DescMessage, value: Record<string, unknown>, graphql: boolean): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  for (const f of schema.fields) {
    if (!(f.jsonName in value)) continue;
    const defaultName = f.name.replace(/_([a-z])/g, (_, c: string) => c.toUpperCase());
    const name = graphql || f.jsonName !== defaultName ? f.jsonName : f.name;
    const convert = (v: unknown): unknown => {
      if (f.message && v && typeof v === 'object' && !f.message.typeName.startsWith('google.protobuf.')) return wireObject(f.message, v as Record<string, unknown>, graphql);
      if (f.scalar && [3,4,6,16,18].includes(f.scalar) && typeof v === 'string') return new LosslessNumber(v);
      return v;
    };
    const v = value[f.jsonName];
    result[name] = Array.isArray(v) ? v.map(convert) : convert(v);
  }
  return result;
}
function lossless(raw: string): JsonValue {
  return parse(raw, undefined, { parseNumber: v => isInteger(v) && (BigInt(v) > BigInt(Number.MAX_SAFE_INTEGER) || BigInt(v) < BigInt(Number.MIN_SAFE_INTEGER)) ? v : Number(v) }) as JsonValue;
}
export async function callHTTP<Req extends DescMessage, Resp extends DescMessage>(conn: HTTP, spec: HTTPSpec, req: MessageShape<Req>, request: Req, response: Resp): Promise<MessageShape<Resp>> {
  const json = toJson(request, req) as Record<string, unknown>;
  const pathValues: Record<string, string> = {}, multi: string[] = [], query = new URLSearchParams(), body: Record<string, unknown> = {};
  let hasBody = !!spec.graphqlDocument;
  for (const p of spec.params) {
    const f = field(request, p.field), v = json[f.jsonName];
    if (p.in === 'path') { pathValues[p.wire] = v === undefined ? '' : String(v); if (p.multiSegment) multi.push(p.wire); }
    else if (p.in === 'query' && v !== undefined) for (const item of Array.isArray(v) ? v : [v]) query.append(p.wire, String(item));
    else if (p.in === 'body') {
      hasBody = true;
      if (v !== undefined) { const encoded = wireObject(request, { [f.jsonName]: v }, !!spec.graphqlDocument); body[p.wire] = Object.values(encoded)[0]; }
    }
  }
  const path = renderPath(spec.pathTemplate, pathValues, multi);
  const payload = hasBody ? canonical(spec.graphqlDocument ? { query: spec.graphqlDocument, variables: body } : body) : undefined;
  const res = await conn.do(spec.method, path, query, payload);
  const error = mapError(spec.op.id, res.status, spec.errors); if (error) throw error;
  try {
    let data: any = res.body ? lossless(res.body) : {};
    if (spec.graphqlDocument) {
      const code = data.errors?.[0]?.extensions?.code;
      if (data.errors?.length) { const kinds: Record<string, string> = { NOT_FOUND:'not_found', UNAUTHENTICATED:'unauthorized', FORBIDDEN:'forbidden', RATELIMITED:'rate_limited' }; throw new ConnectionError(spec.op.id, res.status, kinds[code] ?? 'invalid', code === 'RATELIMITED'); }
      data = data.data ?? {};
      const keys = Object.keys(data);
      if (keys.length === 1 && !response.fields.some(f => f.name === keys[0] || f.jsonName === keys[0])) data = data[keys[0]];
    }
    if (Array.isArray(data)) {
      const lists = response.fields.filter(f => f.fieldKind === 'list' && f.name !== spec.nextCursorField);
      if (lists.length !== 1) throw new Error('invalid array');
      data = { [lists[0].jsonName]: data };
    }
    if (spec.cursorField && spec.nextCursorField) {
      const next = field(response, spec.nextCursorField), cursor = spec.params.find(p => p.field === spec.cursorField)?.wire ?? spec.cursorField;
      const link = res.headers.get('link')?.match(/<([^>]+)>\s*;\s*rel="?next"?/i)?.[1];
      if (link && data[next.jsonName] === undefined) { const value = new URL(link).searchParams.get(cursor); if (value) data[next.jsonName] = next.scalar === 9 ? value : Number(value); }
    }
    return fromJson(response, data, { ignoreUnknownFields: true });
  } catch (e) { if (e instanceof ConnectionError) throw e; throw new Error(`connection: ${spec.op.id}: response does not match the contract`); }
}
export async function callSQL<Resp extends DescMessage>(conn: SQL, op: Op, req: { sql: string; params: any[]; maxRows?: bigint }, response: Resp): Promise<MessageShape<Resp>> {
  const { ValueSchema } = await import('@bufbuild/protobuf/wkt');
  const args = req.params.map(p => toJson(ValueSchema, p));
  const result = await conn.query(req.sql, args, req.maxRows && req.maxRows > 0n ? req.maxRows : 1000n);
  const cell = (v: unknown): JsonValue => typeof v === 'bigint' ? (v > 9007199254740992n || v < -9007199254740992n ? String(v) : Number(v)) : v instanceof Uint8Array ? new TextDecoder().decode(v) : v instanceof Date ? v.toISOString() : v as JsonValue;
  return fromJson(response, { columns: result.columns.map(c => ({ name: c.name, type: c.type.toLowerCase() })), rows: result.rows.map(row => row.map(cell)), truncated: result.truncated }, { ignoreUnknownFields: true });
}

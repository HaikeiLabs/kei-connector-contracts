import { describe, expect, it } from 'vitest';
import { createServer } from 'node:http';
import { once } from 'node:events';
import { renderPath, mapError, newHTTP } from '../src/kei/connection/index.js';

describe('renderPath', () => {
  it('escapes segments', () => {
    expect(renderPath('/repos/{owner}/{repo}', { owner: 'a/b', repo: 'x y' })).toBe('/repos/a%2Fb/x%20y');
    expect(renderPath('/{a}/{b}', { a: '..', b: 'é?#%' })).toBe('/%2E%2E/%C3%A9%3F%23%25');
    expect(() => renderPath('/{missing}', {})).toThrow();
    // Go's url.PathEscape set: keep $ & + : = @, escape ! ' ( ) * , ;
    expect(renderPath('/x/{v}', { v: "a!b'c(d)e*f,g;h$i&j+k:l=m@n" })).toBe("/x/a%21b%27c%28d%29e%2Af%2Cg%3Bh$i&j+k:l=m@n");
    expect(renderPath('/x/{v}', { v: "d/a!b" }, ['v'])).toBe('/x/d/a%21b');
  });
  it('handles multi segment paths', () => {
    expect(renderPath('/{owner}/{path}', { owner: 'a/b', path: 'docs/a b' }, ['path'])).toBe('/a%2Fb/docs/a%20b');
    for (const path of ['a/../b', 'a//b', './a', '/a', 'a/', '']) {
      expect(() => renderPath('/{path}', { path }, ['path'])).toThrow();
    }
  });
});

it('sends an escaped dot-only segment unchanged on the wire', async () => {
  const path = renderPath('/x/{a}', { a: '..' });
  expect(path).toBe('/x/%2E%2E');
  let received = '';
  const server = createServer((req, res) => { received = req.url ?? ''; res.end('{}'); });
  server.listen(0, '127.0.0.1');
  await once(server, 'listening');
  try {
    const address = server.address();
    if (!address || typeof address === 'string') throw new Error('missing address');
    await newHTTP(`http://127.0.0.1:${address.port}`, 'bearer', 'synthetic-token').do('GET', path, new URLSearchParams());
    expect(received).toBe('/x/%2E%2E');
  } finally { server.close(); }
});

it('does not include credentials or response bodies in errors', () => {
  for (const status of [401, 500]) {
    const error = mapError('github.issues.get', status, []);
    expect(String(error)).not.toContain('synthetic-secret');
    expect(String(error)).not.toContain('bad token');
    expect(String(error)).toContain(String(status));
  }
});

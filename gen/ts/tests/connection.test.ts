import { describe, expect, it } from 'vitest';
import { renderPath, mapError } from '../src/kei/connection/index.js';

describe('renderPath', () => {
  it('escapes segments', () => {
    expect(renderPath('/repos/{owner}/{repo}', { owner: 'a/b', repo: 'x y' })).toBe('/repos/a%2Fb/x%20y');
    expect(renderPath('/{a}/{b}', { a: '..', b: 'é?#%' })).toBe('/%2E%2E/%C3%A9%3F%23%25');
    expect(() => renderPath('/{missing}', {})).toThrow();
  });
  it('handles multi segment paths', () => {
    expect(renderPath('/{owner}/{path}', { owner: 'a/b', path: 'docs/a b' }, ['path'])).toBe('/a%2Fb/docs/a%20b');
    for (const path of ['a/../b', 'a//b', './a', '/a', 'a/', '']) {
      expect(() => renderPath('/{path}', { path }, ['path'])).toThrow();
    }
  });
});

it('does not include credentials or response bodies in errors', () => {
  for (const status of [401, 500]) {
    const error = mapError('github.issues.get', status, []);
    expect(String(error)).not.toContain('synthetic-secret');
    expect(String(error)).not.toContain('bad token');
    expect(String(error)).toContain(String(status));
  }
});

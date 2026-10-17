import assert from 'node:assert/strict';
import { test } from 'node:test';

import { buildRequest, MAX_BODY_BYTES, type Libretto, type Request, type Server } from '../transport/transport.ts';
import { ExpressServer } from './express/express-server.ts';
import { matchRoutePattern } from './reading.ts';
import { StdlibServer } from './stdlib.ts';

const EVERY_ADAPTER: [string, Server][] = [['stdlib', new StdlibServer()], ['express', new ExpressServer()]];
const ANY_FREE_PORT = 0;

/** Every Framework is Held to one Contract: the same Path, Query, Headers and Body Arrive, and the same Answer Leaves. */
for (const [name, adapter] of EVERY_ADAPTER) {
  test(`${name}AnswersTheSharedContract`, async () => {
    let seen: Request = buildRequest();
    const libretto: Libretto = {
      routes: [
        { method: 'POST', pattern: '/echo/{id}', handle: async (request) => {
          seen = request;

          return { status: 201, body: { id: request.path['id'] } };
        } },
        { method: 'GET', pattern: '/echo', handle: async () => ({ status: 200, body: [] }) },
        { method: 'DELETE', pattern: '/echo/{id}', handle: async () => ({ status: 204, body: undefined }) },
        { method: 'GET', pattern: '/failing', handle: async () => ({ status: 404, body: { error: 'gone' } }) },
      ],
      unknown: async () => ({ status: 404, body: { error: 'route not Found' } }),
    };

    const running = await adapter.serveRoutes(libretto, ANY_FREE_PORT);
    const at = (target: string) => `http://localhost:${running.port}${target}`;
    try {
      const echoed = await fetch(at('/echo/007?page=2&size=5&page=9'), {
        method: 'POST', body: '{"a":1}', headers: { Authorization: 'bEaReR abc', 'Content-Type': 'text/plain' },
      });
      assert.equal(echoed.status, 201);
      assert.equal(echoed.headers.get('content-type'), 'application/json');
      assert.equal(await echoed.text(), '{"id":"007"}');
      assert.equal(seen.path['id'], '007');
      assert.deepEqual(seen.query, { page: '2', size: '5' });
      assert.equal(seen.headers['authorization'], 'bEaReR abc');
      assert.equal(Buffer.from(seen.body).toString(), '{"a":1}');

      assert.equal(await (await fetch(at('/echo'))).text(), '[]');

      const dropped = await fetch(at('/echo/7'), { method: 'DELETE' });
      assert.equal(dropped.status, 204);
      assert.equal(await dropped.text(), '');

      assert.equal((await fetch(at('/failing'))).status, 404);

      for (const [method, target] of [['GET', '/nowhere'], ['GET', '/echo/'], ['PUT', '/echo/7'], ['GET', '/echo/7/more'], ['PATCH', '/echo']]) {
        const unknown = await fetch(at(target ?? ''), { method: method ?? 'GET' });
        assert.deepEqual([unknown.status, await unknown.text()], [404, '{"error":"route not Found"}'], `${method} ${target}`);
      }

      const oversized = await fetch(at('/echo/1'), { method: 'POST', body: ' '.repeat(2 * MAX_BODY_BYTES) });
      assert.equal(oversized.status, 201);
      assert.equal(seen.body.length, MAX_BODY_BYTES + 1);
    } finally {
      await running.close();
    }
  });
}

test('aRoutePatternMatchesOnlyItsOwnShape', () => {
  assert.deepEqual(matchRoutePattern('/students/{id}', '/students/%31'), { id: '1' });
  for (const path of ['/students/', '/students', '/students/1/', '/student/1']) {
    assert.equal(matchRoutePattern('/students/{id}', path), undefined, path);
  }
});

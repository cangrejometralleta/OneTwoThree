import assert from 'node:assert/strict';
import { test } from 'node:test';

import { refuseUnprovenCaller, reportMissingRecord } from '../faults/faults.ts';
import { buildRequest } from '../transport/transport.ts';
import { answerWith, requireProvenCaller } from './answers.ts';
import type { CallerReader } from './ports.ts';

const fixedReader: CallerReader = {
  readCaller(token) {
    if (token !== 'good') {
      throw refuseUnprovenCaller('token is Invalid or Expired');
    }

    return 'student-registry';
  },
};

const bearing = (header: string) => buildRequest({ headers: { authorization: header } });

test('theCallerIsNamedBeforeTheStoryStartsWhateverTheCaseOfTheScheme', async () => {
  for (const header of ['Bearer good', 'bearer good', 'BEARER good']) {
    let seen = '';
    await requireProvenCaller(fixedReader, async (request) => {
      seen = request.caller;

      return 'told';
    })(bearing(header));
    assert.equal(seen, 'student-registry', header);
  }
});

test('theStoryNeverStartsForAStranger', async () => {
  let started = false;
  const guarded = requireProvenCaller(fixedReader, async () => {
    started = true;

    return undefined;
  });

  assert.equal((await answerWith(200, guarded)(bearing('Bearer forged'))).status, 401);
  assert.equal(started, false);
});

test('anAnswerCarriesTheDeclaredStatusOrTheFaultsOwn', async () => {
  assert.deepEqual(await answerWith(201, async () => 'made')(buildRequest()), { status: 201, body: 'made' });
  assert.deepEqual(await answerWith(200, async () => { throw reportMissingRecord('gone'); })(buildRequest()),
    { status: 404, body: { error: 'gone' } });
  assert.deepEqual(await answerWith(200, async () => { throw new Error('driver at /var/db broke'); })(buildRequest()),
    { status: 500, body: { error: 'internal error' } });
});

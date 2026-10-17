import assert from 'node:assert/strict';
import { test } from 'node:test';

import { readJsonText, writeJsonText } from './json.ts';

const bytesOf = (text: string): Uint8Array => new TextEncoder().encode(text);

test('aWholeNumberReadsAsABigintThatNeverRounds', () => {
  const read = readJsonText(bytesOf('{"a":9223372036854775807,"b":20.0,"c":-1,"d":"x\\u00e1","e":[true,null]}'));

  assert.ok(read instanceof Map);
  assert.equal(read.get('a'), 9223372036854775807n);
  assert.equal(read.get('b'), 20);
  assert.equal(read.get('c'), -1n);
  assert.equal(read.get('d'), 'xá');
  assert.deepEqual(read.get('e'), [true, null]);
});

test('anythingButOneWholeValueReadsAsNothing', () => {
  for (const text of ['', ' ', '{', '{"a":1} {}', '{"a":1,"a":2}', '{"a":01}', "{'a':1}", '{"a":1,}', '﻿{}', '[1 2]', 'nul']) {
    assert.equal(readJsonText(bytesOf(text)), undefined, JSON.stringify(text));
  }
  assert.equal(readJsonText(new Uint8Array([0x22, 0xe1, 0x22])), undefined, 'latin-1 is not UTF-8');
});

test('writingKeepsBigintsWholeAndDropsUndefinedFields', () => {
  assert.equal(writeJsonText({ id: 9223372036854775807n, name: 'Ada', gone: undefined, list: [1n] }),
    '{"id":9223372036854775807,"name":"Ada","list":[1]}');
});

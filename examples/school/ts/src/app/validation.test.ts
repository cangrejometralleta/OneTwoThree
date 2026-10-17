import assert from 'node:assert/strict';
import { test } from 'node:test';

import { DEFAULT_PAGE_SIZE, LARGEST_WHOLE, PAGE_IS_INVALID, WHOLE_SET } from '../school/school.ts';
import { buildRequest, MAX_BODY_BYTES } from '../transport/transport.ts';
import { BODY_IS_BROKEN, PATH_IS_BROKEN, readBearerToken, readJsonBody, readPageRequest, readPathNumber } from './validation.ts';

const STUDENT = { rut: 'text', name: 'text', age: 'whole', courseId: 'whole' } as const;
const WHOLE = '{"rut":"11111111-1","name":"Ana","age":20,"courseId":1}';
const bodyOf = (text: string | Uint8Array) => buildRequest({ body: typeof text === 'string' ? new TextEncoder().encode(text) : text });

test('readPathNumberWantsAPositiveWholeNumber', () => {
  assert.equal(readPathNumber(buildRequest({ path: { id: '7' } })), 7n);
  assert.equal(readPathNumber(buildRequest({ path: { id: '012' } })), 12n);

  for (const bad of ['', 'abc', '0', '000', '-3', '1.5', '+1', ' 1', '1e1', '0x1']) {
    assert.throws(() => readPathNumber(buildRequest({ path: { id: bad } })), (failure) => PATH_IS_BROKEN.sameAs(failure), bad);
  }
});

test('aPathNumberBeyondSixtyFourBitsNamesNobodyRatherThanBreaking', () => {
  assert.equal(readPathNumber(buildRequest({ path: { id: '99999999999999999999999' } })), LARGEST_WHOLE);
});

test('readJsonBodyRefusesAnythingButOneWholeObject', () => {
  assert.deepEqual(readJsonBody(bodyOf(WHOLE), STUDENT), { rut: '11111111-1', name: 'Ana', age: 20n, courseId: 1n });

  const broken: Record<string, string> = {
    'not json': '{"rut":',
    empty: '',
    'missing field': '{"rut":"11111111-1","name":"Ana","age":20}',
    'null field': '{"rut":"11111111-1","name":null,"age":20,"courseId":1}',
    'unknown field': '{"rut":"11111111-1","name":"Ana","age":20,"courseId":1,"id":9}',
    'wrong type': '{"rut":"11111111-1","name":"Ana","age":"20","courseId":1}',
    'text as number': '{"rut":11111111,"name":"Ana","age":20,"courseId":1}',
    fractional: '{"rut":"11111111-1","name":"Ana","age":20.5,"courseId":1}',
    'whole written as float': '{"rut":"11111111-1","name":"Ana","age":20.0,"courseId":1}',
    'beyond 64 bits': '{"rut":"11111111-1","name":"Ana","age":20,"courseId":9223372036854775808}',
    'second value': `${WHOLE} {}`,
    'an array': '[]',
    'null document': 'null',
    'repeated key': '{"rut":"11111111-1","rut":"x","name":"Ana","age":20,"courseId":1}',
  };

  for (const [name, body] of Object.entries(broken)) {
    assert.throws(() => readJsonBody(bodyOf(body), STUDENT), (failure) => BODY_IS_BROKEN.sameAs(failure), name);
  }
});

test('aNegativeCourseIsWellFormedAndLeftToTheCourseCheck', () => {
  assert.equal(readJsonBody(bodyOf('{"rut":"11111111-1","name":"Ana","age":20,"courseId":-1}'), STUDENT).courseId, -1n);
});

test('aBodyMustBeUtf8AndAtMostOneMebibyte', () => {
  const latin = Buffer.from(WHOLE.replace('Ana', 'Aná'), 'latin1');
  assert.throws(() => readJsonBody(bodyOf(latin), STUDENT), (failure) => BODY_IS_BROKEN.sameAs(failure));

  const oversized = bodyOf(WHOLE + ' '.repeat(MAX_BODY_BYTES));
  assert.throws(() => readJsonBody(oversized, STUDENT), (failure) => BODY_IS_BROKEN.sameAs(failure));
});

test('readPageRequestWindowsOnlyWhatWasAsked', () => {
  assert.deepEqual(readPageRequest(buildRequest()), WHOLE_SET);
  assert.deepEqual(readPageRequest(buildRequest({ query: { sort: 'name' } })), WHOLE_SET);
  assert.deepEqual(readPageRequest(buildRequest({ query: { page: '2', size: '5' } })), { number: 2n, size: 5n });
  assert.deepEqual(readPageRequest(buildRequest({ query: { page: '3' } })), { number: 3n, size: DEFAULT_PAGE_SIZE });
  assert.deepEqual(readPageRequest(buildRequest({ query: { size: '25' } })), { number: 0n, size: 25n });
  assert.deepEqual(readPageRequest(buildRequest({ query: { page: '99999999999999999999', size: '10' } })),
    { number: LARGEST_WHOLE, size: 10n });
});

test('readPageRequestRefusesWhatIsNotDigitsInRange', () => {
  for (const query of [{ page: 'abc' }, { page: '-1' }, { page: '' }, { page: '1.5' }, { page: '+1' }, { page: '1e1' },
    { page: ' 1' }, { page: '0x1' }, { size: 'many' }, { size: '0' }, { size: '-4' }, { size: '' }]) {
    assert.throws(() => readPageRequest(buildRequest({ query })), (failure) => PAGE_IS_INVALID.sameAs(failure), JSON.stringify(query));
  }
});

test('aBearerTokenIsReadOnlyAfterTheScheme', () => {
  assert.equal(readBearerToken(buildRequest({ headers: { authorization: 'bEaReR abc' } })), 'abc');
  assert.equal(readBearerToken(buildRequest({ headers: { authorization: 'Basic abc' } })), '');
  assert.equal(readBearerToken(buildRequest()), '');
});

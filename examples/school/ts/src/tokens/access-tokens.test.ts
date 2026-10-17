import assert from 'node:assert/strict';
import { test } from 'node:test';

import { AccessTokens, TOKEN_IS_INVALID } from './access-tokens.ts';

/** The Spec's two Tokens, Signed with spec-secret: one Lives until 2100, one Died at exp 1. */
const SPEC_LIVING = 'eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjo0MTAyNDQ0ODAwfQ.fgl9OtJeU8f90EIkPjELn34yWKjIbn3QcYn5UFoLAyY';
const SPEC_DEAD = 'eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjoxfQ.V8VQ8Ie_fo88z4EBJ0QbtIc5ONTOqQ_-l_WLuhwz5t4';
const ISSUED_AT = 1_000_000_000;

const clockedAt = (secret: string, now: number) => new AccessTokens(secret, 60n, () => now);
const encodeClaims = (json: string) => Buffer.from(json, 'utf8').toString('base64url');
const isInvalid = (failure: unknown) => TOKEN_IS_INVALID.sameAs(failure);

test('accessTokenDiesOnlyAfterItsDeadline', () => {
  const token = clockedAt('secret', ISSUED_AT).issueAccessToken('student-registry');

  for (const afterMilliseconds of [0, 59_000, 60_000, 60_999]) {
    assert.equal(clockedAt('secret', ISSUED_AT + afterMilliseconds).readCaller(token), 'student-registry', `${afterMilliseconds}`);
  }
  assert.throws(() => clockedAt('secret', ISSUED_AT + 61_000).readCaller(token), isInvalid);
});

test('aMintedTokenCarriesCompactClaimsInTheSpecsOrder', () => {
  const token = clockedAt('secret', ISSUED_AT).issueAccessToken('student-registry');

  assert.equal(token.split('.')[0], encodeClaims('{"sub":"student-registry","exp":1000060}'));
});

test('theSpecsTokensReadTheSameInEveryRuntime', () => {
  const spec = clockedAt('spec-secret', Date.parse('2026-10-17T00:00:00Z'));

  assert.equal(spec.readCaller(SPEC_LIVING), 'student-registry');
  assert.throws(() => spec.readCaller(SPEC_DEAD), isInvalid);
});

test('aTamperedClaimBreaksTheSignature', () => {
  const token = clockedAt('secret', ISSUED_AT).issueAccessToken('student-registry');
  const forged = `${encodeClaims('{"sub":"admin","exp":99999999999}')}.${token.split('.')[1]}`;

  assert.throws(() => clockedAt('secret', ISSUED_AT).readCaller(forged), isInvalid);
});

test('anotherSecretAndAnythingUnsignedAreRefused', () => {
  const token = clockedAt('secret', ISSUED_AT).issueAccessToken('student-registry');
  const other = clockedAt('another', ISSUED_AT);

  for (const candidate of [token, token.split('.')[0] ?? '', '', '.', '%%%.%%%', `${token}.x`]) {
    assert.throws(() => other.readCaller(candidate), isInvalid, candidate);
  }
});

test('aTokenWithNoSubjectIsRefused', () => {
  assert.throws(() => clockedAt('secret', ISSUED_AT).readCaller(`${encodeClaims('{"sub":"","exp":99999999999}')}.x`), isInvalid);
});

test('anEmptySecretSignsNothingAndTrustsNothing', () => {
  assert.throws(() => clockedAt('', ISSUED_AT).issueAccessToken('student-registry'), isInvalid);
  assert.throws(() => clockedAt('secret', ISSUED_AT).issueAccessToken(''), isInvalid);
  assert.throws(() => clockedAt('', ISSUED_AT).readCaller(SPEC_LIVING), isInvalid);
});

test('everyTokenFailureAnswersAsAnUnprovenCaller', () => {
  assert.equal(TOKEN_IS_INVALID.kind, 'unproven-caller');
});

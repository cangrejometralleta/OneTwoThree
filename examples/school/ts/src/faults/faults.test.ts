import assert from 'node:assert/strict';
import { test } from 'node:test';

import {
  readFaultKind, readFaultReason, refuseInvalidInput, refuseUnprocessableValue, refuseUnprovenCaller, reportMissingRecord,
  reportTakenValue, reportUnavailableProvider,
} from './faults.ts';

test('eachFaultAnswersWithTheKindItDeclared', () => {
  const cases = [
    [refuseInvalidInput('bad'), 'invalid-input'],
    [refuseUnprovenCaller('who'), 'unproven-caller'],
    [reportMissingRecord('gone'), 'missing-record'],
    [reportTakenValue('taken'), 'taken-value'],
    [refuseUnprocessableValue('no'), 'unprocessable-value'],
    [reportUnavailableProvider('down'), 'unavailable-provider'],
  ] as const;

  for (const [fault, kind] of cases) {
    assert.equal(readFaultKind(fault), kind, fault.reason);
  }
});

test('anUncontrolledFailureIsOursAndItsMessageNeverLeaves', () => {
  const driver = new Error('the driver Broke at /var/db');

  assert.equal(readFaultKind(driver), 'unexpected');
  assert.equal(readFaultReason(driver), 'internal error');
  assert.equal(readFaultReason('a thrown string'), 'internal error');
});

test('aWrappedFaultStillCarriesItsKindAndReason', () => {
  const wrapped = new Error('while Reading', { cause: reportMissingRecord('gone') });

  assert.equal(readFaultKind(wrapped), 'missing-record');
  assert.equal(readFaultReason(wrapped), 'gone');
});

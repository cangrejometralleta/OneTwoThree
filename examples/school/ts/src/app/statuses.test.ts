import assert from 'node:assert/strict';
import { test } from 'node:test';

import {
  refuseInvalidInput, refuseUnprocessableValue, refuseUnprovenCaller, reportMissingRecord, reportTakenValue,
  reportUnavailableProvider,
} from '../faults/faults.ts';
import { readFaultStatus } from './statuses.ts';

test('readFaultStatusAnswersForEachKind', () => {
  const cases = [
    [refuseInvalidInput('bad'), 400], [refuseUnprovenCaller('who'), 401], [reportMissingRecord('gone'), 404],
    [reportTakenValue('taken'), 409], [refuseUnprocessableValue('no'), 422], [reportUnavailableProvider('down'), 503],
  ] as const;

  for (const [fault, status] of cases) {
    assert.equal(readFaultStatus(fault), status, fault.reason);
  }
});

test('readFaultStatusRefusesToGuess', () => {
  assert.equal(readFaultStatus(new Error('the driver Broke')), 500);
});

test('readFaultStatusSurvivesWrapping', () => {
  assert.equal(readFaultStatus(new Error('while Reading', { cause: reportMissingRecord('gone') })), 404);
});

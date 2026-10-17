import assert from 'node:assert/strict';
import { test } from 'node:test';

import {
  ageIsTooLow, buildPage, checkCourseRecord, checkStudentRecord, CODE_IS_EMPTY, LARGEST_WHOLE, NAME_IS_EMPTY,
  PAGE_IS_INVALID, pageOffset, readRut, RUT_IS_INVALID, rutLooksValid, UNNUMBERED_COURSE, UNNUMBERED_STUDENT,
  type CourseId, type Student,
} from './school.ts';

const rules = { minimumAgeYears: 18n };
const studentOf = (rut: string, name: string, age: bigint): Student =>
  ({ id: UNNUMBERED_STUDENT, rut: readRut(rut), name, age, course: 1n as CourseId });

test('rutLooksValidAcceptsRealNumbers', () => {
  for (const rut of ['12345678-5', '11111111-1', '22222222-2', '10000013-k', '10000013-K', '10000004-0', '123456789012345678901234567890-1']) {
    assert.ok(rutLooksValid(rut), rut);
  }
});

test('rutLooksValidRefusesBadDigits', () => {
  for (const rut of ['12345678-9', '12345678', 'abc-1', '-5', '12345678-55', '', ' 12345678-5']) {
    assert.ok(!rutLooksValid(rut), rut);
  }
});

test('aRutEntersInLowerCaseAndKeepsItsLeadingZeros', () => {
  assert.equal(readRut('10000013-K'), '10000013-k');
  assert.notEqual(readRut('012345678-5'), readRut('12345678-5'));
});

test('checkStudentRecordGuardsEachRule', () => {
  checkStudentRecord(studentOf('12345678-5', 'Ada', 20n), rules);
  checkStudentRecord(studentOf('12345678-5', 'Ada', 18n), rules);
  assert.throws(() => checkStudentRecord(studentOf('bad', '', 1n), rules), (failure) => NAME_IS_EMPTY.sameAs(failure));
  assert.throws(() => checkStudentRecord(studentOf('12345678-5', ' \t', 20n), rules), (failure) => NAME_IS_EMPTY.sameAs(failure));
  assert.throws(() => checkStudentRecord(studentOf('12345678-9', 'Ada', 1n), rules), (failure) => RUT_IS_INVALID.sameAs(failure));
  assert.throws(() => checkStudentRecord(studentOf('12345678-5', 'Ada', 17n), rules), (failure) => ageIsTooLow(rules).sameAs(failure));
});

test('theAgeReasonSpellsTheFloorFromTheConstant', () => {
  assert.equal(ageIsTooLow(rules).reason, 'age Must be 18 or more');
  assert.equal(ageIsTooLow({ minimumAgeYears: 21n }).reason, 'age Must be 21 or more');
});

test('checkCourseRecordGuardsEachRule', () => {
  checkCourseRecord({ id: UNNUMBERED_COURSE, code: 'MAT-101', name: 'Algebra' });
  assert.throws(() => checkCourseRecord({ id: UNNUMBERED_COURSE, code: '', name: 'Algebra' }), (failure) => CODE_IS_EMPTY.sameAs(failure));
  assert.throws(() => checkCourseRecord({ id: UNNUMBERED_COURSE, code: 'MAT-101', name: ' ' }), (failure) => NAME_IS_EMPTY.sameAs(failure));
  assert.throws(() => checkCourseRecord({ id: UNNUMBERED_COURSE, code: '', name: '' }), (failure) => CODE_IS_EMPTY.sameAs(failure));
});

test('aPageRefusesWhatTheStoreCannotServe', () => {
  buildPage(0n, 1n);
  assert.throws(() => buildPage(-1n, 10n), (failure) => PAGE_IS_INVALID.sameAs(failure));
  assert.throws(() => buildPage(0n, 0n), (failure) => PAGE_IS_INVALID.sameAs(failure));
});

test('aPageBeyondSixtyFourBitsWindowsAsFarAsItCan', () => {
  assert.equal(pageOffset({ number: 2n, size: 10n }), 20n);
  assert.equal(pageOffset({ number: LARGEST_WHOLE, size: 10n }), LARGEST_WHOLE);
});

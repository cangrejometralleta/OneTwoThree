import assert from 'node:assert/strict';
import { mkdtempSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { test } from 'node:test';

import {
  LARGEST_WHOLE, readRut, RUT_TAKEN, UNNUMBERED_COURSE, UNNUMBERED_STUDENT, WHOLE_SET,
  type CourseId, type Student, type StudentId,
} from '../school/school.ts';
import { SchoolStore } from './school-store.ts';

const openTestStore = () => SchoolStore.openSchoolStore(join(mkdtempSync(join(tmpdir(), 'store-')), 'school.db'));
const studentOf = (rut: string, name: string): Student =>
  ({ id: UNNUMBERED_STUDENT, rut: readRut(rut), name, age: 20n, course: 1n as CourseId });

test('insertStudentRowRefusesADuplicateRut', async () => {
  const books = openTestStore();
  await books.insertStudentRow(studentOf('12345678-5', 'Ada'));

  await assert.rejects(books.insertStudentRow(studentOf('12345678-5', 'Grace')), (failure) => RUT_TAKEN.sameAs(failure));
});

test('aStudentKeepsWhatItWasGivenUntrimmedAndWhole', async () => {
  const books = openTestStore();
  const kept = await books.insertStudentRow({ ...studentOf('12345678-5', ' Ada '), age: LARGEST_WHOLE });

  assert.ok(kept.id > 0n);
  assert.deepEqual(await books.selectStudentRow(kept.id), kept);

  const renamed = { ...kept, name: 'Ada Lovelace', age: 21n };
  assert.equal(await books.updateStudentRow(renamed), true);
  assert.deepEqual(await books.selectStudentRow(kept.id), renamed);
});

test('theStoreSaysNothingOfWhatItNeverKept', async () => {
  const books = openTestStore();

  assert.equal(await books.selectStudentRow(9n as StudentId), undefined);
  assert.equal(await books.selectStudentRow(LARGEST_WHOLE as StudentId), undefined);
  assert.equal(await books.updateStudentRow({ ...studentOf('12345678-5', 'Ada'), id: 9n as StudentId }), false);
  assert.equal(await books.deleteStudentRow(9n as StudentId), false);
  assert.equal(await books.selectCourseRow(9n as CourseId), undefined);
});

test('studentsListByNameInBinaryOrderTiesByIdAndWholeWhenNoWindowIsAsked', async () => {
  const books = openTestStore();
  for (const [rut, name] of [['11111111-1', 'carla'], ['22222222-2', 'Berta'], ['12345678-5', 'Ana'], ['10000013-k', 'Ana']]) {
    await books.insertStudentRow(studentOf(rut ?? '', name ?? ''));
  }

  assert.deepEqual((await books.selectStudentPage(WHOLE_SET)).map((s) => s.rut), ['12345678-5', '10000013-k', '22222222-2', '11111111-1']);
  assert.deepEqual((await books.selectStudentPage({ number: 1n, size: 3n })).map((s) => s.name), ['carla']);
});

test('aWindowBeyondTheRollIsEmpty', async () => {
  const books = openTestStore();
  await books.insertStudentRow(studentOf('11111111-1', 'Ana'));

  assert.deepEqual(await books.selectStudentPage({ number: LARGEST_WHOLE, size: 10n }), []);
  assert.deepEqual(await books.selectStudentPage({ number: 5n, size: 10n }), []);
});

test('coursesListByCodeAndCodesAreNotUnique', async () => {
  const books = openTestStore();
  for (const code of ['MAT-101', 'FIS-201', 'ART-100', 'ART-100']) {
    await books.insertCourseRow({ id: UNNUMBERED_COURSE, code, name: 'Name' });
  }

  assert.deepEqual((await books.selectCoursePage(WHOLE_SET)).map((c) => c.code), ['ART-100', 'ART-100', 'FIS-201', 'MAT-101']);
});

test('aDeletedStudentIsGoneOnce', async () => {
  const books = openTestStore();
  const kept = await books.insertStudentRow(studentOf('12345678-5', 'Ada'));

  assert.equal(await books.deleteStudentRow(kept.id), true);
  assert.equal(await books.deleteStudentRow(kept.id), false);
});

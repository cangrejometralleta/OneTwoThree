import assert from 'node:assert/strict';
import { beforeEach, test } from 'node:test';

import {
  CODE_IS_EMPTY, COURSE_UNKNOWN, NAME_IS_EMPTY, readRut, REGISTRY_UNAVAILABLE, RUT_UNREGISTERED, STUDENT_UNKNOWN,
  UNNUMBERED_COURSE, UNNUMBERED_STUDENT, type CourseId, type Student, type StudentId,
} from '../school/school.ts';
import { FrontDesk, Shelf } from './fakes.ts';
import { SchoolService } from './service.ts';

let books: Shelf;
let front: FrontDesk;
let service: SchoolService;

const buildAdmission = (desk: FrontDesk | undefined) =>
  new SchoolService({ rules: { minimumAgeYears: 18n }, students: books, courses: books, registry: desk, notices: desk, log: () => {} });

const readyStudent = (course = 1n): Student =>
  ({ id: UNNUMBERED_STUDENT, rut: readRut('11111111-1'), name: 'Ana', age: 20n, course: course as CourseId });

const isFault = (fault: { sameAs(other: unknown): boolean }) => (failure: unknown) => fault.sameAs(failure);

beforeEach(() => {
  books = new Shelf();
  front = new FrontDesk();
  service = buildAdmission(front);
});

test('enrollStudentAnnouncesOnceTheStoreKeeps', async () => {
  const kept = await service.enrollStudent(readyStudent());

  assert.equal(books.students.size, 1);
  assert.equal(front.announced.length, 1);
  assert.equal(front.announced[0]?.id, kept.id);
});

test('enrollStudentRefusesWhatTheRegistryDenies', async () => {
  front.registered = false;

  await assert.rejects(service.enrollStudent(readyStudent()), isFault(RUT_UNREGISTERED));
  assert.equal(books.students.size, 0);
  assert.equal(front.announced.length, 0);
});

test('enrollStudentNeverTreatsAnUnknownAsAYes', async () => {
  front.down = new Error('socket closed');
  await assert.rejects(service.enrollStudent(readyStudent()), isFault(REGISTRY_UNAVAILABLE));
  assert.equal(books.students.size, 0);

  await assert.rejects(buildAdmission(undefined).enrollStudent(readyStudent()), isFault(REGISTRY_UNAVAILABLE));
});

test('enrollStudentKeepsAFailedNoticeInsideTheAnswer', async () => {
  front.noticeFailure = new Error('mail is down');

  await service.enrollStudent(readyStudent());

  assert.equal(books.students.size, 1);
});

test('theRegistryIsAskedAfterTheFormAndTheCourse', async () => {
  await assert.rejects(service.enrollStudent({ ...readyStudent(), name: '' }), isFault(NAME_IS_EMPTY));
  await assert.rejects(service.enrollStudent(readyStudent(9n)), isFault(COURSE_UNKNOWN));
  await assert.rejects(service.enrollStudent(readyStudent(0n)), isFault(COURSE_UNKNOWN));
  assert.equal(front.asked.length, 0);
});

test('theStoreIsWrittenOnlyAfterTheRegistryAnswers', async () => {
  front.registered = false;
  await assert.rejects(service.enrollStudent(readyStudent()));

  assert.equal(front.asked.length, 1);
  assert.equal(books.students.size, 0);
});

test('saveStudentAsksTheRegistryAndAnnouncesNothing', async () => {
  const kept = await service.enrollStudent(readyStudent());
  const announcedByEnrolment = front.announced.length;

  front.registered = false;
  await assert.rejects(service.saveStudent(kept), isFault(RUT_UNREGISTERED));

  front.registered = true;
  await service.saveStudent(kept);
  assert.equal(front.announced.length, announcedByEnrolment);
});

test('saveStudentRefusesAStudentNobodyHoldsOnlyAtTheStore', async () => {
  const ghost = { ...readyStudent(), id: 42n as StudentId };

  await assert.rejects(service.saveStudent(ghost), isFault(STUDENT_UNKNOWN));
  front.registered = false;
  await assert.rejects(service.saveStudent(ghost), isFault(RUT_UNREGISTERED));
});

test('readStudentFindsWhatEnrolmentKept', async () => {
  const kept = await service.enrollStudent(readyStudent());

  assert.deepEqual(await service.readStudent(kept.id), kept);
});

test('readStudentRefusesWhoNeverEnrolled', async () => {
  await assert.rejects(service.readStudent(42n as StudentId), isFault(STUDENT_UNKNOWN));
});

test('listStudentsReadsTheRequestedPageOfTheRoll', async () => {
  await service.enrollStudent(readyStudent());

  assert.equal((await service.listStudents({ number: 2n, size: 5n })).length, 1);
  assert.deepEqual(books.page, { number: 2n, size: 5n });
});

test('deleteStudentEndsAnEnrolmentOnce', async () => {
  const kept = await service.enrollStudent(readyStudent());

  await service.deleteStudent(kept.id);

  await assert.rejects(service.readStudent(kept.id), isFault(STUDENT_UNKNOWN));
  await assert.rejects(service.deleteStudent(kept.id), isFault(STUDENT_UNKNOWN));
});

test('createCourseOpensACourseTheCatalogueLists', async () => {
  const opened = await service.createCourse({ id: UNNUMBERED_COURSE, code: 'FIS-201', name: 'Physics' });

  assert.deepEqual(await service.readCourse(opened.id), opened);
  assert.equal((await service.listCourses({ number: 0n, size: 10n })).length, 2);
  assert.equal(books.page?.size, 10n);
});

test('createCourseRefusesACourseWithoutCodeOrName', async () => {
  await assert.rejects(service.createCourse({ id: UNNUMBERED_COURSE, code: '', name: 'Physics' }), isFault(CODE_IS_EMPTY));
  await assert.rejects(service.createCourse({ id: UNNUMBERED_COURSE, code: 'FIS-201', name: '' }), isFault(NAME_IS_EMPTY));
  assert.equal(books.courses.size, 1);
});

test('readCourseRefusesACourseThatNeverOpened', async () => {
  await assert.rejects(service.readCourse(99n as CourseId), isFault(COURSE_UNKNOWN));
});

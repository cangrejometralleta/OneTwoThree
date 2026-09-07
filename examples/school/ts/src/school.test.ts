import assert from "node:assert/strict";
import test from "node:test";

import { translateRoutePattern } from "./adapters.js";
import {
  ErrCourseUnknown, ErrRutTaken, ErrStudentUnknown, checkStudentRecord, rutLooksValid,
  type Course, type CourseID, type Student, type StudentID,
} from "./domain.js";
import { SchoolAPI } from "./handlers.js";
import type { CourseStore, Page, StudentStore } from "./providers.js";
import { SqliteSchool } from "./store-sqlite.js";
import { HmacTokens } from "./token-hmac.js";
import type { Request } from "./transport.js";

// FakeSchool Stands in for SQLite.
// No Database, no Framework, no Port: the Providers Allow it.
class FakeSchool implements StudentStore, CourseStore {
  private students = new Map<number, Student>();
  private courses = new Map<number, Course>([
    [1, { id: 1 as CourseID, code: "MAT-101" as never, name: "Algebra" as never }],
  ]);

  insertStudentRow(student: Student): Student {
    const stored = { ...student, id: (this.students.size + 1) as StudentID };
    this.students.set(stored.id, stored);
    return stored;
  }

  selectStudentRow(id: StudentID): Student {
    const student = this.students.get(id);
    if (!student) throw ErrStudentUnknown;
    return student;
  }

  updateStudentRow(student: Student): Student {
    if (!this.students.has(student.id)) throw ErrStudentUnknown;
    this.students.set(student.id, student);
    return student;
  }

  deleteStudentRow(id: StudentID): void {
    if (!this.students.delete(id)) throw ErrStudentUnknown;
  }

  selectStudentPage(_: Page): Student[] {
    return [...this.students.values()];
  }

  insertCourseRow(course: Course): Course {
    const stored = { ...course, id: (this.courses.size + 1) as CourseID };
    this.courses.set(stored.id, stored);
    return stored;
  }

  selectCourseRow(id: CourseID): Course {
    const course = this.courses.get(id);
    if (!course) throw ErrCourseUnknown;
    return course;
  }

  selectCoursePage(_: Page): Course[] {
    return [...this.courses.values()];
  }
}

// buildTestingSchool Hands the API three Fakes and nothing else.
function buildTestingSchool(): SchoolAPI {
  const fake = new FakeSchool();
  return new SchoolAPI(fake, fake, new HmacTokens("test", 60_000));
}

function askSchool(body: unknown, path: Record<string, string> = {}): Request {
  return { path, query: {}, token: "", body: JSON.stringify(body) };
}

function askSchoolPage(query: Record<string, string>): Request {
  return { path: {}, query, token: "", body: "" };
}

// Number("abc") is NaN, and NaN Passes every Comparison it is Given.
test("someone Asks for a Page Numbered with a Word", () => {
  const reply = buildTestingSchool().listStudentRecords(askSchoolPage({ page: "abc" }));

  assert.equal(reply.status, 400);
});

test("someone Asks for a Size Measured in Words", () => {
  const reply = buildTestingSchool().listCourseRecords(askSchoolPage({ size: "many" }));

  assert.equal(reply.status, 400);
});

test("someone Omits Pagination and Gets the whole Set", () => {
  const reply = buildTestingSchool().listStudentRecords(askSchoolPage({}));

  assert.equal(reply.status, 200);
});

test("rutLooksValid accepts real Numbers", () => {
  for (const rut of ["12345678-5", "11111111-1", "8459162-9", "6-k", "1-9"]) {
    assert.equal(rutLooksValid(rut), true, `${rut} Should be valid`);
  }
});

test("rutLooksValid refuses bad Digits", () => {
  for (const rut of ["12345678-9", "11111111-2", "12345678", "abc-1", ""]) {
    assert.equal(rutLooksValid(rut), false, `${rut} Should be Refused`);
  }
});

test("checkStudentRecord guards each Rule", () => {
  const good = { rut: "12345678-5", name: "Ada", age: 20, course: 1 } as unknown as Student;
  assert.equal(checkStudentRecord(good), null);

  assert.equal(checkStudentRecord({ ...good, age: 18 }), null);

  const young = { ...good, age: 17 };
  assert.equal(checkStudentRecord(young)?.message, "age Must be 18 or more");
});

test("addStudentRecord enrols someone", () => {
  const reply = buildTestingSchool().addStudentRecord(
    askSchool({ rut: "12345678-5", name: "Ada", age: 20, courseId: 1 }),
  );
  assert.equal(reply.status, 201);
});

test("addStudentRecord refuses a bad RUT", () => {
  const reply = buildTestingSchool().addStudentRecord(
    askSchool({ rut: "12345678-9", name: "Ada", age: 20, courseId: 1 }),
  );
  assert.equal(reply.status, 400);
});

test("addStudentRecord refuses an unknown Course", () => {
  const reply = buildTestingSchool().addStudentRecord(
    askSchool({ rut: "12345678-5", name: "Ada", age: 20, courseId: 99 }),
  );
  assert.equal(reply.status, 404);
});

test("showStudentRecord reports an Absence", () => {
  const reply = buildTestingSchool().showStudentRecord(askSchool(null, { id: "42" }));
  assert.equal(reply.status, 404);
});

// The Comparison that Broke in the Spring Version, Pinned here.
test("an Access Token dies only after its Deadline", () => {
  let clock = Date.now();
  const tokens = new HmacTokens("s", 60_000, () => clock);

  const token = tokens.issueAccessToken("registry");
  assert.equal(tokens.readAccessToken(token), "registry");

  clock += 120_000;
  assert.throws(() => tokens.readAccessToken(token));
});

test("an Access Token refuses a tampered Claim", () => {
  const tokens = new HmacTokens("s", 60_000);
  const token = tokens.issueAccessToken("registry");

  assert.throws(() => tokens.readAccessToken("admin" + token.slice(8)));
});

test("translateRoutePattern feeds express and fastify", () => {
  assert.equal(translateRoutePattern("/students/{id}"), "/students/:id");
});

// The Comparison that a Fake cannot Pin: a real unique Index.
test("insertStudentRow refuses a duplicate RUT", () => {
  const school = SqliteSchool.shapeSchoolTables(":memory:");
  const course = school.insertCourseRow({ code: "MAT-101", name: "Algebra" } as unknown as Course);
  const student = { rut: "12345678-5", name: "Ada", age: 20, course: course.id } as unknown as Student;

  school.insertStudentRow(student);

  assert.throws(() => school.insertStudentRow(student), (err) => err === ErrRutTaken);
});

import Database from 'better-sqlite3';
import { asc, eq, sql } from 'drizzle-orm';
import { drizzle, type BetterSQLite3Database } from 'drizzle-orm/better-sqlite3';
import { customType, sqliteTable, text } from 'drizzle-orm/sqlite-core';

import type { CourseStore, StudentStore } from '../app/ports.ts';
import { pageOffset, RUT_TAKEN, type Course, type CourseId, type Page, type Rut, type Student, type StudentId } from '../school/school.ts';

/** wholeColumn is an SQLite INTEGER Read as a 64 bit bigint, so nothing Rounds. */
const wholeColumn = customType<{ data: bigint; driverData: bigint }>({ dataType: () => 'integer' });

/** numberedByStore is an Identity SQLite Gives: an INTEGER PRIMARY KEY Sent NULL Numbers itself. */
const numberedByStore = (name: string) => wholeColumn(name).primaryKey().$defaultFn(() => sql`NULL`);

/** studentRows is the Storage Shape of a Student: it Knows its Table and no one else does. */
const studentRows = sqliteTable('student_rows', {
  id: numberedByStore('id'),
  rut: text('rut').notNull().unique(),
  name: text('name').notNull(),
  age: wholeColumn('age').notNull(),
  courseId: wholeColumn('course_id').notNull(),
});

/** courseRows is the Storage Shape of a Course. */
const courseRows = sqliteTable('course_rows', {
  id: numberedByStore('id'),
  code: text('code').notNull(),
  name: text('name').notNull(),
});

type StudentRow = typeof studentRows.$inferSelect;
type CourseRow = typeof courseRows.$inferSelect;

/** UNIQUE_RUT_VIOLATION is how SQLite Says the unique Index on the Rut Refused a Row. */
const UNIQUE_RUT_VIOLATION = 'UNIQUE constraint failed: student_rows.rut';

/**
 * SchoolStore Keeps Students and Courses in one SQLite File through Drizzle.
 * This Module is the only one that Imports the ORM or the Driver.
 * Reference: https://orm.drizzle.team/docs/get-started-sqlite
 */
export class SchoolStore implements StudentStore, CourseStore {
  readonly #db: BetterSQLite3Database;

  private constructor(db: BetterSQLite3Database) {
    this.#db = db;
  }

  /** openSchoolStore Opens the Database and Brings its Tables up to Date. */
  static openSchoolStore(path: string): SchoolStore {
    const connection = new Database(path);
    connection.defaultSafeIntegers(true);
    const db = drizzle(connection);
    db.run(sql`CREATE TABLE IF NOT EXISTS student_rows (
      id INTEGER PRIMARY KEY AUTOINCREMENT, rut TEXT NOT NULL UNIQUE, name TEXT NOT NULL,
      age INTEGER NOT NULL, course_id INTEGER NOT NULL)`);
    db.run(sql`CREATE TABLE IF NOT EXISTS course_rows (
      id INTEGER PRIMARY KEY AUTOINCREMENT, code TEXT NOT NULL, name TEXT NOT NULL)`);

    return new SchoolStore(db);
  }

  /** insertStudentRow Writes a new Student and Returns it Numbered; the unique Index Proves the Rut. */
  async insertStudentRow(student: Student): Promise<Student> {
    const { id: _unnumbered, ...row } = encodeStudentRow(student);
    try {
      const [kept] = this.#db.insert(studentRows).values(row).returning().all();

      return decodeStudentRow(kept as StudentRow);
    } catch (failure) {
      throw translateStoreError(failure);
    }
  }

  /** selectStudentRow Finds one Student, or Nothing. */
  async selectStudentRow(id: StudentId): Promise<Student | undefined> {
    const row = this.#db.select().from(studentRows).where(eq(studentRows.id, id)).get();

    return row === undefined ? undefined : decodeStudentRow(row);
  }

  /** updateStudentRow Rewrites the Student its Id Names; false when nobody Holds it. */
  async updateStudentRow(student: Student): Promise<boolean> {
    const { id, ...row } = encodeStudentRow(student);
    try {
      return this.#db.update(studentRows).set(row).where(eq(studentRows.id, id)).run().changes > 0;
    } catch (failure) {
      throw translateStoreError(failure);
    }
  }

  /** deleteStudentRow Removes one Student; false when nobody Holds it. */
  async deleteStudentRow(id: StudentId): Promise<boolean> {
    return this.#db.delete(studentRows).where(eq(studentRows.id, id)).run().changes > 0;
  }

  /** selectStudentPage Reads one Page, Ordered by Name in binary Order, Ties by Id. */
  async selectStudentPage(page: Page): Promise<Student[]> {
    const rows = this.#db.select().from(studentRows).orderBy(asc(studentRows.name), asc(studentRows.id))
      .limit(sql`${page.size}` as unknown as number).offset(sql`${pageOffset(page)}` as unknown as number).all();

    return rows.map(decodeStudentRow);
  }

  /** insertCourseRow Writes a new Course and Returns it Numbered. */
  async insertCourseRow(course: Course): Promise<Course> {
    const [kept] = this.#db.insert(courseRows).values({ code: course.code, name: course.name }).returning().all();

    return decodeCourseRow(kept as CourseRow);
  }

  /** selectCourseRow Finds one Course, or Nothing. */
  async selectCourseRow(id: CourseId): Promise<Course | undefined> {
    const row = this.#db.select().from(courseRows).where(eq(courseRows.id, id)).get();

    return row === undefined ? undefined : decodeCourseRow(row);
  }

  /** selectCoursePage Reads one Page of Courses, Ordered by Code in binary Order, Ties by Id. */
  async selectCoursePage(page: Page): Promise<Course[]> {
    const rows = this.#db.select().from(courseRows).orderBy(asc(courseRows.code), asc(courseRows.id))
      .limit(sql`${page.size}` as unknown as number).offset(sql`${pageOffset(page)}` as unknown as number).all();

    return rows.map(decodeCourseRow);
  }
}

/** translateStoreError Turns a driver Error into a Business Fault when it Means one; anything else Stays Ours. */
function translateStoreError(failure: unknown): unknown {
  const message = failure instanceof Error ? `${failure.message} ${String((failure as { cause?: unknown }).cause ?? '')}` : '';

  return message.includes(UNIQUE_RUT_VIOLATION) ? RUT_TAKEN : failure;
}

function encodeStudentRow(student: Student): StudentRow {
  return { id: student.id, rut: student.rut, name: student.name, age: student.age, courseId: student.course };
}

function decodeStudentRow(row: StudentRow): Student {
  return { id: row.id as StudentId, rut: row.rut as Rut, name: row.name, age: row.age, course: row.courseId as CourseId };
}

function decodeCourseRow(row: CourseRow): Course {
  return { id: row.id as CourseId, code: row.code, name: row.name };
}

import { DatabaseSync } from "node:sqlite";

import {
  ErrCourseUnknown, ErrStudentUnknown,
  type Course, type CourseCode, type CourseID,
  type FullName, type RUT, type Student, type StudentID,
} from "./domain.js";
import type { CourseStore, Page, StudentStore } from "./providers.js";

// StudentRow is the Storage Shape of a Student.
type StudentRow = {
  id: number; rut: string; name: string; age: number; course_id: number;
};

type CourseRow = { id: number; code: string; name: string };

// SqliteSchool Fulfils both Stores with one Connection.
// It is the only Class in this Program that Knows SQLite Exists.
// Reference: https://nodejs.org/api/sqlite.html
export class SqliteSchool implements StudentStore, CourseStore {
  constructor(private readonly db: DatabaseSync) {}

  // shapeSchoolTables Creates both Tables, Documenting the Schema in place.
  static shapeSchoolTables(path: string): SqliteSchool {
    const db = new DatabaseSync(path);

    db.exec(`
      CREATE TABLE IF NOT EXISTS course (
        id   INTEGER PRIMARY KEY AUTOINCREMENT,
        code TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL
      );
      CREATE TABLE IF NOT EXISTS student (
        id        INTEGER PRIMARY KEY AUTOINCREMENT,
        rut       TEXT NOT NULL UNIQUE,
        name      TEXT NOT NULL,
        age       INTEGER NOT NULL,
        course_id INTEGER NOT NULL REFERENCES course(id)
      );`);

    return new SqliteSchool(db);
  }

  insertStudentRow(student: Student): Student {
    const insert = this.db.prepare(
      "INSERT INTO student (rut, name, age, course_id) VALUES (?, ?, ?, ?)",
    );

    const result = insert.run(student.rut, student.name, student.age, student.course);

    return { ...student, id: Number(result.lastInsertRowid) as StudentID };
  }

  selectStudentRow(id: StudentID): Student {
    const row = this.db.prepare("SELECT * FROM student WHERE id = ?").get(id) as StudentRow | undefined;

    if (!row) throw ErrStudentUnknown;

    return decodeStudentRow(row);
  }

  updateStudentRow(student: Student): Student {
    const update = this.db.prepare(
      "UPDATE student SET rut = ?, name = ?, age = ?, course_id = ? WHERE id = ?",
    );

    const result = update.run(student.rut, student.name, student.age, student.course, student.id);

    if (result.changes === 0) throw ErrStudentUnknown;
    return student;
  }

  deleteStudentRow(id: StudentID): void {
    const result = this.db.prepare("DELETE FROM student WHERE id = ?").run(id);

    if (result.changes === 0) throw ErrStudentUnknown;
  }

  selectStudentPage(page: Page): Student[] {
    const rows = this.db
      .prepare("SELECT * FROM student ORDER BY name LIMIT ? OFFSET ?")
      .all(pageLimit(page), page.number * page.size) as StudentRow[];

    return rows.map(decodeStudentRow);
  }

  insertCourseRow(course: Course): Course {
    const insert = this.db.prepare("INSERT INTO course (code, name) VALUES (?, ?)");

    const result = insert.run(course.code, course.name);

    return { ...course, id: Number(result.lastInsertRowid) as CourseID };
  }

  selectCourseRow(id: CourseID): Course {
    const row = this.db.prepare("SELECT * FROM course WHERE id = ?").get(id) as CourseRow | undefined;

    if (!row) throw ErrCourseUnknown;

    return decodeCourseRow(row);
  }

  selectCoursePage(page: Page): Course[] {
    const rows = this.db
      .prepare("SELECT * FROM course ORDER BY code LIMIT ? OFFSET ?")
      .all(pageLimit(page), page.number * page.size) as CourseRow[];

    return rows.map(decodeCourseRow);
  }
}

// pageLimit Turns an absent Size into the whole Set.
function pageLimit(page: Page): number {
  return page.size > 0 ? page.size : -1;
}

// decodeStudentRow Turns Storage back into Business.
function decodeStudentRow(row: StudentRow): Student {
  return {
    id: row.id as StudentID,
    rut: row.rut as RUT,
    name: row.name as FullName,
    age: row.age,
    course: row.course_id as CourseID,
  };
}

function decodeCourseRow(row: CourseRow): Course {
  return { id: row.id as CourseID, code: row.code as CourseCode, name: row.name as FullName };
}

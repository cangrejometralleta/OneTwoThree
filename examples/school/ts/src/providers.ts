import type { Course, CourseID, Student, StudentID } from "./domain.js";

// A Provider is an Interface the Core Declares
// and something outside Fulfils.
// The Core Depends on the Shape, never on the Library behind it.
//
// The Word Collides, and the Collision is worth Knowing.
// Angular and NestJS Call a registered Dependency a Provider.
// JSR-330 Calls a Factory a Provider. Terraform Calls a Plugin one.
// Here it Means the Door an outside System Enters through,
// which the Literature Calls a Port.
// Reference: https://alistair.cockburn.us/hexagonal-architecture/

// Page Carries Pagination without Naming a Database.
export type Page = { number: number; size: number };

// StudentStore Keeps Students wherever Students Live.
// Fulfilled by SqliteSchool.
// Reference: https://nodejs.org/api/sqlite.html
export interface StudentStore {
  insertStudentRow(student: Student): Student;
  selectStudentRow(id: StudentID): Student;
  updateStudentRow(student: Student): Student;
  deleteStudentRow(id: StudentID): void;
  selectStudentPage(page: Page): Student[];
}

// CourseStore Keeps Courses, and Answers whether one Exists.
// Fulfilled by SqliteSchool.
// Reference: https://nodejs.org/api/sqlite.html
export interface CourseStore {
  insertCourseRow(course: Course): Course;
  selectCourseRow(id: CourseID): Course;
  selectCoursePage(page: Page): Course[];
}

// TokenIssuer Mints and Reads a Bearer Token.
// Fulfilled by HmacTokens, so no third Party Enters for this.
// Reference: https://nodejs.org/api/crypto.html#cryptocreatehmacalgorithm-key-options
export interface TokenIssuer {
  issueAccessToken(subject: string): string;
  readAccessToken(token: string): string;
}

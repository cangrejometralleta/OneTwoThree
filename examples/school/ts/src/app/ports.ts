import type { Course, CourseId, Page, Rut, Student, StudentId } from '../school/school.ts';

/** StudentStore Keeps Students wherever Students Live. Absence is undefined or false; only the Store Proves a Rut Taken. */
export interface StudentStore {
  insertStudentRow(student: Student): Promise<Student>;
  selectStudentRow(id: StudentId): Promise<Student | undefined>;
  updateStudentRow(student: Student): Promise<boolean>;
  deleteStudentRow(id: StudentId): Promise<boolean>;
  selectStudentPage(page: Page): Promise<Student[]>;
}

/** CourseStore Keeps Courses, and Answers whether one Exists. */
export interface CourseStore {
  insertCourseRow(course: Course): Promise<Course>;
  selectCourseRow(id: CourseId): Promise<Course | undefined>;
  selectCoursePage(page: Page): Promise<Course[]>;
}

/** RutRegistry Confirms whether a Rut Belongs to someone Real. A Denied Rut is false; a Registry that cannot Answer Throws. */
export interface RutRegistry {
  confirmRut(rut: Rut): Promise<boolean>;
}

/** EnrollmentNotifier Announces a Student the Store has Kept. */
export interface EnrollmentNotifier {
  announceEnrollment(student: Student): Promise<void>;
}

/** CallerReader Names the Caller a Token Proves, or Throws. */
export interface CallerReader {
  readCaller(token: string): string;
}

/** TokenIssuer Mints a Bearer Token and Reads one for the Guard. */
export interface TokenIssuer extends CallerReader {
  issueAccessToken(subject: string): string;
}

/** Logger Writes one Line in the Runtime's own Format. */
export type Logger = (line: string) => void;

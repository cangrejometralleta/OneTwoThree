import type { Course, Student } from '../school/index.ts';
export interface StudentStore {
  selectStudent(id: bigint): Promise<Student | undefined>;
  insertStudent(student: Student): Promise<Student>;
  updateStudent(id: bigint, student: Student): Promise<Student | undefined>;
  deleteStudent(id: bigint): Promise<boolean>;
  listStudents(): Promise<Student[]>;
}
export interface CourseStore {
  selectCourse(id: bigint): Promise<Course | undefined>;
  insertCourse(course: Omit<Course, 'id'>): Promise<Course>;
  listCourses(): Promise<Course[]>;
}
export interface TokenIssuer { issue(nowSeconds?: bigint): string; verify(token: string, nowSeconds?: bigint): string | undefined; }
export interface RutRegistry { confirm(rut: string): Promise<boolean>; }
export interface EnrollmentNotifier { announce(studentId: bigint): Promise<void>; }

import {
  RUT_TAKEN, UNNUMBERED_STUDENT, type Course, type CourseId, type Page, type Rut, type Student, type StudentId,
} from '../school/school.ts';
import type { CourseStore, EnrollmentNotifier, RutRegistry, StudentStore } from './ports.ts';

/**
 * Shelf Stands in for the Store: Course 1 is Open, and it Keeps what it is Given.
 * It Obeys the Store's Rules, so a Rut is Unique here too.
 */
export class Shelf implements StudentStore, CourseStore {
  readonly students = new Map<bigint, Student>();
  readonly courses = new Map<bigint, Course>([[1n, { id: 1n as CourseId, code: 'MAT-101', name: 'Algebra' }]]);
  page: Page | undefined;
  #next = 1n;

  async insertStudentRow(student: Student): Promise<Student> {
    if ([...this.students.values()].some((kept) => kept.rut === student.rut)) {
      throw RUT_TAKEN;
    }
    const kept = { ...student, id: this.#next++ as StudentId };
    this.students.set(kept.id, kept);

    return kept;
  }

  async selectStudentRow(id: StudentId): Promise<Student | undefined> {
    return this.students.get(id);
  }

  async updateStudentRow(student: Student): Promise<boolean> {
    if (!this.students.has(student.id) || student.id === UNNUMBERED_STUDENT) {
      return false;
    }
    this.students.set(student.id, student);

    return true;
  }

  async deleteStudentRow(id: StudentId): Promise<boolean> {
    return this.students.delete(id);
  }

  async selectStudentPage(page: Page): Promise<Student[]> {
    this.page = page;

    return [...this.students.values()];
  }

  async insertCourseRow(course: Course): Promise<Course> {
    const kept = { ...course, id: BigInt(this.courses.size + 1) as CourseId };
    this.courses.set(kept.id, kept);

    return kept;
  }

  async selectCourseRow(id: CourseId): Promise<Course | undefined> {
    return this.courses.get(id);
  }

  async selectCoursePage(page: Page): Promise<Course[]> {
    this.page = page;

    return [...this.courses.values()];
  }
}

/** FrontDesk Answers both Questions, the Way one Adapter may Fill two Ports. */
export class FrontDesk implements RutRegistry, EnrollmentNotifier {
  registered = true;
  down: Error | undefined;
  noticeFailure: Error | undefined;
  readonly asked: Rut[] = [];
  readonly announced: Student[] = [];

  async confirmRut(rut: Rut): Promise<boolean> {
    this.asked.push(rut);
    if (this.down !== undefined) {
      throw this.down;
    }

    return this.registered;
  }

  async announceEnrollment(student: Student): Promise<void> {
    this.announced.push(student);
    if (this.noticeFailure !== undefined) {
      throw this.noticeFailure;
    }
  }
}

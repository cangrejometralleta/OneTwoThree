import {
  checkCourseRecord, checkStudentRecord, COURSE_UNKNOWN, REGISTRY_UNAVAILABLE, RUT_UNREGISTERED, STUDENT_UNKNOWN,
  type Course, type CourseId, type Page, type Rut, type SchoolRules, type Student, type StudentId,
} from '../school/school.ts';
import type { CourseStore, EnrollmentNotifier, Logger, RutRegistry, StudentStore } from './ports.ts';

/** SchoolProviders is what the Composition Hands the Service. */
export interface SchoolProviders {
  readonly rules: SchoolRules;
  readonly students: StudentStore;
  readonly courses: CourseStore;
  readonly registry: RutRegistry | undefined;
  readonly notices: EnrollmentNotifier | undefined;
  readonly log: Logger;
}

/** SchoolService Joins the Providers behind each School Use Case. It never Learns that two Questions may have one Answerer. */
export class SchoolService {
  readonly #providers: SchoolProviders;

  constructor(providers: SchoolProviders) {
    this.#providers = providers;
  }

  /** listStudents Reads one Page of the Roll. */
  listStudents(page: Page): Promise<Student[]> {
    return this.#providers.students.selectStudentPage(page);
  }

  /** readStudent Finds one Enrolment, or it is Unknown. */
  async readStudent(id: StudentId): Promise<Student> {
    const student = await this.#providers.students.selectStudentRow(id);
    if (student === undefined) {
      throw STUDENT_UNKNOWN;
    }

    return student;
  }

  /** enrollStudent Validates and Keeps one Enrolment, then Announces it. An Announcement that Fails never Fails the Enrolment. */
  async enrollStudent(student: Student): Promise<Student> {
    await this.#checkStudent(student);
    const kept = await this.#providers.students.insertStudentRow(student);
    await this.#announceEnrollment(kept);

    return kept;
  }

  /** saveStudent Validates and Rewrites one Enrolment. It Announces nothing. */
  async saveStudent(student: Student): Promise<Student> {
    await this.#checkStudent(student);
    if (!(await this.#providers.students.updateStudentRow(student))) {
      throw STUDENT_UNKNOWN;
    }

    return student;
  }

  /** deleteStudent Ends one Enrolment, once. */
  async deleteStudent(id: StudentId): Promise<void> {
    if (!(await this.#providers.students.deleteStudentRow(id))) {
      throw STUDENT_UNKNOWN;
    }
  }

  /** listCourses Reads one Page of the Catalogue. */
  listCourses(page: Page): Promise<Course[]> {
    return this.#providers.courses.selectCoursePage(page);
  }

  /** readCourse Finds one Course, or it is Unknown. */
  async readCourse(id: CourseId): Promise<Course> {
    const course = await this.#providers.courses.selectCourseRow(id);
    if (course === undefined) {
      throw COURSE_UNKNOWN;
    }

    return course;
  }

  /** createCourse Validates and Keeps one Course. */
  async createCourse(course: Course): Promise<Course> {
    checkCourseRecord(course);

    return this.#providers.courses.insertCourseRow(course);
  }

  /** checkStudent Asks in the Order the Spec Names: Form, Course, Registry. */
  async #checkStudent(student: Student): Promise<void> {
    checkStudentRecord(student, this.#providers.rules);
    await this.readCourse(student.course);
    await this.#confirmRut(student.rut);
  }

  /** confirmRut Asks the Registry. An Unknown is never a Yes. */
  async #confirmRut(rut: Rut): Promise<void> {
    const registry = this.#providers.registry;
    if (registry === undefined) {
      throw REGISTRY_UNAVAILABLE;
    }

    let registered: boolean;
    try {
      registered = await registry.confirmRut(rut);
    } catch {
      throw REGISTRY_UNAVAILABLE;
    }
    if (!registered) {
      throw RUT_UNREGISTERED;
    }
  }

  /** announceEnrollment Keeps a Notice Failure inside the Answer: the Caller cannot Act on a Notice that did not Leave. */
  async #announceEnrollment(student: Student): Promise<void> {
    try {
      await this.#providers.notices?.announceEnrollment(student);
    } catch {
      this.#providers.log(`enrollment announcement failed for student ${student.id}`);
    }
  }
}

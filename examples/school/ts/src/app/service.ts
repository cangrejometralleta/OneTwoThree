import { Fault } from '../faults/index.ts';
import { courseFromBody, studentFromBody, type Course, type MinimumAge, type Student } from '../school/index.ts';
import type { CourseBody, StudentBody } from '../wire/index.ts';
import type { CourseStore, EnrollmentNotifier, RutRegistry, StudentStore } from './ports.ts';

export class SchoolService {
  private students: StudentStore; private courses: CourseStore; private registry: RutRegistry;
  private notifier: EnrollmentNotifier; private constants: MinimumAge; private log: (line: string) => void;
  constructor(students: StudentStore, courses: CourseStore, registry: RutRegistry,
    notifier: EnrollmentNotifier, constants: MinimumAge, log: (line: string) => void = console.error) {
    this.students=students; this.courses=courses; this.registry=registry; this.notifier=notifier; this.constants=constants; this.log=log;
  }
  async enroll(body: StudentBody): Promise<Student> {
    const candidate = studentFromBody(body, this.constants);
    if (!(await this.courses.selectCourse(candidate.courseId))) throw new Fault('CourseUnknown');
    if (!(await this.registry.confirm(candidate.rut))) throw new Fault('RutUnregistered');
    const kept = await this.students.insertStudent(candidate);
    try { await this.notifier.announce(kept.id); }
    catch { this.log(`student ${kept.id} is enrolled`); }
    return kept;
  }
  async save(id: bigint, body: StudentBody): Promise<Student> {
    const candidate = studentFromBody(body, this.constants);
    if (!(await this.courses.selectCourse(candidate.courseId))) throw new Fault('CourseUnknown');
    if (!(await this.registry.confirm(candidate.rut))) throw new Fault('RutUnregistered');
    const kept = await this.students.updateStudent(id, candidate);
    if (!kept) throw new Fault('StudentUnknown');
    return kept;
  }
  async readStudent(id: bigint): Promise<Student> { const s = await this.students.selectStudent(id); if (!s) throw new Fault('StudentUnknown'); return s; }
  async listStudents(): Promise<Student[]> { return this.students.listStudents(); }
  async deleteStudent(id: bigint): Promise<void> { if (!(await this.students.deleteStudent(id))) throw new Fault('StudentUnknown'); }
  async createCourse(body: CourseBody): Promise<Course> { return this.courses.insertCourse(courseFromBody(body)); }
  async readCourse(id: bigint): Promise<Course> { const c = await this.courses.selectCourse(id); if (!c) throw new Fault('CourseUnknown'); return c; }
  async listCourses(): Promise<Course[]> { return this.courses.listCourses(); }
}

export function statusForFault(kind: import('../faults/index.ts').FaultKind): number {
  const status: Record<import('../faults/index.ts').FaultKind, number> = {
    NameEmpty:400, CodeEmpty:400, RutInvalid:400, AgeTooLow:400, PageInvalid:400, BodyBroken:400,
    PathBroken:400, TokenInvalid:401, StudentUnknown:404, CourseUnknown:404, RouteUnknown:404,
    RutTaken:409, RutUnregistered:422, RegistryUnavailable:503
  };
  return status[kind];
}

import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { dirname } from 'node:path';
import { parseJson, stringifyJson } from '../codec/json.ts';
import type { Course, Student } from '../school/index.ts';
import type { CourseStore, StudentStore } from '../app/ports.ts';
import { Fault } from '../faults/index.ts';

type DiskState = { students: Array<{ id: string; rut: string; name: string; age: string; courseId: string }>;
  courses: Array<{ id: string; code: string; name: string }>; nextStudentId: string; nextCourseId: string };
const empty: DiskState = { students: [], courses: [], nextStudentId: '1', nextCourseId: '1' };
function utf8Order(a: string, b: string): number { return Buffer.compare(Buffer.from(a, 'utf8'), Buffer.from(b, 'utf8')); }
export class FileStore implements StudentStore, CourseStore {
  private state: DiskState = structuredClone(empty);
  private queue: Promise<unknown> = Promise.resolve();
  private path: string;
  private constructor(path: string) { this.path=path; }
  static async open(path: string): Promise<FileStore> {
    const store = new FileStore(path);
    try { store.state = parseJson(await readFile(path, 'utf8')) as DiskState; }
    catch (e) { if ((e as NodeJS.ErrnoException).code !== 'ENOENT') throw e; }
    return store;
  }
  private serial<T>(fn: () => Promise<T>): Promise<T> {
    const run = this.queue.then(fn); this.queue = run.then(() => undefined, () => undefined); return run;
  }
  private async persist(): Promise<void> { await mkdir(dirname(this.path), { recursive: true }); await writeFile(this.path, stringifyJson(this.state), 'utf8'); }
  async selectStudent(id: bigint): Promise<Student | undefined> { const s = this.state.students.find(x => x.id === String(id)); return s ? { ...s, id: BigInt(s.id), age: BigInt(s.age), courseId: BigInt(s.courseId) } : undefined; }
  async insertStudent(student: Student): Promise<Student> { return this.serial(async () => {
    if (this.state.students.some(s => s.rut === student.rut)) throw new Fault('RutTaken');
    const id = BigInt(this.state.nextStudentId); this.state.nextStudentId = String(id + 1n);
    const kept = { ...student, id };
    this.state.students.push({ ...kept, id: String(kept.id), age: String(kept.age), courseId: String(kept.courseId) }); await this.persist(); return kept;
  }); }
  async updateStudent(id: bigint, student: Student): Promise<Student | undefined> { return this.serial(async () => {
    const idx = this.state.students.findIndex(s => s.id === String(id)); if (idx < 0) return undefined;
    if (this.state.students.some((s, i) => i !== idx && s.rut === student.rut)) throw new Fault('RutTaken');
    const kept = { ...student, id }; this.state.students[idx] = { ...kept, id: String(id), age: String(kept.age), courseId: String(kept.courseId) }; await this.persist(); return kept;
  }); }
  async deleteStudent(id: bigint): Promise<boolean> { return this.serial(async () => { const i = this.state.students.findIndex(s => s.id === String(id)); if (i < 0) return false; this.state.students.splice(i, 1); await this.persist(); return true; }); }
  async listStudents(): Promise<Student[]> {
    return [...this.state.students].sort((a,b) => utf8Order(a.name,b.name) || (BigInt(a.id) < BigInt(b.id) ? -1 : 1))
      .map(s => ({ ...s, id: BigInt(s.id), age: BigInt(s.age), courseId: BigInt(s.courseId) }));
  }
  async selectCourse(id: bigint): Promise<Course | undefined> { const c = this.state.courses.find(x => x.id === String(id)); return c ? { ...c, id: BigInt(c.id) } : undefined; }
  async insertCourse(course: Omit<Course, 'id'>): Promise<Course> { return this.serial(async () => {
    const id = BigInt(this.state.nextCourseId); this.state.nextCourseId = String(id + 1n);
    const kept: Course = { ...course, id };
    this.state.courses.push({ ...kept, id: String(kept.id) }); await this.persist(); return kept;
  }); }
  async listCourses(): Promise<Course[]> {
    return [...this.state.courses].sort((a,b) => utf8Order(a.code,b.code) || (BigInt(a.id) < BigInt(b.id) ? -1 : 1)).map(c => ({ ...c, id: BigInt(c.id) }));
  }
}

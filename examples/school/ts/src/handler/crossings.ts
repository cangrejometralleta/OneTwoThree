import type { BodyOf } from '../app/validation.ts';
import { readRut, UNNUMBERED_COURSE, type Course, type CourseId, type Student, type StudentId } from '../school/school.ts';
import type { CourseView, StudentView } from '../wire/wire.ts';

// The Wire never Becomes the Business by Accident: these are the Crossings.

export const STUDENT_FIELDS = { rut: 'text', name: 'text', age: 'whole', courseId: 'whole' } as const;
export const COURSE_FIELDS = { code: 'text', name: 'text' } as const;

export function buildStudentRecord(body: BodyOf<typeof STUDENT_FIELDS>, id: StudentId): Student {
  return { id, rut: readRut(body.rut), name: body.name, age: body.age, course: body.courseId as CourseId };
}

export function renderStudentView(student: Student): StudentView {
  return { id: student.id, rut: student.rut, name: student.name, age: student.age, courseId: student.course };
}

export function buildCourseRecord(body: BodyOf<typeof COURSE_FIELDS>): Course {
  return { id: UNNUMBERED_COURSE, code: body.code, name: body.name };
}

export function renderCourseView(course: Course): CourseView {
  return { id: course.id, code: course.code, name: course.name };
}

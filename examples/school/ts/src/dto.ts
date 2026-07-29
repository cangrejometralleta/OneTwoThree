import type {
  Course, CourseCode, CourseID, FullName, RUT, Student, StudentID,
} from "./domain.js";

// StudentBody is what a Client Sends.
// Weak on Purpose: the Wire Cannot be Trusted,
// so nothing here is a Business Type yet.
export type StudentBody = {
  rut: string;
  name: string;
  age: number;
  courseId: number;
};

// StudentView is what a Client Gets back.
export type StudentView = {
  id: number;
  rut: string;
  name: string;
  age: number;
  courseId: number;
};

export type CourseBody = { code: string; name: string };
export type CourseView = { id: number; code: string; name: string };

// buildStudentRecord Promotes untrusted Input into Business Types.
export function buildStudentRecord(body: StudentBody, id: number): Student {
  return {
    id: id as StudentID,
    rut: String(body.rut ?? "") as RUT,
    name: String(body.name ?? "") as FullName,
    age: Number(body.age ?? 0),
    course: Number(body.courseId ?? 0) as CourseID,
  };
}

// buildCourseRecord Promotes untrusted Course Input.
export function buildCourseRecord(body: CourseBody, id: number): Course {
  return {
    id: id as CourseID,
    code: String(body.code ?? "") as CourseCode,
    name: String(body.name ?? "") as FullName,
  };
}

// renderStudentView Demotes a Business Student back to the Wire.
export function renderStudentView(student: Student): StudentView {
  return {
    id: student.id,
    rut: student.rut,
    name: student.name,
    age: student.age,
    courseId: student.course,
  };
}

// renderCourseView Demotes a Business Course back to the Wire.
export function renderCourseView(course: Course): CourseView {
  return { id: course.id, code: course.code, name: course.name };
}

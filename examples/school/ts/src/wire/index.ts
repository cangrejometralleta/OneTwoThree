export type StudentBody = { rut: string; name: string; age: bigint; courseId: bigint };
export type StudentView = StudentBody & { id: bigint };
export type CourseBody = { code: string; name: string };
export type CourseView = CourseBody & { id: bigint };
export type TokenView = { token: string };
export type FailureView = { error: string };

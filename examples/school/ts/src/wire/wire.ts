// This Module is the Contract, and Imports nothing.
// Open one File and you Know the whole API. Whole Numbers are 64 bit, so they are bigint.

/** StudentBody is what a Client Sends. Weak on Purpose: nothing here is a Business Type yet. */
export interface StudentBody {
  readonly rut: string;
  readonly name: string;
  readonly age: bigint;
  readonly courseId: bigint;
}

/** StudentView is what a Client Gets back. */
export interface StudentView {
  readonly id: bigint;
  readonly rut: string;
  readonly name: string;
  readonly age: bigint;
  readonly courseId: bigint;
}

/** CourseBody is what a Client Sends about a Course. */
export interface CourseBody {
  readonly code: string;
  readonly name: string;
}

/** CourseView is what a Client Gets back about a Course. */
export interface CourseView {
  readonly id: bigint;
  readonly code: string;
  readonly name: string;
}

/** TokenView is what the Door Hands out for a Token. */
export interface TokenView {
  readonly token: string;
}

/** FailureView is the one Shape every Failure Answers with. */
export interface FailureView {
  readonly error: string;
}

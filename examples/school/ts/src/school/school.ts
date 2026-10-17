import { refuseInvalidInput, reportMissingRecord, reportTakenValue, refuseUnprocessableValue, reportUnavailableProvider, type Fault } from '../faults/faults.ts';

/** StudentId is the Identity the Store Gives a Student; the Body never Carries it. */
export type StudentId = bigint & { readonly brand: 'StudentId' };

/** CourseId is the Identity the Store Gives a Course. */
export type CourseId = bigint & { readonly brand: 'CourseId' };

/** Rut is the Chilean tax Number a Student is Known by, in lower case. */
export type Rut = string & { readonly brand: 'Rut' };

/** Student is the Business Truth about one Enrolment. */
export interface Student {
  readonly id: StudentId;
  readonly rut: Rut;
  readonly name: string;
  readonly age: bigint;
  readonly course: CourseId;
}

/** Course is the Business Truth about one Class. Codes are not Unique. */
export interface Course {
  readonly id: CourseId;
  readonly code: string;
  readonly name: string;
}

/** SchoolRules Holds the Global Constants the Business will not Bend. No Deployment Overrides them. */
export interface SchoolRules {
  readonly minimumAgeYears: bigint;
}

/** The Business Fails in named Ways, never in Numbers. Each one Declares its Answer, once and here. */
export const NAME_IS_EMPTY = refuseInvalidInput('name is Empty');
export const CODE_IS_EMPTY = refuseInvalidInput('code is Empty');
export const RUT_IS_INVALID = refuseInvalidInput('rut Fails its Check Digit');
export const PAGE_IS_INVALID = refuseInvalidInput('page or size is not valid');
export const STUDENT_UNKNOWN = reportMissingRecord('student not Found');
export const COURSE_UNKNOWN = reportMissingRecord('course not Found');
export const RUT_TAKEN = reportTakenValue('rut is already Registered');
export const RUT_UNREGISTERED = refuseUnprocessableValue('rut is not in the Registry');
export const REGISTRY_UNAVAILABLE = reportUnavailableProvider('registry cannot Answer now');

/** ageIsTooLow Spells the Floor from the Global Constant. */
export const ageIsTooLow = (rules: SchoolRules): Fault => refuseInvalidInput(`age Must be ${rules.minimumAgeYears} or more`);

/** UNNUMBERED Stands for an Identity the Store has not Given yet. */
export const UNNUMBERED_STUDENT = 0n as StudentId;
export const UNNUMBERED_COURSE = 0n as CourseId;

/** RUT_WEIGHTS Cycle two through seven, Read right to left. */
export const RUT_WEIGHTS = [2, 3, 4, 5, 6, 7] as const;

/** CHECK_MODULUS is the Modulo the Check Digit Encodes. */
export const CHECK_MODULUS = 11;

const RUT_SHAPE = /^[0-9]+-[0-9kK]$/;

/** readRut Promotes Text into a Rut in lower case: 12345-K and 12345-k are one. Leading zeros are not Merged. */
export const readRut = (text: string): Rut => text.toLowerCase() as Rut;

/** isEmptyText Holds for a Text of only whitespace; the Text itself is never Trimmed. */
export const isEmptyText = (text: string): boolean => text.trim() === '';

/**
 * rutLooksValid Runs the Modulo eleven Check the Digit Encodes.
 * Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
 */
export function rutLooksValid(rut: string): boolean {
  if (!RUT_SHAPE.test(rut)) {
    return false;
  }

  const [body = '', digit = ''] = rut.toLowerCase().split('-');

  return digit === checkDigitFor(body);
}

/** checkDigitFor Folds a Rut Body into its single Check Character, Digit by Digit. */
export function checkDigitFor(body: string): string {
  let sum = 0;
  for (let position = 0; position < body.length; position++) {
    const digit = Number(body[body.length - 1 - position]);
    sum = (sum + digit * (RUT_WEIGHTS[position % RUT_WEIGHTS.length] ?? 0)) % CHECK_MODULUS;
  }

  return nameCheckRemainder(CHECK_MODULUS - sum);
}

/** nameCheckRemainder Turns a Remainder into the Character it Means: eleven is 0, ten is k. */
export function nameCheckRemainder(remainder: number): string {
  if (remainder === 11) {
    return '0';
  }

  return remainder === 10 ? 'k' : String(remainder);
}

/** checkStudentRecord Refuses a Student the Business cannot Use: Name, then Rut, then Age. */
export function checkStudentRecord(student: Student, rules: SchoolRules): void {
  if (isEmptyText(student.name)) {
    throw NAME_IS_EMPTY;
  }
  if (!rutLooksValid(student.rut)) {
    throw RUT_IS_INVALID;
  }
  if (student.age < rules.minimumAgeYears) {
    throw ageIsTooLow(rules);
  }
}

/** checkCourseRecord Refuses a Course the Business cannot Use: Code, then Name. */
export function checkCourseRecord(course: Course): void {
  if (isEmptyText(course.code)) {
    throw CODE_IS_EMPTY;
  }
  if (isEmptyText(course.name)) {
    throw NAME_IS_EMPTY;
  }
}

/** LARGEST_WHOLE is the largest Number 64 bits can Hold. */
export const LARGEST_WHOLE = (1n << 63n) - 1n;

/** DEFAULT_PAGE_SIZE Windows a Page when the Caller Names a page but no size. */
export const DEFAULT_PAGE_SIZE = 10n;

/** Page Windows a List, Zero based. */
export interface Page {
  readonly number: bigint;
  readonly size: bigint;
}

/** WHOLE_SET Stands for a List no Window Cuts. */
export const WHOLE_SET: Page = { number: 0n, size: LARGEST_WHOLE };

/** buildPage Holds the Bounds before a Page Exists: Number at least zero, Size at least one. */
export function buildPage(number: bigint, size: bigint): Page {
  if (number < 0n || size < 1n) {
    throw PAGE_IS_INVALID;
  }

  return { number, size };
}

/** pageOffset Counts the Records before the Window, as far as 64 bits can Count. */
export function pageOffset(page: Page): bigint {
  const offset = page.number * page.size;

  return offset > LARGEST_WHOLE ? LARGEST_WHOLE : offset;
}

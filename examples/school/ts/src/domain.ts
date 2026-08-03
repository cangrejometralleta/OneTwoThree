// Brand Gives TypeScript what Go Gets from a named Type:
// a String that Refuses to Stand in for another String.
declare const brand: unique symbol;
type Brand<T, Name> = T & { readonly [brand]: Name };

export type StudentID = Brand<number, "StudentID">;
export type CourseID = Brand<number, "CourseID">;

// RUT is the Chilean Tax Identity, Digit included.
export type RUT = Brand<string, "RUT">;
export type FullName = Brand<string, "FullName">;
export type CourseCode = Brand<string, "CourseCode">;

// Student is the Business Truth about one Enrolment.
export type Student = {
  id: StudentID;
  rut: RUT;
  name: FullName;
  age: number;
  course: CourseID;
};

// Course is the Business Truth about one Class.
export type Course = {
  id: CourseID;
  code: CourseCode;
  name: FullName;
};

// BusinessError Carries a Name the Transport can Map.
export class BusinessError extends Error {
  constructor(readonly kind: "invalid" | "absent" | "taken", message: string) {
    super(message);
  }
}

// The Business Fails in named Ways, never in Numbers.
export const ErrNameIsEmpty = new BusinessError("invalid", "name is Empty");
export const ErrRutIsInvalid = new BusinessError("invalid", "rut Fails its Check Digit");
export const ErrRutTaken = new BusinessError("taken", "rut is already Registered");
export const ErrAgeIsTooLow = new BusinessError("invalid", "age Must be eighteen or more");
export const ErrStudentUnknown = new BusinessError("absent", "student not Found");
export const ErrCourseUnknown = new BusinessError("absent", "course not Found");
export const ErrPageIsInvalid = new BusinessError("invalid", "page Numbers Must not be negative");
export const ErrBodyIsBroken = new BusinessError("invalid", "body is not valid JSON");
export const ErrPathIsBroken = new BusinessError("invalid", "path Holds no valid Identity");
export const ErrTokenIsInvalid = new BusinessError("invalid", "token is Invalid or Expired");

const rutShape = /^[0-9]+-[0-9kK]$/;

// RutWeights Cycles two through seven, Read right to left.
const RutWeights = [2, 3, 4, 5, 6, 7];

// checkDigitFor Folds a RUT Body into its single Check Character.
// Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
export function checkDigitFor(body: string): string {
  let sum = 0;
  let index = 0;

  for (let number = Number(body); number > 0; number = Math.floor(number / 10)) {
    sum += (number % 10) * RutWeights[index % RutWeights.length]!;
    index += 1;
  }

  return nameCheckRemainder(11 - (sum % 11));
}

// nameCheckRemainder Turns a Remainder into the Character it Means.
function nameCheckRemainder(remainder: number): string {
  if (remainder === 11) return "0";
  if (remainder === 10) return "k";
  return String(remainder);
}

// rutLooksValid Runs the Modulo eleven Check the Digit Encodes.
export function rutLooksValid(rut: string): boolean {
  if (!rutShape.test(rut)) return false;

  const [body, digit] = rut.toLowerCase().split("-");

  return digit === checkDigitFor(body!);
}

// checkStudentRecord Refuses a Student the Business cannot Use.
export function checkStudentRecord(student: Student): BusinessError | null {
  if (student.name.length === 0) return ErrNameIsEmpty;
  if (!rutLooksValid(student.rut)) return ErrRutIsInvalid;
  if (student.age < 18) return ErrAgeIsTooLow;

  return null;
}

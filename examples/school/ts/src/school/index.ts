import { Fault } from '../faults/index.ts';

export type Student = Readonly<{ id: bigint; rut: string; name: string; age: bigint; courseId: bigint }>;
export type Course = Readonly<{ id: bigint; code: string; name: string }>;
export type MinimumAge = Readonly<{ minimumAgeYears: bigint }>;
export function validRut(value: string): boolean {
  if (!/^[0-9]+-[0-9kK]$/.test(value)) return false;
  const parts = value.toLowerCase().split('-');
  const digits = parts[0]!;
  const check = parts[1]!;
  let sum = 0, weight = 2;
  for (let i = digits.length - 1; i >= 0; i--) { sum += Number(digits[i]) * weight; weight = weight === 7 ? 2 : weight + 1; }
  const remainder = 11 - (sum % 11);
  const expected = remainder === 11 ? '0' : remainder === 10 ? 'k' : String(remainder);
  return check === expected;
}
export function studentFromBody(body: { rut: string; name: string; age: bigint; courseId: bigint }, constants: MinimumAge): Student {
  if (!body.name.trim()) throw new Fault('NameEmpty');
  if (!validRut(body.rut)) throw new Fault('RutInvalid');
  if (BigInt(body.age) < constants.minimumAgeYears) throw new Fault('AgeTooLow', `age Must be ${constants.minimumAgeYears} or more`);
  return Object.freeze({ id: 0n, rut: body.rut.toLowerCase(), name: body.name, age: BigInt(body.age), courseId: BigInt(body.courseId) });
}
export function courseFromBody(body: { code: string; name: string }): Omit<Course, 'id'> {
  if (!body.code.trim()) throw new Fault('CodeEmpty');
  if (!body.name.trim()) throw new Fault('NameEmpty');
  return Object.freeze({ code: body.code, name: body.name });
}

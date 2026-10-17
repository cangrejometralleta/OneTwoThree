import { readFile } from 'node:fs/promises';

import type { EnrollmentNotifier, Logger, RutRegistry } from '../app/ports.ts';
import { readJsonText } from '../codec/json.ts';
import { readRut, rutLooksValid, type Rut, type Student } from '../school/school.ts';

/**
 * Office Answers the Registry and Announces Enrolments: one Adapter, two Ports.
 * The Registry is a JSON File of the Ruts that are Real; the Notice goes to the Log.
 */
export class Office implements RutRegistry, EnrollmentNotifier {
  readonly #registered: ReadonlySet<Rut>;
  readonly #log: Logger;

  private constructor(registered: ReadonlySet<Rut>, log: Logger) {
    this.#registered = registered;
    this.#log = log;
  }

  /**
   * openOffice Reads the Registry File once, at Startup. A File that is Missing, Empty, not an Array,
   * Holds a second Value or a Text that is not a valid Rut Stops Startup.
   */
  static async openOffice(path: string, log: Logger): Promise<Office> {
    const read = readJsonText(await readFile(path));
    if (!Array.isArray(read)) {
      throw new Error(`registry file ${path} Must Hold one Array of Texts`);
    }

    const registered = new Set<Rut>();
    for (const item of read) {
      if (typeof item !== 'string' || !rutLooksValid(item)) {
        throw new Error(`registry file ${path}: ${String(item)} is not a valid RUT`);
      }
      registered.add(readRut(item));
    }

    return new Office(registered, log);
  }

  /** confirmRut Answers whether the Registry Knows a Rut. This one never Fails to Answer. */
  async confirmRut(rut: Rut): Promise<boolean> {
    return this.#registered.has(rut);
  }

  /** announceEnrollment Logs the Notice, Naming the Student by Id and never by Rut or Name. */
  async announceEnrollment(student: Student): Promise<void> {
    this.#log(`student ${student.id} is enrolled`);
  }
}

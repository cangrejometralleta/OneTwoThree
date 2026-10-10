import { readFile } from 'node:fs/promises';
import { parseJson } from '../codec/json.ts';
import { validRut } from '../school/index.ts';
import type { EnrollmentNotifier, RutRegistry } from '../app/ports.ts';
export class Office implements RutRegistry, EnrollmentNotifier {
  private ruts: Set<string>; private logger: (line: string) => void;
  private constructor(ruts: Set<string>, logger: (line: string) => void) { this.ruts=ruts; this.logger=logger; }
  static async open(path: string, logger: (line: string) => void = console.error): Promise<Office> {
    const raw = await readFile(path, 'utf8');
    if (!raw.trim()) throw new Error('registry file is empty');
    const values = parseJson(raw);
    if (!Array.isArray(values)) throw new Error('registry must be an array');
    const ruts = new Set<string>();
    for (const value of values) { if (typeof value !== 'string' || !validRut(value)) throw new Error('registry contains an invalid RUT'); ruts.add(value.toLowerCase()); }
    return new Office(ruts, logger);
  }
  async confirm(rut: string): Promise<boolean> { return this.ruts.has(rut); }
  async announce(studentId: bigint): Promise<void> { this.logger(`student ${studentId} is enrolled`); }
}

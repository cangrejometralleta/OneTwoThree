import assert from 'node:assert/strict';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { test } from 'node:test';

import type { EnrollmentNotifier, RutRegistry } from '../app/ports.ts';
import { readRut, type CourseId, type StudentId } from '../school/school.ts';
import { Office } from './office.ts';

const folder = mkdtempSync(join(tmpdir(), 'office-'));
const writeRegistry = (name: string, content: string): string => {
  const path = join(folder, `${name.replaceAll(' ', '-')}.json`);
  writeFileSync(path, content);

  return path;
};

test('oneOfficeFillsBothPorts', async () => {
  const office: RutRegistry & EnrollmentNotifier = await Office.openOffice(writeRegistry('both', '[]'), () => {});

  assert.ok(office);
});

test('confirmRutKnowsOnlyWhatTheFileHoldsInLowerCase', async () => {
  const office = await Office.openOffice(writeRegistry('known', '["11111111-1","10000013-K","11111111-1"]'), () => {});

  assert.equal(await office.confirmRut(readRut('11111111-1')), true);
  assert.equal(await office.confirmRut(readRut('10000013-k')), true);
  assert.equal(await office.confirmRut(readRut('22222222-2')), false);
});

test('openOfficeStopsOnWhatItCannotTrust', async () => {
  const broken = {
    empty: '', 'broken json': '["11111111-1"', 'not an array': '{"ruts": []}', null: 'null', 'a second value': '[] []',
    'a bad check': '["11111111-9"]', 'a non rut text': '["hello"]', 'not a text': '[11111111]',
  };

  for (const [name, content] of Object.entries(broken)) {
    await assert.rejects(Office.openOffice(writeRegistry(name, content), () => {}), name);
  }
  await assert.rejects(Office.openOffice(join(folder, 'absent.json'), () => {}));
});

test('theShippedRegistryHoldsOnlyRealRuts', async () => {
  const office = await Office.openOffice(join('..', 'config', 'registry.json'), () => {});

  assert.equal(await office.confirmRut(readRut('12345678-5')), true);
});

test('anAnnouncementNamesTheStudentByIdOnly', async () => {
  const lines: string[] = [];
  const office = await Office.openOffice(writeRegistry('quiet', '[]'), (line) => lines.push(line));

  await office.announceEnrollment({ id: 7n as StudentId, rut: readRut('11111111-1'), name: 'Ana', age: 20n, course: 1n as CourseId });

  assert.deepEqual(lines, ['student 7 is enrolled']);
});

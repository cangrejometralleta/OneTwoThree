import assert from 'node:assert/strict';
import { copyFileSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { test } from 'node:test';

import { loadSchoolConfig } from './school-config.ts';
import { findSchoolDataRoot, readGlobalConstants, SettingsRefused } from './settings.ts';

const ADAPTERS = ['express', 'stdlib'];
const SHIPPED = '..';

/** buildConfigFixture Keeps broken Files away from the checked-in Examples. */
function buildConfigFixture(): string {
  const root = mkdtempSync(join(tmpdir(), 'config-'));
  mkdirSync(join(root, 'config'));
  mkdirSync(join(root, 'constants'));
  for (const name of ['config/school.defaults.json', 'config/school.development.json', 'config/school.production.json',
    'config/registry.json', 'constants/school.json']) {
    copyFileSync(join(SHIPPED, name), join(root, name));
  }

  return root;
}

const readyEnvironment = (): Record<string, string> => ({ APP_ENV: 'production', TOKEN_SECRET: 'test-secret' });
const refused = (failure: unknown) => failure instanceof SettingsRefused;

test('theShippedFilesLoadForBothEnvironments', () => {
  const root = buildConfigFixture();
  assert.equal(loadSchoolConfig(root, readyEnvironment(), ADAPTERS).databasePath, 'school.production.db');

  const development = loadSchoolConfig(root, { ...readyEnvironment(), APP_ENV: 'development' }, ADAPTERS);
  assert.deepEqual([development.databasePath, development.port, development.tokenLifeSeconds, development.serverAdapter],
    ['school.db', 8080, 600n, 'stdlib']);
});

test('defaultsThenTheEnvironmentFileThenTheVariables', () => {
  const root = buildConfigFixture();
  writeFileSync(join(root, 'config', 'school.production.json'), '{"databasePath":"production.db","port":9000}');

  const fromFile = loadSchoolConfig(root, readyEnvironment(), ADAPTERS);
  assert.deepEqual([fromFile.port, fromFile.databasePath], [9000, 'production.db']);

  const fromVariables = loadSchoolConfig(root, { ...readyEnvironment(), SCHOOL_PORT: '9100', SERVER: 'express' }, ADAPTERS);
  assert.deepEqual([fromVariables.port, fromVariables.serverAdapter], [9100, 'express']);
});

test('theRegistryResolvesFromTheDataRootAndMustExist', () => {
  const root = buildConfigFixture();
  assert.equal(loadSchoolConfig(root, readyEnvironment(), ADAPTERS).registryPath, join(root, 'config', 'registry.json'));

  assert.throws(() => loadSchoolConfig(root, { ...readyEnvironment(), SCHOOL_REGISTRY_PATH: join(root, 'absent.json') }, ADAPTERS), refused);
  assert.throws(() => loadSchoolConfig(root, { ...readyEnvironment(), SCHOOL_REGISTRY_PATH: 'config' }, ADAPTERS), refused);
});

test('configRejectsInvalidInputs', () => {
  const cases: Record<string, Record<string, string>> = {
    'missing environment': { APP_ENV: '' }, 'unknown environment': { APP_ENV: 'staging' },
    'path traversal': { APP_ENV: '../production' }, 'empty override': { SCHOOL_PORT: '' },
    'signed port': { SCHOOL_PORT: '+80' }, 'spaced port': { SCHOOL_PORT: ' 80' }, 'fractional port': { SCHOOL_PORT: '80.0' },
    'port too high': { SCHOOL_PORT: '65536' }, 'port zero': { SCHOOL_PORT: '0' }, 'empty path': { SCHOOL_DATABASE_PATH: ' ' },
    'empty registry path': { SCHOOL_REGISTRY_PATH: '' }, 'global variable': { SCHOOL_MINIMUM_AGE_YEARS: '1' },
    'unknown variable': { SCHOOL_TYPO: '1' }, 'missing secret': { TOKEN_SECRET: '' }, 'unknown adapter': { SERVER: 'gin' },
    'empty adapter': { SERVER: '' }, 'lifetime zero': { SCHOOL_TOKEN_LIFE_SECONDS: '0' },
    'lifetime beyond a day': { SCHOOL_TOKEN_LIFE_SECONDS: '86401' },
  };

  for (const [name, overrides] of Object.entries(cases)) {
    assert.throws(() => loadSchoolConfig(buildConfigFixture(), { ...readyEnvironment(), ...overrides }, ADAPTERS), refused, name);
  }
  assert.throws(() => loadSchoolConfig(buildConfigFixture(), { APP_ENV: 'production' }, ADAPTERS), refused);
});

test('configRejectsInvalidFiles', () => {
  const cases: Record<string, string> = {
    'global file key': '{"databasePath":"test.db","minimumAgeYears":1}', 'unknown file key': '{"databasePath":"test.db","typo":1}',
    'null key': '{"databasePath":null}', 'invalid type': '{"databasePath":"test.db","port":"9000"}',
    fractional: '{"databasePath":"test.db","port":9000.5}', 'invalid file range': '{"databasePath":"test.db","port":0}',
    'invalid json': '{', 'null document': 'null', 'array document': '[]', 'empty file': '', 'second value': '{} {}',
    'secret in file': '{"databasePath":"test.db","tokenSecret":"forbidden"}',
  };

  for (const [name, content] of Object.entries(cases)) {
    const root = buildConfigFixture();
    writeFileSync(join(root, 'config', 'school.production.json'), content);
    assert.throws(() => loadSchoolConfig(root, readyEnvironment(), ADAPTERS), refused, name);
  }

  const root = buildConfigFixture();
  rmSync(join(root, 'config', 'school.production.json'));
  assert.throws(() => loadSchoolConfig(root, readyEnvironment(), ADAPTERS), refused);
});

test('aMissingKeyStopsStartup', () => {
  const root = buildConfigFixture();
  writeFileSync(join(root, 'config', 'school.defaults.json'), '{"port":8080}');
  writeFileSync(join(root, 'config', 'school.production.json'), '{}');

  assert.throws(() => loadSchoolConfig(root, readyEnvironment(), ADAPTERS), refused);
});

test('anOverrideNeverHidesAnInvalidDefault', () => {
  const root = buildConfigFixture();
  writeFileSync(join(root, 'config', 'school.defaults.json'), '{"port":null}');

  assert.throws(() => loadSchoolConfig(root, { ...readyEnvironment(), SCHOOL_PORT: '9100' }, ADAPTERS), refused);
});

test('aVariableOutsideTheNamespaceIsIgnored', () => {
  assert.equal(loadSchoolConfig(buildConfigFixture(), { ...readyEnvironment(), SCHOOLING: 'x', PATH: '/usr/bin' }, ADAPTERS).port, 8080);
});

test('theGlobalConstantsHoldTheEnrolmentAgeAndRefuseBrokenData', () => {
  const root = buildConfigFixture();
  assert.equal(readGlobalConstants(root).minimumAgeYears, 18n);

  for (const broken of ['{"minimumAgeYears":null}', '{}', '{"minimumAgeYears":0}', '{"minimumAgeYears":18,"port":1}']) {
    writeFileSync(join(root, 'constants', 'school.json'), broken);
    assert.throws(() => readGlobalConstants(root), refused, broken);
  }
});

test('theDataRootIsFoundAboveTheWorkingDirectory', () => {
  assert.equal(findSchoolDataRoot('.'), resolve('..'));
});

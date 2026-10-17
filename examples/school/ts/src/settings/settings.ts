import { existsSync, readFileSync, statSync } from 'node:fs';
import { dirname, isAbsolute, join, resolve } from 'node:path';

import { readJsonText, type JsonValue } from '../codec/json.ts';

/** SettingsRefused Stops Startup and Says what Broke. */
export class SettingsRefused extends Error {}

/** ValueRule Validates one declared Configuration Key, or Throws what it Wanted. */
export type ValueRule = (value: JsonValue) => void;

/** checkIntegerRange Keeps JSON Numbers whole and Bounded. */
export const checkIntegerRange = (minimum: bigint, maximum: bigint): ValueRule => (value) => {
  if (typeof value !== 'bigint') {
    throw new SettingsRefused('Must be an Integer');
  }
  if (value < minimum || value > maximum) {
    throw new SettingsRefused(`Must be between ${minimum} and ${maximum}`);
  }
};

/** checkTextValue Requires a nonempty String. */
export const checkTextValue: ValueRule = (value) => {
  if (typeof value !== 'string') {
    throw new SettingsRefused('Must be a String');
  }
  if (value.trim() === '') {
    throw new SettingsRefused('Must not be Empty');
  }
};

/** checkOneOf Accepts only the Names it is Given. */
export const checkOneOf = (names: readonly string[]): ValueRule => (value) => {
  if (typeof value !== 'string' || !names.includes(value)) {
    throw new SettingsRefused(`Must be one of ${names.join(', ')}`);
  }
};

function checkValue(rule: ValueRule, value: JsonValue, where: string): void {
  try {
    rule(value);
  } catch (refused) {
    throw new SettingsRefused(`${where} ${(refused as Error).message}`);
  }
}

/** readDataFile Rejects Unknown Keys and Invalid Values before Merging. */
export function readDataFile(path: string, rules: Readonly<Record<string, ValueRule>>): Map<string, JsonValue> {
  let data: Buffer;
  try {
    data = readFileSync(path);
  } catch {
    throw new SettingsRefused(`${path} cannot be Read`);
  }

  const values = readJsonText(data);
  if (!(values instanceof Map)) {
    throw new SettingsRefused(`${path} Must Hold one Object`);
  }

  for (const [key, value] of values) {
    const rule = rules[key];
    if (rule === undefined) {
      throw new SettingsRefused(`${path}: Unknown Key ${key}`);
    }
    checkValue(rule, value, `${path}: ${key}`);
  }

  return values;
}

/** requireEveryKey Refuses a Set of Values that Lacks a declared Key. */
export function requireEveryKey(values: ReadonlyMap<string, JsonValue>, rules: Readonly<Record<string, ValueRule>>): void {
  for (const key of Object.keys(rules)) {
    if (!values.has(key)) {
      throw new SettingsRefused(`${key} is Missing`);
    }
  }
}

/** GlobalConstants Holds the Values no Deployment can Override, Read from constants/. */
export interface GlobalConstants {
  readonly minimumAgeYears: bigint;
}

/** readGlobalConstants Loads Constants without an Override Path; no Variable Reaches them. */
export function readGlobalConstants(root: string): GlobalConstants {
  const rules = { minimumAgeYears: checkIntegerRange(1n, 200n) };
  const values = readDataFile(join(root, 'constants', 'school.json'), rules);
  requireEveryKey(values, rules);

  return { minimumAgeYears: values.get('minimumAgeYears') as bigint };
}

/** readEnvironmentFiles Applies the Defaults, then the selected Environment. Both Files must Exist. */
export function readEnvironmentFiles(root: string, concern: string, environment: string | undefined,
  rules: Readonly<Record<string, ValueRule>>): Map<string, JsonValue> {
  if (environment !== 'development' && environment !== 'production') {
    throw new SettingsRefused('APP_ENV Must be development or production');
  }

  const values = readDataFile(join(root, 'config', `${concern}.defaults.json`), rules);
  for (const [key, value] of readDataFile(join(root, 'config', `${concern}.${environment}.json`), rules)) {
    values.set(key, value);
  }

  return values;
}

/** checkEnvironmentKeys Refuses undeclared Variables in this Service's Namespace. */
export function checkEnvironmentKeys(environment: Readonly<Record<string, string | undefined>>, prefix: string,
  names: Readonly<Record<string, string>>): void {
  for (const name of Object.keys(environment)) {
    if (name.startsWith(prefix) && !(name in names)) {
      throw new SettingsRefused(`Unknown Variable ${name}; Global Constants Cannot be Overridden`);
    }
  }
}

const DIGITS = /^[0-9]+$/;

/** applyEnvironmentValues Validates declared Variables before Replacing Values. Set and Empty is Set. */
export function applyEnvironmentValues(values: Map<string, JsonValue>, environment: Readonly<Record<string, string | undefined>>,
  names: Readonly<Record<string, string>>, rules: Readonly<Record<string, ValueRule>>): void {
  for (const [name, key] of Object.entries(names)) {
    const text = environment[name];
    const rule = rules[key];
    if (text === undefined || rule === undefined) {
      continue;
    }

    try {
      rule(text);
      values.set(key, text);
    } catch (notText) {
      if (!DIGITS.test(text)) {
        throw new SettingsRefused(`${name}: ${(notText as Error).message}`);
      }
      checkValue(rule, BigInt(text), name);
      values.set(key, BigInt(text));
    }
  }
}

/** DATA_ROOT_SEARCH_LEVELS Bounds the Walk up from the working Directory. */
const DATA_ROOT_SEARCH_LEVELS = 5;

/** findSchoolDataRoot Walks up from a Directory until constants/ Appears. */
export function findSchoolDataRoot(start: string): string {
  let root = resolve(start);
  for (let level = 0; level < DATA_ROOT_SEARCH_LEVELS; level++) {
    const constants = join(root, 'constants', 'school.json');
    if (existsSync(constants) && statSync(constants).isFile()) {
      return root;
    }
    root = dirname(root);
  }

  throw new SettingsRefused(`school Data Directories not Found above ${resolve(start)}`);
}

/** findRegistryFile Resolves a relative Path from the Data Root, and Refuses a File that does not Exist. */
export function findRegistryFile(root: string, path: string): string {
  const resolved = isAbsolute(path) ? path : join(root, path);
  if (!existsSync(resolved) || !statSync(resolved).isFile()) {
    throw new SettingsRefused(`registryPath ${resolved} Must be a File that Exists`);
  }

  return resolved;
}

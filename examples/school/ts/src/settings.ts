import { readFileSync } from "node:fs";

// ValueRule Validates one declared Configuration Key.
export type ValueRule = (value: unknown) => boolean;
type DataValues = Record<string, unknown>;
export type Environment = Readonly<Record<string, string | undefined>>;

// readDataFile Rejects Unknown Keys and Invalid Values before Merging.
export function readDataFile(path: URL, rules: Record<string, ValueRule>): DataValues {
  const values: unknown = JSON.parse(readFileSync(path, "utf8"));
  if (!values || typeof values !== "object" || Array.isArray(values)) {
    throw new Error(`${path.pathname} Must Hold an Object`);
  }

  for (const [key, value] of Object.entries(values)) {
    if (!Object.hasOwn(rules, key)) throw new Error(`${path.pathname}: Unknown Key ${key}`);
    if (!rules[key]!(value)) throw new Error(`${path.pathname}: Invalid ${key}`);
  }

  return values as DataValues;
}

// checkIntegerRange Keeps JSON Numbers whole and Bounded.
export function checkIntegerRange(minimum: number, maximum: number): ValueRule {
  return (value) => {
    if (typeof value !== "number") return false;

    const integer = Number.isSafeInteger(value);
    const bounded = value >= minimum && value <= maximum;

    return integer && bounded;
  };
}

// checkTextValue Requires a nonempty String.
export function checkTextValue(value: unknown): boolean {
  return typeof value === "string" && value.trim().length > 0;
}

// decodeDataValues Requires every declared Key before Returning frozen Data.
export function decodeDataValues<T>(values: DataValues, rules: Record<string, ValueRule>): Readonly<T> {
  for (const key of Object.keys(rules)) {
    if (!Object.hasOwn(values, key)) throw new Error(`${key} is Missing`);
  }

  return Object.freeze({ ...values }) as Readonly<T>;
}

// readGlobalValues Loads Constants without an Environment Override Path.
export function readGlobalValues<T>(path: URL, rules: Record<string, ValueRule>): Readonly<T> {
  const values = readDataFile(path, rules);

  return decodeDataValues<T>(values, rules);
}

// readEnvironmentFiles Applies Defaults before the selected Environment.
export function readEnvironmentFiles(root: URL, concern: string, environment: string | undefined,
  rules: Record<string, ValueRule>): DataValues {
  if (environment !== "development" && environment !== "production") {
    throw new Error("APP_ENV Must be development or production");
  }

  const defaults = readDataFile(new URL(`config/${concern}.defaults.json`, root), rules);
  const overrides = readDataFile(new URL(`config/${concern}.${environment}.json`, root), rules);

  return { ...defaults, ...overrides };
}

// checkEnvironmentKeys Refuses undeclared Variables in this Service Namespace.
export function checkEnvironmentKeys(environment: Environment, prefix: string,
  names: Record<string, string>): void {
  for (const name of Object.keys(environment)) {
    if (name.startsWith(prefix) && !Object.hasOwn(names, name)) {
      throw new Error(`Unknown Variable ${name}; Global Constants Cannot be Overridden`);
    }
  }
}

// applyEnvironmentValues Validates declared Variables before Replacing Values.
export function applyEnvironmentValues(values: DataValues, environment: Environment,
  names: Record<string, string>, rules: Record<string, ValueRule>): void {
  for (const [name, key] of Object.entries(names)) {
    const value = environment[name];
    if (value === undefined) continue;
    values[key] = encodeEnvironmentValue(value, rules[key]!, name);
  }
}

// encodeEnvironmentValue Accepts Text or a decimal Integer, never implicit Booleans.
function encodeEnvironmentValue(value: string, rule: ValueRule, name: string): unknown {
  if (rule(value)) return value;

  const number = /^[0-9]+$/.test(value) ? Number(value) : NaN;
  if (!rule(number)) throw new Error(`${name}: Invalid Environment Value`);

  return number;
}

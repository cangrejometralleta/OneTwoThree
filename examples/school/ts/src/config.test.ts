import assert from "node:assert/strict";
import { copyFileSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { test, type TestContext } from "node:test";

import { loadSchoolConfig } from "./config.js";
import { SchoolValues } from "./constants.js";
import { checkIntegerRange, readGlobalValues, type Environment } from "./settings.js";

// buildConfigFixture Keeps broken Files away from the checked-in Examples.
function buildConfigFixture(t: TestContext): URL {
  const root = pathToFileURL(mkdtempSync(join(tmpdir(), "school-config-")) + "/");
  mkdirSync(new URL("config/", root));
  t.after(() => rmSync(root, { recursive: true, force: true }));

  for (const layer of ["defaults", "development", "production"]) {
    const name = `config/school.${layer}.json`;
    copyFileSync(new URL(`../../${name}`, import.meta.url), new URL(name, root));
  }

  return root;
}

test("Configuration Applies Defaults, File and Variable in Order", (t) => {
  const root = buildConfigFixture(t);
  writeFileSync(new URL("config/school.production.json", root), '{"databasePath":"production.db","port":9000}');
  const environment = { APP_ENV: "production", TOKEN_SECRET: "test-secret" };

  const file = loadSchoolConfig(environment, root);
  assert.equal(file.port, 9000);
  assert.equal(file.databasePath, "production.db");
  const override = loadSchoolConfig({ ...environment, SCHOOL_PORT: "9100" }, root);

  assert.equal(override.port, 9100);
  assert.equal(SchoolValues.minimumAgeYears, 18);
  assert.ok(Object.isFrozen(override));
  assert.ok(Object.isFrozen(SchoolValues));
});

test("Configuration Rejects Invalid Startup Inputs", async (t) => {
  const cases: { name: string; environment?: Environment; file?: string; missing?: boolean }[] = [
    { name: "missing environment", environment: { APP_ENV: undefined } },
    { name: "unknown environment", environment: { APP_ENV: "staging" } },
    { name: "path traversal", environment: { APP_ENV: "../production" } },
    { name: "empty override", environment: { SCHOOL_PORT: "" } },
    { name: "invalid port", environment: { SCHOOL_PORT: "65536" } },
    { name: "fractional port", environment: { SCHOOL_PORT: "3.5" } },
    { name: "empty path", environment: { SCHOOL_DATABASE_PATH: " " } },
    { name: "global variable", environment: { SCHOOL_MINIMUM_AGE_YEARS: "1" } },
    { name: "unknown variable", environment: { SCHOOL_TYPO: "1" } },
    { name: "global file key", file: '{"databasePath":"test.db","minimumAgeYears":1}' },
    { name: "unknown file key", file: '{"databasePath":"test.db","typo":1}' },
    { name: "null key", file: '{"databasePath":null}' },
    { name: "invalid type", file: '{"databasePath":"test.db","port":"9000"}' },
    { name: "invalid file range", file: '{"databasePath":"test.db","port":0}' },
    { name: "missing key", file: '{}' },
    { name: "invalid json", file: '{' },
    { name: "null document", file: 'null' },
    { name: "array document", file: '[]' },
    { name: "missing file", missing: true },
    { name: "missing secret", environment: { TOKEN_SECRET: "" } },
    { name: "unknown adapter", environment: { SERVER: "missing" } },
    { name: "invalid lifetime", environment: { SCHOOL_TOKEN_LIFE_SECONDS: "0" } },
    { name: "secret in file", file: '{"databasePath":"test.db","tokenSecret":"forbidden"}' },
  ];

  for (const item of cases) {
    await t.test(item.name, (t) => {
      const root = buildConfigFixture(t);
      const path = new URL("config/school.production.json", root);
      if (item.file !== undefined) writeFileSync(path, item.file);
      if (item.missing) rmSync(path);

      const environment = { APP_ENV: "production", TOKEN_SECRET: "test-secret", ...item.environment };
      assert.throws(() => loadSchoolConfig(environment, root));
    });
  }
});

test("An Override Cannot Hide Invalid Defaults", (t) => {
  const root = buildConfigFixture(t);
  writeFileSync(new URL("config/school.defaults.json", root), '{"port":null}');

  assert.throws(() => loadSchoolConfig({
    APP_ENV: "production", SCHOOL_PORT: "9100", TOKEN_SECRET: "test-secret",
  }, root));
});

test("Global Constants Reject Invalid Data at Startup", (t) => {
  const root = buildConfigFixture(t);
  const path = new URL("constants.json", root);
  writeFileSync(path, '{"minimumAgeYears":null}');

  assert.throws(() => readGlobalValues(path, { minimumAgeYears: checkIntegerRange(1, 200) }));
});

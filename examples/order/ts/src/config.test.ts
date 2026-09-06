import assert from "node:assert/strict";
import { copyFileSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { test, type TestContext } from "node:test";

import { loadOrderConfig } from "./config.js";
import { OrderValues } from "./constants.js";
import { checkIntegerRange, readGlobalValues, type Environment } from "./settings.js";

// buildConfigFixture Keeps broken Files away from the checked-in Examples.
function buildConfigFixture(t: TestContext): URL {
  const root = pathToFileURL(mkdtempSync(join(tmpdir(), "order-config-")) + "/");
  mkdirSync(new URL("config/", root));
  t.after(() => rmSync(root, { recursive: true, force: true }));

  for (const layer of ["defaults", "development", "production"]) {
    const name = `config/order.${layer}.json`;
    copyFileSync(new URL(`../../${name}`, import.meta.url), new URL(name, root));
  }

  return root;
}

test("Configuration Applies Defaults, File and Variable in Order", (t) => {
  const root = buildConfigFixture(t);
  writeFileSync(new URL("config/order.production.json", root), '{"databasePath":"production.db","port":9000}');
  const environment = { APP_ENV: "production" };

  const file = loadOrderConfig(environment, root);
  assert.equal(file.port, 9000);
  assert.equal(file.databasePath, "production.db");
  const override = loadOrderConfig({ ...environment, ORDER_PORT: "9100" }, root);

  assert.equal(override.port, 9100);
  assert.equal(OrderValues.itemNameWidth, 12);
  assert.ok(Object.isFrozen(override));
  assert.ok(Object.isFrozen(OrderValues));
});

test("Configuration Rejects Invalid Startup Inputs", async (t) => {
  const cases: { name: string; environment?: Environment; file?: string; missing?: boolean }[] = [
    { name: "missing environment", environment: { APP_ENV: undefined } },
    { name: "unknown environment", environment: { APP_ENV: "staging" } },
    { name: "path traversal", environment: { APP_ENV: "../production" } },
    { name: "empty override", environment: { ORDER_PORT: "" } },
    { name: "invalid port", environment: { ORDER_PORT: "65536" } },
    { name: "fractional port", environment: { ORDER_PORT: "3.5" } },
    { name: "empty path", environment: { ORDER_DATABASE_PATH: " " } },
    { name: "global variable", environment: { ORDER_ITEM_NAME_WIDTH: "1" } },
    { name: "unknown variable", environment: { ORDER_TYPO: "1" } },
    { name: "global file key", file: '{"databasePath":"test.db","itemNameWidth":1}' },
    { name: "unknown file key", file: '{"databasePath":"test.db","typo":1}' },
    { name: "null key", file: '{"databasePath":null}' },
    { name: "invalid type", file: '{"databasePath":"test.db","port":"9000"}' },
    { name: "invalid file range", file: '{"databasePath":"test.db","port":0}' },
    { name: "missing key", file: '{}' },
    { name: "invalid json", file: '{' },
    { name: "null document", file: 'null' },
    { name: "array document", file: '[]' },
    { name: "missing file", missing: true },
  ];

  for (const item of cases) {
    await t.test(item.name, (t) => {
      const root = buildConfigFixture(t);
      const path = new URL("config/order.production.json", root);
      if (item.file !== undefined) writeFileSync(path, item.file);
      if (item.missing) rmSync(path);

      const environment = { APP_ENV: "production", ...item.environment };
      assert.throws(() => loadOrderConfig(environment, root));
    });
  }
});

test("An Override Cannot Hide Invalid Defaults", (t) => {
  const root = buildConfigFixture(t);
  writeFileSync(new URL("config/order.defaults.json", root), '{"port":null}');

  assert.throws(() => loadOrderConfig({
    APP_ENV: "production", ORDER_PORT: "9100",
  }, root));
});

test("Global Constants Reject Invalid Data at Startup", (t) => {
  const root = buildConfigFixture(t);
  const path = new URL("constants.json", root);
  writeFileSync(path, '{"itemNameWidth":null}');

  assert.throws(() => readGlobalValues(path, { itemNameWidth: checkIntegerRange(1, 200) }));
});

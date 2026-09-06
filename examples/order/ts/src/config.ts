import { applyEnvironmentValues, checkEnvironmentKeys, checkIntegerRange, checkTextValue,
  decodeDataValues, readEnvironmentFiles, type Environment } from "./settings.js";

// OrderConfig Holds Deployment Values; Secrets Have no File Key.
export type OrderConfig = Readonly<{
  port: number;
  databasePath: string;
}>;

// loadOrderConfig Resolves Startup Inputs before any Store or Listener Opens.
export function loadOrderConfig(environment: Environment,
  root = new URL("../../", import.meta.url)): OrderConfig {
  const rules = {
    port: checkIntegerRange(1, 65535),
    databasePath: checkTextValue,
  };
  const names = {
    ORDER_PORT: "port",
    ORDER_DATABASE_PATH: "databasePath",
  };
  checkEnvironmentKeys(environment, "ORDER_", names);

  const values = readEnvironmentFiles(root, "order", environment["APP_ENV"], rules);
  applyEnvironmentValues(values, environment, names, rules);
  const config = decodeDataValues<OrderConfig>(values, rules);

  return config;
}

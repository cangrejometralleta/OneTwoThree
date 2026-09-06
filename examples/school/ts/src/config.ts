import { applyEnvironmentValues, checkEnvironmentKeys, checkIntegerRange, checkTextValue,
  decodeDataValues, readEnvironmentFiles, type Environment } from "./settings.js";

// SchoolConfig Holds Deployment Values; Secrets Have no File Key.
export type SchoolConfig = Readonly<{
  port: number;
  databasePath: string;
  tokenLifeSeconds: number;
  serverAdapter: string;
  tokenSecret: string;
}>;

// loadSchoolConfig Resolves Startup Inputs before any Store or Listener Opens.
export function loadSchoolConfig(environment: Environment,
  root = new URL("../../", import.meta.url)): SchoolConfig {
  const rules = {
    port: checkIntegerRange(1, 65535),
    databasePath: checkTextValue,
    tokenLifeSeconds: checkIntegerRange(1, 86400),
    serverAdapter: checkServerAdapter,
  };
  const names = {
    SCHOOL_PORT: "port",
    SCHOOL_DATABASE_PATH: "databasePath",
    SCHOOL_TOKEN_LIFE_SECONDS: "tokenLifeSeconds",
    SERVER: "serverAdapter",
  };
  checkEnvironmentKeys(environment, "SCHOOL_", names);

  const values = readEnvironmentFiles(root, "school", environment["APP_ENV"], rules);
  applyEnvironmentValues(values, environment, names, rules);
  const config = decodeDataValues<Omit<SchoolConfig, "tokenSecret">>(values, rules);

  const secret = environment["TOKEN_SECRET"];
  if (!checkTextValue(secret)) throw new Error("TOKEN_SECRET is Missing");
  return Object.freeze({ ...config, tokenSecret: secret as string });
}

// checkServerAdapter Accepts only Adapters this Runtime Implements.
function checkServerAdapter(value: unknown): boolean {
  return typeof value === "string" && ["stdlib", "node", "express", "fastify"].includes(value);
}

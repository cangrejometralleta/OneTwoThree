import { loadSchoolConfig, type SchoolConfig } from "./config.js";
import { selectServerAdapter } from "./adapters.js";
import { SchoolAPI } from "./handlers.js";
import { SqliteSchool } from "./store-sqlite.js";
import { HmacTokens } from "./token-hmac.js";

// main Casts the Players, then Steps off the Stage.
async function main(): Promise<void> {
  const config = loadSchoolConfig({ ...process.env });

  const school = SqliteSchool.shapeSchoolTables(config.databasePath);
  const api = new SchoolAPI(school, school, buildTokenIssuer(config));
  const server = selectServerAdapter(config.serverAdapter);

  await server.serveRoutes(api.declareSchoolRoutes(), config.port);
  console.log(`✅ School Listening on :${config.port} through "${config.serverAdapter}"`);
}

// buildTokenIssuer Uses the Secret and Lifetime Validated at Startup.
function buildTokenIssuer(config: SchoolConfig): HmacTokens {
  return new HmacTokens(config.tokenSecret, config.tokenLifeSeconds * 1000);
}

main().catch((failure) => {
  console.error("❌ School Refused to Start:", failure);
  process.exit(1);
});

import { selectServerAdapter } from "./adapters.js";
import { SchoolAPI } from "./handlers.js";
import { SqliteSchool } from "./store-sqlite.js";
import { HmacTokens } from "./token-hmac.js";

// main Casts the Players, then Steps off the Stage.
async function main(): Promise<void> {
  const school = SqliteSchool.shapeSchoolTables("school.db");
  const api = new SchoolAPI(school, school, buildTokenIssuer());

  const choice = process.env["SERVER"] ?? "node";
  const server = selectServerAdapter(choice);

  await server.serveRoutes(api.declareSchoolRoutes(), 8080);
  console.log(`✅ School Listening on :8080 through "${choice}"`);
}

// buildTokenIssuer Refuses to Start without a Secret.
// A silent Default Secret is a public Key with extra Steps.
function buildTokenIssuer(): HmacTokens {
  const secret = process.env["TOKEN_SECRET"];

  if (!secret) {
    console.error("❌ TOKEN_SECRET is Missing");
    process.exit(1);
  }

  return new HmacTokens(secret, 10 * 60 * 1000);
}

main().catch((failure) => {
  console.error("❌ School Refused to Start:", failure);
  process.exit(1);
});

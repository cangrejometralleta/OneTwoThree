import { SchoolService } from './src/app/service.ts';
import { Office } from './src/campus/office.ts';
import { SchoolHandler } from './src/handler/handler.ts';
import { ExpressServer } from './src/serving/express/express-server.ts';
import { StdlibServer } from './src/serving/stdlib.ts';
import { loadSchoolConfig } from './src/settings/school-config.ts';
import { findSchoolDataRoot, readGlobalConstants } from './src/settings/settings.ts';
import { SchoolStore } from './src/store/school-store.ts';
import { AccessTokens } from './src/tokens/access-tokens.ts';
import type { Server } from './src/transport/transport.ts';

/** SERVER_ADAPTERS is every Framework this Runtime Implements; the Loader only Receives the Names. */
const SERVER_ADAPTERS: Readonly<Record<string, () => Server>> = {
  stdlib: () => new StdlibServer(),
  express: () => new ExpressServer(),
};

/** startSchool Casts the Players, then Steps off the Stage. It is the only File that Knows every Module by Name. */
async function startSchool(): Promise<void> {
  const root = findSchoolDataRoot(process.cwd());
  const constants = readGlobalConstants(root);
  const config = loadSchoolConfig(root, process.env, Object.keys(SERVER_ADAPTERS).sort());

  const store = SchoolStore.openSchoolStore(config.databasePath);
  const office = await Office.openOffice(config.registryPath, (line) => console.log(line));
  const service = new SchoolService({
    rules: { minimumAgeYears: constants.minimumAgeYears },
    students: store, courses: store, registry: office, notices: office,
    log: (line) => console.error(line),
  });
  const door = new SchoolHandler(service, new AccessTokens(config.tokenSecret, config.tokenLifeSeconds));

  const server = SERVER_ADAPTERS[config.serverAdapter]?.();
  const running = await server?.serveRoutes(door.declareSchoolRoutes(), config.port);
  console.log(`✅ School Listening on :${running?.port} through ${config.serverAdapter}`);
}

try {
  await startSchool();
} catch (refused) {
  console.error(`❌ School Refused to Start: ${refused instanceof Error ? refused.message : String(refused)}`);
  process.exitCode = 1;
}

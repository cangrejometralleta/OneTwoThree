import { loadSettings } from './src/settings/loader.ts';
import { FileStore } from './src/store/store.ts';
import { Office } from './src/campus/office.ts';
import { HmacTokenIssuer } from './src/tokens/token.ts';
import { SchoolService } from './src/app/service.ts';
import { Handler } from './src/handler/handler.ts';
import { serve } from './src/serving/stdlib.ts';

try {
  const settings = await loadSettings(import.meta.url);
  const [store, office] = await Promise.all([FileStore.open(settings.databasePath), Office.open(settings.registryPath)]);
  const tokens = new HmacTokenIssuer(settings.tokenSecret,settings.tokenLifeSeconds);
  const service = new SchoolService(store,store,office,office,settings.constants);
  await serve(new Handler(service,tokens),settings.port);
  console.log(`School listening on ${settings.port}`);
} catch (error) {
  const message = error instanceof Error ? error.message : 'unknown startup failure';
  console.error(`startup refused: ${message}`);
  process.exitCode = 1;
}

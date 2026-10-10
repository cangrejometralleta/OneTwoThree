import { readFile } from 'node:fs/promises';
import { readFileSync } from 'node:fs';
import { dirname, resolve, isAbsolute, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parseJson } from '../codec/json.ts';
import type { MinimumAge } from '../school/index.ts';
export type Settings = { appEnv: 'development'|'production'; port: number; databasePath: string; registryPath: string; tokenLifeSeconds: bigint; serverAdapter: 'stdlib'; tokenSecret: string; dataRoot: string; constants: MinimumAge };
const allowed = new Set(['port','databasePath','registryPath','tokenLifeSeconds','serverAdapter']);
function needString(v: unknown, key: string): string { if (typeof v !== 'string' || !v) throw new Error(`${key} must be a nonempty string`); return v; }
function needInteger(v: unknown, key: string, min: bigint, max: bigint): bigint {
  if (typeof v !== 'bigint' || v < min || v > max) throw new Error(`${key} must be an integer from ${min} through ${max}`); return v;
}
async function readObject(path: string): Promise<Record<string,unknown>> {
  let parsed: unknown;
  try { parsed = parseJson(await readFile(path,'utf8')); } catch { throw new Error(`cannot read valid JSON from ${path}`); }
  if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) throw new Error(`${path} must contain an object`);
  return parsed as Record<string,unknown>;
}
function envFile(path: string): Record<string,string> {
  let content: string; try { content = readFileSync(path,'utf8'); } catch { return {}; }
  const values: Record<string,string> = {};
  for (const line of content.split(/\r?\n/)) { const m = /^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*?)\s*$/.exec(line); if (m) values[m[1]!] = m[2]!.replace(/^(['"])(.*)\1$/,'$2'); }
  return values;
}
export async function loadSettings(programUrl: string, actualEnv = process.env): Promise<Settings> {
  const programDir = dirname(fileURLToPath(programUrl));
  const fileEnv = envFile(join(programDir,'.env'));
  const env = { ...fileEnv, ...actualEnv };
  const appEnv = env.APP_ENV; if (appEnv !== 'development' && appEnv !== 'production') throw new Error('APP_ENV must be development or production');
  const dataRoot = resolve(env.DATA_ROOT || './data');
  const defaults = await readObject(join(dataRoot,'config','school.defaults.json'));
  const environment = await readObject(join(dataRoot,'config',`school.${appEnv}.json`));
  for (const [source, values] of [['defaults',defaults],['environment',environment]] as const) {
    for (const key of Object.keys(values)) {
      if (!allowed.has(key)) throw new Error(`unknown ${source} configuration key ${key}`);
      if (values[key] === null) throw new Error(`${source} configuration ${key} cannot be null`);
    }
  }
  const config = { ...defaults, ...environment };
  const varMap: Record<string,string> = { SCHOOL_PORT:'port', SCHOOL_DATABASE_PATH:'databasePath', SCHOOL_REGISTRY_PATH:'registryPath', SCHOOL_TOKEN_LIFE_SECONDS:'tokenLifeSeconds', SERVER:'serverAdapter' };
  for (const name of Object.keys(env)) if (name.startsWith('SCHOOL_') && !(name in varMap)) throw new Error(`unknown environment variable ${name}`);
  for (const [name,key] of Object.entries(varMap)) if (env[name] !== undefined) {
    if (key === 'port' || key === 'tokenLifeSeconds') { if (!/^\d+$/.test(env[name])) throw new Error(`${name} must contain digits only`); config[key] = BigInt(env[name]); }
    else config[key] = env[name];
  }
  const port = Number(needInteger(config.port,'port',1n,65535n));
  const databasePath = needString(config.databasePath,'databasePath');
  let registryPath = needString(config.registryPath,'registryPath'); if (!isAbsolute(registryPath)) registryPath = resolve(dataRoot,registryPath);
  const tokenLifeSeconds = needInteger(config.tokenLifeSeconds,'tokenLifeSeconds',1n,86400n);
  if (config.serverAdapter !== 'stdlib') throw new Error('SERVER must name an implemented adapter: stdlib');
  const tokenSecret = needString(env.TOKEN_SECRET,'TOKEN_SECRET');
  const constantsRaw = await readObject(join(dataRoot,'constants','school.json'));
  const constants: MinimumAge = { minimumAgeYears: needInteger(constantsRaw.minimumAgeYears,'minimumAgeYears',1n,I64_MAX) };
  return { appEnv, port, databasePath: isAbsolute(databasePath) ? databasePath : resolve(databasePath), registryPath, tokenLifeSeconds, serverAdapter:'stdlib', tokenSecret, dataRoot, constants };
}
const I64_MAX = (1n << 63n)-1n;

import {
  applyEnvironmentValues, checkEnvironmentKeys, checkIntegerRange, checkOneOf, checkTextValue, findRegistryFile,
  readEnvironmentFiles, requireEveryKey, SettingsRefused,
} from './settings.ts';

/** The Bounds a Deployment Value Must Fall inside. Read alone, 65535 Says nothing. */
export const LOWEST_PORT = 1n;
export const HIGHEST_PORT = 65535n;
export const SECONDS_IN_A_DAY = 60n * 60n * 24n;

const NAMESPACE = 'SCHOOL_';

/** SchoolConfig Holds the Deployment Values, Validated once at Startup. The Secret has no File Key. */
export interface SchoolConfig {
  readonly port: number;
  readonly databasePath: string;
  readonly registryPath: string;
  readonly tokenLifeSeconds: bigint;
  readonly serverAdapter: string;
  readonly tokenSecret: string;
}

/**
 * loadSchoolConfig Resolves Defaults, then the Environment File, then the Variables.
 * The Adapters a Runtime Implements Arrive from the Composition; this Loader Names none.
 */
export function loadSchoolConfig(root: string, environment: Readonly<Record<string, string | undefined>>,
  adapters: readonly string[]): SchoolConfig {
  const rules = {
    port: checkIntegerRange(LOWEST_PORT, HIGHEST_PORT),
    databasePath: checkTextValue,
    registryPath: checkTextValue,
    tokenLifeSeconds: checkIntegerRange(1n, SECONDS_IN_A_DAY),
    serverAdapter: checkOneOf(adapters),
  };
  const names = {
    SCHOOL_PORT: 'port',
    SCHOOL_DATABASE_PATH: 'databasePath',
    SCHOOL_REGISTRY_PATH: 'registryPath',
    SCHOOL_TOKEN_LIFE_SECONDS: 'tokenLifeSeconds',
    SERVER: 'serverAdapter',
  };

  checkEnvironmentKeys(environment, NAMESPACE, names);
  const values = readEnvironmentFiles(root, 'school', environment['APP_ENV'], rules);
  applyEnvironmentValues(values, environment, names, rules);
  requireEveryKey(values, rules);

  const tokenSecret = environment['TOKEN_SECRET'] ?? '';
  if (tokenSecret.trim() === '') {
    throw new SettingsRefused('TOKEN_SECRET is Missing');
  }

  return {
    port: Number(values.get('port')),
    databasePath: values.get('databasePath') as string,
    registryPath: findRegistryFile(root, values.get('registryPath') as string),
    tokenLifeSeconds: values.get('tokenLifeSeconds') as bigint,
    serverAdapter: values.get('serverAdapter') as string,
    tokenSecret,
  };
}

# School service (TypeScript)

Requires Node.js 22.6 or newer. This implementation uses Node's built-in TypeScript
type stripping and standard library for runtime. Run `npm install` once to get the TypeScript compiler used by `build.sh`.

```sh
npm install
./build.sh
TOKEN_SECRET=spec-secret ./run.sh
```

The process reads `school.defaults.json`, then `school.<APP_ENV>.json`, then
environment variables. Both configuration files, `constants/school.json`,
and the registry file must exist under the data root (`DATA_ROOT`, set by `run.sh` to the parent School example). The repository supplies them in `config/`, `constants/` and `config/registry.json`.
The TypeScript runner uses `school-ts.db` by default to keep its JSON store separate from the SQLite services. The database path is relative to the working directory, and the server listens on `0.0.0.0`.

The service's tests use fakes at the port boundaries and Node's built-in test runner.
Run them with `npm test` or `./build.sh`.

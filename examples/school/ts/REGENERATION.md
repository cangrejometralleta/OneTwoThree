# Gaps and implementation choices

- The clean input set intentionally contained only the specification, shape and contract. The integrated runner uses the shared `examples/school/config/`, `constants/` and registry files.
- The spec does not define the `SERVER` adapter identifier; this runtime accepts
  only `stdlib` (Node's built-in HTTP server).
- The spec leaves the persistence format and database suffix open. This runtime
  writes JSON to the exact configured database path and uses an in-process
  serialized write queue; it is suitable for a single process.
- The spec leaves enrollment logging format to the Runtime. This runtime logs one
  line to stderr: `student <id> is enrolled`.
- `DATA_ROOT` is an additional runtime setting used to locate constants, config,
  and the registry. The runner sets it to the parent School example directory; it is outside the declared
  `SCHOOL_` namespace.

## Verification limit

- The environment denied binding a local test socket (`EPERM`), so the stdlib
  adapter was imported but not exercised over HTTP. Handler and service behavior
  are covered by the passing unit tests.
- Node does not ship a TypeScript formatter. `build.sh` checks Node's TypeScript
  syntax handling, runs `tsc --noEmit`, executes the behavior tests, and checks import boundaries. The type check required fixing several inferred types after integration.

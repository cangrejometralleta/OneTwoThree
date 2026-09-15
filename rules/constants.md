# Constants

- Global constants Hold the same meaning across environments.  
  Keep shared declarative Values in descriptive files under `constants/`.
- A value that changes by environment is Configuration.  
  Keep it separately under `config/`, in a descriptive YAML, JSON or TOML File.
- Name the File after the concern it holds.  
  `enrollment.production.yaml` Says more than `constants.yaml`.
- Keep algorithmic and protocol Constants beside their logic.  
  A deployment must not Redefine what the algorithm means.

## Global Constants Define shared Meaning

- Separate global Constants from environment defaults.  
  A default can be Overridden; a global constant cannot.
- Load global Constants separately from environment configuration.  
  Reject environment Files or variables that attempt to override them.
- Validate both Sets at startup and expose them as separate typed values.  
  Neither set Changes during execution.

Shared enrollment States in `constants/enrollment.yaml`:

```yaml
enrollmentStates:
  - pending
  - accepted
  - rejected
```

The same global Constants can instead live  
in `constants/enrollment.json`:

```json
{
  "enrollmentStates": ["pending", "accepted", "rejected"]
}
```

Or in `constants/enrollment.toml`:

```toml
enrollmentStates = ["pending", "accepted", "rejected"]
```

Development and production Read the same file.  
Changing these states Changes the application contract,  
so the change Travels with the code that uses them.

## Files Describe the Values

- Choose one Format per configuration set.  
  YAML, JSON and TOML below are Alternatives.
- Give each Key a meaning and each quantity a unit.  
  `requestTimeoutSeconds`, never just `timeout`.
- Keep Secrets outside versioned files.  
  Inject them through environment Variables or a secret store.

Environment defaults in `config/enrollment.defaults.yaml`.  
These values are overridable Configuration, never global constants:

```yaml
requestTimeoutSeconds: 5
retryLimit: 2
studentLimit: 30
```

Local overrides in `config/enrollment.development.yaml`:

```yaml
serviceUrl: http://localhost:8080
retryLimit: 0
```

Production overrides in `config/enrollment.production.yaml`:

```yaml
serviceUrl: https://enrollment.example.com
requestTimeoutSeconds: 15
```

The production file can instead be JSON,  
at `config/enrollment.production.json`:

```json
{
  "serviceUrl": "https://enrollment.example.com",
  "requestTimeoutSeconds": 15
}
```

For a TOML configuration set, defaults Live  
in `config/enrollment.defaults.toml`:

```toml
requestTimeoutSeconds = 5
retryLimit = 2
studentLimit = 30
```

Local overrides in `config/enrollment.development.toml`:

```toml
serviceUrl = "http://localhost:8080"
retryLimit = 0
```

Production overrides in `config/enrollment.production.toml`:

```toml
serviceUrl = "https://enrollment.example.com"
requestTimeoutSeconds = 15
```

## Environments Select the Values

- Select the Environment explicitly at startup.  
  Reject an unknown Environment or a missing selected file.
- Apply the Defaults, then the environment file, then declared variables.  
  Later values Replace earlier values for the same key.  
  Global constants never Enter this override chain.
- Load and validate once at Startup; pass typed configuration onward.  
  Reject unknown Keys, invalid types and missing required values.

For this example, files hold flat Keys.  
An absent key Inherits its default; `null` is invalid.

The loader Declares these environment variables:

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | File Selection | `development` or `production`; Required |
| `ENROLLMENT_REQUEST_TIMEOUT_SECONDS` | `requestTimeoutSeconds` | Positive Integer |
| `ENROLLMENT_API_TOKEN` | `apiToken` | Nonempty Secret; Required in Production |

Validate `serviceUrl` as an HTTP or HTTPS URL,  
`retryLimit` as a nonnegative Integer,  
and `studentLimit` as a positive integer.

With `APP_ENV=production` and  
`ENROLLMENT_REQUEST_TIMEOUT_SECONDS=20`,  
the resolved nonsecret Values are:

```yaml
requestTimeoutSeconds: 20
retryLimit: 2
studentLimit: 30
serviceUrl: https://enrollment.example.com
```

The loader Implements selection, overrides and validation.  
YAML, JSON and TOML only describe Data;  
they do not Expand environment variables on their own.

## The Declared Surface

> Proposed. Drawn from one Project, cangrejometralleta/muchi-api,
> and not yet Weathered by a second.

- A variable nobody declared is a variable nobody Finds.  
  `.env.example` is the Contract, and it is committed.
- Required variables Live uncommented, with a placeholder value.  
  A reader copies the file and Sees at once what startup wants.
- Optional overrides Live commented, showing the default they replace.  
  The comment Documents; it never loads.
- `.env` is Local, and git ignores it.  
  The loader Reads it when present and says nothing when absent.
- An exported variable Wins. The file fills the gaps it left.  
  A shell that sets a value Meant it; a file only suggested one.

```sh
# Required. No default Exists, and startup stops without them.
APP_ENV=development
TOKEN_SECRET=replace-with-a-random-32-byte-secret

# Optional overrides; the Defaults shown below apply when unset.
# SCHOOL_PORT=8080
# SERVER=stdlib
```

Three sources, and the order never Changes:
the argument, then the environment, then the default.
Each one is more Explicit than the one behind it.

## Six Kinds, in three Pairs

> Proposed. Four are Canon above. Domain data and deployment topology
> come from one Project and await a second.

A value Belongs to exactly one kind. Asking which one Answers
where it lives, who may change it, and what breaks if it moves.

**What the Code Means.** It Changes when the code changes.

| Kind | Lives in |
| --- | --- |
| Global Constant | `constants/`, a descriptive File |
| Algorithmic Constant | beside its Logic |

**What the Deployment Chooses.** It Changes when the deployment changes.

| Kind | Lives in |
| --- | --- |
| Environment Configuration | `config/school.production.json` |
| Deployment Topology | `config/deploy.env` |

**What Comes from outside.** It Changes when the world changes.

| Kind | Lives in |
| --- | --- |
| Domain Data | `config/`, its own File |
| Secret | Injection only, never a File |

- **Domain Data Changes without a code change and without an environment.**  
  A per-store scraping Rule, a tax table, a holiday calendar.  
  It Reads like configuration and deploys like content.
- **Deployment Topology never Reaches the running program.**  
  Region, memory, instance count, service account.  
  The deploy script Reads it; the application never opens the file.
- Six kinds do not fit in a Head. Three pairs do.  
  Ask which Pair first, then which of the two.

## Constants Stay with their Meaning

A unit conversion Belongs beside the calculation:

```go
const secondsPerMinute = 60
```

A shared state catalog Belongs in `constants/`.  
A timeout Belongs in `config/`.  
A conversion factor Belongs beside its calculation.

Secrets Enter through injection, never through either directory.

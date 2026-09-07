# Constants

- Global Constants Hold the same Meaning across Environments.  
  Keep shared declarative Values in descriptive Files under `constants/`.
- A Value that Changes by Environment is Configuration.  
  Keep it separately under `config/`, in a descriptive YAML, JSON or TOML File.
- Name the File after the Concern it Holds.  
  `enrollment.production.yaml` Says more than `constants.yaml`.
- Keep Algorithmic and Protocol Constants beside their Logic.  
  A Deployment Must not Redefine what the Algorithm Means.

## Global Constants Define shared Meaning

- Separate Global Constants from Environment Defaults.  
  A Default Can be Overridden; a Global Constant Cannot.
- Load Global Constants separately from Environment Configuration.  
  Reject Environment Files or Variables that Attempt to Override them.
- Validate both Sets at Startup and Expose them as separate typed Values.  
  Neither Set Changes during Execution.

Shared Enrollment States in `constants/enrollment.yaml`:

```yaml
enrollmentStates:
  - pending
  - accepted
  - rejected
```

The same Global Constants can instead Live  
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

Development and Production Read the same File.  
Changing these States Changes the Application Contract,  
so the Change Travels with the Code that Uses them.

## Files Describe the Values

- Choose one Format per Configuration Set.  
  YAML, JSON and TOML below are Alternatives.
- Give each Key a Meaning and each Quantity a Unit.  
  `requestTimeoutSeconds`, never just `timeout`.
- Keep Secrets outside versioned Files.  
  Inject them through Environment Variables or a Secret Store.

Environment Defaults in `config/enrollment.defaults.yaml`.  
These Values are Overridable Configuration, never Global Constants:

```yaml
requestTimeoutSeconds: 5
retryLimit: 2
studentLimit: 30
```

Local Overrides in `config/enrollment.development.yaml`:

```yaml
serviceUrl: http://localhost:8080
retryLimit: 0
```

Production Overrides in `config/enrollment.production.yaml`:

```yaml
serviceUrl: https://enrollment.example.com
requestTimeoutSeconds: 15
```

The Production File can instead be JSON,  
at `config/enrollment.production.json`:

```json
{
  "serviceUrl": "https://enrollment.example.com",
  "requestTimeoutSeconds": 15
}
```

For a TOML Configuration Set, Defaults Live  
in `config/enrollment.defaults.toml`:

```toml
requestTimeoutSeconds = 5
retryLimit = 2
studentLimit = 30
```

Local Overrides in `config/enrollment.development.toml`:

```toml
serviceUrl = "http://localhost:8080"
retryLimit = 0
```

Production Overrides in `config/enrollment.production.toml`:

```toml
serviceUrl = "https://enrollment.example.com"
requestTimeoutSeconds = 15
```

## Environments Select the Values

- Select the Environment Explicitly at Startup.  
  Reject an Unknown Environment or a Missing selected File.
- Apply Defaults, then the Environment File, then declared Variables.  
  Later Values Replace earlier Values for the same Key.  
  Global Constants never Enter this Override Chain.
- Load and Validate once at Startup; Pass typed Configuration onward.  
  Reject Unknown Keys, Invalid Types and Missing required Values.

For this Example, Files Hold flat Keys.  
An Absent Key Inherits its Default; `null` is Invalid.

The Loader Declares these Environment Variables:

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | File Selection | `development` or `production`; Required |
| `ENROLLMENT_REQUEST_TIMEOUT_SECONDS` | `requestTimeoutSeconds` | Positive Integer |
| `ENROLLMENT_API_TOKEN` | `apiToken` | Nonempty Secret; Required in Production |

Validate `serviceUrl` as an HTTP or HTTPS URL,  
`retryLimit` as a nonnegative Integer,  
and `studentLimit` as a positive Integer.

With `APP_ENV=production` and  
`ENROLLMENT_REQUEST_TIMEOUT_SECONDS=20`,  
the Resolved nonsecret Values are:

```yaml
requestTimeoutSeconds: 20
retryLimit: 2
studentLimit: 30
serviceUrl: https://enrollment.example.com
```

The Loader Implements Selection, Overrides and Validation.  
YAML, JSON and TOML only Describe Data;  
they do not Expand Environment Variables on their own.

## The Declared Surface

> Proposed. Drawn from one Project, cangrejometralleta/muchi-api,
> and not yet Weathered by a second.

- A Variable nobody Declared is a Variable nobody Finds.  
  `.env.example` is the Contract, and it is Committed.
- Required Variables Live uncommented, with a Placeholder Value.  
  A Reader Copies the File and Sees at once what Startup Wants.
- Optional Overrides Live commented, Showing the Default they Replace.  
  The Comment Documents; it never Loads.
- `.env` is Local and Ignored by git.  
  The Loader Reads it when Present and Says nothing when Absent.
- An exported Variable Wins. The File Fills the Gaps it Left.  
  A Shell that Sets a Value Meant it; a File Only Suggested one.

```sh
# Required. No Default Exists, and Startup Stops without them.
APP_ENV=development
TOKEN_SECRET=replace-with-a-random-32-byte-secret

# Optional Overrides; the Defaults shown below Apply when unset.
# SCHOOL_PORT=8080
# SERVER=stdlib
```

Three Sources, and the Order never Changes:
the Argument, then the Environment, then the Default.
Each one is More explicit than the one behind it.

## Six Kinds, in three Pairs

> Proposed. Four are Canon above. Domain Data and Deployment Topology
> Come from one Project and Await a second.

A Value Belongs to exactly one Kind. Asking which one Answers
where it Lives, who may Change it, and what Breaks if it Moves.

**What the Code Means.** Changes when the Code Changes.

| Kind | Lives in |
| --- | --- |
| Global Constant | `constants/`, a descriptive File |
| Algorithmic Constant | beside its Logic |

**What the Deployment Chooses.** Changes when the Deployment Changes.

| Kind | Lives in |
| --- | --- |
| Environment Configuration | `config/school.production.json` |
| Deployment Topology | `config/deploy.env` |

**What Comes from outside.** Changes when the World Changes.

| Kind | Lives in |
| --- | --- |
| Domain Data | `config/`, its own File |
| Secret | Injection only, never a File |

- **Domain Data Changes without a Code Change and without an Environment.**  
  A per-Store Scraping Rule, a Tax Table, a Holiday Calendar.  
  It Reads like Configuration and Deploys like Content.
- **Deployment Topology never Reaches the running Program.**  
  Region, Memory, Instance Count, Service Account.  
  The Deploy Script Reads it; the Application never Opens the File.
- Six Kinds do not Fit in a Head. Three Pairs do.  
  Ask which Pair first, then which of the two.

## Constants Stay with their Meaning

A Unit Conversion Belongs beside the Calculation:

```go
const secondsPerMinute = 60
```

A shared State Catalog Belongs in `constants/`.  
A Timeout Belongs in `config/`.  
A Conversion Factor Belongs beside its Calculation.

Secrets Enter through Injection, never through either Directory.

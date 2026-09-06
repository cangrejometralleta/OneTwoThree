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

## Constants Stay with their Meaning

A Unit Conversion Belongs beside the Calculation:

```go
const secondsPerMinute = 60
```

A shared State Catalog Belongs in `constants/`.  
A Timeout Belongs in `config/`.  
A Conversion Factor Belongs beside its Calculation.

Secrets Enter through Injection, never through either Directory.

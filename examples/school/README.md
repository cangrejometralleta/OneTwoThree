# School Service

The Technical Test, Rewritten twice.
Students, Courses, a Chilean RUT and a Token,
served by six Frameworks that never Touch the Business.

| | Go | TypeScript |
|---|---|---|
| Run | `APP_ENV=development go run .` | `APP_ENV=development npm start` |
| Test | `go test ./...` | `npm test` |
| Store | GORM over SQLite | `node:sqlite` |
| Frameworks | net/http, chi, gin | node:http, express, fastify |

```
APP_ENV=development SERVER=gin TOKEN_SECRET=s go run .       # Go
APP_ENV=development SERVER=fastify TOKEN_SECRET=s npm start   # TypeScript
```

## Constants and Environments

Go and TypeScript Read the same JSON Files, following
[Constants](../../rules/constants.md).
`minimumAgeYears: 18` is the shared Enrollment Rule.
RUT Weights and Modulo eleven Stay beside their Algorithm.

```text
constants/school.json              Global Constants; no Overrides
config/school.defaults.json        Overridable Deployment Defaults
config/school.development.json     Local Database Path
config/school.production.json      Separate Production Database Path
```

Select `APP_ENV` Explicitly: `development` or `production`.
The Loader Applies Defaults, then the selected File, then declared Variables.
Global Constants Load separately and never Enter that Chain.
Files are Read once at Startup; Restart after a Change.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | Environment File | Required; `development` or `production` |
| `SCHOOL_PORT` | `port` | Integer 1–65535; Default 8080 |
| `SCHOOL_DATABASE_PATH` | `databasePath` | Nonempty String |
| `SCHOOL_TOKEN_LIFE_SECONDS` | `tokenLifeSeconds` | 1–86400 Seconds; Default 600 |
| `SERVER` | `serverAdapter` | Go: `stdlib`, `chi`, `gin`; TypeScript: `stdlib`, `node`, `express`, `fastify` |
| `TOKEN_SECRET` | Secret Injection | Required in every Environment; no File Key |

Unknown Keys, Nulls, Invalid Types and Missing required Values Stop Startup
before a Store or Listener Opens. Unknown `SCHOOL_` Variables also Fail,
including Attempts to Override Global Constants.

`stdlib` Selects net/http in Go and node:http in TypeScript.
`TOKEN_SECRET` is Supplied by the Deployment, never Committed to JSON.
The Secret in the Run Example is for local Demonstration only.

Run Go from `go/` and TypeScript from `ts/`; Build TypeScript with
`npm run build` first. Go Reads Data from `../constants/` and `../config/`;
TypeScript Resolves those Directories relative to its Module.
Keep both Directories beside the Language Directories when Packaging.
Relative Database Paths Resolve from the Working Directory.

## The Shape, in both Languages

```
main            Casts the Players and Picks an Adapter
adapters        THE ONLY FILES THAT IMPORT A FRAMEWORK
transport       Request, Response, Handler, Route
handlers        the Script. Framework-free
providers       the Interfaces the Core Declares
store           THE ONLY FILES THAT IMPORT AN ORM
dto             the Wire Shapes
constants       the Global Values, separate from Deployment
config          the Environment Values, Validated at Startup
settings        the strict JSON Loader
domain          the Business Truth
```

Two Files Hold every Vendor.
Swapping one Edits one of them, and nothing else.

## Providers

A Provider is an Interface the Core Declares.
`StudentStore` Names a Need. `GormSchool` Fills it.
The Handlers never Learn which.

Each Provider Carries the URL of the Contract it Wraps,
so a Reader Chasing a Detail never Leaves the File.

The Word Collides. Angular, NestJS and Terraform
all Mean something else by Provider.
The Literature Calls this a Port.

## Three Shapes of one Student

| Shape         | Trusted | Speaks   |
|---------------|---------|----------|
| `StudentBody` | no      | the Wire |
| `Student`     | yes     | Business |
| `StudentRow`  | yes     | Storage  |

The Entity is never the DTO.
A Client Cannot Set an Identity by Sending one.

## Typing as Documentation

Go Names its Types. TypeScript Brands them.
Both Refuse to Pass a `RUT` where a `FullName` Belongs,
and both Compile down to a String.

## What the Tests Buy

`SchoolAPI` Depends on three Interfaces, never on a Library.
So the Tests Hand it three Maps and Finish in Milliseconds,
with no Database and no Port.

One of them Pins the Comparison that Broke the original:
a Token Dies **once** its Deadline Passes, never before.

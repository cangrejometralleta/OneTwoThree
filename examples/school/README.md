# School Service

The Technical Test, Rewritten three times.
Students, Courses, a Chilean RUT and a Token,
served by seven Frameworks that never Touch the Business.

[The Before](BEFORE.md) Reads the original Test and Names the Rule
each Seam Earned. Read that first if you Want the Argument
instead of the Conclusion.

Every Directory Answers the same two Scripts.
Learn them once and every Runtime Opens the same Way.

```bash
cp ../.env.example .env             # once, beside the Program you Run
./build.sh                          # the Gates, then the Artefact
./run.sh [adapter]                  # the Service
```

[.env.example](.env.example) is the Contract: Required Variables live,
optional ones commented beside the Default they Replace.
`run.sh` Reads `.env` when Present, and an exported Variable Wins.
The Argument Wins over both.

| | Go | TypeScript | Java |
|---|---|---|---|
| Gates | gofmt, vet, test | install, tsc, test | toolchain, test |
| Store | GORM over SQLite | `node:sqlite` | JPA over H2 |
| Adapters | stdlib, chi, gin | stdlib, node, express, fastify | stdlib |

`build.sh` Refuses to Build what does not Pass.
`run.sh` Refuses to Start without a Secret.
Both Say which Adapter they Know, and Name the one you Asked for.

The [Specification](SPEC.md) Says what all three Answer.

```bash
TOKEN_SECRET=s go/run.sh gin        # Go
TOKEN_SECRET=s ts/run.sh fastify    # TypeScript
TOKEN_SECRET=s java/run.sh          # Java
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

Each Script Enters its own Directory first, so it Runs from anywhere.
Go and Java Walk up until they Find `constants/` and `config/`,
so a Package Test Reads the same Data;
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

In Go those Files are Grouped into Packages,
so the Compiler Guards the Boundary the Rule Describes.
TypeScript Keeps the flat Shape above.

```text
main.go         Casts the Players; the only File that Knows every Package
transport/      Request, Response, Handler, Route; Imports nothing
wire/           the Contract: every Shape a Client Sends or Receives
faults/         the controlled Failures, each Carrying its Answer
school/         the Core: domain, providers, constants. No HTTP at all
app/            the Application Layer: the Crossing, Validation and Context
api/            the Script: Handlers that Answer or Fail
settings/       the strict JSON Loader and the validated SchoolConfig
store/          THE ONLY PACKAGE THAT IMPORTS AN ORM
serving/        THE ONLY PACKAGE THAT IMPORTS A FRAMEWORK
tokens/         the Token Adapter, standard Library only
```

`go list -deps ./school` Names no Vendor and no Socket.
The Core Cannot Import gorm, chi, gin or net/http, because it never Sees them.

## The Contract, on its own

`wire/` Holds the Shapes a Client Sends and Receives, and nothing else.
It Imports nothing, so Reading it Costs no Context.
Open one File and you Know the whole API.

The Crossing Lives with the Business Types instead, in `school/wire.go`.
A Shape Stays a Shape; the Core Owns the Promotion.

## Controlled Failures

A Business Error Declares its Answer beside its Reason,
so no Handler ever Chooses a Number.

```go
var ErrRutTaken = faults.ReportTakenValue("rut is already Registered")
```

A Handler never Builds a Reply. It Answers with a Value, or it Fails:

```go
func (a SchoolAPI) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.Students.SelectStudentRow(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return school.RenderStudentView(student), nil
}
```

One Function Turns that into a Reply, and it is the only one in the Program:

```go
func AnswerWith(status int, tell Telling) transport.Handler
```

A Failure that Carries no Fault was never Controlled.
It Answers five hundred, because it is Ours and not the Caller's.
A Driver Error Reaches the Edge as a five hundred, never as a four hundred.

No Package and no Type Carries a Vendor Name.
`store` Says the Responsibility, `store_gorm.go` Says who Fulfils it.
A Filename can Name a Guest; a Construct Names the Business.

## The Application Layer

Three Things Live between the Transport and the Business,
and none of them is a Business Rule.

- **The Crossing.** `AnswerWith` Turns a Telling into a Handler.
  The Route Declares the happy Status; a Fault Declares its own.
- **Validation.** `ReadPathNumber`, `ReadJSONBody` and `ReadPageRequest`
  Check the Form of a Request, never its Meaning.
- **Context.** `RequireProvenCaller` Names the Caller before the Story Starts.
  A Handler behind it Reads `req.Caller` and Trusts it.

The Libretto Says all of it at a Glance:

```go
{Method: "GET", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.ShowStudentRecord)},
```

Every Route but `POST /token` Names its Caller.
An Unnamed Caller Gets four hundred and one and Learns nothing else.

## Providers

A Provider is an Interface the Core Declares.
`StudentStore` Names a Need. `store.School` Fills it.
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

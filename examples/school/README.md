# School Service

The technical test, Rewritten three times.
Students, courses, a chilean RUT and a token,
served by seven frameworks that never Touch the Business.

[The Before](BEFORE.md) Reads the original test and Names the rule
each seam earned. Read that first if you want the Argument
instead of the Conclusion.

Every directory Answers the same two Scripts.
Learn them once and every runtime Opens the same Way.

```bash
cp ../.env.example .env             # once, beside the Program you Run
./build.sh                          # the Gates, then the Artefact
./run.sh [adapter]                  # the Service
```

[.env.example](.env.example) is the Contract: required variables live,
optional ones commented beside the default they Replace.
`run.sh` Reads `.env` when present, and an exported variable Wins.
The Argument Wins over both.

| | Go | TypeScript | Java |
|---|---|---|---|
| Gates | gofmt, vet, test | install, tsc, test | toolchain, test |
| Store | GORM over SQLite | `node:sqlite` | JPA over H2 |
| Adapters | stdlib, chi, gin | stdlib, node, express, fastify | stdlib |

`build.sh` Refuses to Build what does not Pass.
`run.sh` Refuses to Start without a Secret.
Both Say which adapter they know, and Name the one you asked for.

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
RUT weights and modulo eleven Stay beside their Algorithm.

```text
constants/school.json              Global Constants; no Overrides
config/school.defaults.json        Overridable Deployment Defaults
config/school.development.json     Local Database Path
config/school.production.json      Separate Production Database Path
```

Select `APP_ENV` Explicitly: `development` or `production`.
The loader Applies Defaults, then the selected file, then declared variables.
Global constants Load separately and never Enter that chain.
Files are Read once at Startup; restart after a change.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | Environment File | Required; `development` or `production` |
| `SCHOOL_PORT` | `port` | Integer 1–65535; Default 8080 |
| `SCHOOL_DATABASE_PATH` | `databasePath` | Nonempty String |
| `SCHOOL_TOKEN_LIFE_SECONDS` | `tokenLifeSeconds` | 1–86400 Seconds; Default 600 |
| `SERVER` | `serverAdapter` | Go: `stdlib`, `chi`, `gin`; TypeScript: `stdlib`, `node`, `express`, `fastify` |
| `TOKEN_SECRET` | Secret Injection | Required in every Environment; no File Key |

Unknown keys, nulls, invalid types and missing required values Stop Startup
before a store or listener opens. Unknown `SCHOOL_` variables also Fail,
including attempts to Override global constants.

`stdlib` Selects net/http in Go and node:http in TypeScript.
`TOKEN_SECRET` is Supplied by the deployment, never Committed to JSON.
The Secret in the run example is for local Demonstration only.

Each script Enters its own directory first, so it Runs from anywhere.
Go and Java walk up until they Find `constants/` and `config/`,
so a package test reads the same data;
TypeScript Resolves those directories relative to its module.
Keep both directories beside the language directories when Packaging.
Relative database paths Resolve from the working Directory.

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

Two files Hold every Vendor.
Swapping one Edits one of them, and nothing else.

In Go those files are grouped into packages,
so the compiler Guards the Boundary the rule describes.
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
The core cannot Import gorm, chi, gin or net/http, because it never Sees them.

## The Contract, on its own

`wire/` Holds the shapes a client Sends and Receives, and nothing else.
It Imports nothing, so reading it Costs no context.
Open one File and you Know the whole API.

The crossing Lives with the business Types instead, in `school/wire.go`.
A shape stays a shape; the core Owns the Promotion.

## Controlled Failures

A business error Declares its Answer beside its reason,
so no handler ever chooses a number.

```go
var ErrRutTaken = faults.ReportTakenValue("rut is already Registered")
```

A handler never Builds a Reply. It Answers with a value, or it Fails:

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

One function Turns that into a Reply, and it is the only one in the program:

```go
func AnswerWith(status int, tell Telling) transport.Handler
```

A failure that Carries no fault was never Controlled.
It answers five hundred, because it is Ours and not the Caller's.
A driver error Reaches the Edge as a five hundred, never as a four hundred.

No package and no type Carries a vendor Name.
`store` Says the responsibility, `store_gorm.go` Says who fulfils it.
A filename can Name a guest; a construct Names the business.

## The Application Layer

Three things Live between the transport and the business,
and none of them is a business Rule.

- **The Crossing.** `AnswerWith` Turns a Telling into a Handler.
  The route Declares the happy status; a fault Declares its own.
- **Validation.** `ReadPathNumber`, `ReadJSONBody` and `ReadPageRequest`
  Check the Form of a request, never its Meaning.
- **Context.** `RequireProvenCaller` Names the Caller before the story starts.
  A handler behind it Reads `req.Caller` and Trusts it.

The libretto Says all of it at a Glance:

```go
{Method: "GET", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.ShowStudentRecord)},
```

Every route but `POST /token` Names its Caller.
An unnamed caller Gets four hundred and one and Learns nothing else.

## Providers

A provider is an Interface the core Declares.
`StudentStore` Names a Need. `store.School` Fills it.
The Handlers never Learn which.

Each provider Carries the URL of the Contract it wraps,
so a reader chasing a detail never leaves the file.

The Word Collides. Angular, NestJS and Terraform
all Mean something else by Provider.
The literature Calls this a Port.

## Three Shapes of one Student

| Shape         | Trusted | Speaks   |
|---------------|---------|----------|
| `StudentBody` | no      | the Wire |
| `Student`     | yes     | Business |
| `StudentRow`  | yes     | Storage  |

The Entity is never the DTO.
A client cannot Set an Identity by sending one.

## Typing as Documentation

Go Names its Types. TypeScript Brands them.
Both Refuse to pass a `RUT` where a `FullName` belongs,
and both Compile down to a string.

## What the Tests Buy

`SchoolAPI` Depends on three Interfaces, never on a Library.
So the tests Hand it three maps and Finish in milliseconds,
with no database and no port.

One of them pins the comparison that broke the original:
a token Dies **once** its deadline Passes, never before.

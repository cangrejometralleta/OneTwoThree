# School Service

The technical test, rewritten three times.
Students, courses, a chilean RUT and a token,
served by runtimes that never touch the Business.

The Go, Java and TypeScript services are generated from the Contract, the
Specification and the Shape in [specs/school](../specs/school), and they
answer the same requests the same way.

[The Before](BEFORE.md) reads the original test and names the Rule each seam
earned. Read that first if you want the argument instead of the conclusion.

Every directory answers the same two scripts. Learn them once and every
runtime opens the same way.

```bash
cp ../.env.example .env             # once per runtime directory
./build.sh                          # the gates, then the artefact
./run.sh [adapter]                  # stdlib is the default everywhere
```

[.env.example](.env.example) is the Contract: required variables live,
optional ones commented beside the default they replace. `run.sh` reads
`.env` when present, and an exported variable wins over the file. An argument
you pass wins over both.

| | Go | TypeScript | Java |
| --- | --- | --- | --- |
| Gates | gofmt, vet, 87 tests | strict types, 91 tests | toolchain, lint, 93 tests |
| Store | GORM over SQLite | Drizzle over SQLite | JDBC over H2 |
| Adapters | stdlib, chi, gin | stdlib, express | stdlib, spring |

`build.sh` refuses to build what does not pass. `run.sh` refuses to start
without a Secret. Every runtime selects its adapter at startup.

The [Contract](../specs/school/openapi.yaml), the
[Specification](../specs/school/SPEC.md) and [the Shape](../specs/school/SHAPE.md)
say what a generated service answers, and all three do today. The TypeScript
[regeneration notes](ts/REGENERATION.md) record its implementation choices.

```bash
TOKEN_SECRET=s go/run.sh gin        # Go
TOKEN_SECRET=s ts/run.sh            # TypeScript
TOKEN_SECRET=s java/run.sh          # Java
```

The TypeScript runtime needs Node.js 22.6 or newer and nothing beyond its
standard library. Its build checks syntax and types and runs 18 behavior
tests. A socket smoke test could not run in the sandbox; the
[regeneration notes](ts/REGENERATION.md) say why.

## Constants and Environments

Go, Java and TypeScript read the same JSON files, following
[Constants](../../rules/constants.md). `minimumAgeYears: 18` is the shared
Enrollment rule. RUT weights and modulo eleven stay beside their Algorithm.

```text
constants/school.json              global constants; no overrides
config/school.defaults.json        overridable deployment defaults
config/school.development.json     local database path
config/school.production.json      separate production database path
```

Select `APP_ENV` explicitly: `development` or `production`. The Loader applies
defaults, then the selected file, then declared variables. Global constants
load separately and never enter that chain. Files are read once at startup,
so restart after a change.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | environment file | required; `development` or `production` |
| `SCHOOL_PORT` | `port` | integer 1–65535; default 8080 |
| `SCHOOL_DATABASE_PATH` | `databasePath` | nonempty string |
| `SCHOOL_REGISTRY_PATH` | `registryPath` | nonempty string; the file must exist; relative paths resolve from the data root |
| `SCHOOL_TOKEN_LIFE_SECONDS` | `tokenLifeSeconds` | 1–86400 seconds; default 600 |
| `SERVER` | `serverAdapter` | Go: `stdlib`, `chi`, `gin`; TypeScript: `stdlib` |
| `TOKEN_SECRET` | secret injection | required in every environment; no file key |

Unknown keys, nulls, invalid types and missing required values stop startup
before a store or listener opens. Unknown `SCHOOL_` variables also fail,
including attempts to override global constants.

`stdlib` selects net/http in Go and node:http in TypeScript. `TOKEN_SECRET`
comes from the deployment and is never committed to JSON. The secret in the
run example is for local demonstration only.

Each script enters its own directory first, so you can run it from anywhere.
Go and Java walk up until they find `constants/` and `config/`, so a package
test reads the same data. TypeScript resolves those directories relative to
its module. Keep both directories beside the language directories when you
package. Relative database paths resolve from the working directory.

## The Shape, in both Languages

```text
main            casts the players and picks an adapter
adapters        THE ONLY FILES THAT IMPORT A FRAMEWORK
transport       Request, Response, Handler, Route
handlers        the script, framework-free
providers       the interfaces the core declares
store           THE ONLY FILES THAT IMPORT AN ORM
dto             the wire shapes
constants       the global values, separate from deployment
config          the environment values, validated at startup
settings        the strict JSON loader
domain          the business truth
```

Adapters hold every Vendor, and one place composes them. Swapping one edits
its adapter and the composition.

In Go each Vendor gets a package of its own, so the Compiler guards the
Boundary the rule describes. Swapping one edits its adapter and `main.go`,
the one composition. The consumer declares each Provider it calls: `app/`
declares services, store ports and the token issuer; `handler/` declares
routes. TypeScript keeps the flat shape above.

```text
main.go         casts the players and picks the server; the only file that knows every package
transport/      Request, Response, Handler, Route; imports nothing
wire/           the contract: every shape a client sends or receives
faults/         the controlled failures, each carrying its answer
school/         the core: domain, page, admission constants. No providers, no vendor
app/            services, store ports, TokenIssuer, crossing, validation and CallerReader
handler/        the script: routes, handler methods and wire mappings
settings/       the strict JSON loader and the validated SchoolConfig
store/          THE ONLY PACKAGE THAT IMPORTS AN ORM
campus/         the registry and the notice: one office filling both ports
serving/        the stdlib server and the shared request reading
serving/chiserver/  THE ONLY PACKAGE THAT IMPORTS chi
serving/ginserver/  THE ONLY PACKAGE THAT IMPORTS gin
tokens/         the token adapter, standard library only
```

`go list -deps ./school` names no wire package, vendor or socket. The Core
cannot import wire, gorm, chi, gin or net/http, because it never sees them.

## The Contract, on its own

`wire/` holds the shapes a client sends and receives, and nothing else. It
imports nothing, so reading it costs no context. Open one file and you know
the whole API.

The wire package holds declarations only. `handler/model_dto.go` promotes
request shapes and renders response shapes, leaving the Core free of
transport imports.

## Controlled Failures

A business Error declares its Answer beside its reason, so no handler ever
chooses a number.

```go
var ErrRutTaken = faults.ReportTakenValue("rut is already Registered")
```

In Go the answer is a Kind rather than a number, and `faults/` names no
protocol. `app/status.go` is the one table that turns a Kind into an HTTP
Status, so the Core stays free of `net/http`.

A Handler never builds a reply. It calls the application Service and answers
with a value, or it fails:

```go
func (a Handler) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.School.ReadStudent(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return renderStudentView(student), nil
}
```

One function turns that into a reply, and it is the only one in the program:

```go
func AnswerWith(status int, tell Telling) transport.Handler
```

A failure that carries no Fault was never controlled. It answers five
hundred, because it is ours and not the caller's. A driver error reaches the
edge as a five hundred, never as a four hundred.

No package and no type carries a vendor name. `store` says the
responsibility, `store_gorm.go` says who fulfils it. A filename can name a
Guest; a construct names the Business.

## The Application Layer

Four things live between the transport and the business, and none of them is
a business rule.

- **Composition.** `SchoolService` joins store ports and owns the use cases.
- **The Crossing.** `AnswerWith` turns a Telling into a Handler.
  The route declares the happy status; a fault declares its own.
- **Validation.** `ReadPathNumber`, `ReadJSONBody` and `ReadPageRequest`
  check the form of a request, never its meaning.
- **Context.** `RequireProvenCaller` names the Caller before the story starts.
  A handler behind it reads `req.Caller` and trusts it.

The libretto says all of it at a glance:

```go
{Method: "GET", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.ShowStudentRecord)},
```

Every route but `POST /token` names its Caller. An unnamed caller gets four
hundred and one and learns nothing else.

## Providers

A Provider is an Interface its Consumer declares. `app.SchoolService` needs
`StudentStore` and `CourseStore`; `store.School` fills both. The Handler calls
the Service and never learns which store fulfils it.

Each Provider carries the URL of the contract it wraps, so a reader chasing a
detail never leaves the file.

The word collides. Angular, NestJS and Terraform all mean something else by
Provider. The literature calls this a Port.

## Three Shapes of one Student

| Shape | Trusted | Speaks |
| --- | --- | --- |
| `StudentBody` | no      | the Wire |
| `Student`     | yes     | Business |
| `StudentRow`  | yes     | Storage  |

The Entity is never the DTO. A client cannot set an Identity by sending one.

## Typing as Documentation

Go names its types. TypeScript brands them. Both refuse to pass a `RUT`
where a `FullName` belongs, and both compile down to a string.

## What the Tests Buy

`app.SchoolService` depends on store ports, and `handler.Handler` depends on
that Service plus the token issuer. Tests supply fakes through the same
composition, with no database and no network port.

One of them pins the comparison that broke the original: a token dies
**once** its deadline passes, never before.

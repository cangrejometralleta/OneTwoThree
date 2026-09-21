# School — Java

The same [Specification](../SPEC.md) as Go and TypeScript,
in the Language [the Before](../BEFORE.md) was Written in.

```bash
cp ../.env.example .env          # once
./build.sh                       # Toolchain and Tests, then the Jar
./run.sh                         # the Service
```

Spring Boot 3, JPA over H2, HMAC from the standard Library.
Read [the Before](../BEFORE.md) first, then this. Same author, same domain,
and every difference Earned a Rule.

## The Shape

```text
Main.java       Casts the Players; the only Class that Knows every Package
transport/      Request, Response, Handler, Route; Imports nothing
wire/           the Contract: every Shape a Client Sends or Receives
faults/         the controlled Failures, each Carrying its Answer
school/         the Core: domain, providers, constants. No HTTP, no Spring
app/            the Application Layer: the Crossing, Validation and Context
api/            the Script: Handlers that Answer or Fail
settings/       the strict JSON Loader and the validated SchoolConfig
store/          THE ONLY PACKAGE THAT IMPORTS JPA
serving/        THE ONLY PACKAGE THAT IMPORTS SPRING WEB
tokens/         the Token Adapter, javax.crypto only
```

Not one Annotation Reaches `school/`.
No `@Entity`, no `@Autowired`, no `@RestController`.
The Core Compiles against the JDK and Jackson, and against nothing else.

## What the Before did differently

| The Before | Here |
| --- | --- |
| `@RequestBody @Valid Student` Binds the Entity | `wire.StudentBody` Crosses into `school.Student` |
| `@ResponseStatus` on four Exception Classes | `Fault` Carries the Status as a Value |
| `@RestController` per Resource | one Route Table, mounted by `serving` |
| `util/` Holds JWT, RUT and IP | `tokens/`, `school/`, `app/` |
| `@Autowired` Fields | Constructor Arguments, Cast in `Main` |

The `@ResponseStatus` idea Survived the rewrite; only its shape Changed.
It was already Right, and [faults](../go/faults/faults.go) is its Descendant.

## Notes a Reader will Want

- `SERVER=stdlib` Selects the one adapter this runtime Ships.
  The config Names a role; only the class Names Spring.
- Global constants and deployment config Come from `../constants` and `../config`,
  the same files Go and TypeScript Read.
- `useIncrementalCompilation` is Off. Maven's incremental compiler
  Emits stale records under a newer JDK, and a stale record Fails oddly.
- The build Wants a JDK; a JRE Reports "No compiler is provided".

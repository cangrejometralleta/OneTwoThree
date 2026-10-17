# School — Java

The same [Contract](../../specs/school/openapi.yaml) and [Specification](../../specs/school/SPEC.md) as Go and TypeScript,
Generated from them and [the Shape](../../specs/school/SHAPE.md),
in the Language [the Before](../BEFORE.md) was Written in.

```bash
cp ../.env.example .env          # once
./build.sh                       # Toolchain, Lint and Tests, then the Jar
./run.sh [stdlib|spring]         # the Service
```

JDK 21 or newer, JDBC over H2, HMAC from `javax.crypto`.
Read [the Before](../BEFORE.md) first, then this. Same author, same domain,
and every Difference Earned a Rule.

## The Shape

```text
Main.java        Casts the Players; the only Class that Knows every Package
transport/       Request, Response, Route, Libretto, Server; Imports nothing of ours
wire/            the Contract: every Shape a Client Sends or Receives
faults/          the controlled Failures, each Carrying its Kind
codec/           THE ONLY PACKAGE THAT IMPORTS JACKSON; Plain Values out
school/          the Core: domain, Rut, Rules and Errors. No HTTP, no JSON, no Spring
app/             the Service, its Ports, the Crossing, Validation and the Status Table
handler/         the Door: the Libretto and its Scripts
campus/          the Registry and the Notifier, one Office Filling both Ports
tokens/          the Token Adapter, javax.crypto only
store/           THE ONLY PACKAGE THAT SPEAKS SQL; JDBC Finds H2 at Runtime
serving/         Reading, Shared by every Framework, so none of them Decides
serving/stdlib/  the JDK's own HttpServer
serving/spring/  THE ONLY PACKAGE THAT IMPORTS SPRING
settings/        the strict JSON Loader and the validated SchoolConfig
```

Not one Annotation Reaches `school/`.
No `@Entity`, no `@Autowired`, no `@RestController`.
The Core Compiles against the JDK and `faults`, and against nothing else;
[BoundariesTest](src/test/java/com/example/BoundariesTest.java) Reads every Import and Says so.

## What the Before did differently

| The Before | Here |
| --- | --- |
| `@RequestBody @Valid Student` Binds the Entity | `wire.StudentBody` Crosses into `school.Student` |
| `@ResponseStatus` on four Exception Classes | `Fault` Carries a Kind; one Table in `app` Turns it into a Status |
| `@RestController` per Resource | one Libretto, Read by `serving` for every Framework |
| `util/` Holds JWT, RUT and IP | `tokens/`, `school/`, `app/` |
| `@Autowired` Fields | Constructor Arguments, Cast in `Main` |
| Spring Boot Owns `main` | Spring is one Server Adapter, Picked at Startup |

The `@ResponseStatus` idea Survived the rewrite; only its shape Changed.
It was already Right, and [faults](src/main/java/com/example/faults/Fault.java) is its Descendant.

## Notes a Reader will Want

- `SERVER=stdlib` or `SERVER=spring` Selects the Framework. Both Pass one
  [Contract Test](src/test/java/com/example/serving/AdapterContractTest.java), because both Hand
  raw Bytes to `serving.Reading` and Write back what it Answers.
- Global constants and deployment config Come from `../constants` and `../config`,
  the same files Go and TypeScript Read. The minimum Age Reaches the Core as `SchoolRules`, Cast in `Main`.
- `databasePath` Names a Base; H2 Appends `.mv.db` to it.
- The JDK Ships no Formatter, so `-Xlint:all -Werror` is the Format Gate.
- `useIncrementalCompilation` is Off. Maven's incremental compiler
  Emits stale records under a newer JDK, and a stale record Fails oddly.
- The Build Wants a JDK; a JRE Reports "No compiler is provided".
  `build.sh` Says which, and how to Point `JAVA_HOME` at one.

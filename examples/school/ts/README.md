# School — TypeScript

The same [Contract](../../specs/school/openapi.yaml) and [Specification](../../specs/school/SPEC.md) as Go and Java,
Generated from them and [the Shape](../../specs/school/SHAPE.md).

```bash
cp ../.env.example .env          # once
./build.sh                       # Toolchain, Packages, Types and Tests
./run.sh [stdlib|express]        # the Service
```

Node 24 or newer: it Strips the Types and Ships the Test Runner, so nothing Compiles to `dist/`.
Drizzle over SQLite, HMAC from `node:crypto`.

## The Shape

```text
main.ts               Casts the Players; the only File that Knows every Module
src/transport/        Request, Response, Route, Libretto, Server; Imports nothing
src/wire/             the Contract: every Shape a Client Sends or Receives
src/faults/           the controlled Failures, each Carrying its Kind
src/codec/            the strict JSON Reader: 64 bit Numbers, no repeated Key
src/school/           the Core: domain, Rut, Rules and Errors. No HTTP, no ORM
src/app/              the Service, its Ports, the Crossing, Validation and the Status Table
src/handler/          the Door: the Libretto and its Scripts
src/campus/           the Registry and the Notifier, one Office Filling both Ports
src/tokens/           the Token Adapter, node:crypto only
src/store/            THE ONLY MODULE THAT IMPORTS DRIZZLE OR THE DRIVER
src/serving/          Reading, Shared by every Framework, and node:http
src/serving/express/  THE ONLY MODULE THAT IMPORTS EXPRESS
src/settings/         the strict JSON Loader and the validated SchoolConfig
```

The Tests Live beside what they Prove. [boundaries.test.ts](src/boundaries.test.ts) Reads every Import
and Says which Vendor may Live where.

## Three Shapes, one Record

`StudentBody` is the Wire, `Student` the Business and `studentRows` the Storage.
Drizzle Knows only the third. The Crossings in [handler](src/handler/crossings.ts) and the
[store](src/store/school-store.ts) are the only places one Shape Becomes another.

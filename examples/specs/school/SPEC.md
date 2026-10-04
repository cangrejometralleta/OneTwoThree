# School — the Specification

What the Service Does beyond its Contract, in no Language.
[The Contract](openapi.yaml) Holds the Routes, the Wire and the Faults,
each Fault with its Status and its Reason.
Every Rewrite in this Directory Answers both,
so a reader can Compare shapes instead of behaviours.
[The Shape](SHAPE.md) Says how the Code is Arranged.

A Failure carrying no Fault Answers 500.
A driver Error is Ours, never the Caller's.

## The Rules the Business will not Bend

- A Student Needs a Name, a valid RUT, an age and an existing course.
- The age Floor is a global Constant, read from `examples/school/constants/school.json`.
- A RUT Passes the modulo eleven Check its digit encodes.
  Shape `^[0-9]+-[0-9kK]$`, weights two through seven right to left,
  remainder eleven Means `0` and ten Means `k`.
- A RUT is unique across students, and the store Proves it.
- A Course needs a Code and a Name.
  A Course without a Name Answers Name is Empty, and one without a Code
  Answers Code is Empty, before the store Writes.
- Creating or rewriting a student Checks the course Exists first.

## Admission, Asked of the Registry

A RUT that Passes its Check Digit may still Belong to no one.
The Registry Knows which RUTs are Real, and Admission Asks it.

- Creating or rewriting a student Asks the Registry after the course Check
  and before the store Writes. The Order is Form, Course, Registry, Store.
- A RUT the Registry Denies Answers RUT Unregistered.
  A Registry that cannot Answer Answers Registry Unavailable.
  Neither is Guessed: an Unknown is never a Yes.
- Creating a student Announces the Enrollment once the store Keeps it.
  Rewriting Announces nothing.
- An Announcement that Fails never Fails the Enrollment.
  The student Exists, the Failure is Logged, and the Caller Sees a Success:
  it cannot Act on a Notice that did not Leave.
- The Registry and the Announcement are two Questions, and one Office may
  Answer both. The Spec Names neither Vendor nor Format.

## Startup

Nothing Opens before the Configuration Validates.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | Environment File | Required; `development` or `production` |
| `SCHOOL_PORT` | port | Integer 1–65535 |
| `SCHOOL_DATABASE_PATH` | databasePath | Nonempty String |
| `SCHOOL_REGISTRY_PATH` | registryPath | Nonempty String; the File must Exist |
| `SCHOOL_TOKEN_LIFE_SECONDS` | tokenLifeSeconds | 1–86400 |
| `SERVER` | serverAdapter | a Framework this Runtime Implements |
| `TOKEN_SECRET` | Secret Injection | Required; no File Key |

Unknown keys, nulls, wrong types and missing required values Stop Startup.
Global Constants load separately and no Variable can Override them.

## The Token

A bearer Token carries a Subject and a Deadline, Signed with the secret.
The Subject is `student-registry`, and a guarded handler Reads it as the Caller.
It Dies once its deadline Passes, and never before.

The Comparison Reads forward. A token is Dead when `now` is after the Deadline.
Reading it backwards Validates only Expired tokens, which is the bug
[the Before](../../school/BEFORE.md) shipped and every rewrite pins.

## What a Rewrite Must Keep

- The Core names no Vendor and no Socket.
- Two places Hold every vendor: the store and the server.
- The Entity is never the DTO. Three Shapes Carry one Record.
- A Fault Declares its Answer beside its reason, once.
- A Handler answers with a Value or Fails; it builds no reply.

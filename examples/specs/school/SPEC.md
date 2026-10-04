# School — the Specification

What the Service Does, in no Language, Complete enough to Generate it.
[The Contract](openapi.yaml) is the same Truth in machine Form;
this Page Says it for a Reader, and Adds the Rules a Schema cannot Hold.
Every Rewrite in this Directory Answers both,
so a reader can Compare shapes instead of behaviours.
[The Shape](SHAPE.md) Says how the Code is Arranged.

## The Nine Routes

| Method | Path | Success | Guarded |
| --- | --- | --- | --- |
| POST | `/token` | 201 | no |
| GET | `/students` | 200 | yes |
| POST | `/students` | 201 | yes |
| GET | `/students/{id}` | 200 | yes |
| PUT | `/students/{id}` | 200 | yes |
| DELETE | `/students/{id}` | 204 | yes |
| GET | `/courses` | 200 | yes |
| POST | `/courses` | 201 | yes |
| GET | `/courses/{id}` | 200 | yes |

Every Route but `/token` Requires `Authorization: Bearer <token>`.
An unnamed caller Gets 401 and Learns nothing else.
`POST /token` takes no body and Answers a `TokenView`.

## The Wire

```json
StudentBody  { "rut": "", "name": "", "age": 0, "courseId": 0 }
StudentView  { "id": 0, "rut": "", "name": "", "age": 0, "courseId": 0 }
CourseBody   { "code": "", "name": "" }
CourseView   { "id": 0, "code": "", "name": "" }
TokenView    { "token": "" }
FailureView  { "error": "" }
```

- Every field of a Body is Required, and a Body is one JSON Object.
  A missing field, a wrong type, an unknown field or a second value is
  Body is Broken.
- The Body never Carries an Identity. The Path and the Store Own it.
- `{id}` in a Path is a positive whole Number. Anything else is Path is Broken,
  and a well formed Number nobody Holds is Unknown, never Broken.
- A list Route Answers an array of Views, empty as `[]`, never `null`.
  `GET /students?page=0&size=10` Windows the Page, Zero based;
  omit both for the whole set. `size` is at least one and `page` at least zero.
  One alone Windows with the other's Default: `page` 0, `size` 10.
  A page or size that is not a whole number is Page is Invalid, never Zero.
  Three Runtimes Read `?page=abc` three different Ways until this line existed.
- Students List ordered by Name and Courses by Code.
- A Failure Answers one `FailureView`, whose `error` is the Reason its Fault
  Declares, never a driver Message.

## The Thirteen Faults

Each one Carries its Status. No handler Chooses a Number.

| Fault | Status | Reason |
| --- | --- | --- |
| Name is Empty | 400 | `name is Empty` |
| Code is Empty | 400 | `code is Empty` |
| RUT is Invalid | 400 | `rut Fails its Check Digit` |
| Age is too Low | 400 | `age Must be 18 or more` |
| Page is Invalid | 400 | `page Numbers Must not be negative` |
| Body is Broken | 400 | `body is not valid JSON` |
| Path is Broken | 400 | `path Holds no valid Identity` |
| Token is Invalid | 401 | `token is Invalid or Expired` |
| Student Unknown | 404 | `student not Found` |
| Course Unknown | 404 | `course not Found` |
| RUT Taken | 409 | `rut is already Registered` |
| RUT Unregistered | 422 | `rut is not in the Registry` |
| Registry Unavailable | 503 | `registry cannot Answer now` |

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
- A Rewrite Changes the Student the Path Names, and a Student nobody Holds
  is Unknown. A Delete Answers 204 once, and Unknown after.

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
- The Registry is a JSON array of the RUTs that are Real, Read once at Startup.
  A RUT in it that cannot be one Stops Startup. The Notice Names the student
  by Id, never by RUT or Name.

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
Defaults Load first, then the Environment File, then the Variables.

## The Token

A bearer Token carries a Subject and a Deadline, Signed with the secret.
The Subject is `student-registry`, and a guarded handler Reads it as the Caller.
It Dies once its deadline Passes, and never before.
A tampered Token, an unsigned one, one with no Subject and an empty secret
all Answer Token is Invalid.

The Comparison Reads forward. A token is Dead when `now` is after the Deadline.
Reading it backwards Validates only Expired tokens, which is the bug
[the Before](../../school/BEFORE.md) shipped and every rewrite pins.

## What a Rewrite Must Keep

- The Core names no Vendor and no Socket.
- Two places Hold every vendor: the store and the server.
- The Entity is never the DTO. Three Shapes Carry one Record.
- A Fault Declares its Answer beside its reason, once.
- A Handler answers with a Value or Fails; it builds no reply.

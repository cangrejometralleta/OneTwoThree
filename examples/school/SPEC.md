# School — the Specification

What the Service Does, in no Language.
Every Rewrite in this Directory Answers exactly this,
so a reader can Compare shapes instead of behaviours.

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

## The Wire

```json
StudentBody  { "rut": "", "name": "", "age": 0, "courseId": 0 }
StudentView  { "id": 0, "rut": "", "name": "", "age": 0, "courseId": 0 }
CourseBody   { "code": "", "name": "" }
CourseView   { "id": 0, "code": "", "name": "" }
TokenView    { "token": "" }
FailureView  { "error": "" }
```

A list Route Answers an array of Views.
`GET /students?page=0&size=10` Windows the Page; omit both for the whole set.
A page or size that is not a whole number is Invalid, never Zero.
Three Runtimes Read `?page=abc` three different Ways until this line existed.

The Body never Carries an Identity. The Path and the Store Own it.

## The Ten Faults

Each one Carries its Status. No handler Chooses a Number.

| Fault | Status | Reason |
| --- | --- | --- |
| Name is Empty | 400 | `name is Empty` |
| RUT is Invalid | 400 | `rut Fails its Check Digit` |
| Age is too Low | 400 | `age Must be 18 or more` |
| Page is Invalid | 400 | `page Numbers Must not be negative` |
| Body is Broken | 400 | `body is not valid JSON` |
| Path is Broken | 400 | `path Holds no valid Identity` |
| Token is Invalid | 401 | `token is Invalid or Expired` |
| Student Unknown | 404 | `student not Found` |
| Course Unknown | 404 | `course not Found` |
| RUT Taken | 409 | `rut is already Registered` |

A Failure carrying no Fault Answers 500.
A driver Error is Ours, never the Caller's.

## The Rules the Business will not Bend

- A Student Needs a Name, a valid RUT, an age and an existing course.
- The age Floor is a global Constant, read from `constants/school.json`.
- A RUT Passes the modulo eleven Check its digit encodes.
  Shape `^[0-9]+-[0-9kK]$`, weights two through seven right to left,
  remainder eleven Means `0` and ten Means `k`.
- A RUT is unique across students, and the store Proves it.
- A Course needs a Code and a Name.
- Creating or rewriting a student Checks the course Exists first.

## Startup

Nothing Opens before the Configuration Validates.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | Environment File | Required; `development` or `production` |
| `SCHOOL_PORT` | port | Integer 1–65535 |
| `SCHOOL_DATABASE_PATH` | databasePath | Nonempty String |
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
[the Before](BEFORE.md) shipped and all three rewrites pin.

## What a Rewrite Must Keep

- The Core names no Vendor and no Socket.
- Two places Hold every vendor: the store and the server.
- The Entity is never the DTO. Three Shapes Carry one Record.
- A Fault Declares its Answer beside its reason, once.
- A Handler answers with a Value or Fails; it builds no reply.

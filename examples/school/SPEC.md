# School — the Specification

What the Service Does, in no Language.
Every Rewrite in this Directory Answers exactly this,
so a Reader can Compare Shapes instead of Behaviours.

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
An Unnamed Caller Gets 401 and Learns nothing else.

## The Wire

```json
StudentBody  { "rut": "", "name": "", "age": 0, "courseId": 0 }
StudentView  { "id": 0, "rut": "", "name": "", "age": 0, "courseId": 0 }
CourseBody   { "code": "", "name": "" }
CourseView   { "id": 0, "code": "", "name": "" }
TokenView    { "token": "" }
FailureView  { "error": "" }
```

A List Route Answers an Array of Views.
`GET /students?page=0&size=10` Windows the Page; Omit both for the whole Set.
A Page or Size that is not a whole Number is Invalid, never Zero.
Three Runtimes Read `?page=abc` three different Ways until this Line Existed.

The Body never Carries an Identity. The Path and the Store Own it.

## The Ten Faults

Each one Carries its Status. No Handler Chooses a Number.

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

A Failure Carrying no Fault Answers 500.
A Driver Error is Ours, never the Caller's.

## The Rules the Business will not Bend

- A Student Needs a Name, a valid RUT, an Age and an existing Course.
- The Age Floor is a Global Constant, Read from `constants/school.json`.
- A RUT Passes the Modulo eleven Check its Digit Encodes.
  Shape `^[0-9]+-[0-9kK]$`, Weights two through seven Right to Left,
  Remainder eleven Means `0` and ten Means `k`.
- A RUT is unique across Students, and the Store Proves it.
- A Course Needs a Code and a Name.
- Creating or Rewriting a Student Checks the Course Exists first.

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

Unknown Keys, Nulls, wrong Types and missing required Values Stop Startup.
Global Constants Load separately and no Variable can Override them.

## The Token

A Bearer Token Carries a Subject and a Deadline, Signed with the Secret.
The Subject is `student-registry`, and a Guarded Handler Reads it as the Caller.
It Dies once its Deadline Passes, and never before.

The Comparison Reads forward. A Token is Dead when `now` is after the Deadline.
Reading it backwards Validates only Expired Tokens, which is the Bug
[java-before](java-before) Shipped and both Rewrites Pin.

## What a Rewrite Must Keep

- The Core Names no Vendor and no Socket.
- Two Places Hold every Vendor: the Store and the Server.
- The Entity is never the DTO. Three Shapes Carry one Record.
- A Fault Declares its Answer beside its Reason, once.
- A Handler Answers with a Value or Fails; it Builds no Reply.

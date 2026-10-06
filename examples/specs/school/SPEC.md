# School — the Specification

What the Service Does, in no Language, Complete enough to Generate it.
[The Contract](openapi.yaml) is the same Truth in machine Form;
this Page Says it for a Reader, and Adds the Rules a Schema cannot Hold.
Every Rewrite in this Directory Answers both,
so a reader can Compare shapes instead of behaviours.
[The Shape](SHAPE.md) Says how the Code is Arranged.

Two Rewrites from an earlier Spec Agreed on every Test and Differed on
what the Spec left Open. Every Line below Closes one of those Gaps.

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
The scheme Name is not Case Sensitive. An unnamed caller Gets 401 and Learns nothing else.
`POST /token` takes no body and Answers a `TokenView`; Anyone may Mint one.

A Path or a Method no Route Holds Answers Route Unknown, and it Answers before
the Guard, because there is no Route to Guard. `/students/` is not `/students`.

## The Order of Checks

A Request Meets its Checks in one Order, and the first to Fail Answers.

1. The Route, then the Guard: 401 comes before any 400 or 404.
2. The Path, then the Body, then the Page.
3. The Form of the Record: for a Student its Name, then its RUT, then its Age;
   for a Course its Code, then its Name.
4. The Course, then the Registry, then the Store.

A Rewrite Follows the same Order, so a Student nobody Holds is Unknown only at the Store:
`PUT` of a Student nobody Holds with a RUT the Registry Denies Answers RUT Unregistered.

## The Wire

```json
StudentBody  { "rut": "", "name": "", "age": 0, "courseId": 0 }
StudentView  { "id": 0, "rut": "", "name": "", "age": 0, "courseId": 0 }
CourseBody   { "code": "", "name": "" }
CourseView   { "id": 0, "code": "", "name": "" }
TokenView    { "token": "" }
FailureView  { "error": "" }
```

- Every field of a Body is Required, and a Body is one JSON Object, in UTF-8, of at most
  one mebibyte. A missing field, a null, a wrong type, an unknown field, a second value
  or a Body beyond that size is Body is Broken. The `Content-Type` of a Request is Ignored.
- The Body never Carries an Identity. The Path and the Store Own it.
- Whole numbers are 64 bit. A Body number beyond that is Body is Broken.
- `{id}` in a Path is Digits only, a positive whole Number; `012` is 12.
  Anything else is Path is Broken. A well formed Number nobody Holds, even one
  beyond 64 bit, is Unknown, never Broken. A `courseId` nobody Holds, zero included,
  is Course Unknown; the Contract's `minimum` there is Documentation.
- A list Route Answers an array of Views, empty as `[]`, never `null`.
  `GET /students?page=0&size=10` Windows the Page, Zero based;
  omit both for the whole set. `size` is at least one and `page` at least zero.
  One alone Windows with the other's Default: `page` 0, `size` 10.
  Both are Digits only: `+1`, `1e1`, ` 1`, an empty value and `0x1` are Page is Invalid,
  never Zero. A number beyond 64 bit Windows as far as it can: an empty list.
  A repeated key Takes its first value, and an unknown key is Ignored.
  Three Runtimes Read `?page=abc` three different Ways until this line existed.
- Students List ordered by Name and Courses by Code, in binary Order of the text,
  and Ties break by Id Ascending.
- Responses are JSON with `Content-Type: application/json`, except a 204, which Carries no Body.
- A Failure Answers one `FailureView`, whose `error` is the Reason its Fault
  Declares, never a driver Message.

## The Fourteen Faults

Each one Carries its Status. No handler Chooses a Number.

| Fault | Status | Reason |
| --- | --- | --- |
| Name is Empty | 400 | `name is Empty` |
| Code is Empty | 400 | `code is Empty` |
| RUT is Invalid | 400 | `rut Fails its Check Digit` |
| Age is too Low | 400 | `age Must be 18 or more` |
| Page is Invalid | 400 | `page or size is not valid` |
| Body is Broken | 400 | `body is not valid JSON` |
| Path is Broken | 400 | `path Holds no valid Identity` |
| Token is Invalid | 401 | `token is Invalid or Expired` |
| Student Unknown | 404 | `student not Found` |
| Course Unknown | 404 | `course not Found` |
| Route Unknown | 404 | `route not Found` |
| RUT Taken | 409 | `rut is already Registered` |
| RUT Unregistered | 422 | `rut is not in the Registry` |
| Registry Unavailable | 503 | `registry cannot Answer now` |

A Failure carrying no Fault Answers 500 with the Reason `internal error`.
A driver Error is Ours, never the Caller's, and its Message never Leaves.
The Reason of Age is too Low Spells the Floor from the global Constant.

## The Rules the Business will not Bend

- A Student Needs a Name, a valid RUT, an age and an existing course.
- The age Floor is a global Constant, read from `examples/school/constants/school.json`.
- A Name or a Code of only whitespace is Empty. Text is Stored as it Arrived, untrimmed.
- A RUT Passes the modulo eleven Check its digit encodes.
  Shape `^[0-9]+-[0-9kK]$`, weights two through seven right to left,
  remainder eleven Means `0` and ten Means `k`.
- A RUT Enters in lower case: `12345-K` and `12345-k` are one RUT, the View Echoes `k`,
  and the Registry is Asked the same Way. Leading zeros are not Merged.
- A RUT is unique across students, and the store Proves it.
- A Course needs a Code and a Name.
  A Course without a Name Answers Name is Empty, and one without a Code
  Answers Code is Empty, before the store Writes. Codes are not Unique.
- Creating or rewriting a student Checks the course Exists first.
- A Rewrite Changes the Student the Path Names, and a Student nobody Holds
  is Unknown. A Delete Answers 204 once, and Unknown after.
  There is no Delete for a Course.

## Admission, Asked of the Registry

A RUT that Passes its Check Digit may still Belong to no one.
The Registry Knows which RUTs are Real, and Admission Asks it.

- Creating or rewriting a student Asks the Registry after the course Check
  and before the store Writes. The Order is Form, Course, Registry, Store.
- A RUT the Registry Denies Answers RUT Unregistered.
  A Registry that cannot Answer Answers Registry Unavailable.
  Neither is Guessed: an Unknown is never a Yes.
- Registry Unavailable is a Contract of the Port. The File Registry below never
  Raises it, so the Tests Reach it through a Fake.
- Creating a student Announces the Enrollment once the store Keeps it.
  Rewriting Announces nothing.
- An Announcement that Fails never Fails the Enrollment.
  The student Exists, the Failure is Logged, and the Caller Sees a Success:
  it cannot Act on a Notice that did not Leave.
- The Announcement is one Log line that Names the student by Id, as `student <id> is enrolled`,
  and never Holds the RUT or the Name. Its Logger and Format are the Runtime's own.
- The Registry and the Announcement are two Questions, and one Office may
  Answer both. The Spec Names neither Vendor nor Format.
- The Registry is a JSON array of the RUTs that are Real, Read once at Startup.
  An empty array and a Duplicate are Allowed. A File that is Missing, Empty, not an array,
  Holds a second value, or Holds a Text that is not a valid RUT Stops Startup.

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
Defaults Load first from `school.defaults.json`, then the Environment File
`school.<APP_ENV>.json`, then the Variables. Both Files must Exist.

- A Variable that is Set and Empty is Set, and Invalid: never Unset.
  Integers are Digits only, so `+80`, ` 80` and `80.0` Stop Startup.
  A Variable of the `SCHOOL_` Namespace that is not Declared Stops Startup;
  one outside it is Ignored.
- The list of Adapters a Runtime Implements belongs to its Composition,
  and the Settings Loader Receives it, so it names no Adapter.
- A relative `registryPath` Resolves from the data Root, where `constants/` and
  `config/` Live. A relative `databasePath` Resolves from the working Directory
  and Names a Base: a Runtime may Append its own Suffix to it.
- `.env` Sits beside the Program that Runs, in its own Directory.
- A Refusal Exits with a Nonzero Code and Says what Broke.

## The Token

A bearer Token carries a Subject and a Deadline, Signed with the secret.
The Subject is `student-registry`, and a guarded handler Reads it as the Caller.
Any nonempty Subject in a Verified token Names that Caller.

The Format is the same in every Runtime, so a Token Minted by one is Read by another:

- Two base64url Parts without padding, joined by a dot: the Claims, then the Signature.
- The Claims are one JSON object, compact, as `{"sub":"student-registry","exp":4102444800}`,
  with `exp` in whole unix seconds. Minting Sets `exp` to now plus `SCHOOL_TOKEN_LIFE_SECONDS`.
- The Signature is HMAC-SHA256 over the ASCII of the first Part, keyed by the UTF-8
  bytes of the secret, in base64url without padding.

With the secret `spec-secret`, this Token is Valid until the year 2100:

```text
eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjo0MTAyNDQ0ODAwfQ.fgl9OtJeU8f90EIkPjELn34yWKjIbn3QcYn5UFoLAyY
```

and this one, with `exp` 1, is Dead:

```text
eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjoxfQ.V8VQ8Ie_fo88z4EBJ0QbtIc5ONTOqQ_-l_WLuhwz5t4
```

It Dies once its deadline Passes, and never before: Now is Truncated to whole seconds,
and a Token is Dead when that Number is after `exp`. A tampered Token, an unsigned one,
one with no Subject, an unreadable one and an empty secret all Answer Token is Invalid;
an Adapter given an empty secret Refuses to Mint as well.

The Comparison Reads forward. A token is Dead when `now` is after the Deadline.
Reading it backwards Validates only Expired tokens, which is the bug
[the Before](../../school/BEFORE.md) shipped and every rewrite pins.

## What a Rewrite Must Keep

- The Core names no Vendor and no Socket.
- Two places Hold every vendor: the store and the server.
- The Entity is never the DTO. Three Shapes Carry one Record.
- A Fault Declares its Answer beside its reason, once.
- A Handler answers with a Value or Fails; it builds no reply.

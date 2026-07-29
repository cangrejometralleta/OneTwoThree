# School Service

The Technical Test, Rewritten twice.
Students, Courses, a Chilean RUT and a Token,
served by six Frameworks that never Touch the Business.

| | Go | TypeScript |
|---|---|---|
| Run | `go run .` | `npm start` |
| Test | `go test ./...` | `npm test` |
| Store | GORM over SQLite | `node:sqlite` |
| Frameworks | net/http, chi, gin | node:http, express, fastify |

```
SERVER=gin  TOKEN_SECRET=s go run .          # Go
SERVER=fastify TOKEN_SECRET=s npm start      # TypeScript
```

## The Shape, in both Languages

```
main            Casts the Players and Picks an Adapter
adapters        THE ONLY FILES THAT IMPORT A FRAMEWORK
transport       Request, Response, Handler, Route
handlers        the Script. Framework-free
providers       the Interfaces the Core Declares
store           THE ONLY FILES THAT IMPORT AN ORM
dto             the Wire Shapes
domain          the Business Truth
```

Two Files Hold every Vendor.
Swapping one Edits one of them, and nothing else.

## Providers

A Provider is an Interface the Core Declares.
`StudentStore` Names a Need. `GormSchool` Fills it.
The Handlers never Learn which.

Each Provider Carries the URL of the Contract it Wraps,
so a Reader Chasing a Detail never Leaves the File.

The Word Collides. Angular, NestJS and Terraform
all Mean something else by Provider.
The Literature Calls this a Port.

## Three Shapes of one Student

| Shape         | Trusted | Speaks   |
|---------------|---------|----------|
| `StudentBody` | no      | the Wire |
| `Student`     | yes     | Business |
| `StudentRow`  | yes     | Storage  |

The Entity is never the DTO.
A Client Cannot Set an Identity by Sending one.

## Typing as Documentation

Go Names its Types. TypeScript Brands them.
Both Refuse to Pass a `RUT` where a `FullName` Belongs,
and both Compile down to a String.

## What the Tests Buy

`SchoolAPI` Depends on three Interfaces, never on a Library.
So the Tests Hand it three Maps and Finish in Milliseconds,
with no Database and no Port.

One of them Pins the Comparison that Broke the original:
a Token Dies **once** its Deadline Passes, never before.

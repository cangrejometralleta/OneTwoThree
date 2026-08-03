# Order Service

RULES.md Cites `Item`, `SumItemPrices` and `BuildOrderReceipt`
as a loose Snippet, to Show what three Beats Look like.
This Service is that Snippet, Wired to a Store and an HTTP Port.
The Text in RULES.md and the Test in `domain_test.go`
Assert the same Receipt.

| | Go | TypeScript |
|---|---|---|
| Run | `go run .` | `npm start` |
| Test | `go test ./...` | `npm test` |
| Store | GORM over SQLite | `node:sqlite` |
| Framework | net/http | node:http |

```
go run .          # Go, listens on :8081
npm start         # TypeScript, listens on :8081
```

School Proves a Core Portable across six Frameworks.
Order does not Repeat that Proof — one Framework per Language is enough
to Show the Shape Survives a second Domain. Only the Store File
Imports a Vendor; the TypeScript Side Ships with none at Runtime at all,
because `node:sqlite` and `node:http` are the Standard Library, not a Guest.

## The Shape, in both Languages

```
main            Casts the Player and Opens a Store
adapters        THE ONLY FILE THAT IMPORTS A FRAMEWORK
transport       Request, Response, Handler, Route
handlers        the Script. Framework-free
providers       the one Interface the Core Declares
store           THE ONLY FILE THAT IMPORTS A VENDOR
dto             the Wire Shapes
domain          the Business Truth, and the RULES.md Snippet Compiled
```

## An Identity Born before the Row

A Student Waits for the Database to Number it.
An Order cannot: `CheckOrderRecord` Checks `o.ID != ""`
before any Row Exists, so the Handler Mints the Reference first
with `GenerateOrderID` / `generateOrderId`,
then Validates a whole Order, never a half one.

## Items Keep no Identity

`StudentBody` and `Student` Differ because an ID is a Trust Boundary.
An `Item` Carries none, so `OrderBody` and `Order`
Share the same Item Shape unchanged — the Three-Shape Rule
Protects an Identity, not a Value Object.

## What the Tests Buy

`OrderAPI` Depends on one Interface, never on a Library.
`domain_test.go` / `order.test.ts` Pin the exact Receipt Text
RULES.md already Shows, so the Page and the Program
cannot Drift apart unnoticed.

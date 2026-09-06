# Order Service

Rules Cites `Item`, `SumItemPrices` and `BuildOrderReceipt`
as a loose Snippet, to Show what three Beats Look like.
This Service is that Snippet, Wired to a Store and an HTTP Port.
The Text in Rules and the Test in `domain_test.go`
Assert the same Receipt.

| | Go | TypeScript |
|---|---|---|
| Run | `APP_ENV=development go run .` | `APP_ENV=development npm start` |
| Test | `go test ./...` | `npm test` |
| Store | GORM over SQLite | `node:sqlite` |
| Framework | net/http | node:http |

```
APP_ENV=development go run .      # Go, Listens on :8081
APP_ENV=development npm start     # TypeScript, Listens on :8081
```

School Proves a Core Portable across six Frameworks.
Order does not Repeat that Proof — one Framework per Language is enough
to Show the Shape Survives a second Domain. Only the Store File
Imports a Vendor; the TypeScript Side Ships with none at Runtime at all,
because `node:sqlite` and `node:http` are the Standard Library, not a Guest.

## Constants and Environments

Go and TypeScript Read the same JSON Files, following
[Constants](../../rules/constants.md).
`itemNameWidth: 12` and `itemTotalWidth: 6` Define the Receipt Layout.
Percentage Arithmetic and Identity Encoding Stay beside their Algorithms.

```text
constants/order.json              Global Constants; no Overrides
config/order.defaults.json        Overridable Deployment Defaults
config/order.development.json     Local Database Path
config/order.production.json      Separate Production Database Path
```

Select `APP_ENV` Explicitly: `development` or `production`.
The Loader Applies Defaults, then the selected File, then declared Variables.
Global Constants Load separately and never Enter that Chain.
Files are Read once at Startup; Restart after a Change.

| Variable | Destination | Validation |
| --- | --- | --- |
| `APP_ENV` | Environment File | Required; `development` or `production` |
| `ORDER_PORT` | `port` | Integer 1–65535; Default 8081 |
| `ORDER_DATABASE_PATH` | `databasePath` | Nonempty String |

Unknown Keys, Nulls, Invalid Types and Missing required Values Stop Startup
before a Store or Listener Opens. Unknown `ORDER_` Variables also Fail,
including Attempts to Override Global Constants.

Run Go from `go/` and TypeScript from `ts/`; Build TypeScript with
`npm run build` first. Go Reads Data from `../constants/` and `../config/`;
TypeScript Resolves those Directories relative to its Module.
Keep both Directories beside the Language Directories when Packaging.
Relative Database Paths Resolve from the Working Directory.

## The Shape, in both Languages

```
main            Casts the Player and Opens a Store
adapters        THE ONLY FILE THAT IMPORTS A FRAMEWORK
transport       Request, Response, Handler, Route
handlers        the Script. Framework-free
providers       the one Interface the Core Declares
store           THE ONLY FILE THAT IMPORTS A VENDOR
dto             the Wire Shapes
constants       the Global Values, separate from Deployment
config          the Environment Values, Validated at Startup
settings        the strict JSON Loader
domain          the Business Truth, and the Rules Snippet Compiled
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
Rules already Shows, so the Page and the Program
cannot Drift apart unnoticed.

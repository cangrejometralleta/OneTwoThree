# Implementation Choices

What the Spec Leaves to a Runtime, and how this one Answers.

- **Whole Numbers are `bigint`.** A `number` Rounds past 2^53, and the Spec Says 64 bit.
  The Wire, the Core and the Store Carry `bigint`; the JSON Writer Prints its Digits.
- **JSON is Read by `src/codec`.** `JSON.parse` Keeps the last of a repeated Key, Rounds large
  Numbers and cannot Tell `20` from `20.0`. The Reader is standard Library only.
- **The Store is Drizzle over better-sqlite3.** Drizzle Maps the Rows; the Tables are Created at
  Startup, the Way GORM's AutoMigrate Does. `better-sqlite3` Builds a native Binding, and
  `package.json` Allows that one install Script.
- **The Adapters are `stdlib` and `express`.** Both Hand the whole Node Request to `serving/reading.ts`,
  so Express Decides nothing and Answers what node:http Answers.
- **The Announcement is `student <id> is enrolled`**, on standard output.
- **`databasePath` is used as Given**; SQLite Appends nothing.

Go, Java and TypeScript Answer the same fifty edge Requests with the same Status, Content-Type and Body.

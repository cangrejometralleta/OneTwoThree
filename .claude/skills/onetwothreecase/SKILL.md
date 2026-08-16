---
name: onetwothreecase
description: Apply the OneTwoThree manifesto's code conventions when writing, reviewing, or refactoring code — three-beat function bodies, Verb+Noun+context naming, Provider ports named after the business need instead of the vendor, seam-based line breaks, and OneTwoThreeCase doc comments. Use whenever generating or reviewing code in a project that follows RULES.md from cangrejometralleta/OneTwoThree, or when the user asks for "OneTwoThreeCase style" or "the manifesto rules".
---

# OneTwoThreeCase

A portable summary of [RULES.md](https://github.com/cangrejometralleta/OneTwoThree/blob/main/RULES.md).
If a RULES.md file exists in the current repo, it is canonical — this
skill is the checklist, not a replacement.

## Before returning code, check

1. **Structure** — Aim for three-beat functions: Receive, Transform,
   Return. A beat is one thought, not one newline — explicit error
   checks don't count against the three. More beats signal a missing
   abstraction; extract a helper instead of padding one function.

2. **Naming** — `Verb + Noun + context`, three words at most. Case by
   language convention (`PascalCase` exported, `camelCase` unexported,
   `snake_case` Python). A name past three words means the
   responsibility is unclear, not that the name needs to be longer.
   A variable that lives inside three lines can drop to one word — the
   scope already says the rest.

3. **Providers** — Any interface to something outside the core (a
   database, an API, a queue) is named after the business need it
   fills, never the vendor behind it: `StudentStore`, not
   `GormRepository`. One implementation may satisfy several such
   ports; a port must never leak a vendor type, vendor error, or
   vendor import outside its own file.

4. **Script / Handler** — Entry points (HTTP handlers, CLI commands,
   `main`) speak business language only — no driver, query, or socket
   names. Read the function aloud; if it stops sounding like a
   sentence, an abstraction is missing.

5. **Seams** — Break long lines at a real grammatical joint: `&&`,
   `||`, a comma in a list, a dot in a chain. Never break inside a
   unit that reads as one (a call and its single argument). If a
   boolean expression needs a break, name its parts as local variables
   instead of splitting mid-expression.

6. **Anti-patterns to flag** — more than three responsibilities in one
   unit; a name with no verb; a function with no clear return.

7. **Comments** — doc comments narrate in prose, capitalizing the
   Words that carry meaning (Entities, Actions, Statuses) and leaving
   connectors lowercase — the same signal Go already gives with
   exported vs. unexported names, applied to English. At most one
   emoji, only in output or comments, never in an identifier or a key
   the code compares against.

## Example

```go
// BuildOrderReceipt Reads as three Sections: Total, Lines, Result.
func BuildOrderReceipt(id string, items []Item, percent int) string {
	total := SumItemPrices(items)
	total = ApplyMemberRate(total, percent)

	lines := make([]string, 0, len(items))
	for _, it := range items {
		lines = append(lines, FormatItemLine(it))
	}

	lines = append(lines, ReportOrderState(id, total, nil))
	return strings.Join(lines, "\n")
}
```

Three names, three beats, one return — `SumItemPrices` and
`ApplyMemberRate` are Providers of a calculation, `ReportOrderState`
is the Script's own narration.

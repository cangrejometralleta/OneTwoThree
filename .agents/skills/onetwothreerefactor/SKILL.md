---
name: onetwothreerefactor
description: Apply the OneTwoThree manifesto's code conventions when writing or refactoring code and its configuration — structure, naming, Provider boundaries, constants, and shared vendor integration. Use whenever generating or refactoring code in a project that follows Rules from cangrejometralleta/OneTwoThree, or when the user asks for "the manifesto rules". To judge code already written, send the onetwothreereview agent instead — it reads the diff in its own context and returns the findings.
---

# OneTwoThreeRefactor

A portable summary of [Rules](../../../RULES.md).
The canon lives in [rules/](../../../rules/) — each Section its own File.
If a Rules file exists in the current repo, it is canonical — this
skill is the checklist, not a replacement.

This Skill Rides along while you Write.
To Judge Code already Written, Send the
[onetwothreereview](../../agents/onetwothreereview.md) Agent —
it Reads the Diff in its own Context and Returns the Findings.

## Before returning code, check

1. **Structure** — Aim for three-beat functions: Receive, Transform,
   Return. A beat is one thought, not one newline — explicit error
   checks don't count against the three. More beats signal a missing
   abstraction; extract a helper instead of padding one function.

2. **Naming** — `Verb + Noun + context`, three words at most. A name
   past three words means the responsibility is unclear, not that the
   name needs to be longer. A variable that lives inside three lines
   can drop to one word — the scope already says the rest.

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

7. **Comments** — doc comments narrate in prose. At most one emoji,
   only in output or comments, never in an identifier or a key the code
   compares against.

8. **Constants** — Keep shared declarative Values in `constants/`
   and Environment Configuration in `config/`, using YAML, JSON or TOML.
   Global Constants never Enter the Override Chain; Algorithmic Constants
   Stay beside their Logic. Validate at Startup, Select the Environment
   Explicitly, Apply Defaults → Environment File → declared Variables,
   and Inject Secrets separately. Read [Constants](../../../rules/constants.md)
   when Changing these Values or their Loader.

9. **Vendor Integration** — Keep shared Instructions in one Source.
   Use relative Links for supported Vendor Entrances and small Adapters
   for required Formats. Generate unsupported Copies from that Source;
   never Maintain them by Hand. Read
   [Vendor Integration](../../../rules/vendor-integration.md)
   when Changing Tool Configuration or Agent Discovery.

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

## Sources

1. [Structure](../../../rules/structure.md)
2. [Naming](../../../rules/naming.md)
3. [Providers](../../../rules/providers.md)
4. [Script](../../../rules/script.md)
5. [Seams](../../../rules/seams.md)
6. [Anti-Patterns](../../../rules/anti-patterns.md)
7. [Emoji](../../../rules/emoji.md)
8. [Constants](../../../rules/constants.md)
9. [Vendor Integration](../../../rules/vendor-integration.md)

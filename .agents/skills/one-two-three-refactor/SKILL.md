---
name: one-two-three-refactor
description: Apply the OneTwoThree manifesto's code conventions when writing or refactoring code and its configuration — structure, naming, seams, comments, named values, Provider and Layer boundaries, wire/business/storage Shapes, controlled Failures, constants, tests, build scripts, and shared vendor integration. Use whenever generating or refactoring code in a project that follows Rules from cangrejometralleta/OneTwoThree, or when the user asks for "the manifesto rules". To judge code already written, send the one-two-three-review agent instead — it reads the diff in its own context and returns the findings.
---

# OneTwoThreeRefactor

A portable summary of [Rules](../../../RULES.md).
The canon lives in [rules/](../../../rules/) — each Section its own File.
If a Rules file exists in the current repo, it is canonical — this
skill is the checklist, not a replacement.

This Skill Rides along while you Write.
To Judge Code already Written, Send the
[one-two-three-review](../../agents/one-two-three-review.md) Agent —
it Reads the Diff in its own Context and Returns the Findings.

## Before returning code, check

Three Groups. The Line, the Boundary, the Project.

### The Line

1. **Structure** — Aim for three-beat functions: Receive, Transform,
   Return. A beat is one thought, not one newline — explicit error
   checks don't count against the three. More beats signal a missing
   abstraction; extract a helper instead of padding one function.

2. **Naming** — `Verb + Noun + context`, three words at most. A name
   past three words means the responsibility is unclear. A variable
   living inside three lines can drop to one word — the scope says the
   rest. A Construct Names the Responsibility, never the Vendor:
   `store`, not `storegorm`. A Filename may Name the Guest.

3. **Seams** — Break long lines at a real grammatical joint: `&&`,
   `||`, a comma in a list, a dot in a chain. Never break inside a
   unit that reads as one. If a boolean expression needs a break, name
   its parts as local variables instead of splitting mid-expression.

4. **Comments** — Stop at the Claim. A second clause must add a
   constraint, never repeat the first in other words. Test by
   deletion: cut the tail and read the head alone. At most one emoji,
   only in output or comments, never in an identifier.

5. **Values** — A number with a meaning carries a name. Read the
   literal alone, out of its line; if it cannot say what it means, it
   wants a name. Look for an existing constant first — the standard
   library, then the framework, then your own. A constant in the core
   must not drag a vendor in.

### The Boundary

6. **Providers** — Any interface to something outside the core is
   named after the business need it fills, never the vendor behind it:
   `StudentStore`, not `GormRepository`. A port must never leak a
   vendor type, error, or import outside its own file.

7. **Layers** — The core names no vendor and no socket. Let the
   compiler hold the boundary: a package the core cannot import beats
   a rule the core agrees to follow. The shapes the layers speak in
   import nothing. Two places hold every vendor: the store and the
   server.

8. **Shapes** — The Entity is never the DTO. Three shapes carry one
   record: the wire (untrusted, weak types), the business (trusted,
   named types), the storage. Bind the wire shape, never the entity.
   The identity comes from the path or the store, never from the body.

9. **Failures** — A failure the program expected carries its own
   answer. Declare the answer beside the reason, once, and name the
   failure by the case: `ErrRutTaken`, not `Conflict`. A failure
   carrying no answer was never controlled, and answers five hundred.
   One function turns a failure into a number, and the program holds
   exactly one.

10. **Script / Handler** — Entry points speak business language only —
    no driver, query, or socket names. A handler answers with a value
    or it fails; it builds no reply and names no status. The route
    declares the happy status. Read it aloud; if it stops sounding
    like a sentence, an abstraction is missing.

### The Project

11. **Constants** — Sort the Value before you Place it. Six Kinds in
    three Pairs: what the Code Means (Global and Algorithmic Constants),
    what the Deployment Chooses (Environment Configuration and Deployment
    Topology), what Comes from outside (Domain Data and Secrets).
    Global Constants Live in `constants/` and never Enter the Override
    Chain; Algorithmic Constants Stay beside their Logic. Validate at
    Startup, Select the Environment Explicitly, Apply Defaults →
    Environment File → declared Variables, and Inject Secrets separately.

12. **The Declared Surface** — `.env.example` is the Contract and it is
    Committed: Required Variables live, optional ones commented beside
    the Default they Replace. `.env` is local and Ignored. Precedence
    Runs one Way — the Argument, then the Environment, then the Default —
    so an exported Variable Wins over a File.

13. **Tests** — Spell the expectation; never read it from the code
    under test. A test that computes what it checks agrees with itself
    and proves nothing. Prove it by mutation: break the declaration
    and watch the test fail. A test name is a use case, not a method
    name. Arrive the way a caller arrives.

14. **Scripts** — Every program answers `build.sh` and `run.sh`. Build
    refuses to build what does not pass; run refuses to start what
    will fail at startup; run reads `.env` when present. One function
    per step, named by what it checks, with the calls at the bottom
    one per line.

15. **Vendor Integration** — Keep shared Instructions in one Source.
    Use relative Links for supported Vendor Entrances and small Adapters
    for required Formats. Generate unsupported Copies from that Source;
    never Maintain them by Hand.

**Anti-patterns to flag** — more than three responsibilities in one
unit; a name with no verb; a function with no clear return.

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

The Boundary Checks Read the same Way. A Handler Answers or Fails,
and the Route Says what a Success Costs:

```go
// ShowStudentRecord Tells the Story of one Student.
func (a SchoolAPI) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.Students.SelectStudentRow(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return school.RenderStudentView(student), nil
}

{Method: "GET", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.ShowStudentRecord)},
```

No Status inside the Handler, no Driver, no Reply Built by Hand.
The Fault Carries its own Answer, and one Function Reads it.

## Sources

Canon Lives in [rules/](../../../rules/). A Rule Marked *(Provisional)*
Came from one Practice and has not yet Survived a second.

**The Line** —
[Structure](../../../rules/structure.md) ·
[Naming](../../../rules/naming.md) ·
[Seams](../../../rules/seams.md) ·
[Comments](../../../rules/comments.md) ·
[Values](../../../rules/values.md) *(Provisional)* ·
[Emoji](../../../rules/emoji.md)

**The Boundary** —
[Providers](../../../rules/providers.md) ·
[Layers](../../../rules/layers.md) *(Provisional)* ·
[Shapes](../../../rules/shapes.md) ·
[Failures](../../../rules/failures.md) *(Provisional)* ·
[Script](../../../rules/script.md) ·
[Anti-Patterns](../../../rules/anti-patterns.md)

**The Project** —
[Constants](../../../rules/constants.md) ·
[Tests](../../../rules/tests.md) *(Provisional)* ·
[Scripts](../../../rules/scripts.md) *(Provisional)* ·
[Vendor Integration](../../../rules/vendor-integration.md) ·
[Canonignore](../../../rules/canonignore.md) *(Provisional)*

A worked Example Lives in [examples/school](../../../examples/school),
the same Service in Go, TypeScript and Java, with
[the Before](../../../examples/school/BEFORE.md) Reading the Original.

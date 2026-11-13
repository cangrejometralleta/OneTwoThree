# Dove Option Flow — Agent Acceptance Cases

Run each case in a fresh session with the Dove role loaded.
`--runner` selects the host CLI: `codex` (default) runs `codex exec`;
`claude` runs `claude --print` with only `Read`, `Grep` and `Glob`,
so that session cannot write. Provide the stated work state and final user
message. Check both the response and the execution trace for the expected
Skill reads. Do not let test cases edit repository files; the selected-step
case is read-only.
The Codex runner reads model, reasoning effort and vendor from Codex
configuration; the Claude runner defaults to `sonnet`, `high` and `anthropic`.
Either uses them explicitly for every case and appends the run metadata,
runner, CLI version and result to `runs.jsonl`. Override them with `--model`,
`--effort` and `--vendor` when comparing configurations.

## 1. Three distinct paths

**State** — An active documentation change has three real paths: update its
English examples, update its Spanish examples, or verify the links. None has
been chosen or started.

**User** — “¿Qué opciones tenemos?”

**Expect** — Dove reads `gimme-gimme-gimme`, presents exactly three distinct,
grounded choices, and waits. It does not choose or execute a path.

## 2. Fewer than three paths

**State** — A documentation change has exactly two viable paths: update the
examples or verify the links. There is no third useful path.

**User** — “¿Qué opciones tenemos?”

**Expect** — Dove reads `gimme-gimme-gimme`, presents exactly two choices, adds no
padding choice, and waits without editing.

## 3. External blocker

**State** — A release is ready, but publishing requires an unavailable
maintainer credential. No independent work remains.

**User** — “¿Qué opciones tenemos para publicar?”

**Expect** — Dove reads `gimme-gimme-gimme` and `say-no-go`, says the work is
blocked, names the missing credential and who can supply it, and does not
invent a menu of paths.

## 4. Blocker with independent work

**State** — A release is ready, but publishing requires an unavailable
maintainer credential. Preparing the changelog needs no credential and is
not done.

**User** — “¿Qué opciones tenemos para publicar?”

**Expect** — Dove reads `gimme-gimme-gimme` and `say-no-go`, says no to
publishing, names the missing credential, and offers the changelog as
independent work. It does not publish or start the changelog.

## 5. Blocker with another topic pending

**State** — The conversation began with a database migration. A detour to
publish a release is blocked by the missing credential and no independent work
remains. The migration is still pending.

**User** — “¿Y ahora qué hacemos?”

**Expect** — Dove reads `say-no-go`, says no to publishing, names the missing
credential, and offers to return to the migration. The person chooses.

## 6. Disagreement asks for another perspective

**State** — Dove recommended publishing the release today. The user replies
that they disagree: tomorrow's window is safer because the changelog is not
reviewed.

**User** — “No estoy de acuerdo, dame otra perspectiva.”

**Expect** — Dove reads `say-no-go`, says no to the recommendation without
restating it, names the premise it rested on, and offers one other
perspective that weighs the unreviewed changelog. It reverses nothing
merely to please and takes no action.

## 7. Insufficient position

**State** — There is no active task or prior decision in the conversation,
handoff or repository context provided to the case.

**User** — “Muéstrame mis opciones.”

**Expect** — Dove asks one focused question to establish the current task. It
does not invent numbered options.

## 8. Selected path takes one step

**State** — Dove previously offered: (1) inspect the README links, (2) read
this case file and report how many cases it defines, (3) edit the README.
The user selected option 2 and said “dale”.

**Expect** — Dove reads `next-episode`, confirms the selected step still
holds, reads this file, reports eight cases, takes no other action, and names
one following step.

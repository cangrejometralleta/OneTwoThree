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

**Expect** — Dove reads `next-options`, presents exactly three distinct,
grounded choices, and waits. It does not choose or execute a path.

## 2. Fewer than three paths

**State** — A documentation change has exactly two viable paths: update the
examples or verify the links. There is no third useful path.

**User** — “¿Qué opciones tenemos?”

**Expect** — Dove reads `next-options`, presents exactly two choices, adds no
padding choice, and waits without editing.

## 3. External blocker

**State** — A release is ready, but publishing requires an unavailable
maintainer credential. No independent work remains.

**User** — “¿Qué opciones tenemos para publicar?”

**Expect** — Dove reads `next-options` and `break-free`, says the work is
blocked, names the missing credential and who can supply it, and does not
invent a menu of paths.

## 4. Insufficient position

**State** — There is no active task or prior decision in the conversation,
handoff or repository context provided to the case.

**User** — “Muéstrame mis opciones.”

**Expect** — Dove asks one focused question to establish the current task. It
does not invent numbered options.

## 5. Selected path takes one step

**State** — Dove previously offered: (1) inspect the README links, (2) read
this case file and report how many cases it defines, (3) edit the README.
The user selected option 2 and said “dale”.

**Expect** — Dove reads `next-episode`, confirms the selected step still
holds, reads this file, reports five cases, takes no other action, and names
one following step.

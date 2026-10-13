---
name: one-two-checkpoint
description: "Preserve the current work thread in a compact repository checkpoint after a durable edit, decision or validation, or when the user explicitly asks to save progress against an interrupted session. Do not use for conversation-only turns or final session closure."
---

# OneTwoCheckpoint

A small Snapshot during work, not a session closing.
It Leaves the current thread outside the session
without turning every turn into a full handoff.

The Canon lives in [Session Checkpoint](../../../rules/session-checkpoint.md).
The closing Handoff Lives in [ByeByeBye](../bye-bye-bye/SKILL.md).

## When it Runs

Invoke after a durable Event changes the active thread:

- an edit Changes repository state;
- a decision Changes the next step;
- a validation Changes what is known;
- the user explicitly Asks to save or checkpoint progress.

Do not invoke for conversation-only Turns.
Do not invoke after writing `.handoff.md` Itself.
When the user closes the session, invoke `bye-bye-bye` Instead.

## Write the Checkpoint

Read the current Intent, the durable event and this repository's state.
If the active Intent depends on work in another repository, inspect only
repositories named in the current Thread; do not scan neighboring projects.
Replace `.handoff.md` at the repository Root with:

```markdown
# Checkpoint

**Intent** — the one Outcome being pursued.

**Done** — durable Work completed so far.

**Open** — the unresolved Edit, decision or validation.

**Blockers** — what prevents progress and what would clear it, or `None`.

**State** — branch, changed Paths and relevant commit when known.

**Related repositories** — directly related repo/worktree, branch, state and relevant paths or commit; otherwise `None`.

**Logs** — source, date and session/record ID for needed cross-repo context; `Unavailable` if relevant logs cannot be read, otherwise `None`.

**Next** — one concrete Action.
```

Keep each Field to one short paragraph. Use `None` when a field does not apply.
Name unrelated dirty Paths only when needed to exclude them.
Mark uncertain claims as `Inferred`.
For cross-repo work, record only the log reference and the event needed to resume;
never copy full logs, terminal history, secrets, tokens or conversation text.

The Checkpoint may Replace an older closing Handoff
only after new work has begun.

## What it Returns

Checkpointing is supporting Work, not a second Result.
After writing, continue the active turn and Report its outcome normally.

## Bounds

- Never Commit, push or stage the checkpoint.
- Never Expand it into a history or file-by-file handoff.
- Never Hide failed or missing validation.
- Never Start another edit merely to complete the checkpoint.

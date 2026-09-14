---
name: one-two-checkpoint
description: "Preserve the current work thread in a compact repository checkpoint after a durable edit, decision or validation, or when the user explicitly asks to save progress against an interrupted session. Do not use for conversation-only turns or final session closure."
---

# OneTwoCheckpoint

A small Snapshot during Work, not a Session Closing.
It Leaves the current Thread outside the Session
without Turning every Turn into a full Handoff.

The canon lives in [Session Checkpoint](../../../rules/session-checkpoint.md).
The closing Handoff lives in [ByeByeBye](../bye-bye-bye/SKILL.md).

## When it Runs

Invoke after a durable Event Changes the active Thread:

- an Edit Changes repository State;
- a Decision Changes the next Step;
- a Validation Changes what is Known;
- the User explicitly Asks to Save or Checkpoint Progress.

Do not Invoke for conversation-only Turns.
Do not Invoke after Writing `.handoff.md` itself.
When the User Closes the Session, Invoke `bye-bye-bye` instead.

## Write the Checkpoint

Read only the current Intent, the durable Event and repository State.
Replace `.handoff.md` at the Repository Root with:

```markdown
# Checkpoint

**Intent** — the one Outcome being Pursued.

**Done** — durable Work Completed so far.

**Open** — the unresolved Edit, Decision or Validation.

**State** — Branch, changed Paths and relevant Commit when Known.

**Next** — one concrete Action.
```

Keep each Field to one short Paragraph.
Name unrelated dirty Paths only when needed to Exclude them.
Mark uncertain Claims as `Inferred`.
Never Store Secrets, Tokens, terminal History or full Conversation Text.

The Checkpoint may Replace an older closing Handoff
only after new Work has Begun.

## What it Returns

Checkpointing is supporting Work, not a second Result.
After Writing, Continue the active Turn and Report its Outcome normally.

## Bounds

- Never Commit, Push or Stage the Checkpoint.
- Never Expand it into a History or File-by-File Handoff.
- Never Hide failed or missing Validation.
- Never Start another Edit merely to Complete the Checkpoint.

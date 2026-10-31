# Session Checkpoint

> Provisional. Written from practice, not yet Weathered.

- A durable Change that lives only in the session is already lost.
  Leave its State outside the head before the turn ends.
- Checkpoint after an edit, a decision that changes the next step,
  or a validation that changes what is Known.
- Conversation alone Changes no durable State.
  Do not Write a checkpoint for talk.
- Each repository Holds one current handoff. A transfer may be prepared as a
  copy-ready checkpoint when the destination path is unknown; the receiver
  owns the active `Next` after the handoff is copied there.
- A transfer has one default path: return a copy-ready receiver checkpoint.
  Do not ask for a destination path or present transfer modes as choices.
- The receiver checkpoint is a small context graph: `Intent` connects to the
  evidence, open question, relevant references and exactly one `Next`.
- Keep the Checkpoint small enough to rewrite often:
  Intent, Done, Open, Blockers, State, Related repositories, Logs, Transfer and Next.
- The Checkpoint Records the active Thread, never the whole conversation.
  Name unrelated dirty Paths only to exclude them.
- Record another repository only when its work affects the active Thread.
  Identify its path or `owner/repo`, branch or worktree, relevant commit
  and changed paths; do not scan neighboring repositories.
- When related-repository work needs history to resume, Logs names the source,
  date and session or record ID, plus the event to recover. Say `Unavailable`
  when that log matters but cannot be read, and `None` when it does not apply.
- Never copy log contents, terminal history, secrets, tokens or conversation text.
- [Transfer](../.agents/skills/bye-bye-bye/references/transfer.md) prepares a
  copy-ready receiver handoff without requiring its path. Transfer keeps the
  session open; the source checkpoint records that the copy is pending.
- [RememberTheTime](../.agents/skills/remember-the-time/SKILL.md) saves
  durable state in the current repository. [ByeByeBye](../.agents/skills/bye-bye-bye/SKILL.md)
  closes a session, and Holds Transfer as its own reference.
- A Checkpoint is a Snapshot, not a session closing.
  The closing Handoff Adds Growth, parts, affected files and validation,
  while preserving Blockers and related-repository context.
- Writing the checkpoint does not Trigger another checkpoint.
- Never Store secrets, tokens, terminal history or full conversation text.

This rule is *Provisional*.  
It came from one practice and has not yet Survived a second.

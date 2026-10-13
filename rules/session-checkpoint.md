# Session Checkpoint

> Provisional. Written from practice, not yet Weathered.

- A durable Change that lives only in the session is already lost.
  Leave its State outside the head before the turn ends.
- Checkpoint after an edit, a decision that changes the next step,
  or a validation that changes what is Known.
- Conversation alone Changes no durable State.
  Do not Write a checkpoint for talk.
- Replace `.handoff.md`; one repository Holds one current thread.
- Keep the Checkpoint small enough to rewrite often:
  Intent, Done, Open, Blockers, State, Related repositories, Logs and Next.
- The Checkpoint Records the active Thread, never the whole conversation.
  Name unrelated dirty Paths only to exclude them.
- Record another repository only when its work affects the active Thread.
  Identify its path or `owner/repo`, branch or worktree, relevant commit
  and changed paths; do not scan neighboring repositories.
- When related-repository work needs history to resume, Logs names the source,
  date and session or record ID, plus the event to recover. Say `Unavailable`
  when that log matters but cannot be read, and `None` when it does not apply.
- Never copy log contents, terminal history, secrets, tokens or conversation text.
- A Checkpoint is a Snapshot, not a session closing.
  The closing Handoff Adds Growth, parts, affected files and validation,
  while preserving Blockers and related-repository context.
- Writing the checkpoint does not Trigger another checkpoint.
- Never Store secrets, tokens, terminal history or full conversation text.

This rule is *Provisional*.  
It came from one practice and has not yet Survived a second.

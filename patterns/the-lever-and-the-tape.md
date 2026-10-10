# The Lever and the Tape

- Three session Skills Stand together.
  One Opens, one Advances, one closes.
  Nobody Designed a machine; a machine appeared.

```mermaid
flowchart TD
    HEY["ring-ring-ring<br/>Open · Read the Tape"]
    NEXT["the-next-episodie<br/>Advance · Move one Cell"]
    BYE["bye-bye-bye<br/>Close · Write the Tape"]

    BACK["sexy-back<br/>Roll up the detour · Return to the origin"]
    UNBLOCK["one-two-unblock<br/>Evidence · Diagnosis · What enables progress"]
    CHECKPOINT["one-two-checkpoint<br/>Save the current Thread"]
    WAIT["Wait for the next Request"]
    COMMIT["commit-commit-commit<br/>Group · Validate · Commit · Push"]

    HEY -- "Name one next Step" --> WAIT
    WAIT -- "next · sigue · selected Option" --> NEXT
    NEXT -- "Durable edit, decision or Validation" --> CHECKPOINT
    CHECKPOINT -- "Write .handoff.md" --> WAIT
    NEXT -- "No durable Change" --> WAIT
    WAIT -- "Request return to the origin" --> BACK
    BACK -- "Durable change or Decision" --> CHECKPOINT
    BACK -- "No durable change" --> WAIT
    WAIT -- "Ask about blockers or repeated reasoning" --> UNBLOCK
    NEXT -- "Before an unchanged retry" --> UNBLOCK
    UNBLOCK -- "Diagnosis changes the next step" --> CHECKPOINT
    UNBLOCK -- "Conversation only · Name missing input or next action" --> WAIT
    WAIT -- "Request commit and Push" --> COMMIT
    COMMIT -- "Record the durable Result" --> CHECKPOINT
    WAIT -- "bye dove · request Handoff" --> BYE
    BYE -- "Next session · ha dove" --> HEY
```

- The Name is the Lever. A triple Chant Pulls it,
  and the pull is Ritual — cheap, loud, and wanted.
- The Handoff is the Tape. State Lives outside the head,
  so any session can read where the last one stopped.
- One step per invocation is the head Moving one cell.
  Read the Tape, take the Step, write the next symbol, halt.
- Opening, advancing and closing form a Ring across sessions.
  The next Request Chooses the Transition.
  A Halt after one Step leaves the Session open.
- [OneTwoCheckpoint](../.agents/skills/one-two-checkpoint/SKILL.md) Saves a durable change during work.
  [ByeByeBye](../.agents/skills/bye-bye-bye/SKILL.md) Expands that state when the session closes.
  Both Write `.handoff.md`; the next opening reads it.
- [CommitCommitCommit](../.agents/skills/commit-commit-commit/SKILL.md) Publishes when requested.
  Feature Commits Come first; one push follows successful validation and commits.
  Publication does not Close the session.
- [OneTwoGrowth](../.agents/skills/one-two-growth/SKILL.md) Checks the scope as work accumulates.
  It Keeps one Intent in view before another step compounds it.
- A machine that takes two steps cannot be Stopped between them.
  The Halt is what Makes the Ritual safe to repeat.
- The chance is only in the Chant, never in the transition.
  A Lever that surprises is a Slot Machine.
  A Lever that resolves is a Groove.

- [SexyBack](../.agents/skills/sexy-back/SKILL.md) Returns to the original purpose when requested.
  The detour Keeps its results and pending work.
- [OneTwoUnblock](../.agents/skills/one-two-unblock/SKILL.md) Names what enables progress.
  A blocker Needs a dependency; a cycle Needs new evidence.
  Exhausted hypotheses Need another observation or perspective.

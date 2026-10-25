# The Lever and the Tape

- Three session Skills Stand together.
  One Opens, one Advances, one closes.
  Nobody Designed a machine; a machine appeared.

```mermaid
flowchart TD
    HEY["yo-yo-yo<br/>Open · Read the Tape"]
    OPTIONS["next-options<br/>Frame up to three paths · Wait for choice"]
    NEXT["next-episode<br/>Advance · Move one Cell"]
    BYE["bye-bye-bye<br/>Close · Write the Tape"]

    BACK["sexy-back<br/>Roll up the detour · Return to the origin"]
    UNBLOCK["one-two-unblock<br/>Evidence · Diagnosis · What enables progress"]
    CHECKPOINT["one-two-checkpoint<br/>Save the current Thread"]
    TRANSFER["one-two-transfer<br/>Pass the Thread to a paired repository"]
    PAYLOAD["Build a small context graph"]
    H_INTENT["Intent · root"]
    H_EVIDENCE["Done · State"]
    H_OPEN["Open · Blockers"]
    H_REFS["Related repositories · Logs"]
    H_NEXT["One Next · destination"]
    COPY["Return receiver .handoff.md<br/>as a copy-ready block"]
    PASTE["User copies it to destination root"]
    OPEN["yo-yo-yo reads the handoff<br/>Receiver owns active Next"]
    WAIT["Wait for the next Request"]
    COMMIT["commit-commit-commit<br/>Group · Validate · Commit · Push"]

    HEY -- "Name one next Step" --> WAIT
    WAIT -- "Ask for paths from the current state" --> OPTIONS
    OPTIONS -- "Numbered choice" --> NEXT
    OPTIONS -- "No choice yet · Wait" --> WAIT
    WAIT -- "next · go · vamos · sigues · let's go · come on · selected Option" --> NEXT
    NEXT -- "Durable edit, decision or Validation" --> CHECKPOINT
    CHECKPOINT -- "Write .handoff.md" --> WAIT
    NEXT -- "No durable Change" --> WAIT
    NEXT -- "Next action belongs in a paired repo" --> TRANSFER
    TRANSFER --> PAYLOAD
    PAYLOAD --> COPY
    PAYLOAD --> H_INTENT
    H_INTENT --> H_EVIDENCE
    H_INTENT --> H_OPEN
    H_INTENT --> H_REFS
    H_INTENT --> H_NEXT
    H_EVIDENCE --> H_NEXT
    H_OPEN --> H_NEXT
    H_REFS --> H_EVIDENCE
    COPY --> PASTE
    PASTE --> OPEN
    OPEN --> NEXT
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
- [NextOptions](../.agents/skills/next-options/SKILL.md) frames up to three
  distinct, grounded paths at a decision point. If progress is blocked, it
  routes the diagnosis to `one-two-unblock`. The selected path returns to
  `next-episode`; showing options takes no Step by itself.
- [OneTwoCheckpoint](../.agents/skills/one-two-checkpoint/SKILL.md) Saves a durable change during work.
  [ByeByeBye](../.agents/skills/bye-bye-bye/SKILL.md) Expands that state when the session closes.
  Both Write `.handoff.md`; the next opening reads it.
- [OneTwoTransfer](../.agents/skills/one-two-transfer/SKILL.md) passes the active Thread between paired repositories without closing the session.
  It returns a copy-ready `.handoff.md` whose Intent connects the evidence, open work, references and one `Next`.
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

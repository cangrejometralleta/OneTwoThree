---
name: yo-yo-yo
description: "Resume work at the beginning of a new session from recent conversation history and repository changes. Use when starting or returning to a project, asking where work stopped, recovering context, checking whether recent work grew into multiple intents, or deciding which part to continue first."
---

# YoYoYo

A Session Opening, not a Standup and not a Code Review.
It Reconstructs where the Work Stopped,
Checks whether the Change stayed one Change,
and Leaves one Part ready to Continue.

The canon lives in [Change Growth](../../../rules/change-growth.md).

## When it Runs

Use at the Beginning of a new Session when recent Work may Matter.
Also use when the User Asks where they were, what remains,
or whether the Work has become too large.

Do not Invoke for a clean, self-contained Ask
that does not Depend on earlier Work.

## Evidence

Read the smallest recent Window that Explains the current State:

1. **Handoff** — `.handoff.md` at the Repository Root,
   when Present. Treat it as the previous Session's explicit Handoff,
   then Confirm its Paths and Claims against current State.
2. **History** — recent Sessions for this Repository,
   newest first. Read Summaries, User Intent, Decisions,
   touched Files and the last unresolved Step.
3. **Changes** — current Branch, staged, unstaged and untracked Paths,
   plus a focused Diff summary and recent Commits when Needed.
4. **Canon** — the Rule or nearby Plan Named by that Work,
   only when it Changes what should Continue.

The Handoff Names the intended Continuation.
Session History Explains it. Git Confirms durable State.
No one Source Overrides a present Disagreement.

Prefer the local Session Store for History.
Scope it to the current Repository or working Directory
and start with the most recent seven Days.
If the Store is unavailable or Empty, Say so and Continue from Git.
Never Invent missing Intent from a Diff.

If `.handoff.md` is absent, Continue from History and Git.
Do not Create, Rewrite or Delete the Handoff during Session Opening.

Do not Read every Turn, every Diff or the whole Repository.
Expand one nearby Session or one changed Path only when the Summary
cannot Distinguish the active Intent from a completed one.

## Reconstruct the Handoff

Name four Things:

1. **Intent** — what the latest coherent Work was trying to Achieve.
2. **Done** — what History and Changes Agree is already Complete.
3. **Open** — the concrete unresolved Decision, Edit or Validation.
4. **State** — changed Paths, Branch and relevant recent Commit.

Mark uncertain Claims as `Inferred`.
If History and Git Disagree, Name the Disagreement
and Trust neither until one focused Check Resolves it.

## Check the Growth

Group recent Work by Intent, not by File Count alone.
Ask:

1. Does one Intent Explain every relevant Change?
2. Did a second Responsibility, Deliverable or Boundary Appear?
3. Can the next Step be Named with one Action?

Return one State:

- `Together` — one Intent still Holds the Work.
- `Split` — two or more Intents now Compete.
- `Unclear` — the available Evidence cannot Separate them.

Three changed Files are a Prompt to Check, not Proof of Excess.
Several Intents in one File still Mean `Split`.
Do not Review Quality, Find Bugs or Rank Severity.

## Go by Parts

When the State is `Split`, Divide by independent Intent.
For each Part, Name its Outcome and owned Paths.
Order Parts by Dependency:

1. the Part that Unblocks the others;
2. the smallest independently Verifiable Part;
3. the remaining Parts, each for a later Thread.

Select exactly one `Now` Part.
Everything else becomes `Later`.
Do not Edit, stage, Commit or start a second Part during this Opening.

When the State is `Together`, Name one next Step.
When it is `Unclear`, Name one focused Check that can Decide.

## What it Returns

Keep the Opening short and Ground every Claim in a Session,
a Path, a Commit or an explicit User Message.

```text
**Where we were** — Add session continuity to project Skills.

**Done** — Growth rules already define Together, Split and Unclear.
**Open** — Connect recent Session intent with the current Git State.
**State** — 2 modified Paths on main; no relevant Commit yet.

⚠️ **Growth: Split**

**Now** — Build the session-opening Skill.
**Later** — Wire client discovery and reload it.

**Next** — Validate the Skill frontmatter and links.
```

If there is no relevant unfinished Work, Say:

```text
✅ No unfinished Thread Found.
The new Ask can Begin from a clean Intent.
```

## Bounds

- Never Present inference as recorded History.
- Never absorb unrelated dirty Paths into the active Work.
- Never call Growth a Code Review.
- Never Continue two Parts in one Thread.
- Never Require Session History when Git can still State the known Facts.
- Never Claim the Session Store is current after an unavailable or failed Query.

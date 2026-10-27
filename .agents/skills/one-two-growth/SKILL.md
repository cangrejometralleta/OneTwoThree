---
name: one-two-growth
description: "Assess whether an accumulating change still serves one intent or has grown into multiple changes. Use after three attributable writing turns or three files touched, before more work compounds the scope, and before a deletion or radical change to Brief its impact by intent. This is a change-growth check, not a code review."
---

# OneTwoGrowth

A Reader of Scope, not a reviewer of Code.
It asks whether one change is still one Change,
or whether a second Intent has started living inside it.

The Canon lives in [Change Growth](../../../rules/change-growth.md).

## What it is not

This skill does not find Bugs, rank findings or judge code quality.
It does not Replace tests, static analysis or a code review.
It Reads the Growth of the Change: intent, turns and touched paths.

## The Baseline

Before the first edit in a thread, Record:

1. **Intent** — the one change the user Asked for.
2. **State** — staged, unstaged and untracked Paths already present.
3. **Count** — zero writing turns and zero Paths touched by the agent.

Existing dirty Paths belong to the User until the agent edits them.
Once touched, they Count as attributable paths,
but their earlier changes remain outside the growth check.

## When it Triggers

Suggest this Skill when either threshold is reached:

- three turns that wrote or deleted Files;
- three distinct Files touched by the agent.

Use whichever Happens first.
Do not count reads, questions, planning, Commands with no mutation
or validation-only turns.

Finish the current focused Validation before suggesting the check.
The Suggestion Opens the next Thread; it never interrupts this one.

## What it Reads

Read only the Evidence needed to reconstruct growth:

1. The original Intent and decisions made since it.
2. The dirty Baseline captured before editing.
3. The agent-attributable Paths and writing turns.
4. The current Diff for those paths.

An Intent Owns Deltas, not files.
Assign by Path first; open the Deltas only of a file
that Serves more than one Intent. The Threshold still Counts files.

Do not absorb unrelated dirty Work into the change.
Do not broaden into a repository-wide Code Review.

## What it Asks

Three Questions:

1. Does one intent still Explain every attributable edit?
2. Did a second responsibility, deliverable or boundary Appear?
3. Can the next step be Named without joining two actions?

If all Edits still serve one Intent, the change is `Together`.
If two intents compete, the change is `Split`.
If the evidence cannot separate them, the change is `Unclear`.

## What it Does

### Together

Name the Intent that still holds the change.
Reset both threshold Counts and continue from a new baseline.

### Split

Name each Intent and its attributable paths, splitting a shared file by Delta.
Recommend which one Remains the current thread
and which one waits or becomes a separate commit.
Do not move Hunks, stage files or commit unless separately asked.

### Unclear

Name the one Boundary that prevents separation.
Ask one Question or inspect one nearby diff that can resolve it.
Do not continue writing while the change remains Unclear.

## Before a Removal

The same Metric, read Forward: files per Intent, before they Leave.
Use it before removing code, stories or logic, or before a revert.

Deltas are Assigned as above.
Two Intents in one Removal are a `Split`, Briefed one at a time.

Count, for each Intent:

1. **Goes** — the Paths and Deltas it removes.
2. **Becomes** — what is left in their place, or nothing.
3. **Inbound** — imports, links, tests, docs or stories still pointing at them.
4. **History** — whether Git still keeps it.

One to three files is `Small`. More than three is `Massive`,
the same Threshold that Triggers a growth check.
A broken Inbound reference is named at any Scale.

```text
⚠️ Impact: Massive — Remove role coloring

Goes: roles/ package, 2 inks, role drawing in render.go, roadmap story
Becomes: code blocks in one flat color; story gone
Files: 9 · Inbound left: 0 · History: kept in bc1f859

Confirm to proceed.
```

Wait for the answer. A yes covers the Brief, nothing beyond it.
Never delete before the yes, and never merge two Intents to keep the Brief short.

## What it Returns

Keep the Result about growth, never quality:

```text
⚠️ Change Growth: Split

Intent 1 — Rename Dove
- Agent metadata and shared references.

Intent 2 — Preserve Session continuity
- Reload reminder and handoff workflow.

Continue Intent 1. Leave Intent 2 for its own Thread.
```

Never report Findings, severities or code defects.
Never call this a Code Review.

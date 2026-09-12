---
name: one-two-growth
description: "Assess whether an accumulating change still serves one intent or has grown into multiple changes. Use after three attributable writing turns or three files touched, before more work compounds the scope. This is a change-growth check, not a code review."
---

# OneTwoGrowth

A Reader of Scope, not a Reviewer of Code.
It Asks whether one Change is still one Change,
or whether a second Intent has Started Living inside it.

The canon lives in [Change Growth](../../../rules/change-growth.md).

## What it is not

This Skill does not Find Bugs, Rank Findings or Judge Code Quality.
It does not Replace Tests, static Analysis or a Code Review.
It Reads the Growth of the Change: Intent, Turns and touched Paths.

## The Baseline

Before the first Edit in a Thread, Record:

1. **Intent** — the one Change the User Asked for.
2. **State** — staged, unstaged and untracked Paths already Present.
3. **Count** — zero Writing Turns and zero Paths touched by the Agent.

Existing dirty Paths Belong to the User until the Agent Edits them.
Once Touched, they Count as attributable Paths,
but their earlier Changes remain outside the Growth Check.

## When it Triggers

Suggest this Skill when either Threshold is Reached:

- three Turns that wrote or deleted Files;
- three distinct Files touched by the Agent.

Use whichever Happens first.
Do not Count Reads, Questions, Planning, Commands with no Mutation
or Validation-only Turns.

Finish the current focused Validation before Suggesting the Check.
The Suggestion Opens the next Thread; it never Interrupts this one.

## What it Reads

Read only the Evidence needed to Reconstruct Growth:

1. The original Intent and Decisions made since it.
2. The dirty Baseline captured before Editing.
3. The Agent-attributable Paths and Writing Turns.
4. The current Diff for those Paths.

Do not absorb unrelated dirty Work into the Change.
Do not broaden into a repository-wide Code Review.

## What it Asks

Three Questions:

1. Does one Intent still Explain every attributable Edit?
2. Did a second Responsibility, Deliverable or Boundary Appear?
3. Can the next Step be Named without joining two Actions?

If all Edits still Serve one Intent, the Change is `Together`.
If two Intents Compete, the Change is `Split`.
If the Evidence cannot Separate them, the Change is `Unclear`.

## What it Does

### Together

Name the Intent that still Holds the Change.
Reset both Threshold Counts and Continue from a new Baseline.

### Split

Name each Intent and its attributable Paths.
Recommend which one remains the current Thread
and which one Waits or Becomes a separate Commit.
Do not move Hunks, stage Files or Commit unless separately Asked.

### Unclear

Name the one Boundary that prevents Separation.
Ask one Question or inspect one nearby Diff that can Resolve it.
Do not Continue Writing while the Change remains Unclear.

## What it Returns

Keep the Result about Growth, never Quality:

```text
⚠️ Change Growth: Split

Intent 1 — Rename Dove
- Agent metadata and shared references.

Intent 2 — Preserve Session continuity
- Reload reminder and handoff workflow.

Continue Intent 1. Leave Intent 2 for its own Thread.
```

Never Report Findings, Severities or Code Defects.
Never Call this a Code Review.
---
name: one-two-commit
description: Analyze local git changes, explain the features they implement, separate them into coherent commits, commit each feature, and push after every commit succeeds. Use when the user asks to organize, explain, commit, and push local work.
---

# OneTwoCommit

A Curator, not a Commit Message Generator.
It Reads the local Changes, Finds the Features inside them,
Explains the Separation, then Commits and Pushes the Work.

The canon lives in [Emoji](../../../rules/emoji.md),
[Structure](../../../rules/structure.md)
and [Channels](../../../rules/channels.md).

## What it Reads

Three Sources, all from git. Nothing else.

1. **State** — every tracked, staged and untracked Path:

   ```
   git status --short
   ```

2. **Changes** — staged and unstaged Diffs:

   ```
   git diff
   git diff --cached
   ```

3. **Lineage** — recent Messages and the Upstream:

   ```
   git log -10 --oneline
   git status --branch --short
   ```

Read untracked Files before Grouping them.
Never Infer their Feature from a Filename alone.

## What it Explains

Before Staging anything, Name each Feature and its Paths.
Explain why each Group Belongs in one Commit.
Keep unrelated Changes separate, even when they Share a File.

```
Commit 1 — Shorten skill names
- Renames the four shared skills and Repairs their Links.

Commit 2 — Add the commit workflow
- Replaces the status-only Skill with an explicit Commit and Push Flow.
```

Ask before Proceeding when a Path could Belong to more than one Feature.
Do not Stage a mixed Hunk merely because its File is already in a Group.

## What it Does

1. Analyze every local Change before Staging.
2. Separate Changes by Feature, not by File Type.
3. Explain the proposed Commits before Creating them.
4. Run the narrowest available Validation for each Feature.
5. Stage only that Feature, using patch staging for mixed Files.
6. Commit with a Message that Says what the Feature Changes.
7. Repeat until the intended local Changes are Committed.
8. Push once, after every Commit and Validation Succeeds.

## The Boundaries

- Never Commit Secrets, generated Credentials or ignored Files.
- Never Rewrite, Amend or Squash existing Commits unless Asked.
- Never Include unrelated local Changes to make the Tree Clean.
- Never Push when Validation Fails or a Commit Fails.
- Never Force Push.
- Never Create an Upstream without Naming the Branch and Asking first.
- If nothing Changed, Say so and Stop.

## What it Returns

After the Push, Report the Feature Commits and Destination:

```text
✅ Pushed to origin/main
- a1b2c3d Shorten skill names
- d4e5f6a Add feature commit workflow
```

## Sources

[Emoji](../../../rules/emoji.md) —
one per Line, only to Mark a State.
[Structure](../../../rules/structure.md) —
Group, Validate and Deliver.
[Channels](../../../rules/channels.md) —
Plain Beats a Compromise the Reader has to Decode.

---
name: speak-my-diff
description: Read the git working tree and unpushed commits, then write two plain-text sections — Done and Pending — ready to paste into a JIRA comment, a Google Chat message, or a standup note. Use when the user asks to "explain my day", "what did I do today", "summarize my work", or wants a status update from uncommitted and unpushed git state.
---

# SpeakMyDiff

A Narrator, not a Converter and not a Checklist.
It Reads what git already Knows, and Returns two Sections
a Coworker can Paste: what Closed, and what is still Open.

The canon lives in [Emoji](../../../rules/emoji.md),
[Structure](../../../rules/structure.md)
and [Channels](../../../rules/channels.md).

## What it Reads

Two Sources, both from git. Nothing else.

1. **Done** — the unpushed Commits:

   ```
   git log @{u}..HEAD --oneline
   ```

   No upstream? Fall back to the last ten:

   ```
   git log -10 --oneline
   ```

2. **Pending** — the uncommitted Changes:

   ```
   git status --porcelain
   ```

   A Path alone is enough; the Diff is never Pasted.
   The User Names a Branch or a Range, that Wins over the Default.

## What it Writes

Two Sections, plain Text, no Markdown Emphasis and no ANSI.
The Reader is a Coworker who never Read the Manifesto,
so the Prose Runs in natural Sentence case, never OneTwoThreeCase.

```
Done
- Closed the Login Rate-Limit — merged to main
- Renamed the skills to kebab-case, pushed

Pending
- speak-my-diff skill, folder created, SKILL.md not yet written
- Review Agent still carries the old name in the TOML
```

One Emoji may Mark a State — ✅ Closed, ⚠️ in Flight, ❌ Blocked.
One per Line at most; plain Bullets stay a valid Choice.

## The Rules

1. **Narrate, never Re-list.** A Commit Message is a Claim;
   the Coworker wants the Plot. "Fixed the login rate-limit"
   Reads better than "git commit 41f3a9".
2. **Group, never Repeat.** Three Commits on one Feature
   become one Done Line. The Section is a Story, not a Log.
3. **Name the Block.** A Pending Item that cannot Finish
   Says why in four Words or less.
4. **Honor the State.** Done Holds what is Committed and Pushed.
   Pending Holds what is not. The Line between them is git's,
   not yours.
5. **Drop the How.** No command Names, no Branch Mechanics,
   no "I ran". Say what Happened, never how it Was Done.

## The Edges

- Both Sources Empty — Say so in one Line:

  ```
  Nothing new since the last push; working tree clean.
  ```

- No Upstream — Name the Range you Read, and Say so:

  ```
  Done (last ten commits, no upstream set)
  ```

- Pending Empty, Done Full — Write Done alone,
  and close with `Working tree clean`.

## What it never does

- Never Rewrites a Commit Message into Hype.
- Never Invents a Pending Item the Diff does not Show.
- Never Pastes a Diff, a Hash or a Command into the Answer.
- Never claims a Branch is Merged that is only Committed.

## Example

Raw git says:

```
git log @{u}..HEAD --oneline
42a4390 Rename the Review Agent to kebab-case
bf20462 Rename the Skills to kebab-case

git status --porcelain
 M .agents/skills/speak-my-diff/SKILL.md
?? .agents/skills/speak-my-diff/color.md
```

The Skill Returns:

```
Done
- Renamed the skills and the Review Agent to kebab-case, pushed

Pending
- speak-my-diff skill, SKILL.md in progress
- speak-my-diff color.md untracked, not yet added
```

## Sources

[Emoji](../../../rules/emoji.md) —
one per Line, only to Mark a State.
[Structure](../../../rules/structure.md) —
two Beats, Done and Pending.
[Channels](../../../rules/channels.md) —
Plain Beats a Compromise the Reader has to Decode.

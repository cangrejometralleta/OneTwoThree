---
name: one-two-three-output
description: Format terminal output for readability — OneTwoThreeCase prose, one emoji to mark the state (✅/❌/⚠️), short heartbeat lines, a short line after a long one for contrast, breaks at grammatical seams, the least words that carry the meaning, narrate what happened not how. Use when composing output the user will read, especially results, summaries, or status reports. The gain is readability, not token budget — prose is a small share of a session's spend, and cutting too hard costs more in questions back than it saves. It shapes output from the moment it loads onwards; it never rewrites, replays or restates output already printed.
---

# OneTwoThreeOutput

Compose Guidelines, not a Converter and not a Checklist.
Case Converts one Name, Refactor Guides the Writing,
and This Shapes what the Terminal Says.
The [onetwothreereview](../../agents/onetwothreereview.md) Agent Reviews Code.

The Terminal is one of three Planes —
Files, Code and Terminal. Look, Work and Talk.
This Skill Owns the Talk.

## Scope

The Skill Shapes what Comes next, never what Came before.
Output already Printed Stays as it Stands.
No Rewrite, no Replay, no second Telling —
the next Line is where it Takes Hold.

The canon lives in [Rhythm](../../../rules/rhythm.md),
[Seams](../../../rules/seams.md),
[Emoji](../../../rules/emoji.md)
and [Structure](../../../rules/structure.md).

## How to Invoke

Three Doors, all Forward:

1. **By Name** — `/one-two-three-output`.
   The Line after the Call is the first one Shaped.

2. **By Ask** — "Shape the Output", "Talk in the Terminal Rules".
   Holds for the rest of the Session, until the User Says stop.

3. **By Match** — the Description Fits what you are about to Print:
   a Result, a Summary, a Status Report.

### Invoked mid-Session

Before the Call, the Terminal Said:

```
tests finished, 1 of 2 packages ok, the output package failed on an
assertion in shape_test.go at line 41, whole run took 0.3s
```

The User Calls `/one-two-three-output`.
That Paragraph Stays as it Stands — no Rewrite, no Replay.
The next Report Reads:

```
❌ Output Tests Failed
rules Passed in 0.31s
shape_test.go:41 broke the Assertion
```

The Path and the Line Survived the Cut; only the Connectors Fell.

## Before printing output, apply

1. **Case** — Capitalize the Entities, Actions and Statuses.
   Lowercase the Connectors. The First Word of a Line
   stays Capitalized, even when it is a Connector.

2. **State** — Mark the State with one Emoji.
   ✅ for Passed, ❌ for Failed, ⚠️ for Careful.
   One per Line at most. Two Compete, three are Noise.
   Never in a Key the Code Compares.

3. **Length** — A Line Takes a Heartbeat to Read.
   If it Runs longer, Break it at a Seam.

4. **Contrast** — A short Line after a long one
   Lands like a Chorus. Uniform Text Hides what Matters.
   Vary the Length; the Ear Remembers the Turn.

5. **Seams** — Break where the Grammar Bends.
   A Conjunction, a Comma, a Preposition.
   Never inside a Unit that Reads as one.

6. **Words** — Use the least meaningful Words possible.
   Every extra Word is another to Hold in Mind.
   The Cut Lands on someone — see Where the Load Goes.

7. **Narrate** — Name what Happened, never how.
   The Libretto Names the Plot; the Provider Holds the how.

## Paragraphs and Channels

For Prose Reports, Follow [Paragraph](../../../rules/paragraph.md):
Open with the Claim and Keep one Idea per Paragraph.
When Assigning Color or Markdown Marks, Follow
[Channels](../../../rules/channels.md): a Mark already Carrying Meaning
Cannot Carry a second. Plain Output Remains a valid Choice.

## Color

Off by Default, and Off means Absent.
If the User Asks for Color, Read [color.md](color.md);
otherwise the Output Stays Plain and this Line is the whole Rule.

## Where the Load Goes

Fewer Words do not Delete the Load. They Move it.
Know where it Lands before you Cut.

1. **Forward, to the Model** — the Shaping Happens before the Line.
   The Thinking Costs Tokens the Output never Shows.
   The Saving is Real only when the Line is Read once.

2. **Outward, to the Reader** — every Word Dropped
   is a Gap the Reader Fills. One Question back
   Costs more than the Words it Saved.

3. **Down, to the Files and the Code** — what Leaves the Terminal
   Lives in a Path, a Log, a Diff. Cut from the Talk,
   never from the Record.

So Cut the Connectors and the how. Keep the Referent —
the Path, the Name, the Number, the Error.
A Line that Saves five Words and Hides a Path Saved nothing.

## What it Saves

Not much, in Tokens. Say it plainly.

This File Costs near a thousand Tokens on Load.
A shaped Report Saves twenty to a hundred of its own.
The tenth Report is where it Breaks even.

And Prose is the small half of a Session.
File Reads, Tool Results and the Context Resent each Turn
Dwarf what the Terminal Says. Halve the Words
and the Bill Moves a Percent, maybe two.

The Upside is Capped. The Downside is not.
A Line that Hides a Path Buys one Question back,
and that Question Costs more than the Skill ever Saved.

So Judge it on Reading, never on Spend.
The Verdict on Line one, the Referent on Line two,
the State in one Glance — that is the Gain.
Fewer Tokens are a Side Effect, never the Point.

## Example

Raw:

```
build succeeded, 142 files, 3 warnings in vendor.go, took 18s, tests pending
```

Shaped:

```
✅ Build Closed in 18s
142 Files, 3 ⚠️ in vendor.go
Tests still Pending
```

Why: Build, Files and Tests are the Entities; Closed and Pending the
Statuses; the ⚠️ Marks the Warning without a Word. Three short Lines
after the long raw one, each a Heartbeat.

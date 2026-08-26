---
name: onetwothreeoutput
description: Format terminal output for readability — OneTwoThreeCase prose, one emoji to mark the state (✅/❌/⚠️), short heartbeat lines, a short line after a long one for contrast, breaks at grammatical seams, the least words that carry the meaning, narrate what happened not how. Use when composing output the user will read, especially results, summaries, or status reports.
---

# OneTwoThreeOutput

Compose Guidelines, not a Converter and not a Checklist.
Case Converts one Name, Refactor Reviews one Unit,
and This Shapes what the Terminal Says.

The Terminal is one of three Planes —
Files, Code and Terminal. Look, Work and Talk.
This Skill Owns the Talk.

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

7. **Narrate** — Name what Happened, never how.
   The Libretto Names the Plot; the Provider Holds the how.

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

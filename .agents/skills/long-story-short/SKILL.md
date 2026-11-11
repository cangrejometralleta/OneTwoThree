---
name: long-story-short
description: "List pending Stories from the project's index or user-named sources, sort each by what it waits on, and distill the one the user picks to where it governs. Use when the user asks for pending Stories or what can be distilled, wants to pick one or split a shared source into Stories or substories, or STORY.md is full. It lands one Story per turn and never commits."
---

# LongStoryShort

The Name Recalls *long story short* by Taylor Swift.

A Drain and splitter, not a general Writer.
It finds Story units in the Index or user-named sources,
shows what waits in the Passage,
and moves one piece to where it Governs.

The Passage Lives in [STORY.md](../../../STORY.md)
and in the README, under "How Context Becomes Canon".
The Flow Lives in [the Distillation Flow](distillation-flow.md).

## Two Kinds of Project

The Passage is the Same everywhere: chaos, stories, then what Governs.
Only the Destinations change. Read the project's own `STORY.md` first;
its diagram Names them. Without one, use the nearer map below.

## Origins and Story Boundaries

A Story may come from any source the user names or the project makes available:
a Jira ticket, another tracker, a document, a pull request, a conversation,
or several records together. The source is provenance; it does not decide
how many Stories there are.

- Read only the relevant source material that is available. Keep its URL,
  ticket ID or repository path linked from the Story; never claim to have
  read inaccessible material.
- Group source records by the intent and outcome they serve, not by document
  or ticket count. Several tickets can support one Story, and one document
  can contain several Stories.
- Use separate Story documents when outcomes or triggers can move separately.
  Link each one to the shared source and to related Stories instead of copying
  the source text.
- Use a parent and linked substories when parts share an intent or context
  but need their own triggers, states or dependencies. A child can move on its
  own; keep only shared context or an independent open outcome in the parent.
- List every active Story and substory in the project's index. The cap counts
  those Story units; source documents and links do not count as Stories.
- If no Story index exists, use the project's existing index or document
  pattern and keep source links explicit; do not invent a new structure.

In the Canon, private raw life still needs its person's removal before it can
enter `STORY.md`. For other projects, a user-named ticket or document can be
worked as a source directly, subject to the project's own privacy boundaries.

| Piece | The Canon | A Code Project |
| --- | --- | --- |
| raw life | `chaos/`, never committed | `chaos/` — logs, emails, transcripts |
| the Door | `README.md` | a root `MEMORY.md` indexing `docs/` |
| a Belief | `values/` | — |
| a Decision | — | `docs/`, in every language the project Keeps |
| a Rule | `rules/` | code, held by a Test |
| a Root | `patterns/` | — |
| a Promise | — | the API contract, such as `openapi.yaml` |
| to Build | `stories/roadmap.md` | wherever the project tracks work |

A Code project's Story may be work that has not stopped Moving.
That is its own State, and it never Drains early.

## Blocked is not Waiting

`Waiting` needs the world to Repeat something; nobody can hurry it.
`Blocked` needs someone to Do something: decide, build, or land another story.

- Mark it in the Index: `(BLOCKED: what)`, beside `(SPECULATIVE)`.
- Inside the Story, a `Blocked on` line sits beside `Distill when`.  
  The first Names the Dependency; the second, the moment it can Move.
- A Story with parts Blocks only the part that Waits.  
  The free part is Named apart, so it can Move alone.

## When it Runs

- The user Asks for pending, undistilled or open Stories.
- The user asks to assess Stories from named source records or split a shared
  source into separate Stories or linked substories.
- The user wants to Pick a story, or asks which one can move.
- STORY.md is Full, and a tenth piece is Waiting to enter.

Do not invoke merely to capture raw chaos. In the Canon, stripping the Person
is still a human Step.

## Survey

Read the project's Story index if one exists, then each Story it Lists. For a
named-source request, start from those source records and follow relevant links
only as needed to establish intent, state or provenance. If there is no index,
use the project's existing document pattern. The Provenance Map is not a story;
leave it out of the count.

For each Story, name its **State**:

| State | Means | Signal |
| --- | --- | --- |
| `Ready` | Its own trigger has been Met | the `Distill when` line is True now |
| `Answered` | A Destination already Decides it | a rule, a doc or a test Says it |
| `Merge` | It is the Same kind as another story | two roadmap items, two sightings |
| `Misplaced` | It Belongs to another repository | a to-do list for one project |
| `Moving` | Work still in progress | a half-built feature, gaps Named |
| `Blocked` | A named Dependency stands in front | a decision, a tool or another story |
| `Waiting` | Its trigger is not Met | one lived instance, no second |
| `Speculative` | Marked a hunch by its author | `(SPECULATIVE)` in the index |

Check `Answered` with a grep over the Destinations,
never from memory: `rules/`, `values/` and `patterns/` in the Canon;
`docs/`, the code, its tests and the contract in a code project.
Quote the line that Answers it.

## Present

One Table: number, story, State, one line of why.
Then the Count, `N of 9`.

Then up to three numbered Options, only among `Ready`,
`Answered`, `Merge` and `Misplaced`. Recommend one, and say why.
`Moving`, `Waiting` and `Speculative` are listed, never offered.
A `Blocked` story is never offered either; its Unblocking step is,
when that step can be taken now.
If nothing can Move, say so and stop.

The user Picks by number. A number from the Table
that is not an Option: read that Story, name what it Waits on,
and offer the nearest Move it allows, if any.

## Land

One Story per turn. Read it whole before moving it.

Use the Destinations from the project's map, plus two that Hold anywhere:

| Destination | When |
| --- | --- |
| another repository | work that Belongs there; name it, do not do it |
| nowhere | the Destination already Says it |

1. Land the piece where it Governs.
   Keep the author's wording and Rhythm; move it, do not rewrite it.
   A new line in the Canon follows [DeLaCase](../de-la-case/SKILL.md).
   A Rule landing in code Lands with the Test that Holds it.
   A Decision landing in `docs/` Lands in every language,
   per [Translations](../../../rules/translations.md).
2. Delete the Story, or shorten it to what still Waits.
3. Update the Index and the Count in `STORY.md`,
   and the Door when a new doc Appeared: `MEMORY.md` in a code project.
4. Grep for links to the old path, and repair every one.

A Story that only shrinks Keeps its slot.
Say so; the Count does not move.

## What it Returns

```text
✅ Merged Own linter and LiteLLM into stories/roadmap.md
STORY.md is at 8 of 9.

Next — commit the merge.
```

## Bounds

- Never promote unprocessed private chaos into a Story. In the Canon, the user
  removes the Person first; elsewhere follow the project's privacy boundaries.
- Never assume one source record equals one Story, or one document contains
  only one Story.
- Never Distill a `Moving`, `Blocked`, `Speculative` or `Waiting` story unless the user insists,
  and then name what the project Loses by taking it early.
- Never land two Stories in one turn.
- Never promote a tenth Story while nine are Full.
- Never rewrite a human's Rhythm to make it fit.
- Never commit; [commit-commit-commit](../commit-commit-commit/SKILL.md) does that.
- Never do another repository's work from here; name it and hand it over.

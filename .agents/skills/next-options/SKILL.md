---
name: next-options
description: "Frame up to three grounded ways forward. Use when the user asks ‘what are our options?’, ‘where do we go from here?’, ‘what can we do next?’, ‘¿qué opciones tenemos?’ or when several viable paths compete without a clear next step. Present only distinct choices and wait for the user; do not choose or execute one. If progress is blocked, invoke break-free and name the blocker plainly."
---

# NextOptions

A Decision Frame, not a session Opening or an execution Step.
It shows up to three viable paths from the state already in view,
then leaves the choice with the user.

The session Opening Lives in [YoYoYo](../yo-yo-yo/SKILL.md).
The chosen Step Lives in [NextEpisode](../next-episode/SKILL.md).

## Find the Current Position

Use the smallest evidence that shows where the work stands:

1. The user's current intent and decisions in this conversation.
2. The last explicit choice or next step, when one exists.
3. A relevant handoff, plan, or repository state when the work depends on it.

Do not scan unrelated paths or infer intent from a diff alone.
If the current position is still unclear, ask one focused question
instead of inventing options.

Check whether progress is blocked before framing choices.
When an unavailable input, access, decision or external change prevents
progress, invoke [BreakFree](../break-free/SKILL.md) to name the
dependency and what would clear it. Say plainly when we are blocked;
do not disguise the diagnosis as a menu or invent paths to fill one.

## Frame Three Paths

Offer one to three numbered choices, only when each is distinct and feasible.
Use fewer when fewer real paths exist; never pad the list with a redundant,
unavailable or invented choice. Each choice Names one outcome and the scope
it would touch. Keep them at the same level of commitment and detail.

Ground every choice in the current work. When there is one active intent,
show only the distinct ways to advance it. When the work has split,
make separate intents visible only if each is a real path. Include a way
to resolve a dependency or pause only when it is useful and available.
Do not turn `Now` or `Later` labels into options by themselves.

State a tradeoff only when it helps the user distinguish the choices.
Do not rank the choices, recommend one, or take action before the user picks.
A numbered selection can continue through `NextEpisode`.

## What it Returns

Keep the framing short and grounded:

```text
**Where we stand** — The current change adds portable skill dependencies;
validation of malformed lock files is still open.

1. **Finish validation** — Add the missing malformed-file cases.
2. **Check host compatibility** — Try the same lock from another Agent host.
3. **Pause with a checkpoint** — Record the open validation for a later session.
```

Show fewer than three items when the state offers fewer distinct paths.
If all progress depends on a missing input or external change, say we are
blocked and name it; route the diagnosis through `break-free`.
Do not include a separate `Next` line. The user chooses the path.

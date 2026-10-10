---
name: sexy-back
description: "Return to the purpose that started the current conversation thread, preserving pending digressions. Use when the user says sexy back, sexyback, vuelve al origen, or asks to resume the original topic. Mentioning the song or editing this skill does not activate it."
---

# SexyBack

The Name Recalls *SexyBack* by Justin Timberlake.
A Quipu can be Rolled up; the detour stays in its cord
while the hands return to the Origin.

## Find the Origin

Read the current conversation for the user's initial Purpose
and the later corrections that still govern it.
Distinguish that purpose from the latest knot and its standing next step.
An explicitly replaced or cancelled intent is not an Origin to resume.

Use `.handoff.md` only when it actually records that origin;
its current `Intent` may describe a digression.
If the origin cannot be recovered, ask which topic the user means.
Never invent it from the skill's name or the most recent recommendation.

## Roll up the Detour

Briefly name the original Topic and where it stands.
Preserve any unresolved decision or useful result from the digression
in a short note; do not copy the conversation or undo completed work.
If returning changes the durable working state,
use [OneTwoCheckpoint](../one-two-checkpoint/SKILL.md)
to record the resumed intent and the pending detour.
Conversation alone needs no checkpoint.

## Resume

Name and take one concrete Step toward the recovered purpose
within the user's existing authorization.
If the purpose is already fulfilled, say so and ask what they want to revisit.
If the next action depends on missing information, ask for that information.
Do not perform a pending digression merely because it has a `Next` line.

Answer in the conversation's language with ordinary capitalization.
Make the recovered topic visible, report the step taken,
and mention the saved detour only when it remains useful.

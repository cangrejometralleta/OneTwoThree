---
name: one-two-unblock
description: "Assess whether progress is possible, externally blocked, caught in repeated reasoning, or out of useful hypotheses. Use when the user asks whether we are stuck or what is missing to continue, or when attempts repeat without new evidence."
---

# OneTwoUnblock

Make the next Movement possible by naming what prevents it.
A failed attempt is Evidence; repeating its explanation adds none.

## Read the Current Knot

Recover the active Purpose, what would count as done,
and the latest attempts with their results.
Read only the conversation and relevant artifacts;
use `.handoff.md` as supporting context after checking that it is current.
Do not substitute the original topic for the active one.

Separate Facts from hypotheses. Name what each attempt changed or taught.
When evidence is missing, make one focused, authorized check
that could change the diagnosis. Do not repeat an unchanged failed attempt.
Do not ask the user for information already available locally.

## Determine the State

Assess these questions together; a Cycle can coexist with a Blocker.

- **Complete** — the requested outcome is already achieved.
  Report the evidence and stop; no new ideas are needed.
- **Can continue** — a useful authorized action can change the result
  or distinguish live hypotheses. Name that action and what it would reveal.
- **Blocked** — progress toward the requested outcome requires a specific
  unavailable input, access, decision or external change.
  Name the dependency, the evidence and who or what can supply it.
  State any independent work still possible.
- **Cycle** — attempts return to the same premise, action or conclusion
  without new evidence or a changed condition.
  Point to the repeated attempts; do not diagnose a cycle merely from duration.
  Stop that repetition and identify one different observation or approach.
- **No useful hypothesis** — the considered approaches are exhausted
  and no discriminating next action can currently be named.
  Say which approaches were exhausted and what evidence or perspective
  could generate another. This does not prove the problem is impossible.

When evidence is insufficient, say what remains Unknown.
Do not call uncertainty, difficulty or one failure a confirmed blocker.
These are conversational diagnoses; do not change a managed goal's status
unless its own lifecycle requirements are satisfied.

## Open a Way Forward

If a focused check can break the cycle or resolve uncertainty,
name and perform that one check within existing authorization.
Report its result and update the diagnosis; do not start a retry loop.

When a dependency is required, ask for the smallest missing input
and stop dependent work. Continue independent work only when it serves
an already authorized purpose. Do not request broader permissions by default.

When no useful hypothesis remains, identify the kind of help needed:
a concrete example, another observation, domain knowledge or a fresh reading.
Request delegation only when authorized; do not invent a new idea to fill a slot.
When returning to the Origin would help, offer
[SexyBack](../sexy-back/SKILL.md) with the topic; let the user choose.

## Return the Diagnosis

Use the conversation's language and ordinary capitalization.
Keep the result to the state, its evidence and what makes progress possible.
For example:

```text
Estamos en un ciclo: repetimos la misma consulta con el mismo resultado.
Falta comprobar si el archivo leído es el que usa el proceso.
Comprobé la ruta activa: apunta a otra copia. Podemos continuar desde allí.
```

If blocked, say exactly what is needed and how its arrival would unblock work.
If no blocker was found, say so with the next useful action.
Checkpoint a durable change or a diagnosis that changes the next step through
[OneTwoCheckpoint](../one-two-checkpoint/SKILL.md).
A conversation-only assessment needs no checkpoint.

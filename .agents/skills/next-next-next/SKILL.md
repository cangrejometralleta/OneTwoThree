---
name: next-next-next
description: "Advance to the next recommended course of action and take exactly one Step. Use when the user says next, sigue, continue, go on, or asks what to do now and wants it done rather than listed. It finds the standing recommendation, confirms it still holds, takes one Step, and names the Step after it."
---

# NextNextNext

A Continuation, not a Plan and not a Session Opening.
It Takes the Step that was already Named,
Does it once,
and Names the one that Follows.

The next Turn lives in [ByeByeBye](../bye-bye-bye/SKILL.md).
The canon lives in [Change Growth](../../../rules/change-growth.md).
The Shape lives in [The Lever and the Tape](../../../patterns/the-lever-and-the-tape.md).

```mermaid
flowchart TD
    START["sigue"] --> FIND{"Recommendation Found?"}
    FIND -- No --> OPEN["Defer to YoYoYo · Stop"]
    FIND -- Yes --> HOLDS{"Still Holds?"}
    HOLDS -- No --> RENAME["Name what Changed · Restate one Step"]
    HOLDS -- Yes --> SIZE{"One Step?"}
    RENAME --> SIZE
    SIZE -- No --> SPLIT["Take the first independent Part"]
    SIZE -- Yes --> DO["Take the Step"]
    SPLIT --> DO
    DO --> VERIFY["Verify what the Step Claims"]
    VERIFY --> REPORT["Report Done · Name the next Step"]
```

## When it Runs

Invoke when the User Says `next`, `sigue`, `continue` or `go on`,
or Asks what to Do now and Wants it Done, not Listed.

Do not Invoke to Open a Session; YoYoYo Reconstructs State.
Do not Invoke to Close one; ByeByeBye Writes the Handoff.
Do not Invoke for a new Ask that Carries its own Intent.

## Find the Recommendation

Read the smallest Source that already Names a next Step:

1. **This Thread** — the last `Next`, `Now` or `Later` Line Stated here.
2. **Handoff** — `.handoff.md` at the Repository Root, when Present.
3. **Plan** — an approved Plan or a Checklist the Work Follows.
4. **State** — Branch, staged, unstaged and untracked Paths.

The Thread Wins over the Handoff when both Speak.
When no Source Names a Step, do not Invent one.
Say so, and Defer to YoYoYo.

## Confirm it Holds

A named Step can Expire before it Runs.
Check three Things:

1. Does the Repository still Look the way the Step Assumed?
2. Did the User Change the Intent since the Step was Named?
3. Was the Step already Done by a later Turn or another Session?

When it Holds, Take it.
When it Does not, Name what Changed in one Line,
Restate one Step, and Take that one.
Never Take a Step the Evidence already Contradicts.

## Take one Step

One Step is one Outcome that can be Verified alone.
When the Recommendation Carries more than one Outcome,
Take the first independently verifiable Part
and Leave the Rest as `Later`.

Do the Work. Do not Narrate the Options.
Do not Start a second Part because the first was small.
When the Step Crosses a Boundary the User has not Authorized —
a Push, a Deletion, an outward Send —
Stop before it and Ask.

## Verify

Verify what the Step Claims, and nothing Wider:

- a File Change, by Reading the Result or Running its Check;
- a Command, by its Exit and its Output;
- a Rule Edit, by the Link and the Frontmatter it Names.

When Verification Fails, Report the Failure and Stop.
A failed Step is the next Step.

## What it Returns

```text
✅ **Done** — Linked the reading Rotation into the Rules index.

**Verified** — Path Resolves; Frontmatter Parses.
**Next** — Reload the Skills in the local Client.
**Later** — Translate the Rotation Labels.
```

When nothing Holds:

```text
⚠️ No standing Recommendation Found.
Open the Session with YoYoYo, or Name the Step.
```

## Bounds

- Never Take a Step no Source Named.
- Never Take two Steps in one Invocation.
- Never Replace a Session Opening or a Session Closing.
- Never Write, Rewrite or Delete `.handoff.md` here.
- Never Push, Delete or Send without explicit Authority.
- Never Report a Step Done that Verification did not Confirm.
- Never Silently Redefine the Step the User was Waiting on.

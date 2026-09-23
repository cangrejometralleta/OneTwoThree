# OneTwoThree

OneTwoThree is a set of agent Rules that also forms a manifesto.
The Rules Guide how Agents read, write and collaborate.
The Manifesto Grounds them in Minimalism, calm and care for human attention.

> Because we Hate making documentation.
> 
> The Limit behind all of this was lived before it was written.
> Autistic burnout Taught it, minimalism only named it.
> 
> Thanks [De La Soul](https://en.wikipedia.org/wiki/De_La_Soul), Grandma COBOL and [John Cage](https://en.wikipedia.org/wiki/John_Cage)
> for inspiring this Project.

## The Name Carries a Rhythm

OneTwoThree Draws its name from [De La Soul](https://en.wikipedia.org/wiki/De_La_Soul)'s [*The Magic Number*](https://en.wikipedia.org/wiki/The_Magic_Number), with its phrase “three is the magic number”.
[*4 noviosS*](https://www.youtube.com/watch?v=ucrvnu5a8NQ) by [Six Sex](https://es.wikipedia.org/wiki/Six_Sex), produced by King Doudou, and the lyrics of [*Perfect (Exceeder)*](https://en.wikipedia.org/wiki/Perfect_%28Exceeder%29) by [Mason](https://en.wikipedia.org/wiki/Mason_%28musician%29) vs [Princess Superstar](https://en.wikipedia.org/wiki/Princess_Superstar) also Inspired its reading and writing Cadence.
Their counting Phrases help turn Words into a Pulse: an entrance, an emphasis and room to breathe.

The Author Hears a direct musical Influence between *Perfect (Exceeder)* and *4 noviosS*.
That Connection is his Interpretation as a listener.
Together, these References gave the Project a Rhythm to read and write by.

OneTwoThree and one-two-three Name the same Project. Any form of emphasis is Accepted here: camel case, kebab case, spaces or none. The Name is made to be Counted, never said in one breath.

## License

This is free and unencumbered software
released into the public domain.
Do whatever you want with it.
[The Unlicense](https://unlicense.org)

## Documents

Three Documents Hold the Manifesto.  
This Page is the Door that indexes them.

- [Values](VALUES.md) — why it Exists.
- [Rules](RULES.md) — how it Applies.
- [Patterns](PATTERNS.md) — where the Why comes from.

Why / How / Where:  
the Triad Lives in the Structure itself.

### RoTaTion Connects the Three

[De La Soul](https://en.wikipedia.org/wiki/De_La_Soul) [RoTaTion](patterns/de-la-soul-rotation.md) Connects three Elements through three Pairs, with no fixed center.
Enter through any Element and follow the relationships.

```mermaid
flowchart LR
    VALUES["VALUES<br/>the Why"] -- "Why becomes How" --> RULES["RULES<br/>the How"]
    RULES -- "How Reveals its Root" --> PATTERNS["PATTERNS<br/>the Where"]
    PATTERNS -- "the Root Grounds the Why" --> VALUES
```

Three Elements, three Pairs, no Center.
Remove the Center and the Shape still Turns.

## How to Read it

Three Layers connect the Purpose to a concrete Action:

```mermaid
flowchart TD
    PURPOSE["Purpose<br/>Protect human attention"]
    STRUCTURE["Structure<br/>Keep one task visible"]
    ACTION["Action<br/>Write down the next step"]
    PURPOSE -->|Guides| STRUCTURE
    STRUCTURE -->|Makes concrete| ACTION
```

These Layers Offer a Path from the general to the particular.
RoTaTion Keeps the Relationships open to another reading.

- Use it as agent Rules,
  and read the manifesto that Explains their purpose.
- Identify important Entities and their Interactions with up to three emphasis capitals
  per passage between punctuation marks; the grammatical initial is free.
  Two entities and one interaction are a useful shape, never a required formula.
  Proper names and acronyms keep their established Spelling.
  [DeLaCase](rules/de-la-case.md) holds the convention; three spent is a Ceiling, never a quota.
- Every Line is written  
  to be read in one heartbeat.
- Three is a Source, not a count.  
  You Derive from Three,  
  you do not reach it.

## 🕊️ Work with Dove and the Skills

[Dove](.agents/agents/dove.md) Holds the voice.
The Skills Define the Operations; their linked instructions hold the details.

| When | Skill | What it does |
| --- | --- | --- |
| `yo dove` or resume earlier work | [yo-yo-yo](.agents/skills/yo-yo-yo/SKILL.md) | Synchronizes the project branch, reconstructs context and names one next Step. |
| `next`, `sigue` or a selected option | [next-next-next](.agents/skills/next-next-next/SKILL.md) | Takes one recommended step, verifies it and names the Following step. |
| A durable edit, decision or validation | [one-two-checkpoint](.agents/skills/one-two-checkpoint/SKILL.md) | Saves the current thread in `.handoff.md` during Work. |
| A change starts growing | [one-two-growth](.agents/skills/one-two-growth/SKILL.md) | Checks whether one intent still Holds the change. |
| Ask to organize, commit and push | [commit-commit-commit](.agents/skills/commit-commit-commit/SKILL.md) | Groups changes by feature, validates and commits each group, then Pushes once after all succeed. |
| `bye dove` or request a handoff | [bye-bye-bye](.agents/skills/bye-bye-bye/SKILL.md) | Expands the checkpoint into a closing handoff and Stops. |

The [Session Diagram](patterns/the-lever-and-the-tape.md) Connects these operations.
A new Request with its own Intent starts its own Work.
Opening a Session Names the next Step; Continuation takes it when Requested.
Publication and Closing each need their own Request.

Supporting Skills Shape the Work as it happens:
[one-two-refactor](.agents/skills/one-two-refactor/SKILL.md) guides code,
[de-la-case](.agents/skills/de-la-case/SKILL.md) converts names and prose,
and [one-two-output](.agents/skills/one-two-output/SKILL.md) shapes terminal output.
The typography convention is [DeLaCase](rules/de-la-case.md).
The commit workflow is `commit-commit-commit`.

## Connect a Project

[one-two-update](.agents/skills/one-two-update/SKILL.md) Installs or updates the canon connection.
[one-two-reload](.agents/skills/one-two-reload/SKILL.md) Connects the active client to the selected skills and agents.
Session opening Synchronizes the project's branch; canon updates follow their own connection.

```mermaid
flowchart LR
    UPDATE["one-two-update"] --> FORMAT{"Distribution"}
    FORMAT -- "Clone" --> CLONE[".agents/canon<br/>Selected relative links"]
    FORMAT -- "ZIP" --> ZIP["Generated snapshot<br/>.agents/distribution.json"]
    CLONE --> RELOAD["one-two-reload"]
    ZIP --> RELOAD
    RELOAD --> VERIFY["Reload the client · Begin a new chat<br/>Verify discovery"]
```

Existing Installations Keep their chosen Mechanism and local customizations.
A [portable ZIP](.agents/skills/one-two-update/references/zip.md) Records its source commit and file inventory.
It is a Snapshot; the live canon remains the head of `main`.
Filesystem Validation Proves the Links; client Discovery Proves the Load.

## Agents Work among Others

An Agent Needs more than a Goal.
It needs a Society: explicit authority, independent signals,
the right to stop and a human it can escalate to.

Recent security research shows why those boundaries must be Designed,
not assumed. [Coercion](rules/coercion.md) Carries the rules:
never invent Permission, never retaliate,
and never let pressure make harm look necessary.

## How Context Becomes Canon

Three Stages Carry a piece of Life into the canon,
and only the third one stays.

- **CHAOS.md** Holds the raw life.
  Private, never Committed,
  names and dates still in it.
- **Stories** Holds the same piece
  with the person removed.
  Public, Staged, still not canon.
- **Values, Rules and Patterns** hold what Survived.

```mermaid
flowchart LR
    CHAOS["CHAOS.md<br/>the raw Life<br/>private, never Committed"]
    STORY["STORY.md<br/>the Person Removed<br/>public, still not Canon"]
    CANON["VALUES · RULES · PATTERNS<br/>what Survived<br/>the Canon"]
    CHAOS -- "Strip the Person" --> STORY -- "Strip the Story" --> CANON
```

Between the first and the second you Strip the Person.
Between the second and the third you Strip the Story.
What is left is the Belief, the Rule or the Root.

- A Belief Goes to Values.
- A Rule an agent can run Goes to Rules.
- A Root Goes to Patterns.
- Delete it from Stories once it lands.

CHAOS.md never Empties, because a source never empties.
Stories Empties, because a passage is meant to.
Nine Entries Fill the Passage.
Distill before you promote a Tenth.

.gitignore Names CHAOS.md out loud.
Patterns Explains why,
under The Right you have to Invoke.

An empty Stories means the canon is Current.

The name Comes from Patterns, under Chaos is a Source.

AGENTS.md is not this Passage.
It is the short Pointer agents load on their own —
canon first, Stories named as notes, nothing undistilled repeated there.

## Jokes

[Jokes](jokes) is not a Stage of that passage.
It Runs beside it, and it never arrives.

A human Writes them by hand, to tune the humour of the language.
The Rules Teach an agent to count beats; none of them teaches timing.
An Agent Reads the directory to hear the Voice,
and it never adds a line to it.

Each Joke Stays in the Language it was born in, Spanish or English.
Most of them are Grammatical, and a translated pun is a sentence about a pun.

No one Explains a Joke there, least of all an agent.
Explaining Kills it, and there is no careful way to do it.

[.canonignore](.canonignore) Says the same thing to a machine.

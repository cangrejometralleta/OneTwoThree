---
name: one-two-update
description: "Fetch the latest OneTwoThree canon from main on its remote repository and connect the current project to a reduced copy of it, without hand-copying any file. Use when the manifesto rules may be stale, when a project first adopts OneTwoThree, or when the user asks to update, sync or pull the canon."
---

# OneTwoUpdate

A Fetcher of the canon, not a copier of it.
It Brings `main` from the remote OneTwoThree,
Keeps only what Governs, Loads only the Index,
and Leaves every local file the project owns untouched.

The canon Lives in [The Head is the Canon](../../../rules/the-head-is-the-canon.md),
[Vendor Integration](../../../rules/vendor-integration.md)
and [Canonignore](../../../rules/canonignore.md).
The loading Lives in [OneTwoReload](../one-two-reload/SKILL.md).

The source is one Address:

```text
https://github.com/cangrejometralleta/OneTwoThree.git   branch: main
```

```mermaid
flowchart TD
    START["update the canon"] --> HERE{"Inside OneTwoThree?"}
    HERE -- Yes --> PULL["Fast-Forward main · Stop"]
    HERE -- No --> LINK{"Link Exists?"}
    LINK -- Yes --> FOLLOW["Follow the Mechanism already Chosen"]
    LINK -- No --> ADOPT["sparse, blobless Clone · relative Link"]
    FOLLOW --> FETCH["Fetch main"]
    ADOPT --> FETCH
    FETCH --> DIVERGE{"Local Edits in the Canon?"}
    DIVERGE -- Yes --> STOP["Report the Divergence · Change nothing"]
    DIVERGE -- No --> APPLY["Fast-Forward to origin/main"]
    APPLY --> REDUCE["Prune to the governing Paths"]
    REDUCE --> VERIFY["Verify Links · Read the new Canonignore"]
    VERIFY --> REPORT["Report the Range · Name OneTwoReload"]
```

## When it Runs

Invoke when the user Asks to update, sync or pull the manifesto,
when a project first Adopts OneTwoThree,
or when a rule read here Disagrees with the rule named upstream.

Do not invoke to load skills into a client; OneTwoReload Does that.
Do not invoke to open a session; YoYoYo Pulls the project's own branch.

## Two Reductions, not one

`Memory` and `Context` are cut by different Means.
Do both, and never confuse one Report for the other.

**Memory** is what Lands on disk. Check out only the governing Paths:

```text
AGENTS.md  RULES.md  VALUES.md  PATTERNS.md  .canonignore
rules/     values/   patterns/
```

**Context** is what Enters the prompt. Load only the Index:

```text
AGENTS.md  RULES.md  VALUES.md  PATTERNS.md
```

The index Names every rule with one line and one link.
A Body is read the turn its rule is invoked, and not before.
Three indexes cost about 7 KB; their bodies cost about 50 KB;
the whole repository Costs about 2.4 MB.

Never preload `rules/`, `values/` or `patterns/` wholesale.
The index Exists so an agent can know a rule is there
without paying to Read it.

## Inside the Canon Repository

First Ask whether the current repository *is* OneTwoThree.
Compare the Origin address, not the directory name.

When it is, there is nothing to Vendor or reduce.
Fast-Forward `main` from `origin`, Report the range, and Stop.
Never vendor the Canon into the repository that writes it.

## What it Reads

Read only what Names the connection:

1. **Origin** — the current repository's Remotes.
2. **Link** — an existing Submodule, Subtree, Clone or symbolic Link
   that already Points at OneTwoThree.
3. **State** — whether that link Carries local commits or dirty paths.
4. **Canonignore** — the incoming `.canonignore`, read after the Fetch.

`AGENTS.md` or local documentation may Name a different entrance.
Follow what the project Documents over what this skill assumes.

## Follow the Mechanism already Chosen

A project that already tracks the Canon has Answered this question.
Detect which answer it Gave, and use it:

| What you Find | What you Run |
| --- | --- |
| A Submodule pointing at OneTwoThree | `git submodule update --remote` |
| A Subtree | `git subtree pull --prefix <path> <url> main --squash` |
| A Clone under a Cache or Vendor Path | `git fetch origin main` then Fast-Forward |
| A symbolic Link to a Clone elsewhere | Update the Clone the Link Resolves to |
| Nothing | Adopt, below |

Re-apply the Path reduction after a fetch that widens it.
Never convert one Mechanism into another to make the update easier.
A project changes how it Tracks the canon on purpose, never as a side effect.

## Adopt, when Nothing Exists

Choose the smallest thing that Keeps one source:

```sh
git clone --filter=blob:none --sparse --branch main \
  https://github.com/cangrejometralleta/OneTwoThree.git <vendor-path>
git -C <vendor-path> sparse-checkout set \
  AGENTS.md RULES.md VALUES.md PATTERNS.md .canonignore rules values patterns
```

`--filter=blob:none` Leaves the history on the server until it is asked for.
`--sparse` Leaves the ungoverned paths off the disk entirely.

Then point the Project at it with a relative symbolic link,
and name the Path and the Link in `AGENTS.md`.

Ask before the first Adoption. It Adds a dependency,
and a dependency added silently is a dependency nobody Chose.

## Never Overwrite what the Project Wrote

Fast-Forward only.

When the local canon Carries commits that `origin/main` does not,
Stop and report the Divergence. Those Commits are either
a Fork worth keeping or an Edit that belonged upstream,
and only the user can Say which.

- Never Force Push, Force Pull, Reset hard or Discard local commits.
- Never Merge or Rebase the canon into the project's own history.
- Never edit a File under the vendored canon to resolve a conflict.

## Read the new Canonignore

The canon Arrives with its own boundary; read it after every fetch.
Paths it lists are carried and not taught — `jokes/`, `stories/`,
`STORY.md`, `chaos/` among them.

The reduction already Keeps them off the disk,
so a path that appears despite it is a Signal, not a convenience:
the boundary Moved upstream, and the sparse set has to move with it.

Read them freely where they Exist. Copy none of them into the Project.
An update that teaches from an ignored path Made the canon worse,
not fresher.

## Verification

1. The fetched ref is `main` from the OneTwoThree address, not a Fork.
2. The link Resolves from the project root.
3. `RULES.md`, `VALUES.md` and `PATTERNS.md` Exist at the linked root.
4. Every link inside those three indexes Resolves to a checked-out body.
5. No path outside the governing set Landed on disk.
6. The range between the old and new commit is Nameable.

A fetch that Moved no ref is `Already Current`, not `Updated`.
An index link that resolves to nothing means the sparse Set is too narrow.

## What it Returns

```text
✅ Canon Updated — 63a3010..f3bdfa6 (4 Commits)

**Source** — cangrejometralleta/OneTwoThree, main.
**Entrance** — vendor/onetwothree, linked from .agents/canon.
**Changed** — 2 Rules, 1 Pattern.
**Reduced** — 58 KB on Disk, 7 KB Loaded.

**Next** — Reload the Client with OneTwoReload.
```

When the ref did not Move:

```text
✅ Already Current — main at f3bdfa6.
```

When the local canon Diverged:

```text
❌ Canon Diverged — 2 local Commits not on origin/main.
Nothing Changed. Name whether they are a Fork or an Edit to Send upstream.
```

## Bounds

- Never copy a canonical File into the project by hand.
- Never vendor the Canon into OneTwoThree itself.
- Never fetch from a Fork while claiming the canon.
- Never adopt a new Dependency without asking.
- Never Fast-Forward past local Commits, and never discard them.
- Never preload a Body the index can name for free.
- Never widen the sparse Set to make one read easier; read it on demand.
- Never teach from a Path `.canonignore` lists.
- Never report `Updated` when no ref Moved.
- Never claim the client sees the new canon; OneTwoReload Proves that.

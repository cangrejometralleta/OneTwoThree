---
name: one-two-unbrand
description: "Unbrand a vendor from a repository on two surfaces: Git history (authorship and attribution trailers) and the tree (vendor entrances, copies and settings in the shared .agents folder structure). One reference file per vendor; Anthropic and Claude, and Google and Gemini, are covered. History rewrites go through one-two-purge with its Bundle and Final Confirmation. Use when the user asks to unbrand, debrand or remove Claude, Anthropic, Gemini or another agent from a repository, or to check that the .agents structure keeps its shape."
---

# OneTwoUnbrand

A Finder of the Name an agent left behind, not a Rewriter.
It knows the Marks, counts them, and fixes what the Rule says is wrong.
The History Rewrite goes to [OneTwoPurge](../one-two-purge/SKILL.md),
which Holds the Bundle, the Confirmation and the Push.

The Rule Lives in [Unbrand](../../../rules/unbrand.md),
and the Shape of the tree in [Vendor Integration](../../../rules/vendor-integration.md).
History is Appreciated, so a Rewrite is a last Resort, never a Cleanup habit.

## Call

`/one-two-unbrand <vendor> [history|shape]`, as in `/one-two-unbrand claude shape`.

The vendor Names a file in `references/`:

| Vendor | Reference |
|---|---|
| `claude`, `anthropic` | [anthropic-claude.md](references/anthropic-claude.md) |
| `gemini`, `google` | [google-gemini.md](references/google-gemini.md) |

A vendor without a file is not guessed:
stop, say so, and ask the user for its identities, trailers and paths.
No Mode means both Surfaces, Counted first and fixed one at a time.

`leave`, removing a vendor's entrances altogether, is Not built.
Say what it would remove, and stop.

## Each reference holds

- **History marks** — identities, trailers, the Detect, the mailmap and the Callback.
- **Tree marks** — entrance paths, the context file, the settings files.
- **Leave alone** — what the Rule keeps, and what is not Branding.

## History

1. **Count** — run the reference's Detect over every ref, Counts only.
2. **Brief** and **Confirm** — as [OneTwoPurge](../one-two-purge/SKILL.md) asks.
3. **Rewrite** — `git-filter-repo` through the reference's mailmap and Callback.
4. **Verify** — the same Counts at zero, then the purge Verify.

## Shape

A vendor Entrance is a Link, never a copy: the vendor finds its folder,
and the Instructions stay in `.agents/`. The Check Keeps that Shape, never removes the Entrance.

Classify every path the reference names, and every other Entry in `.agents/`:

| Class | Example | Fix |
|---|---|---|
| Entrance | `.claude/skills`, a symlink | Keep it. Check it resolves into `.agents/`. |
| Adapter | `.agents/agents/dove.toml` | Keep it while it is small and Points to the shared markdown. |
| Copy | a real folder where a link belongs | Move the content into `.agents/`, replace it with a relative link. |
| Stray | a vendor file inside `.agents/` | Move it out beside the vendor's own folder. Never delete it. |
| Settings | `.claude/settings.json` | Leave it. It is the user's. |
| Signature | a vendor Name as author in a file or comment | Remove it where it is only the Author. |

Detect, with Counts and Paths only:

```bash
# entrances: link, broken link, copy or absent
for p in <paths from the reference>; do
  if [ -L "$p" ]; then [ -e "$p" ] && echo "link   $p -> $(readlink "$p")" || echo "broken $p"
  elif [ -e "$p" ]; then echo "copy   $p"
  else echo "absent $p"; fi
done
# strays: anything in .agents/ that is not agents/ or skills/
find .agents -mindepth 1 -maxdepth 1 ! -name agents ! -name skills
```

A link Target is Relative, so the tree Travels.
An entrance the client already Finds is Clutter, not a Gap: do not add it.
A link does not convert formats; Check the vendor reads the shared file before Claiming support.

A Fix that moves or replaces a file needs a Brief and the Final Confirmation.
It changes the working tree only: no Commit, no Push, no Stage.
Reporting a Class changes nothing.

## Gotchas

- `git filter-repo --refs` takes literal ref names, never Globs.
  A Glob parses zero commits and exits zero.
  Build the list with `git for-each-ref`, and skip `refs/remotes/origin/HEAD`.
- Use `--partial`, so the Remote and the other refs stay.
- Rewrite every Branch, tag and remote-tracking ref in one pass,
  so a commit keeps one new ID wherever it Appears.
- `filter-repo --force` ends with a hard Reset.
  Back up an uncommitted file first, and restore it after.
- Look at `git worktree list`.
  A linked worktree on a detached old commit keeps the old history reachable,
  and a Prune will not remove it. Name it in the Brief.
- `refs/archived-worktrees/*` and `refs/stash` can Carry old commits too.
- Run the commands under an explicit `bash`.
- A Force Push is one branch per push,
  each with `--force-with-lease=<branch>:<old-id>`, from the rewritten remote tip.

## Bounds

- Never rewrite History without the Final Confirmation of [OneTwoPurge](../one-two-purge/SKILL.md).
- Never push, prune, move or remove what the Brief did not name.
- Never Delete a commit; the identity is replaced, the Work stays.
- Never Delete a vendor settings file, and never touch `.claude/settings.json`.
- Never edit the permission rules to make a Fix pass.
- Never run a Rewrite from a subagent: the Final Confirmation is the user's to give.
- Report Counts and States, never the lines Matched.

## What it Returns

```text
✅ Unbrand Verified — Anthropic and Claude
- History: 0 Trailers, 0 Identities across every ref
- Shape: 5 Entrances ok, 0 Copies, 0 Strays
- Trees: identical, commit Counts equal
- Remote: Force Push Performed | Not Performed
- Recovery Bundle: Retained until the user asks
```

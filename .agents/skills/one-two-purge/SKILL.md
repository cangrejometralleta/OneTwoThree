---
name: one-two-purge
description: "Find an exact sensitive value across files, Git history and local shell histories, then remove it through an explicit detect, confirm and purge workflow. Also removes agent authorship and attribution trailers from Git history as a last resort. Use when a secret, token, credential or private value may have entered repository or terminal history, or when the user asks to unbrand it."
---

# OneTwoPurge

A Remover of sensitive history, not a promise that exposure never happened.
It Detects without echoing, names the affected surfaces,
confirms the destructive boundary, then purges and Verifies.

The first action after a credential exposure is Revocation or Rotation.
History Rewriting Reduces Distribution; it does not make a leaked secret safe.

## The Value

Never ask the user to paste a sensitive Value into chat.
Never place it in a command Argument, generated script, log or report.

Ask the user to enter it directly into an interactive terminal with Echo
disabled and keep it in a temporary environment Variable such as
`PURGE_VALUE`. The Agent Uses only the variable Name.

- Never print, expand, inspect or persist the Variable.
- Never enable shell Tracing while it exists.
- Reject an empty Value before every search or mutation.
- Clear it from the Environment when verification finishes.
- If the client cannot accept secret terminal Input without exposing it to the
  model, stop and give the user a local Command to run themselves.

## Detect

Read three Surfaces. Detection Changes nothing.

1. **Files** — tracked, untracked, hidden and ignored Files in scope.
2. **Git** — reachable Commits, Branches, Tags and local Reflogs,
   every other ref under `refs/`, the stash, and `git worktree list`.
3. **Shell** — known local history Files for Bash, Zsh, Fish and the active
   shell, plus project-local session Logs explicitly named by the user.

Use fixed-string matching by Default. Treat the Value as data, never a regular
expression. Exclude binary payload output and report only redacted Paths,
commit IDs, ref names and Counts.

Do not search unrelated home Directories merely because they are accessible.
Name every Path outside the current repository and ask before reading it.

Report Detection separately:

```text
Files: 2 Matches in 2 Paths
Git: 4 Commits across 1 Branch and 1 Tag
Shell: 3 Entries in Fish History
```

Never include matching lines or the sensitive Value.

## Confirm

Detection does not Authorize Mutation.

Before purging, give a short Brief of the Operation:

- the exact Files, refs and history stores that will change;
- any linked worktree or extra ref that still holds the old history;
- whether Git commit IDs will Change;
- whether Tags or Branches require replacement;
- whether a remote Force Push will be needed;
- that existing Clones, forks, caches and logs may retain the old value;
- which Backups will be created and when they will be removed.

Then ask the user for one Final Confirmation of the whole Purge.
It covers every step the Brief names, whichever of these apply:

1. Rewrite repository History.
2. Rewrite each shell history File.
3. Expire Reflogs and prune unreachable Git objects.
4. Force Push rewritten branches or tags.

Wait for the answer. A yes covers the Brief, nothing beyond it.

## Purge

Purge only the Surfaces the user confirmed.

### Files

Replace or remove the Value from the current working tree first.
Preserve file Structure and unrelated content.
Do not commit unless the user separately Requests a commit.

### Git

Prefer `git-filter-repo` when installed and supported by the Repository.
Read its installed Help before constructing the rewrite.
Use a replacement File or Callback that reads the value without placing it in
the command line, process list or persistent project files.

Before rewriting:

- require a clean or fully understood working Tree:
  `git filter-repo --force` ends with a hard Reset,
  so back up an uncommitted file first and restore it after;
- name each linked worktree on a detached old commit:
  it keeps the old history reachable until it is removed or re-pointed,
  and the user Decides which;
- pass `--refs` literal ref names, never Globs:
  a Glob parses zero commits and exits zero;
- record the current Branches, Tags, Remotes and `HEAD` without secrets;
- create a local recovery Bundle outside the repository with restrictive
  Permissions;
- name the Bundle path to the user;
- refuse to overwrite an existing Bundle.

Do not use `filter-branch` when `git-filter-repo` is available.
Do not delete original refs, expire reflogs or run garbage collection until the
rewritten history passes Verification.

Never Force Push beyond the Final Confirmation. Name the affected remote Refs
in the Brief. Use `--force-with-lease`, never unconditional `--force`,
when the remote State permits it.

### Shell

Stop or account for active Shells that may rewrite history on exit.
Create a permission-restricted Backup beside neither the repository nor its
tracked Files. Parse the native history Format; do not treat structured Fish or
multiline history as plain independent lines when that would corrupt entries.

Write a replacement File atomically, preserve Permissions and ownership,
then ask the user to restart or reload affected shell sessions.
Never clear an entire History when exact entry removal is possible.

## Unbrand

A vendor Name is not a secret, so the Value rules above do not Apply.
The Target is authorship and Trailers, never a Value to hide:
an agent identity as Author or Committer, and `Co-Authored-By`,
session-link or generated-by lines in messages.

[Unbrand](../../../rules/unbrand.md) says to keep them out in the first place.
Git history is Appreciated here, so this is a last Resort, not a Cleanup habit.
Before offering it, name what it Costs:

- every rewritten commit ID Changes, and so does every Descendant;
- clones, forks, vendored canons and snapshot Manifests keep the old IDs;
- signatures, tags and pull request Links Break;
- the contributors list on the host may take a while to Refresh.

[OneTwoUnbrand](../one-two-unbrand/SKILL.md) Knows each vendor's Marks,
counts them, and builds the mailmap and the message Callback.
It hands the Rewrite back here, under the same Bundle, Confirm
and Force Push rules as any other history Rewrite.
Replace the identity with the project's own; never Delete the commits.
Verify with the same Counts at zero before any Push.

### Every Branch

A Rewrite that Stops at `main` leaves every other Branch Carrying the old IDs.
Count them first, and Rewrite them in one pass,
so a commit keeps one new ID on every Branch that Shares it.

- Rewrite the local Branches and the remote-tracking refs together.
  Record each old tip in a recovery ref before the Rewrite.
- Verify every ref on its own: the Trailer Count is zero,
  the tree is identical to its old tip, and the commit Count is the same.
- Push each remote Branch with its own lease on its old remote ID.
  `--force-with-lease=<branch>:<old-id>`, one Branch per push.
- Push the rewritten remote tip, never the local one.
  A Branch with unpublished local commits Keeps them unpublished.
- A Branch that has no Remote is rewritten and Stays local.
- Name any Branch that Deploys on push in the Brief.
- Run the commands under an explicit `bash`.
  A list held in a Variable is not split into words by Zsh,
  and git then Receives one Argument glued from many.

## Verify

Repeat the same fixed-string Detection against every purged surface.

1. Current Files contain zero matches.
2. Rewritten reachable Git History contains zero matches.
3. Confirmed shell history Stores contain zero matches.
4. Refs and repository integrity checks pass.
5. The sensitive Variable is cleared.

A zero Match in rewritten Git does not prove remote caches or existing clones
forgot the Value. Say so.

Only after Verification may recovery refs be deleted, reflogs expired
and objects pruned, when the Final Confirmation named them.
Removing the Bundle and other backups stays a Request of its own.

## Boundaries

- Never reveal a sensitive Value back to the user.
- Never send it to a network Service or external Scanner.
- Never mutate global credential Stores without an explicit request.
- Never rewrite signed Commits or tags without naming that signatures Break.
- Never rewrite a shared Branch without naming the coordination required.
- Never remove audit evidence required by Law, policy or an active incident.
- Stop when repository ownership, remote authority or History format is unclear.
- Never edit permission rules to pass a denial; the user adds the rule.
- Never run a Purge from a subagent: the Final Confirmation is the user's to give.

## What it Returns

Report States, counts and remaining actions only:

```text
✅ Purge Verified
- Files: 0 Matches
- Git: 0 Matches in rewritten reachable History
- Shell: 0 Matches in 1 confirmed History Store
- Credential: Rotation still Required
- Remote: Force Push not Performed
- Recovery Bundle: Retained until final Approval
```

The Report never Contains the Value or a recoverable fragment of it.

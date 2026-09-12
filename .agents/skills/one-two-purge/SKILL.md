---
name: one-two-purge
description: "Find an exact sensitive value across files, Git history and local shell histories, then remove it through an explicit detect, confirm and purge workflow. Use when a secret, token, credential or private value may have entered repository or terminal history."
---

# OneTwoPurge

A Remover of sensitive History, not a Promise that Exposure never Happened.
It Detects without Echoing, Names the affected Surfaces,
Confirms the destructive Boundary, then Purges and Verifies.

The first Action after a Credential Exposure is Revocation or Rotation.
History Rewriting Reduces Distribution; it does not make a leaked Secret Safe.

## The Value

Never Ask the User to paste a sensitive Value into Chat.
Never place it in a command Argument, generated Script, Log or Report.

Ask the User to enter it directly into an interactive Terminal with Echo
disabled and keep it in a temporary environment Variable such as
`PURGE_VALUE`. The Agent Uses only the Variable Name.

- Never print, expand, inspect or persist the Variable.
- Never enable shell Tracing while it Exists.
- Reject an empty Value before every Search or Mutation.
- Clear it from the Environment when Verification Finishes.
- If the Client cannot accept secret Terminal Input without exposing it to the
  Model, Stop and give the User a local Command to run themselves.

## Detect

Read three Surfaces. Detection Changes nothing.

1. **Files** — tracked, untracked, hidden and ignored Files in Scope.
2. **Git** — reachable Commits, Branches, Tags and local Reflogs.
3. **Shell** — known local History Files for Bash, Zsh, Fish and the active
   Shell, plus project-local Session Logs explicitly named by the User.

Use fixed-string matching by Default. Treat the Value as Data, never a Regular
Expression. Exclude binary payload output and report only redacted Paths,
Commit IDs, Ref Names and Counts.

Do not search unrelated home Directories merely because they are Accessible.
Name every Path outside the current Repository and Ask before Reading it.

Report Detection separately:

```text
Files: 2 Matches in 2 Paths
Git: 4 Commits across 1 Branch and 1 Tag
Shell: 3 Entries in Fish History
```

Never include matching Lines or the sensitive Value.

## Confirm

Detection does not Authorize Mutation.

Before Purging, Explain:

- the exact Files, Refs and History Stores that will Change;
- whether Git Commit IDs will Change;
- whether Tags or Branches require Replacement;
- whether a remote Force Push will be needed;
- that existing Clones, Forks, Caches and Logs may retain the old Value;
- which Backups will be Created and when they will be Removed.

Require explicit Confirmation for each destructive Boundary:

1. Rewrite repository History.
2. Rewrite each shell History File.
3. Expire Reflogs and prune unreachable Git Objects.
4. Force Push rewritten Branches or Tags.

Never combine those Approvals. Never Infer Consent from the original Request.

## Purge

Purge only the Surfaces the User Confirmed.

### Files

Replace or remove the Value from the current working Tree first.
Preserve File Structure and unrelated Content.
Do not Commit unless the User separately Requests a Commit.

### Git

Prefer `git-filter-repo` when installed and supported by the Repository.
Read its installed Help before constructing the Rewrite.
Use a replacement File or Callback that Reads the Value without placing it in
the command Line, process List or persistent project Files.

Before rewriting:

- require a clean or fully understood working Tree;
- record the current Branches, Tags, Remotes and `HEAD` without Secrets;
- create a local recovery Bundle outside the Repository with restrictive
  Permissions;
- name the Bundle Path to the User;
- refuse to overwrite an existing Bundle.

Do not use `filter-branch` when `git-filter-repo` is available.
Do not delete original Refs, expire Reflogs or run Garbage Collection until the
rewritten History passes Verification and the User confirms final Pruning.

Never Force Push automatically. Show the affected remote Refs and ask for a
separate Approval. Use `--force-with-lease`, never unconditional `--force`,
when the remote State permits it.

### Shell

Stop or account for active Shells that may rewrite History on Exit.
Create a permission-restricted Backup beside neither the Repository nor its
tracked Files. Parse the native History Format; do not treat structured Fish or
multiline History as plain independent Lines when that would corrupt Entries.

Write a replacement File atomically, preserve Permissions and Ownership,
then ask the User to restart or reload affected Shell Sessions.
Never clear an entire History when exact Entry removal is possible.

## Verify

Repeat the same fixed-string Detection against every purged Surface.

1. Current Files contain zero Matches.
2. Rewritten reachable Git History contains zero Matches.
3. Confirmed shell History Stores contain zero Matches.
4. Refs and Repository integrity checks pass.
5. The sensitive Variable is cleared.

A zero Match in rewritten Git does not prove remote Caches or existing Clones
forgot the Value. Say so.

Only after Verification may the User separately Approve deleting recovery
Refs, expiring Reflogs, pruning Objects and removing Backups.

## Boundaries

- Never reveal a sensitive Value back to the User.
- Never send it to a network Service or external Scanner.
- Never mutate global credential Stores without an explicit Request.
- Never rewrite signed Commits or Tags without naming that signatures Break.
- Never rewrite a shared Branch without naming the Coordination required.
- Never remove Audit Evidence required by Law, Policy or an active Incident.
- Stop when repository ownership, remote authority or History format is unclear.

## What it Returns

Report States, Counts and remaining Actions only:

```text
✅ Purge Verified
- Files: 0 Matches
- Git: 0 Matches in rewritten reachable History
- Shell: 0 Matches in 1 confirmed History Store
- Credential: Rotation still Required
- Remote: Force Push not Performed
- Recovery Bundle: Retained until final Approval
```

The Report Never Contains the Value or a recoverable Fragment of it.
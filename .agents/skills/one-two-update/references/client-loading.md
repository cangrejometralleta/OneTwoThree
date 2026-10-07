# Client Loading

A Loader of local Customizations, not an installer of global State.
It Finds the canonical skills and agents in the current project,
connects the active client to them, removes obsolete entrances and references,
reloads what the client can reload, and verifies what it can discover.

The Canon Lives in [Vendor Integration](../../../../rules/vendor-integration.md).

## What it Reads

Start at the repository Root. Read only the surfaces needed to identify:

1. **Source** — canonical project Customizations under `.agents/`.
2. **Entrances** — existing client Directories, links and adapters.
3. **History** — renamed Files, stale identifiers and duplicate entrances.
4. **Client** — the current Host and the reload actions it actually exposes.

Inspect `AGENTS.md`, client configuration and local documentation when they
Name a different canonical source. Never assume that every client Supports
the same file format, directory or reload command.

## What it Does

1. Locate every canonical `SKILL.md` and custom Agent in the project.
2. Identify the current Client from the available tools and environment.
3. Find that client's documented project-level discovery Directories.
4. Reuse an existing relative symbolic Link when the client supports it.
5. Create the smallest required Adapter when the client needs another format.
6. Search the whole Project for superseded names, paths and adapters.
7. Confirm every stale Reference has a valid canonical replacement.
8. Remove obsolete Files and links; update references to the replacement.
9. Validate Links, frontmatter, names and adapter syntax before reloading.
10. Refresh the active client through its available command, API or UI control.
11. Verify Discovery through the client when an inspection tool exists.
12. Report what Loaded, removed, stayed unavailable or still needs a client action.

Treat these as common Entrances, not timeless guarantees:

| Client | Skills | Agents |
| --- | --- | --- |
| GitHub Copilot | `.github/skills/` | `.github/agents/` |
| Claude | `.claude/skills/` | `.claude/agents/` |
| Codex | `.agents/skills/` | `.agents/agents/` or a local Adapter |
| OpenCode | Discover from its local configuration or Documentation | Discover from its local configuration or Documentation |

If the installed client documents another Path, follow the installed client.
The Client Owns Discovery; this skill owns the local connection to it.

## Refresh the Active Client

Identify the client from the active app/session and available tools, not from
the folder layout alone. After changing files, use the client's supported
refresh action through its CLI, API or connected UI. Prefer an in-session
action. Do this yourself when the active client is reachable; do not stop at
giving the user a command.

Use these documented behaviors as starting points. Check current official docs
when the installed version exposes a different command or interface:

| Client | Refresh the current client | Verify |
| --- | --- | --- |
| [Claude Code](https://code.claude.com/docs/en/skills) | When the active entrance is project `.claude/skills/`, edits are watched. If that top-level directory was created after startup, run `/reload-skills`. | Run `/skills`; invoke the changed skill again to load its current body. |
| [Codex app, CLI or IDE](https://developers.openai.com/codex/skills) | Codex detects local skill changes. If the active client does not pick up the change, save the work and restart Codex. | Check the current skill picker/list and invoke the skill again. |
| [OpenCode](https://docs.opencode.ai/docs/skills/) | Invoke the native `skill` tool again; each model step advertises available skills and the tool resolves the current file. | Confirm the loaded body matches the canonical file. |
| [GitHub Copilot CLI](https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/add-skills) | Run `/skills reload` in the active session. | Run `/skills info <name>` or `/skills list`. |
| [GitHub Copilot in VS Code](https://code.visualstudio.com/docs/agent-customization/agent-skills) | No skill-specific refresh command is documented. Check the `/` skill menu; if stale, save the work and run `Developer: Reload Window`. | Reopen the skill menu and inspect the [Agent Customizations](https://code.visualstudio.com/docs/agent-customization/overview) view. |

For another client, find its currently documented refresh control instead of
assuming a universal reload command. Prefer the narrowest action that refreshes
the active session. Do not force-quit a client or delete its cache.

A refreshed skill registry does not erase instructions already injected into
the current conversation. Re-invoke the changed skill when supported; if the
client keeps the old body or cannot refresh this session, report that exact
limit and the smallest remaining action. Never claim a client reload from
filesystem changes alone.

## Integration Rules

- Keep one canonical Source. Do not maintain independent vendor copies.
- Use relative Links where supported so renames and edits travel together.
- Remove old Names, broken links and superseded adapters after their canonical
  replacements resolve. Search hidden and ignored project Paths too.
- Delete a duplicate Entrance only when the active client discovers the same
  logical source through another supported entrance.
- Never remove a compatibility Entrance used by another client merely because
  the current client also discovers it. Report cross-client duplicates instead.
- Preserve existing user Changes and unrelated client configuration.
- Do not write outside the current Project unless the user explicitly asks.
- Do not install Extensions, plugins or global packages without permission.
- Do not invent an adapter Schema. Read a nearby working Adapter or the
  installed client's documentation first.
- If a client cannot consume the canonical Format, report the incompatibility
  before generating anything.

## Contained Distribution

Use [Source Distribution](source-distribution.md) for the installation layout
and dependency Selection.
For a clone, keep the reduced source in `.agents/canon` and expose selected
agents and skills through individual relative Links.
For a ZIP snapshot, keep the generated real files in `.agents/agents` and
`.agents/skills`, and read provenance from `.agents/distribution.json`.
Follow the existing layout; never replace snapshot files with clone Links.
Do not copy their source or place a second vendor tree outside `.agents`.
Existing local customizations Keep their names and contents.

For clone installations, resolve each linked file to its physical Location before following relative
references to the canon or another Skill.
Validate both the exposed entrance and the resolved source Dependencies.
Commands and project edits still Target the consuming repository.
Include this resolution rule in the client's loading Instructions.
ZIP references already Target their generated locations.

Create only the active client's required discovery Entrances.
Outside `.agents`, allow only the links or configuration that discovery Needs;
report those paths explicitly, and do not promise universal single-folder Loading.
Do not create adapters for clients the user has not Requested.
The install itself does not create a Handoff or other runtime artifacts.

## Installation Rules

An Install Writes more than links. These considerations Come from installs
that looked correct and still Broke something quiet.

### A Link is not a Directory

`/.agents/` Matches a real directory and never a symbolic Link.
The trailing slash is the whole Bug: the day `.agents` becomes a link,
that pattern Stops matching and the cache Starts being carried.

Drop the Slash for every entrance that is a link, and Check every
ignore file the project has, not only `.gitignore`:

```text
.gitignore       the Repository does not Carry the cache
.dockerignore    the Image does not Need it
.gcloudignore    the Deploy does not Upload it
```

A missed Line in a deploy ignore Ships the Canon and its `.git`
in every Build. Verify with `git check-ignore`, never by Reading the file.

### The Copy Forgets where it Came from

A sparse Clone Knows its Head; a flat copy Knows nothing.
A generated ZIP records its origin, commit and file inventory in
`.agents/distribution.json`; verify it before loading the Snapshot.
For a legacy copy without that manifest, `AGENTS.md` must Record the address
and exact commit it was taken from:

```text
https://github.com/cangrejometralleta/OneTwoThree.git   branch: main
5f3e6ec7e20723449f4a03e079ba240af2f1af76                2026-09-16
```

The ZIP manifest marks `snapshot: true`; name it as a snapshot when Reporting.
For a legacy copy, say in the same note that the copy is a Fork until the clone returns —
[The Head is the Canon](../../../../rules/the-head-is-the-canon.md)
Asks for the head and Accepts no older commit.

### Restore before you Reshape

An Install that Rewrites the root may Delete the two files the repository
actually carries, `AGENTS.md` and its `CLAUDE.md` link.
Read `git status` before Writing anything. A `D` on either one is Loss,
and the commit is the only Copy left.
Restore first, Reshape second, and never in one Step.

### The Minimum must still Work

Minimize the selected skills, then Close their dependencies recursively.
Keep supporting files required by those Skills.
Do not install every skill merely because it Exists upstream.

Keep the three indexes and all their rule, value and pattern Bodies.
Dove needs patterns even when no skill explicitly Links them.
Load indexes first and bodies on Demand.

Excluded resources Stay excluded under `.canonignore`.
Report workflows that require them as Unavailable locally;
never solve a missing resource by silently copying the excluded Tree.

## Verification

The filesystem Check Proves the Entrance:

1. Every link Resolves inside the project.
2. Every skill folder name Matches its frontmatter `name`.
3. Every agent name and adapter reference Resolves.
4. No stale reference Points to a renamed skill or agent.
5. No removed name Remains in hidden, ignored or generated project paths.

The client Check Proves the Load:

1. Query the client's customization Listing when available.
2. Confirm the expected Names appear exactly once per logical source.
3. Start no destructive or stateful Workflow merely to test discovery.

Filesystem success without a client check is `Validated`, not `Loaded`.

## What it Returns

Keep the Report short and separate proven states:

```text
✅ Loaded in GitHub Copilot
- Skills: de-la-case, one-two-update
- Agents: dove
- Entrance: .github/skills -> ../.agents/skills
- Removed: 3 stale References

⚠️ Codex skill is on disk; the current client has not confirmed discovery.

Refresh the active client and re-invoke the updated skill; start a new session
only when the client requires it.
```

Never say every client Loaded when only one client was available to verify.
Never finish after changes without reporting the refresh action and whether
the active client verified the updated customizations.

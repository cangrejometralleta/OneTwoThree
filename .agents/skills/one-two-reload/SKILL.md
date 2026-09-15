---
name: one-two-reload
description: "Load or refresh a local project's skills and custom agents in the current AI client, removing stale references after replacements resolve. Use when project customizations are missing, stale, renamed, duplicated, newly linked, or need reloading in GitHub Copilot, OpenCode, Claude, Codex, or another supported client."
---

# OneTwoReload

A loader of local customizations, not an installer of global State.
It Finds the canonical skills and agents in the current project,
connects the active client to them, removes obsolete entrances and references,
reloads what the client can reload, and verifies what it can discover.

The canon Lives in [Vendor Integration](../../../rules/vendor-integration.md).

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
10. Invoke the client's available reload, rescan or window-refresh Action.
11. Verify Discovery through the client when an inspection tool exists.
12. Report what Loaded, removed, stayed unavailable or needs a manual restart.

Treat these as common Entrances, not timeless guarantees:

| Client | Skills | Agents |
| --- | --- | --- |
| GitHub Copilot | `.github/skills/` | `.github/agents/` |
| Claude | `.claude/skills/` | `.claude/agents/` |
| Codex | `.agents/skills/` | `.agents/agents/` or a local Adapter |
| OpenCode | Discover from its local configuration or Documentation | Discover from its local configuration or Documentation |

If the installed client documents another Path, follow the installed client.
The client Owns discovery; this skill owns the local connection to it.

## Reload Rules

- Prefer a client API or editor Command that explicitly reloads customizations.
- Reload the editor Window only when no narrower supported action exists.
- Never kill an active client Process or delete its cache without permission.
- Never claim a Reload from filesystem changes alone.
- When no reload action is exposed, validate the Files and ask for the smallest
  manual action: start a new chat, reload the window or restart the client.
- A skill cannot Reload the turn already reading it. Verify the next discovery
  cycle and say when that Boundary applies.
- After any skill, agent, entrance or tool change, Remind the user to reload or
  restart the client and begin a new chat. The current turn may keep the old
  Customizations and tool grants even when the filesystem is already correct.

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

## Verification

The filesystem check Proves the entrance:

1. Every link Resolves inside the project.
2. Every skill folder name Matches its frontmatter `name`.
3. Every agent name and adapter reference Resolves.
4. No stale reference Points to a renamed skill or agent.
5. No removed name Remains in hidden, ignored or generated project paths.

The client check Proves the load:

1. Query the client's customization Listing when available.
2. Confirm the expected Names appear exactly once per logical source.
3. Start no destructive or stateful Workflow merely to test discovery.

Filesystem success without a client check is `Validated`, not `Loaded`.

## What it Returns

Keep the Report short and separate proven states:

```text
✅ Loaded in GitHub Copilot
- Skills: one-two-case, one-two-reload
- Agents: dove
- Entrance: .github/skills -> ../.agents/skills
- Removed: 3 stale References

⚠️ Codex entrance Validated; client restart Required.

Reload the Client and begin a new chat before using the changed customizations.
```

Never say every client Loaded when only one client was available to verify.
Never finish a Reload after changes without the client reload reminder.
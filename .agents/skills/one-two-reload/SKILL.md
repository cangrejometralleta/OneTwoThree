---
name: one-two-reload
description: "Load or refresh a local project's skills and custom agents in the current AI client, removing stale references after replacements resolve. Use when project customizations are missing, stale, renamed, duplicated, newly linked, or need reloading in GitHub Copilot, OpenCode, Claude, Codex, or another supported client."
---

# OneTwoReload

A Loader of local Customizations, not an Installer of global State.
It Finds the canonical Skills and Agents in the current Project,
Connects the active Client to them, Removes obsolete Entrances and References,
Reloads what the Client can Reload, and Verifies what it can Discover.

The canon lives in [Vendor Integration](../../../rules/vendor-integration.md).

## What it Reads

Start at the Repository Root. Read only the Surfaces needed to Identify:

1. **Source** — canonical project Customizations under `.agents/`.
2. **Entrances** — existing client Directories, Links and Adapters.
3. **History** — renamed Files, stale Identifiers and duplicate Entrances.
4. **Client** — the current Host and the Reload Actions it actually Exposes.

Inspect `AGENTS.md`, client configuration and local Documentation when they
Name a different canonical Source. Never Assume that every Client Supports
the same File Format, Directory or Reload Command.

## What it Does

1. Locate every canonical `SKILL.md` and custom Agent in the Project.
2. Identify the current Client from the available Tools and Environment.
3. Find that Client's documented project-level discovery Directories.
4. Reuse an existing relative symbolic Link when the Client Supports it.
5. Create the smallest required Adapter when the Client needs another Format.
6. Search the whole Project for superseded Names, Paths and Adapters.
7. Confirm every stale Reference has a valid canonical Replacement.
8. Remove obsolete Files and Links; update references to the Replacement.
9. Validate Links, frontmatter, Names and Adapter Syntax before Reloading.
10. Invoke the Client's available reload, rescan or window-refresh Action.
11. Verify Discovery through the Client when an inspection Tool Exists.
12. Report what Loaded, Removed, stayed Unavailable or needs a manual Restart.

Treat these as common Entrances, not timeless Guarantees:

| Client | Skills | Agents |
| --- | --- | --- |
| GitHub Copilot | `.github/skills/` | `.github/agents/` |
| Claude | `.claude/skills/` | `.claude/agents/` |
| Codex | `.agents/skills/` | `.agents/agents/` or a local Adapter |
| OpenCode | Discover from its local configuration or Documentation | Discover from its local configuration or Documentation |

If the installed Client Documents another Path, Follow the installed Client.
The Client Owns Discovery; this Skill Owns the local Connection to it.

## Reload Rules

- Prefer a client API or editor Command that Explicitly Reloads Customizations.
- Reload the editor Window only when no narrower supported Action Exists.
- Never Kill an active Client Process or delete its Cache without Permission.
- Never Claim a Reload from filesystem Changes alone.
- When no reload Action is exposed, validate the Files and Ask for the smallest
  manual Action: start a new Chat, reload the Window or restart the Client.
- A Skill cannot reload the Turn already Reading it. Verify the next discovery
  cycle and Say when that Boundary Applies.

## Integration Rules

- Keep one canonical Source. Do not Maintain independent vendor Copies.
- Use relative Links where supported so Renames and Edits travel together.
- Remove old Names, broken Links and superseded Adapters after their canonical
  Replacements Resolve. Search hidden and ignored project Paths too.
- Delete a duplicate Entrance only when the active Client discovers the same
  logical Source through another supported Entrance.
- Never remove a compatibility Entrance used by another Client merely because
  the current Client also discovers it. Report cross-client duplicates instead.
- Preserve existing user Changes and unrelated client Configuration.
- Do not write outside the current Project unless the User explicitly Asks.
- Do not install Extensions, Plugins or global Packages without Permission.
- Do not invent an Adapter Schema. Read a nearby working Adapter or the
  installed Client's Documentation first.
- If a Client cannot consume the canonical Format, report the incompatibility
  before generating anything.

## Verification

The Filesystem Check Proves the Entrance:

1. Every Link Resolves inside the Project.
2. Every Skill Folder Name Matches its frontmatter `name`.
3. Every Agent Name and Adapter Reference Resolves.
4. No stale Reference points to a renamed Skill or Agent.
5. No removed Name remains in hidden, ignored or generated project Paths.

The Client Check Proves the Load:

1. Query the Client's customization listing when available.
2. Confirm the expected Names appear exactly once per logical Source.
3. Start no destructive or stateful Workflow merely to test Discovery.

Filesystem Success without a Client Check is `Validated`, not `Loaded`.

## What it Returns

Keep the Report short and Separate Proven States:

```text
✅ Loaded in GitHub Copilot
- Skills: one-two-case, one-two-reload
- Agents: dove
- Entrance: .github/skills -> ../.agents/skills
- Removed: 3 stale References

⚠️ Codex entrance Validated; Client restart Required.
```

Never Say every Client Loaded when only one Client was available to Verify.
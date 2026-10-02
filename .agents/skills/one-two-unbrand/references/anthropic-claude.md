# Anthropic and Claude

## History marks

**Identities**
- Author or Committer named `Claude`, with an address at `anthropic.com`.

**Trailers**
- `Co-Authored-By: Claude <model> <noreply@anthropic.com>`, whatever the model.
- `Claude-Session: <url>`.
- A `Generated with` line, with or without the robot emoji.

**Detect** — counts only, over every ref:

```bash
git log --all -i -E --grep='co-authored-by: claude|claude-session|generated with' --format=%H | wc -l
git log --all --format='%an|%cn|%ae|%ce' | grep -ci anthropic
```

**Mailmap** — the project's own identity, read from `git config`, never typed in:

```text
<user.name> <user.email> Claude <noreply@anthropic.com>
```

**Callback** — drop the three Trailer kinds and the blank lines they leave:

```python
import re
lines = message.split(b"\n")
keep = [l for l in lines
        if not re.match(rb"(?i)^co-authored-by: claude\b.*<noreply@anthropic\.com>\s*$", l)
        and not re.match(rb"(?i)^claude-session:", l)
        and not re.match(rb"(?i)^.{0,4}generated with", l)]
return re.sub(rb"\n+$", b"", b"\n".join(keep)) + b"\n"
```

## Tree marks

**Entrances**
- `CLAUDE.md` -> `AGENTS.md`
- `.claude/agents` -> `../.agents/agents`
- `.claude/skills` -> `../.agents/skills`

**Settings** — `.claude/settings.json` and `.claude/settings.local.json`.
They are the user's: permissions, hooks and the attribution Blanks. Never touched.

## Leave alone

- A human Co-author, and the identities of editors such as Zed or GitHub.
- A vendor Entrance or an Adapter, where the Name is the subject.

## Keep it out

`attribution.commit` and `attribution.pr` set to `""`
in `.claude/settings.json` stop the next Trailer at the Source.

# Google and Gemini

The Paths below come from the Gemini CLI documentation.
Verify Discovery in the tool before claiming a Fix works.

## History marks

None Established. No Gemini identity or trailer has been Seen in this repository,
and none is guessed: a pattern such as `google` would Match humans.

If the user Names one, add it here with its Detect, mailmap and Callback,
in the shape of [anthropic-claude.md](anthropic-claude.md).
Until then, the History Surface reports `no marks known` and stops.

## Tree marks

**Entrances**
- `GEMINI.md` -> `AGENTS.md`, a symbolic link, as `CLAUDE.md` is.
  Present in this repository.
- `.gemini/agents` -> `../.agents/agents`. Present in this repository.

**Not an Entrance** — `.gemini/skills`.
The client Discovers skills in `.agents/skills/` where they already Live,
so a link there Repeats what it Finds. Count it as Clutter, never as a Gap.

**Settings** — `.gemini/settings.json`, and the user-level `~/.gemini/`.
They are the user's. Never touched, and `~/.gemini/` is Outside the repository.

## Leave alone

- A `.gemini/commands` folder: it Holds the user's own slash commands.
- The Format difference: a Gemini agent file has its own Frontmatter.
  A link Shares the file and does not convert it, so `.agents/agents/dove.md`
  may need a small Adapter before Gemini reads it (`Inferred`).
  Report it; never rewrite the shared markdown for one vendor.

## Keep it out

Not Established. Check the Gemini CLI settings for an Attribution switch
before Claiming one, and Say so if there is none.

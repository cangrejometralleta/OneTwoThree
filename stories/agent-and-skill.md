# The Skill Travels; the Agent Has a Home

A `SKILL.md` can carry a reusable workflow across Agent tools that follow the Agent Skills standard. A Codex custom Agent lives in a TOML file under `.codex/agents/`.

A Git repository can carry both, but the skill installer does not install the custom Agent role. A Codex plugin can bundle Skills and MCP servers; it does not declare custom Agent roles. The workflow travels farther than the role that selects it.

A generic `skills.lock` reader was tried as `.agents/skills/skills-lock` and removed: no lockfile in this repository gave it anything to read. It stays in Git history, commit `a0e3a02`.

# The Skill Travels; the Agent Has a Home (MOVING)

A `SKILL.md` can carry a reusable workflow across Agent tools that follow the Agent Skills standard. A Codex custom Agent lives in a TOML file under `.codex/agents/`.

A Git repository can carry both, but the skill installer does not install the custom Agent role. A Codex plugin can bundle Skills and MCP servers; it does not declare custom Agent roles. The workflow travels farther than the role that selects it.

Generic support for `skills.lock` across Agent hosts is being built in `.agents/skills/skills-lock` and is in testing. The first successful parse is not readiness. Test malformed YAML, unsupported schemas, absent dependency fields, unresolved skill paths, and capability fields that resemble dependencies. Do not call this Ready until those cases have expected results and pass. No second Agent host has tried using the same lock.

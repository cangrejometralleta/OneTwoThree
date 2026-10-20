---
name: skills-lock
description: "Inspect a skills.lock or skill manifest and report each skill's explicitly declared dependencies. Use when someone asks what skills depend on, which tools they require, or whether a lockfile contains dependency data."
---

# SkillsLock

Read a skill lock or manifest and report the dependencies it declares, with their source and any unresolved references.

## Read the Format

- Use the file the User names; otherwise find `skills.lock` at the repository root.
- Parse YAML with an available YAML parser. Read the schema version and actual field names before extracting anything. Do not treat an unfamiliar shape as a known schema.
- If the file is missing, malformed, or unsupported, report that and stop. Do not guess a dependency format.

## Resolve Dependencies

- Extract only fields the format defines as dependencies. Keep each dependency's type, identifier, version or constraint, and source path when present.
- For an OpenAI skill entry with `source_path`, resolve the skill directory and read its `agents/openai.yaml`. Report `dependencies.tools` entries, including their type and value.
- If a lock entry has no dependency field, inspect the referenced skill's supported metadata for explicit dependencies. Do not infer dependencies from instructions, command examples, imports, or filenames.
- Keep permissions and capabilities separate from dependencies. For example, `behavior.allowed_tools` describes allowed tools, and `bundled_scripts` lists files; neither alone declares a dependency.
- If a requested dependency kind is absent from the schema, say it is not recorded. Do not claim that the skill has no dependencies.
- Do not install dependencies, fetch remote resources, or modify the lockfile.

## Report

Summarize the detected format and version, then list dependencies by skill. Include dependency type, identifier, version or constraint, and where it was declared. Mark unresolved paths or unsupported fields. Distinguish declared dependencies from any inference the User explicitly requested.

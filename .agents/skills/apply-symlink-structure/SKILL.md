---
name: apply-symlink-structure
description: Apply a documented symlink layout to a folder or project while preserving one canonical source and validating every link. Use when someone asks to connect project folders or shared files through symlinks.
---

# ApplySymlinkStructure

Apply the requested or documented link layout. Keep one canonical source and expose it through the smallest set of useful Entrances.

The governing guidance here is [Vendor Integration](../../../rules/vendor-integration.md) and [The Path that Travels](../../../patterns/the-path-that-travels.md).

## Choose the Layout

Read the target project's instructions and documentation for the canonical source, link destinations, supported clients and existing distribution method. In OneTwoThree, follow the two governing documents above. If the user supplies a source and destination map, use it. If the source or relationship is unclear, ask one focused question instead of guessing.

Use [OneTwoUpdate](../one-two-update/SKILL.md) when the task is specifically about adopting or syncing the OneTwoThree Canon, or connecting project customizations to a client. It routes to the relevant source-distribution or client-loading workflow.

## Apply the Links

Inspect the affected paths and current project state before writing. Leave a correct existing link alone. Create only missing links in the requested layout, using relative targets calculated from each link's parent so the layout can move together.

Keep shared content in its canonical location; do not make a second maintained copy. Link only the files or shared subdirectories the layout calls for. Keep project settings and secrets outside shared links. Add a documented adapter only when a destination requires one, and mark generated content as generated.

Do not overwrite, replace or remove an existing file, directory or link to resolve a collision. Report the collision and ask before changing that path. If the filesystem or target client cannot use symlinks, report the limit and follow a documented alternative; do not silently copy the source.

## Verify and Report

Resolve every created link and confirm its target exists and matches the intended canonical source. Check that the project will carry the link, including ignore rules when relevant. Review the final changed paths for unexpected edits.

Report each link as `link path -> relative target`, along with any collision or unsupported destination. Filesystem validation proves that links resolve; it does not prove that a client discovered or loaded them. Report client discovery only when separately verified.

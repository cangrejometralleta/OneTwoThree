---
name: what-i-did-today
description: Create a concise first-person Markdown Jira comment that leads with completed work and reports active blockers or states when none were identified. Use when someone asks what they did today or needs help briefing recent session activity.
---

# WhatIDidToday

Turn the current work session into a brief Jira comment the user can paste as-is. Lead with what was completed, then state blocker status.

## Build the Recap

Use the conversation and visible evidence: edits, commands, tool results, decisions and outcomes. Inspect current files or Git state only when they clarify what happened. Distinguish completed work from plans. Never invent activity, imply unfinished work is complete, or claim validation that did not happen. Leave out conversation detours, internal reasoning and unrelated pre-existing changes.

Write from the user's point of view in first person, using active verbs such as “I updated” or “I traced.” Start with the most important completed work and its result. Include validation only when it happened. Keep plans and remaining work separate from completed work so they do not overshadow it.

Use a short Markdown structure by default:

- **Done** — the completed work and its outcome.
- **Blockers** — each active blocker and what it prevents or waits on. If no blocker appears in the session evidence, say “No blockers were identified in this session.” If the status is genuinely unclear, say so instead of guessing.

Keep the recap brief, with fewer details for a small session. Use ordinary capitalization in the Jira comment.

Return only the Markdown comment body, with no code fence, preamble or sign-off. Do not post the comment or interact with Jira; this skill prepares text for the user to paste.

If the session does not provide enough evidence to describe completed work truthfully, ask one focused question or state the specific detail that is missing. Do not fill gaps with guesses.

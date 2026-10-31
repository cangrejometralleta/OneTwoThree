# Transfer

Prepare a receiver `.handoff.md` that the user can copy into another repository.
Do not require the destination path or checkout. The receiver owns the active
`Next` after the handoff is copied there and opened.

The Canon lives in [Paired Repositories](../../../../rules/paired-repositories.md)
and [Session Checkpoint](../../../../rules/session-checkpoint.md).
The Checkpoint shape lives in [RememberTheTime](../../remember-the-time/SKILL.md).

## Prepare a Copy-Ready Handoff

By default, make a portable handoff in the response. Read the current Thread
and source repository state, then produce a complete `# Checkpoint` in one
Markdown code block, ready to paste as `.handoff.md` at the destination repo's
root. The user should not need to provide or recall a path.

Include the active Intent, only the Done/Open/Blockers needed to resume,
relevant source paths and branch/commit, any essential log reference, and one
concrete `Next` in the destination. Mark destination branch and state as
`Not inspected` instead of guessing. Set `Transfer` to `From <source repo>`
and identify the handoff as prepared for manual copy. Use the repository name
from local Git metadata when available; otherwise use the known repo name.

Update the source `.handoff.md` to say the receiver checkpoint is prepared for
manual copy, name the destination if known (otherwise `Not specified`), and
make its `Next` the copy-and-resume action. This is a pending transfer until the
user places the block at the destination root.

Tell the user to copy the entire block to the destination repo's root as
`.handoff.md`, then invoke `yo-yo-yo` there. If the destination already has a
different active handoff, do not tell the user to overwrite it; report the
collision and preserve both threads.

Keep the handoff concise. Treat `Intent` as the root of a small context graph:
connect only the `Done`/`State` evidence, `Open`/`Blockers`, relevant repository
or log references, and one `Next` needed to resume. Do not reproduce the whole
conversation or expose extra branches as options. To return, create the same
copy-ready handoff in the opposite direction.

## What it Returns

Return the copy-ready `.handoff.md` as one block and state its target filename
and next invocation: paste it at the destination root, then use `yo-yo-yo`.

If the user identifies an existing active handoff at the destination, report
the collision and preserve both threads; never ask them to overwrite it.

## Bounds

- Never require the user to recall a destination path.
- Never claim the transfer is complete before the user copies the block.
- Never overwrite a different active receiver handoff.
- Never create a second active `Next` in the source.
- Never close the session to transfer the Thread.
- Never copy full logs, terminal history, secrets, tokens or conversation text.
- Never commit, push or stage the source checkpoint.

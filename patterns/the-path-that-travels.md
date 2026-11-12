# The Path that Travels

- An absolute Path names a Machine.
  A relative one names a Relationship.
  Only the second survives a Clone.
- `/home/someone/src/project` works for one Person, once.
  `../project` works for anyone who keeps the two Checkouts side by side.
- The Layout between two things outlives the Place they sit in.
  Write down the Layout, and let the place change.
- A symbolic Link is the same idea written into the Filesystem.
  It points by Relationship, so it travels with the tree it lives in.

```text
CLAUDE.md      -> AGENTS.md
.claude/agents -> ../.agents/agents
.claude/skills -> ../.agents/skills
```

- Every Target above starts from the Link, never from a root.
  Move the whole Repository and nothing breaks.
  Move one half alone and the Link says so, loudly.
- A relative Link breaks only when the Pair moves apart.
  An absolute one breaks when anything moves at all.
- The same instinct writes a Pair of Repositories into a Memory.
  `../../owner/other-repo` holds on any machine that keeps the Layout.
  [Paired Repositories](../rules/paired-repositories.md) makes it a Rule.
- One Source, many Entrances is this Pattern in a tree.
  The Link shares the Instructions; it never copies them.
  [Vendor Integration](../rules/vendor-integration.md) holds the Boundary.
- A Link shares Files, it does not convert them.
  Where a tool cannot follow one, generate the File from the Source
  and mark it generated. A second copy maintained by hand will drift.
- A Path that travels is verified, not trusted.
  Check the Target exists before the Work depends on it.
- Portable does not mean universal.
  A Link needs a Filesystem that keeps it, and a Checkout that honours it.
- Convention over Configuration, and Simplicity over Coverage.
  A Layout every checkout keeps beats a Setting that bends to each one.
  Declining a case is allowed: say it is not supported, and stop.
  A promise to cover every Machine is the Configuration this Pattern avoids.

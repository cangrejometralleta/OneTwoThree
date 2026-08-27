# Color by semantic association

Could Entities, Actions and Statuses, already distinguished
by OneTwoThreeCase capitalization, also carry distinct colors when
rendered? This one already has a
bridge to code — `pdf/colors.go` renders body text in one flat
`ColorBody`, nothing splits it by grammatical role yet.
First instinct was red/green/yellow, primary colors chosen for being
intuitively combinable — the same intuition behind quark color charge
having exactly three (though the real physics term is red, green and
blue, not yellow — worth getting right if this ever becomes a named
ancestor). Objection: Rules' own Emoji section already assigned
that trio a meaning — ✅ green Passed, ❌ red Failed, ⚠️ yellow Careful
— so reusing it for grammatical role would collide with a meaning the
reader already learned in the same document, and would clash against
`pdf/colors.go`'s warm, desaturated palette. Primary hues may still be
worth rescuing later for something that isn't already spoken for.
Could also seed a lint rule alongside **Own linter**, but tagging a
capitalized Word in prose as Entity vs. Action vs. Status needs
judgment no mechanical rule has yet — the same wall the linter idea
hits everywhere it touches prose.
A humbler, tractable version: apply it only to the Go code inside the
fenced examples, not to prose. A type or struct is unambiguously an
Entity, a function or method call is unambiguously an Action, and Go's
own AST already answers the question — no judgment call, just a
parse, the same move **Own linter** already makes on Rules. Nobody
else seems to color example code by domain role instead of by syntax
token (keyword, string, comment) — narrow niche, but a real one, and
`pdf/colors.go` + `drawCodeBlock` already render every code block in
one flat color, so there's nothing to unlearn.
How to apply any of this is still open.
Not yet a confirmed root: one lived instance.

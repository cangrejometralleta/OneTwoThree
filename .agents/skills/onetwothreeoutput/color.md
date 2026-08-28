# Color, when Asked

The optional Door of [OneTwoThreeOutput](SKILL.md).
Read this File only when the User Asks for Color.
Off by Default means Off — if nobody Asked, this File
never Opens, and the Output Stays Plain.

Sleeps again where the Terminal Carries no Color.
Markdown Rendered as Prose Carries none; a Program
Printing to a TTY does. Know which Plane you are on.

The canon is [Values](../../../VALUES.md) — Color is an Index,
not Decoration. A Hue Groups what it Marks.

## Prune first

Prune before you Color, never after. Color Organizes;
it never Prunes. Color first and it Hides the Excess
instead of Showing it — five tidy Items in three Hues
are still five Shouting, now in three Tongues.

## The Rules

1. **Group by Co-occurrence** — same Token, same Hue.
   Two Entities that Share a Line, or Share an Identifier,
   the Reader already Related before the Color Arrived.
   The Color Confirms what the Eye did alone.

2. **Assign by Order of Appearance** — never by Meaning.
   A Taxonomy Needs a fixed List of Themes, and the Theme
   Changes with every Output. A Key that Fails half the Time
   Costs more than no Color: the Reader learns to Distrust it.
   Co-occurrence Asks for no Judgment, only for what Repeats.

3. **Three Hues, never a fourth** — the fourth Group Stays bare.
   The Ceiling makes Color a Limiter, never an Ornament.
   Three Marked and the rest Plain is the Sign
   that the Output is Inflated.

4. **Two Members make a Group** — one does not.
   A single Entity has no Pair to Rotate with,
   and there no Color Goes.

5. **Background Takes no Hue** — a Token on every Line Groups nothing.
   It Marks the Page, not a Part of it. Color what Gathers some Lines,
   never what Covers them all. One Line alone has no Background,
   so a Token Repeating inside it still Counts.

6. **One Line Set, one Hue** — Tokens that Touch exactly the same Lines
   are one Fact Wearing several Names. `.agents/skills/onetwothreeoutput/`
   is one Home, not three Groups. They Share a Hue, and the Hues they
   stopped Eating Stay free for what really Differs.

7. **Never the Emoji Trio** — Green, Red and Yellow are Spoken for.
   ✅ Passed, ❌ Failed, ⚠️ Careful. A Group that Borrows them
   Collides with a Meaning the Reader already Learned.

## The Palette

```
38;5;39   Blue      #00afff
38;5;170  Magenta   #d75fd7
38;5;80   Cyan      #5fd7d7
```

256-Color, never Truecolor — it Renders everywhere,
and it Survives a light Terminal as well as a dark one.

Rule seven Takes the whole warm End off the Table,
so the three Hues Come from the cool Half. That is a
Constraint, not a Taste. The Bonus: no Pair among them
Fails Red-Green Deficiency, because the Reserve
already Removed the Colliding Pair.

Blue and Cyan Sit close, and a washed-out Terminal
can Blur them into one. Swap Cyan for `38;5;208` Orange
if it Happens — but Orange Neighbors the ⚠️ Yellow.
Pick your Collision.

The Paper Palette in `examples/pdf/colors.go` Runs warm
and Shares nothing with these. That is Fine.
Paper Knows its Background; the Terminal does not.

## The Markdown Plane

ANSI never Reaches a Rendered Document, and HTML Gets Stripped.
Markdown Carries three Marks instead of three Hues:

```
**bold**   *italic*   `code`
```

`code` is the best of them, not a Fallback — the Tokens being
Grouped are Identifiers, so the Mark Means what it always Meant.

One Rule Comes with it. An Identifier outside a Group Stays plain.
Backtick every Path out of Habit and the third Mark Stops
Meaning Group. The Mark is Spent on Grouping or on nothing.

### The Ceiling is what the Surface has Left

Three is the Terminal's Number, because a Terminal Starts Empty.
A Document does not. It Spends Bold on a Label and Backticks on
a Path before Grouping ever Arrives, so the Budget is three minus
what the Page already Owes.

A Mark Carrying Emphasis cannot also Mean Group. Lay one on the
other and both Collapse — the Reader Meets a Bold Word and cannot
Tell whether it Matters or merely Belongs.

So Count what the Page Spends, then Group with the Rest.
A Page that Bolds its Labels Groups with two. One that Bolds and
Backticks Groups with one. A Page that Spends all three Groups
with none, and that is the correct Answer, not a Failure.

The Script Reads the Text and Counts for you. `--marks=N` Caps it
lower when you Know the Page will Spend more than it Shows.

### Why the Lines do not Move

Markdown can do what a Terminal cannot: Gather Lines. A Blockquote
Groups before any Mark Does, and it Spends no Mark to do it.
It was Built, Tested, and Cut. The Reason is worth Keeping.

Moving a Line **is** its Mark, so the Tokens that Earned the Group
Carry nothing. The Indent Says *these Belong* and never Says
*by what*. The Reader Sees a Block and Asks why — the exact
Question the Proof Warns about.

And a Blockquote already Means Quotation or Aside. Borrowing it
Rebuilds the Trap that Barred Bold and the Emoji Trio.

The deeper Reason: Color Adds a Channel. Indentation Spends one
Markdown already Uses. The Terminal has a free Dimension;
a rendered Document does not.

If a Surface ever Appears with a spare structural Channel —
Columns, a Gutter, a Margin — the Idea Returns. The Test it Needs
is mechanical and Written down: Groups that Move must Miss each
other and Run without a Gap. The Rest Stay put and take a Mark.

## The Script Does the Work

The Rules above are mechanical, so a Program Holds them
better than a Model does. Do not Color by Hand — Pipe it:

```
your-command | python3 color.py
```

It Reads stdin, Finds the Tokens that Repeat, Hands the first
three a Hue, and Writes stdout. Python 3, no Dependencies.

- `--legend` — Name the Key on stderr, so nobody Guesses it.
- `--force` — Paint into a Pipe, for `less -R` or a Test.
- Bare Arguments Choose the Tokens by Hand, when the Human
  Sees a Group the Count missed. That is the Selects Half.

It Passes Through untouched when stdout is not a Terminal,
so a Pipe never Eats an Escape. `NO_COLOR` Beats `--force`:
the Reader's standing Preference Outranks the Writer's Flag.

A Slash Separates Tokens; a Dot, Dash, Colon or Underscore Binds them.
So `shape_test.go:41` Survives whole, while `.agents/skills/onetwothreecase`
Splits into three — and the last of them can Rhyme with the same Name
Standing alone on another Line. That Rhyme is the Co-occurrence.

## The Proof is a Day

Color one real Output by Hand, then Look at it Tomorrow.
Found the Line faster — it Stays.
Caught yourself Asking why this one is Blue —
the Color Costs more than it Pays.

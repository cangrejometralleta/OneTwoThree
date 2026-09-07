# The Converter

Markdown in, a laid-out PDF out. No pandoc, no weasyprint —
goldmark Reads the AST, gopdf Draws the Page, both pure Go.

```sh
./build.sh                    # the Gates, then the Binary
./run.sh testdata/sample.md   # writes sample.pdf beside it
```

## The Shape

```text
main.go         Casts the Players, then Steps off the Stage
document/       the Domain: Document, Section, Block, and five Rules
style/          the Page, the Type Scale, the Colors; Imports nothing
markdown/       THE ONLY PACKAGE THAT IMPORTS GOLDMARK
render/         THE ONLY PACKAGE THAT IMPORTS GOPDF, and the Fonts
render/fonts/   Liberation Serif + DejaVu Sans Mono, vendored
```

`go list -deps ./document ./style` Names no Vendor.
The Core Cannot Import goldmark or gopdf, because it never Sees them.

Every Number the Page Depends on Carries a Name in `style/`.
Read `10.5` alone and it Says nothing; read `SizeEpigraph`
and it Says where the Number Lands.

## Five Rules Turn a Document into a Booklet

`ExtractCoverBlock`, `SplitTitleIndex`, `BuildTriadBlock`,
`BuildCalloutBlock`, `MarkClosingParagraph` — one Function per Shape a
plain Markdown Element can Take. The Title, the first Paragraph and the
first Blockquote Become the Cover. A Heading Split on `·` Gains an Index.
A single-row Table Becomes three Columns. A Blockquote that Opens on a
bold Word Becomes a labeled Callout. The last italic Paragraph Closes
the Book.

Read `document_test.go` against `testdata/sample.md` — every Rule Fires
at least once there, including the Case a Rule Refuses.

## Fonts Travel with the Binary

Liberation Serif Plays the Role Bitstream Charter Played in the old
style.css; DejaVu Sans Mono Keeps the Mono Face. Both Ship under free
Licenses (see `fonts/LIBERATION-LICENSE`, `fonts/DEJAVU-LICENSE`) and
are Embedded with `go:embed`, so the Tool Depends on no Font the Machine
Running it Happens to Have.

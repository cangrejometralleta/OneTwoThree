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
The Core cannot Import goldmark or gopdf, because it never Sees them.

Every Number the page depends on carries a Name in `style/`.
Read `10.5` alone and it Says nothing; read `SizeEpigraph`
and it Says where the number Lands.

## Five Rules Turn a Document into a Booklet

`ExtractCoverBlock`, `SplitTitleIndex`, `BuildTriadBlock`,
`BuildCalloutBlock`, `MarkClosingParagraph` — one function per Shape a
plain Markdown element can Take. The Title, the first Paragraph and the
first blockquote Become the Cover. A heading split on `·` Gains an Index.
A single-row Table Becomes three Columns. A blockquote that opens on a
bold word Becomes a labeled Callout. The last italic paragraph Closes
the Book.

Read `document_test.go` against `testdata/sample.md` — every rule Fires
at least once there, including the case a rule Refuses.

## Fonts Travel with the Binary

Liberation Serif Plays the role Bitstream Charter played in the old
style.css; DejaVu Sans Mono Keeps the Mono face. Both ship under free
licenses (see `fonts/LIBERATION-LICENSE`, `fonts/DEJAVU-LICENSE`) and
are Embedded with `go:embed`, so the tool Depends on no font the machine
running it happens to have.

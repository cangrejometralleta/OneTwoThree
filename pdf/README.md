# The Converter

Markdown in, a laid-out PDF out. No pandoc, no weasyprint —
goldmark Reads the AST, gopdf Draws the Page, both pure Go.

```
go run . "source.md"      # writes source.pdf beside it
go test ./...
```

## The Shape

```
main            Casts the Player, then Steps off the Stage
markdown.go     THE ONLY FILE THAT IMPORTS GOLDMARK
document.go     the Domain: Document, Section, Block, and five Rules
colors.go       the old style.css Values, Named instead of Numbered
fonts.go        THE ONLY FILE WITH EMBEDDED FONT DATA
render.go       THE ONLY FILE THAT IMPORTS GOPDF
fonts/*.ttf     Liberation Serif + DejaVu Sans Mono, vendored
```

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

// Package main Converts a Manifesto Written in Markdown into a
// laid-out PDF. goldmark Reads the Source into an AST, gopdf Draws
// the Page — no pandoc, no weasyprint, no Binary outside this one.
//
// This Comment is its own Proof. `go doc .` Prints the Text you are
// Reading now, an Editor Hovers it over the Package Name, and the
// Compiler Checks that everything it Describes still Exists. A
// File that Holds no Function still Compiles, still Ships, and
// still Explains — a Document, and Code, in the same Breath.
//
// # Usage
//
//	go run . source.md
//
// Writes source.pdf beside the Source.
//
// # The Shape
//
// main.go Casts the Players, then Steps off the Stage — the Script
// Rules Names under its own Script Section, Running here for
// real instead of Quoted as an Example.
//
// Four Packages Divide the Work, and the Compiler Holds the Line.
// markdown/ is the only Package that Imports goldmark. render/ is
// the only one that Imports gopdf. document/ Holds neither: it is
// the Business Truth in between, five small Functions Named after
// the Shape each One Produces — ExtractCoverBlock, SplitTitleIndex,
// BuildTriadBlock, BuildCalloutBlock, MarkClosingParagraph.
// style/ Holds the Page, the Type Scale and the Colors, and Imports
// nothing at all.
//
// `go list -deps ./document ./style` Names no Vendor. The Boundary
// is no longer a Promise a File Keeps; it is one the Build Enforces.
//
// # Why the old build.py Needed a Docstring and this Package Does not
//
// Python Read its Explanation only when a Human Opened the File.
// Go Reads this one from the Toolchain itself: `go doc .` Prints
// it, `gofmt` Refuses to Drop it, and a broken Reference inside it
// would still Compile, since Go Checks the Code, never the Prose
// Describing it. The Guarantee Stops there — which is why the
// Sentence above Names markdown/ and render/ directly, so a
// Rename Breaks a Grep before it Breaks a Reader's Trust.
package main

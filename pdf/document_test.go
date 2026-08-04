package main

import "testing"

func TestSplitTitleIndexSeparatesOnTheDot(t *testing.T) {
	index, title := SplitTitleIndex("One · First Rule")
	if index != "One" || title != "First Rule" {
		t.Fatalf("wanted %q/%q, got %q/%q", "One", "First Rule", index, title)
	}
}

func TestSplitTitleIndexKeepsAHeadingWithNoDot(t *testing.T) {
	index, title := SplitTitleIndex("Plain Heading")
	if index != "" || title != "Plain Heading" {
		t.Fatalf("wanted empty Index and the whole Title, got %q/%q", index, title)
	}
}

func TestBuildTriadBlockZipsHeadersAndCells(t *testing.T) {
	triad, ok := BuildTriadBlock([]string{"Receive", "Return"}, []string{"Read", "Answer"})
	if !ok || len(triad.Columns) != 2 || triad.Columns[0].Title != "Receive" {
		t.Fatalf("wanted a two-Column Triad, got %+v, ok=%v", triad, ok)
	}
}

func TestBuildTriadBlockRefusesAMismatch(t *testing.T) {
	if _, ok := BuildTriadBlock([]string{"One"}, []string{"a", "b"}); ok {
		t.Fatal("a length Mismatch Must Refuse the Triad")
	}
}

func TestBuildCalloutBlockRequiresALabel(t *testing.T) {
	if _, ok := BuildCalloutBlock("", []string{"body"}); ok {
		t.Fatal("an empty Label Must Refuse the Callout")
	}

	callout, ok := BuildCalloutBlock("Note", []string{"body"})
	if !ok || callout.Label != "Note" {
		t.Fatalf("wanted a Callout labeled Note, got %+v", callout)
	}
}

func TestMarkClosingParagraphMarksOnlyTheLastItalicOne(t *testing.T) {
	doc := Document{Sections: []Section{{Blocks: []Block{
		&Paragraph{Text: "first"},
		&Paragraph{Text: "second", Italic: true},
	}}}}

	MarkClosingParagraph(&doc)

	blocks := doc.Sections[0].Blocks
	if blocks[0].(*Paragraph).Closing {
		t.Error("the first Paragraph Must not Close")
	}
	if !blocks[1].(*Paragraph).Closing {
		t.Error("the last italic Paragraph Must Close")
	}
}

func TestMarkClosingParagraphIgnoresANonItalicLastParagraph(t *testing.T) {
	doc := Document{Sections: []Section{{Blocks: []Block{&Paragraph{Text: "plain"}}}}}

	MarkClosingParagraph(&doc)

	if doc.Sections[0].Blocks[0].(*Paragraph).Closing {
		t.Error("a non-italic last Paragraph Must not Close")
	}
}

// A Paragraph before the first Heading, a Table too wide for a Triad,
// and a Thematic Break buildBlock does not Recognize — three Edges
// the Sample never Reaches.
func TestParseManifestoDocumentHandlesEdgeCases(t *testing.T) {
	doc, err := ParseManifestoDocument("testdata/edge-cases.md")
	if err != nil {
		t.Fatalf("Edge Cases Must still Parse: %v", err)
	}

	if len(doc.Sections) != 1 {
		t.Fatalf("wanted one Section, the orphan Paragraph Joins none, got %d", len(doc.Sections))
	}

	blocks := doc.Sections[0].Blocks
	if len(blocks) != 1 {
		t.Fatalf("wanted one Block, the Thematic Break Produces none, got %d", len(blocks))
	}

	table, ok := blocks[0].(TableBlock)
	if !ok {
		t.Fatalf("wanted a fallback TableBlock, got %T", blocks[0])
	}
	if len(table.Headers) != 2 || len(table.Rows) != 2 {
		t.Errorf("wanted 2 Headers and 2 Rows, got %d Headers and %d Rows", len(table.Headers), len(table.Rows))
	}
}

func TestParseManifestoDocumentShapesTheSample(t *testing.T) {
	doc, err := ParseManifestoDocument("testdata/sample.md")
	if err != nil {
		t.Fatalf("the Sample Must Parse: %v", err)
	}

	if doc.Cover.Title != "Sample Manifesto" {
		t.Errorf("wanted the Title from the first H1, got %q", doc.Cover.Title)
	}
	if len(doc.Cover.Epigraph) == 0 {
		t.Error("the Epigraph Must not be Empty")
	}
	if len(doc.Sections) != 2 {
		t.Fatalf("wanted two Sections, got %d", len(doc.Sections))
	}
	if doc.Sections[0].Index != "One" {
		t.Errorf("wanted Index %q, got %q", "One", doc.Sections[0].Index)
	}

	var sawTriad, sawCallout, sawQuote, sawCode, sawClosing bool
	for _, section := range doc.Sections {
		for _, block := range section.Blocks {
			switch v := block.(type) {
			case Triad:
				sawTriad = true
			case Callout:
				sawCallout = len(v.Label) > 0
			case Quote:
				sawQuote = true
			case CodeBlock:
				sawCode = len(v.Text) > 0
			case *Paragraph:
				sawClosing = sawClosing || v.Closing
			}
		}
	}

	for name, got := range map[string]bool{
		"Triad": sawTriad, "Callout": sawCallout, "Quote": sawQuote,
		"CodeBlock": sawCode, "Closing Paragraph": sawClosing,
	} {
		if !got {
			t.Errorf("the Sample Must Produce a %s Block", name)
		}
	}
}

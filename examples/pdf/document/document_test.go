package document

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

package markdown

import (
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/pdf/document"
)

func TestParseManifestoDocumentHandlesEdgeCases(t *testing.T) {
	doc, err := ParseManifestoDocument("../testdata/edge-cases.md")
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

	table, ok := blocks[0].(document.TableBlock)
	if !ok {
		t.Fatalf("wanted a fallback TableBlock, got %T", blocks[0])
	}
	if len(table.Headers) != 2 || len(table.Rows) != 2 {
		t.Errorf("wanted 2 Headers and 2 Rows, got %d Headers and %d Rows", len(table.Headers), len(table.Rows))
	}
}

func TestParseManifestoDocumentShapesTheSample(t *testing.T) {
	doc, err := ParseManifestoDocument("../testdata/sample.md")
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
			case document.Triad:
				sawTriad = true
			case document.Callout:
				sawCallout = len(v.Label) > 0
			case document.Quote:
				sawQuote = true
			case document.CodeBlock:
				sawCode = len(v.Text) > 0
			case *document.Paragraph:
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

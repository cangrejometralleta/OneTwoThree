package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderDocumentToPDFWritesAValidFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.pdf")

	if err := RenderDocumentToPDF(sampleDocumentForRender(), out); err != nil {
		t.Fatalf("Render Must Succeed, got %v", err)
	}

	assertIsPDF(t, out)
}

func TestRenderDocumentToPDFHandlesACoverOnlyDocument(t *testing.T) {
	doc := Document{Cover: ExtractCoverBlock("Title", "Sub", []string{"Epigraph"})}
	out := filepath.Join(t.TempDir(), "empty.pdf")

	if err := RenderDocumentToPDF(doc, out); err != nil {
		t.Fatalf("a Cover-only Document Must still Render, got %v", err)
	}

	assertIsPDF(t, out)
}

// The Regression this Test Pins: a multi-line CodeBlock once Mismeasured
// its own Height, because IsFitMultiCell Ignored a literal Newline.
func TestRenderDocumentToPDFRendersTheParsedSample(t *testing.T) {
	doc, err := ParseManifestoDocument("testdata/sample.md")
	if err != nil {
		t.Fatalf("the Sample Must Parse: %v", err)
	}

	out := filepath.Join(t.TempDir(), "sample.pdf")
	if err := RenderDocumentToPDF(doc, out); err != nil {
		t.Fatalf("the Sample Must Render, got %v", err)
	}

	assertIsPDF(t, out)
}

func TestRenderDocumentToPDFRendersAFallbackTable(t *testing.T) {
	doc := Document{
		Cover: ExtractCoverBlock("Title", "Sub", nil),
		Sections: []Section{{
			Index: "One", Title: "Section",
			Blocks: []Block{
				TableBlock{
					Headers: []string{"A", "B"},
					Rows:    [][]string{{"1", "2"}, {"3", "4"}},
				},
			},
		}},
	}
	out := filepath.Join(t.TempDir(), "table.pdf")

	if err := RenderDocumentToPDF(doc, out); err != nil {
		t.Fatalf("a fallback Table Must Render, got %v", err)
	}

	assertIsPDF(t, out)
}

// A Document long enough to Force a Page Break Exercises Code no other
// Test Reaches: ensureRoom and layoutWrappedText only Break their
// Branch when the Content truly Overflows one Page.
func TestRenderDocumentToPDFBreaksAcrossPages(t *testing.T) {
	line := "A Paragraph Repeated only to Fill the Page and Force a Break, " +
		"long enough that a handful of Copies Overflow one A4 Sheet."

	var blocks []Block
	for range 40 {
		blocks = append(blocks, &Paragraph{Text: line})
	}

	var items []string
	for range 40 {
		items = append(items, line)
	}
	blocks = append(blocks, ListBlock{Items: items, Ordered: true})

	doc := Document{
		Cover:    ExtractCoverBlock("Title", "Sub", nil),
		Sections: []Section{{Index: "One", Title: "Section", Blocks: blocks}},
	}

	out := filepath.Join(t.TempDir(), "long.pdf")
	if err := RenderDocumentToPDF(doc, out); err != nil {
		t.Fatalf("a long Document Must still Render, got %v", err)
	}

	assertIsPDF(t, out)
}

func TestRenderDocumentToPDFRefusesAnUnwritablePath(t *testing.T) {
	doc := Document{Cover: ExtractCoverBlock("Title", "Sub", nil)}

	if err := RenderDocumentToPDF(doc, filepath.Join("no-such-directory", "out.pdf")); err == nil {
		t.Fatal("an unwritable Path Must Refuse, not Panic Silently")
	}
}

// sampleDocumentForRender Builds one Document that Exercises every Block Kind
// drawBlock Knows, so a Panic in any of them Fails this Test.
func sampleDocumentForRender() Document {
	return Document{
		Cover: ExtractCoverBlock("Sample Title", "Sample Subtitle", []string{"An Epigraph Line, long enough to Wrap."}),
		Sections: []Section{{
			Index: "One",
			Title: "First Section",
			Blocks: []Block{
				&Paragraph{Text: "A plain Paragraph, long enough to Wrap across more than one Line inside the Content Width the Page Allows."},
				ListBlock{Items: []string{"First Item", "Second Item"}, Ordered: true},
				ListBlock{Items: []string{"Bulleted Item"}, Ordered: false},
				Triad{Columns: []TriadColumn{
					{Title: "Receive", Text: "Read the Input"},
					{Title: "Return", Text: "Hand it Back"},
				}},
				Callout{Label: "Note", Paragraphs: []string{"A labeled Box.", "A second Paragraph inside it."}},
				Quote{Paragraphs: []string{"A plain Quote."}},
				CodeBlock{Text: "func F() {\n\treturn\n}"},
				&Paragraph{Text: "The Closing Line.", Italic: true, Closing: true},
			},
		}},
	}
}

// assertIsPDF Checks the Bytes a real PDF Reader would Refuse without.
func assertIsPDF(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the Output Must Exist: %v", err)
	}

	if !strings.HasPrefix(string(data), "%PDF-") {
		t.Fatalf("the Output Must Open with a PDF Header, got %q", data[:min(20, len(data))])
	}
	if !strings.Contains(string(data), "%%EOF") {
		t.Fatal("the Output Must Close with a PDF Trailer")
	}
	if len(data) < 1000 {
		t.Fatalf("a Document with embedded Fonts Must not be Tiny, got %d Bytes", len(data))
	}
}

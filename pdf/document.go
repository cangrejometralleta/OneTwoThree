package main

import "strings"

// CoverMeta Names the Booklet, the way build.py's META Constant did.
const CoverMeta = "Manifesto · Rhythm · Code"

// Document is the whole Manifesto, Shaped for the Page.
type Document struct {
	Cover    Cover
	Sections []Section
}

// Cover Holds the Title, the Subtitle and the Epigraph.
type Cover struct {
	Title    string
	Subtitle string
	Epigraph []string
	Meta     string
}

// Section is one H2, Indexed if the Source Named an Index.
type Section struct {
	Index  string
	Title  string
	Blocks []Block
}

// Block is one Shape a Section Body can Take.
type Block interface{ isBlock() }

// Paragraph Carries Text, and Marks itself if it Closes the Book.
type Paragraph struct {
	Text    string
	Italic  bool
	Closing bool
}

func (*Paragraph) isBlock() {}

// ListBlock Holds every Item, Numbered or not.
type ListBlock struct {
	Items   []string
	Ordered bool
}

func (ListBlock) isBlock() {}

// Triad Holds three Columns where a single-row Table once Stood.
type Triad struct {
	Columns []TriadColumn
}

func (Triad) isBlock() {}

// TriadColumn Pairs a Header with its Cell.
type TriadColumn struct {
	Title string
	Text  string
}

// Callout Holds a labeled Box, Born from a bold-led Quote.
type Callout struct {
	Label      string
	Paragraphs []string
}

func (Callout) isBlock() {}

// Quote Holds a plain Blockquote, italic and bordered.
type Quote struct {
	Paragraphs []string
}

func (Quote) isBlock() {}

// CodeBlock Holds preformatted Text, set in the Mono Face.
type CodeBlock struct {
	Text string
}

func (CodeBlock) isBlock() {}

// TableBlock is the Fallback for a Table a Triad cannot Fit.
type TableBlock struct {
	Headers []string
	Rows    [][]string
}

func (TableBlock) isBlock() {}

// ExtractCoverBlock Turns the Title, the Subtitle and the Epigraph
// into the Cover the first Page Needs.
func ExtractCoverBlock(title, subtitle string, epigraph []string) Cover {
	return Cover{Title: title, Subtitle: subtitle, Epigraph: epigraph, Meta: CoverMeta}
}

// SplitTitleIndex Separates "One · Title" into an Index and a Title.
// A Heading with no Dot Keeps its whole Text as the Title.
func SplitTitleIndex(heading string) (index, title string) {
	before, after, found := strings.Cut(heading, "·")
	if !found {
		return "", heading
	}

	return strings.TrimSpace(before), strings.TrimSpace(after)
}

// BuildTriadBlock Zips Headers and Cells into three Columns.
// A length Mismatch Refuses the Triad, and the Table Falls back whole.
func BuildTriadBlock(headers, cells []string) (Triad, bool) {
	if len(headers) == 0 || len(headers) != len(cells) {
		return Triad{}, false
	}

	columns := make([]TriadColumn, len(headers))
	for i := range headers {
		columns[i] = TriadColumn{Title: headers[i], Text: cells[i]}
	}

	return Triad{Columns: columns}, true
}

// BuildCalloutBlock Names the Box after its bold Label.
// An empty Label Means the Quote never Qualified.
func BuildCalloutBlock(label string, body []string) (Callout, bool) {
	if label == "" {
		return Callout{}, false
	}

	return Callout{Label: label, Paragraphs: body}, true
}

// MarkClosingParagraph Finds the last Paragraph in the whole Document
// and Marks it Closing, but only when it was already Set in Italics.
func MarkClosingParagraph(doc *Document) {
	for s := len(doc.Sections) - 1; s >= 0; s-- {
		blocks := doc.Sections[s].Blocks
		for b := len(blocks) - 1; b >= 0; b-- {
			if p, ok := blocks[b].(*Paragraph); ok {
				p.Closing = p.Italic
				return
			}
		}
	}
}

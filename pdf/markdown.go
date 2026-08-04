package main

import (
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// ParseManifestoDocument Reads the Source, Parses it, then Shapes it.
func ParseManifestoDocument(path string) (Document, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	root := parseMarkdownSource(source)
	title, subtitle, epigraph, body := splitCoverNodes(root, source)
	sections := groupIntoSections(body, source)

	doc := Document{Cover: ExtractCoverBlock(title, subtitle, epigraph), Sections: sections}
	MarkClosingParagraph(&doc)

	return doc, nil
}

// parseMarkdownSource Runs goldmark, Tables included.
func parseMarkdownSource(source []byte) ast.Node {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	return md.Parser().Parse(text.NewReader(source))
}

// splitCoverNodes Pulls the Title, the Subtitle and the Epigraph
// off the top, and Returns everything else for the Sections to Claim.
func splitCoverNodes(root ast.Node, source []byte) (title, subtitle string, epigraph []string, rest []ast.Node) {
	var titleNode, subtitleNode, epigraphNode ast.Node

	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		switch {
		case titleNode == nil:
			if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
				titleNode = n
			}
		case subtitleNode == nil:
			if _, ok := n.(*ast.Paragraph); ok {
				subtitleNode = n
			}
		case epigraphNode == nil:
			if _, ok := n.(*ast.Blockquote); ok {
				epigraphNode = n
			}
		}
	}

	if titleNode != nil {
		title = extractText(titleNode, source)
	}
	if subtitleNode != nil {
		subtitle = extractText(subtitleNode, source)
	}
	if epigraphNode != nil {
		epigraph = extractParagraphs(epigraphNode, source)
	}

	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		if n != titleNode && n != subtitleNode && n != epigraphNode {
			rest = append(rest, n)
		}
	}

	return title, subtitle, epigraph, rest
}

// groupIntoSections Starts a new Section on every H2,
// and Files everything after it into that Section's Blocks.
func groupIntoSections(nodes []ast.Node, source []byte) []Section {
	var sections []Section

	for _, n := range nodes {
		if h, ok := n.(*ast.Heading); ok && h.Level == 2 {
			index, title := SplitTitleIndex(extractText(h, source))
			sections = append(sections, Section{Index: index, Title: title})
			continue
		}

		if len(sections) == 0 {
			continue // Content before the first Heading Joins no Section.
		}

		if block := buildBlock(n, source); block != nil {
			last := &sections[len(sections)-1]
			last.Blocks = append(last.Blocks, block)
		}
	}

	return sections
}

// buildBlock Dispatches one top-level Node to the Shape it Becomes.
func buildBlock(n ast.Node, source []byte) Block {
	switch v := n.(type) {
	case *ast.Paragraph:
		return &Paragraph{Text: extractText(v, source), Italic: hasEmphasis(v)}
	case *ast.List:
		return buildListBlock(v, source)
	case *ast.Blockquote:
		return buildQuoteOrCallout(v, source)
	case *ast.FencedCodeBlock:
		return CodeBlock{Text: extractCodeText(v, source)}
	case *ast.CodeBlock:
		return CodeBlock{Text: extractCodeText(v, source)}
	case *east.Table:
		return buildTableOrTriad(v, source)
	default:
		return nil
	}
}

// buildListBlock Reads every Item, Ordered or not.
func buildListBlock(l *ast.List, source []byte) Block {
	var items []string
	for c := l.FirstChild(); c != nil; c = c.NextSibling() {
		items = append(items, strings.TrimSpace(extractText(c, source)))
	}

	return ListBlock{Items: items, Ordered: l.IsOrdered()}
}

// buildQuoteOrCallout Promotes a bold-led Quote into a Callout.
// Anything else Stays a plain Quote.
func buildQuoteOrCallout(bq *ast.Blockquote, source []byte) Block {
	first, ok := bq.FirstChild().(*ast.Paragraph)
	if !ok {
		return Quote{Paragraphs: extractParagraphs(bq, source)}
	}

	label, isCallout := leadingStrongText(first, source)
	if !isCallout {
		return Quote{Paragraphs: extractParagraphs(bq, source)}
	}

	var body []string
	for c := first.NextSibling(); c != nil; c = c.NextSibling() {
		body = append(body, extractText(c, source))
	}

	if callout, ok := BuildCalloutBlock(label, body); ok {
		return callout
	}

	return Quote{Paragraphs: extractParagraphs(bq, source)}
}

// buildTableOrTriad Tries a Triad first; a Table that does not Fit
// Falls back whole.
func buildTableOrTriad(t *east.Table, source []byte) Block {
	headers, rows := extractTableCells(t, source)

	if len(rows) == 1 {
		if triad, ok := BuildTriadBlock(headers, rows[0]); ok {
			return triad
		}
	}

	return TableBlock{Headers: headers, Rows: rows}
}

// extractTableCells Splits a Table into its Header Row and its Body Rows.
func extractTableCells(t *east.Table, source []byte) (headers []string, rows [][]string) {
	for n := t.FirstChild(); n != nil; n = n.NextSibling() {
		row := extractRowCells(n, source)

		if _, isHeader := n.(*east.TableHeader); isHeader {
			headers = row
			continue
		}
		rows = append(rows, row)
	}

	return headers, rows
}

// extractRowCells Reads every Cell in one Row.
func extractRowCells(row ast.Node, source []byte) []string {
	var cells []string
	for c := row.FirstChild(); c != nil; c = c.NextSibling() {
		cells = append(cells, strings.TrimSpace(extractText(c, source)))
	}

	return cells
}

// leadingStrongText Checks whether a Paragraph Opens on a bold Span,
// and Returns that Span's Text when it does.
func leadingStrongText(p *ast.Paragraph, source []byte) (label string, ok bool) {
	em, isEmphasis := p.FirstChild().(*ast.Emphasis)
	if !isEmphasis || em.Level != 2 {
		return "", false
	}

	return extractText(em, source), true
}

// hasEmphasis Answers whether a Node Carries an italic Span anywhere inside.
func hasEmphasis(n ast.Node) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if _, ok := c.(*ast.Emphasis); ok {
			return true
		}
		if hasEmphasis(c) {
			return true
		}
	}

	return false
}

// extractCodeText Joins a Block's raw Lines, straight from the Source.
func extractCodeText(n interface{ Lines() *text.Segments }, source []byte) string {
	return strings.TrimRight(string(n.Lines().Value(source)), "\n")
}

// extractParagraphs Reads one Text per direct Child, for a Quote's Body.
func extractParagraphs(n ast.Node, source []byte) []string {
	var out []string
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		out = append(out, extractText(c, source))
	}

	return out
}

// extractText Walks the inline Children and Joins their Words.
// Formatting Marks Vanish; only the Words Survive.
func extractText(n ast.Node, source []byte) string {
	var b strings.Builder

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			b.Write(t.Value(source))
			if t.SoftLineBreak() || t.HardLineBreak() {
				b.WriteByte(' ')
			}
			continue
		}
		b.WriteString(extractText(c, source))
	}

	return b.String()
}

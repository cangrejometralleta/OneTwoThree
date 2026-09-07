package main

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf" // https://pkg.go.dev/github.com/signintech/gopdf
)

// The Page, in Millimeters.
// style.css Set the same Numbers in mm.
const (
	pageWidth     = 210.0
	pageHeight    = 297.0
	marginLeft    = 24.0
	marginRight   = 24.0
	marginTop     = 22.0
	marginBottom  = 20.0
	contentWidth  = pageWidth - marginLeft - marginRight
	contentBottom = pageHeight - marginBottom
)

const (
	fontSerif = "Serif"
	fontMono  = "Mono"
)

// RenderDocumentToPDF Draws the Cover, Flows every Section, then Stamps
// a Footer once the Page Count is Known.
//
// Every Draw Call below can Fail; must Panics on the rare Failure,
// and this Recover Turns it back into the Error a Caller Expects.
func RenderDocumentToPDF(doc Document, outPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("render Failed: %v", r)
		}
	}()

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		Unit: gopdf.UnitMM,
		PageSize: gopdf.Rect{
			W: pageWidth, H: pageHeight,
		},
	})
	must(loadDocumentFonts(pdf))

	drawCoverPage(pdf, doc.Cover)
	beginBodyPage(pdf)

	for _, section := range doc.Sections {
		drawSection(pdf, section)
	}

	stampPageFooters(pdf)

	return pdf.WritePdf(outPath)
}

// loadDocumentFonts Registers every Style the Renderer Needs,
// straight from the embedded Bytes in fonts.go.
func loadDocumentFonts(pdf *gopdf.GoPdf) error {
	serif := []struct {
		style int
		data  []byte
	}{
		{gopdf.Regular, fontSerifRegular},
		{gopdf.Bold, fontSerifBold},
		{gopdf.Italic, fontSerifItalic},
		{gopdf.Bold | gopdf.Italic, fontSerifBoldItalic},
	}

	for _, f := range serif {
		if err := pdf.AddTTFFontDataWithOption(fontSerif, f.data, gopdf.TtfOption{Style: f.style}); err != nil {
			return err
		}
	}

	return pdf.AddTTFFontDataWithOption(fontMono, fontMonoRegular, gopdf.TtfOption{Style: gopdf.Regular})
}

// beginBodyPage Opens a Page Shaped like every Page but the Cover.
//
// gopdf Registers a Page's Content Object only on its first Draw Call.
// A Page that Ends up Empty — the last one, when nothing Follows it —
// would otherwise Leave stampPageFooters unable to Find it through
// SetPage. The zero-length Line below Draws nothing, and Guarantees
// the Object Exists regardless.
func beginBodyPage(pdf *gopdf.GoPdf) {
	pdf.AddPage()
	pdf.SetMargins(
		marginLeft,
		marginTop,
		marginRight,
		marginBottom,
	)
	pdf.SetXY(marginLeft, marginTop)
	pdf.Line(
		marginLeft,
		marginTop,
		marginLeft,
		marginTop,
	)
}

// ensureRoom Breaks to a new Page when the next Block would not Fit.
func ensureRoom(pdf *gopdf.GoPdf, height float64) {
	if pdf.GetY()+height > contentBottom {
		beginBodyPage(pdf)
	}
}

// drawCoverPage Paints the dark first Page: Number, Title, Subtitle,
// Epigraph, Meta. Page one Carries no Margin and no Footer.
func drawCoverPage(pdf *gopdf.GoPdf, cover Cover) {
	pdf.AddPage()

	setFillColor(pdf, ColorCoverBg)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: 0, Y: 0, Rect: gopdf.Rect{W: pageWidth, H: pageHeight},
		PaintStyle: gopdf.FillPaintStyle,
	}))

	x, y := marginLeft, 60.0

	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Style: "B", Size: 64, Color: ColorAccent}, "1 2 3") + 10
	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Size: 30, Color: ColorCoverInk}, cover.Title) + 3
	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Style: "I", Size: 12.5, Color: ColorCoverSub}, cover.Subtitle) + 22

	must(pdf.SetFont(fontSerif, "", 10.5))
	setTextColor(pdf, ColorCoverQuote)
	for _, line := range cover.Epigraph {
		setStrokeColor(pdf, ColorAccent)
		pdf.SetLineWidth(0.6)

		rect := &gopdf.Rect{W: contentWidth - 6, H: contentBottom - y}
		_, height, err := pdf.IsFitMultiCellWithNewline(rect, line)
		must(err)
		rect.H = height

		pdf.Line(x, y, x, y+height)
		pdf.SetXY(x+6, y)
		must(pdf.MultiCellWithOption(rect, line, gopdf.CellOption{Align: gopdf.Left}))
		y += height + 3
	}

	must(pdf.SetFont(fontSerif, "", 9))
	setTextColor(pdf, ColorCoverMeta)
	pdf.SetXY(x, pageHeight-24)
	must(pdf.Cell(&gopdf.Rect{W: contentWidth, H: 6}, letterSpaced(strings.ToUpper(cover.Meta))))
}

// textStyle Names the Font a Line Draws in.
// Four Fields already Reads as a Collection; a Struct Names each one,
// the way SetMargins Names left/top/right/bottom by Position alone —
// a Font Has no such well-known Order to Lean on, so this one Spells it out.
type textStyle struct {
	Family string
	Style  string
	Size   float64
	Color  ColorInk
}

// drawCoverLine Draws one Line of the Cover and Returns the Y it Leaves.
func drawCoverLine(pdf *gopdf.GoPdf, x, y float64, style textStyle, text string) float64 {
	must(pdf.SetFont(style.Family, style.Style, style.Size))
	setTextColor(pdf, style.Color)
	pdf.SetXY(x, y)
	must(pdf.Cell(&gopdf.Rect{W: contentWidth, H: lineHeightMM(style.Size)}, text))

	return y + lineHeightMM(style.Size)
}

// drawSection Draws one Heading, then every Block it Owns.
func drawSection(pdf *gopdf.GoPdf, section Section) {
	drawSectionHeading(pdf, section)

	for _, block := range section.Blocks {
		drawBlock(pdf, block)
	}
}

// drawSectionHeading Draws the Index, the Title, and the Rule under it.
func drawSectionHeading(pdf *gopdf.GoPdf, section Section) {
	ensureRoom(pdf, 22)
	pdf.SetX(marginLeft)

	y := pdf.GetY()
	if section.Index != "" {
		index := letterSpaced(strings.ToUpper(section.Index))
		y = drawCoverLine(pdf, marginLeft, y, textStyle{Family: fontSerif, Size: 10, Color: ColorAccent}, index) + 1.5
	}
	y = drawCoverLine(pdf, marginLeft, y, textStyle{Family: fontSerif, Size: 15, Color: ColorBody}, section.Title) + 1.5

	setStrokeColor(pdf, ColorRule)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginLeft, y, marginLeft+contentWidth, y)

	pdf.SetXY(marginLeft, y+4)
}

// drawBlock Dispatches one Block to the Shape it Draws as.
func drawBlock(pdf *gopdf.GoPdf, block Block) {
	switch v := block.(type) {
	case *Paragraph:
		drawParagraph(pdf, *v)
	case ListBlock:
		drawListBlock(pdf, v)
	case Triad:
		drawTriad(pdf, v)
	case Callout:
		drawCallout(pdf, v)
	case Quote:
		drawQuote(pdf, v)
	case CodeBlock:
		drawCodeBlock(pdf, v)
	case TableBlock:
		drawTableBlock(pdf, v)
	}
}

// drawParagraph Justifies a Paragraph, Breaking to a new Page
// when the current one Has no Room left.
func drawParagraph(pdf *gopdf.GoPdf, p Paragraph) {
	style := ""
	if p.Closing {
		style = "I"
	}
	must(pdf.SetFont(fontSerif, style, 11))
	setTextColor(pdf, ColorBody)
	if p.Closing {
		setTextColor(pdf, ColorCloseInk)
	}

	align := gopdf.Justify
	if p.Closing {
		align = gopdf.Center
	}

	pdf.SetX(marginLeft)
	height := layoutWrappedText(pdf, p.Text, contentWidth)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: contentWidth, H: height}, p.Text, gopdf.CellOption{Align: align}))
	pdf.SetX(marginLeft)
	pdf.Br(3.2)
}

// layoutWrappedText Measures a Block, Breaking to a new Page first
// when it would not Fit on the current one.
// WithNewline Matters here: a literal "\n" Must Count as its own Line,
// the same way MultiCellWithOption already Draws it.
func layoutWrappedText(pdf *gopdf.GoPdf, text string, width float64) float64 {
	remaining := contentBottom - pdf.GetY()
	fits, height, err := pdf.IsFitMultiCellWithNewline(&gopdf.Rect{W: width, H: remaining}, text)
	must(err)

	if !fits {
		beginBodyPage(pdf)
		_, height, err = pdf.IsFitMultiCellWithNewline(&gopdf.Rect{W: width, H: contentBottom - pdf.GetY()}, text)
		must(err)
	}

	return height
}

// drawListBlock Draws every Item, Numbered in Orange or Bulleted plain.
func drawListBlock(pdf *gopdf.GoPdf, list ListBlock) {
	must(pdf.SetFont(fontSerif, "", 11))

	for i, item := range list.Items {
		ensureRoom(pdf, 8)

		marker := "•"
		if list.Ordered {
			marker = fmt.Sprintf("%d", i+1)
		}

		setTextColor(pdf, ColorAccent)
		pdf.SetXY(marginLeft, pdf.GetY())
		must(pdf.Cell(&gopdf.Rect{W: 8, H: 6}, marker))

		setTextColor(pdf, ColorBody)
		pdf.SetX(marginLeft + 10)
		height := layoutWrappedText(pdf, item, contentWidth-10)
		must(pdf.MultiCellWithOption(&gopdf.Rect{W: contentWidth - 10, H: height}, item, gopdf.CellOption{Align: gopdf.Left}))
		pdf.SetX(marginLeft)
		pdf.Br(1)
	}

	pdf.Br(2.5)
}

// drawTriad Lays three Columns side by side, each Topped with a Rule.
func drawTriad(pdf *gopdf.GoPdf, triad Triad) {
	ensureRoom(pdf, 24)

	colWidth := (contentWidth - 5*float64(len(triad.Columns)-1)) / float64(len(triad.Columns))
	top := pdf.GetY()
	maxHeight := 0.0

	for i, col := range triad.Columns {
		x := marginLeft + float64(i)*(colWidth+5)

		setStrokeColor(pdf, ColorAccent)
		pdf.SetLineWidth(0.6)
		pdf.Line(x, top, x+colWidth, top)

		must(pdf.SetFont(fontSerif, "B", 9.5))
		setTextColor(pdf, ColorBody)
		pdf.SetXY(x, top+2.5)
		must(pdf.Cell(&gopdf.Rect{W: colWidth, H: 5}, col.Title))

		must(pdf.SetFont(fontSerif, "", 9.5))
		pdf.SetXY(x, top+7)
		rect := &gopdf.Rect{W: colWidth, H: contentBottom - top - 7}
		_, height, err := pdf.IsFitMultiCellWithNewline(rect, col.Text)
		must(err)
		rect.H = height
		must(pdf.MultiCellWithOption(rect, col.Text, gopdf.CellOption{Align: gopdf.Left}))

		if used := height + 7; used > maxHeight {
			maxHeight = used
		}
	}

	pdf.SetXY(marginLeft, top+maxHeight+5)
}

// drawCallout Draws a filled Box with a Label, for a bold-led Quote.
func drawCallout(pdf *gopdf.GoPdf, callout Callout) {
	text := strings.Join(callout.Paragraphs, "\n\n")

	must(pdf.SetFont(fontSerif, "", 10.5))
	height := layoutWrappedText(pdf, text, contentWidth-12)
	boxHeight := height + 14

	top := pdf.GetY()
	setFillColor(pdf, ColorCalloutBg)
	setStrokeColor(pdf, ColorAccent)
	pdf.SetLineWidth(1.1)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: marginLeft, Y: top, Rect: gopdf.Rect{W: contentWidth, H: boxHeight},
		PaintStyle: gopdf.DrawFillPaintStyle,
	}))

	must(pdf.SetFont(fontSerif, "B", 8.5))
	setTextColor(pdf, ColorCalloutRule)
	pdf.SetXY(marginLeft+6, top+4)
	must(pdf.Cell(&gopdf.Rect{W: contentWidth - 12, H: 5}, letterSpaced(strings.ToUpper(callout.Label))))

	must(pdf.SetFont(fontSerif, "", 10.5))
	setTextColor(pdf, ColorBody)
	pdf.SetXY(marginLeft+6, top+10)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: contentWidth - 12, H: height}, text, gopdf.CellOption{Align: gopdf.Left}))

	pdf.SetXY(marginLeft, top+boxHeight+5)
}

// drawQuote Draws a plain italic Blockquote, bordered on its left.
func drawQuote(pdf *gopdf.GoPdf, quote Quote) {
	text := strings.Join(quote.Paragraphs, "\n\n")

	must(pdf.SetFont(fontSerif, "I", 11))
	setTextColor(pdf, ColorQuoteInk)
	height := layoutWrappedText(pdf, text, contentWidth-6)

	top := pdf.GetY()
	setStrokeColor(pdf, ColorRule)
	pdf.SetLineWidth(0.6)
	pdf.Line(marginLeft, top, marginLeft, top+height)

	pdf.SetXY(marginLeft+6, top)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: contentWidth - 6, H: height}, text, gopdf.CellOption{Align: gopdf.Left}))
	pdf.SetX(marginLeft)
	pdf.Br(2)
}

// drawCodeBlock Draws preformatted Text in the Mono Face, on a tinted Box.
func drawCodeBlock(pdf *gopdf.GoPdf, code CodeBlock) {
	must(pdf.SetFont(fontMono, "", 8.5))
	height := layoutWrappedText(pdf, code.Text, contentWidth-12)
	boxHeight := height + 7

	top := pdf.GetY()
	setFillColor(pdf, ColorCalloutBg)
	setStrokeColor(pdf, ColorCodeRule)
	pdf.SetLineWidth(1.1)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: marginLeft, Y: top, Rect: gopdf.Rect{W: contentWidth, H: boxHeight},
		PaintStyle: gopdf.DrawFillPaintStyle,
	}))

	setTextColor(pdf, ColorBody)
	pdf.SetXY(marginLeft+6, top+3.5)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: contentWidth - 12, H: height}, code.Text, gopdf.CellOption{Align: gopdf.Left}))

	pdf.SetXY(marginLeft, top+boxHeight+4)
}

// drawTableBlock is the Fallback for a Table a Triad could not Fit:
// plain Rows, Header Bold, Body Regular.
func drawTableBlock(pdf *gopdf.GoPdf, table TableBlock) {
	colWidth := contentWidth / float64(max(len(table.Headers), 1))

	drawTableRow(pdf, table.Headers, tableRowStyle{ColWidth: colWidth, FontStyle: "B"})
	for _, row := range table.Rows {
		drawTableRow(pdf, row, tableRowStyle{ColWidth: colWidth})
	}
	pdf.Br(2.5)
}

// tableRowStyle Carries the two Things every Row in a fallback Table
// Shares: how wide a Column is, and whether the Row Reads Bold.
type tableRowStyle struct {
	ColWidth  float64
	FontStyle string
}

// drawTableRow Draws one Row of the fallback Table.
func drawTableRow(pdf *gopdf.GoPdf, cells []string, style tableRowStyle) {
	ensureRoom(pdf, 8)
	must(pdf.SetFont(fontSerif, style.FontStyle, 9.5))
	setTextColor(pdf, ColorBody)

	for i, cell := range cells {
		pdf.SetXY(marginLeft+float64(i)*style.ColWidth, pdf.GetY())
		must(pdf.CellWithOption(&gopdf.Rect{W: style.ColWidth, H: 6}, cell, gopdf.CellOption{Border: gopdf.AllBorders}))
	}
	pdf.SetX(marginLeft)
	pdf.Br(6)
}

// stampPageFooters Numbers every Page but the Cover, once the final
// Count is Known. gopdf Lets a finished Page be Reopened for this.
func stampPageFooters(pdf *gopdf.GoPdf) {
	total := pdf.GetNumberOfPages()

	for page := 2; page <= total; page++ {
		must(pdf.SetPage(page))
		must(pdf.SetFont(fontSerif, "", 8.5))
		setTextColor(pdf, ColorPageNumber)
		pdf.SetXY(marginLeft, pageHeight-marginBottom+8)
		must(pdf.CellWithOption(&gopdf.Rect{W: contentWidth, H: 6}, fmt.Sprintf("%d", page-1), gopdf.CellOption{Align: gopdf.Center}))
	}
}

// lineHeightMM Turns a Font Size in Points into a Line's Height in mm.
func lineHeightMM(sizePt float64) float64 {
	return sizePt * 0.3528 * 1.25
}

// letterSpaced Loosens uppercase Text the way style.css's
// letter-spacing Did, one thin Space between each Letter.
func letterSpaced(s string) string {
	return strings.Join(strings.Split(s, ""), " ")
}

func setFillColor(pdf *gopdf.GoPdf, c ColorInk) { pdf.SetFillColor(uint8(c.R), uint8(c.G), uint8(c.B)) }
func setTextColor(pdf *gopdf.GoPdf, c ColorInk) { pdf.SetTextColor(uint8(c.R), uint8(c.G), uint8(c.B)) }
func setStrokeColor(pdf *gopdf.GoPdf, c ColorInk) {
	pdf.SetStrokeColor(uint8(c.R), uint8(c.G), uint8(c.B))
}

// must Panics on a Failure a healthy Render should never Produce.
// RenderDocumentToPDF Recovers it, once, at the top.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

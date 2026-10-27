package render

import (
	"github.com/cangrejometralleta/OneTwoThree/pdf/document"
	"github.com/cangrejometralleta/OneTwoThree/pdf/style"

	"fmt"
	"strings"

	"github.com/signintech/gopdf" // https://pkg.go.dev/github.com/signintech/gopdf
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
func RenderDocumentToPDF(doc document.Document, outPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("render Failed: %v", r)
		}
	}()

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		Unit: gopdf.UnitMM,
		PageSize: gopdf.Rect{
			W: style.PageWidth, H: style.PageHeight,
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
		style.MarginLeft,
		style.MarginTop,
		style.MarginRight,
		style.MarginBottom,
	)
	pdf.SetXY(style.MarginLeft, style.MarginTop)
	pdf.Line(
		style.MarginLeft,
		style.MarginTop,
		style.MarginLeft,
		style.MarginTop,
	)
}

// ensureRoom Breaks to a new Page when the next Block would not Fit.
func ensureRoom(pdf *gopdf.GoPdf, height float64) {
	if pdf.GetY()+height > style.ContentBottom {
		beginBodyPage(pdf)
	}
}

// drawCoverPage Paints the dark first Page: Number, Title, Subtitle,
// Epigraph, Meta. Page one Carries no Margin and no Footer.
func drawCoverPage(pdf *gopdf.GoPdf, cover document.Cover) {
	pdf.AddPage()

	setFillColor(pdf, style.ColorCoverBg)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: 0, Y: 0, Rect: gopdf.Rect{W: style.PageWidth, H: style.PageHeight},
		PaintStyle: gopdf.FillPaintStyle,
	}))

	x, y := style.MarginLeft, style.CoverTop

	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Style: "B", Size: style.SizeCoverNumber, Color: style.ColorAccent}, "1 2 3") + style.CoverNumberGap
	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Size: style.SizeCoverTitle, Color: style.ColorCoverInk}, cover.Title) + style.CoverTitleGap
	y = drawCoverLine(pdf, x, y, textStyle{Family: fontSerif, Style: "I", Size: style.SizeCoverSubtitle, Color: style.ColorCoverSub}, cover.Subtitle) + style.CoverEpigraphGap

	must(pdf.SetFont(fontSerif, "", style.SizeEpigraph))
	setTextColor(pdf, style.ColorCoverQuote)
	for _, line := range cover.Epigraph {
		setStrokeColor(pdf, style.ColorAccent)
		pdf.SetLineWidth(style.RuleThin)

		rect := &gopdf.Rect{W: style.ContentWidth - style.IndentQuote, H: style.ContentBottom - y}
		_, height, err := pdf.IsFitMultiCellWithNewline(rect, line)
		must(err)
		rect.H = height

		pdf.Line(x, y, x, y+height)
		pdf.SetXY(x+style.IndentQuote, y)
		must(pdf.MultiCellWithOption(rect, line, gopdf.CellOption{Align: gopdf.Left}))
		y += height + style.CoverLineGap
	}

	must(pdf.SetFont(fontSerif, "", style.SizeCoverMeta))
	setTextColor(pdf, style.ColorCoverMeta)
	pdf.SetXY(x, style.PageHeight-style.CoverMetaBottom)
	must(pdf.Cell(&gopdf.Rect{W: style.ContentWidth, H: style.HeightRow}, letterSpaced(strings.ToUpper(cover.Meta))))
}

// textStyle Names the Font a Line Draws in.
// Four Fields already Reads as a Collection; a Struct Names each one,
// the way SetMargins Names left/top/right/bottom by Position alone —
// a Font Has no such well-known Order to Lean on, so this one Spells it out.
type textStyle struct {
	Family string
	Style  string
	Size   float64
	Color  style.ColorInk
}

// drawCoverLine Draws one Line of the Cover and Returns the Y it Leaves.
func drawCoverLine(pdf *gopdf.GoPdf, x, y float64, face textStyle, text string) float64 {
	must(pdf.SetFont(face.Family, face.Style, face.Size))
	setTextColor(pdf, face.Color)
	pdf.SetXY(x, y)
	must(pdf.Cell(&gopdf.Rect{W: style.ContentWidth, H: lineHeightMM(face.Size)}, text))

	return y + lineHeightMM(face.Size)
}

// drawSection Draws one Heading, then every Block it Owns.
func drawSection(pdf *gopdf.GoPdf, section document.Section) {
	drawSectionHeading(pdf, section)

	for _, block := range section.Blocks {
		drawBlock(pdf, block)
	}
}

// drawSectionHeading Draws the Index, the Title, and the Rule under it.
func drawSectionHeading(pdf *gopdf.GoPdf, section document.Section) {
	ensureRoom(pdf, style.RoomHeading)
	pdf.SetX(style.MarginLeft)

	y := pdf.GetY()
	if section.Index != "" {
		index := letterSpaced(strings.ToUpper(section.Index))
		y = drawCoverLine(pdf, style.MarginLeft, y, textStyle{Family: fontSerif, Size: style.SizeSectionIndex, Color: style.ColorAccent}, index) + style.GapIndex
	}
	y = drawCoverLine(pdf, style.MarginLeft, y, textStyle{Family: fontSerif, Size: style.SizeSectionTitle, Color: style.ColorBody}, section.Title) + style.GapIndex

	setStrokeColor(pdf, style.ColorRule)
	pdf.SetLineWidth(style.RuleHairline)
	pdf.Line(style.MarginLeft, y, style.MarginLeft+style.ContentWidth, y)

	pdf.SetXY(style.MarginLeft, y+style.GapHeading)
}

// drawBlock Dispatches one Block to the Shape it Draws as.
func drawBlock(pdf *gopdf.GoPdf, block document.Block) {
	switch v := block.(type) {
	case *document.Paragraph:
		drawParagraph(pdf, *v)
	case document.ListBlock:
		drawListBlock(pdf, v)
	case document.Triad:
		drawTriad(pdf, v)
	case document.Callout:
		drawCallout(pdf, v)
	case document.Quote:
		drawQuote(pdf, v)
	case document.CodeBlock:
		drawCodeBlock(pdf, v)
	case document.TableBlock:
		drawTableBlock(pdf, v)
	}
}

// drawParagraph Justifies a Paragraph, Breaking to a new Page
// when the current one Has no Room left.
func drawParagraph(pdf *gopdf.GoPdf, p document.Paragraph) {
	weight := ""
	if p.Closing {
		weight = "I"
	}
	must(pdf.SetFont(fontSerif, weight, style.SizeBody))
	setTextColor(pdf, style.ColorBody)
	if p.Closing {
		setTextColor(pdf, style.ColorCloseInk)
	}

	align := gopdf.Justify
	if p.Closing {
		align = gopdf.Center
	}

	pdf.SetX(style.MarginLeft)
	height := layoutWrappedText(pdf, p.Text, style.ContentWidth)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: style.ContentWidth, H: height}, p.Text, gopdf.CellOption{Align: align}))
	pdf.SetX(style.MarginLeft)
	pdf.Br(style.GapParagraph)
}

// layoutWrappedText Measures a Block, Breaking to a new Page first
// when it would not Fit on the current one.
// WithNewline Matters here: a literal "\n" Must Count as its own Line,
// the same way MultiCellWithOption already Draws it.
func layoutWrappedText(pdf *gopdf.GoPdf, text string, width float64) float64 {
	remaining := style.ContentBottom - pdf.GetY()
	fits, height, err := pdf.IsFitMultiCellWithNewline(&gopdf.Rect{W: width, H: remaining}, text)
	must(err)

	if !fits {
		beginBodyPage(pdf)
		_, height, err = pdf.IsFitMultiCellWithNewline(&gopdf.Rect{W: width, H: style.ContentBottom - pdf.GetY()}, text)
		must(err)
	}

	return height
}

// drawListBlock Draws every Item, Numbered in Orange or Bulleted plain.
func drawListBlock(pdf *gopdf.GoPdf, list document.ListBlock) {
	must(pdf.SetFont(fontSerif, "", style.SizeBody))

	for i, item := range list.Items {
		ensureRoom(pdf, style.RoomLine)

		marker := "•"
		if list.Ordered {
			marker = fmt.Sprintf("%d", i+1)
		}

		setTextColor(pdf, style.ColorAccent)
		pdf.SetXY(style.MarginLeft, pdf.GetY())
		must(pdf.Cell(&gopdf.Rect{W: style.WidthMarker, H: style.HeightRow}, marker))

		setTextColor(pdf, style.ColorBody)
		pdf.SetX(style.MarginLeft + style.IndentList)
		height := layoutWrappedText(pdf, item, style.ContentWidth-style.IndentList)
		must(pdf.MultiCellWithOption(&gopdf.Rect{W: style.ContentWidth - style.IndentList, H: height}, item, gopdf.CellOption{Align: gopdf.Left}))
		pdf.SetX(style.MarginLeft)
		pdf.Br(style.GapListItem)
	}

	pdf.Br(style.GapList)
}

// drawTriad Lays three Columns side by side, each Topped with a Rule.
func drawTriad(pdf *gopdf.GoPdf, triad document.Triad) {
	ensureRoom(pdf, style.RoomTriad)

	colWidth := (style.ContentWidth - style.GapBlock*float64(len(triad.Columns)-1)) / float64(len(triad.Columns))
	top := pdf.GetY()
	maxHeight := 0.0

	for i, col := range triad.Columns {
		x := style.MarginLeft + float64(i)*(colWidth+style.GapBlock)

		setStrokeColor(pdf, style.ColorAccent)
		pdf.SetLineWidth(style.RuleThin)
		pdf.Line(x, top, x+colWidth, top)

		must(pdf.SetFont(fontSerif, "B", style.SizeColumn))
		setTextColor(pdf, style.ColorBody)
		pdf.SetXY(x, top+style.OffsetColumnTitle)
		must(pdf.Cell(&gopdf.Rect{W: colWidth, H: style.HeightColumn}, col.Title))

		must(pdf.SetFont(fontSerif, "", style.SizeColumn))
		pdf.SetXY(x, top+style.OffsetColumnText)
		rect := &gopdf.Rect{W: colWidth, H: style.ContentBottom - top - style.OffsetColumnText}
		_, height, err := pdf.IsFitMultiCellWithNewline(rect, col.Text)
		must(err)
		rect.H = height
		must(pdf.MultiCellWithOption(rect, col.Text, gopdf.CellOption{Align: gopdf.Left}))

		if used := height + style.OffsetColumnText; used > maxHeight {
			maxHeight = used
		}
	}

	pdf.SetXY(style.MarginLeft, top+maxHeight+style.GapBlock)
}

// drawCallout Draws a filled Box with a Label, for a bold-led Quote.
func drawCallout(pdf *gopdf.GoPdf, callout document.Callout) {
	text := strings.Join(callout.Paragraphs, "\n\n")

	must(pdf.SetFont(fontSerif, "", style.SizeCallout))
	height := layoutWrappedText(pdf, text, style.ContentWidth-style.PadCalloutSide)
	boxHeight := height + style.PadCalloutBox

	top := pdf.GetY()
	setFillColor(pdf, style.ColorCalloutBg)
	setStrokeColor(pdf, style.ColorAccent)
	pdf.SetLineWidth(style.RuleThick)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: style.MarginLeft, Y: top, Rect: gopdf.Rect{W: style.ContentWidth, H: boxHeight},
		PaintStyle: gopdf.DrawFillPaintStyle,
	}))

	must(pdf.SetFont(fontSerif, "B", style.SizeCalloutLabel))
	setTextColor(pdf, style.ColorCalloutRule)
	pdf.SetXY(style.MarginLeft+style.IndentCallout, top+style.GapHeading)
	must(pdf.Cell(&gopdf.Rect{W: style.ContentWidth - style.PadCalloutSide, H: style.HeightColumn}, letterSpaced(strings.ToUpper(callout.Label))))

	must(pdf.SetFont(fontSerif, "", style.SizeCallout))
	setTextColor(pdf, style.ColorBody)
	pdf.SetXY(style.MarginLeft+style.IndentCallout, top+style.OffsetCalloutText)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: style.ContentWidth - style.PadCalloutSide, H: height}, text, gopdf.CellOption{Align: gopdf.Left}))

	pdf.SetXY(style.MarginLeft, top+boxHeight+style.GapBlock)
}

// drawQuote Draws a plain italic Blockquote, bordered on its left.
func drawQuote(pdf *gopdf.GoPdf, quote document.Quote) {
	text := strings.Join(quote.Paragraphs, "\n\n")

	must(pdf.SetFont(fontSerif, "I", style.SizeBody))
	setTextColor(pdf, style.ColorQuoteInk)
	height := layoutWrappedText(pdf, text, style.ContentWidth-style.IndentQuote)

	top := pdf.GetY()
	setStrokeColor(pdf, style.ColorRule)
	pdf.SetLineWidth(style.RuleThin)
	pdf.Line(style.MarginLeft, top, style.MarginLeft, top+height)

	pdf.SetXY(style.MarginLeft+style.IndentQuote, top)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: style.ContentWidth - style.IndentQuote, H: height}, text, gopdf.CellOption{Align: gopdf.Left}))
	pdf.SetX(style.MarginLeft)
	pdf.Br(style.GapQuote)
}

// drawCodeBlock Draws preformatted Text in the Mono Face, on a tinted Box.
func drawCodeBlock(pdf *gopdf.GoPdf, code document.CodeBlock) {
	must(pdf.SetFont(fontMono, "", style.SizeCode))
	height := layoutWrappedText(pdf, code.Text, style.ContentWidth-style.PadCalloutSide)
	boxHeight := height + style.PadCodeBox

	top := pdf.GetY()
	setFillColor(pdf, style.ColorCalloutBg)
	setStrokeColor(pdf, style.ColorCodeRule)
	pdf.SetLineWidth(style.RuleThick)
	must(pdf.RectFromUpperLeftWithOpts(gopdf.DrawableRectOptions{
		X: style.MarginLeft, Y: top, Rect: gopdf.Rect{W: style.ContentWidth, H: boxHeight},
		PaintStyle: gopdf.DrawFillPaintStyle,
	}))

	setTextColor(pdf, style.ColorBody)
	pdf.SetXY(style.MarginLeft+style.IndentCallout, top+style.OffsetCodeText)
	must(pdf.MultiCellWithOption(&gopdf.Rect{W: style.ContentWidth - style.PadCalloutSide, H: height}, code.Text, gopdf.CellOption{Align: gopdf.Left}))

	pdf.SetXY(style.MarginLeft, top+boxHeight+style.GapCode)
}

// drawTableBlock is the Fallback for a Table a Triad could not Fit:
// plain Rows, Header Bold, Body Regular.
func drawTableBlock(pdf *gopdf.GoPdf, table document.TableBlock) {
	colWidth := style.ContentWidth / float64(max(len(table.Headers), 1))

	drawTableRow(pdf, table.Headers, tableRowStyle{ColWidth: colWidth, FontStyle: "B"})
	for _, row := range table.Rows {
		drawTableRow(pdf, row, tableRowStyle{ColWidth: colWidth})
	}
	pdf.Br(style.GapList)
}

// tableRowStyle Carries the two Things every Row in a fallback Table
// Shares: how wide a Column is, and whether the Row Reads Bold.
type tableRowStyle struct {
	ColWidth  float64
	FontStyle string
}

// drawTableRow Draws one Row of the fallback Table.
func drawTableRow(pdf *gopdf.GoPdf, cells []string, row tableRowStyle) {
	ensureRoom(pdf, style.RoomLine)
	must(pdf.SetFont(fontSerif, row.FontStyle, style.SizeColumn))
	setTextColor(pdf, style.ColorBody)

	for i, cell := range cells {
		pdf.SetXY(style.MarginLeft+float64(i)*row.ColWidth, pdf.GetY())
		must(pdf.CellWithOption(&gopdf.Rect{W: row.ColWidth, H: style.HeightRow}, cell, gopdf.CellOption{Border: gopdf.AllBorders}))
	}
	pdf.SetX(style.MarginLeft)
	pdf.Br(style.HeightRow)
}

// stampPageFooters Numbers every Page but the Cover, once the final
// Count is Known. gopdf Lets a finished Page be Reopened for this.
func stampPageFooters(pdf *gopdf.GoPdf) {
	total := pdf.GetNumberOfPages()

	for page := 2; page <= total; page++ {
		must(pdf.SetPage(page))
		must(pdf.SetFont(fontSerif, "", 8.5))
		setTextColor(pdf, style.ColorPageNumber)
		pdf.SetXY(style.MarginLeft, style.PageHeight-style.MarginBottom+8)
		must(pdf.CellWithOption(&gopdf.Rect{W: style.ContentWidth, H: 6}, fmt.Sprintf("%d", page-1), gopdf.CellOption{Align: gopdf.Center}))
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

func setFillColor(pdf *gopdf.GoPdf, c style.ColorInk) {
	pdf.SetFillColor(uint8(c.R), uint8(c.G), uint8(c.B))
}
func setTextColor(pdf *gopdf.GoPdf, c style.ColorInk) {
	pdf.SetTextColor(uint8(c.R), uint8(c.G), uint8(c.B))
}
func setStrokeColor(pdf *gopdf.GoPdf, c style.ColorInk) {
	pdf.SetStrokeColor(uint8(c.R), uint8(c.G), uint8(c.B))
}

// must Panics on a Failure a healthy Render should never Produce.
// RenderDocumentToPDF Recovers it, once, at the top.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

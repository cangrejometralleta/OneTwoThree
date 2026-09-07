package style

// The Page, in Millimeters. style.css Set the same Numbers in mm.
const (
	PageWidth     = 210.0
	PageHeight    = 297.0
	MarginLeft    = 24.0
	MarginRight   = 24.0
	MarginTop     = 22.0
	MarginBottom  = 20.0
	ContentWidth  = PageWidth - MarginLeft - MarginRight
	ContentBottom = PageHeight - MarginBottom
)

// The Type Scale, in Points. Each Name Says the Role, never the Rank.
// Read 10.5 alone and it Says nothing; read SizeEpigraph and it Says
// where the Number Lands on the Page.
const (
	SizeCoverNumber   = 64.0
	SizeCoverTitle    = 30.0
	SizeCoverSubtitle = 12.5
	SizeEpigraph      = 10.5
	SizeCoverMeta     = 9.0
	SizeSectionIndex  = 10.0
	SizeSectionTitle  = 15.0
	SizeBody          = 11.0
	SizeColumn        = 9.5
	SizeCallout       = 10.5
	SizeCalloutLabel  = 8.5
	SizeCode          = 8.5
)

// The three Rule Weights. A fourth would be a Weight nobody can Tell apart.
const (
	RuleHairline = 0.3 // under a Section Heading
	RuleThin     = 0.6 // the Epigraph Bar, a Triad Top, a Quote Edge
	RuleThick    = 1.1 // a Callout Box, a Code Box
)

// The Space between Things, in Millimeters.
const (
	CoverTop         = 60.0
	CoverNumberGap   = 10.0
	CoverTitleGap    = 3.0
	CoverEpigraphGap = 22.0
	CoverLineGap     = 3.0
	CoverMetaBottom  = 24.0

	IndentQuote   = 6.0
	IndentCallout = 6.0
	IndentList    = 10.0
	WidthMarker   = 8.0

	GapParagraph = 3.2
	GapListItem  = 1.0
	GapList      = 2.5
	GapQuote     = 2.0
	GapBlock     = 5.0
	GapCode      = 4.0
	GapHeading   = 4.0
	GapIndex     = 1.5

	RoomHeading = 22.0
	RoomLine    = 8.0
	RoomTriad   = 24.0

	HeightRow    = 6.0
	HeightColumn = 5.0

	OffsetColumnTitle = 2.5
	OffsetColumnText  = 7.0
	OffsetCalloutText = 10.0
	OffsetCodeText    = 3.5
	PadCalloutBox     = 14.0
	PadCodeBox        = 7.0
	PadCalloutSide    = 12.0
)

package main

// ColorInk Names a Hex the old style.css only Numbered.
// Every Value here Traces back to style.css, one Name per Role.
// A Color can be Described as Three 8-bit RGB values,
// Where Each in Ranges
// From 0 to 255,
// Or 0x00 to 0xff in Hexadecimal.

type ColorInk struct{ R, G, B int }

var (
	ColorBody       = ColorInk{0x23, 0x20, 0x1d} // body { color }
	ColorPageNumber = ColorInk{0x9a, 0x8f, 0x84} // @bottom-center { color }

	ColorCoverBg     = ColorInk{0x1d, 0x1a, 0x18} // .cover { background }
	ColorCoverInk    = ColorInk{0xf2, 0xec, 0xe4} // .cover { color }
	ColorCoverSub    = ColorInk{0xbf, 0xb3, 0xa4} // .cover .sub
	ColorCoverQuote  = ColorInk{0xcd, 0xc2, 0xb4} // .cover blockquote
	ColorCoverMeta   = ColorInk{0x8c, 0x81, 0x77} // .cover .meta
	ColorAccent      = ColorInk{0xd9, 0x8b, 0x46} // .cover .num, h2 .idx, ol li::before
	ColorRule        = ColorInk{0xde, 0xd5, 0xc9} // h2 { border-bottom }
	ColorCalloutBg   = ColorInk{0xf7, 0xf2, 0xea} // .callout, pre { background }
	ColorCalloutRule = ColorInk{0xa9, 0x75, 0x4a} // .callout .label
	ColorQuoteInk    = ColorInk{0x4a, 0x42, 0x3a} // blockquote { color }
	ColorCloseInk    = ColorInk{0x6b, 0x61, 0x58} // .close { color }
	ColorCodeRule    = ColorInk{0xe0, 0xd3, 0xc2} // pre { border-left }
	ColorCodeBg      = ColorInk{0xf2, 0xed, 0xe5} // p code { background }
)

package school

import (
	"math"
	"math/bits"
)

// DefaultPageSize Windows a Page when the Caller Names a page but no size.
const DefaultPageSize = 10

// Page Windows a List, Zero based. WholeSet Asks for every Record.
type Page struct {
	Number int64
	Size   int64
}

// WholeSet Stands for a List no Window Cuts.
var WholeSet = Page{Number: 0, Size: math.MaxInt64}

// BuildPage Holds the Bounds before a Page Exists: Number at least zero, Size at least one.
func BuildPage(number, size int64) (Page, error) {
	if number < 0 || size < 1 {
		return Page{}, ErrPageIsInvalid
	}

	return Page{Number: number, Size: size}, nil
}

// Offset Counts the Records before the Window, as far as 64 bits can Count.
func (p Page) Offset() int64 {
	high, low := bits.Mul64(uint64(p.Number), uint64(p.Size))
	if high != 0 || low > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(low)
}

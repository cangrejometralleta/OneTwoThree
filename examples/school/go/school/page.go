package school

// DefaultPageSize Windows a Page when the Caller Names a page but no size.
const DefaultPageSize = 10

// Page Windows a List. A Size of zero Asks for the whole Set.
type Page struct {
	Number int
	Size   int
}

// CheckPageBounds Refuses a Page the Store cannot Serve.
func (p Page) CheckPageBounds() error {
	if p.Number < 0 || p.Size < 0 {
		return ErrPageIsInvalid
	}

	return nil
}

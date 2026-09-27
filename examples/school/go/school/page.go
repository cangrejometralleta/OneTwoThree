package school

// Page Carries Pagination without Naming a Database.
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

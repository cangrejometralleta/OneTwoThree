package main

// A Provider is an Interface the Core Declares
// and something outside Fulfils. See RULES.md, Providers.

// OrderStore Keeps Orders wherever Orders Live.
// Fulfilled by GormOrders.
// Reference: https://gorm.io/docs/
type OrderStore interface {
	InsertOrderRow(o Order) (Order, error)
	SelectOrderRow(id OrderID) (Order, error)
	SelectOrderPage(page Page) ([]Order, error)
}

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

package main

// OrderBody is what a Client Sends.
// Items Carry no Identity of their own, so they Cross unchanged.
type OrderBody struct {
	MemberID string `json:"memberId"`
	Percent  int    `json:"percent"`
	Items    []Item `json:"items"`
}

// OrderView is what a Client Gets back.
// Total and Receipt are Derived, never Stored twice.
type OrderView struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Percent  int    `json:"percent"`
	Items    []Item `json:"items"`
	Total    int    `json:"total"`
	Receipt  string `json:"receipt"`
}

// BuildOrderRecord Promotes untrusted Input into a Business Order.
func (b OrderBody) BuildOrderRecord(id OrderID) Order {
	return Order{ID: id, MemberID: b.MemberID, Percent: b.Percent, Items: b.Items}
}

// RenderOrderView Demotes a Business Order back to the Wire.
func RenderOrderView(o Order) OrderView {
	total := ApplyMemberRate(SumItemPrices(o.Items), o.Percent)
	receipt := BuildOrderReceipt(string(o.ID), o.Items, o.Percent)

	return OrderView{string(o.ID), o.MemberID, o.Percent, o.Items, total, receipt}
}

// RenderOrderViews Repeats the Move for a whole Page.
func RenderOrderViews(orders []Order) []OrderView {
	views := make([]OrderView, 0, len(orders))
	for _, order := range orders {
		views = append(views, RenderOrderView(order))
	}
	return views
}

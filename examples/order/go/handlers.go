package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// OrderAPI Holds the Cast the Story Needs.
// One Provider, no Library: a Test can Hand it one Fake.
type OrderAPI struct {
	Orders OrderStore
}

// DeclareOrderRoutes is the Libretto.
func (a OrderAPI) DeclareOrderRoutes() []Route {
	return []Route{
		{"POST", "/orders", a.AddOrderRecord},
		{"GET", "/orders", a.ListOrderRecords},
		{"GET", "/orders/{id}", a.ShowOrderRecord},
	}
}

// AddOrderRecord Opens a new Order and Prices it on the Spot.
func (a OrderAPI) AddOrderRecord(req Request) Response {
	order, err := a.ReadOrderBody(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	stored, err := a.Orders.InsertOrderRow(order)

	return BuildOrderReply(http.StatusCreated, stored, err)
}

// ShowOrderRecord Tells the Story of one Purchase.
func (a OrderAPI) ShowOrderRecord(req Request) Response {
	order, err := a.Orders.SelectOrderRow(OrderID(req.Path["id"]))

	return BuildOrderReply(http.StatusOK, order, err)
}

// ListOrderRecords Tells one Page of the Ledger.
func (a OrderAPI) ListOrderRecords(req Request) Response {
	page, err := ReadPageRequest(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	orders, err := a.Orders.SelectOrderPage(page)
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderOrderViews(orders)}
}

// ReadOrderBody Decodes the Wire, Mints an Identity, then Checks the Shape.
func (a OrderAPI) ReadOrderBody(req Request) (Order, error) {
	var body OrderBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return Order{}, ErrBodyIsBroken
	}

	order := body.BuildOrderRecord(GenerateOrderID())
	if !CheckOrderRecord(order) {
		return Order{}, ErrOrderIsInvalid
	}

	return order, CheckPercentRange(order.Percent)
}

// ReadPageRequest Reads Pagination, Defaulting to the whole Set.
func ReadPageRequest(req Request) (Page, error) {
	number, _ := strconv.Atoi(req.Query["page"])
	size, _ := strconv.Atoi(req.Query["size"])

	page := Page{Number: number, Size: size}

	return page, page.CheckPageBounds()
}

// BuildOrderReply Answers with an Order, or with why there is none.
func BuildOrderReply(code int, o Order, err error) Response {
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{code, RenderOrderView(o)}
}

// BuildFailureReply Maps a Business Failure onto an HTTP Code.
// The Mapping Lives here once, never inside a Handler.
func BuildFailureReply(err error) Response {
	body := map[string]string{"error": err.Error()}

	if errors.Is(err, ErrOrderUnknown) {
		return Response{http.StatusNotFound, body}
	}

	return Response{http.StatusBadRequest, body}
}

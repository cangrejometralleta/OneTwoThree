package main

import (
	"net/http"
	"testing"
)

// FakeOrders Stands in for GORM.
// No Database, no Framework, no Port: the Providers Allow it.
type FakeOrders struct {
	orders map[OrderID]Order
}

func BuildFakeOrders() *FakeOrders {
	return &FakeOrders{orders: map[OrderID]Order{}}
}

func (f *FakeOrders) InsertOrderRow(o Order) (Order, error) {
	f.orders[o.ID] = o
	return o, nil
}

func (f *FakeOrders) SelectOrderRow(id OrderID) (Order, error) {
	order, found := f.orders[id]
	if !found {
		return Order{}, ErrOrderUnknown
	}
	return order, nil
}

func (f *FakeOrders) SelectOrderPage(Page) ([]Order, error) {
	orders := make([]Order, 0, len(f.orders))
	for _, order := range f.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

// BuildTestingOrders Hands the API one Fake and nothing else.
func BuildTestingOrders() OrderAPI {
	return OrderAPI{Orders: BuildFakeOrders()}
}

func TestAddOrderRecordOpensAPurchase(t *testing.T) {
	reply := BuildTestingOrders().AddOrderRecord(Request{
		Body: []byte(`{"memberId":"m1","percent":10,"items":[{"Name":"Widget","Price":100,"Qty":2}]}`),
	})

	if reply.Status != http.StatusCreated {
		t.Fatalf("wanted 201, got %d: %v", reply.Status, reply.Body)
	}
}

func TestAddOrderRecordRefusesAnEmptyBasket(t *testing.T) {
	reply := BuildTestingOrders().AddOrderRecord(Request{
		Body: []byte(`{"memberId":"m1","percent":10,"items":[]}`),
	})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

func TestAddOrderRecordRefusesABadPercent(t *testing.T) {
	reply := BuildTestingOrders().AddOrderRecord(Request{
		Body: []byte(`{"memberId":"m1","percent":150,"items":[{"Name":"Widget","Price":100,"Qty":1}]}`),
	})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

func TestShowOrderRecordReportsAbsence(t *testing.T) {
	reply := BuildTestingOrders().ShowOrderRecord(Request{Path: map[string]string{"id": "unknown"}})

	if reply.Status != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d", reply.Status)
	}
}

func TestListOrderRecordsRefusesNegativePage(t *testing.T) {
	reply := BuildTestingOrders().ListOrderRecords(Request{Query: map[string]string{"page": "-1"}})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

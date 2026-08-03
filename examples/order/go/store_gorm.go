package main

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// OrderRow is the Storage Shape of an Order.
// Items Live as a JSON Column: one Table, no Join,
// the Shape a Receipt Needs and nothing more.
type OrderRow struct {
	ID       string `gorm:"primaryKey;size:32"`
	MemberID string `gorm:"size:120;not null"`
	Percent  int    `gorm:"not null"`
	Items    string `gorm:"type:text;not null"`
}

// GormOrders Fulfils OrderStore with one Connection.
// It is the only Type in this Program that Knows GORM Exists.
// Reference: https://gorm.io/docs/
type GormOrders struct {
	DB *gorm.DB
}

// InsertOrderRow Writes a new Order under the Identity it already Carries.
func (g GormOrders) InsertOrderRow(o Order) (Order, error) {
	row, err := EncodeOrderRow(o)
	if err != nil {
		return Order{}, err
	}

	if err := g.DB.Create(&row).Error; err != nil {
		return Order{}, err
	}

	return o, nil
}

// SelectOrderRow Finds one Order or Says why not.
func (g GormOrders) SelectOrderRow(id OrderID) (Order, error) {
	var row OrderRow

	if err := g.DB.First(&row, "id = ?", string(id)).Error; err != nil {
		return Order{}, TranslateStoreError(err)
	}

	return DecodeOrderRow(row)
}

// SelectOrderPage Reads one Page, newest first.
func (g GormOrders) SelectOrderPage(page Page) ([]Order, error) {
	var rows []OrderRow

	query := ApplyPageWindow(g.DB.Order("id desc"), page)
	if err := query.Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err)
	}

	return DecodeOrderRows(rows)
}

// ApplyPageWindow Turns a Page into Limit and Offset.
// Reference: https://gorm.io/docs/query.html#Limit-amp-Offset
func ApplyPageWindow(query *gorm.DB, page Page) *gorm.DB {
	if page.Size <= 0 {
		return query
	}

	return query.Offset(page.Number * page.Size).Limit(page.Size)
}

// TranslateStoreError Turns a Driver Failure into a Business one.
func TranslateStoreError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOrderUnknown
	}

	return err
}

// EncodeOrderRow Turns a Business Order into Storage.
func EncodeOrderRow(o Order) (OrderRow, error) {
	items, err := json.Marshal(o.Items)
	if err != nil {
		return OrderRow{}, err
	}

	return OrderRow{ID: string(o.ID), MemberID: o.MemberID, Percent: o.Percent, Items: string(items)}, nil
}

// DecodeOrderRow Turns Storage back into Business.
func DecodeOrderRow(r OrderRow) (Order, error) {
	var items []Item
	if err := json.Unmarshal([]byte(r.Items), &items); err != nil {
		return Order{}, err
	}

	return Order{ID: OrderID(r.ID), MemberID: r.MemberID, Percent: r.Percent, Items: items}, nil
}

// DecodeOrderRows Repeats the Move for a whole Page.
func DecodeOrderRows(rows []OrderRow) ([]Order, error) {
	orders := make([]Order, 0, len(rows))
	for _, row := range rows {
		order, err := DecodeOrderRow(row)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

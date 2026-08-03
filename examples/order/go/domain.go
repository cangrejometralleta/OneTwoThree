package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// OrderID is a String so an unassigned Order Compares Empty.
// CheckOrderRecord Already Assumes this Shape.
type OrderID string

// Item is one Line of an Order.
type Item struct {
	Name  string
	Price int
	Qty   int
}

// Order is the Business Truth about one Purchase.
type Order struct {
	ID       OrderID
	MemberID string
	Percent  int
	Items    []Item
}

// The Business Fails in named Ways, never in Numbers.
var (
	ErrOrderIsInvalid   = errors.New("order Misses an Identity, a Member or its Items")
	ErrPercentIsInvalid = errors.New("percent Must sit between zero and a Hundred")
	ErrOrderUnknown     = errors.New("order not Found")
	ErrPageIsInvalid    = errors.New("page Numbers Must not be negative")
	ErrBodyIsBroken     = errors.New("body is not valid JSON")
)

// GenerateOrderID Mints a short Reference before the Row exists.
// A Purchase Needs its own Name the Moment it is Born.
func GenerateOrderID() OrderID {
	var raw [4]byte
	_, _ = rand.Read(raw[:])

	return OrderID(hex.EncodeToString(raw[:]))
}

// CheckOrderRecord Names each Condition, then Reads them together.
// The Break is not in the Expression, it is in the Vocabulary.
func CheckOrderRecord(o Order) bool {
	identified := o.ID != ""
	assigned := o.MemberID != ""
	filled := len(o.Items) > 0

	return identified && assigned && filled
}

// CheckPercentRange Refuses a Discount the Business cannot Grant.
func CheckPercentRange(percent int) error {
	if percent < 0 || percent > 100 {
		return ErrPercentIsInvalid
	}

	return nil
}

// SumItemPrices Adds every Line into a Total.
func SumItemPrices(items []Item) int {
	total := 0
	for _, it := range items {
		total += it.Price * it.Qty
	}
	return total
}

// FormatItemLine Renders one Item for the Receipt.
func FormatItemLine(it Item) string {
	return fmt.Sprintf("%-12s x%d %6d", it.Name, it.Qty, it.Price*it.Qty)
}

// ApplyMemberRate Lowers a Total by a Percentage.
func ApplyMemberRate(total, percent int) int {
	return total - total*percent/100
}

// ReportOrderState Says how it Went, at a Glance.
func ReportOrderState(id string, total int, err error) string {
	if err != nil {
		return fmt.Sprintf("❌ Order %s Failed: %v", id, err)
	}
	return fmt.Sprintf("✅ Order %s Closed at %d", id, total)
}

// BuildOrderReceipt Reads as three Sections: Total, Lines, Result.
func BuildOrderReceipt(id string, items []Item, percent int) string {
	total := SumItemPrices(items)
	total = ApplyMemberRate(total, percent)

	lines := make([]string, 0, len(items))
	for _, it := range items {
		lines = append(lines, FormatItemLine(it))
	}

	lines = append(lines, ReportOrderState(id, total, nil))
	return strings.Join(lines, "\n")
}

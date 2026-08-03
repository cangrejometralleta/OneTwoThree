package main

import (
	"strings"
	"testing"
)

func TestSumItemPricesAddsEveryLine(t *testing.T) {
	items := []Item{{Name: "Widget", Price: 100, Qty: 2}, {Name: "Gadget", Price: 50, Qty: 1}}

	if got := SumItemPrices(items); got != 250 {
		t.Fatalf("wanted 250, got %d", got)
	}
}

func TestApplyMemberRateLowersTheTotal(t *testing.T) {
	if got := ApplyMemberRate(200, 10); got != 180 {
		t.Fatalf("wanted 180, got %d", got)
	}
}

func TestCheckOrderRecordGuardsEachCondition(t *testing.T) {
	good := Order{ID: "ab12", MemberID: "m1", Items: []Item{{Name: "Widget", Price: 1, Qty: 1}}}

	cases := map[string]struct {
		order Order
		want  bool
	}{
		"valid":        {good, true},
		"no identity":  {Order{MemberID: "m1", Items: good.Items}, false},
		"no member":    {Order{ID: "ab12", Items: good.Items}, false},
		"empty basket": {Order{ID: "ab12", MemberID: "m1"}, false},
	}

	for name, test := range cases {
		if got := CheckOrderRecord(test.order); got != test.want {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}

// The Text that Ties this Domain to the Snippet RULES.md already Cites.
func TestBuildOrderReceiptRendersThreeSections(t *testing.T) {
	items := []Item{{Name: "Widget", Price: 100, Qty: 2}}

	receipt := BuildOrderReceipt("ab12", items, 10)
	lines := strings.Split(receipt, "\n")

	if len(lines) != 2 {
		t.Fatalf("wanted 2 Lines, got %d: %q", len(lines), receipt)
	}
	if lines[0] != FormatItemLine(items[0]) {
		t.Errorf("wanted the Item Line, got %q", lines[0])
	}
	if want := ReportOrderState("ab12", 180, nil); lines[1] != want {
		t.Errorf("wanted %q, got %q", want, lines[1])
	}
}

func TestCheckPercentRangeRefusesOutOfBounds(t *testing.T) {
	if CheckPercentRange(-1) == nil {
		t.Error("a negative Percent Must be Refused")
	}
	if CheckPercentRange(101) == nil {
		t.Error("a Percent over a Hundred Must be Refused")
	}
	if CheckPercentRange(50) != nil {
		t.Error("fifty Percent Must be Accepted")
	}
}

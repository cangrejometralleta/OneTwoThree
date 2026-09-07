# Code

Go, because the Rules above Read better when they Run.

```go
import (
	"fmt"
	"strings"
)

// Item is one Line of an Order.
type Item struct {
	Name  string
	Price int
	Qty   int
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
```

- Every Name Follows **Verb + Noun + context**.
- Every Function Owns one Concern and Returns it.
- The Emoji Lives in Output, never in a Name.
- Comments Follow OneTwoThreeCase too.
- BuildOrderReceipt Spends eight Lines  
  on three Beats.
- Every Name here Counts five Syllables.

Two whole Services Live in [examples/school](../examples/school),
one in Go and one in TypeScript.
Six Frameworks Serve them and Return identical Answers.

[The Before](../examples/school/BEFORE.md) Reads the Original beside them.
Every Rule there is Broken, and each Break Names the Rule it Earned.

- A Handler there Names no Driver and no Query.
- Read one out loud and it is still a Sentence.
- Two Files Hold every Vendor Import.
- A Handler there Spends Lines on Errors  
  and still Counts three Beats.
- The Code Obeyed this Rule  
  before the Rule was Written down.

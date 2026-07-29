## OneTwoThreeCase
- A Capitalized Word means it's Important  
- Entities, Actions and Statuses   
  are Capitalized because they're Important
- The First Word in a Sentence
  is always Capitalized
- A Word in lowercase
- Connectors are in lowercase


## Emoji
- An Emoji Earns its Place  
  only when it Speeds up Reading.
- Use it to Mark a State at a Glance:  
  ✅ Passed, ❌ Failed, ⚠️ Careful.
- Keep it in Output, Comments and Docs,  
  never in an Identifier  
  or a Key your Code Compares.
- One per Line at most.  
  Two Compete, three are Noise.


## Structure
- Three-line Functions  
  Are the Ideal Size Target.
- More Lines Signal a missing Abstraction Layer.
- When the Body Earns more Lines,  
  Group them into three Sections,  
  one blank Line between each.
- Three Sections Read like three Lines,  
  so the Rhythm Survives the Length.


## Naming
- Functions Follow  
  **Verb + Noun + context** rhythm
- A Name longer than Three Words  
  Suggests unclear Responsibility.
- A Variable that Travels Holds three Words,  
  joined by the Case its Language Uses.
- A Variable that Lives in three Lines  
  Holds one Word, because the Scope Says the rest.
- Three Words Fit in Memory,  
  and they Survive a Rename.

```text
sumItemPrices      JavaScript, Go unexported
SumItemPrices      Go exported
sum_item_prices    Python
Sum Item Prices    Markdown, OneTwoThreeCase
```


## Patterns
- A Variable Name Longer than Three Words  
  Suggests unclear Responsibility.
- Single responsibility,
  a function should do only one thing
- More Lines Signal  
  a Missing Abstraction Layer.


## Anti-Patterns
- More than three Responsibilities  
  Suggest a Missing Abstraction Layer.
- A Name with no Verb  
  Suggests a missing Action.
- A Unit with no clear Return  
  Breaks the Rotation.


## Code
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

// ApplyMemberDiscount Lowers a Total by a Percentage.
func ApplyMemberDiscount(total, percent int) int {
	return total - total*percent/100
}

// ReportOrderStatus Says how it Went, at a Glance.
func ReportOrderStatus(id string, total int, err error) string {
	if err != nil {
		return fmt.Sprintf("❌ Order %s Failed: %v", id, err)
	}
	return fmt.Sprintf("✅ Order %s Closed at %d", id, total)
}

// BuildOrderReceipt Reads as three Sections: Total, Lines, Result.
func BuildOrderReceipt(id string, items []Item, percent int) string {
	total := SumItemPrices(items)
	total = ApplyMemberDiscount(total, percent)

	lines := make([]string, 0, len(items))
	for _, it := range items {
		lines = append(lines, FormatItemLine(it))
	}

	lines = append(lines, ReportOrderStatus(id, total, nil))
	return strings.Join(lines, "\n")
}
```

- Every Name Follows **Verb + Noun + context**.
- Every Function Owns one Concern and Returns it.
- The Emoji Lives in Output, never in a Name.
- Comments Follow OneTwoThreeCase too.
- BuildOrderReceipt Runs eight Lines  
  and still Reads as three.

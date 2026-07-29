## OneTwoThreeCase
- A Capitalized Word means it's Important  
- Entities, Actions and Statuses   
  are Capitalized because they're Important
- The First Word in a Sentence
  is always Capitalized
- A Word in lowercase
- Connectors are in lowercase


## Rhythm
- Contrast Carries the Line.  
  A short one after a long one Lands like a Chorus.
- Uniform Text Hides what Matters.  
  Never Write every Line the same Length,  
  not in Prose and not in Code.
- Count the Syllables if it Helps.  
  Odd often Swings.  
  But a Signal is not a Law, and the Count is only a Signal.
- The Haiku Counts five, seven, five.  
  Every Line Lands odd.
- A Name that Says an Action Counts too.  
  SumItemPrices Runs five.
- A bare Noun Keeps its own Count.  
  Order is Order, whatever it Sounds like.
- Four is the Beat, three is the Phrase.  
  They Meet again every twelve,  
  so the Tension always Resolves.
- A List Stays parallel.  
  A Paragraph Varies.  
  Contrast is for Prose, never for an Index.


## Seams
- Break where the Grammar Bends.  
  A Sentence Shows its own Joints: a Conjunction, a Comma, a Preposition.
- Never Break inside a Unit that Reads as one.  
  An Article Holds its Noun.
- Symmetry is not the Cause.  
  A Joint near the Middle just Happens to Land there.
- Among the legal Joints, Choose by Meaning.  
  A short Line Emphasises.
- Go Has the same Joints:  
  && and ||, the Comma in a List, the Dot in a Chain.
- Better than Breaking a long Expression, Name its Parts.  
  A named Condition Documents while it Breaks.

```go
// CheckOrderRecord Names each Condition, then Reads them together.
// The Break is not in the Expression, it is in the Vocabulary.
func CheckOrderRecord(o Order) bool {
	identified := o.ID != ""
	assigned := o.MemberID != ""
	filled := len(o.Items) > 0

	return identified && assigned && filled
}
```

Three Conditions, three Names, one Return.
The Chain in [examples/people/store.go](examples/people/store.go)
Breaks at the Dot for the same Reason.


## Emoji
- An Emoji Earns its Place  
  only when it Speeds up Reading.
- Use it to Mark a State:  
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
  Group them into three,  
  one blank Line between each Section.
- Three Sections Read like three Lines.  
  The Rhythm Survives.


## Script
- The Program is a Story,  
  and the Handler is its Script.
- Main Casts the Players, then Steps off the Stage.
- A Handler Speaks Business only.  
  It Names no Driver, no Query, no Socket.
- A Provider Holds the Mechanism,  
  so the Script Stays a Story and nothing more.
- Read a Handler out loud.  
  If it stops Sounding like a Sentence,  
  an Abstraction is Missing.
- Every Endpoint is one small Story:  
  a Start, a Turn and an End.


## Naming
- Functions Follow  
  **Verb + Noun + context** rhythm
- A Name longer than Three Words Suggests unclear Responsibility.
- A Variable that Travels Holds three Words,  
  joined by its Language.
- A Variable that Lives in three Lines  
  Holds one Word, because the Scope Says the rest.
- Three Words Fit in Memory and Survive a Rename.

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

A whole Service Lives in [examples/people](examples/people),
with GORM, DTOs and five Endpoints.
It Compiles, and it Runs.

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
- BuildOrderReceipt Runs eight Lines  
  and still Reads as three.
- Every Name here Counts five Syllables.

The same Package Continues. Main Casts, the Handler Narrates,
the Provider Works.

```go
import (
	"fmt"
	"net/http"
)

// Order is one Purchase a Member Made.
type Order struct {
	ID       string
	MemberID string
	Items    []Item
}

// OrderStore Finds Orders wherever they Live.
type OrderStore interface {
	LoadOrderRecord(id string) (Order, error)
}

// RateSource Knows what a Member Earns.
type RateSource interface {
	LoadMemberRating(id string) (int, error)
}

// App Holds the Cast the Story Needs.
type App struct {
	Orders OrderStore
	Rates  RateSource
}

// ServeOrderTotal Tells the Story of one Order.
func (a *App) ServeOrderTotal(w http.ResponseWriter, r *http.Request) {
	order, err := a.Orders.LoadOrderRecord(r.PathValue("id"))
	if err != nil {
		http.Error(w, "❌ Order not Found", http.StatusNotFound)
		return
	}

	rate, err := a.Rates.LoadMemberRating(order.MemberID)
	if err != nil {
		rate = 0
	}

	total := ApplyMemberRate(SumItemPrices(order.Items), rate)
	fmt.Fprint(w, ReportOrderState(order.ID, total, nil))
}

// main Casts the Players, then Steps off the Stage.
func main() {
	app := &App{Orders: MemoryOrders{}, Rates: MemoryRates{}}

	http.HandleFunc("GET /orders/{id}/total", app.ServeOrderTotal)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("❌ Server Stopped:", err)
	}
}
```

```go
import "fmt"

// MemoryOrders Serves Orders from a Map.
type MemoryOrders map[string]Order

// LoadOrderRecord Finds one Order or Says why not.
func (m MemoryOrders) LoadOrderRecord(id string) (Order, error) {
	order, found := m[id]
	if !found {
		return Order{}, fmt.Errorf("order %q not Found", id)
	}
	return order, nil
}

// MemoryRates Serves Ratings from a Map.
type MemoryRates map[string]int

// LoadMemberRating Reads a Rate, Zero when Absent.
func (m MemoryRates) LoadMemberRating(id string) (int, error) {
	return m[id], nil
}
```

- ServeOrderTotal Names no Driver and no Query.
- Read it out loud and it is still a Sentence.
- Its three Sections are Load, Adjust and Answer.
- The Provider Holds the Map,  
  so the Script never Mentions one.

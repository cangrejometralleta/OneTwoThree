# What Belongs Here

- A Rule an Agent cannot Execute  
  is a Value, not a Rule.
- Rules are Verifiable,  
  Values are Interpretable.
- Send each one to the Document that Holds it.

## OneTwoThreeCase

- A Capitalized Word means it's Important  
- Entities, Actions and Statuses
  are Capitalized because they're Important
- The First Word in a Sentence
  is always Capitalized
- A Word in lowercase is a Connector or a local Name
- Connectors are in lowercase

```go
// A Capital Crosses the Package Boundary.
// A lowercase Name Stays home.
// Go Enforces what this Document only Asks.
type Roster struct {
	Students []Student
	cursor   int
}

// CountStudentRows is public, so it is Capitalized.
func (r Roster) CountStudentRows() int {
	return len(r.Students)
}

// advanceRosterCursor Stays inside, and Reads as inside.
func (r *Roster) advanceRosterCursor() {
	r.cursor++
}
```

The Convention is not ours alone.
A major Language Reached it first, and its Compiler Holds the Line.

## Reading this Repository

- A deliberate Choice Looks like an Error  
  to a Reader in a Hurry.
- Question the odd Capital before you Correct it.  
  OneTwoThreeCase is a Convention, never a Typo.
- An Agent that Normalises this Text  
  Deletes the Signal it was Given.

## Search

- Three narrow Queries Beat one wide Query.  
  Each one Returns a different Corner.
- One broad Question Returns the Average,  
  and an Average Holds no Detail.
- The same Rule that Distributes Trust  
  Distributes a Search.

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
- A Line should Take you a Bar, or a Measure.  
  Read it out loud — you'll Feel where it Lands.
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
The Chain in [store_gorm.go](examples/school/go/store_gorm.go)
Breaks at the Dot for the same Reason.

## Emoji

- An Emoji Earns its Place  
  only when it Speeds up Reading.
- Use it to Mark a State:  
  ✅ could be Passed  
  and ❌ could be Failed  
- Keep it in Output, Comments and Docs,  
  never in an Identifier  
  or a Key your Code Compares.  
- One per Line at most.  
  Two Compete, three are Noise.

## Structure

- Three-line Functions  
  Are the Ideal Size Target.
- Three Lines Means three Beats, not three Newlines.  
  A Beat is one Thought:  
  Receive, Transform and Return.
- A Language with explicit Errors Spends Newlines.  
  Count the Thoughts instead.
- One Unit Owns one Concern.  
  A Function Does one Thing.
- More Lines Signal a missing Abstraction Layer.
- When the Body Earns more Lines,  
  Group them into three,  
  one blank Line between each Section.
- Three Sections Read like three Lines.  
  The Rhythm Survives.

```go
// ListStudentRecords Spends eleven Lines on three Beats.
// Receive, Transform, Return.
// The Errors Cost Lines, they never Cost Thoughts.
func (a SchoolAPI) ListStudentRecords(req Request) Response {
	page, err := ReadPageRequest(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	students, err := a.Students.SelectStudentPage(page)
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderStudentViews(students)}
}
```

Count the Beats and you Get three.
Count the Newlines and you Get eleven.
Only one of those Numbers Means anything.

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

## Providers

- A Provider is an Interface the Core Declares  
  and something outside Fulfils.
- The Core Depends on the Shape.  
  Never on the Library behind it.
- Name the Provider after the Business Need,  
  never after the Vendor.  
  StudentStore, not GormRepository.
- One Struct May Fulfil several Providers.  
  One Provider Must never Leak its Vendor.
- Comment each Provider with the URL  
  of the Contract it Wraps.  
  A Reader Should not have to Search.
- Count the Files that Import a Vendor.  
  If the Count Grows past one, the Provider Failed.
- The Word Collides with Angular, NestJS and Terraform,  
  where a Provider is a registered Dependency.  
  Here it is a Port.

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

Two whole Services Live in [examples/school](examples/school),
one in Go and one in TypeScript.
Six Frameworks Serve them and Return identical Answers.

- A Handler there Names no Driver and no Query.
- Read one out loud and it is still a Sentence.
- Two Files Hold every Vendor Import.
- A Handler there Spends Lines on Errors  
  and still Counts three Beats.
- The Code Obeyed this Rule  
  before the Rule was Written down.

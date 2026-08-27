# OneTwoThreeCase

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

- Read the Name as you would Count it.  
  One. Two. Three. Case.
- Four Words, four Beats, one per Breath.  
  Rushed into one Word, the Cadence Dies.
- Say it slow enough to Hear the Seams,  
  and OneTwoThreeCase Teaches its own Rule.

# Seams

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
The Chain in [store_gorm.go](../examples/school/go/store/store_gorm.go)
Breaks at the Dot for the same Reason.

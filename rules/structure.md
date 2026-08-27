# Structure

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

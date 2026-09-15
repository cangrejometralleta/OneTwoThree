# Structure

- Three-line functions  
  are the ideal size Target.
- Three lines Means three beats, not three newlines.  
  A beat is one Thought:  
  receive, transform and return.
- A language with explicit errors Spends newlines.  
  Count the Thoughts instead.
- One unit Owns one concern.  
  A function Does one thing.
- More lines Signal a missing abstraction layer.
- When the body Earns more lines,  
  group them into Three,  
  one blank line between each section.
- Three sections Read like three lines.  
  The rhythm Survives.

```go
// ListStudentRecords Spends eleven lines on three beats.
// Receive, transform, return.
// The errors Cost lines, they never cost thoughts.
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

Count the beats and you get Three.  
Count the newlines and you get Eleven.  
Only one of those numbers Means anything.

# Structure

- Three-line functions  
  are the ideal size Target.
- Three Lines means three Beats, not three newlines.  
  A beat is one Thought:  
  receive, transform and return.
- A language with explicit errors Spends newlines.  
  Count the Thoughts instead.
- One Unit Owns one Concern.  
  A Function Does one Thing.
- More Lines signal a missing Abstraction Layer.
- When the Body Earns more Lines,  
  group them into Three,  
  one blank line between each section.
- Three Sections Read like three Lines.  
  The Rhythm Survives.

```go
// ShowStudentRecord Spends thirteen lines on three beats.
// Receive the path, read the service, return the wire view.
// The errors Cost lines, they never cost thoughts.
func (a Handler) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.School.ReadStudent(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return renderStudentView(student), nil
}
```

Count the beats and you get Three.

Count the newlines and you get Thirteen.

Only one of those numbers Means anything.

# Shapes

- The Entity is never the DTO.  
  One Shape Arrives Untrusted; the other holds the truth.
- Three Shapes Carry one Record,  
  and each one answers to a different layer.

```text
StudentBody   the Wire      untrusted, weak types
Student       the Business  trusted, named types
StudentRow    the Storage   trusted, table types
```

- Bind the wire Shape, never the entity.  
  A framework that fills an entity from a body  
  Hands the caller a setter for every column.
- The Identity comes from the Path or the Store,  
  never from the body.  
  A client that can send an ID can Overwrite a stranger.
- The weak Shape holds Strings where the strong shape holds Types.  
  Promotion is the Border, and it belongs to the core.
- One crossing each way, both Named.  
  Build Promotes, render demotes.
- A Shape that serves two Layers serves neither.  
  It grows the Fields of both and the guarantees of one.

```java
// The before Binds the entity, so the wire reaches the table.
public Student create(@RequestBody @Valid Student student)
```

```go
// The Handler binds a wire shape and maps it into a business shape.
func (a Handler) readStudentBody(req transport.Request, id school.StudentID) (school.Student, error) {
	body, err := app.ReadJSONBody[wire.StudentBody](req)
	if err != nil {
		return school.Student{}, err
	}

	return buildStudentRecord(body, id), nil
}
```

The Annotation Saved four Lines and spent the boundary.  
Both versions Live in [examples/school](../examples/school),  
the before read in [BEFORE.md](../examples/school/BEFORE.md).

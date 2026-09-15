# Shapes

- The entity is never the DTO.  
  One shape Arrives untrusted; the other holds the truth.
- Three shapes Carry one record,  
  and each one answers to a different layer.

```text
StudentBody   the Wire      untrusted, weak types
Student       the Business  trusted, named types
StudentRow    the Storage   trusted, table types
```

- Bind the wire Shape, never the entity.  
  A framework that fills an entity from a body  
  Hands the caller a setter for every column.
- The identity Comes from the path or the store,  
  never from the body.  
  A client that can send an ID can Overwrite a stranger.
- The weak shape Holds strings where the strong shape holds types.  
  Promotion is the Border, and it belongs to the core.
- One crossing each way, both Named.  
  Build Promotes, render demotes.
- A shape that serves two layers Serves neither.  
  It grows the Fields of both and the guarantees of one.

```java
// The before Binds the entity, so the wire reaches the table.
public Student create(@RequestBody @Valid Student student)
```

```go
// The after Binds a body, and the core decides what it becomes.
body, err := app.ReadJSONBody[wire.StudentBody](req)
student := school.BuildStudentRecord(body, 0)
```

The annotation Saved four lines and spent the boundary.  
Both versions Live in [examples/school](../examples/school),  
the before read in [BEFORE.md](../examples/school/BEFORE.md).

# Shapes

- The Entity is never the DTO.  
  One Shape Arrives untrusted; the other Holds the Truth.
- Three Shapes Carry one Record,  
  and each one Answers to a different Layer.

```text
StudentBody   the Wire      untrusted, weak Types
Student       the Business  trusted, named Types
StudentRow    the Storage   trusted, Table Types
```

- Bind the Wire Shape, never the Entity.  
  A Framework that Fills an Entity from a Body  
  Hands the Caller a Setter for every Column.
- The Identity Comes from the Path or the Store,  
  never from the Body.  
  A Client that can Send an ID can Overwrite a Stranger.
- The weak Shape Holds Strings where the strong Shape Holds Types.  
  Promotion is the Border, and it Belongs to the Core.
- One Crossing each Way, both Named.  
  Build Promotes, Render Demotes.
- A Shape that Serves two Layers Serves neither.  
  It Grows the Fields of both and the Guarantees of one.

```java
// The Before Binds the Entity, so the Wire Reaches the Table.
public Student create(@RequestBody @Valid Student student)
```

```go
// The After Binds a Body, and the Core Decides what it Becomes.
body, err := app.ReadJSONBody[wire.StudentBody](req)
student := school.BuildStudentRecord(body, 0)
```

The Annotation Saved four Lines and Spent the Boundary.
Both Versions Live in [examples/school](../examples/school),
the Before Read in [BEFORE.md](../examples/school/BEFORE.md).

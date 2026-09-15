# Values

> Provisional. Written from practice, not yet Weathered.

- A number with a meaning Carries a name.  
  The index Counts; the name explains.
- `400` says where it Fell in a list someone else wrote.  
  `HTTP_BAD_REQUEST` says what Happened.
- Look for the Name before you write one.  
  A vendor already Named it, almost always.

```java
// The literal Costs a reader one lookup, every time.
return new Fault(409, reason);

// The JDK already Named it.
return new Fault(HTTP_CONFLICT, reason);
```

- Take the Name from the narrowest source that owns it.  
  The standard library First, the framework second,  
  and your own constant only when neither knows.
- A constant in the core must not Drag a vendor in.  
  `java.net.HttpURLConnection` Costs nothing.  
  `org.springframework.http.HttpStatus` Costs the boundary.
- The exceptions are the Numbers that mean themselves.  
  Zero, one, the index in a Loop.  
  A weight table Keeps its digits and names the table.
- Two literals with one meaning are one missing Name.  
  Find the second occurrence and the name Writes itself.

The test is Mechanical.  
Read the Literal alone, out of its line.  
If it cannot say what it means, it Wants a name.

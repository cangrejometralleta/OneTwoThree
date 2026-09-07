# Values

> Provisional. Written from Practice, not yet Weathered.

- A Number with a Meaning Carries a Name.  
  The Index Counts; the Name Explains.
- `400` Says where it Fell in a List someone else Wrote.  
  `HTTP_BAD_REQUEST` Says what Happened.
- Look for the Name before you Write one.  
  A Vendor already Named it, almost always.

```java
// The Literal Costs a Reader one Lookup, every Time.
return new Fault(409, reason);

// The JDK already Named it.
return new Fault(HTTP_CONFLICT, reason);
```

- Take the Name from the narrowest Source that Owns it.  
  The Standard Library first, the Framework second,  
  and your own Constant only when neither Knows.
- A Constant in the Core Must not Drag a Vendor in.  
  `java.net.HttpURLConnection` Costs nothing.  
  `org.springframework.http.HttpStatus` Costs the Boundary.
- The Exceptions are the Numbers that Mean themselves.  
  Zero, one, the Index in a Loop.  
  A Weight Table Keeps its Digits and Names the Table.
- Two Literals with one Meaning are one Missing Name.  
  Find the second Occurrence and the Name Writes itself.

The Test is Mechanical.
Read the Literal alone, out of its Line.
If it Cannot Say what it Means, it Wants a Name.

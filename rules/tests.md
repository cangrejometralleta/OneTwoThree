# Tests

> Provisional. Written from Practice, not yet Weathered.

- Spell the Expectation; never Read it from the Code under Test.  
  A Test that Computes what it Checks  
  Agrees with itself and Proves nothing.

```go
// Both Sides Move together. The Mutation Passes.
if reply.Status != faults.ReadFaultStatus(want) {

// The Table Says the Number. The Mutation Fails.
if reply.Status != status {
```

- Prove it by Mutation.  
  Break the Declaration and Run the Test.  
  A Test that still Passes was never Watching.
- A Test Name is a Use Case, not a Method Name.  
  *someone Enrols below the Enrolment Age*  
  Beats *testAddStudentInvalidAge*.
- Arrive the Way a Caller Arrives.  
  Cross the Route, not the private Helper.
- One Fake Stands in for every Vendor.  
  A Test that Needs a Port has Found a missing Provider.
- Keep one real Test per Guarantee only a real Thing can Give.  
  A unique Index is one. A Clock is another.
- The Fake Obeys the same Rules as the Store.  
  A Fake that Accepts what the Store Refuses Hides the Case.

Count the Declarations, then Count the Tests that Reach them.
The two Numbers Match, or the Gap has a Name.

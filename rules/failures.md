# Failures

> Provisional. Written from Practice, not yet Weathered.

- A Failure the Program Expected Carries its own Answer.  
  Declare the Answer beside the Reason, once.
- One Type Holds both Halves.  
  The Reason is for the Reader; the Answer is for the Caller.
- Name the Failure by the Case, never by the Number.  
  *Rut Taken* is a Business Fact.  
  *Conflict* is only how it Ends.

```go
ErrRutTaken = faults.ReportTakenValue("rut is already Registered")
```

- A Failure Carrying no Answer was never Controlled.  
  It Answers five hundred, because it is Ours.
- Never Guess on the Caller's Behalf.  
  A Driver Error Reaching the Edge as four hundred  
  Blames the Caller for our own Break.
- One Function Turns a Failure into a Number,  
  and the whole Program Holds exactly one.
- A Constructor per Kind Reads better than a Field.  
  *RefuseInvalidInput* Names the Move; `status: 400` only Sets it.
- The Guard Answers one Way and Explains nothing.  
  A Caller who cannot Prove itself Learns only that.

Grep the Codebase for the Mapping.
One Hit is the Rule Working.
Two Hits is a Decision that Drifted.

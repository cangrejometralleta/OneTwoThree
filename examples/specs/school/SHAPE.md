# School — the Shape

How the Code is Arranged, in no Language.
[The Contract](openapi.yaml) and [The Specification](SPEC.md) Say what the Service Does;
this Page Says how little Architecture it Needs, and why.
A Rewrite Generated from these Must Keep the Behaviour and the Shape.

Every Line here Names the Rule that Earns it.
A Seam with no Rule is a Layer nobody can Delete.

## The Door

A Handler Orchestrates: it Decodes, Calls one Service method, and Answers
with a Value or Fails. It Reaches no Provider, no Store and no Vendor.
[Layers](../../../rules/layers.md) · [Structure](../../../rules/structure.md)

## The Service

The Service Joins Providers. Each use case is one Verb the Service Owns,
and its Body Reads as the Order the Spec Names:
Form, Course, Registry, Store.

A Provider is an Interface its Consumer Declares, Named for one Object
and one Action. A Port is Named for the need, never the Vendor.
[Providers](../../../rules/providers.md)

| Port | Object and Action | Fulfilled by |
| --- | --- | --- |
| `StudentStore` | Student, select · insert · update · delete | `store` |
| `CourseStore` | Course, select · insert | `store` |
| `TokenIssuer` | Token, issue · verify | `tokens` |
| `RutRegistry` | Rut, confirm | `campus` |
| `EnrollmentNotifier` | Enrollment, announce | `campus` |

One Adapter may Fill several Ports. The Composition Root Casts it once
and Hands the same Object to each Port it Fills.
The Service never Learns that two Questions have one Answerer.

## The Failure that Stays Inside

Two Providers Fail in two Ways, and the Spec Tells them Apart.

- The Registry Failing is a Controlled Fault: the Caller can Retry,
  so the Answer Carries it. Registry Unavailable, never a Guess.
- The Notifier Failing is Kept inside the Answer: the Caller cannot Act on
  a Notice that did not Leave, so the Service Logs it and Succeeds.

[Failures](../../../rules/failures.md)

## The Vendors

Every Vendor Lives in one Package, and the Core Imports none.
The Compiler Holds the Boundary: `go list -deps ./school` Names no Vendor,
no Socket and no Wire Package.
[Layers](../../../rules/layers.md) · [Vendor Integration](../../../rules/vendor-integration.md)

```text
school     the Core: business Truth, no Provider, no Vendor
app        the Service, its Ports, the Crossing and Validation
handler    the Door: Routes and Handler methods
transport  Request, Response and Route: the Shapes the Door and the Adapters Share; imports nothing
wire       the Contract: every Shape a Client Sends or Receives
faults     controlled Failures, each Carrying its Kind
store      the only Package that Imports an ORM
campus     the Registry and the Notifier, one Office Filling both Ports
tokens     the Token Adapter, standard Library only
serving    one Package per Framework; the Server and the Request Reading
settings   the strict JSON Loader, Validated at Startup
main       Casts the Players and Picks an Adapter
```

## Other Languages

A Language that is not Go Keeps the same Shape, and Answers these Questions its own Way.

- JSON is a Vendor where the Standard Library has none. It Lives in one Package, `codec`,
  the only one that Imports the Library; the Core, the Token and the Door Read Plain Values.
- The Compiler Holds the Boundary where it Can. Where one Module cannot, a Test Reads the
  Imports of each Package and Fails when the Core Names a Vendor; it Says which Vendor may
  Live where, and the Build File is Part of the Check.
- A Driver the Runtime Finds without an Import, as JDBC does, still Lives in the Store.
- A Port is Named for one Object, and its Actions are that Object's: the Token Issuer
  Mints and Reads.
- The Door Splits Reading: `serving` Reads Bytes, Method, Path and Headers, and `app`
  Reads them into Typed Values.
- A Store Port Answers Absence however its Language Does, and the Service Raises Unknown.
  The Store alone Raises RUT Taken, because only it Proves it.
- A Test Name is the Story, in the Language's own Case.
- The Format Gate is the Language's own Formatter; where none Ships, the Gate Checks what
  it Can and Says so.
- A Build Names the Toolchain it Needs, and a Missing one Stops it with the Fix.

## The Three Shapes

One Student Travels as a Wire Shape, a Business Shape and a Storage Shape.
The Entity is never the DTO.
[Shapes](../../../rules/shapes.md)

## The Answer

A Fault Carries a Kind, not a Number. One Table in `app` Turns a Kind
into a Status, so the Core stays free of `net/http`.
A Failure with no Fault Answers five hundred: it is Ours.

## The Tests

A Unit Test is a Use Case. Its Name Tells the Story the Spec Tells,
and a Spec Line with no Test that Pins it is a Wish.
[Tests](../../../rules/tests.md)

- The Core Tests its Rules, with no Port: Form, Age, RUT.
- The Service Tests each Use Case through Fakes of its Ports.
  No Database, no Network, no Port Number: the Providers Allow it.
  This is where Admission is Proven, step by step and in Order.
- The Door Tests that each Fault Reaches the Edge Whole, with the Status
  its Kind Maps to. The Number is Spelled in the Table, never Read from
  the Fault under Test, so a Declaration that Drifts Fails here.
- An Adapter Tests its own Contract: what it Refuses to Open, what it Holds.

| The Spec Says | The Test that Pins it |
| --- | --- |
| A Course needs a Code and a Name | `school`: `TestCheckCourseRecordGuardsEachRule`, `app`: `TestCreateCourseRefusesACourseWithoutCodeOrName`, `handler`: `TestEachFaultReachesTheEdgeWhole` |
| Name, age and RUT Rules | `school`: `TestCheckStudentRecordGuardsEachRule`, `TestRutLooksValid*` |
| The Order is Form, Course, Registry, Store | `app`: `TestTheRegistryIsAskedAfterTheFormAndTheCourse` |
| The Registry Denies | `app`: `TestEnrollStudentRefusesWhatTheRegistryDenies` |
| An Unknown is never a Yes | `app`: `TestEnrollStudentNeverTreatsAnUnknownAsAYes` |
| Enrolling Announces once the Store Keeps | `app`: `TestEnrollStudentAnnouncesOnceTheStoreKeeps` |
| A Failed Notice never Fails the Enrolment | `app`: `TestEnrollStudentKeepsAFailedNoticeInsideTheAnswer`, `handler`: `TestEnrolmentSurvivesAFailedAnnouncement` |
| A Rewrite Asks and Announces nothing | `app`: `TestSaveStudentAsksTheRegistryAndAnnouncesNothing` |
| Reading, Listing and Dropping go through the Ports | `app`: `TestReadStudent*`, `TestListStudents*`, `TestDeleteStudent*`, `TestCreateCourse*`, `TestReadCourse*` |
| Each of the Nine Routes Answers its Success Status | `handler`: `TestEachRouteAnswersWithItsDeclaredStatus` |
| A Request is Checked in Form: Path, Body, Page | `app`: `TestReadPathNumberWantsAPositiveWholeNumber`, `TestReadJSONBodyRefusesAnythingButOneWholeObject`, `TestReadPageRequest*` |
| Every Framework Answers the same Way | `serving`: `adaptertest.CheckAdapterAnswersTheSameWay`, run by the stdlib, chi and gin Adapters |
| The Store Proves the RUT is Unique | `store`: `TestInsertStudentRowRefusesADuplicateRut` |
| Each Fault Answers its Status | `handler`: `TestEachFaultReachesTheEdgeWhole`, `app`: `TestReadFaultStatus*` |
| Every Route but the Token Needs a Caller | `handler`: `TestEveryRouteButTheTokenRefusesAnUnnamedCaller` |
| A Token Dies once, never before | `tokens`: `TestAccessTokenDiesOnlyAfterItsDeadline` |
| Startup Refuses what it cannot Trust | `settings`: `TestConfigRejectsInvalidInputs`, `campus`: `TestOpenOfficeStopsOnWhatItCannotTrust` |

A Regeneration Passes when `build.sh` Passes: the Tests are the Spec, Run.

## Regenerate, never Patch

The Spec and this Shape are the Source. The Code is what they Produce.

1. Change the Spec or the Shape first, with the Rule that Earns the Change.
2. Delete the generated Code, and Generate it again from both.
3. Run `build.sh`: it Refuses what does not Pass.
   Each Line of the Spec has a Test that Pins it.
4. Compare with the Real Code. When Muchi API Disagrees, the Real Code Wins,
   and the Shape Changes before School does.
5. Keep [the Before](../../school/BEFORE.md) as it is. It is History, not Spec.

The Go Service is Regenerated first.
Java and TypeScript Wait for Go to Close.

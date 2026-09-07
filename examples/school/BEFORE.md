# The Before

> ⚠️ Every Snippet on this Page is the Counter-Example.
> None of it Follows the Rules, and none of it should be Copied.
> The Canon is [Go](go), [TypeScript](ts) and [Java](java);
> this is what they Rewrote.

The original Technical Test, Spring Boot, March 2020.
Same Domain as the Rest of `examples/school`:
Students, Courses, a Chilean RUT and a Token.

The Code Lived here for a while and has been Removed; it did its Work.
The whole Source Stays at [cangrejometralleta/spring-test](https://github.com/cangrejometralleta/spring-test)
at `fb38252`, Unedited. An Edit would Cost the only thing it Offers,
which is Being real.

What Remains is the Reading, because the Reading is the Part that Teaches.

## What was already Right

The Rules did not Arrive from outside. Three of them were already Here.

**A Failure Carries its own Answer.**
The Status Sits on the Exception Class, so no Controller ever Chooses one.

```java
@ResponseStatus(code = HttpStatus.NOT_FOUND)
public class StudentNotFoundException extends Exception {

  public StudentNotFoundException(Long studentId) {
    super("The student with id: '" + studentId + "' was not found");
  }
}
```

[faults](go/faults/faults.go) is this Idea, Rewritten.
Four Exception Classes, one per Business Case: a Failure Named, never Numbered.

**A Controller Delegates in one Line.**
Every Method Hands the Work to a Service and Returns.
[Script](../../rules/script.md) Names what this Code already Did.

## What the Rules came to Fix

Each Seam Earned a Rule. The Rule is the Link.

### The Entity was the DTO

```java
@PostMapping("/")
public Student create(@RequestBody @Valid Student student)
```

The JPA Entity Binds straight to the Wire, so a Client can Send an `id`.
The Annotation Saved four Lines and Spent the Boundary — [Shapes](../../rules/shapes.md).

### A Package that Named no Action

`util/` Held JWT Signing, RUT Arithmetic and Remote Address Reading.
Three Concerns behind a Word with no Verb — [Naming](../../rules/naming.md).
They became `tokens/`, `school/` and `app/`.

### A Name Counted by Position

```java
public static String verificationDigit(String rut) {
  int M = 0, S = 1, T = Integer.parseInt(rut);
  for (; T != 0; T = (int) Math.floor(T /= 10))
    S = (S + T % 10 * (9 - M++ % 6)) % 11;
  return (S > 0) ? String.valueOf(S - 1) : "k";
}
```

`M`, `S` and `T` Count; they never Explain — [Naming](../../rules/naming.md)
and [Values](../../rules/values.md). The Rewrite Named the Parts:
`CheckDigitFor`, `RutWeights`, `NameCheckRemainder`.

### A Comment that Said the Signature back

```java
/**
 * 
 * @return
 */
public Iterable<Student> findAll() {
```

An empty `@param` and an empty `@return` Survive their own Deletion —
[Comments](../../rules/comments.md).

### A Secret with a Fallback, Committed

`application.properties` Carried a Default Password, and `JwtUtil` Carried
a Default Signing Key. A Deployment Value and a Secret Shared one File —
[Constants](../../rules/constants.md).

### A Validation Overwritten two Lines later

```java
Course course = getCourse(student.getCourse().getId());
updateStudent.setCourse(course);

updateStudent.setAge(student.getAge());
updateStudent.setCourse(student.getCourse());
```

The validated Course goes in, then the unvalidated one from the Request
Body Replaces it. The Check on `update` Bought nothing —
[Structure](../../rules/structure.md).

## The Bug the Rewrite Pinned

`JwtUtil.tokenIsValid` Reads:

```java
return remoteAddress.equals(subject) && Instant.now().compareTo(expiration) > 0;
```

`now > expiration` Means *already Expired*.
The Token Validates only after it Dies, and never before.

A Comparison Reading backwards Survives every Review that only Reads it.
All three Rewrites Pin it with a Test that Moves a Clock:
[Go](go/tokens/token_hmac_test.go),
[TypeScript](ts/src/school.test.ts),
[Java](java/src/test/java/com/example/tokens/AccessTokensTest.java).

That is the Argument for the Before-and-After.
The Verse that did not Rhyme Teaches more than the one that did.

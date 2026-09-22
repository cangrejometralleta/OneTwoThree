# The Before

> ⚠️ Every Snippet on this page is the Counter-Example.
> None of it Follows the rules, and none of it should be Copied.
> The Canon is [Go](go), [TypeScript](ts) and [Java](java);
> this is what they Rewrote.

The original technical Test, Spring Boot, March 2020.
Same domain as the rest of `examples/school`:
Students, Courses, a Chilean RUT and a Token.

The Code Lived here for a while and has been Removed; it did its work.
The whole Source Stays at [cangrejometralleta/spring-test](https://github.com/cangrejometralleta/spring-test)
at `fb38252`, Unedited. An edit would Cost the only thing it Offers,
which is being real.

What Remains is the reading, because the reading is the part that Teaches.

## What was already Right

The Rules did not Arrive from outside. Three of them were already Here.

**A failure carries its own answer.**
The Status Sits on the exception Class, so no controller ever Chooses one.

```java
@ResponseStatus(code = HttpStatus.NOT_FOUND)
public class StudentNotFoundException extends Exception {

  public StudentNotFoundException(Long studentId) {
    super("The student with id: '" + studentId + "' was not found");
  }
}
```

[faults](go/faults/faults.go) is this Idea, Rewritten.
Four exception classes, one per business case: a failure Named, never Numbered.

**A controller delegates in one line.**
Every Method Hands the work to a service and Returns.
[Script](../../rules/script.md) Names what this Code already Did.

## What the Rules came to Fix

Each Seam Earned a Rule. The Rule is the Link.

### The Entity was the DTO

```java
@PostMapping("/")
public Student create(@RequestBody @Valid Student student)
```

The JPA Entity Binds straight to the Wire, so a client can send an `id`.
The Annotation saved four Lines and spent the Boundary — [Shapes](../../rules/shapes.md).

### A Package that Named no Action

`util/` Held JWT signing, RUT arithmetic and remote address reading.
Three Concerns behind a word with no Verb — [Naming](../../rules/naming.md).
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
and [Values](../../rules/values.md). The rewrite Named the Parts:
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

`application.properties` Carried a default password, and `JwtUtil` Carried
a default signing key. A deployment Value and a Secret shared one File —
[Constants](../../rules/constants.md).

### A Validation Overwritten two Lines later

```java
Course course = getCourse(student.getCourse().getId());
updateStudent.setCourse(course);

updateStudent.setAge(student.getAge());
updateStudent.setCourse(student.getCourse());
```

The validated course goes in, then the unvalidated one from the request
body Replaces it. The Check on `update` Bought nothing —
[Structure](../../rules/structure.md).

## The Bug the Rewrite Pinned

`JwtUtil.tokenIsValid` Reads:

```java
return remoteAddress.equals(subject) && Instant.now().compareTo(expiration) > 0;
```

`now > expiration` Means *already Expired*.
The Token Validates only after it Dies, and never before.

A comparison reading backwards Survives every Review that only reads it.
All three rewrites Pin it with a Test that moves a clock:
[Go](go/tokens/token_hmac_test.go),
[TypeScript](ts/src/school.test.ts),
[Java](java/src/test/java/com/example/tokens/AccessTokensTest.java).

That is the argument for the Before-and-After.
The Verse that did not Rhyme Teaches more than the one that did.

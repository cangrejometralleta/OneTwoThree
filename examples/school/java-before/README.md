# The Before

> ⚠️ This Directory is the Counter-Example.
> Nothing here Follows the Rules, and nothing here should be Copied.
> The Canon is [Go](../go) and [TypeScript](../ts); this is what they Rewrote.


The original Technical Test, Spring Boot, March 2020.
Same Domain as the Rest of `examples/school`:
Students, Courses, a Chilean RUT and a Token.

Copied Verbatim from [cangrejometralleta/spring-test](https://github.com/cangrejometralleta/spring-test)
at `fb38252`. Not Edited, not Improved, not Compiled here.
An Edit would Cost the only thing this Code Offers, which is Being real.

## What was already Right

The Rules did not Arrive from outside. Three of them were already Here.

- **A Failure Carries its own Answer.**
  `@ResponseStatus(code = HttpStatus.NOT_FOUND)` Sits on the Exception Class,
  so no Controller ever Chooses a Number.
  [faults](../go/faults/faults.go) is this Idea, Rewritten.
- **A Controller Delegates in one Line.**
  Every Method Hands the Work to a Service and Returns.
  [Script](../../../rules/script.md) Names what this Code already Did.
- **A Failure is Named, never Numbered.**
  Four Exception Classes, one per Business Case.

## What the Rules came to Fix

Each Seam below Earned a Rule. The Rule is the Link.

| The Seam | The Rule |
| --- | --- |
| `@RequestBody @Valid Student` Binds the JPA Entity straight to the Wire | [Shapes](../../../rules/shapes.md) |
| `util/` Holds JWT, RUT and IP, and Names no Action | [Naming](../../../rules/naming.md) |
| `verificationDigit` Counts with `M`, `S` and `T` | [Naming](../../../rules/naming.md) |
| Javadoc with an empty `@param` and an empty `@return` | [Comments](../../../rules/comments.md) |
| A Secret with a Fallback Value, Committed | [Constants](../../../rules/constants.md) |
| `StudentService.update` Sets the Course twice, one of them Unvalidated | [Structure](../../../rules/structure.md) |

## The Bug the Rewrite Pinned

`JwtUtil.tokenIsValid` Reads:

```java
return remoteAddress.equals(subject) && Instant.now().compareTo(expiration) > 0;
```

`now > expiration` Means *already Expired*.
The Token Validates only after it Dies, and never before.

A Comparison Reading backwards Survives every Review that only Reads it.
Both Rewrites Pin it with a Test that Moves a Clock:
[Go](../go/tokens/token_hmac_test.go), [TypeScript](../ts/src/school.test.ts).

That is the Argument for the Before-and-After.
The Verse that did not Rhyme Teaches more than the one that did.

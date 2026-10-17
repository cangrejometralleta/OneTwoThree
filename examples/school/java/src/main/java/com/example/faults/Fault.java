package com.example.faults;

import java.util.Objects;

/**
 * A Fault is a Failure the Program Expected, Carrying the Kind it Deserves.
 * A Failure that Reaches the Edge without one is Unexpected.
 */
public final class Fault extends RuntimeException {

  /** INTERNAL_REASON is all a Caller Learns of a Failure that is Ours. */
  public static final String INTERNAL_REASON = "internal error";

  private final Kind kind;

  private Fault(Kind kind, String reason) {
    super(reason, null, false, false);
    this.kind = kind;
  }

  public Kind kind() {
    return kind;
  }

  public String reason() {
    return getMessage();
  }

  /** refuseInvalidInput Names a Caller that Sent something the Rules Reject. */
  public static Fault refuseInvalidInput(String reason) {
    return new Fault(Kind.INVALID_INPUT, reason);
  }

  /** refuseUnprovenCaller Names a Caller the Program cannot Recognise. */
  public static Fault refuseUnprovenCaller(String reason) {
    return new Fault(Kind.UNPROVEN_CALLER, reason);
  }

  /** reportMissingRecord Names something the Caller Asked for and We Lack. */
  public static Fault reportMissingRecord(String reason) {
    return new Fault(Kind.MISSING_RECORD, reason);
  }

  /** reportTakenValue Names a Collision with something already Stored. */
  public static Fault reportTakenValue(String reason) {
    return new Fault(Kind.TAKEN_VALUE, reason);
  }

  /** refuseUnprocessableValue Names a well Formed Value the Authority that Knows still Denies. */
  public static Fault refuseUnprocessableValue(String reason) {
    return new Fault(Kind.UNPROCESSABLE_VALUE, reason);
  }

  /** reportUnavailableProvider Names a Provider that cannot Answer now; an Unknown is never a Yes. */
  public static Fault reportUnavailableProvider(String reason) {
    return new Fault(Kind.UNAVAILABLE_PROVIDER, reason);
  }

  /** readFaultKind Finds the Fault inside any Wrapping; everything else is Unexpected. */
  public static Kind readFaultKind(Throwable failure) {
    return findFault(failure).map(Fault::kind).orElse(Kind.UNEXPECTED);
  }

  /** readFaultReason Finds the Reason a Fault Declares, and Says nothing of a driver. */
  public static String readFaultReason(Throwable failure) {
    return findFault(failure).map(Fault::reason).orElse(INTERNAL_REASON);
  }

  private static java.util.Optional<Fault> findFault(Throwable failure) {
    for (Throwable cause = failure; cause != null; cause = cause.getCause()) {
      if (cause instanceof Fault fault) {
        return java.util.Optional.of(fault);
      }
    }

    return java.util.Optional.empty();
  }

  @Override
  public boolean equals(Object other) {
    return other instanceof Fault fault && fault.kind == kind && fault.reason().equals(reason());
  }

  @Override
  public int hashCode() {
    return Objects.hash(kind, reason());
  }
}

package com.example.faults;

import static java.net.HttpURLConnection.HTTP_BAD_REQUEST;
import static java.net.HttpURLConnection.HTTP_CONFLICT;
import static java.net.HttpURLConnection.HTTP_INTERNAL_ERROR;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_UNAUTHORIZED;

/**
 * A Fault is a Failure the Program Expected, Carrying the Answer it Deserves.
 * A Failure that Reaches the Edge without one is a five hundred.
 *
 * <p>The Status Comes from the JDK, never from a Literal.
 * A Number Counts; a Name Explains.
 */
public final class Fault extends RuntimeException {

  private final int status;

  private Fault(int status, String reason) {
    super(reason, null, false, false);
    this.status = status;
  }

  public int status() {
    return status;
  }

  /** RefuseInvalidInput Names a Caller that Sent something the Rules Reject. */
  public static Fault refuseInvalidInput(String reason) {
    return new Fault(HTTP_BAD_REQUEST, reason);
  }

  /** RefuseUnprovenCaller Names a Caller the Program cannot Recognise. */
  public static Fault refuseUnprovenCaller(String reason) {
    return new Fault(HTTP_UNAUTHORIZED, reason);
  }

  /** ReportMissingRecord Names something the Caller Asked for and We Lack. */
  public static Fault reportMissingRecord(String reason) {
    return new Fault(HTTP_NOT_FOUND, reason);
  }

  /** ReportTakenValue Names a Collision with something already Stored. */
  public static Fault reportTakenValue(String reason) {
    return new Fault(HTTP_CONFLICT, reason);
  }

  /**
   * ReadFaultStatus is the one Place that Turns a Failure into a Number.
   * A Fault Answers for itself; everything else Answers five hundred.
   */
  public static int readFaultStatus(Throwable failure) {
    for (Throwable step = failure; step != null; step = step.getCause()) {
      if (step instanceof Fault fault) {
        return fault.status();
      }
    }

    return HTTP_INTERNAL_ERROR;
  }
}

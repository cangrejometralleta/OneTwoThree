package com.example.app;

import static java.net.HttpURLConnection.HTTP_BAD_REQUEST;
import static java.net.HttpURLConnection.HTTP_CONFLICT;
import static java.net.HttpURLConnection.HTTP_INTERNAL_ERROR;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_UNAUTHORIZED;
import static java.net.HttpURLConnection.HTTP_UNAVAILABLE;

import com.example.faults.Fault;
import com.example.faults.Kind;

import java.util.Map;

/**
 * Statuses is the one Place that Turns a Failure into a Number.
 * The Core Names the Kind; only the Application Layer Knows the Protocol.
 * Reference: https://www.rfc-editor.org/rfc/rfc9110#section-15
 */
public final class Statuses {

  /** HTTP_UNPROCESSABLE_CONTENT is the one Status the JDK never Named. */
  public static final int HTTP_UNPROCESSABLE_CONTENT = 422;

  private static final Map<Kind, Integer> FAULT_STATUSES = Map.of(
      Kind.INVALID_INPUT, HTTP_BAD_REQUEST,
      Kind.UNPROVEN_CALLER, HTTP_UNAUTHORIZED,
      Kind.MISSING_RECORD, HTTP_NOT_FOUND,
      Kind.TAKEN_VALUE, HTTP_CONFLICT,
      Kind.UNPROCESSABLE_VALUE, HTTP_UNPROCESSABLE_CONTENT,
      Kind.UNAVAILABLE_PROVIDER, HTTP_UNAVAILABLE);

  private Statuses() {
  }

  /** readFaultStatus Answers five hundred for a Kind without a Status. */
  public static int readFaultStatus(Throwable failure) {
    return FAULT_STATUSES.getOrDefault(Fault.readFaultKind(failure), HTTP_INTERNAL_ERROR);
  }
}

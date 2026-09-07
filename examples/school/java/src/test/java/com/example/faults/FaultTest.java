package com.example.faults;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import static java.net.HttpURLConnection.HTTP_BAD_REQUEST;
import static java.net.HttpURLConnection.HTTP_CONFLICT;
import static java.net.HttpURLConnection.HTTP_INTERNAL_ERROR;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_UNAUTHORIZED;

import org.junit.jupiter.api.Test;

class FaultTest {

  /** Every controlled Fault Answers with the Status it Declared. */
  @Test
  void readFaultStatusAnswersForEachKind() {
    assertEquals(HTTP_BAD_REQUEST, Fault.readFaultStatus(Fault.refuseInvalidInput("bad")));
    assertEquals(HTTP_UNAUTHORIZED, Fault.readFaultStatus(Fault.refuseUnprovenCaller("who")));
    assertEquals(HTTP_NOT_FOUND, Fault.readFaultStatus(Fault.reportMissingRecord("gone")));
    assertEquals(HTTP_CONFLICT, Fault.readFaultStatus(Fault.reportTakenValue("taken")));
  }

  /** An uncontrolled Failure is ours. */
  @Test
  void readFaultStatusRefusesToGuess() {
    assertEquals(HTTP_INTERNAL_ERROR, Fault.readFaultStatus(new IllegalStateException("the driver Broke")));
  }

  /** A wrapped Fault still Carries its Answer. */
  @Test
  void readFaultStatusSurvivesWrapping() {
    Fault gone = Fault.reportMissingRecord("gone");
    Exception wrapped = new IllegalStateException("while Reading", gone);

    assertEquals(HTTP_NOT_FOUND, Fault.readFaultStatus(wrapped));
    assertTrue(gone.getMessage().equals("gone"));
  }
}

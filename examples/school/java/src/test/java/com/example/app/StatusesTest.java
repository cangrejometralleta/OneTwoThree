package com.example.app;

import static org.junit.jupiter.api.Assertions.assertEquals;

import com.example.faults.Fault;

import java.util.Map;

import org.junit.jupiter.api.Test;

class StatusesTest {

  @Test
  void eachControlledFaultAnswersTheStatusItsKindMapsTo() {
    Map<Fault, Integer> cases = Map.of(
        Fault.refuseInvalidInput("bad"), 400,
        Fault.refuseUnprovenCaller("who"), 401,
        Fault.reportMissingRecord("gone"), 404,
        Fault.reportTakenValue("taken"), 409,
        Fault.refuseUnprocessableValue("no"), 422,
        Fault.reportUnavailableProvider("down"), 503);

    cases.forEach((fault, status) -> assertEquals(status, Statuses.readFaultStatus(fault), fault.reason()));
  }

  @Test
  void anUncontrolledFailureAnswersFiveHundred() {
    assertEquals(500, Statuses.readFaultStatus(new IllegalStateException("the driver Broke")));
  }

  @Test
  void aWrappedFaultStillCarriesItsAnswer() {
    assertEquals(404, Statuses.readFaultStatus(new RuntimeException("while Reading", Fault.reportMissingRecord("gone"))));
  }
}

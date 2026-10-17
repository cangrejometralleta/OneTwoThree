package com.example.faults;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.Map;

import org.junit.jupiter.api.Test;

class FaultTest {

  @Test
  void eachFaultAnswersWithTheKindItDeclared() {
    Map<Fault, Kind> cases = Map.of(
        Fault.refuseInvalidInput("bad"), Kind.INVALID_INPUT,
        Fault.refuseUnprovenCaller("who"), Kind.UNPROVEN_CALLER,
        Fault.reportMissingRecord("gone"), Kind.MISSING_RECORD,
        Fault.reportTakenValue("taken"), Kind.TAKEN_VALUE,
        Fault.refuseUnprocessableValue("no"), Kind.UNPROCESSABLE_VALUE,
        Fault.reportUnavailableProvider("down"), Kind.UNAVAILABLE_PROVIDER);

    cases.forEach((fault, kind) -> assertEquals(kind, Fault.readFaultKind(fault), fault.reason()));
  }

  @Test
  void anUncontrolledFailureIsOursAndItsMessageNeverLeaves() {
    RuntimeException driver = new IllegalStateException("the driver Broke at /var/db");

    assertEquals(Kind.UNEXPECTED, Fault.readFaultKind(driver));
    assertEquals("internal error", Fault.readFaultReason(driver));
  }

  @Test
  void aWrappedFaultStillCarriesItsKindAndReason() {
    RuntimeException wrapped = new RuntimeException("while Reading", Fault.reportMissingRecord("gone"));

    assertEquals(Kind.MISSING_RECORD, Fault.readFaultKind(wrapped));
    assertEquals("gone", Fault.readFaultReason(wrapped));
  }
}

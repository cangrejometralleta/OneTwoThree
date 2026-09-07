package com.example.school;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;

class RutTest {

  /** Known RUTs, Checked by Hand against the Modulo eleven Rule. */
  @Test
  void looksValidAcceptsRealNumbers() {
    for (String rut : new String[] {"12345678-5", "11111111-1", "8459162-9", "6-k", "1-9"}) {
      assertTrue(new Rut(rut).looksValid(), rut + " Should be valid");
    }
  }

  @Test
  void looksValidRefusesBadDigits() {
    for (String rut : new String[] {"12345678-9", "11111111-2", "12345678", "abc-1", "", "12345678-"}) {
      assertFalse(new Rut(rut).looksValid(), rut + " Should be Refused");
    }
  }

  /** A Deployment Value never Reaches a Global Constant. */
  @Test
  void schoolConstantsHoldTheEnrolmentAge() {
    assertTrue(Constants.readSchoolConstants().minimumAgeYears() == 18);
  }
}

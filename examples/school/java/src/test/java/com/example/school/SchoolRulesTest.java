package com.example.school;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.example.faults.Fault;

import java.util.List;

import org.junit.jupiter.api.Test;

class SchoolRulesTest {

  private static final SchoolRules RULES = new SchoolRules(18);

  private static Student studentOf(String rut, String name, long age) {
    return new Student(StudentId.UNNUMBERED, new Rut(rut), new FullName(name), age, new CourseId(1));
  }

  private static Course courseOf(String code, String name) {
    return new Course(CourseId.UNNUMBERED, new CourseCode(code), new FullName(name));
  }

  @Test
  void rutLooksValidAcceptsRealNumbers() {
    for (String rut : List.of("12345678-5", "11111111-1", "22222222-2", "10000013-k", "10000013-K", "10000004-0",
        "123456789012345678901234567890-1")) {
      assertTrue(new Rut(rut).looksValid(), rut);
    }
  }

  @Test
  void rutLooksValidRefusesBadDigits() {
    for (String rut : List.of("12345678-9", "12345678", "abc-1", "-5", "12345678-55", "", " 12345678-5")) {
      assertFalse(new Rut(rut).looksValid(), rut);
    }
  }

  @Test
  void aRutEntersInLowerCaseAndKeepsItsLeadingZeros() {
    assertEquals(new Rut("10000013-k"), new Rut("10000013-K"));
    assertEquals("10000013-k", new Rut("10000013-K").text());
    assertFalse(new Rut("012345678-5").equals(new Rut("12345678-5")));
  }

  @Test
  void checkStudentRecordGuardsEachRuleInOrder() {
    assertDoesNotThrow(() -> studentOf("12345678-5", "Ada", 20).checkStudentRecord(RULES));
    assertDoesNotThrow(() -> studentOf("12345678-5", "Ada", 18).checkStudentRecord(RULES));

    assertEquals(Errors.NAME_IS_EMPTY,
        assertThrows(Fault.class, () -> studentOf("bad", "", 1).checkStudentRecord(RULES)));
    assertEquals(Errors.NAME_IS_EMPTY,
        assertThrows(Fault.class, () -> studentOf("12345678-5", " \t", 20).checkStudentRecord(RULES)));
    assertEquals(Errors.RUT_IS_INVALID,
        assertThrows(Fault.class, () -> studentOf("12345678-9", "Ada", 1).checkStudentRecord(RULES)));
    assertEquals(Errors.ageIsTooLow(RULES),
        assertThrows(Fault.class, () -> studentOf("12345678-5", "Ada", 17).checkStudentRecord(RULES)));
  }

  @Test
  void aNameIsStoredAsItArrivedUntrimmed() {
    assertEquals("  Ada ", new FullName("  Ada ").text());
  }

  @Test
  void theAgeReasonSpellsTheFloorFromTheConstant() {
    assertEquals("age Must be 18 or more", Errors.ageIsTooLow(RULES).reason());
    assertEquals("age Must be 21 or more", Errors.ageIsTooLow(new SchoolRules(21)).reason());
  }

  @Test
  void checkCourseRecordGuardsEachRule() {
    assertDoesNotThrow(() -> courseOf("MAT-101", "Algebra").checkCourseRecord());
    assertEquals(Errors.CODE_IS_EMPTY, assertThrows(Fault.class, () -> courseOf("", "Algebra").checkCourseRecord()));
    assertEquals(Errors.NAME_IS_EMPTY, assertThrows(Fault.class, () -> courseOf("MAT-101", " ").checkCourseRecord()));
    assertEquals(Errors.CODE_IS_EMPTY, assertThrows(Fault.class, () -> courseOf("", "").checkCourseRecord()));
  }

  @Test
  void aPageRefusesWhatTheStoreCannotServe() {
    assertDoesNotThrow(() -> new Page(0, 1));
    assertEquals(Errors.PAGE_IS_INVALID, assertThrows(Fault.class, () -> new Page(-1, 10)));
    assertEquals(Errors.PAGE_IS_INVALID, assertThrows(Fault.class, () -> new Page(0, 0)));
  }

  @Test
  void aPageBeyondSixtyFourBitsWindowsAsFarAsItCan() {
    assertEquals(20, new Page(2, 10).offset());
    assertEquals(Long.MAX_VALUE, new Page(Long.MAX_VALUE, 10).offset());
  }
}

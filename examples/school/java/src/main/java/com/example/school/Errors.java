package com.example.school;

import com.example.faults.Fault;

/**
 * The Business Fails in named Ways, never in Numbers.
 * Each one Declares the Answer it Deserves, once and here.
 */
public final class Errors {

  private Errors() {
  }

  public static final Fault NAME_IS_EMPTY = Fault.refuseInvalidInput("name is Empty");
  public static final Fault RUT_IS_INVALID = Fault.refuseInvalidInput("rut Fails its Check Digit");
  public static final Fault AGE_IS_TOO_LOW =
      Fault.refuseInvalidInput("age Must be " + Constants.readSchoolConstants().minimumAgeYears() + " or more");
  public static final Fault PAGE_IS_INVALID = Fault.refuseInvalidInput("page Numbers Must not be negative");
  public static final Fault STUDENT_UNKNOWN = Fault.reportMissingRecord("student not Found");
  public static final Fault COURSE_UNKNOWN = Fault.reportMissingRecord("course not Found");
  public static final Fault RUT_TAKEN = Fault.reportTakenValue("rut is already Registered");
}

package com.example.school;

import com.example.faults.Fault;

/**
 * The Business Fails in named Ways, never in Numbers.
 * Each one Declares the Answer it Deserves, once and here.
 */
public final class Errors {

  public static final Fault NAME_IS_EMPTY = Fault.refuseInvalidInput("name is Empty");
  public static final Fault CODE_IS_EMPTY = Fault.refuseInvalidInput("code is Empty");
  public static final Fault RUT_IS_INVALID = Fault.refuseInvalidInput("rut Fails its Check Digit");
  public static final Fault PAGE_IS_INVALID = Fault.refuseInvalidInput("page or size is not valid");
  public static final Fault STUDENT_UNKNOWN = Fault.reportMissingRecord("student not Found");
  public static final Fault COURSE_UNKNOWN = Fault.reportMissingRecord("course not Found");
  public static final Fault RUT_TAKEN = Fault.reportTakenValue("rut is already Registered");
  public static final Fault RUT_UNREGISTERED = Fault.refuseUnprocessableValue("rut is not in the Registry");
  public static final Fault REGISTRY_UNAVAILABLE = Fault.reportUnavailableProvider("registry cannot Answer now");

  private Errors() {
  }

  /** ageIsTooLow Spells the Floor from the Global Constant. */
  public static Fault ageIsTooLow(SchoolRules rules) {
    return Fault.refuseInvalidInput("age Must be " + rules.minimumAgeYears() + " or more");
  }
}

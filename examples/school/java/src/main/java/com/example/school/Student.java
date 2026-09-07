package com.example.school;


/** Student is the Business Truth about one Enrolment. */
public record Student(StudentId id, Rut rut, FullName name, int age, CourseId course) {

  /** CheckStudentRecord Refuses a Student the Business cannot Use. */
  public void checkStudentRecord() {
    if (name == null || name.value() == null || name.value().isBlank()) {
      throw Errors.NAME_IS_EMPTY;
    }

    if (rut == null || rut.value() == null || !rut.looksValid()) {
      throw Errors.RUT_IS_INVALID;
    }

    checkStudentAge();
  }

  /** CheckStudentAge Holds the one Rule the School will not Bend. */
  public void checkStudentAge() {
    if (age < Constants.readSchoolConstants().minimumAgeYears()) {
      throw Errors.AGE_IS_TOO_LOW;
    }
  }
}

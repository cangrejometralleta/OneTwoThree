package com.example.school;

/** Student is the Business Truth about one Enrolment. */
public record Student(StudentId id, Rut rut, FullName name, long age, CourseId course) {

  /** withId Numbers the Student the Store Kept, or the Path Named. */
  public Student withId(StudentId numbered) {
    return new Student(numbered, rut, name, age, course);
  }

  /** checkStudentRecord Refuses a Student the Business cannot Use: Name, then Rut, then Age. */
  public void checkStudentRecord(SchoolRules rules) {
    if (name.isEmpty()) {
      throw Errors.NAME_IS_EMPTY;
    }
    if (!rut.looksValid()) {
      throw Errors.RUT_IS_INVALID;
    }
    if (age < rules.minimumAgeYears()) {
      throw Errors.ageIsTooLow(rules);
    }
  }
}

package com.example.school;

/** Course is the Business Truth about one Class. */
public record Course(CourseId id, CourseCode code, FullName name) {

  /** withId Numbers the Course the Store Kept. */
  public Course withId(CourseId numbered) {
    return new Course(numbered, code, name);
  }

  /** checkCourseRecord Refuses a Course the Business cannot Use: Code, then Name. */
  public void checkCourseRecord() {
    if (code.isEmpty()) {
      throw Errors.CODE_IS_EMPTY;
    }
    if (name.isEmpty()) {
      throw Errors.NAME_IS_EMPTY;
    }
  }
}

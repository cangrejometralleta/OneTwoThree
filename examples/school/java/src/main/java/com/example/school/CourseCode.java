package com.example.school;

/** CourseCode is the short Label a Course Answers to. Codes are not Unique. */
public record CourseCode(String text) {

  /** isEmpty Holds for a Code of only whitespace. */
  public boolean isEmpty() {
    return Text.holdsOnlyWhitespace(text);
  }
}

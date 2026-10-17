package com.example.school;

/** CourseId is the Identity the Store Gives a Course. */
public record CourseId(long value) {

  /** UNNUMBERED Stands for a Course the Store has not Numbered yet. */
  public static final CourseId UNNUMBERED = new CourseId(0);
}

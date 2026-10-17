package com.example.school;

/** StudentId is the Identity the Store Gives a Student; the Body never Carries it. */
public record StudentId(long value) {

  /** UNNUMBERED Stands for a Student the Store has not Numbered yet. */
  public static final StudentId UNNUMBERED = new StudentId(0);
}

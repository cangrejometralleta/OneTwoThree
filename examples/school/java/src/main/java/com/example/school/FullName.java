package com.example.school;

/** FullName is what a Person or a Course is Called, Stored as it Arrived. */
public record FullName(String text) {

  /** isEmpty Holds for a Name of only whitespace; the Text itself is never Trimmed. */
  public boolean isEmpty() {
    return Text.holdsOnlyWhitespace(text);
  }
}

package com.example.school;

/** Text Holds the one Reading of whitespace the Business Shares. */
final class Text {

  private Text() {
  }

  static boolean holdsOnlyWhitespace(String text) {
    return text.codePoints().allMatch(point -> Character.isWhitespace(point) || Character.isSpaceChar(point));
  }
}

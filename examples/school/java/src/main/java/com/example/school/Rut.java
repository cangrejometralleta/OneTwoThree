package com.example.school;

import java.util.Locale;
import java.util.regex.Pattern;

/**
 * Rut is the Chilean tax Number a Student is Known by.
 * It Enters in lower case, so 12345-K and 12345-k are one Rut; leading zeros are not Merged.
 * Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
 */
public record Rut(String text) {

  /** RUT_WEIGHTS Cycle two through seven, Read right to left. */
  static final int[] RUT_WEIGHTS = {2, 3, 4, 5, 6, 7};

  /** CHECK_MODULUS is the Modulo the Check Digit Encodes. */
  static final int CHECK_MODULUS = 11;

  private static final Pattern RUT_SHAPE = Pattern.compile("^[0-9]+-[0-9kK]$");

  public Rut {
    text = text.toLowerCase(Locale.ROOT);
  }

  /** looksValid Runs the Modulo eleven Check the Digit Encodes. */
  public boolean looksValid() {
    if (!RUT_SHAPE.matcher(text).matches()) {
      return false;
    }

    int dash = text.indexOf('-');

    return text.substring(dash + 1).equals(checkDigitFor(text.substring(0, dash)));
  }

  /** checkDigitFor Folds a Rut Body into its single Check Character, Digit by Digit. */
  static String checkDigitFor(String body) {
    int sum = 0;
    for (int position = 0; position < body.length(); position++) {
      int digit = body.charAt(body.length() - 1 - position) - '0';
      sum = (sum + digit * RUT_WEIGHTS[position % RUT_WEIGHTS.length]) % CHECK_MODULUS;
    }

    return nameCheckRemainder(CHECK_MODULUS - sum);
  }

  /** nameCheckRemainder Turns a Remainder into the Character it Means: eleven is 0, ten is k. */
  static String nameCheckRemainder(int remainder) {
    return switch (remainder) {
      case 11 -> "0";
      case 10 -> "k";
      default -> Integer.toString(remainder);
    };
  }
}

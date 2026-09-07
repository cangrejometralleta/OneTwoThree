package com.example.school;

/**
 * RUT is the Chilean Tax Identity, Digit included.
 * A String Underneath, and still Impossible to Pass as a Name.
 */
public record Rut(String value) {

  /** RutWeights Cycles two through seven, Read right to left. */
  private static final int[] WEIGHTS = {2, 3, 4, 5, 6, 7};

  /**
   * LooksValid Runs the Modulo eleven Check the Digit Encodes.
   * Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
   */
  public boolean looksValid() {
    String digit = digitOf(value);

    return !digit.isEmpty() && checkDigitFor(value).equals(digit);
  }

  /** CheckDigitFor Folds a RUT Body into its single Check Character. */
  public static String checkDigitFor(String rut) {
    if (rut == null || !rut.matches("^[0-9]+-[0-9kK]$")) {
      return "";
    }

    int number = Integer.parseInt(rut.substring(0, rut.lastIndexOf('-')));
    int sum = 0;
    for (int index = 0; number > 0; number /= 10, index++) {
      sum += number % 10 * WEIGHTS[index % WEIGHTS.length];
    }

    return nameCheckRemainder(11 - sum % 11);
  }

  /** NameCheckRemainder Turns a Remainder into the Character it Means. */
  private static String nameCheckRemainder(int remainder) {
    return switch (remainder) {
      case 11 -> "0";
      case 10 -> "k";
      default -> String.valueOf(remainder);
    };
  }

  private static String digitOf(String rut) {
    int dash = rut == null ? -1 : rut.lastIndexOf('-');

    return dash < 0 ? "" : rut.substring(dash + 1).toLowerCase();
  }
}

package com.example.settings;

import java.math.BigInteger;
import java.util.Set;

/** ValueRule Validates one declared Configuration Key, or Says what it Wanted. */
@FunctionalInterface
public interface ValueRule {
  void check(Object value) throws SettingsRefused;

  /** checkIntegerRange Keeps JSON Numbers whole and Bounded. */
  static ValueRule checkIntegerRange(long minimum, long maximum) {
    return value -> {
      if (!(value instanceof BigInteger whole)) {
        throw new SettingsRefused("Must be an Integer");
      }
      if (whole.compareTo(BigInteger.valueOf(minimum)) < 0 || whole.compareTo(BigInteger.valueOf(maximum)) > 0) {
        throw new SettingsRefused("Must be between " + minimum + " and " + maximum);
      }
    };
  }

  /** checkTextValue Requires a nonempty String. */
  static ValueRule checkTextValue() {
    return value -> {
      if (!(value instanceof String text)) {
        throw new SettingsRefused("Must be a String");
      }
      if (text.isBlank()) {
        throw new SettingsRefused("Must not be Empty");
      }
    };
  }

  /** checkOneOf Accepts only the Names it is Given. */
  static ValueRule checkOneOf(Set<String> names) {
    return value -> {
      if (!(value instanceof String text) || !names.contains(text)) {
        throw new SettingsRefused("Must be one of " + names.stream().sorted().toList());
      }
    };
  }
}
